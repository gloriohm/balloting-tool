package standards

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func WriteResultExcel(path string, standards []StandardExpanded) error {
	f := excelize.NewFile()
	sheet := "standarder"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Prosjektleder", "Undersøkelsesår", "NS-nummer", "Norsk tittel", "Engelsk tittel", "Komité", "Status komité", "År fastsatt"}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	row := 2

	for _, s := range standards {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "")
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), time.Now().Year())
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Reference)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.TitleNO)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.TitleEN)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.Committee)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), "")
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), s.Year)

		row++
	}

	filepath := path + "/aktualitetsundersokelsen.xlsx"
	if err := f.SaveAs(filepath); err != nil {
		return err
	}

	return nil
}

func WriteOutJSON(in any, filename string) error {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s.json", filename), data, 0644)
}
