package controllers

import (
	"github.com/gemsnote/gemsnote/app/info"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestNotebookOrderValidatesAllTargetsBeforeWriting(t *testing.T) {
	user, other, a, b, parent := bson.NewObjectId(), bson.NewObjectId(), bson.NewObjectId(), bson.NewObjectId(), bson.NewObjectId()
	books := map[string]info.Notebook{a.Hex(): {NotebookId: a, UserId: user}, b.Hex(): {NotebookId: b, UserId: user}}
	lookup := func(id string) info.Notebook { return books[id] }
	order, msg := validateNotebookOrder(user.Hex(), []string{b.Hex(), a.Hex()}, lookup)
	if msg != "" || order[b.Hex()] != 0 || order[a.Hex()] != 1 {
		t.Fatalf("invalid order: %+v %s", order, msg)
	}
	for _, ids := range [][]string{nil, {a.Hex(), a.Hex()}, {"invalid"}} {
		if _, msg := validateNotebookOrder(user.Hex(), ids, lookup); msg == "" {
			t.Fatalf("accepted invalid IDs: %v", ids)
		}
	}
	nb := books[b.Hex()]
	nb.UserId = other
	books[b.Hex()] = nb
	if _, msg = validateNotebookOrder(user.Hex(), []string{a.Hex(), b.Hex()}, lookup); msg != "noAuth" {
		t.Fatal(msg)
	}
	nb.UserId = user
	nb.ParentNotebookId = parent
	books[b.Hex()] = nb
	if _, msg = validateNotebookOrder(user.Hex(), []string{a.Hex(), b.Hex()}, lookup); msg != "invalidParentNotebook" {
		t.Fatal(msg)
	}
	nb.ParentNotebookId = ""
	nb.IsDeleted = true
	books[b.Hex()] = nb
	if _, msg = validateNotebookOrder(user.Hex(), []string{a.Hex(), b.Hex()}, lookup); msg != "noAuth" {
		t.Fatal(msg)
	}
}
