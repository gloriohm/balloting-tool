package ballot

import "ballot-tool/internal/utils"

var rolesHeaderAliases = map[string]string{
	"account_email":       "email",
	"committee_reference": "committee",
	"committee_name":      "committee",
	"name":                "committee",
	"role_description":    "role",
	"commitment_role":     "role",
}

var ballotHeaderAliases = map[string]string{
	"committee_working_group": "committee",
	"opening_date":            "opened",
	"start_date":              "opened",
	"closing_date":            "closes",
	"end_date":                "closes",
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

func normalizeHeaders(rows []map[string]string, aliases map[string]string) []map[string]string {
	normalized := make([]map[string]string, 0, len(rows))

	for _, row := range rows {
		m := make(map[string]string, len(row))

		for header, value := range row {
			if canonical, ok := aliases[header]; ok {
				m[canonical] = value
			} else {
				m[header] = value
			}
		}

		normalized = append(normalized, m)
	}

	return normalized
}
