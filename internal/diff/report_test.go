package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestReport_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	diff.Report(&buf, diff.Result{}, "a.env", "b.env")

	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestReport_MissingInSecond(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{MissingInSecond: []string{"SECRET_KEY"}}
	diff.Report(&buf, r, "dev.env", "prod.env")

	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "dev.env") {
		t.Errorf("expected dev.env label in output, got: %s", out)
	}
}

func TestReport_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{
		Mismatched: []diff.MismatchedKey{
			{Key: "DB_HOST", First: "localhost", Second: "db.prod.example.com"},
		},
	}
	diff.Report(&buf, r, "dev.env", "prod.env")

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected first value in output, got: %s", out)
	}
	if !strings.Contains(out, "db.prod.example.com") {
		t.Errorf("expected second value in output, got: %s", out)
	}
}
