package styles

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/modes/usb/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
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
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)

	// Determine orientation: landscape = width > height, portrait = height > width.
	landscape := hints.PixelWidth > hints.PixelHeight

	// Build rows with tier assignments.
	type row struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	var rows []row

	if !snapshot.HasLast {
		rows = []row{{text: "Waiting for USB...", tier: tiercatalog.TierLarge, entry: largeEntry}}
	} else {
		// Row 0: device name
		name := coreDisplayName(snapshot.Device)
		rows = append(rows, row{text: name, tier: tiercatalog.TierLarge, entry: largeEntry})

		// Row 1: VID:PID
		rows = append(rows, row{
			text: fmt.Sprintf("VID:PID %s:%s", snapshot.Device.VendorID, snapshot.Device.ProductID),
			tier:  tiercatalog.TierSmall,
			entry: smallEntry,
		})

		// Row 2: Bus/Dev
		if snapshot.Device.BusNum != "" || snapshot.Device.DevNum != "" {
			rows = append(rows, row{
				text: fmt.Sprintf("Bus %s Dev %s", coreSafeValue(snapshot.Device.BusNum, "?"), coreSafeValue(snapshot.Device.DevNum, "?")),
				tier:  tiercatalog.TierSmall,
				entry: smallEntry,
			})
		}

		// Row 3: Serial (if present)
		if snapshot.Device.Serial != "" {
			rows = append(rows, row{
				text:  "SN " + snapshot.Device.Serial,
				tier:  tiercatalog.TierSmall,
				entry: smallEntry,
			})
		}

		// Row 4: Status
		status := "connected"
		if !snapshot.Connected {
			status = "unplugged"
		}
		rows = append(rows, row{text: "Status " + status, tier: tiercatalog.TierSmall, entry: smallEntry})
	}

	// Clamp to MaxVisibleRows.
	maxRows := bridge.MaxVisibleRows()
	if maxRows > 0 && len(rows) > maxRows {
		rows = rows[:maxRows]
	}

	// Truncate text per tier and compute offsets.
	items := make([]string, len(rows))
	tiers := make([]tiercatalog.Tier, len(rows))
	offsets := make([]int, len(rows))
	ox, _ := bridge.ContentOrigin()

	for i, r := range rows {
		ga := r.entry.GlyphAdvance
		if ga <= 0 {
			ga = bridge.GlyphAdvance()
		}

		// Truncate to available character width for this tier.
		maxChars := 0
		if ga > 0 {
			maxChars = bridge.AvailableContentWidth() / ga
		}
		text := r.text
		if maxChars > 0 {
			text = textlayout.Truncate(text, maxChars)
		}
		items[i] = text
		tiers[i] = r.tier

		// Alignment based on orientation.
		if landscape {
			offsets[i] = bridge.CenterXWith(utf8.RuneCountInString(text), ga)
		} else {
			offsets[i] = ox
		}
	}

	// Compute vertical offset using FitRows for centering.
	rowHeights := make([]int, len(rows))
	for i, r := range rows {
		rh := r.entry.RowHeight
		if rh <= 0 {
			rh = bridge.RowHeight()
		}
		rowHeights[i] = rh
	}
	_, offsetY, _ := bridge.FitRows(rowHeights)

	// Determine if icon will be rendered (decides layout adjustments).
	const iconSize = 24
	const iconGap = 4
	const iconShift = iconSize + iconGap // 28px total displacement

	var iconSprite *widgets.Sprite
	// Defensive guard: only render icon when panel meets full-layout threshold
	// and a device is actually connected. The threshold check is normally
	// guaranteed by adaptiveBuild routing, but we guard here as well.
	if snapshot.HasLast && hints.PixelWidth >= 240 && hints.PixelHeight >= 240 {
		iconRune := ClassIcon(snapshot.Device.DeviceClass)
		iconFace, ok := font.Get("material-icons-24")
		if ok {
			var iconX, iconY int
			if landscape {
				// Landscape: icon to the left of first row at content origin X.
				iconX = ox
				iconY = offsetY
			} else {
				// Portrait: icon above first row, centered horizontally.
				iconX = (hints.PixelWidth - iconSize) / 2
				iconY = offsetY
			}
			iconSprite = textlabel.Render(textlabel.Config{
				Text:       string(iconRune),
				Bounds:     image.Rect(iconX, iconY, iconX+iconSize, iconY+iconSize),
				Font:       iconFace,
				Foreground: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			})
			if iconSprite != nil {
				iconSprite.Label = "usb-class-icon"
			}
		}
	}

	// Adjust text offsets to account for icon displacement.
	if iconSprite != nil {
		if landscape {
			// Shift all text rows right by iconShift.
			for i := range offsets {
				offsets[i] += iconShift
			}
		} else {
			// Shift text rows down by iconShift.
			offsetY += iconShift
		}
	}

	result := style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
	}

	if iconSprite != nil {
		result.Sprites = append(result.Sprites, *iconSprite)
	}

	return result
}

// ──────────────────────────────────────────────────────────────────────────────
// buildCompact — single-row layout for panels 128–239px wide.
//
// Row 0: "<device-name> <VID>:<PID>" at TierSmall, centered.
// ──────────────────────────────────────────────────────────────────────────────

func buildCompact(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	entry := ctx.Entry(tiercatalog.TierSmall)
	ga := entry.GlyphAdvance
	if ga <= 0 {
		ga = bridge.GlyphAdvance()
	}

	var text string
	if !snapshot.HasLast {
		text = "Waiting for USB..."
	} else {
		name := coreDisplayName(snapshot.Device)
		text = fmt.Sprintf("%s %s:%s", name, snapshot.Device.VendorID, snapshot.Device.ProductID)
	}

	// Truncate to available character width.
	maxChars := 0
	if ga > 0 {
		maxChars = bridge.AvailableContentWidth() / ga
	}
	if maxChars > 0 {
		text = textlayout.Truncate(text, maxChars)
	}

	// Clamp to MaxVisibleRows (always 1 item, but respect the contract).
	maxRows := bridge.MaxVisibleRows()
	items := []string{text}
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	// Centered horizontally.
	offset := bridge.CenterXWith(utf8.RuneCountInString(text), ga)

	return style.ViewData{
		Items:       items,
		Tiers:       []tiercatalog.Tier{tiercatalog.TierSmall},
		LineOffsets: []int{offset},
		OffsetY:     0,
		Static:      true,
	}
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
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	entry := ctx.Entry(tiercatalog.TierLarge)
	ga := entry.GlyphAdvance
	if ga <= 0 {
		ga = bridge.GlyphAdvance()
	}

	var text string
	if !snapshot.HasLast {
		text = "Waiting for USB..."
	} else {
		text = coreDisplayName(snapshot.Device)
	}

	// Truncate to available character width.
	maxChars := 0
	if ga > 0 {
		maxChars = bridge.AvailableContentWidth() / ga
	}
	if maxChars > 0 {
		text = textlayout.Truncate(text, maxChars)
	}

	// Clamp to MaxVisibleRows.
	maxRows := bridge.MaxVisibleRows()
	items := []string{text}
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
	}

	// Left-aligned at content origin X.
	ox, _ := bridge.ContentOrigin()

	return style.ViewData{
		Items:       items,
		Tiers:       []tiercatalog.Tier{tiercatalog.TierLarge},
		LineOffsets: []int{ox},
		OffsetY:     0,
		Static:      true,
	}
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
