package cmd

import (
	"github.com/henge-cli/henge/internal/detector"
	hengeio "github.com/henge-cli/henge/internal/io"
	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto [input]",
	Short: "Auto-detect input format and transform (alias for default behavior)",
	Long:  `Automatically detect the input data format and apply the most appropriate transformation.`,
	RunE:  runAuto,
}

func init() {
	rootCmd.AddCommand(autoCmd)
}

func runAuto(cmd *cobra.Command, args []string) error {
	input, err := hengeio.ReadInput(args)
	if err != nil {
		return err
	}

	result, err := detector.AutoDetectAndTransform(input)
	if err != nil {
		return err
	}

	return writeOutput(result)
}
