package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runApplyRenames applies a rename map to an env file and prints the result,
// or optionally writes it to an output file.
//
// Usage: envdiff rename --map renames.json --input .env [--output .env.renamed]
func runApplyRenames(args []string) {
	fs := flag.NewFlagSet("rename", flag.ExitOnError)
	mapFile := fs.String("map", "", "path to JSON rename map file (required)")
	inputFile := fs.String("input", "", "path to input .env file (required)")
	outputFile := fs.String("output", "", "path to output .env file (optional, prints to stdout if omitted)")
	_ = fs.Parse(args)

	if *mapFile == "" || *inputFile == "" {
		fmt.Fprintln(os.Stderr, "error: --map and --input are required")
		fs.Usage()
		os.Exit(2)
	}

	renames, err := diff.LoadRenameFile(*mapFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading rename map: %v\n", err)
		os.Exit(1)
	}

	env, err := parser.ParseFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env file: %v\n", err)
		os.Exit(1)
	}

	renamed := diff.ApplyRenames(env, renames)

	var lines []string
	for k, v := range renamed {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	output := ""
	for _, l := range lines {
		output += l + "\n"
	}

	if *outputFile == "" {
		fmt.Print(output)
		return
	}

	if err := os.WriteFile(*outputFile, []byte(output), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output file: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "renamed env written to %s\n", *outputFile)
}
