package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

type GpioSnapshot = source.GpioSnapshot
type Policy = source.Policy

// gpioRegistry is the per-mode StyleRegistry for the gpio display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
// Legacy styles (list, icons, detail, dashboard, activity) are registered first as the default.

// Persistent widget instances for the GPIO mode. These survive across frames
// and are reconfigured each frame via the Configurable interface before being
// passed to the Compositor.
//
// Framework pattern demonstrated: PersistentWidgets
var (
	persistentLED       = led.New(led.Config{State: led.Off, Brightness: -1.0, Diameter: 6, Bounds: image.Rect(0, 0, 6, 6)})
	persistentSparkline = sparkline.New(sparkline.Config{Style: sparkline.Bar, Bounds: image.Rect(0, 0, 10, 8)})
)

// MaxVisibleRows computes how many rows fit in the panel's vertical space
// given the current TextHints, minus the specified padding rows.
// Returns 0 if RowHeight is non-positive or usable space is non-positive.
//
// Framework pattern demonstrated: LayoutBridge helper — adaptive row computation
// for resolution-specific rendering decisions (label visibility, row truncation).
func MaxVisibleRows(hints textlayout.TextHints, padding int) int {
	if hints.RowHeight <= 0 {
		return 0
	}
	rows := hints.PixelHeight / hints.RowHeight
	rows -= padding
	if rows < 0 {
		return 0
	}
	return rows
}

// MaxCharsPerRow computes how many monospace glyphs fit in the panel's
// horizontal space given the current TextHints.
// Returns 0 if GlyphAdvance is non-positive or PixelWidth is non-positive.
//
// Framework pattern demonstrated: LayoutBridge helper — adaptive column computation
// for text truncation in resolution-specific styles.
func MaxCharsPerRow(hints textlayout.TextHints) int {
	if hints.GlyphAdvance <= 0 || hints.PixelWidth <= 0 {
		return 0
	}
	return hints.PixelWidth / hints.GlyphAdvance
}

// maxCharsFromBridge computes how many monospace glyphs fit using bridge metrics.
func maxCharsFromBridge(bridge layout.LayoutCalculator) int {
	if bridge.GlyphAdvance() <= 0 || bridge.AvailableContentWidth() <= 0 {
		return 0
	}
	return bridge.AvailableContentWidth() / bridge.GlyphAdvance()
}

// longestItemLen returns the length of the longest string in items.
// Returns 0 for an empty or nil slice.
func longestItemLen(items []string) int {
	maxLen := 0
	for _, s := range items {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen
}

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

// buildListInternal replicates the list-style rendering logic from BuildView.
func buildListInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	pol := snap.Policy
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// Use catalog metrics for layout.
	entry := ctx.Entry(tiercatalog.TierNormal)

	maxRows := bridge.MaxVisibleRows()
	maxChars := 0
	if entry.GlyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / entry.GlyphAdvance
	}

	visiblePins := configuredPins(snap.Pins)
	if len(visiblePins) == 0 {
		return style.ViewData{
			Items:  noConfiguredItems(len(snap.Pins), maxChars),
			Colors: nil,
			Static: true,
		}
	}
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	items := BuildItemsTruncated(visiblePins, maxChars)
	colors := buildAccentColors(visiblePins, pol)

	rowHeight := entry.RowHeight
	if rowHeight <= 0 {
		rowHeight = bridge.RowHeight()
	}
	glyphWidth := entry.GlyphWidth
	if glyphWidth <= 0 {
		glyphWidth = bridge.GlyphAdvance()
	}
	glyphHeight := entry.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = bridge.RowHeight()
	}
	sprites := BuildSpritesFromMetrics(visiblePins, rowHeight, glyphWidth, glyphHeight)

	// TextLabel rendering for list style using tier catalog for face resolution.
	face := ctx.Face("spleen", tiercatalog.TierNormal)
	if face != nil {
		// Build a TextHints-compatible struct for textlabel sprites.
		textLabelHints := textlayout.TextHints{
			PixelWidth:   bridge.AvailableContentWidth(),
			PixelHeight:  bridge.AvailableContentHeight(),
			GlyphWidth:   entry.GlyphWidth,
			GlyphHeight:  entry.GlyphHeight,
			GlyphAdvance: entry.GlyphAdvance,
			RowHeight:    entry.RowHeight,
		}
		textLabelSprites := buildTextLabelSprites(visiblePins, textLabelHints, face, colors)
		if textLabelSprites != nil {
			sprites = append(sprites, textLabelSprites...)
		}
	}

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: sprites,
	}
}

