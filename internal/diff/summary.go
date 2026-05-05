package diff

import "fmt"

// Summary holds aggregated statistics about a diff result.
type Summary struct {
	TotalKeys       int
	MissingInFirst  int
	MissingInSecond int
	Mismatched      int
}

// Summarize computes a Summary from a Result.
func Summarize(r Result) Summary {
	allKeys := make(map[string]struct{})
	for _, k := range r.MissingInFirst {
		allKeys[k] = struct{}{}
	}
	for _, k := range r.MissingInSecond {
		allKeys[k] = struct{}{}
	}
	for _, m := range r.Mismatched {
		allKeys[m.Key] = struct{}{}
	}
	// Also count keys that matched (present in both with same value).
	for _, k := range r.Matching {
		allKeys[k] = struct{}{}
	}

	return Summary{
		TotalKeys:       len(allKeys),
		MissingInFirst:  len(r.MissingInFirst),
		MissingInSecond: len(r.MissingInSecond),
		Mismatched:      len(r.Mismatched),
	}
}

// HasDifferences returns true when any discrepancy exists.
func (s Summary) HasDifferences() bool {
	return s.MissingInFirst > 0 || s.MissingInSecond > 0 || s.Mismatched > 0
}

// String returns a human-readable one-line summary.
func (s Summary) String() string {
	if !s.HasDifferences() {
		return fmt.Sprintf("No differences found (%d keys compared).", s.TotalKeys)
	}
	return fmt.Sprintf(
		"%d key(s) compared: %d missing in first, %d missing in second, %d mismatched.",
		s.TotalKeys, s.MissingInFirst, s.MissingInSecond, s.Mismatched,
	)
}
