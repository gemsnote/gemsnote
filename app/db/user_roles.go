package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

// Role changes, account status changes and login session binding share the
// same transaction lock. A disabled/revoked administrator cannot race another
// administrator to remove all available administrators or issue a new token.
const userAdministrationLock = "SELECT pg_advisory_xact_lock(1196641875, 1)"

var roleChangeMu sync.Mutex

func SetUserRole(actorID, targetID, role string) (bool, string) {
	if msg := validateRoleChange(actorID, targetID, role); msg != "" {
		return false, msg
	}
	return setUserAdministrationField(actorID, targetID, "role", role)
}
func SetUserFlag(actorID, targetID, flag string) (bool, string) {
	if msg := validateFlagChange(actorID, targetID, flag); msg != "" {
		return false, msg
	}
	return setUserAdministrationField(actorID, targetID, "flag", flag)
}
func validateRoleChange(actorID, targetID, role string) string {
	if !bson.IsObjectIdHex(actorID) || !bson.IsObjectIdHex(targetID) {
		return "userNotExists"
	}
	if role != info.RoleAdmin && role != info.RoleUser {
		return "invalidRole"
	}
	if strings.EqualFold(actorID, targetID) && role == info.RoleUser {
		return "cannotDemoteSelf"
	}
	return ""
}
func validateFlagChange(actorID, targetID, flag string) string {
	if !bson.IsObjectIdHex(actorID) || !bson.IsObjectIdHex(targetID) {
		return "userNotExists"
	}
	if flag != info.FlagAvailable && flag != info.FlagDisabled {
		return "invalidUserFlag"
	}
	if strings.EqualFold(actorID, targetID) && flag == info.FlagDisabled {
		return "cannotDisableSelf"
	}
	return ""
}
func UserAvailable(userID string) bool {
	if !bson.IsObjectIdHex(userID) {
		return false
	}
	user := info.User{}
	GetByQWithFields(Users, bson.M{"_id": bson.ObjectIdHex(userID)}, []string{"_id", "Flag"}, &user)
	return user.Available()
}
func setUserAdministrationField(actorID, targetID, column, value string) (bool, string) {
	// Only these internal column names may enter SQL below.
	if column != "role" && column != "flag" {
		return false, "invalidRequest"
	}
	actorID = bson.ObjectIdHex(actorID).Hex()
	targetID = bson.ObjectIdHex(targetID).Hex()
	if postgres, ok := globalDatabase.(*PostgresDatabase); ok {
		return postgres.setUserAdministrationField(actorID, targetID, column, value)
	}
	roleChangeMu.Lock()
	defer roleChangeMu.Unlock()
	actor := info.User{}
	Get(Users, actorID, &actor)
	if actor.Role != info.RoleAdmin || !actor.Available() {
		return false, "noAuth"
	}
	target := info.User{}
	Get(Users, targetID, &target)
	if target.UserId == "" {
		return false, "userNotExists"
	}
	field := "Role"
	if column == "flag" {
		field = "Flag"
	}

	if column == "flag" {
		if !DeleteAll(Sessions, bson.M{"UserId": targetID}) || !DeleteAll(Tokens, bson.M{"_id": target.UserId}) {
			return false, "sessionRevocationFailed"
		}
	}
	if !UpdateByQField(Users, bson.M{"_id": target.UserId}, field, value) {
		return false, "updateFailed"
	}

	return true, ""
}
func (p *PostgresDatabase) setUserAdministrationField(actorID, targetID, column, value string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, "updateFailed"
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, userAdministrationLock); err != nil {
		return false, "updateFailed"
	}
	actor := info.User{UserId: bson.ObjectIdHex(actorID)}
	err = tx.QueryRowContext(ctx, "SELECT role,flag FROM users WHERE id=$1", actorID).Scan(&actor.Role, &actor.Flag)
	if err == sql.ErrNoRows || err == nil && (actor.Role != info.RoleAdmin || !actor.Available()) {
		return false, "noAuth"
	}
	if err != nil {
		return false, "updateFailed"
	}
	result, err := tx.ExecContext(ctx, fmt.Sprintf("UPDATE users SET %s=$1 WHERE id=$2", column), value, targetID)
	if err != nil {
		return false, "updateFailed"
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, "updateFailed"
	}
	if count != 1 {
		return false, "userNotExists"
	}
	if column == "flag" {
		if _, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE user_id=$1", targetID); err != nil {
			return false, "sessionRevocationFailed"
		}
		if _, err = tx.ExecContext(ctx, "DELETE FROM tokens WHERE id=$1", targetID); err != nil {
			return false, "sessionRevocationFailed"
		}
	}
	if err = tx.Commit(); err != nil {
		return false, "updateFailed"
	}
	return true, ""
}

// BindUserSession rechecks availability atomically with the session write.
// An authentication request begun before disabling must not recreate a token
// after the administrator's session deletion has committed.
func BindUserSession(sessionID, userID string) bool {
	if sessionID == "" || !bson.IsObjectIdHex(userID) {
		return false
	}
	userID = bson.ObjectIdHex(userID).Hex()
	if postgres, ok := globalDatabase.(*PostgresDatabase); ok {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tx, err := postgres.db.BeginTx(ctx, nil)
		if err != nil {
			return false
		}
		defer tx.Rollback()
		if _, err = tx.ExecContext(ctx, userAdministrationLock); err != nil {
			return false
		}
		user := info.User{UserId: bson.ObjectIdHex(userID)}
		if err = tx.QueryRowContext(ctx, "SELECT flag FROM users WHERE id=$1", userID).Scan(&user.Flag); err != nil || !user.Available() {
			return false
		}
		result, err := tx.ExecContext(ctx, "UPDATE sessions SET user_id=$1,updated_time=CURRENT_TIMESTAMP WHERE session_id=$2", userID, sessionID)
		if err != nil {
			return false
		}
		count, err := result.RowsAffected()
		return err == nil && count > 0 && tx.Commit() == nil
	}
	roleChangeMu.Lock()
	defer roleChangeMu.Unlock()
	return UserAvailable(userID) && UpdateByQMap(Sessions, bson.M{"SessionId": sessionID}, bson.M{"UserId": userID, "UpdatedTime": time.Now()})
}
