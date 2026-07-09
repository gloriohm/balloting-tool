package filereader

import "slices"

type Filters map[string]Filter

type Filter struct {
	Targets []string
	Func    FilterFunc
}

type FilterFunc func(have string, want []string) bool

func inclusiveFilter(have string, want []string) bool {
	if slices.Contains(want, have) {
		return true
	}

	return false
}

func exclusiveFilter(have string, want []string) bool {
	if slices.Contains(want, have) {
		return false
	}

	return true
}

func passesFilters(row Row, filters map[string]Filter) bool {
	for col, filter := range filters {
		if !filter.Func(row[col], filter.Targets) {
			return false
		}
	}

	return true
}
