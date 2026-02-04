package io

import (
	"ballot-tool/internal/models"
	"bufio"
	"fmt"
	"os"
)

func WriteBallotsTXT(path string, ballots []models.Ballot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, b := range ballots {
		_, err := fmt.Fprintf(w, "%s\t%s\n", b.Committee, b.Closing)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteCommitteesTXT(path string, coms []models.Committee) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, c := range coms {
		_, err := fmt.Fprintf(w, "%s\t%s\n", c.Committee, c.MemberStatus)
		if err != nil {
			return err
		}
	}
	return nil
}
