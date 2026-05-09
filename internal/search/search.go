package search

import (
	"sort"
	"strings"

	"github.com/soranoba/knowledge-qa/internal/indexer"
)

type Result struct {
	Source   string `json:"source"`
	Location string `json:"location"`
	Text     string `json:"text"`
}

func Query(chunks []indexer.Chunk, question string, topK int) []Result {
	if len(chunks) == 0 || strings.TrimSpace(question) == "" {
		return []Result{}
	}

	type scored struct {
		chunk indexer.Chunk
		score float64
	}

	queryBigrams := charBigrams(question)
	queryTerms := splitTerms(question)

	var results []scored
	for _, chunk := range chunks {
		s := score(chunk.Text, queryBigrams, queryTerms)
		if s > 0 {
			results = append(results, scored{chunk, s})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}

	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{
			Source:   r.chunk.Source,
			Location: r.chunk.Location,
			Text:     r.chunk.Text,
		}
	}
	return out
}

func score(text string, queryBigrams map[string]int, queryTerms []string) float64 {
	if len(queryBigrams) == 0 {
		return 0
	}

	textBigrams := charBigrams(text)
	s := 0.0
	for bg, qc := range queryBigrams {
		if tc, ok := textBigrams[bg]; ok {
			match := qc
			if tc < match {
				match = tc
			}
			s += float64(match)
		}
	}

	// Bonus for longer exact substring matches
	lowerText := strings.ToLower(text)
	for _, term := range queryTerms {
		if len([]rune(term)) >= 2 && strings.Contains(lowerText, strings.ToLower(term)) {
			s += float64(len([]rune(term))) * 2
		}
	}

	// Normalize to avoid favoring very long chunks
	textLen := float64(len([]rune(text)))
	if textLen > 0 {
		s = s / textLen * 1000
	}

	return s
}

func charBigrams(text string) map[string]int {
	runes := []rune(strings.ToLower(text))
	counts := make(map[string]int)
	for i := 0; i+1 < len(runes); i++ {
		bg := string(runes[i : i+2])
		counts[bg]++
	}
	return counts
}

// splitTerms splits by spaces and common Japanese particles to extract meaningful terms.
func splitTerms(text string) []string {
	parts := strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case ' ', '\t', '　', 'を', 'は', 'が', 'に', 'で', 'と', 'の', 'も', 'や', 'て':
			return true
		}
		return false
	})
	var terms []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len([]rune(p)) >= 2 {
			terms = append(terms, p)
		}
	}
	return terms
}
