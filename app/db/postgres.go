package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	. "github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PostgresDatabase struct {
	db *sql.DB
}

func newPostgresDB() (Database, error) {
	p := &PostgresDatabase{}
	err := p.Initialize()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PostgresDatabase) Initialize() error {
	config := revel.Config

	var url string
	var ok bool

	url = initURL
	ok = url != ""
	if !ok && config != nil {
		url, ok = config.String("db.url")
	}
	if !ok && config != nil {
		url, ok = config.String("db.urlEnv")
		if ok {
			Log("get PostgreSQL config from db.urlEnv")
		}
	} else {
		Log("get PostgreSQL config from db.url")
	}

	if !ok && config != nil {
		host, _ := config.String("db.host")
		portStr, _ := config.String("db.port")
		port, _ := strconv.Atoi(portStr)
		if port == 0 {
			port = 5432
		}
		username, _ := config.String("db.username")
		password, _ := config.String("db.password")
		dbname, _ := config.String("db.dbname")
		sslmode, _ := config.String("db.sslmode")
		if sslmode == "" {
			sslmode = "disable"
		}
		url = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, port, username, password, dbname, sslmode)
	}
	if url == "" {
		return fmt.Errorf("PostgreSQL connection URL is not configured")
	}

	var err error
	p.db, err = sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	p.db.SetMaxOpenConns(25)
	p.db.SetMaxIdleConns(25)
	p.db.SetConnMaxLifetime(5 * time.Minute)

	err = p.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	Log("Connected to PostgreSQL successfully")
	if err := initializePostgresDatabase(p.db); err != nil {
		p.db.Close()
		return fmt.Errorf("failed to initialize PostgreSQL database: %w", err)
	}
	if _, err := p.db.Exec(`ALTER TABLE notes ADD COLUMN IF NOT EXISTS is_star BOOLEAN NOT NULL DEFAULT FALSE`); err != nil {
		p.db.Close()
		return fmt.Errorf("failed to ensure PostgreSQL star schema: %w", err)
	}
	return nil
}

