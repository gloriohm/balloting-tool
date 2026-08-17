package committee

import (
	"ballot-tool/internal/api/brreg"
	"ballot-tool/internal/filereader"
	"ballot-tool/internal/utils/config"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	brreg *brreg.Client
	cfg   *config.Config
}

func NewService(brreg *brreg.Client, cfg *config.Config) *Service {
	return &Service{brreg: brreg, cfg: cfg}
}

func (s *Service) CountExpertEngagements(path, start, end string) error {
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return errors.New("start time incorrect format")
	}
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return errors.New("end time incorrect format")
	}

	filter := filereader.Filters{}
	filter = append(filter, filereader.NewGreaterThanFilter("commitment_from", startTime))
	filter = append(filter, filereader.NewLessThanFilter("commitment_from", endTime))
	filter = append(filter, filereader.NewInclusiveStringFilter("committee_domain", map[string]struct{}{"regional": {}, "international": {}}))
	//filter = append(filter, filereader.NewInclusiveStringFilter("committee_level", map[string]struct{}{"wg": {}}))

	filter.NewNotSuffixFilter("email", map[string]struct{}{"@standard.no": {}})
	filter.NewEngagementsMinRangeFilter()

	experts, err := filereader.LoadNationalEngagements(path, filter)
	if err != nil {
		return err
	}

	fmt.Printf("Total expert engagements %d\n", len(experts))

	return WriteEngagementsExcel(s.cfg.OutputPath, parseEngagements(experts))
}
