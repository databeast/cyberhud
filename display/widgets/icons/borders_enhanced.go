package icons

// borders_enhanced.go registers additional border theme tiles for the enhanced
// borderframe widget. These complement the base tiles in borders.go.
//
// Themes registered here:
//   - "border/round" aliases (edge tiles aliasing sharp edges)
//   - "doubleline" (parallel-line border tiles)
//   - "circuit" (PCB-trace-style tiles)
//   - "hex" (hexagonal-segment tiles)
//   - Accent tiles for all themes (sharp, rounded, doubleline, circuit, hex)

func init() {
	registerRoundedAliases()
	registerDoubleLineTiles()
	registerCircuitTiles()
	registerHexTiles()
	registerAccentTiles()
}

// registerRoundedAliases registers edge tile aliases for the "border/round" prefix.
// The rounded theme uses the same horizontal and vertical edge tiles as the sharp theme,
// and maps its corner names to the existing round corner tiles.
func registerRoundedAliases() {
	// Edge aliases: border/round/h and border/round/v use the same images as border/h and border/v.
	Register("border/round/h", makeAlpha8x8([]byte{
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00,
	}))
	Register("border/round/v", makeAlpha8x8([]byte{
		0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18, 0x18,
	}))

	// Corner aliases: border/round/corner-* map to the same images as border/round-*.
	Register("border/round/corner-tl", makeAlpha8x8([]byte{
		0x00, 0x00, 0x00, 0x07, 0x0F, 0x18, 0x18, 0x18,
	}))
	Register("border/round/corner-tr", makeAlpha8x8([]byte{
		0x00, 0x00, 0x00, 0xE0, 0xF0, 0x18, 0x18, 0x18,
	}))
	Register("border/round/corner-bl", makeAlpha8x8([]byte{
		0x18, 0x18, 0x18, 0x0F, 0x07, 0x00, 0x00, 0x00,
	}))
	Register("border/round/corner-br", makeAlpha8x8([]byte{
		0x18, 0x18, 0x18, 0xF0, 0xE0, 0x00, 0x00, 0x00,
	}))
}

// registerDoubleLineTiles registers parallel double-line border tiles.
// Horizontal: two 1px lines spaced 2px apart (rows 2 and 5).
// Vertical: two 1px columns spaced 2px apart (bits 2 and 5).
// Corners: double-line right-angle bends.
func registerDoubleLineTiles() {
	// Horizontal edge: two parallel horizontal lines at rows 2 and 5.
	Register("doubleline/h", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b11111111,
		0b00000000,
		0b00000000,
		0b11111111,
		0b00000000,
		0b00000000,
	}))

	// Vertical edge: two parallel vertical lines at columns 2 and 5.
	Register("doubleline/v", makeAlpha8x8([]byte{
		0b00100100,
		0b00100100,
		0b00100100,
		0b00100100,
		0b00100100,
		0b00100100,
		0b00100100,
		0b00100100,
	}))

	// Top-left corner: outer L at row 2/col 2, inner L at row 5/col 5.
	Register("doubleline/corner-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00111111,
		0b00100100,
		0b00100100,
		0b00000111,
		0b00000100,
		0b00000100,
	}))

	// Top-right corner: outer L at row 2/col 5, inner L at row 5/col 2.
	Register("doubleline/corner-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b11111100,
		0b00100100,
		0b00100100,
		0b11100000,
		0b00100000,
		0b00100000,
	}))

	// Bottom-left corner: L bends upward-left.
	Register("doubleline/corner-bl", makeAlpha8x8([]byte{
		0b00000100,
		0b00000100,
		0b00000111,
		0b00100100,
		0b00100100,
		0b00111111,
		0b00000000,
		0b00000000,
	}))

	// Bottom-right corner: L bends upward-right.
	Register("doubleline/corner-br", makeAlpha8x8([]byte{
		0b00100000,
		0b00100000,
		0b11100000,
		0b00100100,
		0b00100100,
		0b11111100,
		0b00000000,
		0b00000000,
	}))
}

