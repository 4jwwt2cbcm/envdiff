package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshot_CopiesEnv(t *testing.T) {
	env := map[string]string{"KEY": "value", "PORT": "8080"}
	snap := NewSnapshot("test", env)

	// Mutate original — snapshot should be unaffected.
	env["KEY"] = "changed"

	if snap.Env["KEY"] != "value" {
		t.Errorf("expected snapshot KEY=value, got %s", snap.Env["KEY"])
	}
	if snap.Label != "test" {
		t.Errorf("expected label 'test', got %s", snap.Label)
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{"DB_HOST": "localhost", "DEBUG": "true"}
	orig := NewSnapshot("prod", env)

	if err := SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if loaded.Label != "prod" {
		t.Errorf("label mismatch: got %s", loaded.Label)
	}
	if loaded.Env["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST mismatch: got %s", loaded.Env["DB_HOST"])
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	origEnv := map[string]string{"KEY": "old", "ONLY_IN_SNAP": "yes"}
	snap := NewSnapshot("base", origEnv)

	currentEnv := map[string]string{"KEY": "new", "ONLY_IN_CURRENT": "yes"}
	result := DiffSnapshot(currentEnv, snap)

	if len(result.Mismatched) != 1 || result.Mismatched[0].Key != "KEY" {
		t.Errorf("expected KEY mismatch, got %+v", result.Mismatched)
	}
	if len(result.MissingInSecond) != 1 || result.MissingInSecond[0] != "ONLY_IN_CURRENT" {
		t.Errorf("expected ONLY_IN_CURRENT missing in second, got %+v", result.MissingInSecond)
	}
}
