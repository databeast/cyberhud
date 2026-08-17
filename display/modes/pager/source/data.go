package source

// PagerSnapshot is the mode-specific snapshot type for the style registry.
// It captures the state needed by style Build methods to produce rendered output.
type PagerSnapshot struct {
	Lines   []string  // visible text lines from the buffer
	Policy  Policy    // current pager configuration
	IsSlow  bool      // true for slow-refresh display surfaces
	Phase   PagePhase // current page transition phase (page strategy)
	Alpha   float64   // current fade alpha (page strategy)
	OffsetY int       // pixel scroll offset (scroll strategy)
}

// Layout holds the computed visible dimensions for text rendering.
// It is derived from the resolved font metrics and the panel's pixel area.
type Layout struct {
	VisibleColumns int // floor(PixelWidth / GlyphAdvance)
	VisibleRows    int // floor(PixelHeight / RowHeight)
	GlyphAdvance   int // effective glyph advance used
	RowHeight      int // effective row height used
	PixelWidth     int // panel pixel width
	PixelHeight    int // panel pixel height
}
