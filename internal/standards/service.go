package standards

import (
	"ballot-tool/internal/sdimport"
	"ballot-tool/internal/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type Counter struct {
	Lang      int
	Addon     int
	LangCode  int
	Selection int
	Duplicate int
}

type Log struct {
	In   int
	Out  int
	Diff int
}

func GenerateAktualitetList(pathToFolder string) error {
	data, err := utils.LoadTabularDataFromFolder(pathToFolder)
	if err != nil {
		return err
	}

	all := rowToStandard(data)
	standards, _ := filterByAdoptionType(all, "national", false)

	filtered := make(map[string]Standard)
	for _, s := range standards {
		if !divisibleByFive(s.Reference) {
			log.Printf("[Aktualitet] skipped %s due to not being released five year period\n", s.Reference)
			continue
		}
		if isAddons(s.Reference) {
			//log.Printf("[Addon] skipped %s due to being an addon\n", s.Reference)
			continue
		}
		if hasLanguageCodeInReference(s.Reference) {
			//log.Printf("[Lang Code] skipped %s due to language code in reference\n", s.Reference)
			continue
		}

		filtered[s.Reference] = s
	}

	var expanded []StandardExpanded
	params := sdimport.NewParameters("", "")
	client := sdimport.NewClient(false, params)
	for _, s := range filtered {
		id := "sn:proj:" + s.URN
		proj, err := client.GetProject(id)
		if err != nil {
			log.Printf("failed fetching metadata for %s: %s\n", s.Reference, err)
			continue
		}

		standard := projToExpanded(proj)

		expanded = append(expanded, standard)
	}

	return WriteResultExcel(pathToFolder, expanded)
}

func projToExpanded(proj sdimport.Project) StandardExpanded {
	var out StandardExpanded
	out.Reference = proj.Reference
	titles := proj.ParseTitles()
	for _, t := range titles {
		switch t.Language {
		case "no":
			out.TitleNO = t.Value
		case "en":
			out.TitleEN = t.Value
		}
	}
	out.Committee = proj.Owner.DisplayName
	year, _ := getYearFromReference(proj.Reference)
	out.Year = year

	return out
}

func CountTotalUniqueProducts(pathToFolder, selection string, nsOnly bool) error {
	data, err := utils.LoadTabularDataFromFolder(pathToFolder)
	if err != nil {
		return err
	}

	count := Counter{}
	all := rowToStandard(data)
	standards, resultLog := filterByAdoptionType(all, selection, nsOnly)

	filtered := make(map[string]struct{})
	var includeDupes []Standard

	for _, s := range standards {
		if !isAllowedLanguage(s.Language, []string{"en", "no", "nb", "nn"}) {
			//log.Printf("[Lang] skipped %s due to language %s\n", s.Reference, s.Language)
			count.Lang++
			continue
		}
		if isAddons(s.Reference) {
			//log.Printf("[Addon] skipped %s due to being an addon\n", s.Reference)
			count.Addon++
			continue
		}
		if hasLanguageCodeInReference(s.Reference) {
			//log.Printf("[Lang Code] skipped %s due to language code in reference\n", s.Reference)
			count.LangCode++
			continue
		}

		_, exists := filtered[s.Reference]
		if exists {
			count.Duplicate++
		}
		includeDupes = append(includeDupes, s)
		filtered[s.Reference] = struct{}{}
	}

	log.Printf("Total in: %d; Total after selection filtration: %d; diff: %d", resultLog.In, resultLog.Out, resultLog.Diff)
	log.Printf("%d standards in selection\n", len(filtered))
	log.Printf("%d standards inkludert oversettelser\n", len(includeDupes))
	log.Printf("total discarded:\nLang: %d\nAddon: %d\nLang code in ref: %d\nDuplicate reference: %d\n", count.Lang, count.Addon, count.LangCode, count.Duplicate)

	if err := WriteResultTXT(pathToFolder, filtered, count); err != nil {
		return err
	}
	return nil
}

func filterByAdoptionType(standards []Standard, choice string, nsOnly bool) ([]Standard, Log) {
	var filtered []Standard
	var allowedPrefixes []string
	switch choice {
	case "national":
		if nsOnly {
			allowedPrefixes = norskStandardNational
		} else {
			allowedPrefixes = allPureNationalPrefixes
		}
	case "adoption":
		if nsOnly {
			allowedPrefixes = norskStandardAdoption
		} else {
			allowedPrefixes = allAdoptionPrefixes
		}
	case "norsok":
		allowedPrefixes = norsokPrefix
	case "all":
		if nsOnly {
			allowedPrefixes = allNorskStandardPrefixes
		} else {
			return standards, Log{In: len(standards), Out: len(standards), Diff: 0}
		}
	default:
		return standards, Log{In: len(standards), Out: len(standards), Diff: 0}
	}

	for _, s := range standards {
		if hasAllowedPrefix(s.Reference, allowedPrefixes) {
			filtered = append(filtered, s)
		}
	}

	log.Println(allowedPrefixes)

	resultLog := Log{In: len(standards), Out: len(filtered), Diff: len(standards) - len(filtered)}

	return filtered, resultLog
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
