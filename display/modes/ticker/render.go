package ticker

import (
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the ticker mode view data with border frame sprites and per-line font selection.
// Font resolution is handled entirely through the tier catalog via hints.Face.
func BuildView(hints textlayout.TextHints) style.ViewData {
	p := PolicySnapshot()

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(tickerRegistry, hints, "ticker", p.Style)

	// Get scroll offset first (acquires lock), then feed snapshot (acquires lock separately).
	scrollOffset := source.CheckAdvance(hints)
	directives := source.FeedSnapshot()

	// Pre-render strip sprites (acquires read lock on feedState).
	stripSprites := source.RenderStripSprites()

	snap := source.TickerSnapshot{
		Directives:   directives,
		Policy:       p,
		ScrollOffset: scrollOffset,
		StripSprites: stripSprites,
		Hints:        hints,
	}

	// Construct StyleContext for the style boundary.
	ctx := style.NewStyleContext(hints)

	// All styles (including BorderedStyle) now use the Compositor pattern
	// internally; no special-case injection needed.
	svd := s.Build(snap, p, ctx)

	// Report style resolution metadata to the registry layer.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	// Guard: ensure Items always contains at least one non-empty string
	// UNLESS there are sprites being rendered (horizontal scrolling).
	if source.AllEmptyItems(svd.Items) && len(svd.Sprites) == 0 {
		svd.Items = []string{"(ticker idle)"}
	}

	return svd
}

// Signature returns a stable signature for feed change detection.
func Signature() string {
	p := source.PolicySnapshot()

	var parts []string
	for _, d := range source.FeedSnapshot() {
		part := d.Text
		if d.Font != "" {
			part += "@" + d.Font
		}
		if d.LineMode != "" {
			part += "#" + d.LineMode
		}
		if d.Scaling != "" {
			part += "^" + d.Scaling
		}
		if d.Scroll != "" {
			part += "~" + d.Scroll
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, "|") +
		";style=" + p.Style +
		";font=" + p.Font +
		";font_tier=" + p.FontTier +
		";line_mode=" + p.LineMode +
		";direction=" + p.Direction +
		";auto_scroll_ms=" + strconv.Itoa(p.AutoScrollMS) +
		";accent=" + p.Accent
}

func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey("ticker", Signature())
}
