package diff

import "testing"

// TestSortedKeys verifies that sortedKeys returns a stable sorted copy.
func TestSortedKeys(t *testing.T) {
	input := []string{"Z", "A", "M", "B"}
	result := sortedKeys(input)
	expected := []string{"A", "B", "M", "Z"}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], v)
		}
	}
	// original must not be mutated
	if input[0] != "Z" {
		t.Error("sortedKeys must not mutate the original slice")
	}
}

// TestBuildKeySet_FromOptions verifies key restriction from opts.Keys.
func TestBuildKeySet_FromOptions(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	keys := buildKeySet([]string{"A", "C"}, src)
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

// TestBuildKeySet_FromSrc verifies that all src keys are used when opts.Keys is empty.
func TestBuildKeySet_FromSrc(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	keys := buildKeySet(nil, src)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

// TestPromoteEnv_DstUnchangedWhenNothingPromoted ensures dst is copied faithfully.
func TestPromoteEnv_DstUnchangedWhenNothingPromoted(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "old", "B": "2"}

	result, pr := PromoteEnv(src, dst, PromoteOptions{Overwrite: false})

	if len(pr.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(pr.Skipped))
	}
	if result["A"] != "old" {
		t.Errorf("A should remain 'old', got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("B should remain '2', got %s", result["B"])
	}
}
