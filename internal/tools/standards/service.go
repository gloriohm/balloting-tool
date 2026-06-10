package standards

import (
	"ballot-tool/internal/api/sdimport"
	"ballot-tool/internal/utils/read"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Service struct {
	sdimport *sdimport.Client
}

func NewService(sdimport *sdimport.Client) *Service {
	return &Service{sdimport: sdimport}
}

func GenerateAktualitetList(pathToFolder string) error {
	data, err := read.LoadTabularDataFromFolder(pathToFolder)
	if err != nil {
		return err
	}

	all := rowToStandard(data)
	standards, _ := filterByAdoptionType(all, "national", false)

	filtered := make(map[string]StandardCore)
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

	var aktStandard []AktualitetStandard
	params := sdimport.NewParameters("", "")
	client := sdimport.NewClient(false, params)
	for _, s := range filtered {
		id := "sn:proj:" + s.URN
		proj, err := client.GetProject(id)
		if err != nil {
			log.Printf("failed fetching metadata for %s: %s\n", s.Reference, err)
			continue
		}

		standard := projToAktualitetStandard(proj)

		aktStandard = append(aktStandard, standard)
	}

	return WriteAktualitetExcel(pathToFolder, aktStandard)
}

func projToAktualitetStandard(proj sdimport.Project) AktualitetStandard {
	var out AktualitetStandard
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
	data, err := read.LoadTabularDataFromFolder(pathToFolder)
	if err != nil {
		return err
	}

	count := Counter{}
	all := rowToStandard(data)
	standards, resultLog := filterByAdoptionType(all, selection, nsOnly)

	filtered := make(map[string]struct{})
	var includeDupes []StandardCore

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

func (s *Service) GetStandards() error {
	var standards []Standard
	pubs, err := s.sdimport.GetPublications()
	if err != nil {
		return err
	}

	for _, p := range pubs {
		standard, err := s.CreateStandardFromPublication(p)
		if err != nil {
			s.sdimport.Logger.Info(fmt.Sprintf("%s: %s", p.Reference, err))
			continue
		}
		standards = append(standards, standard)
	}

	return WriteOutJSON(standards, "get_standard_result")
}

func (s *Service) CreateStandardFromPublication(pub sdimport.Publication) (Standard, error) {
	var standard Standard

	urn := pub.ProjectID.URN
	nationalProject, err := s.sdimport.GetProject(urn)
	if err != nil {
		return standard, fmt.Errorf("failed to get project for urn %s: %w", urn, err)
	}

	standard.ParsePulicationDetails(pub)
	standard.ParseNationalProjectDetails(nationalProject)

	//terrible code, refactor
	parentURN := nationalProject.GetRelationURN("ADOPTED_FROM")
	if parentURN != "" {
		parent, err := s.sdimport.GetProject(parentURN)
		if err != nil {
			return standard, fmt.Errorf("failed to get parent project with urn %s: %w", urn, err)
		}
		standard.ParseInternationalProjectDetails(parent)

		/* grandparentURN := parent.GetRelationURN("VIENNA_AGREEMENT")
		if grandparentURN != "" {
			grandparent, err := s.sdimport.GetProject(grandparentURN)
			if err != nil {
				return standard, fmt.Errorf("failed to get grandparent project with urn %s: %w", urn, err)
			}
			standard.ParseCommittee(grandparent)
		} */
	}

	return standard, nil
}

func (s *Service) FindStandardsWithXML(path string) error {
	data, err := read.LoadTabularDataFromFile(fmt.Sprintf("%s/full.csv", path))
	if err != nil {
		return fmt.Errorf("error loading data at path %s: %w", path, err)
	}

	standards := rowToStandard(data)
	var out []StandardFile
	for _, std := range standards {
		urn := fmt.Sprintf("sn:proj:%s", std.URN)
		pub, err := s.sdimport.GetPublicationByProject(urn, "PUBLISHED")
		if err != nil {
			log.Printf("could not get publication for %s: %s\n", urn, err)
			continue
		}
		_, err = pub.GetReleaseItem("STANDARD", "XML")
		withFile := StandardFile{StandardCore: std}
		if err == nil {
			withFile.HasFile = true
		} else {
			withFile.HasFile = false
		}
		out = append(out, withFile)
	}

	if err := WriteHasFileExcel(path, out); err != nil {
		return err
	}

	return nil
}

func (s *Service) DownloadFiles(in, out string) error {
	data, err := read.LoadTabularDataFromFile(fmt.Sprintf("%s/test.csv", in))
	if err != nil {
		return fmt.Errorf("error loading data at path %s: %w", in, err)
	}

	standards := rowToStandard(data)

	for _, std := range standards {
		urn := fmt.Sprintf("sn:proj:%s", std.URN)
		pub, err := s.sdimport.GetPublicationByProject(urn, "PUBLISHED")
		if err != nil {
			log.Printf("could not get publication for %s: %s\n", urn, err)
			continue
		}
		rel, err := pub.GetReleaseItem("STANDARD", "XML")
		if err != nil {
			log.Printf("could not get release item for %s: %s\n", pub.Reference, err)
			continue
		}

		if err := s.sdimport.GetFile(rel.ContentRef, out); err != nil {
			log.Printf("could not save file %s: %s\n", rel.ContentRef.FileName, err)
			continue
		}
	}

	return nil
}
