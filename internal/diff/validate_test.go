package diff

import (
	"regexp"
	"testing"
)

func TestValidateEnv_AllValid(t *testing.T) {
	env := map[string]string{
		"PORT":     "8080",
		"APP_NAME": "myapp",
	}
	rules := []ValidationRule{
		{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		{Key: "APP_NAME", Required: true},
	}
	result := ValidateEnv(env, rules)
	if result.HasErrors() {
		t.Errorf("expected no errors, got: %v", result.Errors)
	}
}

func TestValidateEnv_MissingRequired(t *testing.T) {
	env := map[string]string{}
	rules := []ValidationRule{
		{Key: "DATABASE_URL", Required: true},
	}
	result := ValidateEnv(env, rules)
	if !result.HasErrors() {
		t.Fatal("expected errors for missing required key")
	}
	if result.Errors[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected error key: %s", result.Errors[0].Key)
	}
}

func TestValidateEnv_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"PORT": "not-a-number",
	}
	rules := []ValidationRule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	result := ValidateEnv(env, rules)
	if !result.HasErrors() {
		t.Fatal("expected pattern mismatch error")
	}
	if result.Errors[0].Key != "PORT" {
		t.Errorf("unexpected error key: %s", result.Errors[0].Key)
	}
}

func TestValidateEnv_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SECRET": "   ",
	}
	rules := []ValidationRule{
		{Key: "SECRET", Required: false},
	}
	result := ValidateEnv(env, rules)
	if !result.HasErrors() {
		t.Fatal("expected error for blank value")
	}
	if result.Errors[0].Message != "value is empty or blank" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidateEnv_OptionalMissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	rules := []ValidationRule{
		{Key: "OPTIONAL_KEY", Required: false, Pattern: regexp.MustCompile(`^\d+$`)},
	}
	result := ValidateEnv(env, rules)
	if result.HasErrors() {
		t.Errorf("expected no errors for missing optional key, got: %v", result.Errors)
	}
}

func TestValidateEnv_MultipleErrors(t *testing.T) {
	env := map[string]string{
		"PORT": "abc",
	}
	rules := []ValidationRule{
		{Key: "PORT", Pattern: regexp.MustCompile(`^\d+$`)},
		{Key: "HOST", Required: true},
	}
	result := ValidateEnv(env, rules)
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}
