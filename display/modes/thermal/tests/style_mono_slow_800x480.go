//go:build legacy_internal_tests
// +build legacy_internal_tests

// TODO: Re-enable after these legacy internal tests are migrated to exported mode/source APIs.

package tests

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/thermal"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

// MonoSlow800x480Style implements the 9-step polish pipeline for the thermal mode
// on a MonoSlow 800×480 e-ink panel. It displays the hottest zone as a header row
// at TierLarge, with remaining zones as secondary rows at TierSmall.
type MonoSlow800x480Style struct{}

// Name returns the style identifier.
func (s MonoSlow800x480Style) Name() string { return "mono-slow-800x480" }

// Requirements returns the surface requirements for this style.
func (s MonoSlow800x480Style) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow}
}

// Supports evaluates how well this style fits the described panel.
func (s MonoSlow800x480Style) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

// Build produces the rendering output using the 9-step polish pipeline.
func (s MonoSlow800x480Style) Build(snap thermal.ThermalSnapshot, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()

	// Step 1 (Guard): LayoutBridge with PaddingPct:2, early return for zero dims.
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2 (Borderframe): Theme "sharp", e-ink safe (all animation/glow/gradient zero).
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3 (Tier resolution): TierLarge for header, TierSmall for secondary zones.
	catalog := ctx.FontCatalog()

	largeEntry, ok := catalog.Get(tiercatalog.TierLarge)
	if !ok {
		largeEntry = tiercatalog.Entry{GlyphAdvance: bridge.GlyphAdvance(), RowHeight: bridge.RowHeight()}
	}
	smallEntry, ok := catalog.Get(tiercatalog.TierSmall)
	if !ok {
		smallEntry = tiercatalog.Entry{GlyphAdvance: bridge.GlyphAdvance(), RowHeight: bridge.RowHeight()}
	}

	// Step 4 (Row assembly): hottest zone as header, remaining zones as secondary.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	pol := thermal.GetPolicy()

	var header []rowSpec
	var secondary []rowSpec

	if len(snap.Zones) == 0 {
		header = []rowSpec{{text: "(no zones)", tier: tiercatalog.TierLarge, entry: largeEntry}}
	} else {
		// Find the hottest zone.
		hottest := snap.Zones[0]
		for _, z := range snap.Zones[1:] {
			if z.TempC > hottest.TempC {
				hottest = z
			}
		}

		header = []rowSpec{
			{text: fmt.Sprintf("%s %s", hottest.Label, formatTemp(hottest.TempC, pol.Unit)), tier: tiercatalog.TierLarge, entry: largeEntry},
		}

		// Secondary: remaining zones (skip the hottest by ZoneID).
		for _, z := range snap.Zones {
			if z.ZoneID == hottest.ZoneID {
				continue
			}
			secondary = append(secondary, rowSpec{
				text:  fmt.Sprintf("%s %s", z.Label, formatTemp(z.TempC, pol.Unit)),
				tier:  tiercatalog.TierSmall,
				entry: smallEntry,
			})
		}
	}

	// Step 5 (GroupGap): half header tier row height, minimum 8px.
	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	// Step 6 (Adaptive fitting): drop secondary rows from end when height insufficient.
	allRows := append(header, secondary...)
	includeSecondary := len(secondary) > 0

	totalHeight := 0
	for _, r := range header {
		totalHeight += r.entry.RowHeight
	}
	for _, r := range secondary {
		totalHeight += r.entry.RowHeight
	}
	if includeSecondary {
		totalHeight += groupGap
	}

	visibleRows := allRows
	if totalHeight > bridge.AvailableContentHeight() {
		// Drop secondary rows from end until fit.
		includeSecondary = false
		visibleRows = header
		for i := len(secondary); i > 0; i-- {
			candidate := append(header, secondary[:i]...)
			h := 0
			for _, r := range candidate {
				h += r.entry.RowHeight
			}
			h += groupGap
			if h <= bridge.AvailableContentHeight() {
				visibleRows = candidate
				includeSecondary = true
				break
			}
		}
	}

	// Step 7 (Truncation): per-tier maxChars.
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 8 (Centering): CenterXWith per-tier, vertical centering.
	offsets := make([]int, len(visibleRows))
	for i, r := range visibleRows {
		offsets[i] = bridge.CenterXWith(len([]rune(r.text)), r.entry.GlyphAdvance)
	}

	var blockHeight int
	for _, r := range visibleRows {
		blockHeight += r.entry.RowHeight
	}
	if includeSecondary && len(visibleRows) > len(header) {
		blockHeight += groupGap
	}
	offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Step 9 (Emit): ViewData with Items, Tiers, LineOffsets, OffsetY, Static:true, Sprites.
	items := make([]string, len(visibleRows))
	tiers := make([]tiercatalog.Tier, len(visibleRows))
	for i, r := range visibleRows {
		items[i] = r.text
		tiers[i] = r.tier
	}

	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
		Sprites:     sprites,
	}
}
func formatTemp(celsius float64, unit string) string {
	if unit == "F" {
		return fmt.Sprintf("%.1f°F", celsius*9.0/5.0+32.0)
	}
	return fmt.Sprintf("%.1f°C", celsius)
}
