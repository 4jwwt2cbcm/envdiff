package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleExportEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/mydb",
		"SECRET":   "s3cr3t",
	}
}

func TestParseExportFormat_Valid(t *testing.T) {
	cases := []struct{ in, want string }{
		{"dotenv", "dotenv"},
		{"shell", "shell"},
		{"json", "json"},
		{"JSON", "json"},
	}
	for _, c := range cases {
		f, err := ParseExportFormat(c.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", c.in, err)
		}
		if string(f) != c.want {
			t.Errorf("got %q, want %q", f, c.want)
		}
	}
}

func TestParseExportFormat_Invalid(t *testing.T) {
	_, err := ParseExportFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExportEnv_Dotenv(t *testing.T) {
	out, err := ExportEnv(sampleExportEnv(), ExportOptions{Format: ExportFormatDotenv, Sorted: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "APP_HOST=") {
		t.Errorf("expected APP_HOST in dotenv output, got:\n%s", out)
	}
}

func TestExportEnv_Shell(t *testing.T) {
	out, err := ExportEnv(sampleExportEnv(), ExportOptions{Format: ExportFormatShell, Sorted: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("expected 'export' prefix, got: %s", line)
		}
	}
}

func TestExportEnv_JSON(t *testing.T) {
	out, err := ExportEnv(sampleExportEnv(), ExportOptions{Format: ExportFormatJSON})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if m["APP_HOST"] != "localhost" {
		t.Errorf("unexpected value for APP_HOST: %s", m["APP_HOST"])
	}
}

func TestExportEnv_PrefixFilter(t *testing.T) {
	out, err := ExportEnv(sampleExportEnv(), ExportOptions{Format: ExportFormatDotenv, Prefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "DB_URL") {
		t.Errorf("DB_URL should be filtered out by prefix APP_")
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Errorf("APP_PORT should be present")
	}
}

func TestWriteExport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "out.env")
	content := "KEY=value\n"
	if err := WriteExport(path, content); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Errorf("got %q, want %q", string(got), content)
	}
}
