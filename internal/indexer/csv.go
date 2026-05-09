package indexer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func extractCSV(path string) ([]Chunk, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open CSV file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	r.Comment = '#'

	// Read header row
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			return nil, nil
		}
		if err != nil {
			continue
		}
		header = record
		break
	}

	var rows [][]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		rows = append(rows, record)
	}

	return chunksFromTable(filepath.Base(path), "row", header, rows), nil
}
