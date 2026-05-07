package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestTemplateCmd_Generate(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	envFile := filepath.Join(dir, ".env")
	writeTempEnv(t, envFile, "APP_ENV=production\nDB_URL=postgres://localhost/db\n")

	outPath := filepath.Join(dir, "out.json")
	cmd := exec.Command(bin, "template", "generate", envFile, "--out", outPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generate failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "template written to") {
		t.Errorf("unexpected output: %s", out)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	var tmpl diff.Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(tmpl.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(tmpl.Vars))
	}
}

func TestTemplateCmd_CheckValid(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	tmplPath := filepath.Join(dir, "template.json")
	tmpl := diff.Template{
		Vars: []diff.TemplateVar{
			{Key: "HOST", Required: true},
			{Key: "PORT", Required: false},
		},
	}
	diff.SaveTemplate(tmplPath, tmpl)

	envFile := filepath.Join(dir, ".env")
	writeTempEnv(t, envFile, "HOST=localhost\nPORT=8080\n")

	cmd := exec.Command(bin, "template", "check", envFile, "--template", tmplPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check failed unexpectedly: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "OK") {
		t.Errorf("expected OK output, got: %s", out)
	}
}

func TestTemplateCmd_CheckViolation(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	tmplPath := filepath.Join(dir, "template.json")
	tmpl := diff.Template{
		Vars: []diff.TemplateVar{
			{Key: "SECRET_KEY", Required: true},
		},
	}
	diff.SaveTemplate(tmplPath, tmpl)

	envFile := filepath.Join(dir, ".env")
	writeTempEnv(t, envFile, "OTHER=value\n")

	cmd := exec.Command(bin, "template", "check", envFile, "--template", tmplPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for violation")
	}
	if !strings.Contains(string(out), "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
}

func TestTemplateCmd_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "template", "generate")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing args")
	}
	if !strings.Contains(string(out), "usage") {
		t.Errorf("expected usage message, got: %s", out)
	}
}
