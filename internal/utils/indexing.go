package utils

import "fmt"

func IndexStrings(s []string, offset int) (map[string]int, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}
	idx := make(map[string]int, len(s))
	for i, st := range s {
		ns := ToLowerCase(st)
		if ns == "" {
			continue
		}

		if _, exists := idx[ns]; !exists {
			idx[ns] = i + offset
		}
	}
	return idx, nil
}