// BuildSpritesFromMetrics returns sprites for output-mode pins using provided metrics.
// This is the bridge-compatible version of BuildSprites.
func BuildSpritesFromMetrics(pins []gpiomgr.PinState, rowHeight, glyphWidth, glyphHeight int) []widgets.Sprite {
	diameter := led.DiameterForRow(rowHeight)

	ctx := widgets.SuppressionContext{
		AvailableWidth:  diameter,
		AvailableHeight: rowHeight * len(pins),
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range pins {
		state := led.Off
		if p.Level {
			state = led.On
		}
		comp.AddIf(p.Mode == gpiomgr.ModeOutput, led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(0, i*rowHeight)},
			Foreground: ColorHigh,
		}))
	}
	if len(comp.Sprites()) == 0 {
		return nil
	}
	return comp.Sprites()
}

// buildIconsInternal replicates the icons-style rendering logic from BuildView.
func buildIconsInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// Use catalog metrics for layout.
	entry := ctx.Entry(tiercatalog.TierNormal)

	sprites := BuildIconGrid(snap.Pins, bridge.AvailableContentWidth(), bridge.AvailableContentHeight(), entry.GlyphWidth, entry.GlyphHeight)
	var viewColors = buildAccentColors(snap.Pins, snap.Policy)

	return style.ViewData{
		Items:   nil,
		Colors:  viewColors,
		Sprites: sprites,
		Static:  true,
	}
}

// buildDetailInternal replicates the detail-style rendering logic from BuildView.
func buildDetailInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	pol := snap.Policy
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// Use catalog metrics for layout.
	entry := ctx.Entry(tiercatalog.TierNormal)

	// Resolve a face for sprite rendering via tier catalog.
	selectedFont := ctx.Face("spleen", tiercatalog.TierNormal)

	// Build TextHints for BuildDetailView compatibility.
	detailHints := textlayout.TextHints{
		PixelWidth:   bridge.AvailableContentWidth(),
		PixelHeight:  bridge.AvailableContentHeight(),
		GlyphWidth:   entry.GlyphWidth,
		GlyphHeight:  entry.GlyphHeight,
		GlyphAdvance: entry.GlyphAdvance,
		RowHeight:    entry.RowHeight,
	}

	sprites := BuildDetailView(snap.Pins, detailHints, pol.PinLabels, selectedFont)
	var viewColors = buildAccentColors(snap.Pins, pol)

	return style.ViewData{
		Items:   nil,
		Colors:  viewColors,
		Sprites: sprites,
		Static:  false,
	}
}

// buildDashboardInternal replicates the dashboard-style rendering logic from BuildView.
func buildDashboardInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// Use catalog metrics for layout.
	entry := ctx.Entry(tiercatalog.TierNormal)

	// Resolve a face for sprite rendering via tier catalog.
	selectedFont := ctx.Face("spleen", tiercatalog.TierNormal)

	// Build TextHints for BuildDashboardView compatibility.
	dashHints := textlayout.TextHints{
		PixelWidth:   bridge.AvailableContentWidth(),
		PixelHeight:  bridge.AvailableContentHeight(),
		GlyphWidth:   entry.GlyphWidth,
		GlyphHeight:  entry.GlyphHeight,
		GlyphAdvance: entry.GlyphAdvance,
		RowHeight:    entry.RowHeight,
	}

	vd := BuildDashboardView(snap.Pins, dashHints, selectedFont)
	return vd
}

// buildActivityInternal replicates the activity-style rendering logic from BuildView.
func buildActivityInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// Use catalog metrics for layout.
	entry := ctx.Entry(tiercatalog.TierNormal)

	// Resolve a face for sprite rendering via tier catalog.
	selectedFont := ctx.Face("spleen", tiercatalog.TierNormal)

	// Build TextHints for BuildActivityView compatibility.
	actHints := textlayout.TextHints{
		PixelWidth:   bridge.AvailableContentWidth(),
		PixelHeight:  bridge.AvailableContentHeight(),
		GlyphWidth:   entry.GlyphWidth,
		GlyphHeight:  entry.GlyphHeight,
		GlyphAdvance: entry.GlyphAdvance,
		RowHeight:    entry.RowHeight,
	}

	vd := BuildActivityView(snap.Pins, actHints, selectedFont)
	return vd
}
