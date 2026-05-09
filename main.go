package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/soranoba/knowledge-qa/internal/indexer"
	"github.com/soranoba/knowledge-qa/internal/search"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: knowledge-qa <directory|file> <question>")
		os.Exit(1)
	}

	path := os.Args[1]
	question := os.Args[2]

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	idx, err := indexer.Load(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := idx.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing index: %v\n", err)
		os.Exit(1)
	}

	results := search.Query(idx.AllChunks(), question, 8)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
