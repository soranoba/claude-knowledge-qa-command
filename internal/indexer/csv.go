package indexer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	source := filepath.Base(path)
	var lines []string
	rowNum := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		rowNum++

		var cells []string
		for _, cell := range record {
			cell = strings.TrimSpace(cell)
			if cell != "" {
				cells = append(cells, cell)
			}
		}
		if len(cells) > 0 {
			lines = append(lines, strings.Join(cells, " | "))
		}
	}

	if len(lines) == 0 {
		return nil, nil
	}

	var chunks []Chunk
	const rowsPerChunk = 20
	for start := 0; start < len(lines); start += rowsPerChunk {
		end := start + rowsPerChunk
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, Chunk{
			Source:   source,
			Location: fmt.Sprintf("行 %d–%d", start+1, end),
			Text:     strings.Join(lines[start:end], "\n"),
		})
	}

	return chunks, nil
}