// initializePostgresDatabase installs the bundled schema and seed only when
// the public schema has no business base tables. The migration marker table
// is intentionally ignored because it can be left behind by an interrupted
// first-run initialization. A database with any other existing table (even if
// it has no rows) must never receive the installation seed.
func initializePostgresDatabase(database *sql.DB) (resultErr error) {
	ctx := context.Background()
	conn, err := database.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Serialize first-run initialization across processes. The lock is held on
	// this dedicated connection and the emptiness check is repeated after it is
	// acquired, so concurrent server starts cannot both install the seed.
	lockCtx, cancelLock := context.WithTimeout(ctx, 10*time.Second)
	defer cancelLock()
	if _, err := conn.ExecContext(lockCtx, postgresBootstrapLockSQL); err != nil {
		return err
	}
	defer conn.ExecContext(ctx, postgresBootstrapUnlockSQL)
	// schema.sql/seed.sql contain session-level SET statements (including an
	// empty search_path). Reset the dedicated connection before it is returned
	// to database/sql, on both success and failure paths.
	defer func() {
		if _, err := conn.ExecContext(context.Background(), `RESET ALL`); err != nil && resultErr == nil {
			resultErr = fmt.Errorf("reset PostgreSQL bootstrap session: %w", err)
		}
	}()

	var tableCount int
	if err := conn.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
		  AND table_name <> 'gemsnote_schema_migrations'
	`).Scan(&tableCount); err != nil {
		return err
	}
	if tableCount != 0 {
		// A failed container init can leave schema.sql applied while seed.sql
		// never ran. Recover only that schema-only state: every public business
		// table must still be empty, so an existing database with real data is
		// never populated with the installation seed.
		empty, err := postgresBusinessTablesEmpty(ctx, conn)
		if err != nil {
			return err
		}
		if !empty {
			return nil
		}
	}

	databaseDir, err := findInstallationDatabaseDir()
	if err != nil {
		return err
	}
	schema, err := os.ReadFile(filepath.Join(databaseDir, "schema.sql"))
	if err != nil {
		return fmt.Errorf("read schema.sql: %w", err)
	}
	seed, err := os.ReadFile(filepath.Join(databaseDir, "seed.sql"))
	if err != nil {
		return fmt.Errorf("read seed.sql: %w", err)
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// The checked-in scripts are also usable directly by Docker/PostgreSQL,
	// so they contain their own BEGIN/COMMIT. Remove only those standalone
	// transaction lines when running both scripts in one outer transaction.
	if _, err := tx.ExecContext(ctx, stripScriptTransactionBoundaries(string(schema))); err != nil {
		tx.Rollback()
		return fmt.Errorf("execute schema.sql: %w", err)
	}
	if _, err := tx.ExecContext(ctx, stripScriptTransactionBoundaries(string(seed))); err != nil {
		tx.Rollback()
		return fmt.Errorf("execute seed.sql: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit PostgreSQL initialization: %w", err)
	}
	Log("Initialized PostgreSQL schema and installation seed")
	return nil
}

func postgresBusinessTablesEmpty(ctx context.Context, conn *sql.Conn) (bool, error) {
	rows, err := conn.QueryContext(ctx, `
		SELECT table_name FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
		  AND table_name <> 'gemsnote_schema_migrations'
	`)
	if err != nil {
		return false, err
	}
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return false, err
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return false, err
	}
	if err := rows.Close(); err != nil {
		return false, err
	}
	for _, table := range tables {
		identifier := `"` + strings.ReplaceAll(table, `"`, `""`) + `"`
		var count int
		if err := conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM public.`+identifier).Scan(&count); err != nil {
			return false, err
		}
		if count != 0 {
			return false, nil
		}
	}
	return true, nil
}

const (
	postgresBootstrapLockSQL   = `SELECT pg_advisory_lock(hashtextextended('gemsnote:postgres-bootstrap:v1', 0))`
	postgresBootstrapUnlockSQL = `SELECT pg_advisory_unlock(hashtextextended('gemsnote:postgres-bootstrap:v1', 0))`
)

var scriptTransactionBoundary = regexp.MustCompile(`(?im)^\s*(BEGIN|COMMIT)\s*;\s*$`)

func stripScriptTransactionBoundaries(script string) string {
	return scriptTransactionBoundary.ReplaceAllString(script, "")
}

func findInstallationDatabaseDir() (string, error) {
	cwd, _ := os.Getwd()
	executable, _ := os.Executable()
	return findInstallationDatabaseDirFrom(strings.TrimSpace(os.Getenv("GEMSNOTE_DATABASE_DIR")), cwd, executable)
}

func findInstallationDatabaseDirFrom(configured, cwd, executable string) (string, error) {
	var candidates []string
	// An explicit override is always considered first. Bundled release files
	// take precedence over cwd so an unrelated directory cannot shadow them.
	if configured != "" {
		candidates = append(candidates, configured)
	}
	if executable != "" {
		dir := filepath.Dir(executable)
		candidates = append(candidates,
			filepath.Join(dir, "database"),
			filepath.Join(dir, "..", "database"),
			filepath.Join(dir, "..", "..", "database"),
		)
	}
	if cwd != "" {
		candidates = append(candidates, filepath.Join(cwd, "database"))
	}
	for _, candidate := range candidates {
		candidate, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(candidate, "schema.sql")); err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(candidate, "seed.sql")); err != nil {
			continue
		}
		return candidate, nil
	}
	return "", fmt.Errorf("cannot locate database/schema.sql and database/seed.sql; set GEMSNOTE_DATABASE_DIR or run from the project/release root")
}

func (p *PostgresDatabase) SetupLegacyVariables() {
	Session = nil
	Notebooks = "notebooks"
	Notes = "notes"
	NoteContents = "note_contents"
	NoteContentHistories = "note_content_histories"
	ShareNotes = "share_notes"
	ShareNotebooks = "share_notebooks"
	HasShareNotes = "has_share_notes"
	Blogs = "blogs"
	Users = "users"
	Groups = "groups"
	GroupUsers = "group_users"
	Tags = "tags"
	NoteTags = "note_tags"
	TagCounts = "tag_count"
	UserBlogs = "user_blogs"
	Tokens = "tokens"
	Suggestions = "suggestions"
	Albums = "albums"
	Files = "files"
	Attachs = "attachs"
	NoteImages = "note_images"
	Configs = "configs"
	EmailLogs = "email_logs"
	BlogLikes = "blog_likes"
	BlogComments = "blog_comments"
	Reports = "reports"
	BlogSingles = "blog_singles"
	Themes = "themes"
	Sessions = "sessions"
}

func (p *PostgresDatabase) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDatabase) Ping() error {
	return p.db.Ping()
}

func (p *PostgresDatabase) CheckConnection() {
	err := p.db.Ping()
	if err != nil {
		Log("Lost connection to database!")
		err = p.db.Ping()
		if err == nil {
			Log("Reconnect to database successful.")
		} else {
			Log("Reconnect failed!!!! Warning")
		}
	}
}

func (p *PostgresDatabase) NewID() string {
	return bson.NewObjectId().Hex()
}

func (p *PostgresDatabase) IsValidID(id string) bool {
	return bson.IsObjectIdHex(id)
}

func (p *PostgresDatabase) GetType() string {
	return "postgresql"
}

func (p *PostgresDatabase) AppliedMigrations() ([]string, error) {
	if _, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS gemsnote_schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return nil, err
	}
	rows, err := p.db.Query("SELECT version FROM gemsnote_schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

func (p *PostgresDatabase) RecordMigration(version string) error {
	_, err := p.db.Exec(`INSERT INTO gemsnote_schema_migrations (version)
		VALUES ($1) ON CONFLICT (version) DO NOTHING`, version)
	return err
}

func getTableName(collection interface{}) string {
	switch c := collection.(type) {
	case *mgo.Collection:
		return c.Name
	case string:
		return c
	default:
		return fmt.Sprintf("%v", collection)
	}
}

func (p *PostgresDatabase) getColl(collection interface{}) string {
	return getTableName(collection)
}

func getDBTag(field reflect.StructField) string {
	dbTag := field.Tag.Get("db")
	if dbTag == "-" {
		return "-"
	}
	if dbTag != "" {
		return dbTag
	}
	bsonTag := field.Tag.Get("bson")
	if bsonTag != "" && bsonTag != "-" {
		parts := strings.Split(bsonTag, ",")
		name := parts[0]
		if name == "_id" {
			return "id"
		}
		if name == "omitempty" || name == "" {
			return postgresColumn(field.Name)
		}
		return postgresColumn(name)
	}
	return postgresColumn(field.Name)
}

func getDBColumns(typ reflect.Type) []string {
	var columns []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := getDBTag(field)
		if tag == "-" {
			continue
		}
		if field.Anonymous {
			subCols := getDBColumns(field.Type)
			columns = append(columns, subCols...)
			continue
		}
		columns = append(columns, tag)
	}
	return columns
}

func (p *PostgresDatabase) Insert(collection interface{}, data interface{}) bool {
	table := p.getColl(collection)
	if dataMap, ok := data.(bson.M); ok {
		keys := make([]string, 0, len(dataMap))
		for key := range dataMap {
			if !strings.HasPrefix(key, "$") {
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)
		columns := make([]string, 0, len(keys))
		placeholders := make([]string, 0, len(keys))
		values := make([]interface{}, 0, len(keys))
		for _, key := range keys {
			columns = append(columns, postgresColumn(key))
			values = append(values, postgresValue(dataMap[key]))
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
		}
		_, err := p.db.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
			strings.Join(columns, ", "), strings.Join(placeholders, ", ")), values...)
		if err != nil {
			Log("Insert error: " + err.Error())
		}
		return err == nil
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}

	typ := val.Type()
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := getDBTag(field)
		if tag == "-" {
			continue
		}
		if field.Anonymous {
			continue
		}

		fieldVal := val.Field(i)

		columns = append(columns, tag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
		values = append(values, postgresValue(fieldVal.Interface()))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err := p.db.Exec(query, values...)
	if err != nil {
		Log("Insert error: " + err.Error())
		return false
	}
	return true
}

func (p *PostgresDatabase) Update(collection interface{}, query interface{}, data interface{}) bool {
	return p.updateWithFilter(collection, query, data, false)
}

func (p *PostgresDatabase) UpdateAll(collection interface{}, query interface{}, data interface{}) bool {
	return p.updateWithFilter(collection, query, data, true)
}

func (p *PostgresDatabase) updateWithFilter(collection interface{}, query interface{}, data interface{}, all bool) bool {
	table := p.getColl(collection)
	setParts, setValues := buildPostgresSet(data, 1)
	if len(setParts) == 0 {
		return false
	}
	whereClause, whereArgs := buildPostgresWhere(table, asBSONMap(query), len(setValues)+1)

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table, strings.Join(setParts, ", "), whereClause)

	args := append(setValues, whereArgs...)

	result, err := p.db.Exec(sql, args...)
	if err != nil {
		Log("Update error: " + err.Error())
		return false
	}
	affected, err := result.RowsAffected()
	return err == nil && affected > 0
}

func (p *PostgresDatabase) Upsert(collection interface{}, query interface{}, data interface{}) bool {
	table := p.getColl(collection)
	setParts, setValues := buildPostgresSet(data, 1)
	if len(setParts) == 0 {
		return false
	}
	whereClause, whereArgs := buildPostgresWhere(table, asBSONMap(query), len(setValues)+1)
	result, err := p.db.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setParts, ", "), whereClause),
		append(setValues, whereArgs...)...)
	if err != nil {
		Log("Upsert update error: " + err.Error())
		return false
	}
	affected, err := result.RowsAffected()
	if err == nil && affected > 0 {
		return true
	}
	return p.Insert(collection, postgresUpsertDocument(asBSONMap(query), data))
}

func (p *PostgresDatabase) Delete(collection interface{}, query interface{}) bool {
	return p.deleteWithFilter(collection, query, false)
}

func (p *PostgresDatabase) DeleteAll(collection interface{}, query interface{}) bool {
	return p.deleteWithFilter(collection, query, true)
}

func (p *PostgresDatabase) deleteWithFilter(collection interface{}, query interface{}, all bool) bool {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, asBSONMap(query), 1)

	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause)

	_, err := p.db.Exec(sql, args...)
	if err != nil {
		Log("Delete error: " + err.Error())
		return false
	}
	return true
}

func (p *PostgresDatabase) Get(collection interface{}, id string, result interface{}) {
	table := p.getColl(collection)
	p.getById(table, id, result)
}

func (p *PostgresDatabase) getById(table, id string, result interface{}) {
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	columns := getDBColumns(typ)

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1",
		strings.Join(columns, ", "), table)

	rows, err := p.db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := scanPostgresStruct(rows, val, typ); err != nil {
			Log("Get scan error: " + err.Error())
		}
	}
}

func (p *PostgresDatabase) GetByQ(collection interface{}, query interface{}, result interface{}) {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, asBSONMap(query), 1)

	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	columns := getDBColumns(typ)

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1",
		strings.Join(columns, ", "), table, whereClause)

	rows, err := p.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := scanPostgresStruct(rows, val, typ); err != nil {
			Log("GetByQ scan error: " + err.Error())
		}
	}
}

func (p *PostgresDatabase) ListByQ(collection interface{}, query interface{}, result interface{}) {
	p.listByQWithOpts(collection, query, result, QueryOptions{})
}

func (p *PostgresDatabase) ListByQLimit(collection interface{}, query interface{}, result interface{}, limit int) {
	p.listByQWithOpts(collection, query, result, QueryOptions{Limit: limit})
}

func (p *PostgresDatabase) ListByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions) {
	p.listByQWithOpts(collection, query, result, options)
}

func (p *PostgresDatabase) GetByQOptions(collection interface{}, query interface{}, result interface{}, options QueryOptions) {
	items := reflect.New(reflect.SliceOf(reflect.TypeOf(result).Elem()))
	p.ListByQOptions(collection, query, items.Interface(), QueryOptions{Sort: options.Sort, Skip: options.Skip, Limit: 1, Fields: options.Fields})
	if items.Elem().Len() > 0 {
		reflect.ValueOf(result).Elem().Set(items.Elem().Index(0))
	}
}

func (p *PostgresDatabase) listByQWithOpts(collection interface{}, query interface{}, result interface{}, options QueryOptions) {
	table := p.getColl(collection)

	sliceVal := reflect.ValueOf(result)
	if sliceVal.Kind() != reflect.Ptr || sliceVal.Elem().Kind() != reflect.Slice {
		return
	}
	sliceVal = sliceVal.Elem()

	elemTyp := sliceVal.Type().Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	columns := getDBColumns(elemTyp)
	if len(options.Fields) > 0 {
		columns = make([]string, len(options.Fields))
		for i, field := range options.Fields {
			columns[i] = postgresColumn(field)
		}
	}

	whereClause, args := buildPostgresWhere(table, asBSONMap(query), 1)

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), table, whereClause)

	if len(options.Sort) > 0 {
		orderParts := make([]string, 0, len(options.Sort))
		for _, sortField := range options.Sort {
			orderDir := "ASC"
			if strings.HasPrefix(sortField, "-") {
				sortField = sortField[1:]
				orderDir = "DESC"
			}
			orderParts = append(orderParts, fmt.Sprintf("%s %s", postgresColumn(sortField), orderDir))
		}
		sql += " ORDER BY " + strings.Join(orderParts, ", ")
	}

	if options.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", options.Limit)
	}
	if options.Skip > 0 {
		sql += fmt.Sprintf(" OFFSET %d", options.Skip)
	}

	rows, err := p.db.Query(sql, args...)
	if err != nil {
		Log("ListByQ error: " + err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		elem := reflect.New(elemTyp).Elem()
		if err := scanPostgresStruct(rows, elem, elemTyp); err != nil {
			Log("ListByQ scan error: " + err.Error())
			continue
		}
		if sliceVal.Type().Elem().Kind() == reflect.Ptr {
			sliceVal.Set(reflect.Append(sliceVal, elem.Addr()))
		} else {
			sliceVal.Set(reflect.Append(sliceVal, elem))
		}
	}
}

func (p *PostgresDatabase) GetByQWithFields(collection interface{}, query bson.M, fields []string, result interface{}) {
	p.getByQWithFieldsInternal(collection, query, fields, result)
}

func (p *PostgresDatabase) getByQWithFieldsInternal(collection interface{}, query bson.M, fields []string, result interface{}) {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, query, 1)

	dbFields := make([]string, len(fields))
	for i, f := range fields {
		dbFields[i] = toSnakeCase(f)
		if f == "_id" {
			dbFields[i] = "id"
		}
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1",
		strings.Join(dbFields, ", "), table, whereClause)

	rows, err := p.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()

	if rows.Next() {
		if err := scanPostgresStruct(rows, val, typ); err != nil {
			Log("GetByQWithFields scan error: " + err.Error())
		}
	}
}

func (p *PostgresDatabase) ListByQWithFields(collection interface{}, query bson.M, fields []string, result interface{}) {
	p.listByQWithFieldsInternal(collection, query, fields, result)
}

func (p *PostgresDatabase) listByQWithFieldsInternal(collection interface{}, query bson.M, fields []string, result interface{}) {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, query, 1)

	dbFields := make([]string, len(fields))
	for i, f := range fields {
		dbFields[i] = toSnakeCase(f)
		if f == "_id" {
			dbFields[i] = "id"
		}
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(dbFields, ", "), table, whereClause)

	rows, err := p.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	sliceVal := reflect.ValueOf(result)
	if sliceVal.Kind() != reflect.Ptr || sliceVal.Elem().Kind() != reflect.Slice {
		return
	}
	sliceVal = sliceVal.Elem()
	elemTyp := sliceVal.Type().Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	for rows.Next() {
		elem := reflect.New(elemTyp).Elem()
		if err := scanPostgresStruct(rows, elem, elemTyp); err != nil {
			Log("ListByQWithFields scan error: " + err.Error())
			continue
		}
		if sliceVal.Type().Elem().Kind() == reflect.Ptr {
			sliceVal.Set(reflect.Append(sliceVal, elem.Addr()))
		} else {
			sliceVal.Set(reflect.Append(sliceVal, elem))
		}
	}
}

func (p *PostgresDatabase) GetByIdAndUserId(collection interface{}, id, userId string, result interface{}) {
	query := bson.M{"_id": id, "UserId": userId}
	p.GetByQ(collection, query, result)
}

func (p *PostgresDatabase) UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool {
	query := bson.M{"_id": id, "UserId": userId}
	return p.Update(collection, query, data)
}

func (p *PostgresDatabase) DeleteByIdAndUserId(collection interface{}, id, userId string) bool {
	query := bson.M{"_id": id, "UserId": userId}
	return p.Delete(collection, query)
}

func (p *PostgresDatabase) DeleteAllByIdAndUserId(collection interface{}, id, userId string) bool {
	query := bson.M{"_id": id, "UserId": userId}
	return p.DeleteAll(collection, query)
}

func (p *PostgresDatabase) UpdateByQField(collection interface{}, q interface{}, field string, value interface{}) bool {
	return p.UpdateByQMap(collection, q, bson.M{field: value})
}

func (p *PostgresDatabase) UpdateByQMap(collection interface{}, q interface{}, v interface{}) bool {
	return p.updateWithFilter(collection, q, bson.M{"$set": postgresUpdateMap(v)}, true)
}

func (p *PostgresDatabase) Count(collection interface{}, query interface{}) int {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, asBSONMap(query), 1)

	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, whereClause)

	var count int
	err := p.db.QueryRow(sql, args...).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (p *PostgresDatabase) Distinct(collection interface{}, q bson.M, field string, result interface{}) {
	table := p.getColl(collection)

	whereClause, args := buildPostgresWhere(table, q, 1)

	sql := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s", toSnakeCase(field), table, whereClause)

	rows, err := p.db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	sliceVal := reflect.ValueOf(result)
	if sliceVal.Kind() != reflect.Ptr || sliceVal.Elem().Kind() != reflect.Slice {
		return
	}
	sliceVal = sliceVal.Elem()

	for rows.Next() {
		elem := reflect.New(sliceVal.Type().Elem()).Elem()
		var raw interface{}
		if err := rows.Scan(&raw); err != nil {
			Log("Distinct scan error: " + err.Error())
			continue
		}
		if err := setPostgresField(elem, raw); err != nil {
			Log("Distinct conversion error: " + err.Error())
			continue
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}
}

func (p *PostgresDatabase) DropIndex(collection interface{}, fields ...string) error {
	// Legacy upgrades drop a MongoDB index by its field list. PostgreSQL indexes
	// are managed by schema migrations, so there is no equivalent runtime action.
	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
