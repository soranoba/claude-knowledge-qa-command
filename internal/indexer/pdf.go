package indexer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/razvandimescu/gopdf/pdf"
)

func extractPDF(path string) ([]Chunk, error) {
	doc, err := pdf.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open PDF: %w", err)
	}

	source := filepath.Base(path)
	var chunks []Chunk

	for i := 0; i < doc.NumPages(); i++ {
		text, err := doc.Page(i).Text()
		if err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		for _, part := range splitIntoChunks(text, 600) {
			chunks = append(chunks, Chunk{
				Source:   source,
				Location: fmt.Sprintf("ページ %d", i+1),
				Text:     part,
			})
		}
	}

	return chunks, nil
}

func splitIntoChunks(text string, maxRunes int) []string {
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return []string{text}
	}

	var chunks []string
	for len(runes) > 0 {
		end := maxRunes
		if end > len(runes) {
			end = len(runes)
		}
		if end < len(runes) {
			lookback := 100
			if lookback > end {
				lookback = end
			}
			for i := end; i > end-lookback; i-- {
				r := runes[i]
				if r == '。' || r == '\n' || r == '.' || r == '！' || r == '？' {
					end = i + 1
					break
				}
			}
		}
		chunks = append(chunks, string(runes[:end]))
		runes = runes[end:]
	}
	return chunks
}
