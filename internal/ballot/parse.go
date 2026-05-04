package ballot

import (
	"ballot-tool/internal/utils"
	"fmt"
	"log"
	"strings"
)

func parseRoles(rows []map[string]string) ([]Role, []error) {
	out := make([]Role, 0, len(rows))
	var errs []error
	for i, row := range rows {
		rawRole := row["role"]
		if !isVoterRole(rawRole) {
			continue
		}
		commitmentStatus := row["commitment_status"]
		if commitmentStatus != "" && commitmentStatus != "active" {
			continue
		}
		com := strings.TrimSpace(row["committee"])
		if com == "" {
			errs = append(errs, fmt.Errorf("row %d: missing committee reference", i+1))
			continue
		}

		out = append(out, Role{
			Committee: com,
			FirstName: row["first_name"],
			LastName:  row["last_name"],
			Email:     row["email"],
		})
	}

	log.Printf("parsed roles, length %d\n", len(out))
	return out, errs
}

func parseBallots(rows []map[string]string) ([]Ballot, []error) {
	out := make([]Ballot, 0, len(rows))
	var errs []error
	for i, row := range rows {
		com := strings.TrimSpace(row["committee"])
		if com == "" {
			errs = append(errs, fmt.Errorf("row %d: missing committee reference", i+1))
			continue
		}
		close := strings.TrimSpace(row["closes"])
		if close == "" {
			continue
		}
		closeTime, err := utils.ParseDate(close)
		if err != nil {
			errs = append(errs, fmt.Errorf("row %d: closing date not formatted as ISO date", i+1))
			continue
		}

		out = append(out, Ballot{
			Source:    row["source"],
			Committee: com,
			Reference: row["reference"],
			Closing:   closeTime,
			Title:     row["title"],
			URL:       row["url"],
		})
	}

	log.Printf("parsed ballots, length %d\n", len(out))
	return out, errs
}

func parseCommittees(rows []map[string]string) ([]Committee, []error) {
	out := make([]Committee, 0, len(rows))
	var errs []error

	for i, row := range rows {
		role := strings.TrimSpace(row["role"])
		if !isMemberStatus(role) {
			continue
		}
		com := strings.TrimSpace(row["committee"])
		if com == "" {
			errs = append(errs, fmt.Errorf("row %d: missing committee reference", i+1))
			continue
		}

		out = append(out, Committee{
			Committee:    com,
			MemberStatus: role,
			Domain:       row["domain"],
		})
	}

	log.Printf("parsed committees, length %d\n", len(out))
	return out, errs
}
