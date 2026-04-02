package cmd

import (
	"github.com/HituziANDO/henge/internal/formatter"
	hengeio "github.com/HituziANDO/henge/internal/io"
	"github.com/spf13/cobra"
)

var jsonIndent int

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format data (json, yaml, xml)",
	Long:  `Format and pretty-print structured data such as JSON, YAML, or XML.`,
}

var formatJSONCmd = &cobra.Command{
	Use:   "json [input]",
	Short: "Format JSON data",
	Long:  `Pretty-print JSON data with configurable indentation.`,
	RunE:  runFormatJSON,
}

var formatYAMLCmd = &cobra.Command{
	Use:   "yaml [input]",
	Short: "Format YAML data",
	Long:  `Normalize and pretty-print YAML data.`,
	RunE:  runFormatYAML,
}

var formatXMLCmd = &cobra.Command{
	Use:   "xml [input]",
	Short: "Format XML data",
	Long:  `Pretty-print XML data with indentation.`,
	RunE:  runFormatXML,
}

func init() {
	formatJSONCmd.Flags().IntVar(&jsonIndent, "indent", 2, "number of spaces for indentation")
	formatCmd.AddCommand(formatJSONCmd)
	formatCmd.AddCommand(formatYAMLCmd)
	formatCmd.AddCommand(formatXMLCmd)
	rootCmd.AddCommand(formatCmd)
}

func runFormatJSON(cmd *cobra.Command, args []string) error {
	input, err := hengeio.ReadInput(args, inputFile)
	if err != nil {
		return err
	}

	var result string
	if compact {
		result, err = formatter.CompactJSON(input)
	} else {
		result, err = formatter.FormatJSON(input, jsonIndent)
	}
	if err != nil {
		return err
	}

	return writeOutput(result)
}

func runFormatYAML(cmd *cobra.Command, args []string) error {
	input, err := hengeio.ReadInput(args, inputFile)
	if err != nil {
		return err
	}

	result, err := formatter.FormatYAML(input)
	if err != nil {
		return err
	}

	return writeOutput(result)
}

func runFormatXML(cmd *cobra.Command, args []string) error {
	input, err := hengeio.ReadInput(args, inputFile)
	if err != nil {
		return err
	}

	result, err := formatter.FormatXML(input)
	if err != nil {
		return err
	}

	return writeOutput(result)
}
