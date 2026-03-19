package encoder

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// Base64Encode encodes the input string to base64.
func Base64Encode(input string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(input)), nil
}

// URLEncode encodes the input string using percent-encoding.
func URLEncode(input string) (string, error) {
	return url.QueryEscape(input), nil
}

// HexEncode encodes the input string to hexadecimal.
func HexEncode(input string) (string, error) {
	return hex.EncodeToString([]byte(input)), nil
}
