// Command migration copies Gemsnote data between MongoDB and PostgreSQL.
// Both databases use the same 24-character ObjectId values, so relationships
// survive a round trip without an external ID mapping table.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var collections = []string{
	"users", "groups", "group_users", "notebooks", "notes", "note_contents", "note_content_histories",
	"tags", "note_tags", "tag_count", "albums", "files", "attachs", "note_images",
	"share_notebooks", "share_notes", "has_share_notes", "blogs", "user_blogs", "blog_singles",
	"blog_likes", "blog_comments", "reports", "configs", "sessions", "tokens", "email_logs",
	"suggestions", "themes",
}

type columnInfo struct {
	Name     string
	DataType string
	UDTName  string
}

type migrator struct {
	mongoSession *mgo.Session
	mongo        *mgo.Database
	postgres     *sql.DB
	dryRun       bool
}

func main() {
	direction := flag.String("direction", "", "mongo_to_pg or pg_to_mongo")
	mongoURL := flag.String("mongo-url", "mongodb://127.0.0.1:27017/gemsnote", "MongoDB URL")
	postgresURL := flag.String("postgres-url", "host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable", "PostgreSQL DSN")
	schema := flag.String("schema", "database/schema.sql", "schema applied before mongo_to_pg; empty disables it")
	dryRun := flag.Bool("dry-run", false, "read and validate without writing")
	validate := flag.Bool("validate", true, "compare collection/table counts after copying")
	flag.Parse()

	if *direction != "mongo_to_pg" && *direction != "pg_to_mongo" {
		log.Fatal("-direction must be mongo_to_pg or pg_to_mongo")
	}

	m, err := connect(*mongoURL, *postgresURL, *dryRun)
	if err != nil {
		log.Fatal(err)
	}
	defer m.close()

	if *direction == "mongo_to_pg" && !*dryRun && *schema != "" {
		if err := m.applySchema(*schema); err != nil {
			log.Fatalf("apply schema: %v", err)
		}
	}

	started := time.Now()
	for _, collection := range collections {
		var count int
		if *direction == "mongo_to_pg" {
			count, err = m.mongoToPostgres(collection)
		} else {
			count, err = m.postgresToMongo(collection)
		}
		if err != nil {
			log.Fatalf("migrate %s: %v", collection, err)
		}
		log.Printf("%s: %d records", collection, count)
	}

	if *validate {
		if err := m.validateCounts(); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("migration completed in %s", time.Since(started).Round(time.Millisecond))
}

func connect(mongoURL, postgresURL string, dryRun bool) (*migrator, error) {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return nil, fmt.Errorf("connect MongoDB: %w", err)
	}
	postgres, err := sql.Open("postgres", postgresURL)
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("open PostgreSQL: %w", err)
	}
	if err := postgres.Ping(); err != nil {
		session.Close()
		postgres.Close()
		return nil, fmt.Errorf("connect PostgreSQL: %w", err)
	}
	return &migrator{mongoSession: session, mongo: session.DB(""), postgres: postgres, dryRun: dryRun}, nil
}

func (m *migrator) close() {
	if m.mongoSession != nil {
		m.mongoSession.Close()
	}
	if m.postgres != nil {
		m.postgres.Close()
	}
}

func (m *migrator) applySchema(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = m.postgres.Exec(string(data))
	return err
}

