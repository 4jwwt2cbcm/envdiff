package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoDifferences(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	b := writeTempEnv(t, "KEY=value\nFOO=bar\n")

	cmd := exec.Command(bin, a, b)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "No differences") {
		t.Errorf("expected 'No differences', got: %s", out)
	}
}

func TestMain_StrictExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "KEY=value\nONLY_A=yes\n")
	b := writeTempEnv(t, "KEY=value\n")

	cmd := exec.Command(bin, "-strict", a, b)
	err := cmd.Run()
	if err == nil {
		t.Error("expected non-zero exit code with -strict and differences")
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("expected non-zero exit when no args provided")
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage message, got: %s", out)
	}
}
