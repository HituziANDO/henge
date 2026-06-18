// Package strcase converts the letter case and identifier case of text.
//
// ToUpper and ToLower operate per character (Unicode-aware) and preserve all
// delimiters, whitespace, and structure. ToSnake, ToCamel, ToKebab, and
// ToPascal tokenize the input — splitting on '_', '-', whitespace, and
// camelCase boundaries (including acronym runs such as "HTTPServer") — and
// re-emit it in the target case. Identifier conversions process each line
// independently, treating one line as one identifier.
package strcase

import (
	"strings"
	"unicode"
)

// ToUpper returns the input with every letter uppercased. Delimiters,
// whitespace, and other characters are preserved (Unicode-aware).
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower returns the input with every letter lowercased. Delimiters,
// whitespace, and other characters are preserved (Unicode-aware).
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToSnake converts each line of the input to snake_case.
func ToSnake(s string) string {
	return mapLines(s, func(line string) string {
		return strings.Join(tokenize(line), "_")
	})
}

// ToKebab converts each line of the input to kebab-case.
func ToKebab(s string) string {
	return mapLines(s, func(line string) string {
		return strings.Join(tokenize(line), "-")
	})
}

// ToCamel converts each line of the input to camelCase.
func ToCamel(s string) string {
	return mapLines(s, func(line string) string {
		tokens := tokenize(line)
		for i := 1; i < len(tokens); i++ {
			tokens[i] = titleToken(tokens[i])
		}
		return strings.Join(tokens, "")
	})
}

// ToPascal converts each line of the input to PascalCase.
func ToPascal(s string) string {
	return mapLines(s, func(line string) string {
		tokens := tokenize(line)
		for i := range tokens {
			tokens[i] = titleToken(tokens[i])
		}
		return strings.Join(tokens, "")
	})
}

// mapLines applies fn to each newline-separated line, preserving line breaks.
func mapLines(s string, fn func(string) string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = fn(line)
	}
	return strings.Join(lines, "\n")
}

// titleToken uppercases the first rune of an already-lowercased token.
func titleToken(tok string) string {
	if tok == "" {
		return ""
	}
	r := []rune(tok)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// isDelimiter reports whether r separates words in an identifier.
func isDelimiter(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

// tokenize splits a single string into lowercased words. It breaks on
// delimiters ('_', '-', whitespace) and camelCase boundaries, also separating
// acronym runs (e.g. "HTTPServer" -> "http", "server"; "parseURL" -> "parse",
// "url"). Digits stay attached to the preceding word. Empty tokens produced by
// consecutive or leading/trailing delimiters are dropped.
func tokenize(s string) []string {
	runes := []rune(s)
	var tokens []string
	var cur []rune

	flush := func() {
		if len(cur) > 0 {
			tokens = append(tokens, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case isDelimiter(r):
			flush()
		case unicode.IsUpper(r):
			if len(cur) > 0 {
				prev := cur[len(cur)-1]
				switch {
				case unicode.IsLower(prev) || unicode.IsDigit(prev):
					// camelCase boundary: fooBar -> foo | Bar
					flush()
				case unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]):
					// acronym boundary: HTTPServer -> HTTP | Server
					flush()
				}
			}
			cur = append(cur, r)
		default:
			cur = append(cur, r)
		}
	}
	flush()

	return tokens
}
