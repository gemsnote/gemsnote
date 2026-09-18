package parser2

import "testing"

func TestShouldSkipSourceDir(t *testing.T) {
	root := "/work/gemsnote"
	for _, path := range []string{
		root + "/data/18/docker",
		root + "/docs/reference",
		root + "/tmp/generated",
		root + "/files/uploads",
	} {
		if !shouldSkipSourceDir(root, path) {
			t.Errorf("should skip %s", path)
		}
	}
	for _, path := range []string{root, root + "/app", root + "/frontend/src", root + "/public"} {
		if shouldSkipSourceDir(root, path) {
			t.Errorf("should not skip %s", path)
		}
	}
}
