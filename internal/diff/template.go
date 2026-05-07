package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// TemplateVar represents a single variable definition in a template.
type TemplateVar struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Example     string `json:"example,omitempty"`
}

// Template is an ordered collection of variable definitions.
type Template struct {
	Vars []TemplateVar `json:"vars"`
}

// GenerateTemplate creates a Template from a parsed env map.
// Keys are sorted alphabetically. All vars are marked required by default.
func GenerateTemplate(env map[string]string) Template {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	vars := make([]TemplateVar, 0, len(keys))
	for _, k := range keys {
		vars = append(vars, TemplateVar{
			Key:      k,
			Example:  env[k],
			Required: true,
		})
	}
	return Template{Vars: vars}
}

// SaveTemplate writes a Template to a JSON file at path.
func SaveTemplate(path string, t Template) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal template: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadTemplate reads a Template from a JSON file at path.
func LoadTemplate(path string) (Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Template{}, fmt.Errorf("template file not found: %s", path)
		}
		return Template{}, fmt.Errorf("read template: %w", err)
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return Template{}, fmt.Errorf("parse template: %w", err)
	}
	return t, nil
}

// CheckAgainstTemplate validates an env map against a Template.
// Returns a list of human-readable violation strings.
func CheckAgainstTemplate(env map[string]string, t Template) []string {
	var violations []string
	for _, v := range t.Vars {
		val, ok := env[v.Key]
		if !ok {
			if v.Required && v.Default == "" {
				violations = append(violations, fmt.Sprintf("missing required key: %s", v.Key))
			}
			continue
		}
		if v.Required && strings.TrimSpace(val) == "" {
			violations = append(violations, fmt.Sprintf("required key has empty value: %s", v.Key))
		}
	}
	return violations
}
