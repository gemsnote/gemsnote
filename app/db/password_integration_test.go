package db

import (
	"os"
	"testing"
	"time"

	"github.com/gemsnote/gemsnote/app/info"
	"github.com/gemsnote/gemsnote/app/lea"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoPasswordStorageContract(t *testing.T) {
	url := os.Getenv("GEMSNOTE_INTEGRATION_MONGO_URL")
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_MONGO_URL to run MongoDB integration tests")
	}

	previousURL, previousDB := initURL, initDBName
	initURL, initDBName = url, ""
	defer func() { initURL, initDBName = previousURL, previousDB }()

	mongo := &MongoDatabase{}
	if err := mongo.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer mongo.Close()
	assertPasswordStorageContract(t, mongo, mongo.Users, mongo.Sessions)
}

func TestPostgresPasswordStorageContract(t *testing.T) {
	url := os.Getenv("GEMSNOTE_INTEGRATION_POSTGRES_URL")
	if url == "" {
		t.Skip("set GEMSNOTE_INTEGRATION_POSTGRES_URL to run PostgreSQL integration tests")
	}

	previousURL, previousDB := initURL, initDBName
	initURL, initDBName = url, "gemsnote"
	defer func() { initURL, initDBName = previousURL, previousDB }()

	postgres := &PostgresDatabase{}
	if err := postgres.Initialize(); err != nil {
		t.Fatal(err)
	}
	defer postgres.Close()
	assertPasswordStorageContract(t, postgres, "users", "sessions")
}

func assertPasswordStorageContract(t *testing.T, database Database, users, sessions interface{}) {
	t.Helper()

	userID := bson.NewObjectId()
	sessionID := bson.NewObjectId()
	oldHash := lea.GenPwd("old-password")
	newHash := lea.GenPwd("new-password")
	if oldHash == "" || newHash == "" {
		t.Fatal("generate password hash failed")
	}

	user := info.User{
		UserId:      userID,
		Email:       userID.Hex() + "@password.test",
		Username:    userID.Hex(),
		UsernameRaw: userID.Hex(),
		Pwd:         oldHash,
		CreatedTime: time.Now().UTC(),
	}
	if !database.Insert(users, user) {
		t.Fatal("insert password test user failed")
	}
	defer database.DeleteAll(users, bson.M{"_id": userID})

	session := info.Session{
		Id:          sessionID,
		SessionId:   sessionID.Hex(),
		UserId:      userID.Hex(),
		CreatedTime: time.Now().UTC(),
		UpdatedTime: time.Now().UTC(),
	}
	if !database.Insert(sessions, session) {
		t.Fatal("insert password test session failed")
	}
	defer database.DeleteAll(sessions, bson.M{"_id": sessionID})

	if database.UpdateByQField(users, bson.M{"_id": bson.NewObjectId()}, "Pwd", newHash) {
		t.Fatal("password update reported success for a missing user")
	}
	if !database.UpdateByQField(users, bson.M{"_id": userID}, "Pwd", newHash) {
		t.Fatal("password update failed")
	}

	var updated info.User
	database.Get(users, userID.Hex(), &updated)
	if !lea.ComparePwd("new-password", updated.Pwd) || lea.ComparePwd("old-password", updated.Pwd) {
		t.Fatal("updated password cannot be verified correctly")
	}

	if !database.DeleteAll(sessions, bson.M{"UserId": userID.Hex()}) {
		t.Fatal("session revocation failed")
	}
	var revoked info.Session
	database.GetByQ(sessions, bson.M{"SessionId": sessionID.Hex()}, &revoked)
	if revoked.Id != "" {
		t.Fatal("password update left an active session")
	}
}
