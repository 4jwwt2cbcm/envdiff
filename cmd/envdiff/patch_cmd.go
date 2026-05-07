package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runGeneratePatch generates a patch from src to dst and writes it out.
func runGeneratePatch(args []string) {
	fs := flag.NewFlagSet("patch", flag.ExitOnError)
	output := fs.String("output", "", "write patch to file (default: stdout)")
	jsonFmt := fs.Bool("json", false, "output patch as JSON")
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff patch [--output file] [--json] <src.env> <dst.env>")
		os.Exit(1)
	}

	src, err := parser.ParseFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading src: %v\n", err)
		os.Exit(1)
	}
	dst, err := parser.ParseFile(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading dst: %v\n", err)
		os.Exit(1)
	}

	p := diff.GeneratePatch(src, dst)

	if *jsonFmt {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(p); err != nil {
			fmt.Fprintf(os.Stderr, "json encode error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *output != "" {
		if err := diff.WritePatch(*output, p); err != nil {
			fmt.Fprintf(os.Stderr, "error writing patch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("patch written to %s (%d ops)\n", *output, len(p.Ops))
		return
	}

	for _, op := range p.Ops {
		switch op.Action {
		case "add":
			fmt.Printf("+ %s=%s\n", op.Key, op.NewVal)
		case "remove":
			fmt.Printf("- %s=%s\n", op.Key, op.OldVal)
		case "update":
			fmt.Printf("~ %s: %s -> %s\n", op.Key, op.OldVal, op.NewVal)
		}
	}
}

// runApplyPatch applies a JSON patch file to a src env and prints the result.
func runApplyPatch(args []string) {
	fs := flag.NewFlagSet("apply-patch", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff apply-patch <src.env> <patch.json>")
		os.Exit(1)
	}

	src, err := parser.ParseFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading src: %v\n", err)
		os.Exit(1)
	}

	data, err := os.ReadFile(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading patch: %v\n", err)
		os.Exit(1)
	}

	var p diff.Patch
	if err := json.Unmarshal(data, &p); err != nil {
		fmt.Fprintf(os.Stderr, "invalid patch JSON: %v\n", err)
		os.Exit(1)
	}

	result, err := diff.ApplyPatch(src, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patch apply error: %v\n", err)
		os.Exit(1)
	}

	for k, v := range result {
		fmt.Printf("%s=%s\n", k, v)
	}
}
