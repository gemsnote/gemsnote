package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"gopkg.in/mgo.v2/bson"
)

var bsonObjectIDType = reflect.TypeOf(bson.ObjectId(""))
var timeType = reflect.TypeOf(time.Time{})

func asBSONMap(value interface{}) bson.M {
	if value == nil {
		return bson.M{}
	}
	if q, ok := value.(bson.M); ok {
		return q
	}
	if q, ok := value.(map[string]interface{}); ok {
		return bson.M(q)
	}
	return bson.M{"$invalid": true}
}

func postgresColumn(field string) string {
	if field == "_id" {
		return "id"
	}
	if field == "Desc" || field == "desc" {
		return "description"
	}
	column := toSnakeCase(field)
	for _, r := range column {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return "invalid_column"
		}
	}
	return column
}

func postgresValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case bson.ObjectId:
		if len(v) == 0 {
			return nil
		}
		return v.Hex()
	case []string:
		if v == nil {
			v = []string{}
		}
		return pq.Array(v)
	case []bson.ObjectId:
		ids := make([]string, 0, len(v))
		for _, id := range v {
			ids = append(ids, id.Hex())
		}
		return pq.Array(ids)
	case time.Time:
		if v.IsZero() {
			return nil
		}
		return v
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		return postgresValue(rv.Elem().Interface())
	}
	if rv.Kind() == reflect.Map || rv.Kind() == reflect.Struct || rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		encoded, err := json.Marshal(value)
		if err == nil {
			return encoded
		}
	}
	return value
}

func postgresUpdateMap(value interface{}) bson.M {
	if value == nil {
		return bson.M{}
	}
	if result, ok := value.(bson.M); ok {
		return result
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return bson.M{}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return bson.M{}
	}
	rt := rv.Type()
	result := bson.M{}
	for i := 0; i < rv.NumField(); i++ {
		fieldType := rt.Field(i)
		if fieldType.Anonymous || getDBTag(fieldType) == "-" {
			continue
		}
		result[fieldType.Name] = rv.Field(i).Interface()
	}
	return result
}

