package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "value1" {
		t.Errorf("KEY1: got %q, want %q", env["KEY1"], "value1")
	}
	if env["KEY2"] != "value2" {
		t.Errorf("KEY2: got %q, want %q", env["KEY2"], "value2")
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=val\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DOUBLE="hello world"
SINGLE='foo bar'
`)
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DOUBLE"] != "hello world" {
		t.Errorf("DOUBLE: got %q", env["DOUBLE"])
	}
	if env["SINGLE"] != "foo bar" {
		t.Errorf("SINGLE: got %q", env["SINGLE"])
	}
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export KEY=exported\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "exported" {
		t.Errorf("KEY: got %q, want %q", env["KEY"], "exported")
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
