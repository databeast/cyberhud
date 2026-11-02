package styles

import (
	"fmt"
	"image"
	"strings"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/modes/usb/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// IconGetter retrieves a named icon image for USB status sprites.
type IconGetter func(name string) (image.Image, bool)

type Params struct {
	BuildFn func(data source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name       string
	reqs       style.SurfaceRequirements
	p          Params
	iconGetter IconGetter
}

var _ style.Style[source.Snapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	return adaptiveBuild(data, pol, ctx, d)
}

func (d *def) SetIconGetter(fn func(name string) (image.Image, bool)) {
	d.iconGetter = fn
}

// ──────────────────────────────────────────────────────────────────────────────
// adaptiveBuild — fallback layout selection for skeleton styles (nil BuildFn).
//
// Selects a core layout based on panel pixel dimensions:
//   - w ≥ 240 && h ≥ 240 → buildFull  (multi-row device detail)
//   - w ≥ 128            → buildCompact (single condensed row)
//   - w < 128            → buildMinimal (device name only)
// ──────────────────────────────────────────────────────────────────────────────

func adaptiveBuild(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	w := hints.PixelWidth
	h := hints.PixelHeight

	var result style.ViewData
	switch {
	case w >= 240 && h >= 240:
		result = buildFull(snapshot, pol, ctx)
	case w >= 128:
		result = buildCompact(snapshot, pol, ctx)
	default:
		result = buildMinimalLayout(snapshot, pol, ctx)
	}

	if len(result.Items) == 0 || allEmptyItems(result.Items) {
		result.Items = []string{"usb"}
	}
	result.Static = true
	return result
}

// ──────────────────────────────────────────────────────────────────────────────
// buildFull — multi-row layout for panels ≥ 240×240.
//
// Rows:
//   - Row 0: Device name at TierLarge
//   - Row 1: VID:PID at TierSmall
//   - Row 2: Bus N Dev M at TierSmall
//   - Row 3: SN <serial> at TierSmall (if serial present)
//   - Row 4: Status at TierSmall
//
// Landscape panels: center each row. Portrait panels: left-align.
// Rows clamped to bridge.MaxVisibleRows().
// Text truncated per tier's available character width.
// ──────────────────────────────────────────────────────────────────────────────

func buildFull(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	// buildFull now delegates to the sprite-plane hero layout (hero.go): a large,
	// class-colored device-class glyph as the focal point, the device name as a
	// big title, a class accent bar, and quiet supporting detail rows. Everything
	// is composited as absolute-positioned sprites with empty ViewData.Items, so
	// the two-plane (row-text vs pixel-sprite) alignment hazard cannot occur.
	return buildHero(snapshot, pol, ctx)
}

// ──────────────────────────────────────────────────────────────────────────────
// buildCompact — single-row layout for panels 128–239px wide.
//
// Row 0: "<device-name> <VID>:<PID>" at TierSmall, centered.
// ──────────────────────────────────────────────────────────────────────────────

func buildCompact(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	// Compact panels get the badge treatment: a class glyph above the name with
	// the VID:PID as a quiet sub-line, vertically centered. On short landscape
	// strips without vertical room, buildCompactBadge drops the glyph and falls
	// back to name + sub-line. This replaces the old lone condensed text line,
	// which rode in the Items plane and looked bare.
	if !snapshot.HasLast {
		return buildCompactBadge(snapshot, ctx, "Waiting for USB...", "")
	}
	name := coreDisplayName(snapshot.Device)
	sub := fmt.Sprintf("%s:%s", snapshot.Device.VendorID, snapshot.Device.ProductID)
	return buildCompactBadge(snapshot, ctx, name, sub)
}

// buildMinimal is the BuildFn-compatible wrapper for buildMinimalLayout.
// It delegates to the 3-arg layout function, ignoring the def parameter.
func buildMinimal(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildMinimalLayout(snapshot, pol, ctx)
}

// ──────────────────────────────────────────────────────────────────────────────
// buildMinimalLayout — single-row layout for panels < 128px wide.
//
// Row 0: device name at TierLarge, left-aligned.
// ──────────────────────────────────────────────────────────────────────────────

func buildMinimalLayout(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	// Minimal panels are narrow but TALL (width < 128, height ≥ 128): a lone
	// centered name floated in a mostly-empty strip and read as "failed to load".
	// Give them the class glyph stacked above the name — it fills the vertical
	// void and still delivers device-type recognition. buildCompactBadge drops
	// the glyph automatically if a given panel lacks the vertical room.
	if !snapshot.HasLast {
		return buildCompactBadge(snapshot, ctx, "Waiting for USB...", "")
	}
	return buildCompactBadge(snapshot, ctx, coreDisplayName(snapshot.Device), "")
}

// ──────────────────────────────────────────────────────────────────────────────
// Core helpers — used by the adaptive layout functions above.
// ──────────────────────────────────────────────────────────────────────────────

// coreDisplayName returns the device name for rendering. Prefers Product,
// then Manufacturer, then a generic fallback.
func coreDisplayName(info source.DeviceInfo) string {
	name := strings.TrimSpace(info.Product)
	if name != "" {
		return name
	}
	name = strings.TrimSpace(info.Manufacturer)
	if name != "" {
		return name
	}
	return "Unknown USB Device"
}

// coreSafeValue returns v trimmed of whitespace, or fallback if v is empty.
func coreSafeValue(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

// nameTierLadder is the descending list of tiers fitName tries for the device
// name — largest first. The name is the hero of the USB bench screen, so it gets
// the biggest font whose whole string fits the width.
var nameTierLadder = []tiercatalog.Tier{
	tiercatalog.TierLarge,
	tiercatalog.TierNormal,
	tiercatalog.TierSmall,
}

// fitName returns the device name sized to the largest tier in nameTierLadder
// whose full string fits within availW pixels, along with that tier. It never
// appends an ellipsis: the bitmap fonts have no U+2026 glyph, so Truncate's
// ellipsis renders as a tofu box (the exact defect this replaces). If the name
// does not fit even at the smallest tier, it is hard-truncated with no indicator
// — a plain clipped word reads as "cut", a trailing tofu box reads as "broken".
func fitName(name string, ctx style.StyleContext, availW int) (string, tiercatalog.Tier) {
	name = strings.TrimSpace(name)
	if availW <= 0 {
		return name, tiercatalog.TierSmall
	}
	n := utf8.RuneCountInString(name)
	for _, tier := range nameTierLadder {
		ga := ctx.Entry(tier).GlyphAdvance
		if ga <= 0 {
			continue
		}
		if n*ga <= availW {
			return name, tier
		}
	}
	// Does not fit at any tier: use the smallest tier and plain-truncate to width.
	small := tiercatalog.TierSmall
	ga := ctx.Entry(small).GlyphAdvance
	if ga <= 0 {
		return name, small
	}
	return truncatePlain(name, availW/ga), small
}

// truncatePlain returns s trimmed to at most max runes, WITHOUT appending an
// ellipsis. Used in place of textlayout.Truncate wherever the rendered font is a
// bitmap face that lacks the U+2026 ellipsis glyph (all the panel fonts), so the
// indicator would otherwise render as a tofu box.
func truncatePlain(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	return string([]rune(s)[:max])
}
