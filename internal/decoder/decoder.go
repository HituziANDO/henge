package decoder

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
)

// Base64Decode decodes a base64 encoded string.
func Base64Decode(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("invalid base64 input: %w", err)
	}
	return string(decoded), nil
}

// URLDecode decodes a percent-encoded string.
func URLDecode(input string) (string, error) {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL-encoded input: %w", err)
	}
	return decoded, nil
}

// HexDecode decodes a hexadecimal encoded string.
func HexDecode(input string) (string, error) {
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("invalid hex input: %w", err)
	}
	return string(decoded), nil
}
