package styles

import (
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a single thermal style.
type Params struct {
	// BuildFn replaces the core layout entirely for bespoke thermal layouts.
	BuildFn func(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

// def is the declaration form of a thermal style: name + surface requirements
// + layout tweaks. It implements style.Style[source.ThermalSnapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.ThermalSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snapshot, pol, ctx, d)
	}
	return adaptiveBuild(snapshot, pol, ctx, d)
}

// adaptiveBuild is the fallback layout for skeleton styles with nil BuildFn.
// It selects the best core layout based on panel dimensions:
//   - Panels ≥ 240px in both dimensions: buildOverview (multi-zone dashboard)
//   - Panels ≥ 128px wide and < 240px tall or wide: buildMonoOLEDCompact (compact dashboard)
//   - Very small panels (< 128px wide): buildMinimal (single temperature)
func adaptiveBuild(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}, Static: true}
	}

	hints := ctx.Hints()
	w := hints.PixelWidth
	h := hints.PixelHeight

	var result style.ViewData
	switch {
	case w >= 240 && h >= 240:
		result = buildOverview(snapshot, pol, ctx)
	case w >= 128:
		result = buildMonoOLEDCompact(snapshot, pol, ctx)
	default:
		result = buildMinimal(snapshot, pol, ctx)
	}

	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"thermal"}
	}
	result.Static = true
	return result
}
