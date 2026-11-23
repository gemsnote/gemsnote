package lea

import (
	"strings"
	"testing"
	"time"
)

func TestNeedsRehashDetectsLegacyMD5(t *testing.T) {
	if NeedsRehash(Md5("gemsnote")) != true {
		t.Fatal("32-char MD5 hash must be flagged for upgrade")
	}
	if NeedsRehash(GenPwd("gemsnote")) {
		t.Fatal("bcrypt hash must not be flagged for upgrade")
	}
	if NeedsRehash("") {
		t.Fatal("empty hash must not be flagged")
	}
}

func TestGenerateHashUsesStrongCost(t *testing.T) {
	hash := GenPwd("correct horse battery staple")
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("unexpected hash format %q", hash)
	}
	start := time.Now()
	if !ComparePwd("correct horse battery staple", hash) {
		t.Fatal("hash did not verify")
	}
	if ComparePwd("wrong password", hash) {
		t.Fatal("wrong password verified")
	}
	if elapsed := time.Since(start); elapsed < 5*time.Millisecond {
		t.Fatalf("bcrypt cost looks too low (%v for two comparisons)", elapsed)
	}
	// Legacy MD5 verification still works so the login upgrade path can run.
	legacy := Md5("legacy-password")
	if !ComparePwd("legacy-password", legacy) || ComparePwd("other", legacy) {
		t.Fatal("legacy MD5 verification broken")
	}
}

func TestPasswordValidationRequiresEightCharacters(t *testing.T) {
	InitVd()
	if ok, _ := Vd("password", "1234567"); ok {
		t.Fatal("7-character password accepted")
	}
	if ok, _ := Vd("password", "12345678"); !ok {
		t.Fatal("8-character password rejected")
	}
}
