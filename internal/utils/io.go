package utils

import (
	"log"
	"os"
	"path/filepath"
)

func LoadTabularDataFromFile(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
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
		panic("file type must be csv or xlsx")
	}

	log.Printf("loaded rows from %s, length %d\n", path, len(rows))

	return rows, nil
}
