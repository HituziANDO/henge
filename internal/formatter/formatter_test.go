package formatter

import (
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		indent  int
		want    string
		wantErr bool
	}{
		{
			name:   "basic object with default indent",
			input:  `{"name":"henge","version":"0.1.0"}`,
			indent: 0,
			want:   "{\n  \"name\": \"henge\",\n  \"version\": \"0.1.0\"\n}",
		},
		{
			name:   "basic object with 4-space indent",
			input:  `{"name":"henge"}`,
			indent: 4,
			want:   "{\n    \"name\": \"henge\"\n}",
		},
		{
			name:   "array",
			input:  `[1,2,3]`,
			indent: 2,
			want:   "[\n  1,\n  2,\n  3\n]",
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			indent:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSON(tt.input, tt.indent)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("FormatJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCompactJSON(t *testing.T) {
	input := `{
  "name": "henge",
  "version": "0.1.0"
}`
	want := `{"name":"henge","version":"0.1.0"}`
	got, err := CompactJSON(input)
	if err != nil {
		t.Fatalf("CompactJSON() error = %v", err)
	}
	if got != want {
		t.Errorf("CompactJSON() = %q, want %q", got, want)
	}

	_, err = CompactJSON(`{invalid}`)
	if err == nil {
		t.Error("CompactJSON() expected error for invalid JSON")
	}
}

func TestFormatYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(string) bool
	}{
		{
			name:  "basic YAML",
			input: "name: henge\nversion: 0.1.0",
			check: func(s string) bool {
				return strings.Contains(s, "name: henge") && strings.Contains(s, "version: 0.1.0")
			},
		},
		{
			name:  "nested YAML",
			input: "server:\n  host: localhost\n  port: 8080",
			check: func(s string) bool {
				return strings.Contains(s, "server:") && strings.Contains(s, "host: localhost")
			},
		},
		{
			name:    "invalid YAML",
			input:   ":\n  :\n  - :\n  - :",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil && !tt.check(got) {
				t.Errorf("FormatYAML() = %q, check failed", got)
			}
		})
	}
}

func TestFormatXML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(string) bool
	}{
		{
			name:  "basic XML",
			input: `<root><name>henge</name><version>0.1.0</version></root>`,
			check: func(s string) bool {
				return strings.Contains(s, "<root>") && strings.Contains(s, "  <name>henge</name>")
			},
		},
		{
			name:  "XML with attributes",
			input: `<root attr="val"><child>text</child></root>`,
			check: func(s string) bool {
				return strings.Contains(s, "<root") && strings.Contains(s, "<child>text</child>")
			},
		},
		{
			name:    "invalid XML",
			input:   `<root><unclosed>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatXML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatXML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil && !tt.check(got) {
				t.Errorf("FormatXML() = %q, check failed", got)
			}
		})
	}
}
