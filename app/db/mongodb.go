package db

import (
	"fmt"
	"strings"
	"time"

	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDatabase struct {
	session *mgo.Session
	db      *mgo.Database

	Notebooks            *mgo.Collection
	Notes                *mgo.Collection
	NoteContents         *mgo.Collection
	NoteContentHistories *mgo.Collection
	ShareNotes           *mgo.Collection
	ShareNotebooks       *mgo.Collection
	HasShareNotes        *mgo.Collection
	Blogs                *mgo.Collection
	Users                *mgo.Collection
	Groups               *mgo.Collection
	GroupUsers           *mgo.Collection
	Tags                 *mgo.Collection
	NoteTags             *mgo.Collection
	TagCounts            *mgo.Collection
	UserBlogs            *mgo.Collection
	Tokens               *mgo.Collection
	Suggestions          *mgo.Collection
	Albums               *mgo.Collection
	Files                *mgo.Collection
	Attachs              *mgo.Collection
	NoteImages           *mgo.Collection
	Configs              *mgo.Collection
	EmailLogs            *mgo.Collection
	BlogLikes            *mgo.Collection
	BlogComments         *mgo.Collection
	Reports              *mgo.Collection
	BlogSingles          *mgo.Collection
	Themes               *mgo.Collection
	Sessions             *mgo.Collection
}

var initURL, initDBName string

func newMongoDB() (Database, error) {
	m := &MongoDatabase{}
	err := m.Initialize()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MongoDatabase) Initialize() error {
	config := revel.Config

	var url string
	var dbname string
	var ok bool

	url, dbname = initURL, initDBName
	if url != "" {
		ok = true
	} else if config != nil {
		url, ok = config.String("db.url")
	}
	if !ok && config != nil {
		url, ok = config.String("db.urlEnv")
		if ok {
			Log("get MongoDB config from db.urlEnv")
		}
	} else {
		Log("get MongoDB config from db.url")
	}
	if ok {
		urls := strings.Split(url, "/")
		dbname = urls[len(urls)-1]
		if strings.Contains(dbname, "?") {
			urls = strings.Split(dbname, "?")
			dbname = urls[0]
		}
	}

	if dbname == "" && config != nil {
		dbname, _ = config.String("db.dbname")
	}

	if !ok && config != nil {
		host, _ := config.String("db.host")
		port, _ := config.String("db.port")
		username, _ := config.String("db.username")
		password, _ := config.String("db.password")
		usernameAndPassword := ""
		if username != "" && password != "" {
			usernameAndPassword = username + ":" + password + "@"
		}
		url = "mongodb://" + usernameAndPassword + host + ":" + port + "/" + dbname
	}
	if url == "" {
		return fmt.Errorf("MongoDB connection URL is not configured")
	}

	var err error
	m.session, err = mgo.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	m.session.SetMode(mgo.Monotonic, true)
	m.db = m.session.DB(dbname)

	m.initCollections()

	Session = m.session

	Log("Connected to MongoDB successfully")
	return nil
}

func (m *MongoDatabase) initCollections() {
	m.Notebooks = m.db.C("notebooks")
	m.Notes = m.db.C("notes")
	m.NoteContents = m.db.C("note_contents")
	m.NoteContentHistories = m.db.C("note_content_histories")
	m.ShareNotes = m.db.C("share_notes")
	m.ShareNotebooks = m.db.C("share_notebooks")
	m.HasShareNotes = m.db.C("has_share_notes")
	m.Blogs = m.db.C("blogs")
	m.Users = m.db.C("users")
	m.Groups = m.db.C("groups")
	m.GroupUsers = m.db.C("group_users")
	m.Tags = m.db.C("tags")
	m.NoteTags = m.db.C("note_tags")
	m.TagCounts = m.db.C("tag_count")
	m.UserBlogs = m.db.C("user_blogs")
	m.Tokens = m.db.C("tokens")
	m.Suggestions = m.db.C("suggestions")
	m.Albums = m.db.C("albums")
	m.Files = m.db.C("files")
	m.Attachs = m.db.C("attachs")
	m.NoteImages = m.db.C("note_images")
	m.Configs = m.db.C("configs")
	m.EmailLogs = m.db.C("email_logs")
	m.BlogLikes = m.db.C("blog_likes")
	m.BlogComments = m.db.C("blog_comments")
	m.Reports = m.db.C("reports")
	m.BlogSingles = m.db.C("blog_singles")
	m.Themes = m.db.C("themes")
	m.Sessions = m.db.C("sessions")
}

func (m *MongoDatabase) SetupLegacyVariables() {
	Notebooks = m.Notebooks
	Notes = m.Notes
	NoteContents = m.NoteContents
	NoteContentHistories = m.NoteContentHistories
	ShareNotes = m.ShareNotes
	ShareNotebooks = m.ShareNotebooks
	HasShareNotes = m.HasShareNotes
	Blogs = m.Blogs
	Users = m.Users
	Groups = m.Groups
	GroupUsers = m.GroupUsers
	Tags = m.Tags
	NoteTags = m.NoteTags
	TagCounts = m.TagCounts
	UserBlogs = m.UserBlogs
	Tokens = m.Tokens
	Suggestions = m.Suggestions
	Albums = m.Albums
	Files = m.Files
	Attachs = m.Attachs
	NoteImages = m.NoteImages
	Configs = m.Configs
	EmailLogs = m.EmailLogs
	BlogLikes = m.BlogLikes
	BlogComments = m.BlogComments
	Reports = m.Reports
	BlogSingles = m.BlogSingles
	Themes = m.Themes
	Sessions = m.Sessions
}

func (m *MongoDatabase) Close() error {
	if m.session != nil {
		m.session.Close()
	}
	return nil
}

func (m *MongoDatabase) Ping() error {
	return m.session.Ping()
}

func (m *MongoDatabase) CheckConnection() {
	err := m.session.Ping()
	if err != nil {
		Log("Lost connection to db!")
		m.session.Refresh()
		err = m.session.Ping()
		if err == nil {
			Log("Reconnect to db successful.")
		} else {
			Log("重连失败!!!! 警告")
		}
	}
}

func (m *MongoDatabase) NewID() string {
	return bson.NewObjectId().Hex()
}

func (m *MongoDatabase) IsValidID(id string) bool {
	return bson.IsObjectIdHex(id)
}

func (m *MongoDatabase) GetType() string {
	return "mongodb"
}

type mongoMigration struct {
	Version   string    `bson:"_id"`
	AppliedAt time.Time `bson:"applied_at"`
}

func (m *MongoDatabase) AppliedMigrations() ([]string, error) {
	var records []mongoMigration
	if err := m.db.C("gemsnote_schema_migrations").Find(nil).All(&records); err != nil {
		return nil, err
	}
	versions := make([]string, 0, len(records))
	for _, record := range records {
		versions = append(versions, record.Version)
	}
	return versions, nil
}

func (m *MongoDatabase) RecordMigration(version string) error {
	_, err := m.db.C("gemsnote_schema_migrations").UpsertId(version, mongoMigration{
		Version: version, AppliedAt: time.Now().UTC(),
	})
	return err
}

func (m *MongoDatabase) getColl(collection interface{}) *mgo.Collection {
	if coll, ok := collection.(*mgo.Collection); ok {
		return coll
	}
	return nil
}

func (m *MongoDatabase) Insert(collection interface{}, data interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	err := coll.Insert(data)
	return Err(err)
}

func (m *MongoDatabase) Update(collection interface{}, query interface{}, data interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	err := coll.Update(query, data)
	return Err(err)
}

func (m *MongoDatabase) UpdateAll(collection interface{}, query interface{}, data interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	_, err := coll.UpdateAll(query, data)
	return Err(err)
}

func (m *MongoDatabase) Upsert(collection interface{}, query interface{}, data interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	_, err := coll.Upsert(query, data)
	return Err(err)
}

func (m *MongoDatabase) Delete(collection interface{}, query interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	err := coll.Remove(query)
	return Err(err)
}

func (m *MongoDatabase) DeleteAll(collection interface{}, query interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	_, err := coll.RemoveAll(query)
	return Err(err)
}

func (m *MongoDatabase) Get(collection interface{}, id string, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.FindId(bson.ObjectIdHex(id)).One(result)
}

func (m *MongoDatabase) GetByQ(collection interface{}, query interface{}, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.Find(query).One(result)
}

func (m *MongoDatabase) ListByQ(collection interface{}, query interface{}, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.Find(query).All(result)
}

func (m *MongoDatabase) ListByQLimit(collection interface{}, query interface{}, result interface{}, limit int) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.Find(query).Limit(limit).All(result)
}

