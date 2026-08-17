// Package style defines the shared contract for display mode visual variants.
// This package has zero dependencies on any mode-specific packages.
package style

import (
	"image/color"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
)

// Style is the interface that all display mode style implementations satisfy.
// S is the mode-specific snapshot type, P is mode policy for the style to obey
type Style[C any, P catalog.ConfigPolicy] interface {
	// Name returns the style's identifier: lowercase ASCII [a-z0-9-], 1-32 chars.
	Name() string

	// Requirements returns the style's panel surface needs.
	Requirements() SurfaceRequirements

	// Supports evaluates how well this style fits the described panel.
	Supports(hints textlayout.TextHints) Fitness

	// Build produces the rendering output for this style.
	// snapshot is the mode-specific typed snapshot; ctx provides layout geometry and capability flags.
	Build(snapshot C, pol P, ctx StyleContext) ViewData
}

// StyleReport carries metadata about which style was resolved during BuildView.
// Zero value means "no style resolution occurred" (mode has no style registry).
type StyleReport struct {
	Name   string // Resolved style name (e.g., "color-240x320")
	Reason string // "configured" or "default"
}

// ViewData is the rendering output produced by a Style's Build method.
//
// # The layout contract
//
// A style owns layout. The renderer owns drawing. ViewData is the whole of the
// interface between them, so any geometry the renderer needs must be carried here.
//
// This is not a stylistic preference; it is the lesson of a specific class of bug
// that recurred in this codebase. The style solves row layout once, against the
// real per-row font metrics, via [layout.LayoutCalculator.LayoutRows]. Whenever a
// field of that solution was omitted from ViewData, the renderer re-derived it from
// whatever it had to hand — a single global row height, its own LayoutBridge, its
// own padding assumption — and that derivation disagreed with the OffsetY the style
// had already centred the block with. The symptoms were text of the wrong size,
// rows touching each other, and blocks that were neither centred nor clipped
// correctly.
//
// So: when adding a field to RowsResult, add it here and consume it in the
// renderer. When writing a renderer, read geometry from ViewData; do not
// recompute it.
type ViewData struct {
	Items  []string //row-based text data
	Colors []color.Color

	// Tiers is per-row tier intent, parallel to Items. Nil defaults to the mode's
	// configured tier.
	//
	// Tiers is intent; FontIDs is the resolution of that intent. Styles set Tiers
	// and leave FontIDs nil; the renderer resolves one into the other. Declaring a
	// tier is therefore sufficient to get that font — a style does not need, and
	// should not have, its own font resolution step.
	Tiers []tiercatalog.Tier

	// FontIDs are the resolved per-row font IDs.
	//
	// Normally left nil by styles and populated by the renderer from Tiers against
	// the region's catalog. Set it directly only to override that resolution; the
	// renderer will not overwrite a non-nil value.
	//
	// Historical note: nothing performed the Tiers to FontIDs conversion for a
	// period, because it was expected of each mode individually and the only
	// reference implementation of that step (validators/resolvefont.go) is
	// commented out. Styles set Tiers, the renderer read FontIDs, and the two never
	// met — so every row rendered in the surface's default 5x8 face regardless of
	// the tier the catalog had correctly chosen. Keeping the resolution in the
	// renderer means it cannot be forgotten per-mode.
	FontIDs []string

	LineOffsets []int // Per-line horizontal pixel offset for marquee; nil if inactive.
	OffsetY     int   // Vertical pixel offset for content block (used for centering).

	// RowHeights is the per-row pixel height, parallel to Items. Nil means the
	// renderer derives each row's height from that row's font metrics.
	//
	// Supply this whenever rows use different tiers. The style's OffsetY was
	// centred against these exact heights, so letting the renderer guess them
	// puts the block off centre.
	RowHeights []int

	// Spacing is the inter-row gap in pixels, as computed by the style's row
	// fitting. Zero means rows are drawn flush, which is a legitimate choice for a
	// single-row view and wrong for a multi-row one — an omitted Spacing is why
	// stacked rows once rendered touching.
	Spacing int

	// VisibleCount is how many of Items actually fit the region, from
	// [layout.RowsResult.VisibleCount]. Zero means "unspecified", in which case the
	// renderer falls back to its own row budget.
	//
	// Set this for tiered views. The renderer's fallback divides region height by a
	// single row height, which cannot be correct when rows use different fonts.
	VisibleCount int

	// PaddingPct is the percentage inset the style laid out against, matching the
	// argument it passed to [StyleContext.Layout].
	//
	// The renderer builds its own LayoutCalculator to find the content origin and
	// width it draws into. If that calculator's padding differs from the style's,
	// the style's horizontal offsets were measured against a different content
	// width than the renderer draws into, and text drifts. Carrying the value keeps
	// the two in agreement. Zero is both the default and the common case.
	PaddingPct int

	Cursor      int
	TopRow      int
	Static      bool
	Sprites     []widgets.Sprite
	StyleReport StyleReport // Style resolution metadata; zero value = no style dispatch.
}