func (m *migrator) mongoToPostgres(table string) (int, error) {
	columns, err := m.tableColumns(table)
	if err != nil {
		return 0, err
	}
	iter := m.mongo.C(table).Find(nil).Iter()
	defer iter.Close()

	count := 0
	var document bson.M
	for iter.Next(&document) {
		count++
		if m.dryRun {
			continue
		}
		row := make(map[string]interface{})
		for key, value := range document {
			column := mongoKeyToColumn(key)
			info, exists := columns[column]
			if !exists {
				continue
			}
			converted, include, err := mongoValueForPostgres(value, info)
			if err != nil {
				return count - 1, fmt.Errorf("document %v column %s: %w", document["_id"], column, err)
			}
			if include {
				row[column] = converted
			}
		}
		if _, ok := row["id"]; !ok {
			return count - 1, fmt.Errorf("document has no valid _id")
		}
		if err := m.upsertPostgres(table, row); err != nil {
			return count - 1, err
		}
		document = nil
	}
	if err := iter.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func (m *migrator) upsertPostgres(table string, row map[string]interface{}) error {
	keys := make([]string, 0, len(row))
	for key := range row {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Update by MongoDB _id first. A retried migration can otherwise hit a
	// secondary unique constraint before PostgreSQL applies ON CONFLICT (id).
	updateAssignments := make([]string, 0, len(keys)-1)
	updateValues := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		if key == "id" {
			continue
		}
		updateValues = append(updateValues, row[key])
		updateAssignments = append(updateAssignments,
			fmt.Sprintf("%s = $%d", pq.QuoteIdentifier(key), len(updateValues)))
	}
	updateValues = append(updateValues, row["id"])
	assignment := pq.QuoteIdentifier("id") + " = " + pq.QuoteIdentifier("id")
	if len(updateAssignments) > 0 {
		assignment = strings.Join(updateAssignments, ", ")
	}
	result, err := m.postgres.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		pq.QuoteIdentifier(table), assignment, len(updateValues)), updateValues...)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected > 0 {
		return nil
	}

	columns := make([]string, len(keys))
	placeholders := make([]string, len(keys))
	values := make([]interface{}, len(keys))
	updates := []string{}
	for i, key := range keys {
		quoted := pq.QuoteIdentifier(key)
		columns[i] = quoted
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		values[i] = row[key]
		if key != "id" {
			updates = append(updates, quoted+" = EXCLUDED."+quoted)
		}
	}
	conflict := "DO NOTHING"
	if len(updates) > 0 {
		conflict = "DO UPDATE SET " + strings.Join(updates, ", ")
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (id) %s",
		pq.QuoteIdentifier(table), strings.Join(columns, ", "), strings.Join(placeholders, ", "), conflict)
	_, err = m.postgres.Exec(query, values...)
	return err
}

