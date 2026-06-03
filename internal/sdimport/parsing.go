package sdimport

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func (s *Standard) ParsePulicationDetails(pub Publication) error {
	s.PublicationDate = pub.PublicationDate

	pageCount, err := pub.GetPageNumber()
	if err != nil {
		log.Printf("could not get pages for publication %s\n", pub.Reference)
	}
	s.Pages = pageCount

	return nil
}

func (s *Standard) ParseNationalProjectDetails(proj Project) error {
	s.Reference = proj.Reference
	s.Title = proj.ParseTitles()
	abstract := proj.ParseAbstract("en")
	if abstract == "" {
		abstract = proj.ParseAbstract("no")
	}
	s.Abstract = abstract
	s.ICS = proj.ParseClassification("ICS")
	s.ParseCommittee(proj)

	return nil
}

func (s *Standard) ParseInternationalProjectDetails(proj Project) {
	edition, err := strconv.Atoi(proj.Edition)
	if err != nil {
		log.Printf("edition not a number: got %s", proj.Edition)
		edition = 1
	}
	s.Edition = edition
	if s.Abstract == "" {
		s.Abstract = proj.ParseAbstract("en")
	}

	s.SusDevGoals = proj.ParseClassification("SUSTAINABLE_DEVELOPMENT_GOAL")
}

func (s *Standard) ParseCommittee(proj Project) {
	if len(proj.Owner.DisplayName) > 0 {
		s.Owner = append(s.Owner, proj.Owner.DisplayName)
	}
	if len(proj.Developer.DisplayName) > 0 {
		s.Developer = append(s.Developer, proj.Developer.DisplayName)
	}
}

func (p *Publication) GetPageNumber() (int, error) {
	for _, item := range p.ReleaseItems {
		if item.Type == "STANDARD" && item.Format == "PDF" {
			return item.Pages, nil
		}
	}

	return 0, fmt.Errorf("no main pdf")
}

func (p *Project) ParseTitles() []Title {
	var out []Title

	for _, t := range p.Title {
		switch t.Language {
		case "no", "en":
			out = append(out, t)
		}
	}

	return out
}

func (p *Project) ParseAbstract(langCode string) string {
	for _, a := range p.Abstract {
		if a.Language == langCode {
			switch a.Format {
			case "text/plain":
				return a.Content
			case "text/html":
				abstract, err := decodeHTML(a.Content)
				if err != nil {
					log.Printf("failed to decode abstract for project %s; using encoded text %s\n", p.Reference, err)
					return a.Content
				}
				return abstract
			default:
				return a.Content
			}
		}
	}

	return ""
}

func decodeHTML(in string) (string, error) {
	doc, err := html.Parse(strings.NewReader(in))
	if err != nil {
		return "", err
	}

	var b strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			b.WriteByte(' ')
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	return strings.TrimSpace(b.String()), nil
}

func (p *Project) ParseClassification(target string) []string {
	var out []string
	for _, c := range p.Classifications {
		switch target {
		case "SUSTAINABLE_DEVELOPMENT_GOAL":
			out = append(out, parseSusDevGoal(c.Value))
		case "ICS":
			out = append(out, c.Value)
		default:
			out = append(out, c.Value)
		}
	}

	return out
}

func parseSusDevGoal(value string) string {
	parsedID, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("could not parse id for goal with value %s", value)
		return value
	}

	goal, ok := sustainableDevelopmentGoals[parsedID]
	if ok {
		return goal
	} else {
		return value
	}
}

func (p *Project) GetRelationURN(relationType string) string {
	for _, r := range p.ProjectRelations {
		if r.Type == relationType {
			if relationType == "ADOPTED_FROM" {
				return r.ExternalProject.ProjectID
			} else {
				return r.URN
			}
		}
	}

	return ""
}
