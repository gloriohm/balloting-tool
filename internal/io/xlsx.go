package io

import (
	"ballot-tool/internal/models"
	"ballot-tool/internal/processing"
	"ballot-tool/internal/utils"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadXLSX(r io.Reader, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("xlsx: no sheets found")
		}
		sheetName = sheets[0]
	}

	allRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(allRows) == 0 {
		return []map[string]string{}, nil
	}

	rawHeader := allRows[0]
	headers := make([]string, len(rawHeader))
	for i, h := range rawHeader {
		nh := utils.NormalizeString(h)
		if canon, ok := utils.HeaderAliases[nh]; ok {
			headers[i] = canon
		} else {
			headers[i] = nh
		}
	}

	rows := make([]map[string]string, 0, max(0, len(allRows)-1))
	for _, rec := range allRows[1:] {
		empty := true
		for _, c := range rec {
			if strings.TrimSpace(c) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}

		m := make(map[string]string, len(headers))
		for i := 0; i < len(rec) && i < len(headers); i++ {
			m[headers[i]] = strings.TrimSpace(rec[i])
		}
		for i := len(rec); i < len(headers); i++ {
			m[headers[i]] = ""
		}
		rows = append(rows, m)
	}

	return rows, nil
}

func WriteBallotWithRoleXLSX(path string, rows []models.BallotWithRole, centralizedVoters []string) error {
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

func setBallotCells(f *excelize.File, sheet string, rows []models.BallotWithRole) error {
	processing.SortByCloses(rows)

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

func filterRowsByVoter(rows []models.BallotWithRole, voters []string) [][]models.BallotWithRole {
	voterIdx, _ := utils.IndexStrings(voters, 1)

	out := make([][]models.BallotWithRole, len(voters)+1)

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
