package diff

import (
	"os"
	"testing"
)

func TestInterpolateEnv_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	results := InterpolateEnv(env)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestInterpolateEnv_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST":     "db.example.com",
		"DATABASE_URL": "postgres://${HOST}/mydb",
	}
	results := InterpolateEnv(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Key != "DATABASE_URL" {
		t.Errorf("unexpected key %q", r.Key)
	}
	if r.Resolved != "postgres://db.example.com/mydb" {
		t.Errorf("unexpected resolved value %q", r.Resolved)
	}
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", r.Missing)
	}
}

func TestInterpolateEnv_BareStyle(t *testing.T) {
	env := map[string]string{
		"USER":    "admin",
		"WELCOME": "Hello $USER",
	}
	results := InterpolateEnv(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Resolved != "Hello admin" {
		t.Errorf("got %q", results[0].Resolved)
	}
}

func TestInterpolateEnv_FallsBackToOS(t *testing.T) {
	os.Setenv("_ENVDIFF_TEST_VAR", "fromOS")
	defer os.Unsetenv("_ENVDIFF_TEST_VAR")

	env := map[string]string{
		"MSG": "value=${_ENVDIFF_TEST_VAR}",
	}
	results := InterpolateEnv(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Resolved != "value=fromOS" {
		t.Errorf("got %q", results[0].Resolved)
	}
}

func TestInterpolateEnv_MissingReference(t *testing.T) {
	env := map[string]string{
		"URL": "http://${UNKNOWN_HOST}/path",
	}
	results := InterpolateEnv(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.Missing) != 1 || r.Missing[0] != "UNKNOWN_HOST" {
		t.Errorf("expected missing UNKNOWN_HOST, got %v", r.Missing)
	}
}

func TestInterpolateMissing_Deduplicates(t *testing.T) {
	results := []InterpolateResult{
		{Key: "A", Missing: []string{"X", "Y"}},
		{Key: "B", Missing: []string{"Y", "Z"}},
	}
	missing := InterpolateMissing(results)
	if len(missing) != 3 {
		t.Fatalf("expected 3 unique missing keys, got %d: %v", len(missing), missing)
	}
}

func TestFormatInterpolateResult_WithMissing(t *testing.T) {
	r := InterpolateResult{
		Key:      "URL",
		Original: "http://${HOST}",
		Resolved: "http://${HOST}",
		Missing:  []string{"HOST"},
	}
	out := FormatInterpolateResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
