package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestPatchCmd_GenerateAndApply(t *testing.T) {
	bin := buildBinary(t)

	srcFile := writeTempEnv(t, "A=1\nB=2\nC=3\n")
	dstFile := writeTempEnv(t, "A=1\nB=99\nD=4\n")

	out, err := runCmd(t, bin, "patch", srcFile, dstFile)
	if err != nil {
		t.Fatalf("patch command failed: %v\noutput: %s", err, out)
	}

	for _, want := range []string{"+ D=4", "- C=3", "~ B: 2 -> 99"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestPatchCmd_JSONOutput(t *testing.T) {
	bin := buildBinary(t)

	srcFile := writeTempEnv(t, "X=hello\n")
	dstFile := writeTempEnv(t, "X=world\n")

	out, err := runCmd(t, bin, "patch", "--json", srcFile, dstFile)
	if err != nil {
		t.Fatalf("patch --json failed: %v\noutput: %s", err, out)
	}

	var p diff.Patch
	if err := json.Unmarshal([]byte(out), &p); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	if len(p.Ops) != 1 || p.Ops[0].Action != "update" {
		t.Errorf("expected 1 update op, got %+v", p.Ops)
	}
}

func TestPatchCmd_WriteToFile(t *testing.T) {
	bin := buildBinary(t)

	srcFile := writeTempEnv(t, "A=1\n")
	dstFile := writeTempEnv(t, "A=2\n")
	out := filepath.Join(t.TempDir(), "my.patch")

	_, err := runCmd(t, bin, "patch", "--output", out, srcFile, dstFile)
	if err != nil {
		t.Fatalf("patch --output failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("patch file not written: %v", err)
	}
	if !strings.Contains(string(data), "~ A: 1 -> 2") {
		t.Errorf("unexpected patch content:\n%s", data)
	}
}

func TestPatchCmd_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	_, err := runCmd(t, bin, "patch")
	if err == nil {
		t.Error("expected non-zero exit for missing args")
	}
}
