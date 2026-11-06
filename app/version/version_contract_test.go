package version

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDeploymentImagesMatchApplicationVersion(t *testing.T) {
	t.Helper()
	root := filepath.Join("..", "..")
	files := []string{
		"docker-compose.yml",
		"docker-compose.postgres.yml",
		"docker-compose.mongodb.yml",
	}
	imageVersion := regexp.MustCompile(`gemsnote:(?:\$\{GEMSNOTE_VERSION:-)?([0-9]+\.[0-9]+\.[0-9]+)`) // optional Compose default wrapper
	for _, name := range files {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		matches := imageVersion.FindAllStringSubmatch(string(data), -1)
		if len(matches) == 0 {
			t.Fatalf("%s has no versioned Gemsnote image", name)
		}
		for _, match := range matches {
			if match[1] != Current {
				t.Errorf("%s image version = %q, application version = %q", name, match[1], Current)
			}
		}
	}
}

func TestReleaseImageReceivesCanonicalVersion(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(data)
	for _, required := range []string{
		"app/version/VERSION",
		"GEMSNOTE_VERSION=${{ needs.validate.outputs.version }}",
		"org.opencontainers.image.version=${{ needs.validate.outputs.version }}",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release workflow does not derive image metadata from the canonical version: missing %q", required)
		}
	}
}
