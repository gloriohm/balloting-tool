package app

import (
	"ballot-tool/internal/config"
	"ballot-tool/internal/io"
	"ballot-tool/internal/models"
	"ballot-tool/internal/parse"
	"log"
	"path/filepath"
)

func getBallots(cfg *config.Config) ([]models.Ballot, error) {
	isoBallotRows, err := io.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Ballot1))
	if err != nil {
		return nil, err
	}
	cenBallotRows, err := io.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Ballot2))
	if err != nil {
		return nil, err
	}

	isoBallotsParsed, isoBallotErrs := parse.ParseBallots(isoBallotRows)
	cenBallotsParsed, cenBallotErrs := parse.ParseBallots(cenBallotRows)
	if err := parse.CombineErrs(
		"iso ballots", isoBallotErrs,
		"cen ballots", cenBallotErrs,
	); err != nil {
		log.Printf("parsing errors: %s", err)
	}

	allBallots := make([]models.Ballot, 0, len(isoBallotsParsed)+len(cenBallotsParsed))
	allBallots = append(allBallots, isoBallotsParsed...)
	allBallots = append(allBallots, cenBallotsParsed...)

	return allBallots, nil
}

func getVoterRoles(cfg *config.Config) ([]models.Role, error) {
	roleRows, err := io.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Voters))
	if err != nil {
		return nil, err
	}

	rolesParsed, roleErrs := parse.ParseRoles(roleRows)
	parse.ErrorPrinter(roleErrs, 5)

	return rolesParsed, nil
}

func getOrgRoles(cfg *config.Config) ([]models.Committee, error) {
	orgRows, err := io.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.OrgRoles))
	if err != nil {
		return nil, err
	}

	orgParsed, orgErrs := parse.ParseCommittees(orgRows)
	parse.ErrorPrinter(orgErrs, 5)

	return orgParsed, nil
}
