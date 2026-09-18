package db

import (
	"fmt"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var globalDatabase Database
var dbType string

var Session *mgo.Session

var Notebooks, Notes, NoteContents interface{}
var NoteContentHistories, ShareNotes, ShareNotebooks interface{}
var HasShareNotes, Blogs, Users, Groups interface{}
var GroupUsers, Tags, NoteTags, TagCounts interface{}
var UserBlogs, Tokens, Suggestions interface{}
var Albums, Files, Attachs, NoteImages interface{}
var Configs, EmailLogs interface{}
var BlogLikes, BlogComments, Reports interface{}
var BlogSingles, Themes, Sessions interface{}

func Init(url, dbname string) {
	initURL, initDBName = url, dbname
	dbType = getDBType()
	var err error
	switch dbType {
	case "mongodb":
		globalDatabase, err = newMongoDB()
	case "postgresql":
		globalDatabase, err = newPostgresDB()
	default:
		panic("Unknown database type: " + dbType)
	}
	if err != nil {
		panic(err)
	}
	if err = migrate(globalDatabase); err != nil {
		globalDatabase.Close()
		globalDatabase = nil
		panic(err)
	}
	globalDatabase.SetupLegacyVariables()
}

func Close() {
	if globalDatabase != nil {
		globalDatabase.Close()
	}
}

func CheckConnection() {
	if globalDatabase != nil {
		globalDatabase.CheckConnection()
	}
}

func NewID() string {
	if globalDatabase != nil {
		return globalDatabase.NewID()
	}
	return ""
}

// NewUUID is retained for source compatibility. It returns a MongoDB ObjectId
// hex string for both database backends.
func NewUUID() string {
	return NewID()
}

func Insert(collection interface{}, i interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.Insert(collection, i)
}

func Update(collection interface{}, query interface{}, i interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.Update(collection, query, i)
}

func Upsert(collection interface{}, query interface{}, i interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.Upsert(collection, query, i)
}

func UpdateAll(collection interface{}, query interface{}, i interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.UpdateAll(collection, query, i)
}

func UpdateByIdAndUserId(collection interface{}, id, userId string, i interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.UpdateByIdAndUserId(collection, id, userId, i)
}

func UpdateByIdAndUserId2(collection interface{}, id, userId bson.ObjectId, i interface{}) bool {
	return UpdateByIdAndUserId(collection, id.Hex(), userId.Hex(), i)
}

func UpdateByIdAndUserIdField(collection interface{}, id, userId, field string, value interface{}) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}

func UpdateByIdAndUserIdMap(collection interface{}, id, userId string, v bson.M) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection interface{}, id, userId bson.ObjectId, field string, value interface{}) bool {
	return UpdateByIdAndUserIdField(collection, id.Hex(), userId.Hex(), field, value)
}

func UpdateByIdAndUserIdMap2(collection interface{}, id, userId bson.ObjectId, v bson.M) bool {
	return UpdateByIdAndUserIdMap(collection, id.Hex(), userId.Hex(), v)
}

func UpdateByQField(collection interface{}, q interface{}, field string, value interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.UpdateByQField(collection, q, field, value)
}

func UpdateByQI(collection interface{}, q interface{}, v interface{}) bool {
	return UpdateByQMap(collection, q, v)
}

func UpdateByQMap(collection interface{}, q interface{}, v interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.UpdateByQMap(collection, q, v)
}

func Delete(collection interface{}, q interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.Delete(collection, q)
}

func DeleteByIdAndUserId(collection interface{}, id, userId string) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.DeleteByIdAndUserId(collection, id, userId)
}

func DeleteByIdAndUserId2(collection interface{}, id, userId bson.ObjectId) bool {
	return DeleteByIdAndUserId(collection, id.Hex(), userId.Hex())
}

func DeleteAllByIdAndUserId(collection interface{}, id, userId string) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.DeleteAllByIdAndUserId(collection, id, userId)
}

func DeleteAllByIdAndUserId2(collection interface{}, id, userId bson.ObjectId) bool {
	return DeleteAllByIdAndUserId(collection, id.Hex(), userId.Hex())
}

func DeleteAll(collection interface{}, q interface{}) bool {
	if globalDatabase == nil {
		return false
	}
	return globalDatabase.DeleteAll(collection, q)
}

func Get(collection interface{}, id string, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.Get(collection, id, i)
}

func Get2(collection interface{}, id bson.ObjectId, i interface{}) {
	Get(collection, id.Hex(), i)
}

func GetByQ(collection interface{}, q interface{}, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.GetByQ(collection, q, i)
}

func ListByQ(collection interface{}, q interface{}, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.ListByQ(collection, q, i)
}

func ListByQLimit(collection interface{}, q interface{}, i interface{}, limit int) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.ListByQLimit(collection, q, i, limit)
}

func ListByQOptions(collection interface{}, q interface{}, i interface{}, options QueryOptions) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.ListByQOptions(collection, q, i, options)
}

func GetByQOptions(collection interface{}, q interface{}, i interface{}, options QueryOptions) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.GetByQOptions(collection, q, i, options)
}

func GetByQWithFields(collection interface{}, q bson.M, fields []string, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.GetByQWithFields(collection, q, fields, i)
}

func ListByQWithFields(collection interface{}, q bson.M, fields []string, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.ListByQWithFields(collection, q, fields, i)
}

func GetByIdAndUserId(collection interface{}, id, userId string, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.GetByIdAndUserId(collection, id, userId, i)
}

func GetByIdAndUserId2(collection interface{}, id, userId bson.ObjectId, i interface{}) {
	GetByIdAndUserId(collection, id.Hex(), userId.Hex(), i)
}

func Distinct(collection interface{}, q bson.M, field string, i interface{}) {
	if globalDatabase == nil {
		return
	}
	globalDatabase.Distinct(collection, q, field, i)
}

func DropIndex(collection interface{}, fields ...string) error {
	if globalDatabase == nil {
		return nil
	}
	return globalDatabase.DropIndex(collection, fields...)
}

func Count(collection interface{}, q interface{}) int {
	if globalDatabase == nil {
		return 0
	}
	return globalDatabase.Count(collection, q)
}

func Has(collection interface{}, q interface{}) bool {
	return Count(collection, q) > 0
}

func GetIdAndUserIdQ(id, userId string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}

func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
	return bson.M{"_id": id, "UserId": userId}
}

func Err(err error) bool {
	if err != nil {
		fmt.Println(err)
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

func getDBType() string {
	if revel.Config != nil {
		if t, ok := revel.Config.String("db.type"); ok {
			return t
		}
	}
	return "mongodb"
}
