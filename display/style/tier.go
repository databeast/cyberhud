package style

import "github.com/databeast/cyberhud/display/surface/tiercatalog"

// Default tier selection heuristics.
//
// # Why these live here
//
// "Pick a font tier based on how tall the region is" is a decision every tiered
// display mode has to make, and it is not mode-specific. The clock mode grew its
// own copy as selectTier in display/modes/clock/styles/layouts.go, with a ladder of
// height thresholds. Any second mode wanting the same behaviour would have
// reinvented it, and the two copies would drift — which is exactly how the display
// system accumulated its other duplicated-layout-math bugs.
//
// These helpers are defaults, not policy. A style that wants a specific tier still
// sets Params.Tier (or its mode's equivalent) and overrides the heuristic entirely.
// The framework never forces a tier on a style; it only answers the question when
// the style declines to.

// TierForHeight returns a sensible primary tier for a region of the given pixel
// height.
//
// The thresholds are inherited from the clock mode's original selectTier so that
// adopting this helper is behaviour-preserving for it. They ascend with room to
// display: a region under 128px tall can only afford small text once date and
// weekday rows are stacked, while a 600px-tall region can give the primary row a
// colossal glyph and still fit secondary rows beneath it.
//
// Note that the returned tier is a request, not a guarantee. The region's advance
// budget still constrains what the catalog can actually supply, so a very narrow
// but tall region may resolve colossal to a modest font. That is the correct
// division of responsibility: the style says how big it wants the text, the catalog
// says how big the region can afford.
//
// Non-positive heights return TierSmall, the safest request.
func TierForHeight(pixelHeight int) tiercatalog.Tier {
	switch {
	case pixelHeight >= 600:
		return tiercatalog.TierColossal
	case pixelHeight >= 400:
		return tiercatalog.TierHuge
	case pixelHeight >= 240:
		return tiercatalog.TierLarge
	case pixelHeight >= 128:
		return tiercatalog.TierNormal
	default:
		return tiercatalog.TierSmall
	}
}

// SecondaryTier returns the tier one step below primary, for rows that should read
// as subordinate to the main row (a date under a time, a label under a value).
//
// TierFull and TierColossal both step down to TierHuge: TierFull is "the largest
// font that fits" rather than a size, so its neighbour below is the largest named
// tier. Unrecognized tiers step down to TierSmall, which cannot overflow.
func SecondaryTier(primary tiercatalog.Tier) tiercatalog.Tier {
	switch tiercatalog.NormalizeTier(primary) {
	case tiercatalog.TierFull, tiercatalog.TierColossal:
		return tiercatalog.TierHuge
	case tiercatalog.TierHuge:
		return tiercatalog.TierLarge
	case tiercatalog.TierLarge:
		return tiercatalog.TierNormal
	case tiercatalog.TierNormal:
		return tiercatalog.TierSmall
	default:
		return tiercatalog.TierSmall
	}
}