func buildPostgresSet(data interface{}, start int) ([]string, []interface{}) {
	update := postgresUpdateMap(data)
	keys := make([]string, 0, len(update))
	for key := range update {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := []string{}
	args := []interface{}{}
	next := start
	appendSet := func(field string, value interface{}) {
		parts = append(parts, fmt.Sprintf("%s = $%d", postgresColumn(field), next))
		args = append(args, postgresValue(value))
		next++
	}
	for _, key := range keys {
		value := update[key]
		if !strings.HasPrefix(key, "$") {
			appendSet(key, value)
			continue
		}
		fields := asBSONMap(value)
		fieldNames := make([]string, 0, len(fields))
		for field := range fields {
			fieldNames = append(fieldNames, field)
		}
		sort.Strings(fieldNames)
		for _, field := range fieldNames {
			column := postgresColumn(field)
			fieldValue := fields[field]
			switch key {
			case "$set":
				appendSet(field, fieldValue)
			case "$inc":
				parts = append(parts, fmt.Sprintf("%s = COALESCE(%s, 0) + $%d", column, column, next))
				args = append(args, postgresValue(fieldValue))
				next++
			case "$push":
				parts = append(parts, fmt.Sprintf("%s = array_append(COALESCE(%s, '{}'), $%d)", column, column, next))
				args = append(args, postgresValue(fieldValue))
				next++
			case "$pull":
				parts = append(parts, fmt.Sprintf("%s = array_remove(COALESCE(%s, '{}'), $%d)", column, column, next))
				args = append(args, postgresValue(fieldValue))
				next++
			case "$addToSet":
				parts = append(parts, fmt.Sprintf("%s = CASE WHEN $%d = ANY(COALESCE(%s, '{}')) THEN %s ELSE array_append(COALESCE(%s, '{}'), $%d) END", column, next, column, column, column, next))
				args = append(args, postgresValue(fieldValue))
				next++
			}
		}
	}
	return parts, args
}

func postgresUpsertDocument(query bson.M, data interface{}) bson.M {
	document := bson.M{}
	for key, value := range query {
		if strings.HasPrefix(key, "$") {
			continue
		}
		if _, operator := value.(bson.M); !operator {
			document[key] = value
		}
	}
	update := postgresUpdateMap(data)
	for key, value := range update {
		if !strings.HasPrefix(key, "$") {
			document[key] = value
			continue
		}
		for field, fieldValue := range asBSONMap(value) {
			switch key {
			case "$set", "$inc":
				document[field] = fieldValue
			case "$push", "$addToSet":
				document[field] = []string{fmt.Sprint(fieldValue)}
			}
		}
	}
	return document
}

func buildPostgresWhere(table string, query bson.M, start int) (string, []interface{}) {
	next := start
	clause, args := buildPostgresExpression(table, query, &next)
	if clause == "" {
		return "TRUE", args
	}
	return clause, args
}

func buildPostgresExpression(table string, query bson.M, next *int) (string, []interface{}) {
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	conditions := make([]string, 0, len(keys))
	args := []interface{}{}
	for _, key := range keys {
		value := query[key]
		switch key {
		case "$or", "$and":
			items := valueSlice(value)
			parts := make([]string, 0, len(items))
			for _, item := range items {
				sub, subArgs := buildPostgresExpression(table, asBSONMap(item), next)
				if sub != "" {
					parts = append(parts, "("+sub+")")
					args = append(args, subArgs...)
				}
			}
			if len(parts) == 0 {
				conditions = append(conditions, "FALSE")
			} else {
				joiner := " OR "
				if key == "$and" {
					joiner = " AND "
				}
				conditions = append(conditions, "("+strings.Join(parts, joiner)+")")
			}
			continue
		}
		if strings.HasPrefix(key, "$") {
			conditions = append(conditions, "FALSE")
			continue
		}

		column := postgresColumn(key)
		operatorMap, isOperator := value.(bson.M)
		if !isOperator {
			if raw, ok := value.(map[string]interface{}); ok {
				operatorMap, isOperator = bson.M(raw), true
			}
		}
		if !isOperator {
			if value == nil {
				conditions = append(conditions, column+" IS NULL")
				continue
			}
			conditions = append(conditions, fmt.Sprintf("%s = $%d", column, *next))
			args = append(args, postgresValue(value))
			*next++
			continue
		}

		opKeys := make([]string, 0, len(operatorMap))
		for op := range operatorMap {
			opKeys = append(opKeys, op)
		}
		sort.Strings(opKeys)
		for _, op := range opKeys {
			opValue := operatorMap[op]
			switch op {
			case "$gt", "$gte", "$lt", "$lte":
				symbol := map[string]string{"$gt": ">", "$gte": ">=", "$lt": "<", "$lte": "<="}[op]
				conditions = append(conditions, fmt.Sprintf("%s %s $%d", column, symbol, *next))
				args = append(args, postgresValue(opValue))
				*next++
			case "$ne":
				conditions = append(conditions, fmt.Sprintf("%s IS DISTINCT FROM $%d", column, *next))
				args = append(args, postgresValue(opValue))
				*next++
			case "$in", "$nin":
				values := valueSlice(opValue)
				if len(values) == 0 {
					conditions = append(conditions, map[bool]string{true: "TRUE", false: "FALSE"}[op == "$nin"])
					continue
				}
				if isPostgresArrayColumn(table, column) {
					stringsValue := make([]string, len(values))
					for i, v := range values {
						stringsValue[i] = fmt.Sprint(postgresValue(v))
					}
					expression := fmt.Sprintf("%s && $%d", column, *next)
					if op == "$nin" {
						expression = "NOT (" + expression + ")"
					}
					conditions = append(conditions, expression)
					args = append(args, pq.Array(stringsValue))
					*next++
					continue
				}
				placeholders := make([]string, len(values))
				for i, v := range values {
					placeholders[i] = fmt.Sprintf("$%d", *next)
					args = append(args, postgresValue(v))
					*next++
				}
				not := ""
				if op == "$nin" {
					not = " NOT"
				}
				conditions = append(conditions, fmt.Sprintf("%s%s IN (%s)", column, not, strings.Join(placeholders, ", ")))
			case "$all":
				values := valueSlice(opValue)
				stringValues := make([]string, len(values))
				for i, v := range values {
					stringValues[i] = fmt.Sprint(postgresValue(v))
				}
				conditions = append(conditions, fmt.Sprintf("%s @> $%d", column, *next))
				args = append(args, pq.Array(stringValues))
				*next++
			case "$exists":
				if exists, _ := opValue.(bool); exists {
					conditions = append(conditions, column+" IS NOT NULL")
				} else {
					conditions = append(conditions, column+" IS NULL")
				}
			case "$regex":
				pattern, insensitive := postgresRegex(opValue)
				symbol := "~"
				if insensitive {
					symbol = "~*"
				}
				conditions = append(conditions, fmt.Sprintf("%s %s $%d", column, symbol, *next))
				args = append(args, pattern)
				*next++
			default:
				// Failing closed is important for UPDATE and DELETE callers: an
				// unsupported Mongo operator must never broaden a SQL query.
				conditions = append(conditions, "FALSE")
			}
		}
	}
	return strings.Join(conditions, " AND "), args
}

func valueSlice(value interface{}) []interface{} {
	if value == nil {
		return nil
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return []interface{}{value}
	}
	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result
}

func isPostgresArrayColumn(table, column string) bool {
	return (table == "notes" && column == "tags") ||
		(table == "blog_comments" && column == "like_user_ids") ||
		(table == "user_blogs" && column == "cate_ids") ||
		(table == "tags" && column == "tags")
}

func postgresRegex(value interface{}) (string, bool) {
	switch regex := value.(type) {
	case bson.RegEx:
		return regex.Pattern, strings.Contains(regex.Options, "i")
	case *bson.RegEx:
		return regex.Pattern, strings.Contains(regex.Options, "i")
	default:
		return fmt.Sprint(value), false
	}
}

func scanPostgresStruct(rows *sql.Rows, target reflect.Value, targetType reflect.Type) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columns))
	destinations := make([]interface{}, len(columns))
	for i := range values {
		destinations[i] = &values[i]
	}
	if err := rows.Scan(destinations...); err != nil {
		return err
	}
	for i, column := range columns {
		field, ok := postgresFieldByColumn(target, targetType, column)
		if !ok || !field.CanSet() {
			continue
		}
		if err := setPostgresField(field, values[i]); err != nil {
			return fmt.Errorf("scan column %s: %w", column, err)
		}
	}
	return nil
}

