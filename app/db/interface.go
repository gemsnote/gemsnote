package db

import (
	"gopkg.in/mgo.v2/bson"
)

type Database interface {
	Initialize() error
	Close() error
	Ping() error
	CheckConnection()
	NewID() string
	IsValidID(id string) bool
	SetupLegacyVariables()

	Insert(collection interface{}, data interface{}) bool
	Update(collection interface{}, query interface{}, data interface{}) bool
	UpdateAll(collection interface{}, query interface{}, data interface{}) bool
	Upsert(collection interface{}, query interface{}, data interface{}) bool
	Delete(collection interface{}, query interface{}) bool
	DeleteAll(collection interface{}, query interface{}) bool

	Get(collection interface{}, id string, result interface{})
	GetByQ(collection interface{}, query interface{}, result interface{})
	ListByQ(collection interface{}, query interface{}, result interface{})
	ListByQLimit(collection interface{}, query interface{}, result interface{}, limit int)
	ListByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions)
	GetByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions)
	GetByQWithFields(collection interface{}, query bson.M, fields []string, result interface{})
	ListByQWithFields(collection interface{}, query bson.M, fields []string, result interface{})

	GetByIdAndUserId(collection interface{}, id, userId string, result interface{})
	UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool
	DeleteByIdAndUserId(collection interface{}, id, userId string) bool
	DeleteAllByIdAndUserId(collection interface{}, id, userId string) bool
	UpdateByQField(collection interface{}, q interface{}, field string, value interface{}) bool
	UpdateByQMap(collection interface{}, q interface{}, v interface{}) bool

	Count(collection interface{}, query interface{}) int
	Distinct(collection interface{}, q bson.M, field string, result interface{})
	DropIndex(collection interface{}, fields ...string) error
	AppliedMigrations() ([]string, error)
	RecordMigration(version string) error

	GetType() string
}

// QueryOptions is the common subset of query options used by Gemsnote services.
// Sort fields use mgo's convention: "Field" for ascending and "-Field" for
// descending order.
type QueryOptions struct {
	Sort   []string
	Skip   int
	Limit  int
	Fields []string
}