// registerCircuitTiles registers PCB-trace-style border tiles.
// Edges are 2px-wide traces with periodic stubs (via pads).
// Corners are 90° trace bends with a solder-dot at the vertex.
func registerCircuitTiles() {
	// Horizontal edge: 2px-wide trace (rows 3-4) with stub marks at columns 1 and 6.
	Register("circuit/h", makeAlpha8x8([]byte{
		0b00000000,
		0b01000010,
		0b01000010,
		0b11111111,
		0b11111111,
		0b00000000,
		0b00000000,
		0b00000000,
	}))

	// Vertical edge: 2px-wide trace (cols 3-4) with stub marks at rows 1 and 6.
	Register("circuit/v", makeAlpha8x8([]byte{
		0b00000000,
		0b00011000,
		0b00011000,
		0b00011000,
		0b00011000,
		0b00011000,
		0b00111100,
		0b00011000,
	}))

	// Top-left corner: 90° trace bend with solder dot at vertex (2x2 filled at row 3-4, col 3-4).
	Register("circuit/corner-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00000000,
		0b00011111,
		0b00011111,
		0b00011000,
		0b00011000,
		0b00011000,
	}))

	// Top-right corner: 90° trace bend with solder dot.
	Register("circuit/corner-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00000000,
		0b11111000,
		0b11111000,
		0b00011000,
		0b00011000,
		0b00011000,
	}))

	// Bottom-left corner: 90° trace bend with solder dot.
	Register("circuit/corner-bl", makeAlpha8x8([]byte{
		0b00011000,
		0b00011000,
		0b00011000,
		0b00011111,
		0b00011111,
		0b00000000,
		0b00000000,
		0b00000000,
	}))

	// Bottom-right corner: 90° trace bend with solder dot.
	Register("circuit/corner-br", makeAlpha8x8([]byte{
		0b00011000,
		0b00011000,
		0b00011000,
		0b11111000,
		0b11111000,
		0b00000000,
		0b00000000,
		0b00000000,
	}))
}

// registerHexTiles registers hexagonal-segment border tiles.
// Edges are angled line segments evoking tessellated hexagons.
// Corners are 120°-angle junctions rather than 90°.
func registerHexTiles() {
	// Horizontal edge: angled segments suggesting hex tessellation.
	Register("hex/h", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b01000010,
		0b00100100,
		0b00011000,
		0b00100100,
		0b01000010,
		0b00000000,
	}))

	// Vertical edge: angled segments for hex tessellation.
	Register("hex/v", makeAlpha8x8([]byte{
		0b00010000,
		0b00100000,
		0b01000000,
		0b00100000,
		0b00010000,
		0b00001000,
		0b00000100,
		0b00001000,
	}))

	// Top-left corner: 120° angle junction.
	Register("hex/corner-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00000100,
		0b00001000,
		0b00010000,
		0b00100000,
		0b01000000,
		0b00100000,
	}))

	// Top-right corner: 120° angle junction.
	Register("hex/corner-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00100000,
		0b00010000,
		0b00001000,
		0b00000100,
		0b00000010,
		0b00000100,
	}))

	// Bottom-left corner: 120° angle junction.
	Register("hex/corner-bl", makeAlpha8x8([]byte{
		0b00100000,
		0b01000000,
		0b00100000,
		0b00010000,
		0b00001000,
		0b00000100,
		0b00000000,
		0b00000000,
	}))

	// Bottom-right corner: 120° angle junction.
	Register("hex/corner-br", makeAlpha8x8([]byte{
		0b00000100,
		0b00000010,
		0b00000100,
		0b00001000,
		0b00010000,
		0b00100000,
		0b00000000,
		0b00000000,
	}))
}

