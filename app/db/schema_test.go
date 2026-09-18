package db

import (
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/gemsnote/gemsnote/app/info"
)

func TestPostgresSchemaMatchesPersistedModels(t *testing.T) {
	schema, err := os.ReadFile("../../database/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	specs := map[string]interface{}{
		"users": info.User{}, "notebooks": info.Notebook{}, "notes": info.Note{},
		"note_contents": info.NoteContent{}, "note_content_histories": info.NoteContentHistory{},
		"share_notebooks": info.ShareNotebook{}, "share_notes": info.ShareNote{}, "has_share_notes": info.HasShareNote{},
		"groups": info.Group{}, "group_users": info.GroupUser{}, "tags": info.Tag{}, "note_tags": info.NoteTag{},
		"tag_count": info.TagCount{}, "blogs": info.BlogStat{}, "user_blogs": info.UserBlog{},
		"blog_singles": info.BlogSingle{}, "blog_likes": info.BlogLike{}, "blog_comments": info.BlogComment{},
		"reports": info.Report{}, "albums": info.Album{}, "files": info.File{}, "attachs": info.Attach{},
		"note_images": info.NoteImage{}, "configs": info.Config{}, "sessions": info.Session{},
		"tokens": info.Token{}, "email_logs": info.EmailLog{}, "suggestions": info.Suggestion{}, "themes": info.Theme{},
	}
	for table, model := range specs {
		tablePattern := regexp.MustCompile(`(?s)CREATE TABLE IF NOT EXISTS ` + regexp.QuoteMeta(table) + `\s*\((.*?)\);`)
		match := tablePattern.FindSubmatch(schema)
		if match == nil {
			t.Errorf("schema has no table %s", table)
			continue
		}
		for _, column := range getDBColumns(reflect.TypeOf(model)) {
			columnPattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(column) + `\s+`)
			if !columnPattern.Match(match[1]) {
				t.Errorf("table %s is missing model column %s", table, column)
			}
		}
	}
}
