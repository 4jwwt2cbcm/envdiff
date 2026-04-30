package diff

import "sort"

func sortStrings(s []string) {
	sort.Strings(s)
}

func sortMismatched(m []MismatchedKey) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Key < m[j].Key
	})
}
