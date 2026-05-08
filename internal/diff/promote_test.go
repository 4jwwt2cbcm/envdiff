package diff

import (
	"strings"
	"testing"
)

func TestPromoteEnv_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"D": "4"}

	result, pr := PromoteEnv(src, dst, PromoteOptions{})

	if len(pr.Promoted) != 3 {
		t.Errorf("expected 3 promoted, got %d", len(pr.Promoted))
	}
	if result["A"] != "1" || result["B"] != "2" || result["C"] != "3" || result["D"] != "4" {
		t.Error("unexpected result values")
	}
}

func TestPromoteEnv_NoOverwrite(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}

	result, pr := PromoteEnv(src, dst, PromoteOptions{Overwrite: false})

	if len(pr.Skipped) != 1 || pr.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", pr.Skipped)
	}
	if result["A"] != "old" {
		t.Errorf("expected A to remain 'old', got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Error("expected B to be promoted")
	}
}

func TestPromoteEnv_WithOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	result, pr := PromoteEnv(src, dst, PromoteOptions{Overwrite: true})

	if len(pr.Promoted) != 1 {
		t.Errorf("expected 1 promoted, got %d", len(pr.Promoted))
	}
	if result["A"] != "new" {
		t.Errorf("expected A='new', got %s", result["A"])
	}
}

func TestPromoteEnv_KeyRestriction(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	_, pr := PromoteEnv(src, dst, PromoteOptions{Keys: []string{"A", "C"}})

	if len(pr.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(pr.Promoted))
	}
	for _, k := range pr.Promoted {
		if k != "A" && k != "C" {
			t.Errorf("unexpected key promoted: %s", k)
		}
	}
}

func TestPromoteEnv_NotFound(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}

	_, pr := PromoteEnv(src, dst, PromoteOptions{Keys: []string{"A", "MISSING"}})

	if len(pr.NotFound) != 1 || pr.NotFound[0] != "MISSING" {
		t.Errorf("expected MISSING in NotFound, got %v", pr.NotFound)
	}
}

func TestFormatPromoteResult(t *testing.T) {
	pr := PromoteResult{
		Promoted: []string{"A", "B"},
		Skipped:  []string{"C"},
		NotFound: []string{"D"},
	}
	out := FormatPromoteResult(pr)
	if !strings.Contains(out, "Promoted: 2") {
		t.Error("expected promoted count in output")
	}
	if !strings.Contains(out, "Skipped") {
		t.Error("expected skipped section")
	}
	if !strings.Contains(out, "Not found") {
		t.Error("expected not found section")
	}
}
