package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyRenames_Basic(t *testing.T) {
	env := map[string]string{
		"OLD_HOST": "localhost",
		"OLD_PORT": "5432",
		"UNCHANGED": "value",
	}
	renames := RenameMap{
		"OLD_HOST": "DB_HOST",
		"OLD_PORT": "DB_PORT",
	}
	result := ApplyRenames(env, renames)

	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", result["DB_PORT"])
	}
	if result["UNCHANGED"] != "value" {
		t.Errorf("expected UNCHANGED=value, got %q", result["UNCHANGED"])
	}
	if _, ok := result["OLD_HOST"]; ok {
		t.Error("expected OLD_HOST to be removed after rename")
	}
}

func TestApplyRenames_Empty(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	result := ApplyRenames(env, RenameMap{})
	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", result["KEY"])
	}
}

func TestRenameMap_Inverse(t *testing.T) {
	rm := RenameMap{"OLD": "NEW", "FOO": "BAR"}
	inv := rm.Inverse()
	if inv["NEW"] != "OLD" {
		t.Errorf("expected NEW->OLD, got %q", inv["NEW"])
	}
	if inv["BAR"] != "FOO" {
		t.Errorf("expected BAR->FOO, got %q", inv["BAR"])
	}
}

func TestSaveAndLoadRenameFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "renames.json")

	rm := RenameMap{"OLD_KEY": "NEW_KEY", "LEGACY_URL": "APP_URL"}
	if err := SaveRenameFile(path, rm); err != nil {
		t.Fatalf("SaveRenameFile: %v", err)
	}

	loaded, err := LoadRenameFile(path)
	if err != nil {
		t.Fatalf("LoadRenameFile: %v", err)
	}
	if loaded["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected OLD_KEY->NEW_KEY, got %q", loaded["OLD_KEY"])
	}
	if loaded["LEGACY_URL"] != "APP_URL" {
		t.Errorf("expected LEGACY_URL->APP_URL, got %q", loaded["LEGACY_URL"])
	}
}

func TestLoadRenameFile_Missing(t *testing.T) {
	_, err := LoadRenameFile("/nonexistent/renames.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRenameFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := LoadRenameFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	_ = json.Unmarshal // suppress unused import if needed
}
