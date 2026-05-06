package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleReport() Report {
	return Report{
		MissingInFirst:  []string{"KEY_A"},
		MissingInSecond: []string{"KEY_B", "KEY_C"},
		Mismatched:      []MismatchedKey{{Key: "DB_URL", ValueA: "x", ValueB: "y"}},
	}
}

func TestNewAuditEntry(t *testing.T) {
	r := sampleReport()
	before := time.Now().UTC()
	e := NewAuditEntry("a.env", "b.env", r)
	after := time.Now().UTC()

	if e.FileA != "a.env" || e.FileB != "b.env" {
		t.Errorf("unexpected file names: %s %s", e.FileA, e.FileB)
	}
	if e.MissingInA != 1 {
		t.Errorf("expected MissingInA=1, got %d", e.MissingInA)
	}
	if e.MissingInB != 2 {
		t.Errorf("expected MissingInB=2, got %d", e.MissingInB)
	}
	if e.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", e.Mismatched)
	}
	if !e.HasDrift {
		t.Error("expected HasDrift=true")
	}
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
}

func TestNewAuditEntry_NoDrift(t *testing.T) {
	r := Report{}
	e := NewAuditEntry("a.env", "b.env", r)
	if e.HasDrift {
		t.Error("expected HasDrift=false for empty report")
	}
}

func TestAppendAndLoadAuditLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	e1 := NewAuditEntry("a.env", "b.env", sampleReport())
	if err := AppendAuditLog(path, e1); err != nil {
		t.Fatalf("first append: %v", err)
	}

	e2 := NewAuditEntry("c.env", "d.env", Report{})
	if err := AppendAuditLog(path, e2); err != nil {
		t.Fatalf("second append: %v", err)
	}

	log, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].FileA != "a.env" {
		t.Errorf("unexpected first entry file: %s", log.Entries[0].FileA)
	}
}

func TestLoadAuditLog_MissingFile(t *testing.T) {
	_, err := LoadAuditLog("/nonexistent/audit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadAuditLog_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadAuditLog(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
