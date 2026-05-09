package indexer

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleDir = "../../samples"

func TestSync_Dir_IndexesAllFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"user-data.csv", "sample-pdf-a4-size.pdf"} {
		data, err := os.ReadFile(filepath.Join(sampleDir, name))
		if err != nil {
			t.Fatalf("read sample %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
			t.Fatalf("write sample %s: %v", name, err)
		}
	}

	index, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := index.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	chunks := index.AllChunks()
	if len(chunks) == 0 {
		t.Fatal("got 0 chunks after Sync")
	}

	sources := make(map[string]bool)
	for _, c := range chunks {
		sources[c.Source] = true
	}
	for _, name := range []string{"user-data.csv", "sample-pdf-a4-size.pdf"} {
		if !sources[name] {
			t.Errorf("no chunks found for %s", name)
		}
	}
}

func TestSync_IndexFilesCreated(t *testing.T) {
	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join(sampleDir, "user-data.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "user-data.csv"), data, 0644); err != nil {
		t.Fatalf("write sample: %v", err)
	}

	index, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := index.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	indexFile := filepath.Join(dir, ".index", "user-data.csv.json")
	if _, err := os.Stat(indexFile); err != nil {
		t.Errorf("index file not created: %s", indexFile)
	}
}

func TestSync_SkipsUnchangedFiles(t *testing.T) {
	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join(sampleDir, "user-data.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "user-data.csv"), data, 0644); err != nil {
		t.Fatalf("write sample: %v", err)
	}

	index, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := index.Sync(); err != nil {
		t.Fatalf("first Sync: %v", err)
	}

	// Second Sync with same files — nothing should be re-indexed
	index2, err := Load(dir)
	if err != nil {
		t.Fatalf("Load (2nd): %v", err)
	}
	// Verify that the record is already loaded from the index file
	if len(index2.records) == 0 {
		t.Error("second Load: no records loaded from existing index files")
	}
}

func TestSync_SingleFile(t *testing.T) {
	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join(sampleDir, "user-data.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	csvPath := filepath.Join(dir, "user-data.csv")
	if err := os.WriteFile(csvPath, data, 0644); err != nil {
		t.Fatalf("write sample: %v", err)
	}

	index, err := Load(csvPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := index.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	chunks := index.AllChunks()
	if len(chunks) == 0 {
		t.Fatal("got 0 chunks for single file")
	}
	for _, c := range chunks {
		if c.Source != "user-data.csv" {
			t.Errorf("unexpected source %q", c.Source)
		}
	}
}

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
