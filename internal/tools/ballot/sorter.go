package ballot

import (
	"ballot-tool/internal/utils/normalization"

	"fmt"
	"log"
	"sort"

	"github.com/xuri/excelize/v2"
)

func sortByCloses(rows []BallotMatched) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Ballot.Closing.Before(rows[j].Ballot.Closing)
	})
}

func setBallotCells(f *excelize.File, sheet string, rows []BallotMatched) error {
	sortByCloses(rows)

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
		if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), br.Voter.FirstName); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), br.Voter.LastName); err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), br.Voter.Email); err != nil {
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

func filterRowsByVoter(rows []BallotMatched, voters []string) [][]BallotMatched {
	voterIdx, _ := normalization.IndexStrings(voters, 1)

	out := make([][]BallotMatched, len(voters)+1)

	for _, row := range rows {
		email := normalization.ToLowerCase(row.Voter.Email)

		if i, ok := voterIdx[email]; ok {
			out[i] = append(out[i], row)
		} else {
			out[0] = append(out[0], row)
		}
	}

	return out
}
