package styles

import (
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// FormatDirective applies rendering constraints to a single LineDirective
// given panel hints and the global policy as fallback. Font resolution uses
// hints.Face for catalog-validated font resolution.
//
// Font resolution chain (tier-based):
//   - directive font → validate against catalog via hints.Face
//   - policy font → validate against catalog via hints.Face
//   - auto → use the resolved ticker font tier for this panel
//
// Named fonts in directives are validated: if the named font violates catalog
// constraints, resolution falls back to the mode's default tier.
func FormatDirective(d source.LineDirective, hints textlayout.TextHints, policy source.Policy, borderInset int) source.FormattedLine {
	// Step 1: Resolve effective scaling.
	scaling := d.Scaling
	if scaling == "" {
		scaling = "fixed"
	}

	// Step 2: Resolve font via hints.Face (tier-based).
	var face font.Face
	var tier tiercatalog.Tier
	var fontWarning string
	fontTier := source.ResolveFontTier(policy, hints)

	if scaling == "fit" {
		// Fit-to-width: use the largest tier whose face fits the full text.
		face, tier = fitFontToWidthTier(d.Text, hints.PixelWidth-borderInset*2, hints)
	} else {
		// Fixed: resolve directive font → policy font → auto (all through tier catalog).
		face, tier, fontWarning = resolveDirectiveFont(d, policy, hints, fontTier)
	}

	// Step 3: Compute max chars for this line's font.
	metrics := face.Metrics()
	usableWidth := hints.PixelWidth - borderInset*2
	maxChars := 0
	if metrics.GlyphAdvance > 0 {
		maxChars = usableWidth / metrics.GlyphAdvance
	}

	// Step 4: Apply line_mode.
	lineMode := d.LineMode
	if lineMode == "" {
		lineMode = policy.LineMode
	}
	if lineMode == "" {
		lineMode = textlayout.LineModeTruncate
	}

	var formatted string
	if maxChars <= 0 {
		formatted = ""
	} else {
		switch lineMode {
		case textlayout.LineModeClip:
			formatted = clipToChars(d.Text, maxChars)
		default:
			formatted = textlayout.Truncate(d.Text, maxChars)
		}
	}

	return source.FormattedLine{Text: formatted, Tier: tier, FontWarning: fontWarning}
}

// resolveDirectiveFont resolves the font face for a directive using the
// tier catalog. It follows the chain: directive font → policy font → auto.
// Named fonts are validated against catalog constraints via hints.Face.
// If a named font is not available or violates constraints, falls back to
// tier-based resolution.
func resolveDirectiveFont(d source.LineDirective, policy source.Policy, hints textlayout.TextHints, tier tiercatalog.Tier) (font.Face, tiercatalog.Tier, string) {
	fontID := d.Font
	if fontID == "" {
		fontID = policy.Font
	}

	if fontID == "" || fontID == "auto" {
		face := source.ResolveFace(hints, "spleen", tier)
		return face, tier, ""
	}

	// Named font: validate against catalog constraints via hints.Face.
	// hints.Face only returns catalog-validated faces. We attempt to
	// find the named font in the catalog by trying to resolve it as a family
	// preference. If the named font family is recognized, use it; otherwise
	// fall back to the default tier.
	face := hints.Face(fontID, tier)
	if face != nil && face.ID() != "" {
		// Validate: the face must satisfy width constraints (hints.Face guarantees this).
		return face, tier, ""
	}

	// Unrecognized font family: fall back to default tier resolution.
	face = source.ResolveFace(hints, "spleen", tier)
	return face, tier, "unrecognized font " + fontID
}

// PartitionScroll divides the feed into pinned (static) and scrollable lines.
// Returns the visible items in display order for the current scroll offset.
// Pinned lines remain at their declared position indices; scrollable lines
// rotate based on scrollOffset % len(scrollableLines). If all lines are pinned,
// scroll offset is forced to 0.
func PartitionScroll(directives []source.LineDirective, hints textlayout.TextHints, policy source.Policy, borderInset int, scrollOffset int) []source.FormattedLine {
	// Separate pinned and scrollable directives by index.
	type indexed struct {
		index int
		d     source.LineDirective
	}
	var pinned []indexed
	var scrollable []indexed

	for i, d := range directives {
		scroll := d.Scroll
		if scroll == "" {
			scroll = "normal"
		}
		entry := indexed{index: i, d: d}
		if scroll == "pinned" {
			pinned = append(pinned, entry)
		} else {
			scrollable = append(scrollable, entry)
		}
	}

	// If all pinned, format all in order, no scrolling.
	if len(scrollable) == 0 {
		result := make([]source.FormattedLine, len(directives))
		for i, d := range directives {
			result[i] = FormatDirective(d, hints, policy, borderInset)
		}
		return result
	}

	// Rotate scrollable pool by scrollOffset.
	effectiveOffset := scrollOffset % len(scrollable)
	rotated := make([]indexed, len(scrollable))
	for i := range scrollable {
		rotated[i] = scrollable[(i+effectiveOffset)%len(scrollable)]
	}

	// Reconstruct display order: pinned at their indices, scrollable fill gaps.
	result := make([]source.FormattedLine, len(directives))
	pinnedSet := make(map[int]bool)
	for _, p := range pinned {
		pinnedSet[p.index] = true
	}

	scrollIdx := 0
	for i, d := range directives {
		if pinnedSet[i] {
			result[i] = FormatDirective(d, hints, policy, borderInset)
		} else {
			if scrollIdx < len(rotated) {
				result[i] = FormatDirective(rotated[scrollIdx].d, hints, policy, borderInset)
				scrollIdx++
			}
		}
	}
	return result
}

// formatNonScrollingLines returns FormattedLine entries only for lines that
// do not scroll (pinned or text fits within usable width). Scrolling lines
// are omitted entirely (they exist as Sprites, not Items).
func formatNonScrollingLines(directives []source.LineDirective, hints textlayout.TextHints, policy source.Policy, borderInset int) []source.FormattedLine {
	usableWidth := hints.PixelWidth - 2*borderInset
	fontTier := source.ResolveFontTier(policy, hints)

	var result []source.FormattedLine
	for _, d := range directives {
		fontID := d.Font
		if fontID == "" {
			fontID = policy.Font
		}
		if fontID == "" || fontID == "auto" {
			fontID = "spleen"
		}
		face := source.ResolveFace(hints, fontID, fontTier)
		metrics := face.Metrics()
		textWidth := utf8.RuneCountInString(d.Text) * metrics.GlyphAdvance

		scroll := d.Scroll
		if scroll == "" {
			scroll = "normal"
		}

		// Include line if: pinned, zero-advance font, or text fits within usable width.
		if scroll == "pinned" || metrics.GlyphAdvance == 0 || textWidth <= usableWidth {
			result = append(result, FormatDirective(d, hints, policy, borderInset))
		}
		// Otherwise: scrolling line, skip (represented as a Sprite).
	}
	return result
}

// fitFontToWidthTier selects the largest tier from the catalog where the full text
// fits within the given pixel width. Falls back to the smallest tier if none fits.
func fitFontToWidthTier(text string, pixelWidth int, hints textlayout.TextHints) (font.Face, tiercatalog.Tier) {
	textLen := utf8.RuneCountInString(text)
	if textLen == 0 {
		face := source.ResolveFace(hints, "spleen", tiercatalog.TierNormal)
		return face, tiercatalog.TierNormal
	}

	// Try tiers from largest to smallest.
	tiersDesc := []tiercatalog.Tier{
		tiercatalog.TierFullsize,
		tiercatalog.TierLarge,
		tiercatalog.TierNormal,
		tiercatalog.TierSmall,
	}

	for _, tier := range tiersDesc {
		face := source.ResolveFace(hints, "spleen", tier)
		m := face.Metrics()
		if m.GlyphAdvance*textLen <= pixelWidth {
			return face, tier
		}
	}

	// No tier fits: return smallest.
	face := source.ResolveFace(hints, "spleen", tiercatalog.TierSmall)
	return face, tiercatalog.TierSmall
}

func clipToChars(s string, maxChars int) string {
	r := []rune(s)
	if len(r) <= maxChars {
		return s
	}
	return string(r[:maxChars])
}
