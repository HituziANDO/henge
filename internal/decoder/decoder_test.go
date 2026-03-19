package decoder

import (
	"testing"
)

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple string", "aGVsbG8=", "hello", false},
		{"empty string", "", "", false},
		{"with spaces", "aGVsbG8gd29ybGQ=", "hello world", false},
		{"unicode", "44GT44KT44Gr44Gh44Gv", "こんにちは", false},
		{"invalid base64", "not-valid-base64!", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Base64Decode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Base64Decode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestURLDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple string", "hello", "hello", false},
		{"empty string", "", "", false},
		{"with plus", "hello+world", "hello world", false},
		{"percent encoded", "foo%40bar.com", "foo@bar.com", false},
		{"unicode", "%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%81%AF", "こんにちは", false},
		{"invalid encoding", "%GG", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("URLDecode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("URLDecode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHexDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple string", "68656c6c6f", "hello", false},
		{"empty string", "", "", false},
		{"with spaces", "68656c6c6f20776f726c64", "hello world", false},
		{"uppercase hex", "48454C4C4F", "HELLO", false},
		{"invalid hex", "zzzz", "", true},
		{"odd length", "abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("HexDecode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("HexDecode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
