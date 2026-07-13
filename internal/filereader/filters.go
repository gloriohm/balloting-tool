package filereader

import (
	"ballot-tool/internal/utils/normalization"
	"strings"
)

type Filters map[string]Filter

type Filter struct {
	Targets []string
	Func    FilterFunc
}

type FilterFunc func(have string, want []string) bool

func inclusiveFilter(have string, want []string) bool {
	have = normalization.NormalizeString(have)
	for _, target := range want {
		if normalization.NormalizeString(target) == have {
			return true
		}
	}

	return false
}

func inclusiveHasPrefixFilter(have string, want []string) bool {
	have = normalization.NormalizeString(have)
	for _, target := range want {
		if strings.HasPrefix(have, normalization.NormalizeString(target)) {
			return true
		}
	}

	return false
}

func exclusiveFilter(have string, want []string) bool {
	have = normalization.NormalizeString(have)
	for _, target := range want {
		if normalization.NormalizeString(target) == have {
			return false
		}
	}

	return true
}

func exclusiveHasPrefixFilter(have string, want []string) bool {
	have = normalization.NormalizeString(have)
	for _, target := range want {
		if strings.HasPrefix(have, normalization.NormalizeString(target)) {
			return false
		}
	}

	return true
}

func passesFilters(row Row, filters Filters) bool {
	for col, filter := range filters {
		value, exists := row[col]
		if !exists {
			return false
		}

		if filter.Func == nil {
			return false
		}

		if !filter.Func(value, filter.Targets) {
			return false
		}
	}

	return true
}
