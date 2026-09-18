package route

import "testing"

func TestRemovedWebPaths(t *testing.T) {
	for _, path := range []string{"/blog", "/blog/post/id", "/preview/index", "/member/blog/theme", "/memberblog/updateTheme", "/adminblog/index"} {
		if !removedWebPath(path) { t.Errorf("publishing route still enabled: %s", path) }
	}
	for _, path := range []string{"/api/note/updateNote", "/share/addShareNote", "/member/user/password", "/images/blog/default_avatar.png", "/blogger"} {
		if removedWebPath(path) { t.Errorf("unrelated route disabled: %s", path) }
	}
}
