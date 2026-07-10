package ballot

import (
	"ballot-tool/internal/filereader"
	"bufio"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func writeMissing(path string, ballots []filereader.BallotRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, b := range ballots {
		_, err := fmt.Fprintf(w, "%s\t%s\n", b.Committee, b.ClosingDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeBallotsMatchedXLSX(path string, rows []BallotMatched, centralizedVoters []string) error {
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
