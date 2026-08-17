package styles

// Core layout assembly for ticker styles.
//
// Each concrete ticker style is declared as a def value in its styles_*.go file:
// a stable name, surface requirements, and an optional Params block. The shared
// Build method below covers the fast and slow ticker bodies; Params.BuildFn is
// the escape hatch for bespoke layouts.

import (
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
)

// Params holds hand-tweakable knobs for a declared ticker style.
// The zero value produces the fast/shared ticker layout.
type Params struct {
	// Slow selects the slow-refresh shared layout: disable autoscroll,
	// ignore strip sprites, and mark the ViewData static.
	Slow bool

	// BuildFn, when set, replaces the shared layout entirely for this style.
	BuildFn func(source.TickerSnapshot, source.Policy, style.StyleContext, def) style.ViewData
}

// def is the declaration form of a ticker style: name + surface requirements +
// layout tweaks. It implements style.Style[source.TickerSnapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.TickerSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snapshot source.TickerSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snapshot, pol, ctx, d)
	}
	if d.p.Slow {
		return slowBuild(snapshot, pol, ctx, d)
	}
	return fastBuild(snapshot, pol, ctx, d)
}

func fastBuild(snapshot source.TickerSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := snapshot.Hints
	borderInset := 0
	effective := source.EffectivePolicy(pol, hints)

	// --- Conditional border ---
	var borderSprites []widgets.Sprite
	if pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16 {
		borderSprites = buildBorderSprites(hints, effective)
	}

	var formatted []source.FormattedLine
	if effective.Direction == "horizontal" && effective.AutoScrollMS > 0 {
		formatted = formatNonScrollingLines(snapshot.Directives, hints, pol, borderInset)
	} else {
		formatted = PartitionScroll(snapshot.Directives, hints, pol, borderInset, snapshot.ScrollOffset)
	}

	items, tiers := itemsAndTiers(formatted)
	offsetY := source.ComputeOffsetY(formatted, hints, effective)

	// --- Conditional glow (ColorFast panels only) ---
	var glowSprites []widgets.Sprite
	if pol.ShowGlow && isColorFast(hints) {
		accent := sharedcolor.ResolveAccent(effective.Accent)
		glowSprites = append(glowSprites, buildGlowBackground(hints.PixelWidth, hints.PixelHeight, accent))
		glowSprites = append(glowSprites, buildGlowSprites(formatted, hints, offsetY, accent)...)
	}

	allSprites := make([]widgets.Sprite, 0, len(borderSprites)+len(glowSprites)+len(snapshot.StripSprites))
	allSprites = append(allSprites, borderSprites...)
	allSprites = append(allSprites, glowSprites...)
	allSprites = append(allSprites, snapshot.StripSprites...)

	return style.ViewData{
		Items:   items,
		Tiers:   tiers,
		Sprites: allSprites,
		OffsetY: offsetY,
	}
}

func slowBuild(snapshot source.TickerSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := snapshot.Hints
	borderInset := 0
	effective := source.EffectivePolicy(pol, hints)
	effective.AutoScrollMS = 0

	// --- Conditional border ---
	var borderSprites []widgets.Sprite
	if pol.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16 {
		borderSprites = buildBorderSprites(hints, effective)
	}

	formatted := PartitionScroll(snapshot.Directives, hints, effective, borderInset, 0)
	items, tiers := itemsAndTiers(formatted)
	offsetY := source.ComputeOffsetY(formatted, hints, effective)

	return style.ViewData{
		Items:   items,
		Tiers:   tiers,
		Sprites: borderSprites,
		OffsetY: offsetY,
		Static:  true,
	}
}

func itemsAndTiers(formatted []source.FormattedLine) ([]string, []tiercatalog.Tier) {
	items := make([]string, len(formatted))
	tiers := make([]tiercatalog.Tier, len(formatted))
	for i, fl := range formatted {
		items[i] = fl.Text
		tiers[i] = fl.Tier
	}
	return items, tiers
}

// BorderInset returns the pixel inset for the current style.
// After the layout-padding-refactor, border is decorative only and does not
// affect content layout — always returns 0.
func BorderInset() int {
	return 0
}
