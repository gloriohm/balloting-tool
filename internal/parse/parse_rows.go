package parse

import (
	"ballot-tool/internal/models"
	"ballot-tool/internal/utils"
	"fmt"
	"strings"
)

func ParseRoles(rows []map[string]string) ([]models.Role, []error) {
	out := make([]models.Role, 0, len(rows))
	var errs []error
	for i, row := range rows {
		rawRole := row["role"]
		if !utils.IsVoterRole(rawRole) {
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

		out = append(out, models.Role{
			Committee: com,
			FirstName: row["first_name"],
			LastName:  row["last_name"],
			Email:     row["email"],
		})
	}

	return out, errs
}

func ParseBallots(rows []map[string]string) ([]models.Ballot, []error) {
	out := make([]models.Ballot, 0, len(rows))
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

		out = append(out, models.Ballot{
			Source:    row["source"],
			Committee: com,
			Reference: row["reference"],
			Closing:   closeTime,
			Title:     row["title"],
			URL:       row["url"],
		})
	}

	return out, errs
}
