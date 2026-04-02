package detector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/HituziANDO/henge/internal/decoder"
	"github.com/HituziANDO/henge/internal/timconv"
	"gopkg.in/yaml.v3"
)

// AutoDetectAndTransform detects the input format and applies the best transformation.
func AutoDetectAndTransform(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	// 1. UNIX timestamp → date string (before JSON, since bare numbers are valid JSON)
	if isTimestamp(input) {
		result, err := timconv.AutoConvert(input)
		if err == nil {
			return result, nil
		}
	}

	// 2. JSON → pretty print
	if isJSON(input) {
		return formatJSON(input)
	}

	// 3. Base64 → decode
	if isBase64(input) {
		decoded, err := decoder.Base64Decode(input)
		if err == nil && isPrintable([]byte(decoded)) {
			return decoded, nil
		}
	}

	// 4. YAML → JSON
	if isYAML(input) {
		return yamlToJSON(input)
	}

	// 5. URL encoded → decode
	if isURLEncoded(input) {
		decoded, err := decoder.URLDecode(input)
		if err == nil {
			return decoded, nil
		}
	}

	// 6. Hex → decode
	if isHex(input) {
		decoded, err := decoder.HexDecode(input)
		if err == nil && isPrintable([]byte(decoded)) {
			return decoded, nil
		}
	}

	return "", fmt.Errorf("could not auto-detect input format")
}

// isTimestamp reports whether input looks like a UNIX timestamp.
// Only matches 10-digit (seconds, ~2001-2286) or 13-digit (milliseconds, ~2001-2286)
// numeric strings to avoid false positives with other formats.
func isTimestamp(s string) bool {
	if !timconv.IsTimestamp(s) {
		return false
	}
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil || v < 0 {
		return false
	}
	digits := len(strings.TrimSpace(s))
	return digits == 10 || digits == 13
}

func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func formatJSON(s string) (string, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return "", err
	}
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isBase64(s string) bool {
	if len(s) < 4 {
		return false
	}
	match, _ := regexp.MatchString(`^[A-Za-z0-9+/]+=*$`, s)
	if !match {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func isPrintable(b []byte) bool {
	for _, c := range b {
		if c < 0x20 && c != '\n' && c != '\r' && c != '\t' {
			return false
		}
	}
	return true
}

func isYAML(s string) bool {
	// Must contain ':' to look like YAML key-value
	if !strings.Contains(s, ":") {
		return false
	}
	var obj interface{}
	err := yaml.Unmarshal([]byte(s), &obj)
	if err != nil || obj == nil {
		return false
	}
	// Only consider structured data (map or slice) as YAML, not scalar values
	switch obj.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}

func yamlToJSON(s string) (string, error) {
	var obj interface{}
	if err := yaml.Unmarshal([]byte(s), &obj); err != nil {
		return "", err
	}
	obj = convertYAMLToJSON(obj)
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// convertYAMLToJSON converts YAML map[interface{}]interface{} to map[string]interface{} for JSON compat.
func convertYAMLToJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = convertYAMLToJSON(val)
		}
		return result
	case []interface{}:
		for i, val := range v {
			v[i] = convertYAMLToJSON(val)
		}
		return v
	default:
		return v
	}
}

func isURLEncoded(s string) bool {
	if !strings.Contains(s, "%") {
		return false
	}
	decoded, err := url.QueryUnescape(s)
	return err == nil && decoded != s
}

func isHex(s string) bool {
	if len(s) < 2 || len(s)%2 != 0 {
		return false
	}
	match, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, s)
	return match
}

