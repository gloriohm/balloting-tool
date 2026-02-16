package app

import (
	"ballot-tool/internal/config"
	"os"
	"path/filepath"
)

func setConfig() (*config.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return nil, err
	}

	if cfg.InputPath == "" {
		cfg.InputPath = filepath.Join(home, "downloads")
	}

	if cfg.OutputPath == "" {
		cfg.OutputPath = filepath.Join(cfg.InputPath, "ballot_resultat")
	}

	if cfg.Files.Ballot1 == "" {
		cfg.Files.Ballot1 = "iso_ballots.xlsx"
	}

	if cfg.Files.Ballot2 == "" {
		cfg.Files.Ballot2 = "cen_ballots.xlsx"
	}

	if cfg.Files.Voters == "" {
		cfg.Files.Voters = "roles.csv"
	}

	if cfg.Files.OrgRoles == "" {
		cfg.Files.OrgRoles = "org_roles.xlsx"
	}

	return cfg, nil
}
