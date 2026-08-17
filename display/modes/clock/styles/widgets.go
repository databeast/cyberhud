package styles

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
)

// Persistent widget instances for the clock mode. These survive across frames
// and are reconfigured each frame via the Configurable interface before being
// passed to the Compositor.
var (
	clockBorderFrameWidget = borderframe.New(borderframe.Config{Bounds: image.Rect(0, 0, 32, 32), Theme: "rounded"})
	clockLEDWidget         = led.New(led.Config{State: led.Off, Diameter: 6, Brightness: -1.0, Bounds: image.Rect(0, 0, 6, 6)})
	clockProgressBarWidget = progressbar.New(progressbar.Config{Bounds: image.Rect(0, 0, 10, 4)})
	clockSparklineWidget   = sparkline.New(sparkline.Config{Style: sparkline.Bar, Bounds: image.Rect(0, 0, 10, 8)})
)

// resolveBorderColor resolves the border color configuration to a concrete RGBA value.
//   - fgColor "none" (mono panel) → opaque white {255, 255, 255, 255} unconditionally
//   - "auto" → delegate to resolveFGColor(fgColor) (inherits the clock's foreground color)
//   - "none" → opaque white {255, 255, 255, 255}
//   - named color → delegate to resolveFGColor(borderColor)
func resolveBorderColor(borderColor, fgColor string) color.RGBA {
	// Monochrome panels always get white border tint regardless of borderColor setting.
	if fgColor == "none" {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	switch borderColor {
	case "auto":
		return resolveFGColor(fgColor)
	case "none":
		return color.RGBA{255, 255, 255, 255}
	default:
		return resolveFGColor(borderColor)
	}
}

// appendClockWidgets appends border, LED, progress bar, and sparkline sprites
// to the ViewData. Called from within each style's Build() method.
//
// The bridge parameter is the style's self-constructed LayoutBridge — this
// function does NOT construct its own layout. All positioning uses
// layout.ContentOrigin() as base.
//
// Border logic: showBorder is computed from policy and panel dimensions.
// Mono auto-enablement: when ShowBorder is false and no explicit user override
// exists, auto-enable for mono panels (FGColor="none") with dimensions ≥ 16×16.
//
// Uses the Compositor pattern with SuppressOnEink rule for LED/progressbar.
// The sparkline daybar is not suppressed on e-ink (changes at minute granularity).
func appendClockWidgets(vd *style.ViewData, bridge layout.LayoutCalculator, hints textlayout.TextHints, p source.Policy, now time.Time, isColor bool, reqs style.SurfaceRequirements) {
	effectiveW := bridge.AvailableContentWidth()
	effectiveH := bridge.AvailableContentHeight()

	// --- Border decision ---
	showBorder := p.ShowBorder && hints.PixelWidth >= 16 && hints.PixelHeight >= 16

	// Panel-aware auto-enablement for mono panels: when ShowBorder is false
	// (from DefaultPolicy) and no explicit user override exists, auto-enable
	// the border for mono panels (FGColor="none") with dimensions ≥ 16×16.
	if !showBorder && !isColor && !p.ShowBorderExplicit && hints.PixelWidth >= 16 && hints.PixelHeight >= 16 {
		showBorder = true
	}

	if showBorder {
		// Border bounds are panel-covering (from 0,0 to full panel dimensions).
		var tint color.RGBA
		if !isColor {
			tint = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		} else {
			tint = resolveBorderColor(p.BorderColor, p.FGColor)
		}
		clockBorderFrameWidget.(widgets.Configurable).Configure(borderframe.Config{
			Theme:     "rounded",
			ColorTint: tint,
			Bounds:    image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		})
	}

	// Construct SuppressionContext from bridge geometry.
	isEink := reqs.Capability <= style.GrayscaleSlow
	ctx := widgets.SuppressionContext{
		IsEink:          isEink,
		AvailableWidth:  effectiveW,
		AvailableHeight: effectiveH,
	}

	// Create Compositor with SuppressOnEink rule.
	comp := widgets.NewCompositor(ctx, widgets.SuppressOnEink())

	// --- Border frame widget ---
	comp.AddIf(showBorder, clockBorderFrameWidget)

	// --- LED widget ---
	showLED := p.ShowLED && !p.ShowSeconds && effectiveW >= 6
	if showLED {
		showLED = configureLEDWidget(vd, now, effectiveW, effectiveH, p, hints)
	}
	comp.AddIf(showLED, clockLEDWidget)

	// --- Progress bar widget ---
	showProgressBar := p.SecondsBar != "none" && effectiveH >= 48
	if showProgressBar {
		configureProgressBarWidget(now, effectiveW, effectiveH, p, isColor)
	}
	comp.AddIf(showProgressBar, clockProgressBarWidget)

	// --- Sparkline daybar widget ---
	showDaybar := p.ShowDaybar && hints.PixelHeight >= 128
	if showDaybar {
		configureSparklineWidget(now, effectiveW, hints.PixelHeight, p, isColor)
	}

	// Collect sprites from the suppression-aware compositor.
	vd.Sprites = append(vd.Sprites, comp.Sprites()...)

	// Sparkline uses a separate compositor without SuppressOnEink since it is
	// not suppressed on e-ink panels (changes at minute granularity only).
	if showDaybar {
		sparkComp := widgets.NewCompositor(ctx)
		sparkComp.Add(clockSparklineWidget)
		vd.Sprites = append(vd.Sprites, sparkComp.Sprites()...)
	}
}

// configureLEDWidget reconfigures the persistent LED widget instance for the
// current frame. Returns false if the LED should be suppressed due to text
// collision. The LED is positioned at the top-right of the effective content area.
func configureLEDWidget(vd *style.ViewData, now time.Time, effectiveW, effectiveH int, p source.Policy, hints textlayout.TextHints) bool {
	const ledSize = 6

	// LED position: top-right of effective area.
	ledX := effectiveW - ledSize
	ledY := 0

	// LED bounding rectangle (relative to content area origin).
	ledRect := image.Rect(ledX, ledY, ledX+ledSize, ledY+ledSize)

	// Check intersection with text rows. If any row overlaps, suppress LED.
	if ledIntersectsTextRows(ledRect, vd, hints) {
		return false
	}

	// Determine LED state based on current second.
	sec := now.Second()
	state := led.On
	if sec%2 != 0 {
		state = led.Off
	}

	// Reconfigure the persistent LED widget with current frame parameters.
	clockLEDWidget.(widgets.Configurable).Configure(led.Config{
		State:      state,
		Brightness: -1.0,
		Diameter:   ledSize,
		Bounds:     image.Rect(ledX, ledY, ledX+ledSize, ledY+ledSize),
	})

	return true
}

// configureProgressBarWidget reconfigures the persistent progress bar widget
// for the current frame. Positioned at bottom (horizontal) or bottom-right (pie)
// of the effective content area.
func configureProgressBarWidget(now time.Time, effectiveW, effectiveH int, p source.Policy, isColor bool) {
	second := now.Second()
	value := float64(second) / 60.0

	// Determine foreground color based on panel type.
	var fg color.RGBA
	if isColor {
		fg = resolveFGColor(p.FGColor)
	} else {
		fg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	var cfg progressbar.Config

	switch p.SecondsBar {
	case "horizontal":
		// Position at bottom of effective area, height=4, width=effectiveWidth.
		cfg = progressbar.Config{
			Style:      progressbar.Linear,
			Value:      value,
			Bounds:     image.Rect(0, effectiveH-4, effectiveW, effectiveH),
			Foreground: fg,
			Background: bg,
		}
	case "pie":
		// Position at bottom-right corner, diameter=16.
		diameter := 16
		cfg = progressbar.Config{
			Style:      progressbar.Pie,
			Value:      value,
			Bounds:     image.Rect(effectiveW-diameter, effectiveH-diameter, effectiveW, effectiveH),
			Foreground: fg,
			Background: bg,
		}
	default:
		return
	}

	clockProgressBarWidget.(widgets.Configurable).Configure(cfg)
}

// configureSparklineWidget reconfigures the persistent sparkline widget for the
// current frame. Positioned at the bottom of the panel area.
func configureSparklineWidget(now time.Time, effectiveWidth, panelHeight int, p source.Policy, isColor bool) {
	// Compute timezone-adjusted time for data point.
	t := source.ApplyTimezone(now, p.Timezone)
	hour := t.Hour()
	minute := t.Minute()

	// Data: fraction of day elapsed as single-element array.
	dayProgress := float64(hour*60+minute) / 1440.0
	data := []float64{dayProgress}

	// Determine sparkline foreground color.
	var fg color.RGBA
	if isColor {
		accent := resolveFGColor(p.FGColor)
		// 30% brightness: each RGB channel multiplied by 0.30, truncated to int.
		fg = color.RGBA{
			R: uint8(int(float64(accent.R) * 0.30)),
			G: uint8(int(float64(accent.G) * 0.30)),
			B: uint8(int(float64(accent.B) * 0.30)),
			A: 255,
		}
	} else {
		// Monochrome: native foreground (white).
		fg = color.RGBA{255, 255, 255, 255}
	}

	// Position: bottom of panel minus 12px for y, with specified dimensions.
	y := panelHeight - 12
	bounds := image.Rect(0, y, effectiveWidth, y+8)

	clockSparklineWidget.(widgets.Configurable).Configure(sparkline.Config{
		Data:       data,
		Style:      sparkline.Bar,
		Bounds:     bounds,
		Foreground: fg,
		Background: color.RGBA{0, 0, 0, 0}, // Transparent background.
	})
}

// ledIntersectsTextRows checks whether the LED's bounding rectangle overlaps
// with any text row's bounding rectangle. Row positions are derived from the
// ViewData's OffsetY, per-row OffsetXs, and font metrics from the tier catalog via hints.
func ledIntersectsTextRows(ledRect image.Rectangle, vd *style.ViewData, hints textlayout.TextHints) bool {
	if len(vd.Items) == 0 {
		return false
	}

	// Use tier catalog metrics from hints for all rows (single-tier mode).
	metrics := defaultMetrics(hints)

	spacing := 0
	if len(vd.Items) > 1 {
		spacing = 1
	}

	// Check each row's bounding rectangle against the LED rect.
	currentY := vd.OffsetY
	for i, item := range vd.Items {
		rowX := 0
		if i < len(vd.LineOffsets) {
			rowX = vd.LineOffsets[i]
		}

		rowWidth := metrics.GlyphAdvance * len(item)
		rowHeight := metrics.RowHeight

		rowRect := image.Rect(rowX, currentY, rowX+rowWidth, currentY+rowHeight)

		if ledRect.Overlaps(rowRect) {
			return true
		}

		currentY += rowHeight + spacing
	}

	return false
}

// defaultMetrics returns fallback font metrics derived from TextHints
// when the font ID cannot be resolved from the registry.
func defaultMetrics(hints textlayout.TextHints) font.Metrics {
	return font.Metrics{
		GlyphWidth:   hints.GlyphWidth,
		GlyphHeight:  hints.GlyphHeight,
		GlyphAdvance: hints.GlyphAdvance,
		RowHeight:    hints.RowHeight,
	}
}
