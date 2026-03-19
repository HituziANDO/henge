package hasher

import "testing"

func TestMD5Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
	}
	for _, tt := range tests {
		got := MD5Hash(tt.input)
		if got != tt.expected {
			t.Errorf("MD5Hash(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSHA1Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
	}
	for _, tt := range tests {
		got := SHA1Hash(tt.input)
		if got != tt.expected {
			t.Errorf("SHA1Hash(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSHA256Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}
	for _, tt := range tests {
		got := SHA256Hash(tt.input)
		if got != tt.expected {
			t.Errorf("SHA256Hash(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSHA512Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	}
	for _, tt := range tests {
		got := SHA512Hash(tt.input)
		if got != tt.expected {
			t.Errorf("SHA512Hash(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
