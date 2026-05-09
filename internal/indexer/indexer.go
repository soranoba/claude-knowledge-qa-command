package indexer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const indexDirName = ".index"

type Chunk struct {
	Source   string `json:"source"`
	Location string `json:"location"`
	Text     string `json:"text"`
}

type fileIndex struct {
	MTime  time.Time `json:"mtime"`
	Chunks []Chunk   `json:"chunks"`
}

// Index holds the target file list and per-file records.
// It has no knowledge of a base directory; all paths are absolute.
type Index struct {
	files   []string
	records map[string]*fileIndex // key: absolute source file path
}

// Load accepts a directory or a single file path.
func Load(path string) (*Index, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", absPath)
	}

	index := &Index{
		files:   make([]string, 0),
		records: make(map[string]*fileIndex),
	}

	if fileInfo.IsDir() {
		files, err := ScanTargetFiles(absPath)
		if err != nil {
			return nil, err
		}
		index.files = files
	} else {
		index.files = []string{absPath}
	}

	index.loadExistingRecords()
	return index, nil
}

// indexFilePath derives the index file path from the source file path.
// e.g. /docs/manual.pdf -> /docs/.index/manual.pdf.json
func indexFilePath(filePath string) string {
	return filepath.Join(filepath.Dir(filePath), indexDirName, filepath.Base(filePath)+".json")
}

// Sync updates the index for any file in index.files that is new or modified.
func (index *Index) Sync() error {
	for _, absPath := range index.files {
		info, err := os.Stat(absPath)
		if err != nil {
			continue
		}

		rec, exists := index.records[absPath]
		if exists && !info.ModTime().After(rec.MTime) {
			continue
		}

		chunks, err := extractChunks(absPath)
		if err != nil {
			return err
		}

		fi := &fileIndex{MTime: info.ModTime(), Chunks: chunks}
		if err := saveFileIndex(absPath, fi); err != nil {
			return err
		}
		index.records[absPath] = fi
	}
	return nil
}

// AllChunks returns chunks from every indexed file.
func (index *Index) AllChunks() []Chunk {
	var all []Chunk
	for _, rec := range index.records {
		all = append(all, rec.Chunks...)
	}
	return all
}

func (index *Index) loadExistingRecords() {
	for _, absPath := range index.files {
		data, err := os.ReadFile(indexFilePath(absPath))
		if err != nil {
			continue
		}
		var fi fileIndex
		if err := json.Unmarshal(data, &fi); err == nil {
			index.records[absPath] = &fi
		}
	}
}

func saveFileIndex(filePath string, fi *fileIndex) error {
	path := indexFilePath(filePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create index directory: %w", err)
	}
	data, err := json.MarshalIndent(fi, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ScanTargetFiles walks dir and returns absolute paths of all PDF/Excel files.
func ScanTargetFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if _, ok := fileTypeFromExt(filepath.Ext(path)); ok {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func extractChunks(path string) ([]Chunk, error) {
	ft, ok := fileTypeFromExt(filepath.Ext(path))
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", filepath.Ext(path))
	}
	switch ft {
	case fileTypePDF:
		return extractPDF(path)
	case fileTypeExcel:
		return extractExcel(path)
	case fileTypeCSV:
		return extractCSV(path)
	default:
		return nil, fmt.Errorf("unsupported format: %s", filepath.Ext(path))
	}
}
