package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// ScopeRule defines which keys belong to a named scope (e.g. "database", "auth").
type ScopeRule struct {
	Name    string   `json:"name"`
	Prefixes []string `json:"prefixes"`
	Keys    []string `json:"keys"`
}

// ScopeReport holds per-scope diff results.
type ScopeReport struct {
	Scope   string
	Report  Report
}

// LoadScopes reads a JSON file containing a list of ScopeRule entries.
func LoadScopes(path string) ([]ScopeRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("scope: read file: %w", err)
	}
	var rules []ScopeRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("scope: parse JSON: %w", err)
	}
	return rules, nil
}

// SaveScopes writes scope rules to a JSON file.
func SaveScopes(path string, rules []ScopeRule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// FilterByScope returns only the env entries that match the given ScopeRule.
func FilterByScope(env map[string]string, rule ScopeRule) map[string]string {
	out := make(map[string]string)
	keySet := make(map[string]bool, len(rule.Keys))
	for _, k := range rule.Keys {
		keySet[k] = true
	}
	for k, v := range env {
		if keySet[k] {
			out[k] = v
			continue
		}
		for _, p := range rule.Prefixes {
			if len(k) >= len(p) && k[:len(p)] == p {
				out[k] = v
				break
			}
		}
	}
	return out
}

// CompareByScopes runs Compare for each scope and returns a slice of ScopeReports.
func CompareByScopes(a, b map[string]string, rules []ScopeRule) []ScopeReport {
	results := make([]ScopeReport, 0, len(rules))
	sort.Slice(rules, func(i, j int) bool { return rules[i].Name < rules[j].Name })
	for _, rule := range rules {
		sa := FilterByScope(a, rule)
		sb := FilterByScope(b, rule)
		report := Compare(sa, sb)
		results = append(results, ScopeReport{Scope: rule.Name, Report: report})
	}
	return results
}
