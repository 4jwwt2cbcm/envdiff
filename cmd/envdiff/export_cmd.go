package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runExport handles the `export` sub-command.
//
// Usage:
//
//	envdiff export <file> [--format=dotenv|shell|json] [--prefix=PREFIX] [--sorted] [--out=FILE]
func runExport(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff export <file> [--format=dotenv|shell|json] [--prefix=PREFIX] [--sorted] [--out=FILE]")
		os.Exit(2)
	}

	envFile := args[0]
	formatStr := "dotenv"
	prefix := ""
	sorted := false
	outPath := ""

	for _, arg := range args[1:] {
		switch {
		case hasFlag(arg, "--format"):
			formatStr = flagValue(arg)
		case hasFlag(arg, "--prefix"):
			prefix = flagValue(arg)
		case arg == "--sorted":
			sorted = true
		case hasFlag(arg, "--out"):
			outPath = flagValue(arg)
		}
	}

	fmt_, err := diff.ParseExportFormat(formatStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(2)
	}

	env, err := parser.ParseFile(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export: cannot read %s: %v\n", envFile, err)
		os.Exit(1)
	}

	opts := diff.ExportOptions{
		Format:  fmt_,
		Prefix:  prefix,
		Sorted:  sorted,
	}

	content, err := diff.ExportEnv(env, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "export:", err)
		os.Exit(1)
	}

	if outPath != "" {
		if err := diff.WriteExport(outPath, content); err != nil {
			fmt.Fprintln(os.Stderr, "export: write:", err)
			os.Exit(1)
		}
		fmt.Printf("exported %s → %s\n", envFile, outPath)
		return
	}

	fmt.Print(content)
}

// hasFlag returns true when arg starts with "name=".
func hasFlag(arg, name string) bool {
	return len(arg) > len(name)+1 && arg[:len(name)+1] == name+"="
}

// flagValue returns the value part of a "--flag=value" argument.
func flagValue(arg string) string {
	for i, c := range arg {
		if c == '=' {
			return arg[i+1:]
		}
	}
	return ""
}