func (m *MongoDatabase) ListByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	q := coll.Find(query)
	if len(options.Sort) > 0 {
		q = q.Sort(options.Sort...)
	}
	if options.Skip > 0 {
		q = q.Skip(options.Skip)
	}
	if options.Limit > 0 {
		q = q.Limit(options.Limit)
	}
	if len(options.Fields) > 0 {
		selector := bson.M{}
		for _, field := range options.Fields {
			selector[field] = true
		}
		q = q.Select(selector)
	}
	q.All(result)
}

func (m *MongoDatabase) GetByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	q := coll.Find(query)
	if len(options.Sort) > 0 {
		q = q.Sort(options.Sort...)
	}
	if options.Skip > 0 {
		q = q.Skip(options.Skip)
	}
	if len(options.Fields) > 0 {
		selector := bson.M{}
		for _, field := range options.Fields {
			selector[field] = true
		}
		q = q.Select(selector)
	}
	q.One(result)
}

func (m *MongoDatabase) GetByQWithFields(collection interface{}, query bson.M, fields []string, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	coll.Find(query).Select(selector).One(result)
}

func (m *MongoDatabase) ListByQWithFields(collection interface{}, query bson.M, fields []string, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	coll.Find(query).Select(selector).All(result)
}

func (m *MongoDatabase) GetByIdAndUserId(collection interface{}, id, userId string, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.Find(GetIdAndUserIdQ(id, userId)).One(result)
}

