package indexer

import (
	"path/filepath"
	"testing"
)

func TestIndexFilePath(t *testing.T) {
	tests := []struct {
		filePath string
		want     string
	}{
		{
			filePath: "/docs/manual.pdf",
			want:     "/docs/.index/manual.pdf.json",
		},
		{
			filePath: "/docs/sub/report.xlsx",
			want:     "/docs/sub/.index/report.xlsx.json",
		},
		{
			filePath: "/manual.pdf",
			want:     "/.index/manual.pdf.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			got := indexFilePath(tt.filePath)
			want := filepath.FromSlash(tt.want)
			if got != want {
				t.Errorf("indexFilePath(%q) = %q, want %q", tt.filePath, got, want)
			}
		})
	}
}
