package route

import "testing"

func TestLegacyAPIRewriteExcludesAPI2(t *testing.T) {
	for _, tc := range []struct {
		path string
		want bool
	}{
		{path: "/api/note/deleteNote", want: true},
		{path: "/api/system/version", want: true},
		{path: "/api2/note/deleteNote", want: false},
		{path: "/api2/notebook/addNotebook", want: false},
		{path: "/api2/attachments", want: false},
		{path: "/api", want: false},
	} {
		if got := isLegacyAPIPath(tc.path); got != tc.want {
			t.Errorf("isLegacyAPIPath(%q) = %t, want %t", tc.path, got, tc.want)
		}
	}
}
