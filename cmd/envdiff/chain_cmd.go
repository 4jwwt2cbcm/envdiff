package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runChain applies a named sequence of built-in transforms to an env file
// and prints a per-step summary.
//
// Usage: envdiff chain <file> <step1,step2,...>
// Available steps: trim, uppercase-keys, lowercase-keys, prefix:<PREFIX>
func runChain(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff chain <file> <step1,step2,...>")
		os.Exit(1)
	}

	filePath := args[0]
	stepNames := strings.Split(args[1], ",")

	env, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", filePath, err)
		os.Exit(1)
	}

	steps, err := resolveSteps(stepNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	results, err := diff.Chain(env, steps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chain error: %v\n", err)
		os.Exit(1)
	}

	for _, line := range diff.ChainSummary(results, env) {
		fmt.Println(line)
	}

	final := results[len(results)-1].Output
	fmt.Println("\nFinal keys:")
	for k, v := range final {
		fmt.Printf("  %s=%s\n", k, v)
	}
}

func resolveSteps(names []string) ([]diff.ChainStep, error) {
	steps := make([]diff.ChainStep, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		switch {
		case name == "trim":
			steps = append(steps, diff.ChainStep{Name: "trim", Fn: diff.ApplyTransformFunc(strings.TrimSpace, nil)})
		case name == "uppercase-keys":
			steps = append(steps, diff.ChainStep{Name: "uppercase-keys", Fn: diff.ApplyTransformFunc(nil, strings.ToUpper)})
		case name == "lowercase-keys":
			steps = append(steps, diff.ChainStep{Name: "lowercase-keys", Fn: diff.ApplyTransformFunc(nil, strings.ToLower)})
		case strings.HasPrefix(name, "prefix:"):
			pfx := strings.TrimPrefix(name, "prefix:")
			captured := pfx
			steps = append(steps, diff.ChainStep{Name: name, Fn: diff.ApplyTransformFunc(nil, func(k string) string {
				return captured + k
			})})
		default:
			return nil, fmt.Errorf("unknown chain step %q", name)
		}
	}
	return steps, nil
}
