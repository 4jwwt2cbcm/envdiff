package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

const usage = `envdiff - Compare .env files across environments

Usage:
  envdiff [flags] <file1> <file2>

Flags:
`

func main() {
	mask := flag.Bool("mask", true, "Mask values in mismatch output")
	strict := flag.Bool("strict", false, "Exit with non-zero code if differences found")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	file1, file2 := args[0], args[1]

	env1, err := parser.ParseFile(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", file1, err)
		os.Exit(1)
	}

	env2, err := parser.ParseFile(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", file2, err)
		os.Exit(1)
	}

	result := diff.Compare(env1, env2)

	opts := diff.ReportOptions{
		FileA:     file1,
		FileB:     file2,
		MaskValues: *mask,
	}

	diff.Report(os.Stdout, result, opts)

	if *strict && result.HasDifferences() {
		os.Exit(1)
	}
}
