package formatter

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// FormatJSON pretty-prints JSON with the given indent (number of spaces).
// If indent is 0, it defaults to 2 spaces.
func FormatJSON(input string, indent int) (string, error) {
	if indent == 0 {
		indent = 2
	}

	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	out, err := json.MarshalIndent(obj, "", strings.Repeat(" ", indent))
	if err != nil {
		return "", fmt.Errorf("formatting JSON: %w", err)
	}

	return string(out), nil
}

// CompactJSON returns compact (no whitespace) JSON.
func CompactJSON(input string) (string, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	return buf.String(), nil
}

// FormatYAML parses and re-serializes YAML to normalize formatting.
func FormatYAML(input string) (string, error) {
	var obj interface{}
	if err := yaml.Unmarshal([]byte(input), &obj); err != nil {
		return "", fmt.Errorf("invalid YAML: %w", err)
	}

	out, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("formatting YAML: %w", err)
	}

	return strings.TrimRight(string(out), "\n"), nil
}

// FormatXML pretty-prints XML with indentation.
func FormatXML(input string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(input))
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("invalid XML: %w", err)
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", fmt.Errorf("formatting XML: %w", err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", fmt.Errorf("formatting XML: %w", err)
	}

	return buf.String(), nil
}
