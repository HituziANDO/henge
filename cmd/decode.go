package cmd

import (
	"fmt"
	"strings"

	"github.com/henge-cli/henge/internal/decoder"
	hengeio "github.com/henge-cli/henge/internal/io"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode <format>",
	Short: "Decode input data (base64, url, hex)",
	Long: `Decode input data from the specified format.

Supported formats:
  base64    Base64 decoding
  url       URL percent-decoding
  hex       Hexadecimal decoding

Examples:
  echo "aGVsbG8=" | henge decode base64
  echo "hello%20world" | henge decode url
  echo "68656c6c6f" | henge decode hex`,
}

var decodeBase64Cmd = &cobra.Command{
	Use:   "base64 [input]",
	Short: "Decode base64 input",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := decoder.Base64Decode(input)
		if err != nil {
			return fmt.Errorf("base64 decode failed: %w", err)
		}
		return writeOutput(result)
	},
}

var decodeURLCmd = &cobra.Command{
	Use:   "url [input]",
	Short: "Decode URL percent-encoded input",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := decoder.URLDecode(input)
		if err != nil {
			return fmt.Errorf("url decode failed: %w", err)
		}
		return writeOutput(result)
	},
}

var decodeHexCmd = &cobra.Command{
	Use:   "hex [input]",
	Short: "Decode hexadecimal input",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := decoder.HexDecode(input)
		if err != nil {
			return fmt.Errorf("hex decode failed: %w", err)
		}
		return writeOutput(result)
	},
}

func init() {
	decodeCmd.AddCommand(decodeBase64Cmd)
	decodeCmd.AddCommand(decodeURLCmd)
	decodeCmd.AddCommand(decodeHexCmd)
	rootCmd.AddCommand(decodeCmd)
}
