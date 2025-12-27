package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

// FormatJSON formats a JSON string with indentation
func FormatJSON(input string) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(input), "", "  "); err != nil {
		return input, err
	}
	return buf.String(), nil
}

// IsValidJSON checks if a string is valid JSON
func IsValidJSON(input string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(input), &js) == nil
}

// CompactJSON removes whitespace from JSON
func CompactJSON(input string) (string, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return input, err
	}
	return buf.String(), nil
}

// ParseHeaders parses a string of headers (key: value format) into a map
func ParseHeaders(input string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				headers[key] = value
			}
		}
	}
	return headers
}

// FormatHeaders formats a map of headers as a string
func FormatHeaders(headers map[string]string) string {
	var lines []string
	for key, value := range headers {
		lines = append(lines, key+": "+value)
	}
	return strings.Join(lines, "\n")
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
