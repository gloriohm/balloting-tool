package standards

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

func WriteAktualitetExcel(path string, standards []AktualitetStandard) error {
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

func WriteHasFileExcel(path string, standards map[bool][]string) error {
	f := excelize.NewFile()
	sheet := "resultat"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"referanse", "har_xml"}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	row := 2

	for col, s := range standards {
		for _, ref := range s {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), ref)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), col)
			row++
		}
	}

	if err := f.SaveAs(path); err != nil {
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

func WriteResultTXT(path string, references map[string]struct{}, count Counter) error {
	out := filepath.Join(path, "result.txt")
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = fmt.Fprintf(w, "total discarded:\nLang: %d\nAddon: %d\nLang code in ref: %d\nDuplicate reference: %d\n", count.Lang, count.Addon, count.LangCode, count.Duplicate)

	for key := range references {
		_, err := fmt.Fprintf(w, "%s\n", key)
		if err != nil {
			return err
		}
	}
	return nil
}
