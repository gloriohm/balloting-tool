package sdimport

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func (c *Client) PrintProject(urn string) error {
	proj, err := c.GetProject(urn)
	if err != nil {
		return fmt.Errorf("failed to print project: %w", err)
	}

	fmt.Println(proj)
	return nil
}

func (c *Client) PrintPublication(urn string) error {
	pub, err := c.GetPublication(urn)
	if err != nil {
		return fmt.Errorf("failed to print publication: %w", err)
	}

	fmt.Println(pub)
	return nil
}

func (c *Client) GetStandards() error {
	var standards []Standard
	pubs, err := c.GetPublications()
	if err != nil {
		return err
	}

	for _, p := range pubs {
		standard, err := c.CreateStandardFromPublication(p)
		if err != nil {
			log.Printf("error generating standard for publication %s: %s", p.Reference, err)
			continue
		}
		standards = append(standards, standard)
	}

	return WriteOutJSON(standards)
}

func (c *Client) CreateStandardFromPublication(pub Publication) (Standard, error) {
	var standard Standard

	urn := pub.ProjectID.URN
	nationalProject, err := c.GetProject(urn)
	if err != nil {
		return standard, fmt.Errorf("failed to get project for urn %s: %w", urn, err)
	}

	standard.ParsePulicationDetails(pub)
	standard.ParseNationalProjectDetails(nationalProject)

	return standard, nil
}

func WriteOutJSON(standards []Standard) error {
	data, err := json.MarshalIndent(standards, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile("result.json", data, 0644)
}
