package converter

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"json object", `{"key": "value"}`, "json"},
		{"json array", `[1, 2, 3]`, "json"},
		{"yaml", "name: test\nvalue: 123", "yaml"},
		{"toml", "key = \"value\"\nnum = 42", "toml"},
		{"csv", "name,age\nAlice,30\nBob,25", "csv"},
		{"empty", "", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormat(tt.input)
			if result != tt.expected {
				t.Errorf("DetectFormat() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestJSONToYAML(t *testing.T) {
	input := `{"name": "test", "value": 123}`
	result, err := ToYAML(input, "json")
	if err != nil {
		t.Fatalf("ToYAML() error: %v", err)
	}

	// Verify it's valid YAML
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid YAML: %v", err)
	}
	if data["name"] != "test" {
		t.Errorf("name = %v, want %q", data["name"], "test")
	}
	if data["value"] != 123 {
		t.Errorf("value = %v, want 123", data["value"])
	}
}

func TestYAMLToJSON(t *testing.T) {
	input := "name: test\nvalue: 123"
	result, err := ToJSON(input, "yaml")
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	// Verify it's valid JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if data["name"] != "test" {
		t.Errorf("name = %v, want %q", data["name"], "test")
	}
	if data["value"] != float64(123) {
		t.Errorf("value = %v, want 123", data["value"])
	}
}

func TestJSONToTOML(t *testing.T) {
	input := `{"name": "test", "value": 123}`
	result, err := ToTOML(input, "json")
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}

	if result == "" {
		t.Fatal("ToTOML() returned empty string")
	}

	// Result should contain key-value pairs
	if !containsSubstring(result, "name") || !containsSubstring(result, "test") {
		t.Errorf("TOML output missing expected content: %s", result)
	}
}

func TestCSVToJSON(t *testing.T) {
	input := "name,age,city\nAlice,30,Tokyo\nBob,25,Osaka"
	result, err := ToJSON(input, "csv")
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	// Verify it's valid JSON array
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid JSON array: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 records, got %d", len(data))
	}
	if data[0]["name"] != "Alice" {
		t.Errorf("first record name = %v, want %q", data[0]["name"], "Alice")
	}
	if data[0]["age"] != "30" {
		t.Errorf("first record age = %v, want %q", data[0]["age"], "30")
	}
	if data[1]["city"] != "Osaka" {
		t.Errorf("second record city = %v, want %q", data[1]["city"], "Osaka")
	}
}

func TestAutoDetectConversion(t *testing.T) {
	// Auto-detect JSON and convert to YAML
	input := `{"hello": "world"}`
	result, err := ToYAML(input, "")
	if err != nil {
		t.Fatalf("ToYAML() with auto-detect error: %v", err)
	}
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("result is not valid YAML: %v", err)
	}
	if data["hello"] != "world" {
		t.Errorf("hello = %v, want %q", data["hello"], "world")
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
