package diff

import (
	"sort"
	"strings"
)

// NormalizeOptions controls how env map normalization is performed.
type NormalizeOptions struct {
	TrimValues    bool
	LowercaseKeys bool
	UppercaseKeys bool
	RemoveEmpty   bool
}

// NormalizeResult holds the output of a normalization pass.
type NormalizeResult struct {
	Env      map[string]string
	Changed  []string // keys whose values or names were altered
	Removed  []string // keys that were dropped (e.g. empty value removal)
}

// NormalizeEnv applies the given options to env and returns a NormalizeResult.
// The original map is not modified.
func NormalizeEnv(env map[string]string, opts NormalizeOptions) NormalizeResult {
	result := NormalizeResult{
		Env: make(map[string]string, len(env)),
	}

	changedSet := map[string]struct{}{}

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.LowercaseKeys {
			newKey = strings.ToLower(k)
		} else if opts.UppercaseKeys {
			newKey = strings.ToUpper(k)
		}

		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
		}

		if opts.RemoveEmpty && newVal == "" {
			result.Removed = append(result.Removed, k)
			continue
		}

		if newKey != k || newVal != v {
			changedSet[k] = struct{}{}
		}

		result.Env[newKey] = newVal
	}

	for k := range changedSet {
		result.Changed = append(result.Changed, k)
	}

	sort.Strings(result.Changed)
	sort.Strings(result.Removed)

	return result
}
