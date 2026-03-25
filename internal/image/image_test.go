package image

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

// Minimal 1x1 pixel PNG
var testPNGData = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00,
	0x0c, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x00, 0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc, 0x33, 0x00, 0x00, 0x00,
	0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

// writeTempPNG creates a temporary PNG file and returns its path.
func writeTempPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.png")
	if err := os.WriteFile(p, testPNGData, 0644); err != nil {
		t.Fatalf("failed to write temp PNG: %v", err)
	}
	return p
}

func TestEncodeFileToBase64(t *testing.T) {
	p := writeTempPNG(t)

	got, err := EncodeFileToBase64(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := base64.StdEncoding.EncodeToString(testPNGData)
	if got != want {
		t.Errorf("EncodeFileToBase64: got %q, want %q", got, want)
	}
}

func TestEncodeFileToBase64_FileNotFound(t *testing.T) {
	_, err := EncodeFileToBase64("/nonexistent/file.png")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestEncodeFileToDataURI(t *testing.T) {
	p := writeTempPNG(t)

	got, err := EncodeFileToDataURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b64 := base64.StdEncoding.EncodeToString(testPNGData)
	want := "data:image/png;base64," + b64
	if got != want {
		t.Errorf("EncodeFileToDataURI: got %q, want %q", got, want)
	}
}

func TestDecodeBase64ToFile(t *testing.T) {
	// Encode the test PNG data, then decode it back and verify round-trip.
	b64 := base64.StdEncoding.EncodeToString(testPNGData)
	dir := t.TempDir()
	outPath := filepath.Join(dir, "output.png")

	if err := DecodeBase64ToFile(b64, outPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(got) != len(testPNGData) {
		t.Fatalf("decoded file size mismatch: got %d, want %d", len(got), len(testPNGData))
	}
	for i := range got {
		if got[i] != testPNGData[i] {
			t.Fatalf("decoded file differs at byte %d: got 0x%02x, want 0x%02x", i, got[i], testPNGData[i])
		}
	}
}

func TestDecodeBase64ToFile_DataURI(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString(testPNGData)
	dataURI := "data:image/png;base64," + b64
	dir := t.TempDir()
	outPath := filepath.Join(dir, "output.png")

	if err := DecodeBase64ToFile(dataURI, outPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(got) != len(testPNGData) {
		t.Fatalf("decoded file size mismatch: got %d, want %d", len(got), len(testPNGData))
	}
	for i := range got {
		if got[i] != testPNGData[i] {
			t.Fatalf("decoded file differs at byte %d: got 0x%02x, want 0x%02x", i, got[i], testPNGData[i])
		}
	}
}

func TestDecodeBase64ToFile_InvalidBase64(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "output.png")

	err := DecodeBase64ToFile("!!!not-valid-base64!!!", outPath)
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestDecodeBase64ToFile_InvalidOutputPath(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString(testPNGData)

	err := DecodeBase64ToFile(b64, "/nonexistent-dir/subdir/output.png")
	if err == nil {
		t.Fatal("expected error for invalid output path, got nil")
	}
}

func TestWrapString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{
			name:  "wrap at 4",
			input: "abcdefgh",
			width: 4,
			want:  "abcd\nefgh",
		},
		{
			name:  "no wrap when width is 0",
			input: "abcdefgh",
			width: 0,
			want:  "abcdefgh",
		},
		{
			name:  "no wrap when width is negative",
			input: "abcdefgh",
			width: -1,
			want:  "abcdefgh",
		},
		{
			name:  "input shorter than width",
			input: "abc",
			width: 10,
			want:  "abc",
		},
		{
			name:  "exact width",
			input: "abcd",
			width: 4,
			want:  "abcd",
		},
		{
			name:  "empty string",
			input: "",
			width: 4,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapString(tt.input, tt.width)
			if got != tt.want {
				t.Errorf("WrapString(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
			}
		})
	}
}

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		filePath string
		want     string
	}{
		{"image.png", "image/png"},
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.gif", "image/gif"},
		{"image.webp", "image/webp"},
		{"image.bmp", "image/bmp"},
		{"image.svg", "image/svg+xml"},
		{"image.ico", "image/x-icon"},
		{"image.unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			got := DetectMIMEType(tt.filePath)
			if got != tt.want {
				t.Errorf("DetectMIMEType(%q) = %q, want %q", tt.filePath, got, tt.want)
			}
		})
	}
}

func TestParseDataURI(t *testing.T) {
	input := "data:image/png;base64,iVBORw0KGgo="
	mime, data := ParseDataURI(input)
	if mime != "image/png" {
		t.Errorf("ParseDataURI mime: got %q, want %q", mime, "image/png")
	}
	if data != "iVBORw0KGgo=" {
		t.Errorf("ParseDataURI data: got %q, want %q", data, "iVBORw0KGgo=")
	}
}

func TestParseDataURI_NotDataURI(t *testing.T) {
	input := "iVBORw0KGgo="
	mime, data := ParseDataURI(input)
	if mime != "" {
		t.Errorf("ParseDataURI_NotDataURI mime: got %q, want empty", mime)
	}
	if data != input {
		t.Errorf("ParseDataURI_NotDataURI data: got %q, want %q", data, input)
	}
}
