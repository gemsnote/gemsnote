package service

import (
	"github.com/gemsnote/gemsnote/app/info"
	"testing"
)

func TestUsedNoteTagsWithoutIndex(t *testing.T) {
	tags := usedNoteTags([]info.Note{
		{Tags: []string{"new", "new", "", " ", "中文"}},
		{Tags: []string{"new", "100%"}},
	})
	counts := map[string]int{}
	for _, tag := range tags {
		counts[tag.Tag] = tag.Count
	}
	if len(counts) != 3 || counts["new"] != 2 || counts["中文"] != 1 || counts["100%"] != 1 {
		t.Fatalf("unexpected tag counts: %v", counts)
	}
	if tags := usedNoteTags(nil); tags == nil || len(tags) != 0 {
		t.Fatalf("want empty list: %v", tags)
	}
}
