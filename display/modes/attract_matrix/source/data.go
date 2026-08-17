package source

// MatrixSnapshot is the mode-specific snapshot type consumed by style Build
// methods via style.Style[MatrixSnapshot]. It captures the rendering parameters
// needed to produce a single frame of the matrix rain effect.
type MatrixSnapshot struct {
	Policy       Policy // Active policy snapshot at frame start
	PanelWidth   int    // Panel pixel width
	PanelHeight  int    // Panel pixel height
	GlyphAdvance int    // Matrix-code font glyph advance (typically 14px)
	RowHeight    int    // Matrix-code font row height (typically 14px)
	Mono         bool   // True for monochrome OLED panels
	Eink         bool   // True for e-ink panels (static frame, no animation)
}
