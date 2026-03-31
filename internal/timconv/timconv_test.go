package timconv

import (
	"strings"
	"testing"
	"time"
)

func TestDateToUnix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		tz      string
		want    int64
		wantErr string
	}{
		{
			name:  "RFC3339 input",
			input: "2025-01-01T00:00:00Z",
			tz:    "",
			want:  1735689600,
		},
		{
			name:  "RFC3339 with timezone offset",
			input: "2025-01-01T09:00:00+09:00",
			tz:    "",
			want:  1735689600,
		},
		{
			name:  "DateOnly",
			input: "2025-01-01",
			tz:    "",
			want:  1735689600,
		},
		{
			name:  "DateSlash",
			input: "2025/01/01",
			tz:    "",
			want:  1735689600,
		},
		{
			name:  "DateTime with timezone flag Asia/Tokyo",
			input: "2025-01-01 09:00:00",
			tz:    "Asia/Tokyo",
			want:  1735689600,
		},
		{
			name:  "RFC1123",
			input: "Wed, 01 Jan 2025 00:00:00 UTC",
			tz:    "",
			want:  1735689600,
		},
		{
			name:  "RFC3339Nano",
			input: "2025-01-01T00:00:00.123456789Z",
			tz:    "",
			want:  1735689600,
		},
		{
			name:    "Invalid date string",
			input:   "not-a-date",
			tz:      "",
			wantErr: "unable to parse date string",
		},
		{
			name:    "Invalid timezone",
			input:   "2025-01-01",
			tz:      "Invalid/TZ",
			wantErr: "unknown timezone",
		},
		{
			name:  "DateOnly without timezone should be UTC",
			input: "2025-01-01",
			tz:    "",
			want:  1735689600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DateToUnix(tt.input, tt.tz)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("DateToUnix(%q, %q) = %d, want %d", tt.input, tt.tz, got, tt.want)
			}
		})
	}
}

func TestUnixToDate(t *testing.T) {
	tests := []struct {
		name    string
		ts      int64
		format  string
		tz      string
		want    string
		wantErr string
	}{
		{
			name:   "Timestamp to RFC3339 default",
			ts:     1735689600,
			format: "rfc3339",
			tz:     "",
			want:   "2025-01-01T00:00:00Z",
		},
		{
			name:   "Timestamp to Asia/Tokyo",
			ts:     1735689600,
			format: "rfc3339",
			tz:     "Asia/Tokyo",
			want:   "2025-01-01T09:00:00+09:00",
		},
		{
			name:   "Timestamp to custom format",
			ts:     1735689600,
			format: "2006/01/02",
			tz:     "",
			want:   "2025/01/01",
		},
		{
			name:   "Timestamp to datetime preset",
			ts:     1735689600,
			format: "datetime",
			tz:     "",
			want:   "2025-01-01 00:00:00",
		},
		{
			name:   "Timestamp to date preset",
			ts:     1735689600,
			format: "date",
			tz:     "",
			want:   "2025-01-01",
		},
		{
			name:    "Invalid timezone",
			ts:      1735689600,
			format:  "rfc3339",
			tz:      "Invalid/TZ",
			wantErr: "unknown timezone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnixToDate(tt.ts, tt.format, tt.tz)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("UnixToDate(%d, %q, %q) = %q, want %q", tt.ts, tt.format, tt.tz, got, tt.want)
			}
		})
	}
}

func TestParseUnixTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		forceMillis bool
		want        int64
		wantErr     string
	}{
		{
			name:        "Seconds 10-digit",
			input:       "1735689600",
			forceMillis: false,
			want:        1735689600,
		},
		{
			name:        "Milliseconds auto-detect 13-digit",
			input:       "1735689600000",
			forceMillis: false,
			want:        1735689600,
		},
		{
			name:        "Force millis",
			input:       "1735689600000",
			forceMillis: true,
			want:        1735689600,
		},
		{
			name:        "Non-numeric input",
			input:       "abc",
			forceMillis: false,
			wantErr:     "invalid UNIX timestamp",
		},
		{
			name:        "Input with whitespace",
			input:       "  1735689600  ",
			forceMillis: false,
			want:        1735689600,
		},
		{
			name:        "Negative timestamp",
			input:       "-1",
			forceMillis: false,
			want:        -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUnixTimestamp(tt.input, tt.forceMillis)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParseUnixTimestamp(%q, %v) = %d, want %d", tt.input, tt.forceMillis, got, tt.want)
			}
		})
	}
}

func TestResolveFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
	}{
		{
			name:  "rfc3339 lowercase",
			input: "rfc3339",
			want:  time.RFC3339,
		},
		{
			name:  "RFC3339 case insensitive",
			input: "RFC3339",
			want:  time.RFC3339,
		},
		{
			name:  "datetime preset",
			input: "datetime",
			want:  "2006-01-02 15:04:05",
		},
		{
			name:  "date preset",
			input: "date",
			want:  "2006-01-02",
		},
		{
			name:  "rfc1123 preset",
			input: "rfc1123",
			want:  time.RFC1123,
		},
		{
			name:  "rfc822 preset",
			input: "rfc822",
			want:  time.RFC822,
		},
		{
			name:  "Custom layout passthrough",
			input: "2006/01/02 15:04",
			want:  "2006/01/02 15:04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveFormat(tt.input)
			if got != tt.want {
				t.Errorf("ResolveFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
