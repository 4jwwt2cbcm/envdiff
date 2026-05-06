package diff

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy string

const (
	StrategyFirst  MergeStrategy = "first"  // keep value from first file
	StrategySecond MergeStrategy = "second" // keep value from second file
	StrategyUnion  MergeStrategy = "union"  // include all keys from both files
)

// MergeResult holds the merged key-value pairs and metadata.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string // keys that had differing values
}

// Merge combines two env maps using the given strategy.
// Keys unique to either map are always included.
// Conflicting keys are resolved by strategy.
func Merge(first, second map[string]string, strategy MergeStrategy) MergeResult {
	result := MergeResult{
		Merged: make(map[string]string),
	}

	// Copy all keys from first
	for k, v := range first {
		result.Merged[k] = v
	}

	for k, v := range second {
		existing, exists := result.Merged[k]
		if !exists {
			result.Merged[k] = v
			continue
		}
		if existing != v {
			result.Conflicts = append(result.Conflicts, k)
			if strategy == StrategySecond {
				result.Merged[k] = v
			}
			// StrategyFirst and StrategyUnion keep the existing (first) value
		}
	}

	sort.Strings(result.Conflicts)
	return result
}

// WriteMerged writes merged key-value pairs to a file in .env format.
func WriteMerged(path string, merged map[string]string) error {
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := merged[k]
		if strings.ContainsAny(v, " \t#") {
			v = fmt.Sprintf(`"%s"`, v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
