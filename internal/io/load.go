package io

import (
	"log"
	"os"
	"path/filepath"
)

func LoadFile(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := validateTimestamp(f); err != nil {
		return nil, err
	}

	fileType := filepath.Ext(path)

	var rows []map[string]string
	switch fileType {
	case ".csv":
		rows, err = ReadCSV(f)
	case ".xlsx":
		rows, err = ReadXLSX(f, "")
	default:
		panic("unsupported file type")
	}

	log.Printf("loaded rows from %s, length %d\n", path, len(rows))

	return rows, nil
}
