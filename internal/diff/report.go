package diff

import (
	"fmt"
	"io"
	"strings"
)

// Report writes a human-readable diff report to w.
// firstName and secondName are labels for the two env files being compared.
func Report(w io.Writer, r Result, firstName, secondName string) {
	if !r.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(r.MissingInSecond) > 0 {
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", firstName, secondName)
		for _, key := range r.MissingInSecond {
			fmt.Fprintf(w, "  - %s\n", key)
		}
	}

	if len(r.MissingInFirst) > 0 {
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", secondName, firstName)
		for _, key := range r.MissingInFirst {
			fmt.Fprintf(w, "  - %s\n", key)
		}
	}

	if len(r.Mismatched) > 0 {
		fmt.Fprintln(w, "Mismatched values:")
		for _, m := range r.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n", m.Key)
			fmt.Fprintf(w, "      %s: %s\n", firstName, maskValue(m.First))
			fmt.Fprintf(w, "      %s: %s\n", secondName, maskValue(m.Second))
		}
	}
}

// maskValue replaces the value with asterisks if it looks like a secret.
func maskValue(v string) string {
	lower := strings.ToLower(v)
	_ = lower
	// Placeholder: return value as-is (masking can be added later)
	return v
}
