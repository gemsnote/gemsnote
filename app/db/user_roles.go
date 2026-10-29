package db

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

// The SQL transaction also serializes competing administrators across server
// instances. Authority is rechecked after the lock, not only in the controller.
func SetUserRole(actorID, targetID, role string) (bool, string) {
	if msg := validateRoleChange(actorID, targetID, role); msg != "" {
		return false, msg
	}
	actorID = bson.ObjectIdHex(actorID).Hex()
	targetID = bson.ObjectIdHex(targetID).Hex()
	if postgres, ok := globalDatabase.(*PostgresDatabase); ok {
		return postgres.setUserRole(actorID, targetID, role)
	}
	roleChangeMu.Lock()
	defer roleChangeMu.Unlock()
	actor := info.User{}
	Get(Users, actorID, &actor)
	if actor.Role != info.RoleAdmin {
		return false, "noAuth"
	}
	target := info.User{}
	Get(Users, targetID, &target)
	if target.UserId == "" {
		return false, "userNotExists"
	}
	if !UpdateByQField(Users, bson.M{"_id": bson.ObjectIdHex(targetID)}, "Role", role) {
		return false, "updateFailed"
	}
	return true, ""
}

var roleChangeMu sync.Mutex

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
func (p *PostgresDatabase) setUserRole(actorID, targetID, role string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, "updateFailed"
	}
	defer tx.Rollback()
	// Fixed namespace lock reserved for user-role changes, released on commit.
	if _, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(1196641875, 1)"); err != nil {
		return false, "updateFailed"
	}
	var actorRole string
	err = tx.QueryRowContext(ctx, "SELECT role FROM users WHERE id=$1", actorID).Scan(&actorRole)
	if err == sql.ErrNoRows || err == nil && actorRole != info.RoleAdmin {
		return false, "noAuth"
	}
	if err != nil {
		return false, "updateFailed"
	}
	result, err := tx.ExecContext(ctx, "UPDATE users SET role=$1 WHERE id=$2", role, targetID)
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
	if err = tx.Commit(); err != nil {
		return false, "updateFailed"
	}
	return true, ""
}
