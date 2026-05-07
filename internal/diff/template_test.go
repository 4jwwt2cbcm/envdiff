package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTemplate_Basic(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"APP_PORT": "8080",
		"SECRET":  "abc123",
	}
	tmpl := GenerateTemplate(env)
	if len(tmpl.Vars) != 3 {
		t.Fatalf("expected 3 vars, got %d", len(tmpl.Vars))
	}
	// Keys should be sorted
	if tmpl.Vars[0].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT first, got %s", tmpl.Vars[0].Key)
	}
	if !tmpl.Vars[0].Required {
		t.Error("expected Required=true by default")
	}
	if tmpl.Vars[0].Example != "8080" {
		t.Errorf("expected example 8080, got %s", tmpl.Vars[0].Example)
	}
}

func TestSaveAndLoadTemplate(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "FOO", Description: "foo var", Required: true, Example: "bar"},
			{Key: "BAZ", Required: false, Default: "default_val"},
		},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "template.json")

	if err := SaveTemplate(path, tmpl); err != nil {
		t.Fatalf("SaveTemplate: %v", err)
	}
	loaded, err := LoadTemplate(path)
	if err != nil {
		t.Fatalf("LoadTemplate: %v", err)
	}
	if len(loaded.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(loaded.Vars))
	}
	if loaded.Vars[0].Key != "FOO" || loaded.Vars[1].Default != "default_val" {
		t.Error("loaded template content mismatch")
	}
}

func TestLoadTemplate_MissingFile(t *testing.T) {
	_, err := LoadTemplate("/nonexistent/template.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadTemplate_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadTemplate(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCheckAgainstTemplate_AllValid(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "HOST", Required: true},
			{Key: "PORT", Required: false},
		},
	}
	env := map[string]string{"HOST": "localhost", "PORT": "9090"}
	violations := CheckAgainstTemplate(env, tmpl)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckAgainstTemplate_MissingRequired(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "SECRET", Required: true},
		},
	}
	env := map[string]string{}
	violations := CheckAgainstTemplate(env, tmpl)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheckAgainstTemplate_EmptyRequiredValue(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "TOKEN", Required: true},
		},
	}
	env := map[string]string{"TOKEN": "   "}
	violations := CheckAgainstTemplate(env, tmpl)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation for empty required value, got %d", len(violations))
	}
}

func TestCheckAgainstTemplate_OptionalMissingOK(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "OPTIONAL_KEY", Required: false},
		},
	}
	env := map[string]string{}
	violations := CheckAgainstTemplate(env, tmpl)
	if len(violations) != 0 {
		t.Errorf("expected no violations for missing optional key, got %v", violations)
	}
}

func TestCheckAgainstTemplate_DefaultSkipsMissing(t *testing.T) {
	tmpl := Template{
		Vars: []TemplateVar{
			{Key: "WITH_DEFAULT", Required: true, Default: "fallback"},
		},
	}
	env := map[string]string{}
	violations := CheckAgainstTemplate(env, tmpl)
	if len(violations) != 0 {
		t.Errorf("expected no violations when default is set, got %v", violations)
	}
}

// ensure json round-trip preserves all fields
func TestTemplateVar_JSONRoundTrip(t *testing.T) {
	v := TemplateVar{Key: "X", Description: "desc", Required: true, Default: "d", Example: "e"}
	data, _ := json.Marshal(v)
	var v2 TemplateVar
	json.Unmarshal(data, &v2)
	if v != v2 {
		t.Errorf("round-trip mismatch: %+v vs %+v", v, v2)
	}
}
