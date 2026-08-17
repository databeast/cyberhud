package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

// MonoSlow800x480Style implements the polished 9-step pipeline for pager
// on monochrome slow-refresh 800x480 e-ink panels. Uses flat-tier mode
// (all TierNormal, no GroupGap).
var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: monoSlow800x480},
}

func monoSlow800x480(snap source.PagerSnapshot, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	// Step 1 (Guard): LayoutBridge with PaddingPct:2, early return for zero dims.
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2 (Borderframe): Theme "sharp", e-ink safe.
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3 (Tier resolution): TierNormal for content rows, TierLarge for empty state.
	normalEntry := ctx.Entry(tiercatalog.TierNormal)
	largeEntry := ctx.Entry(tiercatalog.TierLarge)

	// Step 4 (Row assembly): ignore snap.Phase and snap.Alpha (static rendering).
	// Empty state: single "(no content)" row at TierLarge.
	if len(snap.Lines) == 0 {
		emptyText := "(no content)"

		// Truncation
		ga := largeEntry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				emptyText = textlayout.Truncate(emptyText, maxChars)
			}
		}

		// Horizontal centering
		offsetX := bridge.CenterXWith(len([]rune(emptyText)), ga)

		// Vertical centering
		blockHeight := largeEntry.RowHeight
		offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
		if offsetY < 0 {
			offsetY = 0
		}

		return style.ViewData{
			Items:       []string{emptyText},
			Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
			LineOffsets: []int{offsetX},
			OffsetY:     offsetY,
			Static:      true,
			Sprites:     sprites,
		}
	}

	// All lines at TierNormal.
	lines := make([]string, len(snap.Lines))
	copy(lines, snap.Lines)

	// Step 5 (Adaptive fitting — flat): all lines at TierNormal, no GroupGap.
	rowHeight := normalEntry.RowHeight
	if rowHeight <= 0 {
		rowHeight = 10
	}
	maxRows := bridge.AvailableContentHeight() / rowHeight
	if maxRows > 0 && len(lines) > maxRows {
		lines = lines[:maxRows]
	}

	// Step 6 (Truncation): uniform TierNormal glyph advance for all rows.
	ga := normalEntry.GlyphAdvance
	if ga > 0 {
		maxChars := bridge.AvailableContentWidth() / ga
		if maxChars > 0 {
			for i := range lines {
				lines[i] = textlayout.Truncate(lines[i], maxChars)
			}
		}
	}

	// Step 7 (Horizontal centering): CenterXWith per line.
	offsets := make([]int, len(lines))
	for i, line := range lines {
		offsets[i] = bridge.CenterXWith(len([]rune(line)), ga)
	}

	// Step 8 (Vertical centering): center the block of lines.
	blockHeight := len(lines) * rowHeight
	offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Step 9 (Emit): build tiers slice (all TierNormal).
	tiers := make([]tiercatalog.Tier, len(lines))
	for i := range tiers {
		tiers[i] = tiercatalog.TierNormal
	}

	return style.ViewData{
		Items:       lines,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
		Sprites:     sprites,
	}
}
