package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeExportEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestExportCmd_DefaultDotenv(t *testing.T) {
	bin := buildBinary(t)
	env := writeExportEnv(t, "APP_HOST=localhost\nAPP_PORT=9000\n")

	out, err := runBinary(t, bin, "export", env, "--sorted")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "APP_HOST=") {
		t.Errorf("expected APP_HOST in output, got:\n%s", out)
	}
}

func TestExportCmd_ShellFormat(t *testing.T) {
	bin := buildBinary(t)
	env := writeExportEnv(t, "DB_URL=postgres://localhost/db\n")

	out, err := runBinary(t, bin, "export", env, "--format=shell")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "export ") {
		t.Errorf("expected shell export prefix, got:\n%s", out)
	}
}

func TestExportCmd_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	env := writeExportEnv(t, "KEY=value\n")

	out, err := runBinary(t, bin, "export", env, "--format=json")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected JSON key in output, got:\n%s", out)
	}
}

func TestExportCmd_WriteToFile(t *testing.T) {
	bin := buildBinary(t)
	env := writeExportEnv(t, "X=1\nY=2\n")
	outFile := filepath.Join(t.TempDir(), "result.env")

	_, err := runBinary(t, bin, "export", env, "--out="+outFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, readErr := os.ReadFile(outFile)
	if readErr != nil {
		t.Fatalf("output file not created: %v", readErr)
	}
	if !strings.Contains(string(data), "X=") {
		t.Errorf("expected X= in output file, got:\n%s", string(data))
	}
}

func TestExportCmd_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, err := runBinary(t, bin, "export")
	if err == nil {
		t.Fatal("expected non-zero exit for missing args")
	}
}
