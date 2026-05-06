package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaFile(t *testing.T, entries []SchemaEntry) string {
	t.Helper()
	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	path := filepath.Join(t.TempDir(), "schema.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write schema file: %v", err)
	}
	return path
}

func TestLoadSchema_Basic(t *testing.T) {
	entries := []SchemaEntry{
		{Key: "DB_HOST", Required: true, Desc: "database host"},
		{Key: "API_KEY", Required: true, Secret: true},
	}
	path := writeSchemaFile(t, entries)

	schema, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema) != 2 {
		t.Errorf("expected 2 entries, got %d", len(schema))
	}
	if !schema["API_KEY"].Secret {
		t.Error("expected API_KEY to be secret")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSchema_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte(`{not valid json`), 0644)
	_, err := LoadSchema(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadSchema_MissingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schema.json")
	os.WriteFile(path, []byte(`[{"required":true}]`), 0644)
	_, err := LoadSchema(path)
	if err == nil {
		t.Error("expected error for entry missing key field")
	}
}

func TestSaveAndLoadSchema(t *testing.T) {
	original := Schema{
		"PORT":   {Key: "PORT", Required: false, Default: "8080"},
		"DB_URL": {Key: "DB_URL", Required: true, Secret: true},
	}
	path := filepath.Join(t.TempDir(), "schema.json")
	if err := SaveSchema(path, original); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	loaded, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("LoadSchema: %v", err)
	}
	if len(loaded) != len(original) {
		t.Errorf("expected %d entries, got %d", len(original), len(loaded))
	}
	if !loaded["DB_URL"].Secret {
		t.Error("expected DB_URL to be secret after round-trip")
	}
}

func TestGenerateSchema(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "3000",
	}
	schema := GenerateSchema(env)
	if len(schema) != 2 {
		t.Errorf("expected 2 schema entries, got %d", len(schema))
	}
	if schema["HOST"].Default != "localhost" {
		t.Errorf("expected default 'localhost', got %q", schema["HOST"].Default)
	}
	if !schema["PORT"].Required {
		t.Error("expected generated entries to be required by default")
	}
}
