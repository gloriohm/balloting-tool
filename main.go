package main

import (
	"ballot-tool/internal/config"
	"ballot-tool/internal/io"
	"ballot-tool/internal/models"
	"ballot-tool/internal/parse"
	"ballot-tool/internal/processing"
	"log"
	"path/filepath"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	dir := cfg.OutputPath
	centralizedVoters := cfg.CentralizedVoters

	cenBallotRows := io.LoadFile(filepath.Join(dir, "cen_ballots.csv"))
	isoBallotRows := io.LoadFile(filepath.Join(dir, "iso_ballots.xlsx"))
	roleRows := io.LoadFile(filepath.Join(dir, "roles.csv"))
	orgRows := io.LoadFile(filepath.Join(dir, "org_roles.xlsx"))

	cenBallotsParsed, cenBallotErrs := parse.ParseBallots(cenBallotRows)
	isoBallotsParsed, isoBallotErrs := parse.ParseBallots(isoBallotRows)
	rolesParsed, roleErrs := parse.ParseRoles(roleRows)
	orgParsed, orgErrs := parse.ParseCommittees(orgRows)

	if err := parse.CombineErrs(
		"cen ballots", cenBallotErrs,
		"iso ballots", isoBallotErrs,
		"roles", roleErrs,
		"org", orgErrs,
	); err != nil {
		log.Printf("parsing errors: %s", err)
	}

	cenJoined, cenMissing := processing.JoinBallotRole(rolesParsed, cenBallotsParsed)
	isoJoined, isoMissing := processing.JoinBallotRole(rolesParsed, isoBallotsParsed)
	memberWithoutVoter := processing.JoinCommitteeRole(rolesParsed, orgParsed)

	allJoined := make([]models.BallotWithRole, 0, len(cenJoined)+len(isoJoined))
	allJoined = append(allJoined, cenJoined...)
	allJoined = append(allJoined, isoJoined...)

	allMissing := make([]models.Ballot, 0, len(cenMissing)+len(isoMissing))
	allMissing = append(allMissing, cenMissing...)
	allMissing = append(allMissing, isoMissing...)

	outPath := filepath.Join(dir, "utestående_avstemninger.xlsx")
	missingPath := filepath.Join(dir, "missing.txt")
	memberStatus := filepath.Join(dir, "member_status.txt")
	io.WriteBallotWithRoleXLSX(outPath, allJoined, centralizedVoters)
	io.WriteBallotsTXT(missingPath, allMissing)
	io.WriteCommitteesTXT(memberStatus, memberWithoutVoter)
}
