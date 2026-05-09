package indexer

import (
	"fmt"
	"path/filepath"

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
		if err != nil || len(rows) < 2 {
			continue
		}
		header := rows[0]
		chunks = append(chunks, chunksFromTable(source, fmt.Sprintf("sheet: %s, row", sheet), header, rows[1:])...)
	}

	return chunks, nil
}
