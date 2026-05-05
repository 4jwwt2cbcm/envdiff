package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runSaveBaseline saves a baseline snapshot from the given .env file.
func runSaveBaseline(envFile, outputPath, environment string) error {
	keys, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse %q: %w", envFile, err)
	}
	if err := diff.SaveBaseline(outputPath, environment, keys); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Baseline saved to %s (env: %s, %d keys)\n", outputPath, environment, len(keys))
	return nil
}

// runDiffBaseline compares a .env file against a saved baseline and prints drift.
func runDiffBaseline(envFile, baselinePath string, strict bool) error {
	baseline, err := diff.LoadBaseline(baselinePath)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	current, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse %q: %w", envFile, err)
	}

	report := diff.NewBaselineDriftReport(baseline, current)

	if !report.HasDrift() {
		fmt.Fprintln(os.Stdout, "No drift detected from baseline.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Drift detected from baseline (env: %s):\n", report.Environment)
	formatted, err := diff.FormatReport(diff.FormatText, diff.Report(report.Result, false))
	if err != nil {
		return fmt.Errorf("format report: %w", err)
	}
	fmt.Fprint(os.Stdout, formatted)

	if strict {
		return fmt.Errorf("baseline drift found")
	}
	return nil
}
