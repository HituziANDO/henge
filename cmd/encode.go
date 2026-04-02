package cmd

import (
	"fmt"

	"github.com/HituziANDO/henge/internal/encoder"
	hengeimg "github.com/HituziANDO/henge/internal/image"
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
  image     Encode image file to base64

Examples:
  echo "hello" | henge encode base64
  echo "hello world" | henge encode url
  echo "hello" | henge encode hex
  henge encode image --file logo.png --data-uri`,
}

var encodeBase64Cmd = &cobra.Command{
	Use:   "base64 [input]",
	Short: "Encode input to base64",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args, inputFile)
		if err != nil {
			return err
		}
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
		input, err := hengeio.ReadInput(args, inputFile)
		if err != nil {
			return err
		}
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
		input, err := hengeio.ReadInput(args, inputFile)
		if err != nil {
			return err
		}
		result, err := encoder.HexEncode(input)
		if err != nil {
			return fmt.Errorf("hex encode failed: %w", err)
		}
		return writeOutput(result)
	},
}

var encodeImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Encode image file to base64",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputFile == "" {
			return fmt.Errorf("--file flag is required for image encoding (e.g. henge encode image --file logo.png)")
		}
		dataURI, _ := cmd.Flags().GetBool("data-uri")
		wrapWidth, _ := cmd.Flags().GetInt("wrap")

		var result string
		var err error
		if dataURI {
			result, err = hengeimg.EncodeFileToDataURI(inputFile)
		} else {
			result, err = hengeimg.EncodeFileToBase64(inputFile)
		}
		if err != nil {
			return fmt.Errorf("image encode failed: %w", err)
		}

		if wrapWidth > 0 {
			result = hengeimg.WrapString(result, wrapWidth)
		}

		return writeOutput(result)
	},
}

func init() {
	encodeImageCmd.Flags().BoolP("data-uri", "d", false, "output as Data URI scheme")
	encodeImageCmd.Flags().IntP("wrap", "w", 0, "wrap output at specified width (0=no wrap)")

	encodeCmd.AddCommand(encodeBase64Cmd)
	encodeCmd.AddCommand(encodeURLCmd)
	encodeCmd.AddCommand(encodeHexCmd)
	encodeCmd.AddCommand(encodeImageCmd)
	rootCmd.AddCommand(encodeCmd)
}
