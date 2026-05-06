package diff

import (
	"testing"
)

func TestShouldRedact_DefaultPatterns(t *testing.T) {
	opts := NewRedactOptions()

	sensitive := []string{
		"DB_PASSWORD",
		"API_SECRET",
		"AUTH_TOKEN",
		"STRIPE_API_KEY",
		"PRIVATE_KEY",
		"AWS_CREDENTIAL",
	}
	for _, key := range sensitive {
		if !opts.ShouldRedact(key) {
			t.Errorf("expected key %q to be redacted", key)
		}
	}

	safe := []string{
		"APP_ENV",
		"LOG_LEVEL",
		"PORT",
		"DATABASE_HOST",
	}
	for _, key := range safe {
		if opts.ShouldRedact(key) {
			t.Errorf("expected key %q NOT to be redacted", key)
		}
	}
}

func TestRedactEnv_ReplacesValues(t *testing.T) {
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_TOKEN":   "tok_abc123",
	}

	opts := NewRedactOptions()
	result := RedactEnv(env, opts)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("API_TOKEN should be redacted, got %q", result["API_TOKEN"])
	}
}

func TestRedactReport_RedactsMismatchedValues(t *testing.T) {
	r := Report{
		MissingInFirst:  []string{"GONE_KEY"},
		MissingInSecond: []string{"EXTRA_KEY"},
		Mismatched: []MismatchedKey{
			{Key: "DB_PASSWORD", Value1: "old", Value2: "new"},
			{Key: "APP_ENV", Value1: "staging", Value2: "production"},
		},
	}

	opts := NewRedactOptions()
	result := RedactReport(r, opts)

	if len(result.MissingInFirst) != 1 || result.MissingInFirst[0] != "GONE_KEY" {
		t.Errorf("MissingInFirst should be preserved")
	}

	for _, m := range result.Mismatched {
		if m.Key == "DB_PASSWORD" {
			if m.Value1 != "[REDACTED]" || m.Value2 != "[REDACTED]" {
				t.Errorf("DB_PASSWORD values should be redacted, got %q / %q", m.Value1, m.Value2)
			}
		}
		if m.Key == "APP_ENV" {
			if m.Value1 != "staging" || m.Value2 != "production" {
				t.Errorf("APP_ENV values should not be redacted")
			}
		}
	}
}

func TestRedactOptions_CustomPattern(t *testing.T) {
	opts := RedactOptions{Patterns: []string{"internal"}}
	if !opts.ShouldRedact("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to match custom pattern 'internal'")
	}
	if opts.ShouldRedact("DB_PASSWORD") {
		t.Error("DB_PASSWORD should not match custom pattern 'internal'")
	}
}
