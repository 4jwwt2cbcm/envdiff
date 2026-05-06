package diff

import (
	"testing"
)

func TestLintEnv_AllValid(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
		"_PRIVATE":     "secret",
	}
	violations := LintEnv(env)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestLintEnv_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost/db",
	}
	violations := LintEnv(env)
	if !containsRule(violations, "no-lowercase-key") {
		t.Errorf("expected no-lowercase-key violation, got: %+v", violations)
	}
}

func TestLintEnv_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	violations := LintEnv(env)
	if !containsRule(violations, "no-empty-value") {
		t.Errorf("expected no-empty-value violation, got: %+v", violations)
	}
}

func TestLintEnv_SpaceInKey(t *testing.T) {
	env := map[string]string{
		"BAD KEY": "value",
	}
	violations := LintEnv(env)
	if !containsRule(violations, "no-space-in-key") {
		t.Errorf("expected no-space-in-key violation, got: %+v", violations)
	}
}

func TestLintEnv_KeyMustStartWithLetter(t *testing.T) {
	env := map[string]string{
		"1INVALID": "value",
	}
	violations := LintEnv(env)
	if !containsRule(violations, "key-must-start-with-letter") {
		t.Errorf("expected key-must-start-with-letter violation, got: %+v", violations)
	}
}

func TestLintEnv_MultipleViolationsSorted(t *testing.T) {
	env := map[string]string{
		"bad_key": "",
		"GOOD_KEY": "value",
	}
	violations := LintEnv(env)
	// bad_key triggers lowercase + empty value
	if len(violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(violations))
	}
	// verify sorted order: all violations for bad_key should appear first
	for i := 1; i < len(violations); i++ {
		if violations[i-1].Key > violations[i].Key {
			t.Errorf("violations not sorted: %s > %s", violations[i-1].Key, violations[i].Key)
		}
	}
}

func containsRule(violations []LintViolation, rule string) bool {
	for _, v := range violations {
		if v.Rule == rule {
			return true
		}
	}
	return false
}
