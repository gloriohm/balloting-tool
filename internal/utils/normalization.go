package utils

import (
	"regexp"
	"strings"
	"time"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func NormalizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = nonAlnum.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	return s
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func ToLowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func SanitizeFilename(s string) string {
	s = strings.TrimSpace(s)

	replacer := strings.NewReplacer(
		".", "",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	s = replacer.Replace(s)

	if s == "" {
		return "file"
	}
	return s
}

func StripLabel(s string, label string) string {
	return strings.TrimSpace(strings.TrimPrefix(s, label))
}

func NormalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
