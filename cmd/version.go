package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of henge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("henge version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
