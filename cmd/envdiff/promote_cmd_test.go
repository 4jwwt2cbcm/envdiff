package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePromoteEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPromoteCmd_Basic(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	src := writePromoteEnv(t, dir, "src.env", "A=1\nB=2\n")
	dst := writePromoteEnv(t, dir, "dst.env", "C=3\n")

	out, err := runBinary(bin, "promote", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Promoted: 2") {
		t.Errorf("expected 'Promoted: 2' in output, got:\n%s", out)
	}
}

func TestPromoteCmd_NoOverwrite(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	src := writePromoteEnv(t, dir, "src.env", "A=new\nB=2\n")
	dst := writePromoteEnv(t, dir, "dst.env", "A=old\n")

	out, err := runBinary(bin, "promote", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Skipped") {
		t.Errorf("expected 'Skipped' in output, got:\n%s", out)
	}
}

func TestPromoteCmd_WithOverwrite(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	src := writePromoteEnv(t, dir, "src.env", "A=new\n")
	dst := writePromoteEnv(t, dir, "dst.env", "A=old\n")

	out, err := runBinary(bin, "promote", "--overwrite", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Promoted: 1") {
		t.Errorf("expected 'Promoted: 1' in output, got:\n%s", out)
	}
}

func TestPromoteCmd_WriteOutput(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	src := writePromoteEnv(t, dir, "src.env", "A=1\nB=2\n")
	dst := writePromoteEnv(t, dir, "dst.env", "C=3\n")
	out_file := filepath.Join(dir, "out.env")

	_, err := runBinary(bin, "promote", src, dst, out_file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out_file)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "A=1") {
		t.Errorf("expected A=1 in output file, got:\n%s", string(data))
	}
}

func TestPromoteCmd_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, err := runBinary(bin, "promote")
	if err == nil {
		t.Error("expected non-zero exit for missing args")
	}
}
