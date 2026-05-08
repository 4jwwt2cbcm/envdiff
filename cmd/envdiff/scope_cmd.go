package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runScopeCompare loads two env files and a scopes definition file, then
// prints per-scope diff reports to stdout.
//
// Usage: envdiff scope <scopes.json> <file-a> <file-b>
func runScopeCompare(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envdiff scope <scopes.json> <file-a> <file-b>")
	}
	scopesPath, pathA, pathB := args[0], args[1], args[2]

	rules, err := diff.LoadScopes(scopesPath)
	if err != nil {
		return fmt.Errorf("load scopes: %w", err)
	}

	envA, err := parser.ParseFile(pathA)
	if err != nil {
		return fmt.Errorf("parse %s: %w", pathA, err)
	}
	envB, err := parser.ParseFile(pathB)
	if err != nil {
		return fmt.Errorf("parse %s: %w", pathB, err)
	}

	reports := diff.CompareByScopes(envA, envB, rules)

	hasDiff := false
	for _, sr := range reports {
		r := sr.Report
		if len(r.MissingInFirst) == 0 && len(r.MissingInSecond) == 0 && len(r.Mismatched) == 0 {
			fmt.Printf("[%s] no differences\n", sr.Scope)
			continue
		}
		hasDiff = true
		fmt.Printf("[%s]\n", sr.Scope)
		for _, k := range r.MissingInSecond {
			fmt.Printf("  missing in %s: %s\n", pathB, k)
		}
		for _, k := range r.MissingInFirst {
			fmt.Printf("  missing in %s: %s\n", pathA, k)
		}
		for _, m := range r.Mismatched {
			fmt.Printf("  mismatch: %s (%q vs %q)\n", m.Key, m.ValueA, m.ValueB)
		}
	}

	if hasDiff {
		os.Exit(1)
	}
	return nil
}
