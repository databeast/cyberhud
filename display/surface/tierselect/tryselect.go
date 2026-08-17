package tierselect

import (
	"fmt"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// Non-panicking selection for the render path.
//
// # Why this exists alongside Select
//
// [Select] panics on three conditions: an unknown tier, a family with no
// registered variants, and an exhausted fallback chain. That contract is
// deliberate and is asserted by dedicated tests (Requirements 6.3-6.5), so it is
// left exactly as it is. It suits build-time and configuration-time callers, where
// a panic surfaces a programming error loudly and early.
//
// It does not suit the render loop. RegionRenderer.Render resolves a tier for every
// row of every frame on real hardware. Two of Select's three panic conditions are
// not purely programmer error from the renderer's point of view:
//
//   - "family has no registered variants" depends on which generated font files are
//     compiled into the binary, which varies by build configuration.
//   - "unknown tier" depends on the tier strings a display mode emits, which is
//     mode-authored data flowing through ViewData.Tiers.
//
// A mode shipping a typo'd tier string should degrade to a legible fallback and log,
// not take down the display process mid-frame. TrySelect provides that: it never
// panics, and it always returns a usable face when any font is registered.
//
// # Guidance for future agents
//
// Use TrySelect anywhere inside a render or frame-production path. Use Select for
// setup, validation and tooling, where failing loudly is the point.

// TrySelect resolves a tier+family request to a concrete font face without
// panicking.
//
// It returns the face plus a reason string that is empty on a clean resolution and
// otherwise describes the degradation that occurred. Callers that log should log
// the reason at most once per condition rather than per frame.
//
// Resolution order:
//
//  1. [Select]'s normal strategies, when the tier is present in the catalog and the
//     family has registered variants. This is the common path and produces results
//     identical to Select.
//  2. The catalog's own entry for the tier via [tiercatalog.Catalog.Entry], resolved
//     back to its registered face. Entry cannot fail and its FontID names a real
//     font, so this covers both an unknown tier and a missing family.
//  3. [font.ByHeight] against the entry's glyph height, bounded by the catalog's
//     MaxAdvance so the result still fits the region horizontally.
//  4. [font.Default].
//
// Returns (nil, reason) only when the font registry is completely empty, which
// cannot occur in a normal build because faces self-register in init. Callers in
// the render path should skip font switching on a nil face and let the surface keep
// its current one.
func TrySelect(catalog tiercatalog.Catalog, req Request) (face font.Face, reason string) {
	// Fast path: delegate to Select, converting any panic into a reason. Select
	// panics with a string, but recover is typed as any, so handle both.
	face, panicked := selectRecovering(catalog, req)
	if !panicked {
		return face, ""
	}

	// Select rejected the request. Fall back to the catalog's own answer, which is
	// guaranteed usable and names a real registered face.
	entry := catalog.Entry(req.Tier)
	if entry.FontID != "" {
		if f, ok := font.Get(entry.FontID); ok {
			return f, fmt.Sprintf("tier %q family %q unresolvable; using catalog entry %q",
				req.Tier, req.Family, entry.FontID)
		}
	}

	// The catalog named a font that is no longer registered, or carried no ID.
	// Fall back to a height-based match that still respects the region's width.
	maxAdvance := catalog.MaxAdvance()
	if f := byHeightWithin(entry.GlyphHeight, maxAdvance); f != nil {
		return f, fmt.Sprintf("tier %q family %q unresolvable; using height match %q",
			req.Tier, req.Family, f.ID())
	}

	if f := font.Default(); f != nil {
		return f, fmt.Sprintf("tier %q family %q unresolvable; using default face %q",
			req.Tier, req.Family, f.ID())
	}

	return nil, fmt.Sprintf("tier %q family %q unresolvable and no fonts registered",
		req.Tier, req.Family)
}

// selectRecovering calls Select and reports whether it panicked, so TrySelect can
// reuse Select's full strategy chain without duplicating it. Duplicating the
// strategies would be worse than recovering: two implementations of font choice
// would drift, and the whole point of these accessors is that layout and drawing
// agree on one face.
func selectRecovering(catalog tiercatalog.Catalog, req Request) (face font.Face, panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			face = nil
			panicked = true
		}
	}()
	f := Select(catalog, req)
	if f == nil {
		return nil, true
	}
	return f, false
}

// byHeightWithin returns the tallest registered face no taller than targetHeight
// whose advance fits maxAdvance, or nil when nothing qualifies.
//
// font.ByHeight alone is not sufficient here: it ignores advance, so on a narrow
// region it can return a face too wide to fit, which is the defect the catalog's
// advance budget exists to prevent.
func byHeightWithin(targetHeight, maxAdvance int) font.Face {
	var best font.Face
	for _, f := range font.List() {
		m := f.Metrics()
		if m.GlyphAdvance > maxAdvance {
			continue
		}
		if targetHeight > 0 && m.GlyphHeight > targetHeight {
			continue
		}
		if best == nil || m.GlyphHeight > best.Metrics().GlyphHeight {
			best = f
		}
	}
	return best
}
