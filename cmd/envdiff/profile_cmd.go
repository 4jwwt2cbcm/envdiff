package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runProfileCompare filters one or two env files by a named profile and reports differences.
// Usage: envdiff profile <profile-file> <profile-name> <file1> [file2]
func runProfileCompare(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envdiff profile <profile-file> <profile-name> <file1> [file2]")
	}
	profilePath := args[0]
	profileName := args[1]
	file1 := args[2]

	pm, err := diff.LoadProfiles(profilePath)
	if err != nil {
		return fmt.Errorf("load profiles: %w", err)
	}

	env1, err := parser.ParseFile(file1)
	if err != nil {
		return fmt.Errorf("parse %s: %w", file1, err)
	}

	filtered1, err := diff.FilterByProfile(env1, pm, profileName)
	if err != nil {
		return err
	}

	if len(args) < 4 {
		// Single file: just list keys present in the profile.
		fmt.Printf("Profile %q keys found in %s:\n", profileName, file1)
		for k, v := range filtered1 {
			fmt.Printf("  %s=%s\n", k, v)
		}
		return nil
	}

	file2 := args[3]
	env2, err := parser.ParseFile(file2)
	if err != nil {
		return fmt.Errorf("parse %s: %w", file2, err)
	}

	filtered2, err := diff.FilterByProfile(env2, pm, profileName)
	if err != nil {
		return err
	}

	report := diff.Compare(filtered1, filtered2)
	diff.Report(os.Stdout, report, false)
	return nil
}

// runListProfiles lists all profile names from a profile file.
// Usage: envdiff profiles list <profile-file>
func runListProfiles(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff profiles list <profile-file>")
	}
	pm, err := diff.LoadProfiles(args[0])
	if err != nil {
		return fmt.Errorf("load profiles: %w", err)
	}
	names := diff.ListProfileNames(pm)
	if len(names) == 0 {
		fmt.Println("No profiles defined.")
		return nil
	}
	fmt.Printf("Profiles (%d):\n", len(names))
	for _, name := range names {
		p := pm[name]
		fmt.Printf("  %-20s keys: %s\n", name, strings.Join(p.Keys, ", "))
	}
	return nil
}
