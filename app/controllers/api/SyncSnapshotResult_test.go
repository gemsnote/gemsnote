package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/revel/revel"
)

func TestSnapshotGzipNegotiation(t *testing.T) {
	for _, tc := range []struct {
		header string
		want   bool
	}{
		{"", false}, {"br", false}, {"gzip", true}, {"br, gzip;q=0.5", true},
		{"gzip;q=0", false}, {"*;q=1, gzip;q=0", false}, {"*;q=0", false},
		{"*", true}, {"gzip;q=invalid", false}, {"gzip;q=2", false},
	} {
		if got := acceptsSnapshotGzip(tc.header); got != tc.want {
			t.Errorf("%q: gzip=%v want=%v", tc.header, got, tc.want)
		}
	}
}

func TestSnapshotResultWritesNegotiatedHTTPResponse(t *testing.T) {
	for _, accept := range []string{"gzip", "identity"} {
		request := httptest.NewRequest("GET", "/api2/note/getSyncNotesWithContent", nil)
		request.Header.Set("Accept-Encoding", accept)
		recorder := httptest.NewRecorder()
		context := revel.NewGoContext(nil)
		context.Request.SetRequest(request)
		context.Response.SetResponse(recorder)
		controller := revel.NewController(context)
		syncSnapshotResult{data: []map[string]string{{"Content": "body"}}}.Apply(controller.Request, controller.Response)
		if recorder.Code != 200 || recorder.Header().Get("Vary") != "Accept-Encoding" || !strings.Contains(recorder.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("incorrect snapshot response: %v", recorder)
		}
		body := recorder.Body.Bytes()
		if accept == "gzip" {
			if recorder.Header().Get("Content-Encoding") != "gzip" {
				t.Fatal("missing gzip header")
			}
			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			body, err = io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
		} else if recorder.Header().Get("Content-Encoding") != "" {
			t.Fatal("compressed response despite identity negotiation")
		}
		if string(body) != `[{"Content":"body"}]` {
			t.Fatalf("wrong snapshot: %s", body)
		}
	}
}

func TestSnapshotGzipPreservesBodiesAndEmptyContent(t *testing.T) {
	data := []map[string]interface{}{{"NoteId": "note1", "Content": strings.Repeat("珠玑笔记 <p>long note body</p>\n", 10000)}, {"NoteId": "empty", "Content": ""}}
	plain, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, accept := range []string{"", "gzip", "gzip;q=0"} {
		body, compressed, err := encodeSyncSnapshot(accept, data)
		if err != nil {
			t.Fatal(err)
		}
		if compressed {
			if len(body) >= len(plain)/4 {
				t.Fatalf("snapshot not effectively compressed: %d / %d bytes", len(body), len(plain))
			}
			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			body, err = io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
		}
		if !bytes.Equal(body, plain) {
			t.Fatalf("snapshot content changed with accept=%q", accept)
		}
	}
}
