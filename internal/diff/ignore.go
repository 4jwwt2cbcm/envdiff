package diff

import (
	"bufio"
	"os"
	"strings"
)

// IgnoreList holds a set of keys that should be excluded from diff results.
type IgnoreList struct {
	keys map[string]struct{}
}

// NewIgnoreList creates an empty IgnoreList.
func NewIgnoreList() *IgnoreList {
	return &IgnoreList{keys: make(map[string]struct{})}
}

// Add inserts a key into the ignore list.
func (il *IgnoreList) Add(key string) {
	il.keys[strings.TrimSpace(key)] = struct{}{}
}

// Contains reports whether the given key is in the ignore list.
func (il *IgnoreList) Contains(key string) bool {
	_, ok := il.keys[strings.TrimSpace(key)]
	return ok
}

// Len returns the number of keys in the ignore list.
func (il *IgnoreList) Len() int {
	return len(il.keys)
}

// LoadIgnoreFile reads a file where each non-blank, non-comment line is a key
// to ignore, and returns a populated IgnoreList.
func LoadIgnoreFile(path string) (*IgnoreList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	il := NewIgnoreList()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		il.Add(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return il, nil
}

// ApplyIgnoreList removes entries whose keys appear in the IgnoreList from a
// Result, returning a new filtered Result.
func ApplyIgnoreList(r Result, il *IgnoreList) Result {
	if il == nil || il.Len() == 0 {
		return r
	}

	filtered := Result{}

	for _, k := range r.MissingInFirst {
		if !il.Contains(k) {
			filtered.MissingInFirst = append(filtered.MissingInFirst, k)
		}
	}
	for _, k := range r.MissingInSecond {
		if !il.Contains(k) {
			filtered.MissingInSecond = append(filtered.MissingInSecond, k)
		}
	}
	for _, m := range r.Mismatched {
		if !il.Contains(m.Key) {
			filtered.Mismatched = append(filtered.Mismatched, m)
		}
	}
	return filtered
}
