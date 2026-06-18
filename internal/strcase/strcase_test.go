package strcase

import "testing"

func TestToUpper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "hello", "HELLO"},
		{"preserve delimiters", "foo_bar Baz", "FOO_BAR BAZ"},
		{"empty", "", ""},
		{"unicode passthrough", "café こんにちは", "CAFÉ こんにちは"},
		{"symbols preserved", "a-b.c/d", "A-B.C/D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUpper(tt.input); got != tt.want {
				t.Errorf("ToUpper(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "HELLO", "hello"},
		{"preserve delimiters", "Foo-BAR Baz", "foo-bar baz"},
		{"empty", "", ""},
		{"unicode passthrough", "CAFÉ", "café"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLower(tt.input); got != tt.want {
				t.Errorf("ToLower(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSnake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"from camel", "fooBar", "foo_bar"},
		{"from pascal", "FooBar", "foo_bar"},
		{"from kebab", "foo-bar", "foo_bar"},
		{"from spaces", "User First Name", "user_first_name"},
		{"acronym run", "HTTPServer", "http_server"},
		{"trailing acronym", "parseURL", "parse_url"},
		{"id acronym", "userID", "user_id"},
		{"digits attach", "user42name", "user42name"},
		{"digit then upper", "foo2Bar", "foo2_bar"},
		{"idempotent", "already_snake", "already_snake"},
		{"collapse delimiters", "__foo--bar  baz__", "foo_bar_baz"},
		{"empty", "", ""},
		{"multiline", "fooBar\nbazQux", "foo_bar\nbaz_qux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnake(tt.input); got != tt.want {
				t.Errorf("ToSnake(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToKebab(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"from camel", "fooBar", "foo-bar"},
		{"from snake", "foo_bar", "foo-bar"},
		{"acronym run", "HTTPServer", "http-server"},
		{"idempotent", "already-kebab", "already-kebab"},
		{"empty", "", ""},
		{"multiline", "fooBar\nbazQux", "foo-bar\nbaz-qux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToKebab(tt.input); got != tt.want {
				t.Errorf("ToKebab(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToCamel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"from snake", "foo_bar", "fooBar"},
		{"from kebab", "foo-bar", "fooBar"},
		{"from pascal", "FooBar", "fooBar"},
		{"acronym not preserved", "userID", "userId"},
		{"three words", "user_first_name", "userFirstName"},
		{"single word", "hello", "hello"},
		{"empty", "", ""},
		{"multiline", "foo_bar\nbaz_qux", "fooBar\nbazQux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamel(tt.input); got != tt.want {
				t.Errorf("ToCamel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"from snake", "foo_bar", "FooBar"},
		{"from kebab", "foo-bar", "FooBar"},
		{"from camel", "fooBar", "FooBar"},
		{"acronym not preserved", "parseURL", "ParseUrl"},
		{"idempotent", "FooBar", "FooBar"},
		{"empty", "", ""},
		{"multiline", "foo_bar\nbaz_qux", "FooBar\nBazQux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPascal(tt.input); got != tt.want {
				t.Errorf("ToPascal(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRoundTrips(t *testing.T) {
	// snake -> camel -> snake is stable for already-tokenized identifiers.
	in := "user_first_name"
	if got := ToSnake(ToCamel(in)); got != in {
		t.Errorf("round-trip snake->camel->snake: got %q, want %q", got, in)
	}
	// snake -> kebab -> snake is stable.
	if got := ToSnake(ToKebab(in)); got != in {
		t.Errorf("round-trip snake->kebab->snake: got %q, want %q", got, in)
	}
}
