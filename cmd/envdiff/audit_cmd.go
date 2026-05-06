package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
)

// runShowAuditLog prints the audit log at the given path to stdout.
func runShowAuditLog(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff audit <audit-log-path>")
		os.Exit(1)
	}
	path := args[0]

	log, err := diff.LoadAuditLog(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading audit log: %v\n", err)
		os.Exit(1)
	}

	if len(log.Entries) == 0 {
		fmt.Println("No audit entries found.")
		return
	}

	for i, e := range log.Entries {
		fmt.Printf("[%d] %s  %s <-> %s\n", i+1, e.Timestamp.Format("2006-01-02 15:04:05"), e.FileA, e.FileB)
		fmt.Printf("     missing_in_a=%d  missing_in_b=%d  mismatched=%d  drift=%v\n",
			e.MissingInA, e.MissingInB, e.Mismatched, e.HasDrift)
	}
}

// runAuditDiff performs a diff between two env files, records the result in the
// audit log, and exits non-zero when drift is detected and --strict is set.
func runAuditDiff(fileA, fileB, auditPath string, strict bool) {
	envA, err := parseEnvFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileA, err)
		os.Exit(1)
	}
	envB, err := parseEnvFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	report := diff.Compare(envA, envB)
	entry := diff.NewAuditEntry(fileA, fileB, report)

	if auditPath != "" {
		if err := diff.AppendAuditLog(auditPath, entry); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not write audit log: %v\n", err)
		}
	}

	data, _ := json.MarshalIndent(entry, "", "  ")
	fmt.Println(string(data))

	if strict && entry.HasDrift {
		os.Exit(1)
	}
}

// parseEnvFile is a thin wrapper used by audit_cmd to avoid import cycles.
func parseEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_ = data
	// Delegate to the parser package via the shared helper in main.go.
	return nil, fmt.Errorf("use parser.ParseFile directly")
}
