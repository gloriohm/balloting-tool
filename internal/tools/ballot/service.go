package ballot

import (
	"ballot-tool/internal/filereader"
	"ballot-tool/internal/utils/config"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

func GenerateBallotReport(cfg *config.Config) error {
	isoPath := filepath.Join(cfg.InputPath, cfg.Files.Ballot1)
	cenPath := filepath.Join(cfg.InputPath, cfg.Files.Ballot2)
	rolesPath := filepath.Join(cfg.InputPath, cfg.Files.Voters)

	iso, err := filereader.LoadBallots(isoPath, filereader.Filters{})
	if err != nil {
		return fmt.Errorf("failed getting iso ballots: %w", err)
	}

	log.Printf("retreived %d iso ballots", len(iso))

	cen, err := filereader.LoadBallots(cenPath, filereader.Filters{})
	if err != nil {
		return fmt.Errorf("failed getting cen ballots: %w", err)
	}

	log.Printf("retreived %d cen ballots", len(cen))

	ballots := make([]filereader.BallotRow, 0, len(iso)+len(cen))
	ballots = append(ballots, iso...)
	ballots = append(ballots, cen...)

	voterFilter, err := filereader.NewFilter("commitment_role==Voter&committee_domain!=national&commitment_status==active")

	voters, err := filereader.LoadNationalEngagements(rolesPath, voterFilter)
	if err != nil {
		return fmt.Errorf("failed getting voter roles: %w", err)
	}

	log.Printf("retreived %d voter roles", len(voters))

	matched, unmatched := matchBallotVoter(ballots, voters)

	fileName := fmt.Sprintf("utestående_avstemninger-%s.xlsx", time.Now().Format("2006-01-02"))
	outPath := filepath.Join(cfg.OutputPath, fileName)
	if err := writeBallotsMatchedXLSX(outPath, matched, cfg.CentralizedVoters); err != nil {
		return err
	}

	missingVoterPath := filepath.Join(cfg.OutputPath, "missing.txt")
	if err = writeMissing(missingVoterPath, unmatched); err != nil {
		return err
	}

	return nil
}
