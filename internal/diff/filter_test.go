package diff

import (
	"testing"
)

func TestFilter_NoOptions(t *testing.T) {
	r := Result{
		MissingInFirst:  []string{"A"},
		MissingInSecond: []string{"B"},
		Mismatched:      []Mismatch{{Key: "C", First: "1", Second: "2"}},
	}
	out := Filter(r, FilterOptions{})
	if len(out.MissingInFirst) != 1 || out.MissingInFirst[0] != "A" {
		t.Errorf("expected MissingInFirst=[A], got %v", out.MissingInFirst)
	}
	if len(out.MissingInSecond) != 1 || out.MissingInSecond[0] != "B" {
		t.Errorf("expected MissingInSecond=[B], got %v", out.MissingInSecond)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "C" {
		t.Errorf("expected Mismatched=[C], got %v", out.Mismatched)
	}
}

func TestFilter_OnlyMissing(t *testing.T) {
	r := Result{
		MissingInFirst:  []string{"A"},
		MissingInSecond: []string{"B"},
		Mismatched:      []Mismatch{{Key: "C", First: "1", Second: "2"}},
	}
	out := Filter(r, FilterOptions{OnlyMissing: true})
	if len(out.Mismatched) != 0 {
		t.Errorf("expected no mismatches, got %v", out.Mismatched)
	}
	if len(out.MissingInFirst) != 1 {
		t.Errorf("expected MissingInFirst to have 1 entry, got %v", out.MissingInFirst)
	}
}

func TestFilter_OnlyMismatched(t *testing.T) {
	r := Result{
		MissingInFirst:  []string{"A"},
		MissingInSecond: []string{"B"},
		Mismatched:      []Mismatch{{Key: "C", First: "1", Second: "2"}},
	}
	out := Filter(r, FilterOptions{OnlyMismatched: true})
	if len(out.MissingInFirst) != 0 || len(out.MissingInSecond) != 0 {
		t.Errorf("expected no missing keys, got first=%v second=%v", out.MissingInFirst, out.MissingInSecond)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "C" {
		t.Errorf("expected Mismatched=[C], got %v", out.Mismatched)
	}
}

func TestFilter_KeyRestriction(t *testing.T) {
	r := Result{
		MissingInSecond: []string{"A", "B", "D"},
		Mismatched:      []Mismatch{{Key: "C", First: "1", Second: "2"}, {Key: "E", First: "x", Second: "y"}},
	}
	out := Filter(r, FilterOptions{Keys: []string{"A", "C"}})
	if len(out.MissingInSecond) != 1 || out.MissingInSecond[0] != "A" {
		t.Errorf("expected MissingInSecond=[A], got %v", out.MissingInSecond)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "C" {
		t.Errorf("expected Mismatched=[C], got %v", out.Mismatched)
	}
}

func TestFilter_EmptyResult(t *testing.T) {
	out := Filter(Result{}, FilterOptions{OnlyMissing: true, Keys: []string{"X"}})
	if len(out.MissingInFirst) != 0 || len(out.MissingInSecond) != 0 || len(out.Mismatched) != 0 {
		t.Errorf("expected empty result, got %+v", out)
	}
}
