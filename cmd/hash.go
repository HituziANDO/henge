package cmd

import (
	"strings"

	"github.com/henge-cli/henge/internal/hasher"
	hengeio "github.com/henge-cli/henge/internal/io"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash <algorithm>",
	Short: "Hash input data with the specified algorithm",
	Long:  `Compute a cryptographic hash of the input using md5, sha1, sha256, or sha512.`,
}

var md5Cmd = &cobra.Command{
	Use:   "md5 [input]",
	Short: "Compute MD5 hash",
	RunE:  runHash(hasher.MD5Hash),
}

var sha1Cmd = &cobra.Command{
	Use:   "sha1 [input]",
	Short: "Compute SHA-1 hash",
	RunE:  runHash(hasher.SHA1Hash),
}

var sha256Cmd = &cobra.Command{
	Use:   "sha256 [input]",
	Short: "Compute SHA-256 hash",
	RunE:  runHash(hasher.SHA256Hash),
}

var sha512Cmd = &cobra.Command{
	Use:   "sha512 [input]",
	Short: "Compute SHA-512 hash",
	RunE:  runHash(hasher.SHA512Hash),
}

func init() {
	hashCmd.AddCommand(md5Cmd, sha1Cmd, sha256Cmd, sha512Cmd)
	rootCmd.AddCommand(hashCmd)
}

func runHash(hashFn func(string) string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args)
		if err != nil {
			return err
		}

		// Remove trailing newline from stdin/pipe input
		input = strings.TrimRight(input, "\n")

		result := hashFn(input)
		return writeOutput(result)
	}
}
