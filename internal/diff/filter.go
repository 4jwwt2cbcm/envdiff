package diff

// FilterOptions controls which differences are included in results.
type FilterOptions struct {
	// OnlyMissing limits results to keys missing in either file.
	OnlyMissing bool
	// OnlyMismatched limits results to keys present in both files but with different values.
	OnlyMismatched bool
	// Keys, if non-empty, restricts comparison to only these keys.
	Keys []string
}

// Filter applies FilterOptions to a Result, returning a new Result containing
// only the differences that match the specified criteria.
func Filter(r Result, opts FilterOptions) Result {
	filtered := Result{}

	allowedKeys := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowedKeys[k] = true
	}

	include := func(key string) bool {
		if len(allowedKeys) == 0 {
			return true
		}
		return allowedKeys[key]
	}

	if !opts.OnlyMismatched {
		for _, k := range r.MissingInSecond {
			if include(k) {
				filtered.MissingInSecond = append(filtered.MissingInSecond, k)
			}
		}
		for _, k := range r.MissingInFirst {
			if include(k) {
				filtered.MissingInFirst = append(filtered.MissingInFirst, k)
			}
		}
	}

	if !opts.OnlyMissing {
		for _, m := range r.Mismatched {
			if include(m.Key) {
				filtered.Mismatched = append(filtered.Mismatched, m)
			}
		}
	}

	return filtered
}
