package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func InitConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg, err := loadConfig("config.json")
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

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
