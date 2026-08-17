package bdf

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Glyph represents a single parsed BDF glyph.
type Glyph struct {
	Codepoint rune
	Width     int      // BBX width
	Height    int      // BBX height
	XOff      int      // BBX x-offset
	YOff      int      // BBX y-offset
	DWidth    int      // DWIDTH advance width
	Rows      []uint32 // Bitmap rows, MSB-left aligned within glyph width
}

// Font represents a parsed BDF font containing retained glyphs.
type Font struct {
	Name   string
	Glyphs map[rune]Glyph
}

// CodepointRange defines an inclusive range of Unicode codepoints to retain.
type CodepointRange struct {
	Low  rune // Inclusive lower bound
	High rune // Inclusive upper bound
}

// ParseConfig controls which glyphs are retained from the BDF source.
type ParseConfig struct {
	Ranges []CodepointRange // Codepoint ranges to include (nil/empty defaults to ASCII 32-126)
}

// Parse reads a BDF font from r and returns the parsed Font.
// Only ASCII codepoints 32–126 are retained; all others are skipped.
// Returns an error if the BDF data is malformed or contains no matching glyphs.
func Parse(r io.Reader) (*Font, error) {
	return ParseWithConfig(r, ParseConfig{})
}

// ParseWithConfig reads a BDF font, retaining only glyphs within configured ranges.
// Falls back to ASCII 32-126 when cfg.Ranges is nil or empty.
func ParseWithConfig(r io.Reader, cfg ParseConfig) (*Font, error) {
	ranges := cfg.Ranges
	if len(ranges) == 0 {
		ranges = []CodepointRange{{Low: 32, High: 126}}
	}
	scanner := bufio.NewScanner(r)

	// First non-empty line must be STARTFONT.
	foundStart := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "STARTFONT") {
			foundStart = true
		}
		break
	}
	if !foundStart {
		return nil, fmt.Errorf("not a valid BDF file (missing STARTFONT)")
	}

	font := &Font{
		Glyphs: make(map[rune]Glyph),
	}

	// Parse the rest of the file.
	var (
		inChar   bool
		inBitmap bool
		glyph    Glyph
		bbxSet   bool
		encSet   bool
		rows     []string
		fontName string
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !inChar {
			if strings.HasPrefix(trimmed, "FONT ") {
				fontName = strings.TrimPrefix(trimmed, "FONT ")
			} else if strings.HasPrefix(trimmed, "STARTCHAR") {
				inChar = true
				inBitmap = false
				glyph = Glyph{}
				bbxSet = false
				encSet = false
				rows = nil
			}
			continue
		}

		// Inside a character definition.
		if inBitmap {
			if trimmed == "ENDCHAR" {
				// Finalize glyph.
				if !encSet {
					inChar = false
					inBitmap = false
					continue
				}
				if !bbxSet {
					inChar = false
					inBitmap = false
					continue
				}

				// Validate row count.
				if len(rows) != glyph.Height {
					return nil, fmt.Errorf("glyph U+%04X has %d bitmap rows, expected %d",
						glyph.Codepoint, len(rows), glyph.Height)
				}

				// Convert hex rows to uint32 with MSB-left alignment.
				glyph.Rows = make([]uint32, glyph.Height)
				for i, hexRow := range rows {
					val, err := parseHexRow(hexRow)
					if err != nil {
						return nil, fmt.Errorf("glyph U+%04X row %d: invalid hex %s",
							glyph.Codepoint, i, hexRow)
					}
					// The hex data represents pixels MSB-left within the byte-padded width.
					// Shift to align MSB of hex data to bit 31 of uint32.
					hexBits := len(hexRow) * 4
					shifted := val << (32 - hexBits)
					// Apply xoff: shift right by xoff to position within glyph width.
					if glyph.XOff > 0 {
						shifted >>= uint(glyph.XOff)
					} else if glyph.XOff < 0 {
						shifted <<= uint(-glyph.XOff)
					}
					glyph.Rows[i] = shifted
				}

				// Retain only glyphs within configured ranges.
				if inRanges(glyph.Codepoint, ranges) {
					font.Glyphs[glyph.Codepoint] = glyph
				}

				inChar = false
				inBitmap = false
			} else {
				// Collect bitmap hex row.
				rows = append(rows, trimmed)
			}
			continue
		}

		// Not yet in bitmap section.
		switch {
		case strings.HasPrefix(trimmed, "ENCODING "):
			val, err := strconv.Atoi(strings.TrimPrefix(trimmed, "ENCODING "))
			if err == nil {
				glyph.Codepoint = rune(val)
				encSet = true
			}

		case strings.HasPrefix(trimmed, "BBX "):
			parts := strings.Fields(trimmed)
			if len(parts) >= 5 {
				w, _ := strconv.Atoi(parts[1])
				h, _ := strconv.Atoi(parts[2])
				xo, _ := strconv.Atoi(parts[3])
				yo, _ := strconv.Atoi(parts[4])
				glyph.Width = w
				glyph.Height = h
				glyph.XOff = xo
				glyph.YOff = yo
				bbxSet = true
			}

		case strings.HasPrefix(trimmed, "DWIDTH "):
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				dw, _ := strconv.Atoi(parts[1])
				glyph.DWidth = dw
			}

		case trimmed == "BITMAP":
			inBitmap = true

		case trimmed == "ENDCHAR":
			// ENDCHAR without BITMAP — skip glyph.
			inChar = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	font.Name = fontName

	if len(font.Glyphs) == 0 {
		return nil, fmt.Errorf("no glyphs found in configured ranges")
	}

	return font, nil
}

// inRanges reports whether cp falls within at least one of the given ranges.
func inRanges(cp rune, ranges []CodepointRange) bool {
	for _, r := range ranges {
		if cp >= r.Low && cp <= r.High {
			return true
		}
	}
	return false
}

// parseHexRow parses a hex string into a uint32 value.
// The hex string may be 1-8 hex characters (representing 4-32 bits).
func parseHexRow(s string) (uint32, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0, fmt.Errorf("empty hex string")
	}
	// Validate it is valid hex.
	if len(s) > 8 {
		return 0, fmt.Errorf("hex string too long: %s", s)
	}
	// Use hex.DecodeString for validation (needs even length).
	normalized := s
	if len(normalized)%2 != 0 {
		normalized = "0" + normalized
	}
	if _, err := hex.DecodeString(normalized); err != nil {
		return 0, fmt.Errorf("invalid hex: %s", s)
	}
	val, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(val), nil
}
