package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "diff":
		runDiff(os.Args[2:])
	case "baseline":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(2)
		}
		switch os.Args[2] {
		case "save":
			runSaveBaseline(os.Args[3:])
		case "diff":
			runDiffBaseline(os.Args[3:])
		default:
			printUsage()
			os.Exit(2)
		}
	case "validate":
		runValidate(os.Args[2:])
	case "lint":
		runLint(os.Args[2:])
	case "rename":
		runApplyRenames(os.Args[2:])
	case "patch":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(2)
		}
		switch os.Args[2] {
		case "generate":
			runGeneratePatch(os.Args[3:])
		case "apply":
			runApplyPatch(os.Args[3:])
		default:
			printUsage()
			os.Exit(2)
		}
	case "audit":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(2)
		}
		switch os.Args[2] {
		case "show":
			runShowAuditLog(os.Args[3:])
		case "diff":
			runAuditDiff(os.Args[3:])
		default:
			printUsage()
			os.Exit(2)
		}
	case "template":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(2)
		}
		switch os.Args[2] {
		case "generate":
			runGenerateTemplate(os.Args[3:])
		case "check":
			runCheckTemplate(os.Args[3:])
		default:
			printUsage()
			os.Exit(2)
		}
	default:
		printUsage()
		os.Exit(2)
	}
}

func runDiff(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff diff <file1> <file2> [--strict] [--format text|json|csv]")
		os.Exit(2)
	}
	a, err := parser.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	b, err := parser.ParseFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	report := diff.Compare(a, b)
	fmt.Print(diff.Report(report, false))
	if len(os.Args) > 4 && os.Args[4] == "--strict" {
		if len(report.MissingInSecond)+len(report.MissingInFirst)+len(report.Mismatched) > 0 {
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `envdiff — compare and manage .env files

Usage:
  envdiff diff <file1> <file2> [--strict]
  envdiff baseline save <env-file> --out <baseline.json>
  envdiff baseline diff <env-file> --baseline <baseline.json>
  envdiff validate <env-file> --rules <rules.json>
  envdiff lint <env-file>
  envdiff rename <env-file> --map <renames.json>
  envdiff patch generate <file1> <file2>
  envdiff patch apply <env-file> --patch <patch.json>
  envdiff audit show --log <audit.json>
  envdiff audit diff <env-file>
  envdiff template generate <env-file> [--out template.json]
  envdiff template check <env-file> [--template template.json]`)
}
