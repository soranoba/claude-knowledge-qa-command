package indexer

import "strings"

type fileType int

const (
	fileTypePDF   fileType = iota
	fileTypeExcel fileType = iota
	fileTypeCSV   fileType = iota
)

// extFileTypes is the single source of truth for supported extensions.
var extFileTypes = map[string]fileType{
	".pdf":  fileTypePDF,
	".xlsx": fileTypeExcel,
	".xls":  fileTypeExcel,
	".xlsm": fileTypeExcel,
	".csv":  fileTypeCSV,
}

func fileTypeFromExt(ext string) (fileType, bool) {
	ft, ok := extFileTypes[strings.ToLower(ext)]
	return ft, ok
}
