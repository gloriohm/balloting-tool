package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"log"
	"path/filepath"
	"strings"
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
			f[key] = Filter{Func: inclusiveFilter, Targets: values}

			log.Printf("key: %s; values: %s", key, values)
		case strings.Contains(claim, "!="):
			keyValues := strings.SplitN(claim, "!=", 2)
			key := keyValues[0]
			values := splitBySeparator(keyValues[1], ";")
			f[key] = Filter{Func: exclusiveFilter, Targets: values}

			log.Printf("key: %s; values: %s", key, values)
		default:
			return nil, ErrUnknownOperator
		}
	}

	return f, nil
}

func (f Filters) NewBeginsWith(column string, targets []string, inclusive bool) {
	column = normalization.NormalizeString(column)

	fn := exclusiveHasPrefixFilter
	if inclusive {
		fn = inclusiveHasPrefixFilter
	}

	f[column] = Filter{
		Targets: targets,
		Func:    fn,
	}
}

func NewProjectsFilter() Filters {
	f := make(Filters, 1)
	f["stage"] = Filter{Targets: []string{"working"}, Func: inclusiveFilter}

	return f
}
