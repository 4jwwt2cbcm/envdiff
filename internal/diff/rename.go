package diff

import (
	"encoding/json"
	"fmt"
	"os"
)

// RenameMap holds a mapping from old key names to new key names.
type RenameMap map[string]string

// Inverse returns a new RenameMap with keys and values swapped.
func (r RenameMap) Inverse() RenameMap {
	inv := make(RenameMap, len(r))
	for oldKey, newKey := range r {
		inv[newKey] = oldKey
	}
	return inv
}

// Validate checks the RenameMap for duplicate target keys, which would cause
// ambiguous renames. It returns an error listing the first duplicate found.
func (r RenameMap) Validate() error {
	seen := make(map[string]string, len(r))
	for oldKey, newKey := range r {
		if prev, ok := seen[newKey]; ok {
			return fmt.Errorf("rename: duplicate target key %q (mapped from both %q and %q)", newKey, prev, oldKey)
		}
		seen[newKey] = oldKey
	}
	return nil
}

// ApplyRenames rewrites the keys in an env map according to the rename map.
// Keys not present in the rename map are left unchanged.
func ApplyRenames(env map[string]string, renames RenameMap) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if newKey, ok := renames[k]; ok {
			result[newKey] = v
		} else {
			result[k] = v
		}
	}
	return result
}

// LoadRenameFile reads a JSON file mapping old key names to new key names.
// The file format is: {"OLD_KEY": "NEW_KEY", ...}
func LoadRenameFile(path string) (RenameMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("rename: read file %q: %w", path, err)
	}
	var rm RenameMap
	if err := json.Unmarshal(data, &rm); err != nil {
		return nil, fmt.Errorf("rename: parse file %q: %w", path, err)
	}
	if err := rm.Validate(); err != nil {
		return nil, err
	}
	return rm, nil
}

// SaveRenameFile writes a RenameMap to a JSON file at the given path.
func SaveRenameFile(path string, rm RenameMap) error {
	data, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return fmt.Errorf("rename: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("rename: write file %q: %w", path, err)
	}
	return nil
}
