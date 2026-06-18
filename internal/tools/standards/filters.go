package standards

import (
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var suffixRef = regexp.MustCompile(`:\d{4}/[^/]+:\d{4}$`)

var (
	norskStandardNational  = []string{"NS "}
	norsokPrefix           = []string{"NORSOK"}
	norskStandardAdoption  = []string{"NS-"}
	technicalOtherNational = []string{"SN_", "SN-NSPEK", "P-", "NHS"}
	otherAdoptions         = []string{"SN-CEN", "SN-ISO"}

	allNorskStandardPrefixes = append(append([]string{}, norskStandardNational...), norskStandardAdoption...)
	allAdoptionPrefixes      = append(append([]string{}, norskStandardAdoption...), otherAdoptions...)
	allPureNationalPrefixes  = append(append([]string{}, norskStandardNational...), technicalOtherNational...)
)

func isAddons(ref string) bool {
	// returns true if input ends with :year/{suffix}:{year}
	if suffixRef.MatchString(ref) {
		return true
	}
	return false
}

func hasLanguageCodeInReference(ref string) bool {
	re := regexp.MustCompile(`\.[A-Z]:`)
	if re.MatchString(ref) {
		return true
	}
	return false
}

func isAllowedLanguage(lang string, allowed []string) bool {
	// returns true if input is in list of allowed languages
	if slices.Contains(allowed, lang) {
		return true
	}

	return false
}

func hasAllowedPrefix(ref string, allowed []string) bool {
	for _, a := range allowed {
		if strings.HasPrefix(ref, a) {
			return true
		}
	}

	return false
}

func divisibleByFive(ref string) bool {
	year, err := getYearFromReference(ref)
	if err != nil {
		return false
	}
	currentYear := time.Now().Year()

	diff := currentYear - year
	if diff > 0 && diff%5 == 0 {
		return true
	}

	return false
}

func getYearFromReference(ref string) (int, error) {
	yearString := ref[len(ref)-4:]
	year, err := strconv.Atoi(yearString)
	if err != nil {
		log.Printf("last four characters not a valid year %s", ref)
		return year, err
	}
	return year, nil
}

func filterByAdoptionType(standards []StandardCore, choice string, nsOnly bool) ([]StandardCore, Log) {
	var filtered []StandardCore
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

func createDownloadJob(job string) []Target {
	var targets []Target
	switch job {
	case "all":
		targets = append(targets, pdfTarget())
		targets = append(targets, xmlTarget())
		targets = append(targets, wordTarget())
		targets = append(targets, anyTarget())
	case "source":
		targets = append(targets, xmlTarget())
		targets = append(targets, wordTarget())
	case "xml":
		targets = append(targets, xmlTarget())
	case "word":
		targets = append(targets, wordTarget())
	case "pdf":
		targets = append(targets, pdfTarget())
	case "other":
		targets = append(targets, anyTarget())
	case "pdfsource":
		targets = append(targets, xmlTarget())
		targets = append(targets, wordTarget())
		targets = append(targets, pdfTarget())
	}

	return targets
}
