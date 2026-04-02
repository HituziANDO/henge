package io

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadInput reads data from a file (if filePath is set), args (literal), or stdin.
func ReadInput(args []string, filePath string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}

	if len(args) > 0 {
		// Treat arguments as literal input
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
