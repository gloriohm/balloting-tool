package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"ballot-tool/internal/utils/parsing"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func LoadNationalEngagements(path string, filter Filters) ([]NationalEngagementsRow, error) {
	return parseCSV(path, filter, parseNationalEngagementsRow)
}

func LoadBallots(path string, filter Filters) ([]BallotRow, error) {
	filetype := filepath.Ext(path)

	switch filetype {
	case ".csv":
		return parseCSV(path, filter, parseBallotRow)
	case ".xlsx":
		return parseExcel(path, filter, parseBallotRow)
	default:
		return nil, ErrInvalidFileType
	}
}

func LoadStandardsDashboard(path string, filters Filters) ([]StandardDashboardRow, error) {
	return parseCSV(path, filters, parseStandardDashboardRow)
}

func NewFilter(s string) (Filters, error) {
	// filterstring should follow format key (identical to normalized column name)
	// operator == for inclusive or != for exclusive
	// phrases separated by ;
	// parameters separated by &
	// example: commitment_status==active&committee_status!=in_progress;suspended

	if strings.TrimSpace(s) == "" {
		return Filters{}, nil
	}

	claims := strings.Split(s, "&")
	f := make(Filters, len(claims))

	for _, claim := range claims {
		switch {
		case strings.Contains(claim, "=="):
			keyValues := strings.SplitN(claim, "==", 2)
			key := keyValues[0]
			values := splitBySeparator(keyValues[1], ";")
			set := parsing.ToSet(values)
			f = append(f, InclusiveStringFilter{Field: key, Values: set})

			log.Printf("key: %s; values: %s", key, values)
		case strings.Contains(claim, "!="):
			keyValues := strings.SplitN(claim, "!=", 2)
			key := keyValues[0]
			values := splitBySeparator(keyValues[1], ";")
			set := parsing.ToSet(values)
			f = append(f, ExclusiveStringFilter{Field: key, Values: set})

			log.Printf("key: %s; values: %s", key, values)
		default:
			return nil, ErrUnknownOperator
		}
	}

	return f, nil
}

func NewInclusiveStringFilter(column string, targets map[string]struct{}) InclusiveStringFilter {
	return InclusiveStringFilter{Field: column, Values: targets}
}

func NewExclusiveStringFilter(column string, targets map[string]struct{}) ExclusiveStringFilter {
	return ExclusiveStringFilter{Field: column, Values: targets}
}

func NewGreaterThanFilter(column string, target time.Time) GreaterOrEqualTime {
	return GreaterOrEqualTime{Field: column, Value: target}
}

func NewLessThanFilter(column string, target time.Time) LessOrEqualTime {
	return LessOrEqualTime{Field: column, Value: target}
}

func (f *Filters) NewBeginsWith(column string, targets map[string]struct{}) {
	column = normalization.NormalizeString(column)
	*f = append(*f, HasPrefixFilter{Field: column, Values: targets})
}

func (f *Filters) NewNotSuffixFilter(column string, targets map[string]struct{}) {
	column = normalization.NormalizeString(column)
	*f = append(*f, NotSuffixFilter{Field: column, Values: targets})
}

func (f *Filters) NewEngagementsMinRangeFilter() {
	*f = append(*f, EngagementMinRangeFilter{})
}

func NewProjectsFilter() Filters {
	f := make(Filters, 1)
	f = append(f, InclusiveStringFilter{Field: "stage", Values: map[string]struct{}{"working": {}}})

	return f
}
