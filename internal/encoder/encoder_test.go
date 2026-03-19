package encoder

import (
	"testing"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", "aGVsbG8="},
		{"empty string", "", ""},
		{"with spaces", "hello world", "aGVsbG8gd29ybGQ="},
		{"unicode", "こんにちは", "44GT44KT44Gr44Gh44Gv"},
		{"special chars", "foo@bar.com", "Zm9vQGJhci5jb20="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Encode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Base64Encode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestURLEncode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", "hello"},
		{"empty string", "", ""},
		{"with spaces", "hello world", "hello+world"},
		{"special chars", "foo@bar.com", "foo%40bar.com"},
		{"query string", "key=value&key2=value2", "key%3Dvalue%26key2%3Dvalue2"},
		{"unicode", "こんにちは", "%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%81%AF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLEncode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("URLEncode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHexEncode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", "68656c6c6f"},
		{"empty string", "", ""},
		{"with spaces", "hello world", "68656c6c6f20776f726c64"},
		{"special chars", "!@#", "214023"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexEncode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("HexEncode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
