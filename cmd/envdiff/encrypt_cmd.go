package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// runEncrypt reads a .env file, encrypts all values, and writes JSON to stdout or a file.
// Usage: envdiff encrypt <file> --passphrase <pass> [--output <out.json>]
func runEncrypt(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff encrypt <file> --passphrase <pass> [--output <out.json>]")
		os.Exit(2)
	}

	envFile := args[0]
	passphrase := ""
	outputFile := ""

	for i := 1; i < len(args)-1; i++ {
		switch args[i] {
		case "--passphrase":
			passphrase = args[i+1]
		case "--output":
			outputFile = args[i+1]
		}
	}

	if passphrase == "" {
		fmt.Fprintln(os.Stderr, "encrypt: --passphrase is required")
		os.Exit(2)
	}

	env, err := parser.ParseFile(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encrypt: parse error: %v\n", err)
		os.Exit(1)
	}

	enc, err := diff.EncryptEnv(env, diff.EncryptOptions{Passphrase: passphrase})
	if err != nil {
		fmt.Fprintf(os.Stderr, "encrypt: %v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(enc, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encrypt: marshal error: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, data, 0600); err != nil {
			fmt.Fprintf(os.Stderr, "encrypt: write error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("encrypted env written to %s\n", outputFile)
	} else {
		fmt.Println(string(data))
	}
}

// runDecrypt reads an encrypted JSON file and prints decrypted env to stdout.
// Usage: envdiff decrypt <encrypted.json> --passphrase <pass>
func runDecrypt(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff decrypt <encrypted.json> --passphrase <pass>")
		os.Exit(2)
	}

	encFile := args[0]
	passphrase := ""

	for i := 1; i < len(args)-1; i++ {
		if args[i] == "--passphrase" {
			passphrase = args[i+1]
		}
	}

	if passphrase == "" {
		fmt.Fprintln(os.Stderr, "decrypt: --passphrase is required")
		os.Exit(2)
	}

	raw, err := os.ReadFile(encFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decrypt: read error: %v\n", err)
		os.Exit(1)
	}

	var enc diff.EncryptedEnv
	if err := json.Unmarshal(raw, &enc); err != nil {
		fmt.Fprintf(os.Stderr, "decrypt: invalid JSON: %v\n", err)
		os.Exit(1)
	}

	env, err := diff.DecryptEnv(enc, diff.EncryptOptions{Passphrase: passphrase})
	if err != nil {
		fmt.Fprintf(os.Stderr, "decrypt: %v\n", err)
		os.Exit(1)
	}

	for k, v := range env {
		fmt.Printf("%s=%s\n", k, v)
	}
}
