package db

import "testing"

func TestEnvValue(t *testing.T) {
	t.Setenv("GEMSNOTE_TEST_ENV_VALUE", " secret ")
	if value, ok := envValue("GEMSNOTE_TEST_ENV_VALUE"); !ok || value != "secret" {
		t.Fatalf("envValue = %q, %v", value, ok)
	}
	t.Setenv("GEMSNOTE_TEST_ENV_EMPTY", "   ")
	if _, ok := envValue("GEMSNOTE_TEST_ENV_EMPTY"); ok {
		t.Fatal("blank environment value must not override configuration")
	}
	if _, ok := envValue("GEMSNOTE_TEST_ENV_MISSING"); ok {
		t.Fatal("missing environment value reported as set")
	}
}
