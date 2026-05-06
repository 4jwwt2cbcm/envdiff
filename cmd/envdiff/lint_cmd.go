package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runLint parses the given .env file and runs lint rules against it.
// Prints all violations to stdout. Exits non-zero if any violations are found.
func runLint(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff lint <file.env>")
		os.Exit(2)
	}

	filePath := args[0]
	env, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", filePath, err)
		os.Exit(2)
	}

	violations := diff.LintEnv(env)
	if len(violations) == 0 {
		fmt.Printf("✔  No lint violations found in %s\n", filePath)
		return
	}

	fmt.Printf("Lint violations in %s:\n", filePath)
	for _, v := range violations {
		fmt.Printf("  [%s] %s\n", v.Rule, v.Message)
	}
	fmt.Printf("\n%d violation(s) found.\n", len(violations))
	os.Exit(1)
}
