package diff

// MismatchedKey holds a key whose value differs between two env files.
type MismatchedKey struct {
	Key         string
	FirstValue  string
	SecondValue string
}

// Result holds the full comparison outcome between two env maps.
type Result struct {
	Matching        []string
	MissingInFirst  []string
	MissingInSecond []string
	Mismatched      []MismatchedKey
}

// Compare compares two parsed env maps and returns a Result.
func Compare(first, second map[string]string) Result {
	var result Result

	for key, val1 := range first {
		val2, exists := second[key]
		if !exists {
			result.MissingInSecond = append(result.MissingInSecond, key)
		} else if val1 != val2 {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:         key,
				FirstValue:  val1,
				SecondValue: val2,
			})
		} else {
			result.Matching = append(result.Matching, key)
		}
	}

	for key := range second {
		if _, exists := first[key]; !exists {
			result.MissingInFirst = append(result.MissingInFirst, key)
		}
	}

	result.Matching = sortStrings(result.Matching)
	result.MissingInFirst = sortStrings(result.MissingInFirst)
	result.MissingInSecond = sortStrings(result.MissingInSecond)
	result.Mismatched = sortMismatched(result.Mismatched)

	return result
}
