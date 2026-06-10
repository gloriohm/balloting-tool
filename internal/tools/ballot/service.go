package ballot

import (
	"ballot-tool/internal/utils/read"
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func GetBallots(path string) ([]Ballot, error) {
	rawRows, err := read.LoadTabularDataFromFile(path)
	if err != nil {
		return nil, err
	}

	rowNormalizedHeaders := normalizeHeaders(rawRows, ballotHeaderAliases)

	ballots, errs := parseBallots(rowNormalizedHeaders)
	ErrorPrinter(errs, 5)

	return ballots, nil
}

func GetVoterRoles(path string) ([]Role, error) {
	rawRows, err := read.LoadTabularDataFromFile(path)
	if err != nil {
		return nil, err
	}

	rowNormalizedHeaders := normalizeHeaders(rawRows, rolesHeaderAliases)

	rolesParsed, roleErrs := parseRoles(rowNormalizedHeaders)
	ErrorPrinter(roleErrs, 5)

	return rolesParsed, nil
}

func GetOrgRoles(path string) ([]Committee, error) {
	orgRows, err := read.LoadTabularDataFromFile(path)
	if err != nil {
		return nil, err
	}

	orgParsed, orgErrs := parseCommittees(orgRows)
	ErrorPrinter(orgErrs, 5)

	return orgParsed, nil
}
func GetCombinedBallots(paths []string) ([]Ballot, error) {
	var combinedBallots []Ballot
	for _, path := range paths {
		ballots, err := GetBallots(path)
		if err != nil {
			return nil, err
		}
		combinedBallots = append(combinedBallots, ballots...)
	}

	return combinedBallots, nil
}

func JoinBallotRole(roles []Role, ballots []Ballot) ([]BallotWithRole, []Ballot) {
	roleCommitteeIdx := createRoleComIdx(roles)

	matches := make([]BallotWithRole, 0, len(ballots))
	missing := make([]Ballot, 0, len(ballots))
	for _, b := range ballots {
		match, ok := roleCommitteeIdx[b.Committee]
		if !ok {
			missing = append(missing, b)
			continue
		}

		matches = append(matches, BallotWithRole{
			Ballot: b,
			Role:   *match,
		})
	}

	log.Printf("matched %d ballots \n", len(matches))
	log.Printf("found %d ballots without voter \n", len(missing))

	return matches, missing
}

func JoinCommitteeRole(roles []Role, coms []Committee) []Committee {
	roleCommitteeIdx := createRoleComIdx(roles)

	missing := make([]Committee, 0, len(coms))

	for _, c := range coms {
		_, ok := roleCommitteeIdx[c.Committee]
		if !ok {
			missing = append(missing, c)
			continue
		}
	}

	log.Printf("found %d committees with P or O membership without Voter\n", len(missing))
	return missing
}

func WriteBallotsTXT(path string, ballots []Ballot) error {
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

func WriteCommitteesTXT(path string, coms []Committee) error {
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

func WriteBallotWithRoleXLSX(path string, rows []BallotWithRole, centralizedVoters []string) error {
	f := excelize.NewFile()
	sheets := []string{"Utestående"}
	f.SetSheetName("Sheet1", sheets[0])

	for i := range centralizedVoters {
		newSheet := fmt.Sprintf("Centralized%d", i+1)
		f.NewSheet(newSheet)
		sheets = append(sheets, newSheet)
	}

	splitRows := filterRowsByVoter(rows, centralizedVoters)

	for i, rows := range splitRows {
		if err := setBallotCells(f, sheets[i], rows); err != nil {
			return err
		}
	}

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}
