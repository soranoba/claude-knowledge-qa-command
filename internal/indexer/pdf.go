package indexer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/razvandimescu/gopdf/pdf"
)

func extractPDF(path string) ([]Chunk, error) {
	doc, err := pdf.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open PDF: %w", err)
	}

	source := filepath.Base(path)

	type lineWithPage struct {
		line pdf.TextLine
		page int
	}

	var allLines []lineWithPage
	var allFontSizes []float64

	for i := 0; i < doc.NumPages(); i++ {
		lines, err := doc.Page(i).TextLines()
		if err != nil {
			continue
		}
		for _, line := range lines {
			allLines = append(allLines, lineWithPage{line, i + 1})
			for _, span := range line.Spans {
				if span.FontSize > 0 {
					allFontSizes = append(allFontSizes, span.FontSize)
				}
			}
		}
	}

	if len(allLines) == 0 {
		return nil, nil
	}

	headingThreshold := estimateBodyFontSize(allFontSizes) * 1.3

	type section struct {
		location string
		lines    []string
	}

	current := section{location: fmt.Sprintf("p%d", allLines[0].page)}
	var sections []section

	for _, lp := range allLines {
		text := strings.TrimSpace(lp.line.Text)
		if text == "" {
			continue
		}

		if maxFontSize(lp.line) >= headingThreshold && len([]rune(text)) <= 100 {
			if len(current.lines) > 0 {
				sections = append(sections, current)
			}
			current = section{location: text}
		} else {
			current.lines = append(current.lines, text)
		}
	}
	if len(current.lines) > 0 {
		sections = append(sections, current)
	}

	var chunks []Chunk
	for _, sec := range sections {
		text := strings.Join(sec.lines, "\n")
		for _, part := range splitIntoChunks(text, 600) {
			chunks = append(chunks, Chunk{
				Source:   source,
				Location: sec.location,
				Text:     part,
			})
		}
	}

	return chunks, nil
}

func maxFontSize(line pdf.TextLine) float64 {
	maxFontSize := 0.0
	for _, span := range line.Spans {
		if span.FontSize > maxFontSize {
			maxFontSize = span.FontSize
		}
	}
	return maxFontSize
}

func estimateBodyFontSize(values []float64) float64 {
	if len(values) == 0 {
		return 12
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	return sorted[len(sorted)/2]
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
				if unicode.Is(unicode.STerm, r) || r == '\n' {
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
