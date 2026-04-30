package diff

// Result holds the comparison result between two env files.
type Result struct {
	MissingInSecond []string          // keys present in first but not in second
	MissingInFirst  []string          // keys present in second but not in first
	Mismatched      []MismatchedKey   // keys present in both but with different values
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key    string
	First  string
	Second string
}

// Compare compares two parsed env maps and returns a Result.
func Compare(first, second map[string]string) Result {
	result := Result{}

	for key, val1 := range first {
		val2, ok := second[key]
		if !ok {
			result.MissingInSecond = append(result.MissingInSecond, key)
			continue
		}
		if val1 != val2 {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				First:  val1,
				Second: val2,
			})
		}
	}

	for key := range second {
		if _, ok := first[key]; !ok {
			result.MissingInFirst = append(result.MissingInFirst, key)
		}
	}

	sortStrings(result.MissingInFirst)
	sortStrings(result.MissingInSecond)
	sortMismatched(result.Mismatched)

	return result
}

// HasDifferences returns true if the result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInFirst) > 0 ||
		len(r.MissingInSecond) > 0 ||
		len(r.Mismatched) > 0
}
