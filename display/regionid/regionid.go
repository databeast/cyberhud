// Package regionid provides parsing and resolution of region identifiers
// used by both CLI and daemon.
//
// A region identifier uses the notation <surface_name>.<index> (e.g., "main.0",
// "left-aux.0"). Bare integers (e.g., "0", "2") are detected separately via
// ParseBareInt and resolve to the coordinator panel at that index — the caller
// is responsible for that lookup.
package regionid

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// surfaceNameRe matches valid surface names: starts with lowercase letter,
// followed by zero or more lowercase letters, digits, or hyphens.
var surfaceNameRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ID represents a resolved region identifier.
type ID struct {
	Surface string // e.g., "main", "left-aux"
	Index   int    // e.g., 0 (currently always 0, future: subdivisions)
}

// String returns the canonical <surface>.<index> form.
func (id ID) String() string {
	return fmt.Sprintf("%s.%d", id.Surface, id.Index)
}

// Parse interprets a region identifier string.
//
// Accepts:
//   - "main.0"      → ID{Surface: "main", Index: 0}
//   - "left-aux.0"  → ID{Surface: "left-aux", Index: 0}
//
// Bare integers (e.g., "0") are NOT handled here — they are detected by
// ParseBareInt. If a bare integer is passed to Parse, it will fail validation
// as an invalid surface name.
func Parse(s string) (ID, error) {
	if s == "" {
		return ID{}, errors.New("empty region identifier")
	}

	// Check for multiple dots.
	dotCount := strings.Count(s, ".")
	if dotCount > 1 {
		return ID{}, errors.New("invalid region identifier format")
	}

	// No dot — treat as surface-only (which must still be valid).
	if dotCount == 0 {
		if !surfaceNameRe.MatchString(s) {
			return ID{}, fmt.Errorf("invalid surface name %q: must match [a-z][a-z0-9-]*", s)
		}
		// A surface name alone without an index is ambiguous. However, the design
		// only shows <surface>.<index> as the valid form. A bare valid surface name
		// without a dot is not a complete region ID — but we allow it with index 0
		// for ergonomics? No — per the design, the canonical form is always
		// <surface>.<index>. A bare surface name with no dot should produce an error
		// if it's not a bare integer.
		//
		// Actually, re-reading the design: Parse accepts "main.0" and "left-aux.0".
		// Bare integers go through ParseBareInt. A bare surface name like "main"
		// (no dot) isn't listed as a valid input to Parse.
		//
		// Return an error for a bare surface name without index.
		return ID{}, fmt.Errorf("invalid region identifier %q: missing index (expected format: <surface>.<index>)", s)
	}

	// Exactly one dot — split into surface and index parts.
	parts := strings.SplitN(s, ".", 2)
	surface := parts[0]
	indexStr := parts[1]

	// Validate surface name.
	if !surfaceNameRe.MatchString(surface) {
		return ID{}, fmt.Errorf("invalid surface name %q: must match [a-z][a-z0-9-]*", surface)
	}

	// Parse the index.
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return ID{}, fmt.Errorf("invalid region index %q", indexStr)
	}

	if index < 0 {
		return ID{}, errors.New("region index must be non-negative")
	}

	return ID{Surface: surface, Index: index}, nil
}

// ParseBareInt detects if the string is a bare non-negative integer.
// Returns the integer and true if so, 0 and false otherwise.
func ParseBareInt(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	if n < 0 {
		return 0, false
	}
	return n, true
}
