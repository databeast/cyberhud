package styles

// Placement helpers for the pager core layout.
//
// These answer "where does this go on screen?" for an already-resolved frame.
// Nothing here reads a clock or advances animation state — the PagerSnapshot
// arrives with its lines, scroll offset, and page phase/alpha already settled.
//
// The pager shows text and nothing else: no border, no decorative chrome. A
// frame is committed to exactly one plane. Fast (smooth-scroll) frames use the
// renderer's row plane: Items plus a sub-line OffsetY. Slow (page-transition)
// frames rasterize the page to a single sprite so the fade alpha can be baked in,
// which the row plane cannot express. The two are never emitted together.

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
)

// pagerPaddingPct is the inset every pager style lays out against: a small, even
// margin that keeps text off the panel edge. It is passed both to ctx.Layout (so
// the style centres against the inset content width) and to ViewData.PaddingPct
// (so the renderer draws into the same content origin).
const pagerPaddingPct = 2

// placeScroll places a smooth-scroll frame on the row plane.
//
// The lines are already the visible window, already truncated to the frame's
// column budget by BuildView. Placement centres each line horizontally and hands
// the renderer the resolved sub-line OffsetY so the block scrolls with sub-pixel
// vertical motion. Output is non-static: the scroll offset changes every frame.
func placeScroll(snap source.PagerSnapshot, _ source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	bridge := ctx.Layout(pagerPaddingPct)
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return tooSmallView()
	}

	m := ctx.Metrics(tiercatalog.TierNormal)

	if len(snap.Lines) == 0 {
		return centeredEmptyView(ctx, bridge, false)
	}

	lines := snap.Lines
	offsets := make([]style.TTYX, len(lines))
	tiers := make([]tiercatalog.Tier, len(lines))
	for i, line := range lines {
		offsets[i] = bridge.CenterX(runeLen(line), m)
		tiers[i] = tiercatalog.TierNormal
	}

	return style.ViewData{
		Items:       lines,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     style.TTYY(snap.OffsetY), // pager scroll offset is TTY-block space
		PaddingPct:  pagerPaddingPct,
		Static:      false,
	}
}

// placePage places a page-transition frame on the sprite plane.
//
// The page text is rasterized to one image with the frame's fade alpha baked in,
// because opacity is a per-pixel property the row plane cannot carry. Output is
// static: a settled page holds until the next transition, and mid-fade frames are
// driven by BuildView advancing Alpha, not by the renderer.
func placePage(snap source.PagerSnapshot, _ source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(pagerPaddingPct)
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return tooSmallView()
	}

	page := renderPageImage(snap, bridge, hints)
	if page == nil {
		return centeredEmptyView(ctx, bridge, true)
	}

	return style.ViewData{
		Static:     true,
		PaddingPct: pagerPaddingPct,
		Sprites: []widgets.Sprite{{
			Image:    page,
			Position: image.Point{},
			Label:    "pager-page",
		}},
	}
}

// renderPageImage rasterizes the resolved page lines into a panel-sized image
// with the frame's fade alpha applied to the text color. Returns nil when there
// is nothing to draw.
func renderPageImage(snap source.PagerSnapshot, bridge layout.LayoutCalculator, hints textlayout.TextHints) *image.RGBA {
	if len(snap.Lines) == 0 {
		return nil
	}

	w := hints.PixelWidth
	h := hints.PixelHeight
	if w <= 0 || h <= 0 {
		return nil
	}

	face := resolvePageFace(snap.RowHeight)
	if face == nil {
		return nil
	}
	metrics := face.Metrics()
	rowHeight := snap.RowHeight
	if rowHeight <= 0 {
		rowHeight = metrics.RowHeight
	}
	if rowHeight <= 0 {
		return nil
	}

	alpha := snap.Alpha
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	fg := color.RGBA{R: 255, G: 255, B: 255, A: uint8(alpha * 255)}

	ox, oy := bridge.ContentOrigin()
	contentW := bridge.AvailableContentWidth()
	contentH := bridge.AvailableContentHeight()
	maxRows := contentH / rowHeight
	if maxRows <= 0 {
		maxRows = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for row := 0; row < len(snap.Lines) && row < maxRows; row++ {
		y := oy + row*rowHeight
		yOffset := 0
		if rowHeight > metrics.GlyphHeight {
			yOffset = (rowHeight - metrics.GlyphHeight) / 2
		}
		font.DrawText(img, face, snap.Lines[row], ox, y+yOffset, fg, ox+contentW)
	}
	return img
}

// centeredEmptyView places a single "(no content)" message, centred. Used by
// both planes when the frame carries no lines.
func centeredEmptyView(ctx style.StyleContext, bridge layout.LayoutCalculator, static bool) style.ViewData {
	const msg = "(no content)"
	m := ctx.Metrics(tiercatalog.TierNormal)

	text := msg
	if max := bridge.MaxChars(m); max > 0 {
		text = textlayout.Truncate(text, max)
	}
	offsetX := bridge.CenterX(runeLen(text), m)
	blockH := m.RowHeight()
	offsetY := style.TTYY((bridge.AvailableContentHeight() - blockH) / 2)
	if offsetY < 0 {
		offsetY = 0
	}

	return style.ViewData{
		Items:       []string{text},
		Tiers:       []tiercatalog.Tier{tiercatalog.TierNormal},
		LineOffsets: []style.TTYX{offsetX},
		OffsetY:     offsetY,
		PaddingPct:  pagerPaddingPct,
		Static:      static,
	}
}

// tooSmallView is the degenerate placement for a panel with no usable content
// area.
func tooSmallView() style.ViewData {
	return style.ViewData{Items: []string{"(too small)"}, Static: true, PaddingPct: pagerPaddingPct}
}

// resolvePageFace picks the font face for the page sprite, matching the row
// height the frame was resolved against.
func resolvePageFace(rowHeight int) font.Face {
	if rowHeight <= 0 {
		return font.Default()
	}
	return font.ByHeight(rowHeight)
}

// runeLen counts runes so centering is correct for multi-byte content.
func runeLen(s string) int {
	n := 0
	for range s {
		n++
	}
	return n
}
