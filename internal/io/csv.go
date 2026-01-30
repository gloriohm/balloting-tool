package io

import (
	"ballot-tool/internal/utils"
	"encoding/csv"
	"io"
	"strings"
)

func ReadCSV(r io.Reader) ([]map[string]string, error) {
	delimiter, br, err := detectDelimiter(r)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(br)
	cr.Comma = delimiter
	cr.TrimLeadingSpace = true
	rawHeader, err := cr.Read()
	if err != nil {
		return nil, err
	}
	headers := make([]string, len(rawHeader))

	for i, h := range rawHeader {
		nh := utils.NormalizeString(h)
		if canon, ok := utils.HeaderAliases[nh]; ok {
			headers[i] = canon
		} else {
			headers[i] = nh
		}
	}

	var rows []map[string]string
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		m := make(map[string]string)
		for i, cell := range rec {
			m[headers[i]] = strings.TrimSpace(cell)
		}
		rows = append(rows, m)
	}

	return rows, nil
}
