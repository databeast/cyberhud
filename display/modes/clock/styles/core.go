package styles

// Core layout assembly for clock styles.
//
// Every per-resolution style in this package is declared as a `def` value in
// its styles_WxH.go file: a name, the surface requirements it targets, and a
// Params block of hand-tweakable knobs. The def's Build method assembles the
// shared core layout, adjusted by those Params.
//
// To hand-tweak the clock for a specific display (surface/panel combo), edit
// that resolution's Params — or set Params.BuildFn for a fully bespoke layout
// (see MonoSlow800x480Style's e-ink poster in styles_800x600.go).
// Full tweak guide: README.md in this directory.

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// clockPaddingPct is the percentage inset the clock layouts use.
//
// It is a named constant used in exactly two places — the argument to ctx.Layout
// and the ViewData.PaddingPct handed to the renderer — because those two values
// must agree. The style centres its rows against the content width implied by this
// padding, and the renderer draws into the content origin implied by it. When the
// renderer hardcoded zero and a style padded, text drifted by the inset.
//
// Zero means the clock uses the full panel; its border frame is drawn as a sprite
// at the panel edge rather than as a layout inset.
const clockPaddingPct = 0

// Toggle is a tri-state override: Auto defers to the layout's default rule.
type Toggle int8

const (
	Auto Toggle = iota
	On
	Off
)

// enabled resolves a Toggle against the layout's default for that knob.
func (t Toggle) enabled(auto bool) bool {
	switch t {
	case On:
		return true
	case Off:
		return false
	default:
		return auto
	}
}

// Params holds the hand-tweakable knobs for a single per-resolution style.
// The zero value produces the adaptive core layout unchanged.
type Params struct {
	// Tier overrides the font tier. Empty = auto (adaptive layout picks by
	// panel height via selectTier; fast layout uses TierNormal).
	Tier tiercatalog.Tier

	// Date and Weekday control secondary-row inclusion. Auto derives from
	// policy (adaptive layout) or panel-height thresholds (fast layout).
	Date    Toggle
	Weekday Toggle

	// Fast selects the fast-panel layout: threshold-driven rows, fixed tier,
	// non-static output for continuous refresh.
	Fast bool

	// Color enables the accent color scheme on the fast layout
	// (accent time row, dimmed secondary rows).
	Color bool

	// BuildFn, when set, replaces the core layout entirely for this style.
	BuildFn func(data source.ClockData, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData
}

// styleDef is the declaration form of a clock style: name + surface requirements
// + layout tweaks. It implements style.Style[source.ClockData, source.Policy].
type styleDef struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.ClockData, source.Policy] = styleDef{}

func (d styleDef) Name() string { return d.name }

func (d styleDef) Requirements() style.SurfaceRequirements { return d.reqs }

func (d styleDef) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d styleDef) Build(data source.ClockData, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	if d.p.Fast {
		return fastBuild(data, pol, ctx, d)
	}
	return adaptiveBuild(data, pol, ctx, d)
}

