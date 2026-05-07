package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runGenerateTemplate parses args for the "template generate" sub-command.
// Usage: envdiff template generate <env-file> --out <template.json>
func runGenerateTemplate(args []string) {
	fs := flag.NewFlagSet("template generate", flag.ExitOnError)
	out := fs.String("out", "template.json", "output path for the generated template")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff template generate <env-file> [--out template.json]")
		os.Exit(2)
	}

	envPath := fs.Arg(0)
	env, err := parser.ParseFile(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading env file: %v\n", err)
		os.Exit(1)
	}

	tmpl := diff.GenerateTemplate(env)
	if err := diff.SaveTemplate(*out, tmpl); err != nil {
		fmt.Fprintf(os.Stderr, "error saving template: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("template written to %s (%d vars)\n", *out, len(tmpl.Vars))
}

// runCheckTemplate validates an env file against a saved template.
// Usage: envdiff template check <env-file> --template <template.json>
func runCheckTemplate(args []string) {
	fs := flag.NewFlagSet("template check", flag.ExitOnError)
	tmplPath := fs.String("template", "template.json", "path to template file")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff template check <env-file> [--template template.json]")
		os.Exit(2)
	}

	envPath := fs.Arg(0)
	env, err := parser.ParseFile(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading env file: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := diff.LoadTemplate(*tmplPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading template: %v\n", err)
		os.Exit(1)
	}

	violations := diff.CheckAgainstTemplate(env, tmpl)
	if len(violations) == 0 {
		fmt.Println("OK: env file satisfies template")
		return
	}
	fmt.Fprintf(os.Stderr, "template violations (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v)
	}
	os.Exit(1)
}
