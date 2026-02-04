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

var HeaderAliases = map[string]string{
	// ballot variants
	"committee_working_group": "committee",
	"opening_date":            "opened",
	"start_date":              "opened",
	"closing_date":            "closes",
	"end_date":                "closes",

	// role variants
	"account_email":       "email",
	"committee_reference": "committee",
	"committee_name":      "committee",
	"name":                "committee",
	"role_description":    "role",
	"commitment_role":     "role",
}

func IsVoterRole(raw string) bool {
	switch NormalizeString(raw) {
	case "voter", "cen_voter", "obligated_voter":
		return true
	default:
		return false
	}
}

func IsMemberStatus(raw string) bool {
	switch NormalizeString(raw) {
	case "p_member", "o_member":
		return true
	default:
		return false
	}
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func ToLowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
