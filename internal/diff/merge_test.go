package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	first := map[string]string{"A": "1", "B": "2"}
	second := map[string]string{"C": "3", "D": "4"}

	result := Merge(first, second, StrategyFirst)

	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", result.Conflicts)
	}
	if result.Merged["A"] != "1" || result.Merged["C"] != "3" {
		t.Errorf("unexpected merged values: %v", result.Merged)
	}
}

func TestMerge_ConflictStrategyFirst(t *testing.T) {
	first := map[string]string{"KEY": "from_first"}
	second := map[string]string{"KEY": "from_second"}

	result := Merge(first, second, StrategyFirst)

	if len(result.Conflicts) != 1 || result.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", result.Conflicts)
	}
	if result.Merged["KEY"] != "from_first" {
		t.Errorf("expected 'from_first', got %q", result.Merged["KEY"])
	}
}

func TestMerge_ConflictStrategySecond(t *testing.T) {
	first := map[string]string{"KEY": "from_first"}
	second := map[string]string{"KEY": "from_second"}

	result := Merge(first, second, StrategySecond)

	if result.Merged["KEY"] != "from_second" {
		t.Errorf("expected 'from_second', got %q", result.Merged["KEY"])
	}
}

func TestMerge_UnionKeepsBothUniqueKeys(t *testing.T) {
	first := map[string]string{"A": "1", "SHARED": "x"}
	second := map[string]string{"B": "2", "SHARED": "y"}

	result := Merge(first, second, StrategyUnion)

	if _, ok := result.Merged["A"]; !ok {
		t.Error("expected key A in merged result")
	}
	if _, ok := result.Merged["B"]; !ok {
		t.Error("expected key B in merged result")
	}
	if result.Merged["SHARED"] != "x" {
		t.Errorf("union should keep first value for conflicts, got %q", result.Merged["SHARED"])
	}
}

func TestWriteMerged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "merged.env")

	merged := map[string]string{
		"HOST": "localhost",
		"DESC": "hello world",
		"PORT": "8080",
	}

	if err := WriteMerged(path, merged); err != nil {
		t.Fatalf("WriteMerged failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read merged file: %v", err)
	}

	content := string(data)
	for _, expected := range []string{"HOST=localhost", "PORT=8080", `DESC="hello world"`} {
		if !containsLine(content, expected) {
			t.Errorf("expected line %q in output:\n%s", expected, content)
		}
	}
}

func containsLine(content, line string) bool {
	for _, l := range splitLines(content) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			if i > start {
				lines = append(lines, s[start:i])
			}
			start = i + 1
		}
	}
	return lines
}
