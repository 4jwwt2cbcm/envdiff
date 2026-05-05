package diff

// Mismatch represents a key that exists in both env files but has different values.
type Mismatch struct {
	Key    string
	First  string
	Second string
}

// Result holds the full comparison outcome between two env files.
type Result struct {
	// MissingInSecond contains keys present in the first file but absent in the second.
	MissingInSecond []string
	// MissingInFirst contains keys present in the second file but absent in the first.
	MissingInFirst []string
	// Mismatched contains keys present in both files but with differing values.
	Mismatched []Mismatch
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInFirst) > 0 ||
		len(r.MissingInSecond) > 0 ||
		len(r.Mismatched) > 0
}

// Compare compares two maps of env key/value pairs and returns a Result
// describing the differences between them.
func Compare(first, second map[string]string) Result {
	var result Result

	for k, v1 := range first {
		if v2, ok := second[k]; !ok {
			result.MissingInSecond = append(result.MissingInSecond, k)
		} else if v1 != v2 {
			result.Mismatched = append(result.Mismatched, Mismatch{
				Key:    k,
				First:  v1,
				Second: v2,
			})
		}
	}

	for k := range second {
		if _, ok := first[k]; !ok {
			result.MissingInFirst = append(result.MissingInFirst, k)
		}
	}

	sortStrings(result.MissingInSecond)
	sortStrings(result.MissingInFirst)
	sortMismatched(result.Mismatched)

	return result
}
