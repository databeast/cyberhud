package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// IconEntry represents a parsed icon name → codepoint mapping.
type IconEntry struct {
	Name      string // snake_case icon name from codepoints file
	Codepoint rune   // Unicode codepoint value
}

// ParseCodepoints reads a codepoints file and returns the parsed entries.
// Lines must have exactly two whitespace-separated fields: name and hex codepoint.
// Blank lines and malformed lines are skipped.
// Duplicate names use last-occurrence-wins semantics.
func ParseCodepoints(r io.Reader) ([]IconEntry, error) {
	type entryInfo struct {
		entry     IconEntry
		firstLine int // line number of first occurrence (for warnings)
		lastLine  int // line number of last occurrence (determines order)
		lastPos   int // sequential position of last occurrence among valid lines
	}

	scanner := bufio.NewScanner(r)
	nameMap := make(map[string]*entryInfo)
	lineNum := 0
	validPos := 0 // counts valid parsed lines to track ordering

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		fields := strings.Fields(line)
		if len(fields) != 2 {
			// Skip blank lines and lines without exactly 2 fields.
			continue
		}

		name := fields[0]
		hexStr := fields[1]

		// Parse hex codepoint.
		val, err := strconv.ParseUint(hexStr, 16, 32)
		if err != nil || val > 0x10FFFF {
			fmt.Fprintf(os.Stderr, "gen-icons: warning: invalid codepoint at line %d: %q\n", lineNum, line)
			continue
		}

		cp := rune(val)
		validPos++

		if prev, exists := nameMap[name]; exists {
			// Duplicate name: warn and use last-occurrence-wins.
			fmt.Fprintf(os.Stderr, "gen-icons: warning: duplicate name %q at lines %d and %d\n", name, prev.firstLine, lineNum)
			prev.entry.Codepoint = cp
			prev.lastLine = lineNum
			prev.lastPos = validPos
		} else {
			nameMap[name] = &entryInfo{
				entry:     IconEntry{Name: name, Codepoint: cp},
				firstLine: lineNum,
				lastLine:  lineNum,
				lastPos:   validPos,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Build result slice ordered by last-seen position.
	infos := make([]*entryInfo, 0, len(nameMap))
	for _, info := range nameMap {
		infos = append(infos, info)
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].lastPos < infos[j].lastPos
	})

	entries := make([]IconEntry, len(infos))
	for i, info := range infos {
		entries[i] = info.entry
	}

	return entries, nil
}
