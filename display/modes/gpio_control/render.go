package gpio_control

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

func (i *instance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(Signature(i.pins()), PolicySnapshot().Fingerprint())
}

// BuildView constructs the full GPIO Control view result including Colors,
// Sprites, and TextHints-based truncation. It uses registry-based dispatch
// instead of conditional branching on style name strings.
//
// Font resolution is handled by the registry wrapper's tier-based resolution
// (configured as terminus/TierNormal for gpio-control). BuildView no longer
// selects fonts directly — it uses hints.Face for sprite widget faces and
// relies on the registry wrapper for the Items path.
func BuildView(pins []gpiomgr.PinState, hints textlayout.TextHints, cursor ...int) style.ViewData {
	pol := PolicySnapshot()
	pol = normalizePolicy(pol)

	// Obtain a catalog-validated face via hints.Face for sprite rendering.
	// The tier catalog guarantees GlyphAdvance ≤ maxAdvance.
	var spriteFace font.Face
	if hints.Catalog.PixelWidth() > 0 {
		spriteFace = hints.Face("terminus", tiercatalog.TierNormal)
	}

	// Rebind glyph metrics to match the tier-resolved font.
	if spriteFace != nil {
		hints = textlayout.WithFont(hints, spriteFace)
	}

	// Build snapshot for registry-based dispatch.
	cur := 0
	if len(cursor) > 0 {
		cur = cursor[0]
	}
	snap := source.Data{
		Pins:   pins,
		Cursor: cur,
		TopRow: 0,
	}

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(gpioControlRegistry, hints, "gpio-control", pol.Style)

	// Call Build on the looked-up style.
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, pol, ctx)

	// Title/Hint/Static/StyleReport ONLY — no sprites, no geometry, no layout.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	return svd
}

// Signature returns a change token for UI refresh.
func Signature(pins []gpiomgr.PinState) string {
	sig := "gpio-control:"
	for _, p := range pins {
		lv := 0
		if p.Level {
			lv = 1
		}
		sig += fmt.Sprintf("%d:%s:%d;", p.Number, p.Mode, lv)
	}
	return sig
}
