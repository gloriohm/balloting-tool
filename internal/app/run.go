package app

import (
	"ballot-tool/internal/api/sdimport"
	"ballot-tool/internal/tools/ballot"
	"ballot-tool/internal/tools/standards"
	"ballot-tool/internal/utils/config"
	"fmt"
	"log"
	"os"
)

func RunBallotTool() error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.OutputPath, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := ballot.GenerateBallotReport(cfg); err != nil {
		return fmt.Errorf("failed to genereate ballot report: %w", err)
	}

	return nil
}

func RunStandardsTool(job, from, to, filename, opts string, nsOnly, dev bool) error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}
	params := sdimport.NewParameters(from, to)
	client := sdimport.NewClient(dev, params)
	stdSvc := standards.NewService(client, cfg)

	switch job {
	case "count":
		if err := stdSvc.CountTotalUniqueProducts(filename, job, nsOnly); err != nil {
			return err
		}
	case "fetch":
		if err := stdSvc.GetStandards(); err != nil {
			return err
		}
	case "aktualitet":
		if err := stdSvc.GenerateAktualitetList(filename); err != nil {
			return err
		}
	case "xml":
		if err := stdSvc.FindStandardsWithXML(filename); err != nil {
			return err
		}
	case "download":
		if err := stdSvc.DownloadFiles(filename, opts); err != nil {
			return err
		}
	default:
		log.Println("unknown job")
	}

	return nil
}
