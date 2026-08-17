package tiercatalog

import "github.com/databeast/cyberhud/display/surface/fonts"

// Infallible catalog access.
//
// # Why this file exists
//
// Before these accessors, every style that wanted tier-based metrics wrote this:
//
//	entry, ok := catalog.Get(tier)
//	if !ok {
//	    entry = tiercatalog.Entry{
//	        GlyphAdvance: hints.GlyphAdvance,
//	        RowHeight:    hints.RowHeight,
//	    }
//	}
//
// That block was duplicated in clock (twice), serial, zmq and pager. It has two
// serious problems:
//
//  1. The fallback substitutes textlayout.GlyphAdvance / textlayout.RowHeight
//     (6 and 10), which correspond to NO registered font. The style then performs
//     all of its centering and row-fitting arithmetic against metrics that the
//     renderer will never draw with, so text lands in the wrong place. This is
//     precisely the failure that made 32px-dimension panels render a few pixels
//     of text jammed into a corner: the catalog build failed, every Get returned
//     false, and every style silently laid out for a font that does not exist.
//
//  2. It is opt-in boilerplate. A style author who forgets it gets a zero Entry
//     (GlyphAdvance 0, RowHeight 0) and divides or positions by zero.
//
// [Catalog.Entry] removes both problems by guaranteeing a usable Entry for any
// tier, on any catalog, including the zero Catalog. Every fallback level resolves
// to a REAL registered font, so a style's arithmetic always matches what the
// renderer can actually draw.
//
// # Migration note for future agents
//
// [Catalog.Get] is retained unchanged. Its (Entry, bool) signature is asserted by
// a large body of tests in this package that encode the "tier absent" contract,
// and callers outside the display pipeline may still want to distinguish
// "populated" from "substituted". Prefer [Catalog.Entry] in styles and renderers;
// use Get only when the absence itself is meaningful.

// Entry returns the metrics for tier, always yielding a usable value.
//
// Resolution order, most to least specific:
//
//  1. The tier's own catalog entry, when populated. A catalog produced by a
//     successful [Build] populates every tier in [Catalog.Tiers], so this is the
//     overwhelmingly common path.
//  2. The nearest populated tier below the request, then the nearest above.
//     Descending first is deliberate: rendering text smaller than asked for
//     degrades legibility, while rendering it larger risks overflowing the region
//     and clipping. When in doubt, undershoot.
//  3. The metrics of [font.Default], for the zero Catalog (no catalog was ever
//     built for this region). This is a real registered face, so the caller's
//     layout math stays consistent with what the renderer draws.
//  4. A minimal non-zero Entry, only reachable if the font registry is empty,
//     which cannot happen in a normal build because faces self-register in init.
//     Returned solely so callers never divide by zero.
//
// The returned Entry's FontID identifies the face whose metrics it carries.
// Renderers should draw with exactly that face; see resolveTierFonts in
// runtime/ui. Any divergence between the face used for layout and the face used
// for drawing reintroduces the mis-positioning class of bug described above.
func (c Catalog) Entry(tier Tier) Entry {
	tier = NormalizeTier(tier)

	if e, ok := c.entries[tier]; ok && e.usable() {
		return e
	}

	// Walk the canonical order to find the nearest populated neighbour,
	// preferring smaller over larger.
	if idx := tierIndex(tier); idx >= 0 {
		for i := idx - 1; i >= 0; i-- {
			if e, ok := c.entries[tierOrder[i]]; ok && e.usable() {
				return e
			}
		}
		for i := idx + 1; i < len(tierOrder); i++ {
			if e, ok := c.entries[tierOrder[i]]; ok && e.usable() {
				return e
			}
		}
	} else {
		// Unknown tier string: any populated entry beats a fabricated one.
		for _, t := range tierOrder {
			if e, ok := c.entries[t]; ok && e.usable() {
				return e
			}
		}
	}

	if face := font.Default(); face != nil {
		m := face.Metrics()
		return Entry{
			GlyphWidth:   m.GlyphWidth,
			GlyphHeight:  m.GlyphHeight,
			GlyphAdvance: m.GlyphAdvance,
			RowHeight:    m.RowHeight,
			FontID:       face.ID(),
		}
	}

	// Empty font registry. Unreachable in a normal build; exists so that callers
	// dividing by GlyphAdvance or RowHeight cannot panic.
	return Entry{GlyphWidth: 1, GlyphHeight: 1, GlyphAdvance: 1, RowHeight: 1}
}

// usable reports whether an Entry carries metrics a caller can safely divide by
// and position with. A zero Entry is what a map miss yields, so this is the guard
// that distinguishes "populated" from "absent".
func (e Entry) usable() bool {
	return e.GlyphAdvance > 0 && e.RowHeight > 0
}

// MaxAdvance returns the widest GlyphAdvance permitted in this region, which is
// PixelWidth / MinChars — the same budget [Build] filtered candidates against.
//
// Exposed because tierselect was recomputing this division from PixelWidth and
// MinChars on every call. Two copies of a derivation invite drift; the catalog
// already holds both inputs, so it should own the result.
//
// Never returns less than 1, so callers can use it as a divisor or bound
// unconditionally.
func (c Catalog) MaxAdvance() int {
	mc := c.minChars
	if mc <= 0 {
		mc = defaultMinChars
	}
	if adv := c.width / mc; adv > 0 {
		return adv
	}
	return 1
}

// Relaxed reports whether this catalog was built with a MinChars lower than the
// caller requested because no registered font was narrow enough to satisfy the
// request. See Params.AllowRelaxedMinChars for the rationale.
//
// [Catalog.MinChars] reports the value actually achieved, not the value asked
// for, so the width-safety invariant GlyphAdvance*MinChars <= PixelWidth holds
// against the catalog's own reported constraint in both the strict and relaxed
// cases. RequestedMinChars reports what was originally asked.
func (c Catalog) Relaxed() bool { return c.relaxed }

// RequestedMinChars returns the MinChars originally passed to [Build], which
// differs from [Catalog.MinChars] only when [Catalog.Relaxed] is true.
func (c Catalog) RequestedMinChars() int {
	if c.requestedMinChars <= 0 {
		return c.minChars
	}
	return c.requestedMinChars
}

// NormalizeTier maps tier aliases onto their canonical constant.
//
// TierFullsize is a backward-compatibility alias for TierFull. It was previously
// normalized inline inside Get only, which meant any other code path comparing
// tier strings — iterating Tiers(), keying a map, logging — saw an alias that
// matched nothing. Normalization belongs in one exported place so all paths agree.
//
// Unrecognized tiers are returned unchanged; callers decide whether to treat that
// as an error ([Catalog.Get] reports absence) or to substitute ([Catalog.Entry]).
func NormalizeTier(tier Tier) Tier {
	if tier == TierFullsize {
		return TierFull
	}
	return tier
}

// tierIndex returns the position of tier within the canonical ascending order,
// or -1 when the tier is not a known constant.
func tierIndex(tier Tier) int {
	for i, t := range tierOrder {
		if t == tier {
			return i
		}
	}
	return -1
}
