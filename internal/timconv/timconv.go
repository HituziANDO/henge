package timconv

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// dateFormats defines the date formats to try in order when parsing date strings.
var dateFormats = []struct {
	layout string
	hasTZ  bool
}{
	{time.RFC3339, true},
	{time.RFC3339Nano, true},
	{time.RFC1123, true},
	{time.RFC822, true},
	{"2006-01-02 15:04:05", false},
	{"2006-01-02", false},
	{"2006/01/02", false},
	{"2006/01/02 15:04:05", false},
}

// formatPresets maps preset names (lowercase) to Go time layout strings.
var formatPresets = map[string]string{
	"rfc3339":  time.RFC3339,
	"rfc1123":  time.RFC1123,
	"rfc822":   time.RFC822,
	"datetime": "2006-01-02 15:04:05",
	"date":     "2006-01-02",
}

// millisThreshold is the threshold above which a timestamp is considered
// to be in milliseconds rather than seconds.
const millisThreshold int64 = 10000000000

// DateToUnix converts a date string to a UNIX timestamp (seconds).
// It tries multiple date formats in order and returns the first match.
// If tzName is non-empty, it's applied to inputs that lack timezone info.
func DateToUnix(input string, tzName string) (int64, error) {
	var loc *time.Location
	if tzName != "" {
		var err error
		loc, err = time.LoadLocation(tzName)
		if err != nil {
			return 0, fmt.Errorf("unknown timezone: %q", tzName)
		}
	}

	for _, f := range dateFormats {
		if !f.hasTZ && loc != nil {
			t, err := time.ParseInLocation(f.layout, input, loc)
			if err == nil {
				return t.Unix(), nil
			}
		} else if !f.hasTZ {
			t, err := time.ParseInLocation(f.layout, input, time.UTC)
			if err == nil {
				return t.Unix(), nil
			}
		} else {
			t, err := time.Parse(f.layout, input)
			if err == nil {
				return t.Unix(), nil
			}
		}
	}

	return 0, fmt.Errorf("unable to parse date string: %q", input)
}

// UnixToDate converts a UNIX timestamp (seconds) to a formatted date string.
// format can be a preset name ("rfc3339", "rfc1123", "rfc822", "datetime", "date")
// or a Go time layout string. If tzName is empty, UTC is used.
func UnixToDate(timestamp int64, format string, tzName string) (string, error) {
	loc := time.UTC
	if tzName != "" {
		var err error
		loc, err = time.LoadLocation(tzName)
		if err != nil {
			return "", fmt.Errorf("unknown timezone: %q", tzName)
		}
	}

	layout := ResolveFormat(format)
	t := time.Unix(timestamp, 0).In(loc)
	return t.Format(layout), nil
}

// ParseUnixTimestamp parses a string as a UNIX timestamp in seconds.
// Auto-detects milliseconds (values >= 10000000000) and normalizes to seconds.
// If forceMillis is true, always interprets as milliseconds.
func ParseUnixTimestamp(input string, forceMillis bool) (int64, error) {
	input = strings.TrimSpace(input)
	v, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid UNIX timestamp: %q", input)
	}

	if forceMillis {
		return v / 1000, nil
	}

	if v >= millisThreshold || v <= -millisThreshold {
		return v / 1000, nil
	}

	return v, nil
}

// IsTimestamp reports whether input looks like a UNIX timestamp (purely numeric,
// optionally with a leading minus sign).
func IsTimestamp(input string) bool {
	s := strings.TrimSpace(input)
	if s == "" {
		return false
	}
	if s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// AutoConvert auto-detects the input type and converts it.
// If the input is a numeric UNIX timestamp, it converts to an RFC3339 date string.
// If the input is a date string, it converts to a UNIX timestamp.
func AutoConvert(input string) (string, error) {
	if IsTimestamp(input) {
		ts, err := ParseUnixTimestamp(input, false)
		if err != nil {
			return "", err
		}
		return UnixToDate(ts, "rfc3339", "")
	}
	ts, err := DateToUnix(input, "")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ts, 10), nil
}

// ResolveFormat resolves a preset name to a Go time layout string.
// If the input doesn't match any preset, it's returned as-is.
func ResolveFormat(format string) string {
	if layout, ok := formatPresets[strings.ToLower(format)]; ok {
		return layout
	}
	return format
}
