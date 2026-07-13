package filereader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func parseCSV[T any](path string, filters Filters, mapper func(Row) (T, error)) ([]T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rawHeaders, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headers := normalizeHeaders(rawHeaders)

	var out []T
	rowNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		rowNum++

		if err != nil {
			return nil, err
		}

		row := make(Row)

		for i, header := range headers {
			if i < len(record) {
				row[header] = strings.TrimSpace(record[i])
			}
		}

		if !passesFilters(row, filters) {
			continue
		}

		item, err := mapper(row)
		if errors.Is(err, ErrMalformedData) || errors.Is(err, ErrMissingData) {
			// upgrade to actual logging
			//log.Printf("%s: skipping row with index %d: %s", path, rowNum, err)
			continue
		}
		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

func parseExcel[T any](path string, filter Filters, mapper func(Row) (T, error)) ([]T, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrEmptyFile
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrEmptyFile
	}

	headers := normalizeHeaders(rows[0])

	var out []T

	for i, record := range rows[1:] {
		rowNum := i + 2

		row := make(Row)

		for colIdx, header := range headers {
			if colIdx < len(record) {
				row[header] = strings.TrimSpace(record[colIdx])
			}
		}

		if !passesFilters(row, filter) {
			continue
		}

		item, err := mapper(row)
		if errors.Is(err, ErrMalformedData) || errors.Is(err, ErrMissingData) {
			// upgrade to actual logging
			log.Printf("%s: skipping row with index %d: %s", path, rowNum, err)
			continue
		}
		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

func parseStandardDashboardRow(row Row) (StandardDashboardRow, error) {
	return StandardDashboardRow{
		ImportID:  row["id"],
		Reference: row["reference"],
		PubStatus: row["pub_status"],
		Language:  row["lang"],
		Title:     row["title"],
		SDO:       row["sdo"],
	}, nil
}

func parseBallotRow(row Row) (BallotRow, error) {
	openingDate, err := time.Parse("2006-01-02", row["opening_date"])
	if err != nil {
		return BallotRow{}, ErrMalformedData
	}

	closingDate, err := time.Parse("2006-01-02", row["closing_date"])
	if err != nil {
		return BallotRow{}, ErrMalformedData
	}

	ballot := BallotRow{
		BallotType:      row["type"],
		Committee:       row["committee"],
		BallotReference: row["reference"],
		OpeningDate:     openingDate,
		ClosingDate:     closingDate,
		BallotName:      row["title"],
		URL:             row["url"],
	}

	return ballot, nil
}

func parseNationalEngagementsRow(row Row) (NationalEngagementsRow, error) {
	committee, err := parseCommittee(row)
	if err != nil {
		return NationalEngagementsRow{}, err
	}

	commitment, err := parseCommitment(row)
	if err != nil {
		return NationalEngagementsRow{}, err
	}

	out := NationalEngagementsRow{
		Committee:  committee,
		Commitment: commitment,
		Person:     parsePerson(row),
	}

	return out, nil
}

func parseCommittee(row Row) (Committee, error) {
	established, err := time.Parse("2006-01-02", row["committee_established_date"])
	if err != nil {
		return Committee{}, fmt.Errorf("committee_established: %w", ErrMalformedData)
	}

	status, err := parseCommitteeStatus(row["committee_status"])
	if err != nil {
		return Committee{}, fmt.Errorf("committee_status: %w", ErrMalformedData)
	}

	domain, err := parseCommitteeDomain(row["committee_domain"])
	if err != nil {
		return Committee{}, fmt.Errorf("committee_domain: %w", ErrMalformedData)
	}

	level, err := parseCommitteeLevel(row["committee_level"])
	if err != nil {
		return Committee{}, fmt.Errorf("committee_level: %w", ErrMalformedData)
	}

	committee := Committee{
		Reference:       row["committee_name"],
		Title:           row["committee_title"],
		Domain:          domain,
		Level:           level,
		Status:          status,
		Established:     established,
		MirrorCommittee: row["is_mirror_committee"] == "TRUE",
	}

	return committee, nil
}

func parseCommitment(row Row) (Commitment, error) {
	role := row["commitment_role"]
	if role == "" {
		return Commitment{}, fmt.Errorf("commitment_role: %w", ErrMissingData)
	}

	status, err := parseCommitmentStatus(row["commitment_status"])
	if err != nil {
		return Commitment{}, fmt.Errorf("commitment_status: %w", ErrMalformedData)
	}

	start, err := time.Parse("2006-01-02", row["commitment_from"])
	if err != nil {
		log.Printf("[WARNING]: commitment for %s - %s parsed without starting date", row["email"], row["committee_name"])
	}

	end, err := time.Parse("2006-01-02", row["commitment_to"])
	if err != nil {
		if row["commitment_to"] != "" {
			return Commitment{}, fmt.Errorf("commitment_to: %w", ErrMalformedData)
		}
	}

	commitment := Commitment{
		Role:    role,
		Status:  status,
		Start:   start,
		End:     end,
		Company: parseCommitmentCompany(row),
	}

	return commitment, nil
}

func parsePerson(row Row) Person {
	return Person{
		FirstName: row["first_name"],
		LastName:  row["last_name"],
		Email:     row["email"],
		Company:   parseEmployedByCompany(row),
	}
}

func parseCommitmentCompany(row Row) Company {
	return Company{
		Name:     row["commitment_company_name"],
		Category: row["commitment_company_category"],
		OrgForm:  row["commitment_company_organization_form"],
	}
}

func parseEmployedByCompany(row Row) Company {
	return Company{
		Name:     row["employed_by_company_name"],
		Category: row["employed_by_company_category"],
		OrgForm:  row["employed_by_company_organization_form"],
	}
}

func parseCommitteeLevel(s string) (CommitteeLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(CommitteeLevelTC):
		return CommitteeLevelTC, nil
	case string(CommitteeLevelSC):
		return CommitteeLevelSC, nil
	case string(CommitteeLevelWG):
		return CommitteeLevelWG, nil
	default:
		return "", ErrMalformedData
	}
}

func parseCommitteeDomain(s string) (CommitteeDomain, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(CommitteeDomainNational):
		return CommitteeDomainNational, nil
	case string(CommitteeDomainRegional):
		return CommitteeDomainRegional, nil
	case string(CommitteeDomainInternational):
		return CommitteeDomainInternational, nil
	default:
		return "", ErrMalformedData
	}
}

func parseCommitteeStatus(s string) (CommitteeStatus, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(CommitteeStatusActive):
		return CommitteeStatusActive, nil
	case string(CommitteeStatusInactive):
		return CommitteeStatusInactive, nil
	case string(CommitteeStatusTerminated):
		return CommitteeStatusTerminated, nil
	case string(CommitteeStatusSuspended):
		return CommitteeStatusSuspended, nil
	case string(CommitteeStatusInProgress):
		return CommitteeStatusInProgress, nil
	default:
		return "", ErrMalformedData
	}
}

func parseCommitmentStatus(s string) (CommitmentStatus, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(CommitmentStatusActive):
		return CommitmentStatusActive, nil
	case string(CommitmentStatusTerminated):
		return CommitmentStatusTerminated, nil
	default:
		return "", ErrMalformedData
	}
}
