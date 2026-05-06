package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuditEntry records a single diff operation for audit logging.
type AuditEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	FileA       string    `json:"file_a"`
	FileB       string    `json:"file_b"`
	MissingInA  int       `json:"missing_in_a"`
	MissingInB  int       `json:"missing_in_b"`
	Mismatched  int       `json:"mismatched"`
	HasDrift    bool      `json:"has_drift"`
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// NewAuditEntry creates an AuditEntry from a Report and file paths.
func NewAuditEntry(fileA, fileB string, r Report) AuditEntry {
	return AuditEntry{
		Timestamp:  time.Now().UTC(),
		FileA:      fileA,
		FileB:      fileB,
		MissingInA: len(r.MissingInFirst),
		MissingInB: len(r.MissingInSecond),
		Mismatched: len(r.Mismatched),
		HasDrift:   len(r.MissingInFirst)+len(r.MissingInSecond)+len(r.Mismatched) > 0,
	}
}

// AppendAuditLog loads an existing audit log from path (or starts fresh),
// appends the entry, and writes it back.
func AppendAuditLog(path string, entry AuditEntry) error {
	log, err := LoadAuditLog(path)
	if err != nil {
		log = &AuditLog{}
	}
	log.Entries = append(log.Entries, entry)

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadAuditLog reads and parses an audit log from the given path.
func LoadAuditLog(path string) (*AuditLog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return &log, nil
}
