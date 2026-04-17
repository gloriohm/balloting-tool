package ballot

import (
	"ballot-tool/internal/utils"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadCSV(r io.Reader) ([]map[string]string, error) {
	delimiter, br, err := utils.DetectDelimiter(r)
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
		if canon, ok := headerAliases[nh]; ok {
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

func ReadXLSX(r io.Reader, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

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

	rawHeader := allRows[0]
	headers := make([]string, len(rawHeader))
	for i, h := range rawHeader {
		nh := utils.NormalizeString(h)
		if canon, ok := headerAliases[nh]; ok {
			headers[i] = canon
		} else {
			headers[i] = nh
		}
	}

	rows := make([]map[string]string, 0, max(0, len(allRows)-1))
	for _, rec := range allRows[1:] {
		empty := true
		for _, c := range rec {
			if strings.TrimSpace(c) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}

		m := make(map[string]string, len(headers))
		for i := 0; i < len(rec) && i < len(headers); i++ {
			m[headers[i]] = strings.TrimSpace(rec[i])
		}
		for i := len(rec); i < len(headers); i++ {
			m[headers[i]] = ""
		}
		rows = append(rows, m)
	}

	return rows, nil
}