func (m *MongoDatabase) UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	err := coll.Update(GetIdAndUserIdQ(id, userId), data)
	return Err(err)
}

func (m *MongoDatabase) DeleteByIdAndUserId(collection interface{}, id, userId string) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	err := coll.Remove(GetIdAndUserIdQ(id, userId))
	return Err(err)
}

func (m *MongoDatabase) DeleteAllByIdAndUserId(collection interface{}, id, userId string) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	_, err := coll.RemoveAll(GetIdAndUserIdQ(id, userId))
	return Err(err)
}

func (m *MongoDatabase) UpdateByQField(collection interface{}, q interface{}, field string, value interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	info, err := coll.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
	return Err(err) && info != nil && info.Matched > 0
}

func (m *MongoDatabase) UpdateByQMap(collection interface{}, q interface{}, v interface{}) bool {
	coll := m.getColl(collection)
	if coll == nil {
		return false
	}
	info, err := coll.UpdateAll(q, bson.M{"$set": v})
	return Err(err) && info != nil && info.Matched > 0
}

func (m *MongoDatabase) Count(collection interface{}, query interface{}) int {
	coll := m.getColl(collection)
	if coll == nil {
		return 0
	}
	cnt, err := coll.Find(query).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func (m *MongoDatabase) Distinct(collection interface{}, q bson.M, field string, result interface{}) {
	coll := m.getColl(collection)
	if coll == nil {
		return
	}
	coll.Find(q).Distinct(field, result)
}

func (m *MongoDatabase) DropIndex(collection interface{}, fields ...string) error {
	coll := m.getColl(collection)
	if coll == nil {
		return nil
	}
	return coll.DropIndex(fields...)
}
