package cmd

import (
	"strings"

	hengeio "github.com/HituziANDO/henge/internal/io"
	"github.com/HituziANDO/henge/internal/strcase"
	"github.com/spf13/cobra"
)

var caseCmd = &cobra.Command{
	Use:   "case <type>",
	Short: "Convert the case of input text (upper, lower, snake, camel, kebab, pascal)",
	Long: `Convert the letter case or identifier case of the input text.

Supported types:
  upper     Uppercase every letter (delimiters preserved)
  lower     Lowercase every letter (delimiters preserved)
  snake     snake_case
  camel     camelCase
  kebab     kebab-case
  pascal    PascalCase

snake, camel, kebab, and pascal tokenize the input — splitting on _, -,
whitespace, and camelCase boundaries (including acronym runs) — and re-emit it
in the target case, processing each line as one identifier.

Examples:
  echo "foo_bar Baz" | henge case upper    # FOO_BAR BAZ
  echo "FOO-BAR"      | henge case lower    # foo-bar
  echo "fooBar"       | henge case snake    # foo_bar
  echo "foo_bar"      | henge case camel    # fooBar
  echo "fooBar"       | henge case kebab    # foo-bar
  echo "foo_bar"      | henge case pascal   # FooBar`,
}

var caseUpperCmd = &cobra.Command{
	Use:   "upper [input]",
	Short: "Uppercase every letter (delimiters preserved)",
	RunE:  runCase(strcase.ToUpper),
}

var caseLowerCmd = &cobra.Command{
	Use:   "lower [input]",
	Short: "Lowercase every letter (delimiters preserved)",
	RunE:  runCase(strcase.ToLower),
}

var caseSnakeCmd = &cobra.Command{
	Use:   "snake [input]",
	Short: "Convert to snake_case",
	RunE:  runCase(strcase.ToSnake),
}

var caseCamelCmd = &cobra.Command{
	Use:   "camel [input]",
	Short: "Convert to camelCase",
	RunE:  runCase(strcase.ToCamel),
}

var caseKebabCmd = &cobra.Command{
	Use:   "kebab [input]",
	Short: "Convert to kebab-case",
	RunE:  runCase(strcase.ToKebab),
}

var casePascalCmd = &cobra.Command{
	Use:   "pascal [input]",
	Short: "Convert to PascalCase",
	RunE:  runCase(strcase.ToPascal),
}

func init() {
	caseCmd.AddCommand(caseUpperCmd, caseLowerCmd, caseSnakeCmd, caseCamelCmd, caseKebabCmd, casePascalCmd)
	rootCmd.AddCommand(caseCmd)
}

func runCase(convertFn func(string) string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		input, err := hengeio.ReadInput(args, inputFile)
		if err != nil {
			return err
		}

		// Remove trailing newline (including CRLF) from stdin/pipe input so it
		// does not produce a spurious empty final line.
		input = strings.TrimRight(input, "\r\n")

		result := convertFn(input)
		return writeOutput(result)
	}
}
