package diff

import (
	"testing"
)

func TestSummarize_NoDifferences(t *testing.T) {
	r := Result{
		Matching: []string{"A", "B", "C"},
	}
	s := Summarize(r)
	if s.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", s.TotalKeys)
	}
	if s.HasDifferences() {
		t.Error("expected no differences")
	}
	expected := "No differences found (3 keys compared)."
	if s.String() != expected {
		t.Errorf("expected %q, got %q", expected, s.String())
	}
}

func TestSummarize_MissingInSecond(t *testing.T) {
	r := Result{
		Matching:        []string{"A"},
		MissingInSecond: []string{"B", "C"},
	}
	s := Summarize(r)
	if s.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", s.TotalKeys)
	}
	if s.MissingInSecond != 2 {
		t.Errorf("expected MissingInSecond=2, got %d", s.MissingInSecond)
	}
	if !s.HasDifferences() {
		t.Error("expected differences")
	}
}

func TestSummarize_Mismatched(t *testing.T) {
	r := Result{
		Mismatched: []MismatchedKey{
			{Key: "PORT", FirstValue: "8080", SecondValue: "9090"},
		},
	}
	s := Summarize(r)
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
}

func TestSummarize_Mixed(t *testing.T) {
	r := Result{
		Matching:        []string{"A"},
		MissingInFirst:  []string{"B"},
		MissingInSecond: []string{"C"},
		Mismatched: []MismatchedKey{
			{Key: "D", FirstValue: "x", SecondValue: "y"},
		},
	}
	s := Summarize(r)
	if s.TotalKeys != 4 {
		t.Errorf("expected TotalKeys=4, got %d", s.TotalKeys)
	}
	if s.MissingInFirst != 1 {
		t.Errorf("expected MissingInFirst=1, got %d", s.MissingInFirst)
	}
	expected := "4 key(s) compared: 1 missing in first, 1 missing in second, 1 mismatched."
	if s.String() != expected {
		t.Errorf("expected %q, got %q", expected, s.String())
	}
}

func TestSummarize_EmptyResult(t *testing.T) {
	s := Summarize(Result{})
	if s.TotalKeys != 0 {
		t.Errorf("expected TotalKeys=0, got %d", s.TotalKeys)
	}
	if s.HasDifferences() {
		t.Error("expected no differences for empty result")
	}
}
