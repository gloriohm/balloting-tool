package standards

import (
	"ballot-tool/internal/api/sdimport"
	"ballot-tool/internal/filereader"
	"ballot-tool/internal/utils/config"
	"ballot-tool/internal/utils/normalization"
	"ballot-tool/internal/utils/read"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Service struct {
	sdimport *sdimport.Client
	cfg      *config.Config
}

func NewService(sdimport *sdimport.Client, cfg *config.Config) *Service {
	return &Service{sdimport: sdimport, cfg: cfg}
}

func (s *Service) GenerateAktualitetList(pathToFolder string) error {
	data, err := read.LoadTabularDataFromFolder(pathToFolder)
	if err != nil {
		return err
	}

	all := rowToStandard(data)
	standards, _ := filterByAdoptionType(all, "national", false)

	filtered := make(map[string]StandardCore)
	for _, std := range standards {
		if !divisibleByFive(std.Reference) {
			log.Printf("[Aktualitet] skipped %s due to not being released five year period\n", std.Reference)
			continue
		}
		if isAddons(std.Reference) {
			//log.Printf("[Addon] skipped %s due to being an addon\n", std.Reference)
			continue
		}
		if hasLanguageCodeInReference(std.Reference) {
			//log.Printf("[Lang Code] skipped %s due to language code in reference\n", std.Reference)
			continue
		}

		filtered[std.Reference] = std
	}

	var aktStandard []AktualitetStandard
	for _, std := range filtered {
		id := "sn:proj:" + std.URN
		proj, err := s.sdimport.GetProject(id)
		if err != nil {
			log.Printf("failed fetching metadata for %s: %s\n", std.Reference, err)
			continue
		}

		standard := projToAktualitetStandard(proj)

		aktStandard = append(aktStandard, standard)
	}

	return WriteAktualitetExcel(pathToFolder, aktStandard)
}

func (s *Service) CountTotalUniqueProducts(pathToFolder, selection string, nsOnly bool) error {
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
		pub, err := s.sdimport.GetPublicationByProject(urn, "published")
		if err != nil {
			log.Printf("could not get publication for %s: %s\n", urn, err)
			continue
		}
		_, err = pub.GetReleaseItems(sdimport.ReleaseItemTypeStandard, sdimport.ReleaseItemFormatXML)
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

func (s *Service) DownloadFiles(in string, opts string) error {
	path := filepath.Join(s.cfg.InputPath, in)
	filter, err := filereader.NewFilter("stage==Working")
	standards, err := filereader.LoadStandardsDashboard(path, filter)
	if err != nil {
		return fmt.Errorf("error loading data at path %s: %w", in, err)
	}

	targets := createDownloadJob(opts)

	downloadsFolder := filepath.Join(s.cfg.OutputPath, "downloads")

	for _, std := range standards {
		urn := fmt.Sprintf("sn:proj:%s", std.ImportID)
		pub, err := s.sdimport.GetPublicationByProject(urn, "published")
		if err != nil {
			log.Printf("could not get publication for %s: %s\n", urn, err)
			continue
		}

		stdFolder := filepath.Join(downloadsFolder, normalization.NormalizeString(std.Reference))
		err = os.MkdirAll(stdFolder, 0755)
		if err != nil {
			return err
		}

		for _, t := range targets {
			items, err := pub.GetReleaseItems(t.itemType, t.itemFormat)
			if err != nil {
				log.Printf("%s: %s\n", pub.Reference, err)
				continue
			}

			for _, i := range items {
				lang := extractLanguage(i.Language)
				languageFolder := filepath.Join(stdFolder, lang)
				attachmentsFolder := filepath.Join(languageFolder, "attachments")
				if err := os.MkdirAll(languageFolder, 0755); err != nil {
					return err
				}

				switch i.Type {
				case string(sdimport.ReleaseItemTypeOther):
					if err := os.MkdirAll(attachmentsFolder, 0755); err != nil {
						return err
					}
					if err := s.sdimport.GetFile(i.ContentRef, attachmentsFolder); err != nil {
						log.Printf("could not save file %s: %s\n", i.ContentRef.FileName, err)
						continue
					}
				default:
					if err := s.sdimport.GetFile(i.ContentRef, languageFolder); err != nil {
						log.Printf("could not save file %s: %s\n", i.ContentRef.FileName, err)
						continue
					}
				}
			}
		}
	}

	return nil
}
