package diff_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	first := map[string]string{"FOO": "bar", "BAZ": "qux"}
	second := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(first, second)

	if result.HasDifferences() {
		t.Error("expected no differences, got some")
	}
}

func TestCompare_MissingInSecond(t *testing.T) {
	first := map[string]string{"FOO": "bar", "ONLY_FIRST": "val"}
	second := map[string]string{"FOO": "bar"}

	result := diff.Compare(first, second)

	if len(result.MissingInSecond) != 1 || result.MissingInSecond[0] != "ONLY_FIRST" {
		t.Errorf("expected ONLY_FIRST missing in second, got %v", result.MissingInSecond)
	}
}

func TestCompare_MissingInFirst(t *testing.T) {
	first := map[string]string{"FOO": "bar"}
	second := map[string]string{"FOO": "bar", "ONLY_SECOND": "val"}

	result := diff.Compare(first, second)

	if len(result.MissingInFirst) != 1 || result.MissingInFirst[0] != "ONLY_SECOND" {
		t.Errorf("expected ONLY_SECOND missing in first, got %v", result.MissingInFirst)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	first := map[string]string{"FOO": "bar"}
	second := map[string]string{"FOO": "baz"}

	result := diff.Compare(first, second)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "FOO" || m.First != "bar" || m.Second != "baz" {
		t.Errorf("unexpected mismatch: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	first := map[string]string{"Z_KEY": "v", "A_KEY": "v", "M_KEY": "v"}
	second := map[string]string{}

	result := diff.Compare(first, second)

	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, key := range result.MissingInSecond {
		if key != expected[i] {
			t.Errorf("expected sorted key %s at index %d, got %s", expected[i], i, key)
		}
	}
}

func TestHasDifferences(t *testing.T) {
	empty := diff.Result{}
	if empty.HasDifferences() {
		t.Error("empty result should have no differences")
	}

	withMissing := diff.Result{MissingInFirst: []string{"KEY"}}
	if !withMissing.HasDifferences() {
		t.Error("result with missing keys should have differences")
	}
}
