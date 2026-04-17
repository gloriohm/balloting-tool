package app

import (
	"ballot-tool/internal/ballot"
	"ballot-tool/internal/config"
	"log"
	"path/filepath"
)

func getBallots(cfg *config.Config) ([]ballot.Ballot, error) {
	isoBallotRows, err := ballot.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Ballot1))
	if err != nil {
		return nil, err
	}
	cenBallotRows, err := ballot.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Ballot2))
	if err != nil {
		return nil, err
	}

	isoBallotsParsed, isoBallotErrs := ballot.ParseBallots(isoBallotRows)
	cenBallotsParsed, cenBallotErrs := ballot.ParseBallots(cenBallotRows)
	if err := ballot.CombineErrs(
		"iso ballots", isoBallotErrs,
		"cen ballots", cenBallotErrs,
	); err != nil {
		log.Printf("parsing errors: %s", err)
	}

	allBallots := make([]ballot.Ballot, 0, len(isoBallotsParsed)+len(cenBallotsParsed))
	allBallots = append(allBallots, isoBallotsParsed...)
	allBallots = append(allBallots, cenBallotsParsed...)

	return allBallots, nil
}

func getVoterRoles(cfg *config.Config) ([]ballot.Role, error) {
	roleRows, err := ballot.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.Voters))
	if err != nil {
		return nil, err
	}

	rolesParsed, roleErrs := ballot.ParseRoles(roleRows)
	ballot.ErrorPrinter(roleErrs, 5)

	return rolesParsed, nil
}

func getOrgRoles(cfg *config.Config) ([]ballot.Committee, error) {
	orgRows, err := ballot.LoadFile(filepath.Join(cfg.InputPath, cfg.Files.OrgRoles))
	if err != nil {
		return nil, err
	}

	orgParsed, orgErrs := ballot.ParseCommittees(orgRows)
	ballot.ErrorPrinter(orgErrs, 5)

	return orgParsed, nil
}
