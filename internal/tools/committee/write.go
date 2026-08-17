package committee

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func WriteEngagementsExcel(path string, engagements []Engagement) error {
	f := excelize.NewFile()
	sheet := "engagements"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Komite", "Fornavn", "Etternavn", "E-post", "Innmeldt", "Utmeldt"}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	row := 2

	for _, e := range engagements {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), e.Committee)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), e.FirstName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), e.LastName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), e.Email)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), e.From)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), e.To)

		row++
	}

	filepath := path + "/engagements.xlsx"
	if err := f.SaveAs(filepath); err != nil {
		return err
	}

	return nil
}
