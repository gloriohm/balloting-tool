package ballot

import (
	"bufio"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func WriteBallotsTXT(path string, ballots []Ballot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, b := range ballots {
		_, err := fmt.Fprintf(w, "%s\t%s\n", b.Committee, b.Closing)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteCommitteesTXT(path string, coms []Committee) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, c := range coms {
		_, err := fmt.Fprintf(w, "%s\t%s\n", c.Committee, c.MemberStatus)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteBallotWithRoleXLSX(path string, rows []BallotWithRole, centralizedVoters []string) error {
	f := excelize.NewFile()
	sheets := []string{"Utestående"}
	f.SetSheetName("Sheet1", sheets[0])

	for i := range centralizedVoters {
		newSheet := fmt.Sprintf("Centralized%d", i+1)
		f.NewSheet(newSheet)
		sheets = append(sheets, newSheet)
	}

	splitRows := filterRowsByVoter(rows, centralizedVoters)

	for i, rows := range splitRows {
		if err := setBallotCells(f, sheets[i], rows); err != nil {
			return err
		}
	}

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}
