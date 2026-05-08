package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExportFormat represents a supported export target format.
type ExportFormat string

const (
	ExportFormatDotenv ExportFormat = "dotenv"
	ExportFormatShell  ExportFormat = "shell"
	ExportFormatJSON   ExportFormat = "json"
)

// ExportOptions controls how an env map is exported.
type ExportOptions struct {
	Format  ExportFormat
	Prefix  string // optional key prefix filter
	Sorted  bool
}

// ParseExportFormat parses and validates an export format string.
func ParseExportFormat(s string) (ExportFormat, error) {
	switch ExportFormat(strings.ToLower(s)) {
	case ExportFormatDotenv:
		return ExportFormatDotenv, nil
	case ExportFormatShell:
		return ExportFormatShell, nil
	case ExportFormatJSON:
		return ExportFormatJSON, nil
	}
	return "", fmt.Errorf("unsupported export format %q: must be dotenv, shell, or json", s)
}

// ExportEnv serialises env according to opts and returns the result as a string.
func ExportEnv(env map[string]string, opts ExportOptions) (string, error) {
	keys := filteredKeys(env, opts.Prefix)
	if opts.Sorted {
		sort.Strings(keys)
	}

	switch opts.Format {
	case ExportFormatShell:
		return exportShell(env, keys), nil
	case ExportFormatJSON:
		return exportJSON(env, keys)
	default:
		return exportDotenv(env, keys), nil
	}
}

// WriteExport writes the exported content to path, creating parent dirs as needed.
func WriteExport(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("export: create dirs: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func filteredKeys(env map[string]string, prefix string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys
}

func exportDotenv(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportShell(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportJSON(env map[string]string, keys []string) (string, error) {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = env[k]
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("export json: %w", err)
	}
	return string(b) + "\n", nil
}
