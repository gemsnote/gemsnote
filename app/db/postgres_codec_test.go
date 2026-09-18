package db

import (
	"database/sql/driver"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestPostgresNilTagsBecomeEmptyArray(t *testing.T) {
	value, ok := postgresValue([]string(nil)).(driver.Valuer)
	if !ok {
		t.Fatal("nil tag slice was not converted to a PostgreSQL array")
	}
	encoded, err := value.Value()
	if err != nil || encoded != "{}" {
		t.Fatalf("nil tags encoded as %v, err %v", encoded, err)
	}
}

func TestPostgresUsesMongoObjectIDs(t *testing.T) {
	id := bson.NewObjectId()
	if got := postgresValue(id); got != id.Hex() {
		t.Fatalf("postgresValue(ObjectId) = %v, want %s", got, id.Hex())
	}
	p := &PostgresDatabase{}
	if !p.IsValidID(id.Hex()) || p.IsValidID("550e8400-e29b-41d4-a716-446655440000") {
		t.Fatal("PostgreSQL ID validation must use MongoDB ObjectId syntax")
	}
}

func TestBuildPostgresWhereSupportsServiceOperators(t *testing.T) {
	id1, id2 := bson.NewObjectId(), bson.NewObjectId()
	query := bson.M{
		"UserId": id1,
		"_id":    bson.M{"$in": []bson.ObjectId{id1, id2}},
		"$or": []bson.M{
			{"Title": bson.M{"$regex": bson.RegEx{Pattern: "hello", Options: "i"}}},
			{"IsDeleted": bson.M{"$exists": false}},
		},
	}
	clause, args := buildPostgresWhere("notes", query, 1)
	for _, expected := range []string{"user_id =", "id IN", " OR ", "title ~*", "is_deleted IS NULL"} {
		if !strings.Contains(clause, expected) {
			t.Fatalf("where clause %q does not contain %q", clause, expected)
		}
	}
	if len(args) != 4 {
		t.Fatalf("got %d args, want 4: %#v", len(args), args)
	}
}

func TestUnsupportedWhereOperatorFailsClosed(t *testing.T) {
	clause, _ := buildPostgresWhere("notes", bson.M{"Title": bson.M{"$unknown": "x"}}, 1)
	if clause != "FALSE" {
		t.Fatalf("unsupported operator produced %q, want FALSE", clause)
	}
}

func TestBuildPostgresSetSupportsAtomicUpdates(t *testing.T) {
	parts, args := buildPostgresSet(bson.M{
		"$inc":  bson.M{"ReadNum": 1},
		"$push": bson.M{"LikeUserIds": "abc"},
	}, 1)
	joined := strings.Join(parts, " ")
	if !strings.Contains(joined, "read_num = COALESCE(read_num, 0) + $1") ||
		!strings.Contains(joined, "array_append") || len(args) != 2 {
		t.Fatalf("unexpected update: %q %#v", joined, args)
	}
}
