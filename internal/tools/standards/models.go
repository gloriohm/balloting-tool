package standards

import (
	"ballot-tool/internal/api/sdimport"
	"log"
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

type Standard struct {
	Reference       string
	Title           []Title
	Abstract        string
	ICS             []string
	SusDevGoals     []string
	Edition         int
	Pages           int
	Developer       []string
	Owner           []string
	PublicationDate string
}

type Title struct {
	Language string
	Value    string
}

type StandardCore struct {
	Reference string
	Language  string
	Title     string
	URN       string
}

type StandardFile struct {
	StandardCore
	HasFile bool
}

type AktualitetStandard struct {
	Reference string
	TitleNO   string
	TitleEN   string
	Committee string
	Year      int
}

type AktualitetCommittee struct {
	ProjectManager  string
	CommitteeStatus string
}

type Aktualitet struct {
	AktualitetStandard
	AktualitetCommittee
}

func (s *Standard) ParsePulicationDetails(pub sdimport.Publication) error {
	s.PublicationDate = pub.PublicationDate

	pageCount, err := pub.GetPageNumber()
	if err != nil {
		log.Printf("could not get pages for publication %s\n", pub.Reference)
	}
	s.Pages = pageCount

	return nil
}

func (s *Standard) ParseNationalProjectDetails(proj sdimport.Project) error {
	s.Reference = proj.Reference
	/* s.Title = proj.ParseTitles()
	abstract := proj.ParseAbstract("en")
	if abstract == "" {
		abstract = proj.ParseAbstract("no")
	}
	s.Abstract = abstract
	s.ICS = proj.ParseClassification("ICS")
	s.ParseCommittee(proj) */

	return nil
}

func (s *Standard) ParseInternationalProjectDetails(proj sdimport.Project) {
	/* s.Edition = edition
	if s.Abstract == "" {
		s.Abstract = proj.ParseAbstract("en")
	}

	s.ParseCommittee(proj) */

	s.SusDevGoals = proj.ParseClassification("SUSTAINABLE_DEVELOPMENT_GOAL")
}
