package indexer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func extractExcel(path string) ([]Chunk, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open Excel file: %w", err)
	}
	defer f.Close()

	source := filepath.Base(path)
	var chunks []Chunk

	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}

		var lines []string
		for _, row := range rows {
			var cells []string
			for _, cell := range row {
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
			continue
		}

		const rowsPerChunk = 20
		for start := 0; start < len(lines); start += rowsPerChunk {
			end := start + rowsPerChunk
			if end > len(lines) {
				end = len(lines)
			}
			chunks = append(chunks, Chunk{
				Source:   source,
				Location: fmt.Sprintf("シート: %s (行 %d–%d)", sheet, start+1, end),
				Text:     strings.Join(lines[start:end], "\n"),
			})
		}
	}

	return chunks, nil
}
