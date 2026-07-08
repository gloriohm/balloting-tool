package committee

import (
	"ballot-tool/internal/api/brreg"
	"ballot-tool/internal/utils/config"
)

type Service struct {
	brreg *brreg.Client
	cfg   *config.Config
}

func NewService(brreg *brreg.Client, cfg *config.Config) *Service {
	return &Service{brreg: brreg, cfg: cfg}
}
