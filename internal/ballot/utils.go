package ballot

import "ballot-tool/internal/utils"

var headerAliases = map[string]string{
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

func isVoterRole(raw string) bool {
	switch utils.NormalizeString(raw) {
	case "voter", "cen_voter", "obligated_voter":
		return true
	default:
		return false
	}
}

func isMemberStatus(raw string) bool {
	switch utils.NormalizeString(raw) {
	case "p_member", "o_member":
		return true
	default:
		return false
	}
}
