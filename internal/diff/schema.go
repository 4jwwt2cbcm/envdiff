package diff

import (
	"encoding/json"
	"fmt"
	"os"
)

// SchemaEntry defines the expected metadata for a single env key.
type SchemaEntry struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	Default  string `json:"default,omitempty"`
	Secret   bool   `json:"secret"`
	Desc     string `json:"desc,omitempty"`
}

// Schema is a collection of SchemaEntry items keyed by env variable name.
type Schema map[string]SchemaEntry

// LoadSchema reads a JSON schema file and returns a Schema.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read file: %w", err)
	}

	var entries []SchemaEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("schema: parse JSON: %w", err)
	}

	schema := make(Schema, len(entries))
	for _, e := range entries {
		if e.Key == "" {
			return nil, fmt.Errorf("schema: entry missing key field")
		}
		schema[e.Key] = e
	}
	return schema, nil
}

// SaveSchema writes a Schema to a JSON file.
func SaveSchema(path string, schema Schema) error {
	entries := make([]SchemaEntry, 0, len(schema))
	for _, e := range schema {
		entries = append(entries, e)
	}
	sortSchemaEntries(entries)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("schema: marshal JSON: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// GenerateSchema builds a Schema from an env map with sensible defaults.
func GenerateSchema(env map[string]string) Schema {
	schema := make(Schema, len(env))
	for k, v := range env {
		schema[k] = SchemaEntry{
			Key:      k,
			Required: true,
			Default:  v,
		}
	}
	return schema
}

// sortSchemaEntries sorts entries by key for deterministic output.
func sortSchemaEntries(entries []SchemaEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
