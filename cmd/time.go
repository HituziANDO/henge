package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HituziANDO/henge/internal/timconv"
	hengeio "github.com/HituziANDO/henge/internal/io"
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time <subcommand>",
	Short: "Convert between UNIX timestamps and date strings",
	Long: `Convert between UNIX timestamps and date strings.

Supported subcommands:
  unix    Convert date string to UNIX timestamp
  date    Convert UNIX timestamp to date string

Examples:
  henge time unix "2025-01-01T00:00:00Z"
  henge time date 1735689600
  henge time date --timezone Asia/Tokyo 1735689600`,
}

var timeUnixCmd = &cobra.Command{
	Use:   "unix [input]",
	Short: "Convert date string to UNIX timestamp",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		tzName, _ := cmd.Flags().GetString("timezone")
		millis, _ := cmd.Flags().GetBool("millis")

		timestamp, err := timconv.DateToUnix(input, tzName)
		if err != nil {
			return fmt.Errorf("time unix failed: %w", err)
		}

		if millis {
			timestamp = timestamp * 1000
		}

		return writeOutput(strconv.FormatInt(timestamp, 10))
	},
}

var timeDateCmd = &cobra.Command{
	Use:   "date [input]",
	Short: "Convert UNIX timestamp to date string",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		format, _ := cmd.Flags().GetString("format")
		tzName, _ := cmd.Flags().GetString("timezone")
		forceMillis, _ := cmd.Flags().GetBool("millis")

		timestamp, err := timconv.ParseUnixTimestamp(input, forceMillis)
		if err != nil {
			return fmt.Errorf("time date failed: %w", err)
		}

		result, err := timconv.UnixToDate(timestamp, format, tzName)
		if err != nil {
			return fmt.Errorf("time date failed: %w", err)
		}

		return writeOutput(result)
	},
}

func init() {
	timeUnixCmd.Flags().BoolP("millis", "m", false, "output millisecond UNIX timestamp")
	timeUnixCmd.Flags().StringP("timezone", "z", "", "timezone for inputs without timezone info (e.g. Asia/Tokyo)")

	timeDateCmd.Flags().StringP("format", "F", "rfc3339", "output date format (preset name or Go layout)")
	timeDateCmd.Flags().StringP("timezone", "z", "UTC", "output timezone (e.g. Asia/Tokyo, Local)")
	timeDateCmd.Flags().BoolP("millis", "m", false, "force interpret input as millisecond timestamp")

	timeCmd.AddCommand(timeUnixCmd, timeDateCmd)
	rootCmd.AddCommand(timeCmd)
}