func (m *migrator) postgresToMongo(table string) (int, error) {
	columns, err := m.tableColumns(table)
	if err != nil {
		return 0, err
	}
	rows, err := m.postgres.Query("SELECT * FROM " + pq.QuoteIdentifier(table))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	names, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	count := 0
	for rows.Next() {
		raw := make([]interface{}, len(names))
		dest := make([]interface{}, len(names))
		for i := range raw {
			dest[i] = &raw[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return count, err
		}
		document := bson.M{}
		for i, name := range names {
			value, include, err := postgresValueForMongo(raw[i], columns[name])
			if err != nil {
				return count, fmt.Errorf("column %s: %w", name, err)
			}
			if include {
				document[columnToMongoKey(name)] = value
			}
		}
		id, ok := document["_id"].(bson.ObjectId)
		if !ok || !id.Valid() {
			return count, fmt.Errorf("row has no valid id")
		}
		count++
		if !m.dryRun {
			if _, err := m.mongo.C(table).UpsertId(id, document); err != nil {
				return count - 1, err
			}
		}
	}
	return count, rows.Err()
}

func (m *migrator) tableColumns(table string) (map[string]columnInfo, error) {
	rows, err := m.postgres.Query(`SELECT column_name, data_type, udt_name
        FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]columnInfo{}
	for rows.Next() {
		var info columnInfo
		if err := rows.Scan(&info.Name, &info.DataType, &info.UDTName); err != nil {
			return nil, err
		}
		result[info.Name] = info
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("table does not exist; apply database/schema.sql first")
	}
	return result, rows.Err()
}

func mongoValueForPostgres(value interface{}, column columnInfo) (interface{}, bool, error) {
	value = normalizeMongoValue(value)
	if value == nil {
		return nil, false, nil
	}
	if column.UDTName == "_text" {
		items := interfaceSlice(value)
		stringsValue := make([]string, len(items))
		for i, item := range items {
			stringsValue[i] = fmt.Sprint(item)
		}
		return pq.Array(stringsValue), true, nil
	}
	if column.DataType == "jsonb" {
		encoded, err := json.Marshal(value)
		return encoded, true, err
	}
	return value, true, nil
}

func normalizeMongoValue(value interface{}) interface{} {
	switch v := value.(type) {
	case bson.ObjectId:
		return v.Hex()
	case bson.M:
		result := map[string]interface{}{}
		for key, item := range v {
			result[key] = normalizeMongoValue(item)
		}
		return result
	}
	rv := reflect.ValueOf(value)
	if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
		result := make([]interface{}, rv.Len())
		for i := range result {
			result[i] = normalizeMongoValue(rv.Index(i).Interface())
		}
		return result
	}
	return value
}

func postgresValueForMongo(raw interface{}, column columnInfo) (interface{}, bool, error) {
	if raw == nil {
		return nil, false, nil
	}
	if column.UDTName == "bpchar" {
		value := strings.TrimSpace(rawString(raw))
		if value == "" {
			return nil, false, nil
		}
		if !bson.IsObjectIdHex(value) {
			return nil, false, fmt.Errorf("invalid ObjectId %q", value)
		}
		return bson.ObjectIdHex(value), true, nil
	}
	if column.UDTName == "_text" {
		var values pq.StringArray
		if err := values.Scan(raw); err != nil {
			return nil, false, err
		}
		return []string(values), true, nil
	}
	if column.DataType == "jsonb" {
		var value interface{}
		if err := json.Unmarshal([]byte(rawString(raw)), &value); err != nil {
			return nil, false, err
		}
		return value, true, nil
	}
	if bytes, ok := raw.([]byte); ok {
		return string(bytes), true, nil
	}
	return raw, true, nil
}

func (m *migrator) validateCounts() error {
	for _, table := range collections {
		mongoCount, err := m.mongo.C(table).Find(nil).Count()
		if err != nil {
			return err
		}
		var postgresCount int
		if err := m.postgres.QueryRow("SELECT COUNT(*) FROM " + pq.QuoteIdentifier(table)).Scan(&postgresCount); err != nil {
			return err
		}
		if mongoCount != postgresCount {
			return fmt.Errorf("validation failed for %s: MongoDB=%d PostgreSQL=%d", table, mongoCount, postgresCount)
		}
	}
	log.Print("validation passed")
	return nil
}

var firstCapital = regexp.MustCompile("(.)([A-Z][a-z]+)")
var allCapital = regexp.MustCompile("([a-z0-9])([A-Z])")

func mongoKeyToColumn(key string) string {
	if key == "_id" {
		return "id"
	}
	if strings.EqualFold(key, "desc") {
		return "description"
	}
	value := firstCapital.ReplaceAllString(key, "${1}_${2}")
	return strings.ToLower(allCapital.ReplaceAllString(value, "${1}_${2}"))
}

func columnToMongoKey(column string) string {
	if column == "id" {
		return "_id"
	}
	special := map[string]string{
		"description": "Desc", "max_image_nums": "MaxImageNums",
	}
	if value, ok := special[column]; ok {
		return value
	}
	parts := strings.Split(column, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func interfaceSlice(value interface{}) []interface{} {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return []interface{}{value}
	}
	result := make([]interface{}, rv.Len())
	for i := range result {
		result[i] = rv.Index(i).Interface()
	}
	return result
}

func rawString(value interface{}) string {
	if bytes, ok := value.([]byte); ok {
		return string(bytes)
	}
	return fmt.Sprint(value)
}
