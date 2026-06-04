package app

import (
	"ballot-tool/internal/ballot"
	"ballot-tool/internal/config"
	"ballot-tool/internal/sdimport"
	"ballot-tool/internal/standards"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RunBallotTool(opt bool) error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.OutputPath, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	isoPath := filepath.Join(cfg.InputPath, cfg.Files.Ballot1)
	cenPath := filepath.Join(cfg.InputPath, cfg.Files.Ballot2)
	ballotPaths := []string{isoPath, cenPath}
	ballots, err := ballot.GetCombinedBallots(ballotPaths)
	if err != nil {
		return err
	}

	rolesPath := filepath.Join(cfg.InputPath, cfg.Files.Voters)
	roles, err := ballot.GetVoterRoles(rolesPath)
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
		orgPath := filepath.Join(cfg.InputPath, cfg.Files.OrgRoles)
		orgRoles, err := ballot.GetOrgRoles(orgPath)
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

func RunStandardsTool(job string, nsOnly, aktualitet bool) error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	if aktualitet {
		if err := standards.GenerateAktualitetList(cfg.InputPath); err != nil {
			return err
		}
	} else {
		if err := standards.CountTotalUniqueProducts(cfg.InputPath, job, nsOnly); err != nil {
			return err
		}
	}

	return nil
}

func RunImportTool(from, to string, dev bool) error {
	params := sdimport.NewParameters(from, to)
	client := sdimport.NewClient(dev, params)

	err := client.GetStandards()
	if err != nil {
		return err
	}

	return nil
}
