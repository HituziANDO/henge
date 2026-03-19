package io

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadInput reads data from args (file path), stdin, or returns an error.
func ReadInput(args []string) (string, error) {
	if len(args) > 0 {
		// Try to read as file first
		data, err := os.ReadFile(args[0])
		if err == nil {
			return string(data), nil
		}
		// If not a file, treat as literal input
		return strings.Join(args, " "), nil
	}

	// Read from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no input provided. Pipe data or provide a file argument")
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}

	return string(data), nil
}
