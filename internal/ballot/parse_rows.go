package ballot

import (
	"ballot-tool/internal/utils"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ParseRoles(rows []map[string]string) ([]Role, []error) {
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

func ParseBallots(rows []map[string]string) ([]Ballot, []error) {
	out := make([]Ballot, 0, len(rows))
	var errs []error
	for i, row := range rows {
		rawVoter := strings.TrimSpace(row["role"])
		if !isVoterRole(rawVoter) {
			errs = append(errs, fmt.Errorf("row %d: ballot role is not Voter", i+1))
			continue
		}
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

func ParseCommittees(rows []map[string]string) ([]Committee, []error) {
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

func setBallotCells(f *excelize.File, sheet string, rows []BallotWithRole) error {
	SortByCloses(rows)

	headers := []string{
		"Committee",
		"Reference",
		"Closes",
		"First Name",
		"Last Name",
		"Email",
		"URL",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return fmt.Errorf("set header cell %s: %w", cell, err)
		}
	}

	for rIdx, br := range rows {
		rowNum := rIdx + 2
		if err := f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), br.Ballot.Committee); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), br.Ballot.Reference); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), br.Ballot.Closing); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), br.Role.FirstName); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), br.Role.LastName); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), br.Role.Email); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), br.Ballot.URL); err != nil {
			return err
		}
	}

	if err := f.SetColWidth(sheet, "A", "G", 20); err != nil {
		log.Println(err)
	}

	return nil
}

func filterRowsByVoter(rows []BallotWithRole, voters []string) [][]BallotWithRole {
	voterIdx, _ := utils.IndexStrings(voters, 1)

	out := make([][]BallotWithRole, len(voters)+1)

	for _, row := range rows {
		email := utils.ToLowerCase(row.Role.Email)

		if i, ok := voterIdx[email]; ok {
			out[i] = append(out[i], row)
		} else {
			out[0] = append(out[0], row)
		}
	}

	return out
}
