package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// accentName is the schedule mode's house accent: amber reads as "now / live"
// without the forbidden cyan, and dims cleanly for the secondary rows.
const accentName = "amber"

// adaptiveBuild is the single schedule layout. It resolves the view model, sizes
// the title row to the panel, truncates each row to the content width of its own
// font, and hands the renderer the complete layout solution.
func adaptiveBuild(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(schedulePaddingPct)
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		vd := style.ViewData{Items: []string{"(too small)"}, Cursor: -1}
		style.GuardEmptyItems(&vd)
		return vd
	}

	v := resolveView(snap, pol)

	// Title uses the largest tier that suits the panel height; the secondary and
	// upcoming rows step one tier down so the headline dominates. Each row is then
	// sized DOWN until its whole text fits the width — never truncated. A clipped
	// "Design re..." tells a glance nothing, so a row that cannot fit even at the
	// smallest tier is dropped entirely (the title excepted: it is the one row that
	// must always show, so it falls back to the smallest tier and hard-clips at the
	// panel edge rather than vanishing).
	cw := bridge.AvailableContentWidth()
	titleCeil := style.TierForHeight(hints.PixelHeight)
	subCeil := style.SecondaryTier(titleCeil)

	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}
	rows := make([]rowSpec, 0, 3)

	// Title: fit down from its ceiling; if nothing fits, keep the smallest tier the
	// catalog offers and let it hard-clip at the edge — the headline is
	// non-negotiable and must never vanish.
	titleTier, ok := ctx.FitTier(v.Title, cw, titleCeil)
	if !ok {
		if all := ctx.FontCatalog().Tiers(); len(all) > 0 {
			titleTier = all[0] // ascending order: smallest first
		}
	}
	rows = append(rows, rowSpec{text: v.Title, tier: titleTier, entry: ctx.Entry(titleTier)})

	// Secondary rows: fit down from subCeil; drop the row if even the smallest
	// tier cannot render it whole.
	addFitRow := func(text string) {
		if text == "" {
			return
		}
		if tier, ok := ctx.FitTier(text, cw, subCeil); ok {
			rows = append(rows, rowSpec{text: text, tier: tier, entry: ctx.Entry(tier)})
		}
	}
	addFitRow(v.Subtitle)
	addFitRow(v.Upcoming)

	items := make([]string, len(rows))
	tiers := make([]tiercatalog.Tier, len(rows))
	rowInputs := make([]layout.RowInput, len(rows))
	for i, r := range rows {
		items[i] = r.text
		tiers[i] = r.tier
		rowInputs[i] = layout.RowInput{
			TextLen: len([]rune(r.text)),
			Font:    layout.NewMetrics(r.entry.GlyphAdvance, r.entry.GlyphWidth, r.entry.RowHeight, r.entry.FontID),
		}
	}

	result := bridge.LayoutRows(rowInputs)

	var colors []color.Color
	if d.p.Color && ctx.Cap() >= style.ColorSlow {
		accent := sharedcolor.ResolveAccent(accentName)
		if !v.Accent {
			// Idle / error states read in a restrained white rather than the live
			// accent, so the accent means "something is on right now".
			accent = sharedcolor.ResolveAccent("white")
		}
		colors = make([]color.Color, len(items))
		colors[0] = accent
		for i := 1; i < len(items); i++ {
			colors[i] = sharedcolor.Dim(accent)
		}
	}

	vd := style.ViewData{
		Items:        items,
		Tiers:        tiers,
		Colors:       colors,
		LineOffsets:  result.Offsets,
		OffsetY:      result.OffsetY,
		RowHeights:   result.RowHeights,
		Spacing:      result.Spacing,
		VisibleCount: result.VisibleCount,
		PaddingPct:   schedulePaddingPct,
		Cursor:       -1, // status view: never a navigable list (see ViewData.IsNavigable).
	}
	style.GuardEmptyItems(&vd)
	return vd
}
