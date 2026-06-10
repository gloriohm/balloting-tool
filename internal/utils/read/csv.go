package read

import (
	"ballot-tool/internal/utils/normalization"
	"encoding/csv"
	"io"
	"strings"
)

func ReadCSV(r io.Reader) ([]map[string]string, error) {
	delimiter, br, err := DetectDelimiter(r)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(br)
	cr.Comma = delimiter
	cr.TrimLeadingSpace = true
	rawHeaders, err := cr.Read()
	if err != nil {
		return nil, err
	}

	headers := make([]string, len(rawHeaders))

	for i, h := range rawHeaders {
		headers[i] = normalization.NormalizeString(h)
	}

	var rows []map[string]string
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		m := make(map[string]string, len(headers))

		for i, header := range headers {
			if i < len(row) {
				m[header] = strings.TrimSpace(row[i])
			} else {
				m[header] = ""
			}
		}

		rows = append(rows, m)
	}

	return rows, nil
}