// registerAccentTiles registers visually heavier/decorative corner accent tiles
// for all five themes. These are used when CornerAccent is enabled in the Config.
func registerAccentTiles() {
	// Sharp accents: filled triangles in each corner.
	Register("border/accent-tl", makeAlpha8x8([]byte{
		0b11111111,
		0b11111110,
		0b11111100,
		0b11111000,
		0b11110000,
		0b11100000,
		0b11000000,
		0b10000000,
	}))
	Register("border/accent-tr", makeAlpha8x8([]byte{
		0b11111111,
		0b01111111,
		0b00111111,
		0b00011111,
		0b00001111,
		0b00000111,
		0b00000011,
		0b00000001,
	}))
	Register("border/accent-bl", makeAlpha8x8([]byte{
		0b10000000,
		0b11000000,
		0b11100000,
		0b11110000,
		0b11111000,
		0b11111100,
		0b11111110,
		0b11111111,
	}))
	Register("border/accent-br", makeAlpha8x8([]byte{
		0b00000001,
		0b00000011,
		0b00000111,
		0b00001111,
		0b00011111,
		0b00111111,
		0b01111111,
		0b11111111,
	}))

	// Rounded accents: filled quarter-circle arcs (heavier than the 1px line version).
	Register("border/round/accent-tl", makeAlpha8x8([]byte{
		0b00001111,
		0b00011111,
		0b00111111,
		0b01111111,
		0b01111000,
		0b11110000,
		0b11100000,
		0b11000000,
	}))
	Register("border/round/accent-tr", makeAlpha8x8([]byte{
		0b11110000,
		0b11111000,
		0b11111100,
		0b11111110,
		0b00011110,
		0b00001111,
		0b00000111,
		0b00000011,
	}))
	Register("border/round/accent-bl", makeAlpha8x8([]byte{
		0b11000000,
		0b11100000,
		0b11110000,
		0b01111000,
		0b01111111,
		0b00111111,
		0b00011111,
		0b00001111,
	}))
	Register("border/round/accent-br", makeAlpha8x8([]byte{
		0b00000011,
		0b00000111,
		0b00001111,
		0b00011110,
		0b11111110,
		0b11111100,
		0b11111000,
		0b11110000,
	}))

	// Double-line accents: thicker double-L corners with filled inner area.
	Register("doubleline/accent-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00111111,
		0b00111111,
		0b00110100,
		0b00110111,
		0b00110111,
		0b00110100,
	}))
	Register("doubleline/accent-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b11111100,
		0b11111100,
		0b00101100,
		0b11101100,
		0b11101100,
		0b00101100,
	}))
	Register("doubleline/accent-bl", makeAlpha8x8([]byte{
		0b00110100,
		0b00110111,
		0b00110111,
		0b00110100,
		0b00111111,
		0b00111111,
		0b00000000,
		0b00000000,
	}))
	Register("doubleline/accent-br", makeAlpha8x8([]byte{
		0b00101100,
		0b11101100,
		0b11101100,
		0b00101100,
		0b11111100,
		0b11111100,
		0b00000000,
		0b00000000,
	}))

	// Circuit accents: junction nodes with solder dots (filled circles at vertex).
	Register("circuit/accent-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00000000,
		0b00111111,
		0b00111111,
		0b00111100,
		0b00111100,
		0b00011000,
	}))
	Register("circuit/accent-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b00000000,
		0b00000000,
		0b11111100,
		0b11111100,
		0b00111100,
		0b00111100,
		0b00011000,
	}))
	Register("circuit/accent-bl", makeAlpha8x8([]byte{
		0b00011000,
		0b00111100,
		0b00111100,
		0b00111111,
		0b00111111,
		0b00000000,
		0b00000000,
		0b00000000,
	}))
	Register("circuit/accent-br", makeAlpha8x8([]byte{
		0b00011000,
		0b00111100,
		0b00111100,
		0b11111100,
		0b11111100,
		0b00000000,
		0b00000000,
		0b00000000,
	}))

	// Hex accents: hex rosette patterns (heavier hex junction nodes).
	Register("hex/accent-tl", makeAlpha8x8([]byte{
		0b00000000,
		0b00001110,
		0b00011111,
		0b00111110,
		0b01111100,
		0b00111000,
		0b00010000,
		0b00100000,
	}))
	Register("hex/accent-tr", makeAlpha8x8([]byte{
		0b00000000,
		0b01110000,
		0b11111000,
		0b01111100,
		0b00111110,
		0b00011100,
		0b00001000,
		0b00000100,
	}))
	Register("hex/accent-bl", makeAlpha8x8([]byte{
		0b00100000,
		0b00010000,
		0b00111000,
		0b01111100,
		0b00111110,
		0b00011111,
		0b00001110,
		0b00000000,
	}))
	Register("hex/accent-br", makeAlpha8x8([]byte{
		0b00000100,
		0b00001000,
		0b00011100,
		0b00111110,
		0b01111100,
		0b11111000,
		0b01110000,
		0b00000000,
	}))
}
