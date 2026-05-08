package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeScopeFile(t *testing.T, rules []ScopeRule) string {
	t.Helper()
	data, err := json.Marshal(rules)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "scopes.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadScopes_Basic(t *testing.T) {
	rules := []ScopeRule{
		{Name: "db", Prefixes: []string{"DB_"}, Keys: []string{"DATABASE_URL"}},
	}
	p := writeScopeFile(t, rules)
	loaded, err := LoadScopes(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Name != "db" {
		t.Errorf("expected 1 rule named 'db', got %+v", loaded)
	}
}

func TestLoadScopes_MissingFile(t *testing.T) {
	_, err := LoadScopes("/nonexistent/scopes.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoadScopes(t *testing.T) {
	rules := []ScopeRule{
		{Name: "auth", Prefixes: []string{"AUTH_", "JWT_"}, Keys: []string{}},
	}
	p := filepath.Join(t.TempDir(), "out.json")
	if err := SaveScopes(p, rules); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadScopes(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded[0].Name != "auth" {
		t.Errorf("expected auth scope, got %s", loaded[0].Name)
	}
}

func TestFilterByScope_PrefixAndKey(t *testing.T) {
	env := map[string]string{
		"DB_HOST":      "localhost",
		"DB_PORT":      "5432",
		"DATABASE_URL": "postgres://",
		"APP_NAME":     "myapp",
	}
	rule := ScopeRule{Name: "db", Prefixes: []string{"DB_"}, Keys: []string{"DATABASE_URL"}}
	out := FilterByScope(env, rule)
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(out), out)
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("APP_NAME should not be in scope")
	}
}

func TestCompareByScopes_DetectsMissing(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost", "AUTH_TOKEN": "abc"}
	b := map[string]string{"DB_HOST": "localhost"}
	rules := []ScopeRule{
		{Name: "auth", Prefixes: []string{"AUTH_"}, Keys: []string{}},
		{Name: "db", Prefixes: []string{"DB_"}, Keys: []string{}},
	}
	reports := CompareByScopes(a, b, rules)
	if len(reports) != 2 {
		t.Fatalf("expected 2 scope reports, got %d", len(reports))
	}
	// auth scope should show AUTH_TOKEN missing in second
	if reports[0].Scope != "auth" {
		t.Errorf("expected auth first (sorted), got %s", reports[0].Scope)
	}
	if len(reports[0].Report.MissingInSecond) != 1 {
		t.Errorf("expected 1 missing in second for auth scope, got %v", reports[0].Report.MissingInSecond)
	}
	// db scope should have no differences
	if len(reports[1].Report.MissingInSecond) != 0 || len(reports[1].Report.MissingInFirst) != 0 {
		t.Errorf("expected no diff in db scope, got %+v", reports[1].Report)
	}
}
