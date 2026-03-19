package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "0.1.0"
	outputFile string
	fromFormat string
	compact    bool
	noNewline  bool
)

var rootCmd = &cobra.Command{
	Use:   "henge [file]",
	Short: "henge - universal data transformation tool",
	Long: `henge (変化) - CyberChef for Terminal

A universal CLI tool for data transformation:
  - Encode/Decode: base64, url, hex
  - Hash: md5, sha1, sha256, sha512
  - Format: json, yaml, xml
  - Convert: json↔yaml, json↔toml, csv→json
  - Auto-detect: intelligent format detection

Usage:
  echo "aGVsbG8=" | henge              # Auto-detect and transform
  echo "hello" | henge encode base64   # Explicit encoding
  cat file.json | henge format json    # Format JSON`,
	Version: version,
	RunE:    runAuto,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "output file (default: stdout)")
	rootCmd.PersistentFlags().StringVarP(&fromFormat, "from", "f", "", "input format (override auto-detection)")
	rootCmd.PersistentFlags().BoolVarP(&compact, "compact", "c", false, "compact output (no indentation)")
	rootCmd.PersistentFlags().BoolVarP(&noNewline, "no-newline", "n", false, "do not append newline to output")
}

// writeOutput handles writing result to stdout or file, respecting global flags.
func writeOutput(result string) error {
	if !noNewline && len(result) > 0 && result[len(result)-1] != '\n' {
		result += "\n"
	}

	if outputFile != "" {
		return os.WriteFile(outputFile, []byte(result), 0644)
	}

	fmt.Print(result)
	return nil
}
