package tiercatalog

import "math"

// maxTierHeightFraction bounds any tier's pixel target to this fraction of the
// region's height, expressed as a percentage.
//
// Why a bound is needed: both target tables are absolute. DefaultTargetsPx states
// colossal is 40px and DefaultTargetsMM states 18mm, neither of which knows how
// tall the region is. On a 64px-tall panel, colossal resolves to 40px (63% of the
// panel) or, at the assumed 96 PPI, 68px — taller than the panel itself. Nothing
// downstream corrects this: bestFit then simply picks the largest font available,
// so the *request* is meaningless and the result is decided entirely by the
// advance-budget filter.
//
// Bounding the target to half the region height keeps the request meaningful on
// small panels while leaving normal-sized panels untouched: at 240px tall the
// bound is 120px, far above every tier target, so it never binds. This is a
// sanity clamp, not a design parameter — it exists so the numbers mean what they
// say, and so a future agent reading a resolved target can trust it.
const maxTierHeightFraction = 50

// Small text is allowed a tighter cap on short panels so it does not occupy too
// much of the available height. This keeps 128x64-class mono panels from
// selecting a visually heavy "small" face while leaving taller panels unchanged.
const maxSmallTierHeightFraction = 20

// resolvePixelTargets computes the pixel height target for each tier (except
// TierFull) based on the given Params. When PPI is positive, mm targets are
// converted to pixels via round(mm × PPI / 25.4). When PPI is zero or
// negative, fixed pixel fallback targets are used instead.
//
// # A caution about PPI for future agents
//
// The mm path is the design intent: legibility is a physical property, so a tier
// should mean a real-world glyph size. That intent only holds when PPI is real.
//
// region.resolvePPI currently substitutes an assumed desktop-monitor PPI when the
// panel does not report one (see AssumedPPI there), which means the mm targets are
// converted using a number that is usually wrong for small embedded panels — a
// 1.3in 240x240 display is roughly 260 PPI, not 96. The conversion is therefore
// self-consistent but not physically meaningful on most hardware.
//
// The fix is to have panel drivers report real PPI, not to change these tables.
// Switching the fallback to the px path instead would trade one guess for another
// while visibly changing font selection on every existing panel, so it was left
// alone deliberately. See the note on AssumedPPI in display/region.
func resolvePixelTargets(p Params) map[Tier]int {
	// Normalize negative PPI to zero.
	ppi := p.PPI
	if ppi < 0 {
		ppi = 0
	}

	targets := make(map[Tier]int, len(tierOrder)-1)

	if ppi > 0 {
		// MM-to-pixel conversion path.
		for _, tier := range tierOrder {
			if tier == TierFull {
				continue
			}
			mm, ok := p.TierTargetsMM[tier]
			if !ok {
				mm = DefaultTargetsMM[tier]
			}
			targets[tier] = int(math.Round(mm * ppi / 25.4))
		}
	} else {
		// Pixel fallback path (PPI == 0).
		for _, tier := range tierOrder {
			if tier == TierFull {
				continue
			}
			px, ok := p.TierTargetsPx[tier]
			if !ok {
				px = DefaultTargetsPx[tier]
			}
			targets[tier] = px
		}
	}

	clampTargetsToRegion(targets, p.PixelHeight)

	return targets
}

// clampTargetsToRegion bounds every target to maxTierHeightFraction of
// pixelHeight, preserving ascending tier order.
//
// Clamping can collapse several tiers onto the same value on a short panel. That
// is correct and intended: if the region cannot accommodate distinct huge and
// colossal sizes, those tiers genuinely are the same size here. enforceMonotonicity
// already permits such collapsing.
//
// A non-positive pixelHeight leaves targets untouched, because there is no region
// to bound against and the caller is already in an error path.
func clampTargetsToRegion(targets map[Tier]int, pixelHeight int) {
	if pixelHeight <= 0 {
		return
	}
	for tier, px := range targets {
		limit := pixelHeight * maxTierHeightFraction / 100
		if tier == TierSmall {
			smallLimit := pixelHeight * maxSmallTierHeightFraction / 100
			if smallLimit > 0 && smallLimit < limit {
				limit = smallLimit
			}
		}
		if limit < 1 {
			limit = 1
		}
		if px > limit {
			targets[tier] = limit
		}
	}
}
