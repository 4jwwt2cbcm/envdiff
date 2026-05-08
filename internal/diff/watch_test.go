package diff

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func writeWatchEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeWatchEnv: %v", err)
	}
	return p
}

func TestWatchFiles_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	fileA := writeWatchEnv(t, dir, ".env.a", "KEY=value1\n")
	fileB := writeWatchEnv(t, dir, ".env.b", "KEY=value1\n")

	var mu sync.Mutex
	var results []*Result

	done := make(chan struct{})
	opts := WatchOptions{
		Interval: 50 * time.Millisecond,
		OnChange: func(_ string, r *Result) {
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		},
	}

	go WatchFiles(fileA, fileB, opts, done)

	// Allow initial check to settle.
	time.Sleep(80 * time.Millisecond)

	// Mutate fileB to trigger a change.
	if err := os.WriteFile(fileB, []byte("KEY=value2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	time.Sleep(150 * time.Millisecond)
	close(done)

	mu.Lock()
	defer mu.Unlock()
	if len(results) < 2 {
		t.Fatalf("expected at least 2 OnChange calls (initial + change), got %d", len(results))
	}
	last := results[len(results)-1]
	if len(last.Mismatched) == 0 {
		t.Errorf("expected mismatch on KEY, got none")
	}
}

func TestWatchFiles_NoSpuriousFires(t *testing.T) {
	dir := t.TempDir()
	fileA := writeWatchEnv(t, dir, ".env.a", "FOO=bar\n")
	fileB := writeWatchEnv(t, dir, ".env.b", "FOO=bar\n")

	var mu sync.Mutex
	callCount := 0

	done := make(chan struct{})
	opts := WatchOptions{
		Interval: 40 * time.Millisecond,
		OnChange: func(_ string, _ *Result) {
			mu.Lock()
			callCount++
			mu.Unlock()
		},
	}

	go WatchFiles(fileA, fileB, opts, done)
	time.Sleep(200 * time.Millisecond)
	close(done)

	mu.Lock()
	defer mu.Unlock()
	// Only the initial call should have fired (files unchanged after that).
	if callCount > 1 {
		t.Errorf("expected at most 1 OnChange call, got %d", callCount)
	}
}

func TestWatchFiles_ErrorOnMissingFile(t *testing.T) {
	dir := t.TempDir()
	fileA := writeWatchEnv(t, dir, ".env.a", "KEY=val\n")
	fileB := filepath.Join(dir, "nonexistent.env")

	var mu sync.Mutex
	var errs []error

	done := make(chan struct{})
	opts := WatchOptions{
		Interval: 40 * time.Millisecond,
		OnError: func(_ string, err error) {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		},
	}

	go WatchFiles(fileA, fileB, opts, done)
	time.Sleep(100 * time.Millisecond)
	close(done)

	mu.Lock()
	defer mu.Unlock()
	if len(errs) == 0 {
		t.Error("expected at least one error for missing file")
	}
}
