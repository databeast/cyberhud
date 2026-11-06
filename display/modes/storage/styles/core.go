package styles

import (
	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Params holds the hand-tweakable knobs for a single storage style.
type Params struct {
	// BuildFn replaces the core layout entirely for a bespoke storage layout.
	// When nil, the style uses adaptiveBuild, which picks a layout by panel size.
	BuildFn func(snapshot source.StorageSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

// def is the declaration form of a storage style: name + surface requirements
// + layout tweaks. It implements style.Style[source.StorageSnapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.StorageSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snapshot source.StorageSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snapshot, pol, ctx, d)
	}
	return adaptiveBuild(snapshot, pol, ctx, d)
}

// adaptiveBuild is the fallback layout for skeleton styles (nil BuildFn),
// selected by fitness when the user has not picked a specific variant. It
// chooses a layout by panel dimensions:
//   - Very small panels (width < 96 or height < 48): the single-focus readout
//     of the most-used filesystem.
//   - Everything else: the multi-mount bars overview.
func adaptiveBuild(snapshot source.StorageSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	if hints.PixelWidth < 96 || hints.PixelHeight < 48 {
		return buildFocusStyle(snapshot, pol, ctx, d)
	}
	return buildBarsStyle(snapshot, pol, ctx, d)
}
