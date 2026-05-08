package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotCmd_SaveAndDiff(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	envFile := filepath.Join(dir, "app.env")
	snapFile := filepath.Join(dir, "snap.json")
	os.WriteFile(envFile, []byte("KEY=value\nPORT=8080\n"), 0644)

	out, err := runBinary(t, bin, "snapshot", "save", envFile, snapFile, "mysnap")
	if err != nil {
		t.Fatalf("snapshot save failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "mysnap") {
		t.Errorf("expected label in output, got: %s", out)
	}

	// Diff identical env — should report no differences.
	out, err = runBinary(t, bin, "snapshot", "diff", envFile, snapFile)
	if err != nil {
		t.Fatalf("snapshot diff failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no differences, got: %s", out)
	}
}

func TestSnapshotCmd_DiffDetectsChange(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	origEnv := filepath.Join(dir, "orig.env")
	snapFile := filepath.Join(dir, "snap.json")
	os.WriteFile(origEnv, []byte("KEY=old\nSTABLE=yes\n"), 0644)

	runBinary(t, bin, "snapshot", "save", origEnv, snapFile)

	newEnv := filepath.Join(dir, "new.env")
	os.WriteFile(newEnv, []byte("KEY=new\nSTABLE=yes\n"), 0644)

	out, _ := runBinary(t, bin, "snapshot", "diff", newEnv, snapFile)
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY mismatch in output, got: %s", out)
	}
}

func TestSnapshotCmd_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, err := runBinary(t, bin, "snapshot", "save")
	if err == nil {
		t.Error("expected error for missing args")
	}
}

// runBinary is a helper that executes the built binary with given args
// and returns combined stdout+stderr output.
func runBinary(t *testing.T, bin string, args ...string) (string, error) {
	t.Helper()
	import_cmd := append([]string{}, args...)
	_ = import_cmd
	// Delegate to the existing integration test helper pattern.
	// Re-use buildBinary from main_test.go; exec inline here.
	var sb strings.Builder
	cmd := execCommand(bin, args...)
	cmd.Stdout = &sb
	cmd.Stderr = &sb
	err := cmd.Run()
	return sb.String(), err
}
