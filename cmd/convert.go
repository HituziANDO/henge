package cmd

import (
	"fmt"

	"github.com/HituziANDO/henge/internal/converter"
	hengeio "github.com/HituziANDO/henge/internal/io"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert <target-format>",
	Short: "Convert data between formats (json, yaml, toml)",
	Long: `Convert input data to the specified target format.

Supported target formats: json, yaml, toml
Supported input formats (auto-detected or via --from): json, yaml, toml, csv

Examples:
  cat file.yaml | henge convert json
  cat file.json | henge convert yaml
  cat file.json | henge convert toml
  cat data.csv  | henge convert json`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	target := args[0]

	input, err := hengeio.ReadInput(args[1:])
	if err != nil {
		return err
	}

	var result string
	switch target {
	case "json":
		result, err = converter.ToJSON(input, fromFormat)
	case "yaml":
		result, err = converter.ToYAML(input, fromFormat)
	case "toml":
		result, err = converter.ToTOML(input, fromFormat)
	default:
		return fmt.Errorf("unsupported target format: %s (supported: json, yaml, toml)", target)
	}
	if err != nil {
		return err
	}

	return writeOutput(result)
}
