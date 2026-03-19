package converter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// DetectFormat detects the format of the input string.
func DetectFormat(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "unknown"
	}

	// Try JSON
	var js json.RawMessage
	if json.Unmarshal([]byte(trimmed), &js) == nil {
		return "json"
	}

	// Try TOML (before YAML since YAML is very permissive)
	var tomlVal interface{}
	if _, err := toml.Decode(trimmed, &tomlVal); err == nil && tomlVal != nil {
		// TOML should have key=value pairs; check for '=' to avoid false positives
		if strings.Contains(trimmed, "=") {
			return "toml"
		}
	}

	// Try CSV: must have multiple lines and commas
	if strings.Contains(trimmed, ",") && strings.Contains(trimmed, "\n") {
		r := csv.NewReader(strings.NewReader(trimmed))
		records, err := r.ReadAll()
		if err == nil && len(records) >= 2 {
			// All rows should have same number of fields
			cols := len(records[0])
			if cols >= 2 {
				consistent := true
				for _, row := range records[1:] {
					if len(row) != cols {
						consistent = false
						break
					}
				}
				if consistent {
					return "csv"
				}
			}
		}
	}

	// Try YAML
	var yamlVal interface{}
	if err := yaml.Unmarshal([]byte(trimmed), &yamlVal); err == nil && yamlVal != nil {
		if strings.Contains(trimmed, ":") {
			return "yaml"
		}
	}

	return "unknown"
}

// ToJSON converts input to JSON from the specified or auto-detected format.
func ToJSON(input string, fromFormat string) (string, error) {
	if fromFormat == "" {
		fromFormat = DetectFormat(input)
	}
	data, err := parseInput(input, fromFormat)
	if err != nil {
		return "", err
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(out), nil
}

// ToYAML converts input to YAML from the specified or auto-detected format.
func ToYAML(input string, fromFormat string) (string, error) {
	if fromFormat == "" {
		fromFormat = DetectFormat(input)
	}
	data, err := parseInput(input, fromFormat)
	if err != nil {
		return "", err
	}
	out, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshaling YAML: %w", err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// ToTOML converts input to TOML from the specified or auto-detected format.
func ToTOML(input string, fromFormat string) (string, error) {
	if fromFormat == "" {
		fromFormat = DetectFormat(input)
	}
	data, err := parseInput(input, fromFormat)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return "", fmt.Errorf("marshaling TOML: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// parseInput parses input string according to the given format.
func parseInput(input string, format string) (interface{}, error) {
	switch format {
	case "json":
		return jsonToInterface(input)
	case "yaml":
		return yamlToInterface(input)
	case "toml":
		return tomlToInterface(input)
	case "csv":
		return csvToInterface(input)
	default:
		return nil, fmt.Errorf("unsupported input format: %s", format)
	}
}

func jsonToInterface(input string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return data, nil
}

func yamlToInterface(input string) (interface{}, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(input), &data); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	return data, nil
}

func tomlToInterface(input string) (interface{}, error) {
	var data interface{}
	if _, err := toml.Decode(input, &data); err != nil {
		return nil, fmt.Errorf("parsing TOML: %w", err)
	}
	return data, nil
}

func csvToInterface(input string) (interface{}, error) {
	r := csv.NewReader(strings.NewReader(strings.TrimSpace(input)))
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must have a header row and at least one data row")
	}

	headers := records[0]
	var result []map[string]string
	for _, row := range records[1:] {
		obj := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				obj[header] = row[i]
			}
		}
		result = append(result, obj)
	}
	return result, nil
}
