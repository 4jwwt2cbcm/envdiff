package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func writeProfileFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadProfiles_Basic(t *testing.T) {
	path := writeProfileFile(t, `[
		{"name": "backend", "keys": ["DB_HOST", "DB_PORT"]},
		{"name": "frontend", "keys": ["API_URL"]}
	]`)
	pm, err := LoadProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pm) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(pm))
	}
	if pm["backend"].Keys[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", pm["backend"].Keys[0])
	}
}

func TestLoadProfiles_MissingFile(t *testing.T) {
	pm, err := LoadProfiles("/nonexistent/profiles.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(pm) != 0 {
		t.Errorf("expected empty map, got %v", pm)
	}
}

func TestSaveAndLoadProfiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")
	pm := ProfileMap{
		"infra": {Name: "infra", Keys: []string{"AWS_KEY", "AWS_SECRET"}},
	}
	if err := SaveProfiles(path, pm); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadProfiles(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, ok := loaded["infra"]; !ok {
		t.Error("expected infra profile")
	}
}

func TestFilterByProfile_Basic(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "SECRET": "abc"}
	pm := ProfileMap{"db": {Name: "db", Keys: []string{"DB_HOST", "DB_PORT"}}}
	result, err := FilterByProfile(env, pm, "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["SECRET"]; ok {
		t.Error("SECRET should not be in result")
	}
}

func TestFilterByProfile_NotFound(t *testing.T) {
	_, err := FilterByProfile(map[string]string{}, ProfileMap{}, "missing")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestListProfileNames(t *testing.T) {
	pm := ProfileMap{
		"zebra":   {Name: "zebra"},
		"alpha":   {Name: "alpha"},
		"backend": {Name: "backend"},
	}
	names := ListProfileNames(pm)
	if names[0] != "alpha" || names[1] != "backend" || names[2] != "zebra" {
		t.Errorf("unexpected order: %v", names)
	}
}
