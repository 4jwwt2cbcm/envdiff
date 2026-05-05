package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	keys := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"LOG_LEVEL": "info",
	}

	if err := SaveBaseline(path, "production", keys); err != nil {
		t.Fatalf("SaveBaseline error: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline error: %v", err)
	}

	if b.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", b.Environment)
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
	for k, v := range keys {
		if got := b.Keys[k]; got != v {
			t.Errorf("key %q: expected %q, got %q", k, v, got)
		}
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSaveBaseline_InvalidPath(t *testing.T) {
	err := SaveBaseline("/nonexistent/dir/baseline.json", "test", map[string]string{})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
