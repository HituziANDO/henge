package detector

import (
	"testing"
)

func TestAutoDetectAndTransform(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Empty input
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   \n\t  ",
			wantErr: true,
		},

		// 1. JSON → pretty print
		{
			name:  "JSON object compact to pretty",
			input: `{"name":"henge","version":"0.1.0"}`,
			want:  "{\n  \"name\": \"henge\",\n  \"version\": \"0.1.0\"\n}",
		},
		{
			name:  "JSON array",
			input: `[1,2,3]`,
			want:  "[\n  1,\n  2,\n  3\n]",
		},

		// 2. Base64 → decode
		{
			name:  "base64 decode",
			input: "aGVsbG8gd29ybGQ=",
			want:  "hello world",
		},
		{
			name:  "base64 decode no padding",
			input: "aGVsbG8=",
			want:  "hello",
		},

		// 3. YAML → JSON
		{
			name:  "YAML to JSON",
			input: "name: henge\nversion: 0.1.0",
			want:  "{\n  \"name\": \"henge\",\n  \"version\": \"0.1.0\"\n}",
		},
		{
			name:  "YAML nested to JSON",
			input: "server:\n  host: localhost\n  port: 8080",
			want:  "{\n  \"server\": {\n    \"host\": \"localhost\",\n    \"port\": 8080\n  }\n}",
		},

		// 4. URL encoded → decode
		{
			name:  "URL decode percent encoding",
			input: "hello%20world%21",
			want:  "hello world!",
		},
		{
			name:  "URL decode complex",
			input: "key%3Dvalue%26foo%3Dbar",
			want:  "key=value&foo=bar",
		},

		// 5. Hex → decode
		{
			name:  "hex decode lowercase",
			input: "68656c6c6f",
			want:  "hello",
		},
		{
			name:  "hex decode uppercase",
			input: "48454C4C4F",
			want:  "HELLO",
		},

		// Unrecognized input
		{
			name:    "unrecognized input",
			input:   "just plain text",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AutoDetectAndTransform(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AutoDetectAndTransform(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("AutoDetectAndTransform(%q) =\n%s\nwant:\n%s", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{`{"key": "value"}`, true},
		{`[1, 2, 3]`, true},
		{`"string"`, true},
		{`not json`, false},
		{`{invalid}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isJSON(tt.input); got != tt.want {
				t.Errorf("isJSON(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsBase64(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"aGVsbG8=", true},
		{"aGVsbG8gd29ybGQ=", true},
		{"abc", false},  // too short
		{"!!!", false},  // invalid chars
		{"", false},     // empty
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isBase64(tt.input); got != tt.want {
				t.Errorf("isBase64(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsURLEncoded(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"hello%20world", true},
		{"foo%3Dbar", true},
		{"nopercent", false},
		{"hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isURLEncoded(tt.input); got != tt.want {
				t.Errorf("isURLEncoded(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"68656c6c6f", true},
		{"48454C4C4F", true},
		{"abcdef", true},
		{"xyz", false},
		{"", false},
		{"a", false},   // odd length
		{"GG", false},  // invalid hex chars
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isHex(tt.input); got != tt.want {
				t.Errorf("isHex(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{"printable ASCII", []byte("hello world"), true},
		{"with newline", []byte("hello\nworld"), true},
		{"with tab", []byte("hello\tworld"), true},
		{"with null byte", []byte("hello\x00world"), false},
		{"with control char", []byte("hello\x01world"), false},
		{"empty", []byte{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrintable(tt.input); got != tt.want {
				t.Errorf("isPrintable(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
