// Package tierselect resolves (family, tier) requests into concrete font faces.
// It bridges mode intent (aesthetic preference) with region constraints (what fits).
package tierselect

import (
	"fmt"
	"log"
	"strings"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// Request describes what a display mode wants.
type Request struct {
	Family string           // Font family prefix (e.g., "spleen", "terminus", "cozette")
	Tier   tiercatalog.Tier // Desired size tier
}

// Select resolves a tier+family request to a concrete font face.
// Resolution order:
//  1. Exact match: family variant at the tier's target GlyphHeight.
//  2. Closest family variant with GlyphHeight ≤ target.
//  3. Cross-family best-fit at target height (logs a warning).
//
// Panics if:
//   - The requested family has no registered variants at all.
//   - The fallback chain is exhausted with no font found.
func Select(catalog tiercatalog.Catalog, req Request) font.Face {
	// Panic for genuinely unrecognized tier strings (not a known constant).
	// Known tiers on an empty catalog fall through to catalog.Entry's safe fallbacks.
	if !isKnownTier(req.Tier) {
		panic(fmt.Sprintf("tierselect: invalid tier %q", req.Tier))
	}
	entry := catalog.Entry(req.Tier)
	targetHeight := entry.GlyphHeight

	// Compute maxAdvance from catalog constraints for width-safe selection.
	maxAdvance := catalog.PixelWidth() / catalog.MinChars()
	if maxAdvance <= 0 {
		maxAdvance = 1
	}

	// Strategy 1: Exact match within family (respecting width constraint).
	candidate := findFamilyVariant(req.Family, targetHeight, maxAdvance)
	if candidate != nil {
		return candidate
	}

	// Guard: if the requested family has no registered variants at all, panic immediately.
	if !familyHasAnyVariant(req.Family) {
		panic(fmt.Sprintf("tierselect: family %q has no registered variants (tier=%q, targetHeight=%d)",
			req.Family, req.Tier, targetHeight))
	}

	// Strategy 2: Closest family variant with GlyphHeight ≤ targetHeight
	// that also satisfies the width constraint.
	candidate = closestFamilyVariant(req.Family, targetHeight, maxAdvance)
	if candidate != nil {
		return candidate
	}

	// Strategy 3: Best-fit face at target height regardless of family.
	// Region constraints take priority: only consider faces that satisfy the
	// catalog's width constraint (GlyphAdvance * MinChars ≤ PixelWidth).
	candidate = crossFamilyFallback(targetHeight, maxAdvance)
	if candidate != nil {
		log.Printf("tierselect: family %q has no variant at height≤%d, using %q", req.Family, targetHeight, candidate.ID())
		return candidate
	}

	panic(fmt.Sprintf("tierselect: no font found for family=%q tier=%q targetHeight=%d; registry has %d fonts",
		req.Family, req.Tier, targetHeight, len(font.List())))
}

// SelectMulti resolves multiple tier+family requests in one call.
// Returns a slice of exactly len(reqs) faces, one per request.
func SelectMulti(catalog tiercatalog.Catalog, reqs []Request) []font.Face {
	result := make([]font.Face, len(reqs))
	for i, req := range reqs {
		result[i] = Select(catalog, req)
	}
	return result
}

// findFamilyVariant looks for an exact family match at the given target height
// that also satisfies GlyphAdvance ≤ maxAdvance (width constraint).
// When multiple variants of the same family share the target height, the one with
// the smallest font ID (lexicographic ascending) is returned for determinism.
func findFamilyVariant(family string, targetHeight, maxAdvance int) font.Face {
	var best font.Face
	for _, f := range font.List() {
		if extractFamily(f.ID()) != family {
			continue
		}
		m := f.Metrics()
		if m.GlyphHeight != targetHeight {
			continue
		}
		if m.GlyphAdvance > maxAdvance {
			continue
		}
		if best == nil || f.ID() < best.ID() {
			best = f
		}
	}
	return best
}

// closestFamilyVariant finds the family variant with the largest GlyphHeight ≤ maxHeight
// that also satisfies GlyphAdvance ≤ maxAdvance (width constraint).
// When multiple variants share the same best height, ties are broken by font ID ascending.
func closestFamilyVariant(family string, maxHeight, maxAdvance int) font.Face {
	var best font.Face
	for _, f := range font.List() {
		if extractFamily(f.ID()) != family {
			continue
		}
		m := f.Metrics()
		if m.GlyphHeight > maxHeight {
			continue
		}
		if m.GlyphAdvance > maxAdvance {
			continue
		}
		if best == nil || m.GlyphHeight > best.Metrics().GlyphHeight || (m.GlyphHeight == best.Metrics().GlyphHeight && f.ID() < best.ID()) {
			best = f
		}
	}
	return best
}

// familyHasAnyVariant checks whether any font in the registry belongs to the given family.
func familyHasAnyVariant(family string) bool {
	for _, f := range font.List() {
		if extractFamily(f.ID()) == family {
			return true
		}
	}
	return false
}

// crossFamilyFallback finds the best font at or below targetHeight whose
// GlyphAdvance satisfies the catalog's width constraint. This ensures Strategy 3
// never returns a font that's too wide for the region.
// Selection priority: largest GlyphHeight ≤ targetHeight, then highest family
// priority, then font ID ascending.
func crossFamilyFallback(targetHeight, maxAdvance int) font.Face {
	var best font.Face
	for _, f := range font.List() {
		m := f.Metrics()
		if m.GlyphHeight > targetHeight {
			continue
		}
		if m.GlyphAdvance > maxAdvance {
			continue
		}
		if best == nil {
			best = f
			continue
		}
		bm := best.Metrics()
		// Prefer larger height, then higher family priority, then smaller ID.
		if m.GlyphHeight > bm.GlyphHeight {
			best = f
		} else if m.GlyphHeight == bm.GlyphHeight {
			fp := familyPriority(f.ID())
			bp := familyPriority(best.ID())
			if fp > bp || (fp == bp && f.ID() < best.ID()) {
				best = f
			}
		}
	}
	return best
}

// familyPriority returns the tie-breaking priority for a font face.
// Spleen has highest priority (3), then Terminus (2), then Cozette (1), others (0).
func familyPriority(id string) int {
	lower := strings.ToLower(id)
	switch {
	case strings.HasPrefix(lower, "spleen-"):
		return 3
	case strings.HasPrefix(lower, "terminus-"):
		return 2
	case strings.HasPrefix(lower, "cozette-"):
		return 1
	default:
		return 0
	}
}

// extractFamily splits a font ID on the first "-" and returns the family prefix.
// For example, "spleen-8x16" → "spleen", "terminus-10x20" → "terminus".
// If no "-" exists, the entire ID is returned as the family name.
func extractFamily(id string) string {
	idx := strings.Index(id, "-")
	if idx < 0 {
		return id
	}
	return id[:idx]
}

// isKnownTier reports whether tier is one of the canonical tier constants.
// Unrecognized strings (e.g. "bogus") return false; known tiers return true
// even if the catalog was not built with that tier populated.
func isKnownTier(tier tiercatalog.Tier) bool {
	n := tiercatalog.NormalizeTier(tier)
	return n == tiercatalog.TierSmall ||
		n == tiercatalog.TierNormal ||
		n == tiercatalog.TierLarge ||
		n == tiercatalog.TierHuge ||
		n == tiercatalog.TierColossal ||
		n == tiercatalog.TierFull
}
