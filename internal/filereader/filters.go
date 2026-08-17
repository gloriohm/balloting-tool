package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"strings"
	"time"
)

type Filter interface {
	Match(Row) (bool, error)
}

type Filters []Filter

type InclusiveStringFilter struct {
	Field  string
	Values map[string]struct{}
}

type ExclusiveStringFilter struct {
	Field  string
	Values map[string]struct{}
}

type BeginsWithFilter struct {
	Field  string
	Values map[string]struct{}
}

type GreaterOrEqualTime struct {
	Field string
	Value time.Time
}

type LessOrEqualTime struct {
	Field string
	Value time.Time
}

type HasPrefixFilter struct {
	Field  string
	Values map[string]struct{}
}

type NotPrefixFilter struct {
	Field  string
	Values map[string]struct{}
}

type HasSuffixFilter struct {
	Field  string
	Values map[string]struct{}
}

type NotSuffixFilter struct {
	Field  string
	Values map[string]struct{}
}

type EngagementMinRangeFilter struct {
}

func (f EngagementMinRangeFilter) Match(row Row) (bool, error) {
	from, err := time.Parse("2006-01-02", row["commitment_from"])
	if err != nil {
		return false, nil
	}

	to, err := time.Parse("2006-01-02", row["commitment_to"])
	if err != nil {
		return true, nil
	}

	fromPlusOneWeek := from.AddDate(0, 0, 7)

	if !to.After(fromPlusOneWeek) {
		return false, nil
	}

	return true, nil
}

func (f InclusiveStringFilter) Match(row Row) (bool, error) {
	_, ok := f.Values[row[f.Field]]
	return ok, nil
}

func (f ExclusiveStringFilter) Match(row Row) (bool, error) {
	_, ok := f.Values[row[f.Field]]
	return !ok, nil
}

func (f GreaterOrEqualTime) Match(row Row) (bool, error) {
	v, err := time.Parse("2006-01-02", row[f.Field])
	if err != nil {
		return false, nil
	}

	return !v.Before(f.Value), nil
}

func (f LessOrEqualTime) Match(row Row) (bool, error) {
	v, err := time.Parse("2006-01-02", row[f.Field])
	if err != nil {
		return false, nil
	}

	return !v.After(f.Value), nil
}

func (f HasPrefixFilter) Match(row Row) (bool, error) {
	have := normalization.NormalizeString(row[f.Field])
	for v := range f.Values {
		if strings.HasPrefix(have, normalization.NormalizeString(v)) {
			return true, nil
		}
	}

	return false, nil
}

func (f NotPrefixFilter) Match(row Row) (bool, error) {
	have := normalization.NormalizeString(row[f.Field])
	for v := range f.Values {
		if strings.HasPrefix(have, normalization.NormalizeString(v)) {
			return false, nil
		}
	}

	return true, nil
}

func (f HasSuffixFilter) Match(row Row) (bool, error) {
	have := normalization.NormalizeString(row[f.Field])
	for v := range f.Values {
		if strings.HasSuffix(have, normalization.NormalizeString(v)) {
			return true, nil
		}
	}

	return false, nil
}

func (f NotSuffixFilter) Match(row Row) (bool, error) {
	have := normalization.NormalizeString(row[f.Field])
	for v := range f.Values {
		if strings.HasSuffix(have, normalization.NormalizeString(v)) {
			return false, nil
		}
	}

	return true, nil
}

func (fs Filters) Match(row Row) (bool, error) {
	for _, f := range fs {
		ok, err := f.Match(row)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}
