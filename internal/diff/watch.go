package diff

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the file watcher.
type WatchOptions struct {
	Interval  time.Duration
	OnChange  func(file string, result *Result)
	OnError   func(file string, err error)
}

// fileHash returns the MD5 hash of a file's contents.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// WatchFiles monitors two .env files for changes and calls OnChange
// with a fresh diff result whenever either file is modified.
// It blocks until the done channel is closed.
func WatchFiles(fileA, fileB string, opts WatchOptions, done <-chan struct{}) {
	if opts.Interval <= 0 {
		opts.Interval = 2 * time.Second
	}

	hashes := map[string]string{fileA: "", fileB: ""}

	check := func() {
		changed := false
		for _, f := range []string{fileA, fileB} {
			h, err := fileHash(f)
			if err != nil {
				if opts.OnError != nil {
					opts.OnError(f, err)
				}
				return
			}
			if h != hashes[f] {
				hashes[f] = h
				changed = true
			}
		}
		if !changed {
			return
		}
		envA, err := parseEnvMap(fileA)
		if err != nil {
			if opts.OnError != nil {
				opts.OnError(fileA, err)
			}
			return
		}
		envB, err := parseEnvMap(fileB)
		if err != nil {
			if opts.OnError != nil {
				opts.OnError(fileB, err)
			}
			return
		}
		result := Compare(envA, envB)
		if opts.OnChange != nil {
			opts.OnChange(fileA, result)
		}
	}

	// Run once immediately.
	check()

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			check()
		case <-done:
			return
		}
	}
}

// parseEnvMap is a thin wrapper used internally by the watcher.
func parseEnvMap(path string) (map[string]string, error) {
	env, err := ParseEnvFile(path)
	if err != nil {
		return nil, fmt.Errorf("watch: parsing %s: %w", path, err)
	}
	return env, nil
}
