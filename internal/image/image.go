package image

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EncodeFileToBase64 reads an image file and returns its Base64 encoded string.
func EncodeFileToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// EncodeFileToDataURI reads an image file and returns a Data URI string.
func EncodeFileToDataURI(filePath string) (string, error) {
	b64, err := EncodeFileToBase64(filePath)
	if err != nil {
		return "", err
	}
	mime := DetectMIMEType(filePath)
	return fmt.Sprintf("data:%s;base64,%s", mime, b64), nil
}

// DecodeBase64ToFile decodes a Base64 string and writes the result to a file.
// It automatically handles Data URI scheme inputs and strips whitespace/newlines.
func DecodeBase64ToFile(input string, outputPath string) error {
	_, data := ParseDataURI(input)

	// Remove whitespace and newlines
	data = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, data)

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid base64 input: %w", err)
	}

	if err := os.WriteFile(outputPath, decoded, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WrapString wraps a string at the specified width.
// If width <= 0, the string is returned unchanged.
func WrapString(s string, width int) string {
	if width <= 0 {
		return s
	}

	var builder strings.Builder
	for i := 0; i < len(s); i += width {
		end := i + width
		if end > len(s) {
			end = len(s)
		}
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(s[i:end])
	}
	return builder.String()
}

// DetectMIMEType determines the MIME type from the file extension.
func DetectMIMEType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

// ParseDataURI parses a Data URI string and returns the MIME type and Base64 data.
// If the input is not a Data URI, it returns an empty MIME type and the original string.
func ParseDataURI(input string) (mimeType string, data string) {
	const prefix = "data:"
	if !strings.HasPrefix(input, prefix) {
		return "", input
	}

	// Find ";base64," separator
	idx := strings.Index(input, ";base64,")
	if idx < 0 {
		return "", input
	}

	mimeType = input[len(prefix):idx]
	data = input[idx+len(";base64,"):]
	return mimeType, data
}
