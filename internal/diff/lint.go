package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// LintRule defines a single linting rule for env keys or values.
type LintRule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// LintViolation represents a single rule violation.
type LintViolation struct {
	Key     string
	Rule    string
	Message string
}

var defaultLintRules = []LintRule{
	{
		Name:    "no-lowercase-key",
		Message: "env key should be uppercase",
		Check: func(key, _ string) bool {
			return key != strings.ToUpper(key)
		},
	},
	{
		Name:    "no-space-in-key",
		Message: "env key must not contain spaces",
		Check: func(key, _ string) bool {
			return strings.Contains(key, " ")
		},
	},
	{
		Name:    "no-empty-value",
		Message: "env value must not be empty",
		Check: func(_, value string) bool {
			return strings.TrimSpace(value) == ""
		},
	},
	{
		Name:    "key-must-start-with-letter",
		Message: "env key must start with a letter or underscore",
		Check: func(key, _ string) bool {
			if len(key) == 0 {
				return true
			}
			matched, _ := regexp.MatchString(`^[A-Za-z_]`, key)
			return !matched
		},
	},
}

// LintEnv runs all default lint rules against the provided env map.
// Returns a slice of violations (empty if all pass).
func LintEnv(env map[string]string) []LintViolation {
	var violations []LintViolation
	for key, value := range env {
		for _, rule := range defaultLintRules {
			if rule.Check(key, value) {
				violations = append(violations, LintViolation{
					Key:     key,
					Rule:    rule.Name,
					Message: fmt.Sprintf("%s: %s", key, rule.Message),
				})
			}
		}
	}
	sortLintViolations(violations)
	return violations
}

// sortLintViolations sorts violations by key then rule name for deterministic output.
func sortLintViolations(violations []LintViolation) {
	for i := 1; i < len(violations); i++ {
		for j := i; j > 0; j-- {
			a, b := violations[j-1], violations[j]
			if a.Key > b.Key || (a.Key == b.Key && a.Rule > b.Rule) {
				violations[j-1], violations[j] = violations[j], violations[j-1]
			}
		}
	}
}
