package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func writeIgnoreFile(t *testing.T, lines string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(p, []byte(lines), 0o644); err != nil {
		t.Fatalf("failed to write ignore file: %v", err)
	}
	return p
}

func TestIgnoreList_AddAndContains(t *testing.T) {
	il := NewIgnoreList()
	il.Add("SECRET_KEY")
	il.Add("  DB_PASS  ") // should be trimmed

	if !il.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be in ignore list")
	}
	if !il.Contains("DB_PASS") {
		t.Error("expected DB_PASS (trimmed) to be in ignore list")
	}
	if il.Contains("OTHER_KEY") {
		t.Error("OTHER_KEY should not be in ignore list")
	}
	if il.Len() != 2 {
		t.Errorf("expected length 2, got %d", il.Len())
	}
}

func TestLoadIgnoreFile_Basic(t *testing.T) {
	p := writeIgnoreFile(t, "# comment\nSECRET_KEY\n\nDB_PASS\n")
	il, err := LoadIgnoreFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if il.Len() != 2 {
		t.Errorf("expected 2 keys, got %d", il.Len())
	}
	if !il.Contains("SECRET_KEY") || !il.Contains("DB_PASS") {
		t.Error("expected both keys to be present")
	}
}

func TestLoadIgnoreFile_Missing(t *testing.T) {
	_, err := LoadIgnoreFile("/nonexistent/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestApplyIgnoreList_RemovesKeys(t *testing.T) {
	r := Result{
		MissingInFirst:  []string{"A", "SECRET_KEY"},
		MissingInSecond: []string{"B", "DB_PASS"},
		Mismatched:      []Mismatch{{Key: "C", Value1: "x", Value2: "y"}, {Key: "SECRET_KEY", Value1: "1", Value2: "2"}},
	}

	il := NewIgnoreList()
	il.Add("SECRET_KEY")
	il.Add("DB_PASS")

	out := ApplyIgnoreList(r, il)

	if len(out.MissingInFirst) != 1 || out.MissingInFirst[0] != "A" {
		t.Errorf("unexpected MissingInFirst: %v", out.MissingInFirst)
	}
	if len(out.MissingInSecond) != 1 || out.MissingInSecond[0] != "B" {
		t.Errorf("unexpected MissingInSecond: %v", out.MissingInSecond)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "C" {
		t.Errorf("unexpected Mismatched: %v", out.Mismatched)
	}
}

func TestApplyIgnoreList_NilIgnoreList(t *testing.T) {
	r := Result{
		MissingInFirst: []string{"A"},
	}
	out := ApplyIgnoreList(r, nil)
	if len(out.MissingInFirst) != 1 {
		t.Error("nil ignore list should not remove any keys")
	}
}
