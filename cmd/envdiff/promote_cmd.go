package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runPromote handles the `promote` subcommand.
// Usage: envdiff promote [--overwrite] [--keys KEY1,KEY2] <src.env> <dst.env> [output.env]
//
// It reads key-value pairs from src.env and merges them into dst.env.
// By default, only keys missing from dst.env are added. Use --overwrite to
// replace existing keys. Use --keys to restrict promotion to specific keys.
// If output.env is provided, the merged result is written there; otherwise
// only the promotion summary is printed to stdout.
func runPromote(args []string) {
	fs := flag.NewFlagSet("promote", flag.ExitOnError)
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	keysFlag := fs.String("keys", "", "comma-separated list of keys to promote (default: all)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error parsing flags:", err)
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff promote [--overwrite] [--keys K1,K2] <src.env> <dst.env> [output.env]")
		os.Exit(1)
	}

	srcFile := positional[0]
	dstFile := positional[1]
	outFile := ""
	if len(positional) >= 3 {
		outFile = positional[2]
	}

	src, err := parser.ParseFile(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file %s: %v\n", srcFile, err)
		os.Exit(1)
	}

	dst, err := parser.ParseFile(dstFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading destination file %s: %v\n", dstFile, err)
		os.Exit(1)
	}

	opts := diff.PromoteOptions{Overwrite: *overwrite}
	if *keysFlag != "" {
		for _, k := range strings.Split(*keysFlag, ",") {
			if k = strings.TrimSpace(k); k != "" {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	merged, pr := diff.PromoteEnv(src, dst, opts)
	fmt.Print(diff.FormatPromoteResult(pr))

	if outFile != "" {
		if err := diff.WriteMerged(merged, outFile); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Written to %s\n", outFile)
	}
}
