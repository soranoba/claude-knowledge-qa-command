package indexer

import (
	"testing"
)

const pdfSamplePath = "../../samples/sample-pdf-a4-size.pdf"

func TestExtractPDF_ChunksGenerated(t *testing.T) {
	chunks, err := extractPDF(pdfSamplePath)
	if err != nil {
		t.Fatalf("extractPDF: %v", err)
	}
	if len(chunks) == 0 {
		t.Fatal("got 0 chunks, want at least 1")
	}
}

func TestExtractPDF_SectionDetected(t *testing.T) {
	chunks, err := extractPDF(pdfSamplePath)
	if err != nil {
		t.Fatalf("extractPDF: %v", err)
	}

	wantSections := []string{"Introduction", "Objectives", "Conclusion"}
	found := make(map[string]bool)
	for _, chunk := range chunks {
		for _, section := range wantSections {
			if chunk.Location == section {
				found[section] = true
			}
		}
	}
	for _, section := range wantSections {
		if !found[section] {
			t.Errorf("section %q not found in any chunk location", section)
		}
	}
}

func TestExtractPDF_TextNotEmpty(t *testing.T) {
	chunks, err := extractPDF(pdfSamplePath)
	if err != nil {
		t.Fatalf("extractPDF: %v", err)
	}
	for _, chunk := range chunks {
		if chunk.Text == "" {
			t.Errorf("chunk %q has empty text", chunk.Location)
		}
	}
}

func TestExtractPDF_Source(t *testing.T) {
	chunks, err := extractPDF(pdfSamplePath)
	if err != nil {
		t.Fatalf("extractPDF: %v", err)
	}
	for _, chunk := range chunks {
		if chunk.Source != "sample-pdf-a4-size.pdf" {
			t.Errorf("chunk source = %q, want %q", chunk.Source, "sample-pdf-a4-size.pdf")
		}
	}
}
