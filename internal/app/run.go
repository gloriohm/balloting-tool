package app

import (
	"ballot-tool/internal/io"
	"ballot-tool/internal/processing"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

	requiredBallots, balloutsWithoutVoter := processing.JoinBallotRole(roles, ballots)

	fileName := fmt.Sprintf("utestående_avstemninger-%s.xlsx", time.Now().Format("2006-01-02"))
	outPath := filepath.Join(cfg.OutputPath, fileName)
	if err := io.WriteBallotWithRoleXLSX(outPath, requiredBallots, cfg.CentralizedVoters); err != nil {
		return err
	}

	missingVoterPath := filepath.Join(cfg.OutputPath, "missing.txt")
	if err = io.WriteBallotsTXT(missingVoterPath, balloutsWithoutVoter); err != nil {
		return err
	}

	if opt {
		orgRoles, err := getOrgRoles(cfg)
		if err != nil {
			return err
		}
		memberWithoutVoter := processing.JoinCommitteeRole(roles, orgRoles)
		memberStatus := filepath.Join(cfg.OutputPath, "member_status.txt")
		if err = io.WriteCommitteesTXT(memberStatus, memberWithoutVoter); err != nil {
			return err
		}
	}

	return nil
}
