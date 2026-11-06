package source

// PagerSnapshot is the fully-resolved frame handed to a pager style's Build
// method. It is the whole interface between the pager's mode logic and its
// styles.
//
// Every decision that depends on time, animation state, or the buffer is made
// before a style sees this struct. BuildView tails the buffer, classifies the
// panel, advances the scroll or page animation, and writes the settled result
// here. A style then does nothing but place these already-chosen lines on the
// panel — it answers "where does this go on screen?" and nothing else. No style
// reads a clock, a *ScrollState, or a *PageState; those never leave BuildView.
//
// The frame is single-plane by construction: IsSlow picks exactly one of the
// two render strategies, so a style never mixes the scrolling row plane with
// the page-fade sprite plane.
type PagerSnapshot struct {
	// Lines are the visible text rows for this frame, already selected from the
	// buffer and truncated to VisibleColumns. A style places them as-is; it does
	// not re-slice, re-truncate, or re-order them.
	Lines []string

	// Policy is the current pager configuration, carried so styles can read
	// presentation-affecting fields (e.g. border intent) without a package-level
	// lookup.
	Policy Policy

	// IsSlow selects the render strategy this frame was resolved for:
	//   false → smooth-scroll frame (row plane, uses OffsetY)
	//   true  → page-transition frame (sprite plane, uses Phase/Alpha)
	// It reflects the panel the frame was resolved against, which BuildView
	// decides from the hardware's refresh class.
	IsSlow bool

	// Phase is the page-transition phase for slow frames (idle | fadeOut |
	// fadeIn). Unused by scroll frames.
	Phase PagePhase

	// Alpha is the page-transition opacity for slow frames, 0.0–1.0. Unused by
	// scroll frames.
	Alpha float64

	// OffsetY is the resolved sub-line pixel scroll offset for fast frames. It is
	// negative to scroll content upward, matching ViewData.OffsetY semantics.
	// Unused by slow frames.
	OffsetY int

	// RowHeight and VisibleColumns are the resolved text metrics for this frame,
	// derived once by ComputeLayout so a style never recomputes them from font
	// metrics it would have to re-fetch.
	RowHeight      int
	VisibleColumns int
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
