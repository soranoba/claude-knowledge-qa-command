package indexer

import (
	"encoding/json"
	"strings"
	"testing"
)

const csvSamplePath = "../../samples/user-data.csv"

func TestExtractCSV_ChunkCount(t *testing.T) {
	chunks, err := extractCSV(csvSamplePath)
	if err != nil {
		t.Fatalf("extractCSV: %v", err)
	}
	if len(chunks) != 100 {
		t.Errorf("got %d chunks, want 100", len(chunks))
	}
}

func TestExtractCSV_JSONFormat(t *testing.T) {
	chunks, err := extractCSV(csvSamplePath)
	if err != nil {
		t.Fatalf("extractCSV: %v", err)
	}
	for _, chunk := range chunks {
		var m map[string]string
		if err := json.Unmarshal([]byte(chunk.Text), &m); err != nil {
			t.Errorf("chunk %q: text is not valid JSON: %v", chunk.Location, err)
		}
	}
}

func TestExtractCSV_HeaderKeys(t *testing.T) {
	chunks, err := extractCSV(csvSamplePath)
	if err != nil {
		t.Fatalf("extractCSV: %v", err)
	}
	wantKeys := []string{"ID", "Name", "Age", "Country", "Email"}
	for _, chunk := range chunks {
		var m map[string]string
		if err := json.Unmarshal([]byte(chunk.Text), &m); err != nil {
			continue
		}
		for _, key := range wantKeys {
			if _, ok := m[key]; !ok {
				t.Errorf("chunk %q: missing key %q", chunk.Location, key)
			}
		}
	}
}

func TestExtractCSV_CommentLineExcluded(t *testing.T) {
	chunks, err := extractCSV(csvSamplePath)
	if err != nil {
		t.Fatalf("extractCSV: %v", err)
	}
	for _, chunk := range chunks {
		if strings.Contains(chunk.Text, "sample-files.com") {
			t.Errorf("comment line leaked into chunk: %s", chunk.Text)
		}
	}
}

func TestExtractCSV_Source(t *testing.T) {
	chunks, err := extractCSV(csvSamplePath)
	if err != nil {
		t.Fatalf("extractCSV: %v", err)
	}
	for _, chunk := range chunks {
		if chunk.Source != "user-data.csv" {
			t.Errorf("chunk source = %q, want %q", chunk.Source, "user-data.csv")
		}
	}
}
