package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func LoadTabularDataFromFile(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file with path %s: %w", path, err)
	}
	defer f.Close()

	fileType := filepath.Ext(path)

	var rows []map[string]string
	switch fileType {
	case ".csv":
		rows, err = ReadCSV(f)
	case ".xlsx":
		rows, err = ReadXLSX(f, "")
	default:
		return nil, fmt.Errorf("file type must be csv or xlsx")
	}

	log.Printf("loaded rows from %s, length %d\n", path, len(rows))

	return rows, nil
}

func LoadTabularDataFromFolder(folder string) ([]map[string]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var rows []map[string]string

	for _, entry := range entries {
		// Skip subdirectories
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(folder, entry.Name())

		data, err := LoadTabularDataFromFile(path)
		if err != nil {
			log.Printf("[SKIPPED] could not open file with path %s: %s", path, err)
			continue
		}

		rows = append(rows, data...)
	}

	return rows, nil
}
