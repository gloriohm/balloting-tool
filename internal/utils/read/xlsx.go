package read

import (
	"ballot-tool/internal/utils/normalization"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadXLSX(r io.Reader, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("xlsx: no sheets found")
		}
		sheetName = sheets[0]
	}

	allRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(allRows) == 0 {
		return []map[string]string{}, nil
	}

	rawHeaders := allRows[0]
	headers := make([]string, len(rawHeaders))

	for i, h := range rawHeaders {
		headers[i] = normalization.NormalizeString(h)
	}

	rows := make([]map[string]string, 0, max(0, len(allRows)-1))
	for _, row := range allRows[1:] {
		empty := true
		for _, c := range row {
			if strings.TrimSpace(c) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
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
