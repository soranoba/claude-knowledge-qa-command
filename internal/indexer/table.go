package indexer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chunksFromTable converts tabular data (header + rows) into one Chunk per row.
// Each chunk's text is a JSON object keyed by header values.
// locationPrefix is prepended to the row number (e.g. "sheet: Sheet1, row").
func chunksFromTable(source, locationPrefix string, header []string, rows [][]string) []Chunk {
	var chunks []Chunk
	rowNum := 1 // 1-indexed data rows (row 1 = first data row after header)
	for _, record := range rows {
		row := make(map[string]string)
		for i, cell := range record {
			cell = strings.TrimSpace(cell)
			if cell == "" {
				continue
			}
			key := fmt.Sprintf("col%d", i+1)
			if i < len(header) && header[i] != "" {
				key = header[i]
			}
			row[key] = cell
		}
		if len(row) == 0 {
			rowNum++
			continue
		}
		b, _ := json.Marshal(row)
		chunks = append(chunks, Chunk{
			Source:   source,
			Location: fmt.Sprintf("%s %d", locationPrefix, rowNum+1), // +1 to account for header row
			Text:     string(b),
		})
		rowNum++
	}
	return chunks
}
