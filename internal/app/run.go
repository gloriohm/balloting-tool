package app

import (
	"ballot-tool/internal/ballot"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	url          = "https://isotc.iso.org/livelink/eb3/part/viewMyBallots.do?method=doVoteRequired&org.apache.struts.taglib.html.CANCEL=true&startIndex=0"
	cenBallotURL = "https://cen.iso.org/livelink/eb33/part/exportBallotListXLS.do?noreset=true"
	isoBallotURL = "https://isotc.iso.org/livelink/eb3/part/exportBallotListXLS.do?noreset=true"
)

func Run(opt bool) error {
	cfg, err := setConfig()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.OutputPath, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	ballots, err := getBallots(cfg)
	if err != nil {
		return err
	}

	roles, err := getVoterRoles(cfg)
	if err != nil {
		return err
	}

	requiredBallots, balloutsWithoutVoter := ballot.JoinBallotRole(roles, ballots)

	fileName := fmt.Sprintf("utestående_avstemninger-%s.xlsx", time.Now().Format("2006-01-02"))
	outPath := filepath.Join(cfg.OutputPath, fileName)
	if err := ballot.WriteBallotWithRoleXLSX(outPath, requiredBallots, cfg.CentralizedVoters); err != nil {
		return err
	}

	missingVoterPath := filepath.Join(cfg.OutputPath, "missing.txt")
	if err = ballot.WriteBallotsTXT(missingVoterPath, balloutsWithoutVoter); err != nil {
		return err
	}

	if opt {
		orgRoles, err := getOrgRoles(cfg)
		if err != nil {
			return err
		}
		memberWithoutVoter := ballot.JoinCommitteeRole(roles, orgRoles)
		memberStatus := filepath.Join(cfg.OutputPath, "member_status.txt")
		if err = ballot.WriteCommitteesTXT(memberStatus, memberWithoutVoter); err != nil {
			return err
		}
	}

	return nil
}
