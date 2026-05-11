package standards

import (
	"ballot-tool/internal/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