func postgresFieldByColumn(value reflect.Value, typ reflect.Type, column string) (reflect.Value, bool) {
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := value.Field(i)
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			if nested, ok := postgresFieldByColumn(fieldValue, fieldType.Type, column); ok {
				return nested, true
			}
		}
		if getDBTag(fieldType) == column {
			return fieldValue, true
		}
	}
	return reflect.Value{}, false
}

func setPostgresField(field reflect.Value, raw interface{}) error {
	if raw == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}
	if field.Type() == bsonObjectIDType {
		id := strings.TrimSpace(rawString(raw))
		if id == "" {
			field.SetString("")
			return nil
		}
		if !bson.IsObjectIdHex(id) {
			return fmt.Errorf("invalid ObjectId %q", id)
		}
		field.Set(reflect.ValueOf(bson.ObjectIdHex(id)))
		return nil
	}
	if field.Type() == timeType {
		if value, ok := raw.(time.Time); ok {
			field.Set(reflect.ValueOf(value))
			return nil
		}
		return fmt.Errorf("expected time.Time, got %T", raw)
	}
	if field.Kind() == reflect.Slice {
		if field.Type().Elem().Kind() == reflect.String {
			var values pq.StringArray
			if err := values.Scan(raw); err != nil {
				return err
			}
			field.Set(reflect.ValueOf([]string(values)).Convert(field.Type()))
			return nil
		}
		return json.Unmarshal([]byte(rawString(raw)), field.Addr().Interface())
	}
	if field.Kind() == reflect.Map || (field.Kind() == reflect.Struct && field.Type() != timeType) {
		return json.Unmarshal([]byte(rawString(raw)), field.Addr().Interface())
	}
	if field.Kind() == reflect.String {
		field.SetString(rawString(raw))
		return nil
	}
	if field.Kind() == reflect.Bool {
		if value, ok := raw.(bool); ok {
			field.SetBool(value)
			return nil
		}
		value, err := strconv.ParseBool(rawString(raw))
		if err != nil {
			return err
		}
		field.SetBool(value)
		return nil
	}
	if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
		var value int64
		switch number := raw.(type) {
		case int64:
			value = number
		case int32:
			value = int64(number)
		case int:
			value = int64(number)
		default:
			parsed, err := strconv.ParseInt(rawString(raw), 10, 64)
			if err != nil {
				return err
			}
			value = parsed
		}
		field.SetInt(value)
		return nil
	}
	value := reflect.ValueOf(raw)
	if value.Type().AssignableTo(field.Type()) {
		field.Set(value)
		return nil
	}
	if value.Type().ConvertibleTo(field.Type()) {
		field.Set(value.Convert(field.Type()))
		return nil
	}
	return fmt.Errorf("cannot assign %T to %s", raw, field.Type())
}

func rawString(value interface{}) string {
	if bytes, ok := value.([]byte); ok {
		return string(bytes)
	}
	return fmt.Sprint(value)
}
