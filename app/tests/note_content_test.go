package tests

import (
	"testing"
	//  . "github.com/gemsnote/gemsnote/app/lea"
	"github.com/gemsnote/gemsnote/app/service"
	// "regexp"
	//  "gopkg.in/mgo.v2"
	// "fmt"
	// "strings"
)

// 可在server端调试

func TestApiFixNoteContent2(t *testing.T) {
	requireMongoIntegration(t)
	note2 := service.NoteS.GetNote("585df83771c1b17e8a000000", "585df81199c37b6176000004")
	note := service.NoteS.GetNoteContent("585df83771c1b17e8a000000", "585df81199c37b6176000004")
	contentFixed := service.NoteS.FixContent(note.Content, false)
	t.Log(note2.Title)
	t.Log(note.Content)
	t.Log(contentFixed)
}
