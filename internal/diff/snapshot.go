package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// NewSnapshot creates a snapshot from an env map with a label.
func NewSnapshot(label string, env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Env:       copy,
	}
}

// SaveSnapshot writes a snapshot to a JSON file.
func SaveSnapshot(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("snapshot: file not found: %s", path)
		}
		return Snapshot{}, fmt.Errorf("snapshot: read: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return s, nil
}

// DiffSnapshot compares a current env map against a saved snapshot and
// returns a Result describing what has changed since the snapshot was taken.
func DiffSnapshot(current map[string]string, s Snapshot) Result {
	return Compare(s.Env, current)
}
