package indexer

import (
	"testing"
)

func TestFileTypeFromExt(t *testing.T) {
	tests := []struct {
		ext      string
		wantType fileType
		wantOK   bool
	}{
		{".pdf", fileTypePDF, true},
		{".PDF", fileTypePDF, true},
		{".xlsx", fileTypeExcel, true},
		{".xls", fileTypeExcel, true},
		{".xlsm", fileTypeExcel, true},
		{".XLSX", fileTypeExcel, true},
		{".txt", 0, false},
		{".doc", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got, ok := fileTypeFromExt(tt.ext)
			if ok != tt.wantOK {
				t.Errorf("fileTypeFromExt(%q) ok = %v, want %v", tt.ext, ok, tt.wantOK)
			}
			if ok && got != tt.wantType {
				t.Errorf("fileTypeFromExt(%q) = %v, want %v", tt.ext, got, tt.wantType)
			}
		})
	}
}
