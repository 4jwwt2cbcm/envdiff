package diff

import (
	"fmt"
	"sort"
)

// PromoteOptions controls how promotion between environments is performed.
type PromoteOptions struct {
	// Overwrite allows values in the destination to be overwritten.
	Overwrite bool
	// Keys restricts promotion to only the specified keys. If empty, all keys are promoted.
	Keys []string
}

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Promoted  []string // keys that were promoted (added or updated)
	Skipped   []string // keys skipped because they already exist and Overwrite is false
	NotFound  []string // keys requested via Options.Keys that were absent in source
}

// PromoteEnv copies keys from src into dst according to opts.
// It returns a PromoteResult describing what happened.
func PromoteEnv(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult) {
	result := map[string]string{}
	for k, v := range dst {
		result[k] = v
	}

	keySet := buildKeySet(opts.Keys, src)

	var pr PromoteResult
	for _, key := range sortedKeys(keySet) {
		srcVal, inSrc := src[key]
		if !inSrc {
			pr.NotFound = append(pr.NotFound, key)
			continue
		}
		_, inDst := dst[key]
		if inDst && !opts.Overwrite {
			pr.Skipped = append(pr.Skipped, key)
			continue
		}
		result[key] = srcVal
		pr.Promoted = append(pr.Promoted, key)
	}
	return result, pr
}

// FormatPromoteResult returns a human-readable summary of a PromoteResult.
func FormatPromoteResult(pr PromoteResult) string {
	out := fmt.Sprintf("Promoted: %d key(s)\n", len(pr.Promoted))
	for _, k := range pr.Promoted {
		out += fmt.Sprintf("  + %s\n", k)
	}
	if len(pr.Skipped) > 0 {
		out += fmt.Sprintf("Skipped (already exists): %d key(s)\n", len(pr.Skipped))
		for _, k := range pr.Skipped {
			out += fmt.Sprintf("  ~ %s\n", k)
		}
	}
	if len(pr.NotFound) > 0 {
		out += fmt.Sprintf("Not found in source: %d key(s)\n", len(pr.NotFound))
		for _, k := range pr.NotFound {
			out += fmt.Sprintf("  ! %s\n", k)
		}
	}
	return out
}

func buildKeySet(keys []string, src map[string]string) []string {
	if len(keys) > 0 {
		return keys
	}
	return sortedKeys(func() []string {
		ks := make([]string, 0, len(src))
		for k := range src {
			ks = append(ks, k)
		}
		return ks
	}())
}

func sortedKeys(keys []string) []string {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	return sorted
}
