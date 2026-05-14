package diff

import (
	"fmt"
	"sort"
)

// ChainStep represents a single transformation step in a pipeline.
type ChainStep struct {
	Name string
	Fn   func(map[string]string) (map[string]string, error)
}

// ChainResult holds the output of each step in the chain.
type ChainResult struct {
	Step   string
	Output map[string]string
	Err    error
}

// Chain runs a sequence of transformation steps against an env map,
// passing the output of each step as the input to the next.
func Chain(env map[string]string, steps []ChainStep) ([]ChainResult, error) {
	results := make([]ChainResult, 0, len(steps))
	current := copyEnv(env)

	for _, step := range steps {
		out, err := step.Fn(current)
		results = append(results, ChainResult{
			Step:   step.Name,
			Output: copyEnv(out),
			Err:    err,
		})
		if err != nil {
			return results, fmt.Errorf("chain step %q failed: %w", step.Name, err)
		}
		current = out
	}

	return results, nil
}

// ChainSummary returns a human-readable summary of keys added/removed per step.
func ChainSummary(results []ChainResult, original map[string]string) []string {
	lines := []string{}
	prev := original
	for _, r := range results {
		added, removed := diffKeys(prev, r.Output)
		lines = append(lines, fmt.Sprintf("[%s] +%d keys, -%d keys", r.Step, len(added), len(removed)))
		prev = r.Output
	}
	return lines
}

func copyEnv(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func diffKeys(before, after map[string]string) (added, removed []string) {
	for k := range after {
		if _, ok := before[k]; !ok {
			added = append(added, k)
		}
	}
	for k := range before {
		if _, ok := after[k]; !ok {
			removed = append(removed, k)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return
}
