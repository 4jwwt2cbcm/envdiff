package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// ruleSpec is the JSON structure for a validation rules file.
type ruleSpec struct {
	Key      string `json:"key"`
	Pattern  string `json:"pattern,omitempty"`
	Required bool   `json:"required"`
}

func runValidate(envFile, rulesFile string) error {
	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	rules, err := loadRuleSpecs(rulesFile)
	if err != nil {
		return fmt.Errorf("loading rules file: %w", err)
	}

	result := diff.ValidateEnv(env, rules)
	if !result.HasErrors() {
		fmt.Println("validation passed: no issues found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "validation failed: %d issue(s) found\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", e.Key, e.Message)
	}
	return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
}

func loadRuleSpecs(path string) ([]diff.ValidationRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var specs []ruleSpec
	if err := json.Unmarshal(data, &specs); err != nil {
		return nil, fmt.Errorf("invalid rules JSON: %w", err)
	}

	var rules []diff.ValidationRule
	for _, s := range specs {
		rule := diff.ValidationRule{
			Key:      s.Key,
			Required: s.Required,
		}
		if s.Pattern != "" {
			re, err := regexp.Compile(s.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern for key %q: %w", s.Key, err)
			}
			rule.Pattern = re
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
