package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationRule defines a rule applied to env values.
type ValidationRule struct {
	Key     string
	Pattern *regexp.Regexp
	Required bool
}

// ValidationError represents a single validation failure.
type ValidationError struct {
	Key     string
	Value   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// ValidateEnv checks a parsed env map against a set of rules.
func ValidateEnv(env map[string]string, rules []ValidationRule) ValidationResult {
	result := ValidationResult{}

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			result.Errors = append(result.Errors, ValidationError{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if strings.TrimSpace(val) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     rule.Key,
				Value:   val,
				Message: "value is empty or blank",
			})
			continue
		}

		if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     rule.Key,
				Value:   val,
				Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern.String()),
			})
		}
	}

	return result
}
