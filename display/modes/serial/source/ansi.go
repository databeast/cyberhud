package source

import (
	"image/color"
	"strconv"
	"strings"
)

// ColorSegment represents a colored portion of parsed text.
type ColorSegment struct {
	Start  int         // Byte offset in stripped text
	Length int         // Byte length of this segment
	Color  color.Color // nil means default foreground
}

// StripAnsi removes all ANSI escape sequences (ESC[...m pattern) from input,
// returning clean text. Only sequences matching ESC[ <digits/semicolons> m are
// stripped. Malformed or non-SGR sequences pass through unchanged.
func StripAnsi(input []byte) []byte {
	if len(input) == 0 {
		return input
	}

	out := make([]byte, 0, len(input))
	i := 0
	for i < len(input) {
		if input[i] == 0x1b && i+1 < len(input) && input[i+1] == '[' {
			// Found ESC[ — scan for the terminating 'm'.
			j := i + 2
			matched := false
			for j < len(input) {
				if input[j] == 'm' {
					// Valid SGR sequence ESC[...m — skip it entirely.
					i = j + 1
					matched = true
					break
				}
				// SGR parameters are digits and semicolons.
				if (input[j] >= '0' && input[j] <= '9') || input[j] == ';' {
					j++
					continue
				}
				// Non-SGR CSI sequence or malformed — pass through unchanged.
				break
			}
			if !matched {
				// Not a valid ESC[...m sequence — copy current byte and advance.
				out = append(out, input[i])
				i++
			}
		} else {
			out = append(out, input[i])
			i++
		}
	}
	return out
}

// ansiToColor maps ANSI SGR foreground codes to color.Color values.
// Codes 30–37: standard 8 colors. Codes 90–97: bright 8 colors.
// Code 0 (reset) returns nil. Unrecognized codes return nil.
func ansiToColor(code int) color.Color {
	// Standard colors (30–37)
	standard := [8]color.Color{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},       // 30: Black
		color.RGBA{R: 170, G: 0, B: 0, A: 255},     // 31: Red
		color.RGBA{R: 0, G: 170, B: 0, A: 255},     // 32: Green
		color.RGBA{R: 170, G: 170, B: 0, A: 255},   // 33: Yellow
		color.RGBA{R: 0, G: 0, B: 170, A: 255},     // 34: Blue
		color.RGBA{R: 170, G: 0, B: 170, A: 255},   // 35: Magenta
		color.RGBA{R: 0, G: 170, B: 170, A: 255},   // 36: Cyan
		color.RGBA{R: 170, G: 170, B: 170, A: 255}, // 37: White
	}

	// Bright colors (90–97)
	bright := [8]color.Color{
		color.RGBA{R: 85, G: 85, B: 85, A: 255},    // 90: Bright Black
		color.RGBA{R: 255, G: 85, B: 85, A: 255},   // 91: Bright Red
		color.RGBA{R: 85, G: 255, B: 85, A: 255},   // 92: Bright Green
		color.RGBA{R: 255, G: 255, B: 85, A: 255},  // 93: Bright Yellow
		color.RGBA{R: 85, G: 85, B: 255, A: 255},   // 94: Bright Blue
		color.RGBA{R: 255, G: 85, B: 255, A: 255},  // 95: Bright Magenta
		color.RGBA{R: 85, G: 255, B: 255, A: 255},  // 96: Bright Cyan
		color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 97: Bright White
	}

	switch {
	case code == 0:
		return nil
	case code >= 30 && code <= 37:
		return standard[code-30]
	case code >= 90 && code <= 97:
		return bright[code-90]
	default:
		return nil
	}
}

// ParseLine processes a single line of serial data, stripping ANSI escape
// sequences and returning clean text with a color segment map.
func ParseLine(raw string) (text string, segments []ColorSegment) {
	if raw == "" {
		return "", nil
	}

	var cleaned strings.Builder
	cleaned.Grow(len(raw))

	var currentColor color.Color
	var segStart int

	i := 0
	for i < len(raw) {
		if raw[i] == 0x1b && i+1 < len(raw) && raw[i+1] == '[' {
			// Found ESC[ — parse the SGR sequence.
			j := i + 2
			matched := false
			for j < len(raw) {
				if raw[j] == 'm' {
					// Extract parameter string between ESC[ and m.
					params := raw[i+2 : j]
					newColor := parseSGRParams(params, currentColor)

					// Close current segment if we have accumulated text.
					pos := cleaned.Len()
					if pos > segStart {
						segments = append(segments, ColorSegment{
							Start:  segStart,
							Length: pos - segStart,
							Color:  currentColor,
						})
					}
					currentColor = newColor
					segStart = pos

					i = j + 1
					matched = true
					break
				}
				if (raw[j] >= '0' && raw[j] <= '9') || raw[j] == ';' {
					j++
					continue
				}
				// Non-SGR sequence — pass through unchanged.
				break
			}
			if !matched {
				// Not a valid ESC[...m — copy current byte and advance.
				cleaned.WriteByte(raw[i])
				i++
			}
		} else {
			cleaned.WriteByte(raw[i])
			i++
		}
	}

	text = cleaned.String()

	// Close final segment if there's remaining text.
	if cleaned.Len() > segStart {
		segments = append(segments, ColorSegment{
			Start:  segStart,
			Length: cleaned.Len() - segStart,
			Color:  currentColor,
		})
	}

	return text, segments
}

// parseSGRParams parses semicolon-separated SGR parameters and returns the
// resulting color. Only foreground color codes (30–37, 90–97) and reset (0)
// are handled. Unrecognized parameters are ignored.
func parseSGRParams(params string, current color.Color) color.Color {
	// ESC[m is equivalent to ESC[0m (reset).
	if params == "" {
		return nil
	}

	parts := strings.Split(params, ";")
	result := current
	for _, part := range parts {
		code, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			// Malformed parameter — ignore.
			continue
		}
		if code == 0 {
			result = nil
		} else if (code >= 30 && code <= 37) || (code >= 90 && code <= 97) {
			result = ansiToColor(code)
		}
		// Unrecognized codes (bold, underline, background, etc.) are ignored.
	}
	return result
}
