package cmd

import (
	"fmt"
	"strings"

	"github.com/HituziANDO/henge/internal/encoder"
	hengeio "github.com/HituziANDO/henge/internal/io"
	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode <format>",
	Short: "Encode input data (base64, url, hex)",
	Long: `Encode input data into the specified format.

Supported formats:
  base64    Base64 encoding
  url       URL percent-encoding
  hex       Hexadecimal encoding

Examples:
  echo "hello" | henge encode base64
  echo "hello world" | henge encode url
  echo "hello" | henge encode hex`,
}

var encodeBase64Cmd = &cobra.Command{
	Use:   "base64 [input]",
	Short: "Encode input to base64",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := encoder.Base64Encode(input)
		if err != nil {
			return fmt.Errorf("base64 encode failed: %w", err)
		}
		return writeOutput(result)
	},
}

var encodeURLCmd = &cobra.Command{
	Use:   "url [input]",
	Short: "Encode input with URL percent-encoding",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := encoder.URLEncode(input)
		if err != nil {
			return fmt.Errorf("url encode failed: %w", err)
		}
		return writeOutput(result)
	},
}

var encodeHexCmd = &cobra.Command{
	Use:   "hex [input]",
	Short: "Encode input to hexadecimal",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		result, err := encoder.HexEncode(input)
		if err != nil {
			return fmt.Errorf("hex encode failed: %w", err)
		}
		return writeOutput(result)
	},
}

func init() {
	encodeCmd.AddCommand(encodeBase64Cmd)
	encodeCmd.AddCommand(encodeURLCmd)
	encodeCmd.AddCommand(encodeHexCmd)
	rootCmd.AddCommand(encodeCmd)
}
