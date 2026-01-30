package io

import (
	"os"
	"path/filepath"
)

func LoadFile(path string) []map[string]string {
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
		panic("unsupported file type")
	}

	return rows
}
