package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/yourorg/envdiff/internal/diff"
)

// runWatch starts a live watch on two .env files and prints diffs on change.
// Usage: envdiff watch <fileA> <fileB> [interval_seconds]
func runWatch(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("watch requires two file arguments")
	}
	fileA, fileB := args[0], args[1]

	interval := 2 * time.Second
	if len(args) >= 3 {
		secs, err := strconv.ParseFloat(args[2], 64)
		if err != nil || secs <= 0 {
			return fmt.Errorf("invalid interval %q: must be a positive number", args[2])
		}
		interval = time.Duration(secs * float64(time.Second))
	}

	fmt.Fprintf(os.Stderr, "Watching %s and %s (interval: %s). Press Ctrl+C to stop.\n",
		fileA, fileB, interval)

	done := make(chan struct{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nStopping watcher.")
		close(done)
	}()

	opts := diff.WatchOptions{
		Interval: interval,
		OnChange: func(file string, result *diff.Result) {
			timestamp := time.Now().Format(time.RFC3339)
			if len(result.MissingInSecond) == 0 &&
				len(result.MissingInFirst) == 0 &&
				len(result.Mismatched) == 0 {
				fmt.Printf("[%s] No differences.\n", timestamp)
				return
			}
			report := diff.Report(result, false)
			fmt.Printf("[%s] Change detected:\n%s\n", timestamp, report)
		},
		OnError: func(file string, err error) {
			fmt.Fprintf(os.Stderr, "[watch error] %s: %v\n", file, err)
		},
	}

	diff.WatchFiles(fileA, fileB, opts, done)
	return nil
}
