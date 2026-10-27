package db

import (
	"encoding/json"
	"errors"
	"html"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
)

// Reads in the cleanup protocol must report failures; the older convenience
// query helpers deliberately discard errors and cannot prove safe deletion.
type imageUploadStore interface {
	checkedImageFiles(bson.M) ([]info.File, error)
}

func (m *MongoDatabase) checkedImageFiles(query bson.M) ([]info.File, error) {
	var files []info.File
	err := m.db.C("files").Find(query).All(&files)
	return files, err
}

func (p *PostgresDatabase) checkedImageFiles(query bson.M) ([]info.File, error) {
	where, args := buildPostgresWhere("files", query, 1)
	typ := reflect.TypeOf(info.File{})
	rows, err := p.db.Query("SELECT "+strings.Join(getDBColumns(typ), ", ")+" FROM files WHERE "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []info.File
	for rows.Next() {
		var file info.File
		if err := scanPostgresStruct(rows, reflect.ValueOf(&file).Elem(), typ); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, rows.Err()
}

var imageIDPattern = regexp.MustCompile(`[a-fA-F0-9]{24}`)
var imagePathPattern = regexp.MustCompile(`(?i)files/[a-z0-9/_.-]+`)

func imageReferenceQuery(data interface{}) (bson.M, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var value interface{}
	if err := json.Unmarshal(encoded, &value); err != nil {
		return nil, err
	}
	ids, paths := map[bson.ObjectId]bool{}, map[string]bool{}
	var walk func(interface{})
	walk = func(value interface{}) {
		switch v := value.(type) {
		case string:
			for i := 0; i < 4; i++ {
				v = html.UnescapeString(strings.ReplaceAll(v, `\/`, `/`))
				if decoded, err := url.PathUnescape(v); err == nil {
					v = decoded
				}
			}
			for _, id := range imageIDPattern.FindAllString(v, -1) {
				ids[bson.ObjectIdHex(id)] = true
			}
			for _, name := range imagePathPattern.FindAllString(v, -1) {
				paths[path.Clean(name)] = true
			}
		case []interface{}:
			for _, item := range v {
				walk(item)
			}
		case map[string]interface{}:
			for _, item := range v {
				walk(item)
			}
		}
	}
	walk(value)
	var alternatives []bson.M
	if len(ids) > 0 {
		values := make([]bson.ObjectId, 0, len(ids))
		for id := range ids {
			values = append(values, id)
		}
		alternatives = append(alternatives, bson.M{"_id": bson.M{"$in": values}})
	}
	if len(paths) > 0 {
		values := make([]string, 0, len(paths))
		for name := range paths {
			values = append(values, name)
		}
		alternatives = append(alternatives, bson.M{"Path": bson.M{"$in": values}})
	}
	if len(alternatives) == 0 {
		return nil, nil
	}
	return bson.M{"$or": alternatives}, nil
}

// Claim BEFORE any current body, cover, membership, profile or history write.
// A database CAS arbitrates against cleanup across processes on both backends.
// Retention is permanent (also after history trimming): conservative leaks are
// preferable to deleting an original that a saved version may still need.
func protectImageReferences(database Database, collection, data interface{}) bool {
	switch getTableName(collection) {
	case "notes", "note_contents", "note_content_histories", "note_images", "users":
	default:
		return true
	}
	query, err := imageReferenceQuery(data)
	if err != nil {
		return false
	}
	if query == nil {
		return true
	}
	store, ok := database.(imageUploadStore)
	if !ok {
		return false
	}
	files, err := store.checkedImageFiles(query)
	if err != nil {
		return false
	}
	for _, file := range files {
		switch file.UploadState {
		case "", "retained":
			continue
		case "pending":
			if database.UpdateByQMap(imageCollection(database), bson.M{"_id": file.FileId, "UploadState": "pending"}, bson.M{"UploadState": "retained"}) {
				continue
			}
			current, err := store.checkedImageFiles(bson.M{"_id": file.FileId})
			if err == nil && len(current) == 1 && current[0].UploadState == "retained" {
				continue
			}
		}
		return false // deleting/deleted or an ambiguous database result
	}
	return true
}

func imageCollection(database Database) interface{} {
	if mongo, ok := database.(*MongoDatabase); ok {
		return mongo.db.C("files")
	}
	return "files"
}

func ImageUpload(userID, fileID string) (info.File, error) {
	if !bson.IsObjectIdHex(fileID) {
		return info.File{}, errors.New("invalidRequest")
	}
	return imageUploadByQuery(userID, bson.M{"_id": bson.ObjectIdHex(fileID)})
}

func ImageUploadByKey(userID, key string) (info.File, error) {
	return imageUploadByQuery(userID, bson.M{"UploadKey": key})
}

func imageUploadByQuery(userID string, query bson.M) (info.File, error) {
	store, ok := globalDatabase.(imageUploadStore)
	if !ok || !bson.IsObjectIdHex(userID) {
		return info.File{}, errors.New("invalidRequest")
	}
	query["UserId"] = bson.ObjectIdHex(userID)
	files, err := store.checkedImageFiles(query)
	if err != nil {
		return info.File{}, errors.New("databaseUnavailable")
	}
	if len(files) != 1 || files[0].UploadState == "" {
		return info.File{}, errors.New("uploadNotFound")
	}
	return files[0], nil
}

func ReserveImageCleanup(userID, fileID string) (info.File, error) {
	file, err := ImageUpload(userID, fileID)
	if err != nil {
		return file, err
	}
	if file.UploadState == "pending" {
		if globalDatabase.UpdateByQMap(imageCollection(globalDatabase), bson.M{"_id": file.FileId, "UserId": file.UserId, "UploadState": "pending"}, bson.M{"UploadState": "deleting"}) {
			file.UploadState = "deleting"
			return file, nil
		}
		file, err = ImageUpload(userID, fileID)
		if err != nil {
			return file, err
		}
	}
	if file.UploadState == "deleting" || file.UploadState == "deleted" {
		return file, nil
	}
	return info.File{}, errors.New("imageAlreadyReferenced")
}

func CompleteImageCleanup(file info.File) bool {
	return globalDatabase.UpdateByQMap(imageCollection(globalDatabase), bson.M{"_id": file.FileId, "UserId": file.UserId, "UploadState": "deleting"}, bson.M{"UploadState": "deleted"})
}
