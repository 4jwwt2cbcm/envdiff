package diff

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateResult holds the outcome of interpolating a single value.
type InterpolateResult struct {
	Key      string
	Original string
	Resolved string
	Missing  []string
}

// InterpolateEnv resolves variable references within env values using the
// provided env map as the source of substitutions. Falls back to OS
// environment when a key is absent from the map.
func InterpolateEnv(env map[string]string) []InterpolateResult {
	results := make([]InterpolateResult, 0, len(env))

	for key, value := range env {
		if !strings.Contains(value, "$") {
			continue
		}

		missing := []string{}
		resolved := interpolatePattern.ReplaceAllStringFunc(value, func(match string) string {
			refKey := extractKey(match)
			if v, ok := env[refKey]; ok {
				return v
			}
			if v, ok := os.LookupEnv(refKey); ok {
				return v
			}
			missing = append(missing, refKey)
			return match
		})

		results = append(results, InterpolateResult{
			Key:      key,
			Original: value,
			Resolved: resolved,
			Missing:  missing,
		})
	}

	return results
}

// InterpolateMissing returns the set of keys that could not be resolved
// across all interpolation results.
func InterpolateMissing(results []InterpolateResult) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, r := range results {
		for _, m := range r.Missing {
			if _, ok := seen[m]; !ok {
				seen[m] = struct{}{}
				out = append(out, m)
			}
		}
	}
	sortStrings(out)
	return out
}

// extractKey strips the ${ } wrapper or leading $ from a match.
func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}

// FormatInterpolateResult returns a human-readable description of the result.
func FormatInterpolateResult(r InterpolateResult) string {
	if len(r.Missing) > 0 {
		return fmt.Sprintf("%s: %q => %q (unresolved: %s)",
			r.Key, r.Original, r.Resolved, strings.Join(r.Missing, ", "))
	}
	return fmt.Sprintf("%s: %q => %q", r.Key, r.Original, r.Resolved)
}
