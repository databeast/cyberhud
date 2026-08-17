package main

import (
	"testing"
)

func TestRunMissingRequiredFlags(t *testing.T) {
	tests := []struct {
		name      string
		bdf       string
		out       string
		id        string
		pkg       string
		width     int
		height    int
		advance   int
		rowheight int
		clone     string
		wantErr   string
	}{
		{
			name:      "missing -bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing required flag: -bdf",
		},
		{
			name:      "missing -out",
			bdf:       "spleen-5x8.bdf",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing required flag: -out",
		},
		{
			name:      "missing -id",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			pkg:       "font",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing required flag: -id",
		},
		{
			name:      "missing -pkg",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing required flag: -pkg",
		},
		{
			name:      "missing -width",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing or invalid flag: -width",
		},
		{
			name:      "missing -height",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			advance:   6,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing or invalid flag: -height",
		},
		{
			name:      "missing -advance",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			height:    8,
			rowheight: 10,
			clone:     "spleen",
			wantErr:   "missing or invalid flag: -advance",
		},
		{
			name:    "missing -rowheight",
			bdf:     "spleen-5x8.bdf",
			out:     "out.go",
			id:      "spleen-5x8",
			pkg:     "font",
			width:   5,
			height:  8,
			advance: 6,
			clone:   "spleen",
			wantErr: "missing or invalid flag: -rowheight",
		},
		{
			name:      "missing -clone and -download",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			wantErr:   "must specify either -ttf or -clone/-download",
		},
		{
			name:      "invalid -clone value",
			bdf:       "spleen-5x8.bdf",
			out:       "out.go",
			id:        "spleen-5x8",
			pkg:       "font",
			width:     5,
			height:    8,
			advance:   6,
			rowheight: 10,
			clone:     "unknown",
			wantErr:   "unknown clone target",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.bdf, "", tt.out, tt.id, tt.pkg, tt.width, tt.height, tt.advance, tt.rowheight, tt.clone, "", 0, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := err.Error(); !containsSubstring(got, tt.wantErr) {
				t.Errorf("error = %q, want substring %q", got, tt.wantErr)
			}
		})
	}
}

func TestIdToIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"spleen-5x8", "spleen5x8"},
		{"cozette-6x13", "cozette6x13"},
		{"terminus-16x32", "terminus16x32"},
		{"abc", "abc"},
		{"a-b-c", "abc"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := idToIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("idToIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIdToExportedIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"spleen-5x8", "Spleen5x8"},
		{"cozette-6x13", "Cozette6x13"},
		{"terminus-16x32", "Terminus16x32"},
		{"abc", "Abc"},
		{"a-b-c", "ABC"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := idToExportedIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("idToExportedIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDerivedNames(t *testing.T) {
	// Verify the naming convention for derived identifiers matches the design document.
	// e.g. fontID "spleen-5x8" → struct: spleen5x8Face, const: Spleen5x8ID, array: spleen5x8Rows
	tests := []struct {
		fontID     string
		wantStruct string
		wantConst  string
		wantArray  string
	}{
		{"spleen-5x8", "spleen5x8Face", "Spleen5x8ID", "spleen5x8Rows"},
		{"spleen-6x12", "spleen6x12Face", "Spleen6x12ID", "spleen6x12Rows"},
		{"cozette-6x13", "cozette6x13Face", "Cozette6x13ID", "cozette6x13Rows"},
		{"terminus-8x16", "terminus8x16Face", "Terminus8x16ID", "terminus8x16Rows"},
	}

	for _, tt := range tests {
		t.Run(tt.fontID, func(t *testing.T) {
			structName := idToIdentifier(tt.fontID) + "Face"
			constName := idToExportedIdentifier(tt.fontID) + "ID"
			arrayName := idToIdentifier(tt.fontID) + "Rows"

			if structName != tt.wantStruct {
				t.Errorf("structName = %q, want %q", structName, tt.wantStruct)
			}
			if constName != tt.wantConst {
				t.Errorf("constName = %q, want %q", constName, tt.wantConst)
			}
			if arrayName != tt.wantArray {
				t.Errorf("arrayName = %q, want %q", arrayName, tt.wantArray)
			}
		})
	}
}

// containsSubstring checks if s contains substr.
func containsSubstring(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
