package diff

// CompareWithBaseline compares a current env map against a saved Baseline
// and returns a Result describing any differences.
func CompareWithBaseline(baseline *Baseline, current map[string]string) Result {
	return Compare(baseline.Keys, current)
}

// BaselineDriftReport builds a human-readable report describing drift
// between a baseline and the current environment map.
type BaselineDriftReport struct {
	Environment string
	Result      Result
}

// NewBaselineDriftReport creates a BaselineDriftReport from a baseline and current keys.
func NewBaselineDriftReport(baseline *Baseline, current map[string]string) BaselineDriftReport {
	return BaselineDriftReport{
		Environment: baseline.Environment,
		Result:      CompareWithBaseline(baseline, current),
	}
}

// HasDrift returns true if there are any differences from the baseline.
func (r BaselineDriftReport) HasDrift() bool {
	return len(r.Result.MissingInSecond) > 0 ||
		len(r.Result.MissingInFirst) > 0 ||
		len(r.Result.Mismatched) > 0
}
