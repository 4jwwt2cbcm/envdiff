package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runSaveSnapshot parses an env file and saves it as a named snapshot.
// Usage: envdiff snapshot save <env-file> <snapshot-file> [label]
func runSaveSnapshot(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff snapshot save <env-file> <snapshot-file> [label]")
	}
	envFile := args[0]
	snapFile := args[1]
	label := envFile
	if len(args) >= 3 {
		label = args[2]
	}

	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", envFile, err)
	}

	snap := diff.NewSnapshot(label, env)
	if err := diff.SaveSnapshot(snapFile, snap); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot '%s' saved to %s (%d keys)\n", label, snapFile, len(env))
	return nil
}

// runDiffSnapshot compares a current env file against a saved snapshot.
// Usage: envdiff snapshot diff <env-file> <snapshot-file> [--strict]
func runDiffSnapshot(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff snapshot diff <env-file> <snapshot-file> [--strict]")
	}
	envFile := args[0]
	snapFile := args[1]
	strict := len(args) >= 3 && args[2] == "--strict"

	currentEnv, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", envFile, err)
	}

	snap, err := diff.LoadSnapshot(snapFile)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	result := diff.DiffSnapshot(currentEnv, snap)
	report := diff.Report(snap.Label, envFile, result)
	fmt.Print(report)

	if strict && (len(result.MissingInFirst) > 0 || len(result.MissingInSecond) > 0 || len(result.Mismatched) > 0) {
		os.Exit(1)
	}
	return nil
}
