package regionid

import (
	"testing"
)

func TestParse_ValidIDs(t *testing.T) {
	tests := []struct {
		input string
		want  ID
	}{
		{"main.0", ID{Surface: "main", Index: 0}},
		{"left-aux.0", ID{Surface: "left-aux", Index: 0}},
		{"right-aux.0", ID{Surface: "right-aux", Index: 0}},
		{"a.0", ID{Surface: "a", Index: 0}},
		{"main.1", ID{Surface: "main", Index: 1}},
		{"main.99", ID{Surface: "main", Index: 99}},
		{"abc123.5", ID{Surface: "abc123", Index: 5}},
		{"a-b-c.0", ID{Surface: "a-b-c", Index: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParse_Errors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{"", "empty region identifier"},
		{"Main.0", "invalid surface name"},
		{"MAIN.0", "invalid surface name"},
		{"left aux.0", "invalid surface name"},
		{"123surface.0", "invalid surface name"},
		{"-hyphen.0", "invalid surface name"},
		{"main.-1", "region index must be non-negative"},
		{"main.abc", "invalid region index"},
		{"main.1.2", "invalid region identifier format"},
		{"a.b.c", "invalid region identifier format"},
		{"main", "missing index"},
		{".0", "invalid surface name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatalf("Parse(%q) returned nil error, want error containing %q", tt.input, tt.wantErr)
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Errorf("Parse(%q) error = %q, want it to contain %q", tt.input, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestParseBareInt(t *testing.T) {
	tests := []struct {
		input  string
		wantN  int
		wantOK bool
	}{
		{"0", 0, true},
		{"1", 1, true},
		{"42", 42, true},
		{"100", 100, true},
		{"-1", 0, false},
		{"-99", 0, false},
		{"abc", 0, false},
		{"1.0", 0, false},
		{"main.0", 0, false},
		{"", 0, false},
		{"3a", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			n, ok := ParseBareInt(tt.input)
			if ok != tt.wantOK || n != tt.wantN {
				t.Errorf("ParseBareInt(%q) = (%d, %v), want (%d, %v)", tt.input, n, ok, tt.wantN, tt.wantOK)
			}
		})
	}
}

func TestID_String(t *testing.T) {
	tests := []struct {
		id   ID
		want string
	}{
		{ID{Surface: "main", Index: 0}, "main.0"},
		{ID{Surface: "left-aux", Index: 0}, "left-aux.0"},
		{ID{Surface: "right-aux", Index: 2}, "right-aux.2"},
		{ID{Surface: "a", Index: 99}, "a.99"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.id.String()
			if got != tt.want {
				t.Errorf("ID{%q, %d}.String() = %q, want %q", tt.id.Surface, tt.id.Index, got, tt.want)
			}
		})
	}
}

func TestParse_RoundTrip(t *testing.T) {
	ids := []ID{
		{Surface: "main", Index: 0},
		{Surface: "left-aux", Index: 0},
		{Surface: "right-aux", Index: 3},
		{Surface: "abc123", Index: 42},
	}

	for _, id := range ids {
		t.Run(id.String(), func(t *testing.T) {
			parsed, err := Parse(id.String())
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", id.String(), err)
			}
			if parsed != id {
				t.Errorf("round-trip failed: got %+v, want %+v", parsed, id)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
