package diff

import (
	"regexp"
	"strings"
)

// RedactOptions controls which keys have their values redacted in output.
type RedactOptions struct {
	// Patterns is a list of glob-style substrings; keys containing any of these
	// (case-insensitive) will have their values replaced with [REDACTED].
	Patterns []string
}

// DefaultRedactPatterns contains common sensitive key name fragments.
var DefaultRedactPatterns = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
	"auth",
}

// NewRedactOptions returns RedactOptions using the default sensitive patterns.
func NewRedactOptions() RedactOptions {
	return RedactOptions{Patterns: DefaultRedactPatterns}
}

// ShouldRedact reports whether the given key matches any redaction pattern.
func (r RedactOptions) ShouldRedact(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.Patterns {
		matched, err := regexp.MatchString(strings.ToLower(p), lower)
		if err == nil && matched {
			return true
		}
		// fallback: plain substring match
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// RedactEnv returns a copy of env with sensitive values replaced by [REDACTED].
func RedactEnv(env map[string]string, opts RedactOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.ShouldRedact(k) {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactReport returns a copy of the Report with sensitive mismatch values redacted.
func RedactReport(r Report, opts RedactOptions) Report {
	redacted := Report{
		MissingInFirst:  r.MissingInFirst,
		MissingInSecond: r.MissingInSecond,
	}
	for _, m := range r.Mismatched {
		if opts.ShouldRedact(m.Key) {
			m.Value1 = "[REDACTED]"
			m.Value2 = "[REDACTED]"
		}
		redacted.Mismatched = append(redacted.Mismatched, m)
	}
	return redacted
}
