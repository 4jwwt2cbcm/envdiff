package diff

import (
	"testing"
	"time"
)

func sampleBaseline() *Baseline {
	return &Baseline{
		CreatedAt:   time.Now().UTC(),
		Environment: "staging",
		Keys: map[string]string{
			"APP_ENV":   "staging",
			"DB_HOST":   "db.staging.internal",
			"LOG_LEVEL": "debug",
		},
	}
}

func TestCompareWithBaseline_NoDrift(t *testing.T) {
	b := sampleBaseline()
	current := map[string]string{
		"APP_ENV":   "staging",
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
	result := CompareWithBaseline(b, current)
	if len(result.MissingInFirst) != 0 || len(result.MissingInSecond) != 0 || len(result.Mismatched) != 0 {
		t.Errorf("expected no drift, got %+v", result)
	}
}

func TestCompareWithBaseline_MissingKey(t *testing.T) {
	b := sampleBaseline()
	current := map[string]string{
		"APP_ENV": "staging",
		// DB_HOST missing
		"LOG_LEVEL": "debug",
	}
	result := CompareWithBaseline(b, current)
	if len(result.MissingInSecond) != 1 || result.MissingInSecond[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST missing in second, got %v", result.MissingInSecond)
	}
}

func TestNewBaselineDriftReport_HasDrift(t *testing.T) {
	b := sampleBaseline()
	current := map[string]string{
		"APP_ENV":   "production", // mismatched
		"DB_HOST":   "db.staging.internal",
		"LOG_LEVEL": "debug",
	}
	report := NewBaselineDriftReport(b, current)
	if report.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %q", report.Environment)
	}
	if !report.HasDrift() {
		t.Error("expected HasDrift() to return true")
	}
}

func TestNewBaselineDriftReport_NoDrift(t *testing.T) {
	b := sampleBaseline()
	report := NewBaselineDriftReport(b, b.Keys)
	if report.HasDrift() {
		t.Errorf("expected no drift, got result: %+v", report.Result)
	}
}
