package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of a diff result for future comparison.
type Baseline struct {
	CreatedAt   time.Time         `json:"created_at"`
	Environment string            `json:"environment"`
	Keys        map[string]string `json:"keys"`
}

// SaveBaseline writes a baseline snapshot of the given env map to a file.
func SaveBaseline(path, environment string, keys map[string]string) error {
	b := Baseline{
		CreatedAt:   time.Now().UTC(),
		Environment: environment,
		Keys:        keys,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write baseline: %w", err)
	}
	return nil
}

// LoadBaseline reads a baseline snapshot from a file.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read baseline: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshal baseline: %w", err)
	}
	return &b, nil
}
