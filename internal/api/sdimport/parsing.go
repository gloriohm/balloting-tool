package sdimport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func (p *Project) ParseEdition() (int, error) {
	return strconv.Atoi(p.Edition)
}

func (p *Project) ParseCommittee(level string) string {
	switch level {
	case "owner":
		return p.Owner.DisplayName
	case "developer":
		return p.Developer.DisplayName
	default:
		return ""
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
		if c.Type == target {
			switch target {
			case "SUSTAINABLE_DEVELOPMENT_GOAL":
				out = append(out, parseSusDevGoal(c.Value))
			case "ICS":
				out = append(out, c.Value)
			}
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
			if relationType == "ADOPTED_FROM" || relationType == "VIENNA_AGREEMENT" {
				return r.ExternalProject.ProjectID
			} else {
				return r.URN
			}
		}
	}

	return ""
}

func (p *Publication) GetReleaseItem(itemType ReleaseItemType, itemFormat ReleaseItemFormat) (ReleaseItem, error) {
	for _, r := range p.ReleaseItems {
		if r.Type == string(itemType) && r.Format == string(itemFormat) {
			return r, nil
		}
	}

	return ReleaseItem{}, fmt.Errorf("found no items of type %s and format %s on publication %s", itemType, itemFormat, p.Reference)
}

func dumpResponse(resp *http.Response) (Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}

func (p *Project) getPublicationURNs() []string {
	var pubURNs []string
	for _, p := range p.PubLink {
		pubURNs = append(pubURNs, p.URN)
	}

	return pubURNs
}