// adaptiveBuild is the core clock layout for policy-driven, mostly-static
// panels (e-ink, slow refresh, and generic fallbacks). Rows follow policy,
// the font tier scales with panel height, and colors engage only when the
// panel capability and policy allow.
func adaptiveBuild(data source.ClockData, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(clockPaddingPct)
	p := pol

	// Resolve primary tier for the time row.
	primaryTier := d.p.Tier
	if primaryTier == "" {
		primaryTier = selectTier(hints.PixelHeight)
	}
	secTier := secondaryTier(primaryTier)

	// Look up catalog entries for each tier.
	//
	// ctx.Entry cannot fail, so the hand-written "if !ok substitute hints.GlyphAdvance"
	// fallback that used to sit here is gone. That fallback was actively harmful: it
	// substituted textlayout's 6x10 defaults, which belong to no registered font, so
	// the layout arithmetic below measured a font the renderer would never draw with.
	// Entry always names a real face, and the renderer draws each row with exactly
	// the face named here.
	primaryEntry := ctx.Entry(primaryTier)
	secEntry := ctx.Entry(secTier)

	// Build row list with per-row tier assignments.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	rows := []rowSpec{
		{text: data.Time, tier: primaryTier, entry: primaryEntry},
	}
	if d.p.Date.enabled(p.DateFormat != "none") && data.Date != "" {
		rows = append(rows, rowSpec{text: data.Date, tier: secTier, entry: secEntry})
	}
	if d.p.Weekday.enabled(p.ShowWeekday) && data.Weekday != "" {
		rows = append(rows, rowSpec{text: data.Weekday, tier: secTier, entry: secEntry})
	}

	// Compute per-row centering using each row's own glyph advance.
	items := make([]string, len(rows))
	tiers := make([]tiercatalog.Tier, len(rows))
	rowInputs := make([]layout.RowInput, len(rows))
	for i, r := range rows {
		items[i] = r.text
		tiers[i] = r.tier
		rowInputs[i] = layout.RowInput{
			TextLen:      len(r.text),
			GlyphAdvance: r.entry.GlyphAdvance,
			RowHeight:    r.entry.RowHeight,
		}
	}

	result := bridge.LayoutRows(rowInputs)

	var colors []color.Color
	if ctx.Cap() >= style.ColorSlow && p.FGColor != "none" {
		accent := resolveFGColor(p.FGColor)
		colors = make([]color.Color, len(items))
		colors[0] = accent
		for i := 1; i < len(items); i++ {
			colors[i] = sharedcolor.Dim(accent)
		}
	}

	// Hand the renderer the entire layout solution. Every field of RowsResult must
	// travel: OffsetY was centred against these exact RowHeights and Spacing, and
	// VisibleCount is the fit decision FitRows made. A renderer that re-derives any
	// of them disagrees with OffsetY. See the layout contract on style.ViewData.
	vd := style.ViewData{
		Items:        items,
		Tiers:        tiers,
		LineOffsets:  result.Offsets,
		OffsetY:      result.OffsetY,
		RowHeights:   result.RowHeights,
		Spacing:      result.Spacing,
		VisibleCount: result.VisibleCount,
		PaddingPct:   clockPaddingPct,
		Colors:       colors,
		Static:       true,
	}
	guardEmptyItems(&vd)

	reqs := style.SurfaceRequirements{
		MinWidth:   hints.PixelWidth,
		MinHeight:  hints.PixelHeight,
		Capability: ctx.Cap(),
	}
	appendClockWidgets(&vd, bridge, hints, p, data.Now, p.FGColor != "none", reqs)

	return vd
}

// fastBuild is the core clock layout for fast-refresh TFT/OLED panels.
// Secondary rows engage by panel-height threshold (date at ≥64px, weekday at
// ≥96px), the font tier is fixed, and output is non-static so seconds and
// widget animations refresh continuously.
func fastBuild(data source.ClockData, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(clockPaddingPct)
	p := pol

	// Resolve primary tier for the time row.
	primaryTier := d.p.Tier
	if primaryTier == "" {
		primaryTier = tiercatalog.TierNormal
	}
	secTier := secondaryTier(primaryTier)

	// Look up catalog entries for each tier. See the equivalent note in
	// adaptiveBuild: ctx.Entry is infallible and names a real face, replacing the
	// per-field zero-patching that used to substitute bridge defaults.
	primaryEntry := ctx.Entry(primaryTier)
	secEntry := ctx.Entry(secTier)

	// Build row list with per-row tier assignments.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	rows := []rowSpec{
		{text: data.Time, tier: primaryTier, entry: primaryEntry},
	}
	if d.p.Date.enabled(d.reqs.MinHeight >= 64) && data.Date != "" {
		rows = append(rows, rowSpec{text: data.Date, tier: secTier, entry: secEntry})
	}
	if d.p.Weekday.enabled(d.reqs.MinHeight >= 96) && data.Weekday != "" {
		rows = append(rows, rowSpec{text: data.Weekday, tier: secTier, entry: secEntry})
	}

	// Compute per-row centering and heights using each row's tier entry.
	items := make([]string, len(rows))
	tiers := make([]tiercatalog.Tier, len(rows))
	rowInputs := make([]layout.RowInput, len(rows))
	for i, r := range rows {
		items[i] = r.text
		tiers[i] = r.tier
		rowInputs[i] = layout.RowInput{
			TextLen:      len(r.text),
			GlyphAdvance: r.entry.GlyphAdvance,
			RowHeight:    r.entry.RowHeight,
		}
	}

	result := bridge.LayoutRows(rowInputs)

	var colors []color.Color
	if d.p.Color {
		accent := resolveFGColor(p.FGColor)
		colors = make([]color.Color, len(items))
		colors[0] = accent
		for i := 1; i < len(items); i++ {
			colors[i] = sharedcolor.Dim(accent)
		}
	}

	// See the note in adaptiveBuild: the complete RowsResult travels to the renderer.
	vd := style.ViewData{
		Items:        items,
		Tiers:        tiers,
		LineOffsets:  result.Offsets,
		OffsetY:      result.OffsetY,
		RowHeights:   result.RowHeights,
		Spacing:      result.Spacing,
		VisibleCount: result.VisibleCount,
		PaddingPct:   clockPaddingPct,
		Colors:       colors,
	}

	appendClockWidgets(&vd, bridge, hints, p, data.Now, p.FGColor != "none", d.reqs)

	return vd
}
