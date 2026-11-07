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
	"github.com/databeast/cyberhud/display/widgets/textlabel"
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

// maxCharsFromBridge computes how many monospace glyphs of font m fit across the
// content width.
func maxCharsFromBridge(bridge layout.LayoutCalculator, m layout.Metrics) int {
	return bridge.MaxChars(m)
}

// buildListInternal replicates the list-style rendering logic from BuildView.
// buildListInternal renders the GPIO LIVE-FUNCTIONS view — the one idea this
// mode exists to deliver on a headless box's only screen: "what is the header
// actually doing, right now?"
//
// It shows one row per live function (see source.LiveFunctions): standalone
// driven/read pins individually with their level, each bus (SPI0, I2C1, …)
// collapsed to a single "up" row, and unconfigured pins absent entirely. The
// resting state is a calm legible list; because the screen always reflects
// current truth, a line going dark or an output flipping simply shows — the
// human's memory supplies the delta, with no animation.
//
// It renders ENTIRELY on the pixel (sprite) plane: one indicator dot and one
// text label per function, positioned by the SAME row arithmetic so they can
// never desync. ViewData.Items is left empty so the renderer's TTY text flow is
// never engaged (mixing the two planes is what smeared earlier versions). Rows
// distribute across the full content height so a tall strip fills top-to-bottom
// instead of clumping, and text truncates to the width budget so nothing runs
// off a narrow panel.
func buildListInternal(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	pol := snap.Policy
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	face := ctx.Face("spleen", tiercatalog.TierNormal)

	fns := source.LiveFunctions(snap.Pins)

	// Idle fallback: nothing on the header is doing anything. A single centered
	// line via the TTY path (no sprites) — honest and does not engage the
	// two-plane machinery.
	if len(fns) == 0 {
		maxChars := maxCharsFromBridge(bridge, ctx.DefaultMetrics())
		msg := "GPIO idle"
		if maxChars > 0 {
			msg = textlayout.Truncate(msg, maxChars)
		}
		return style.ViewData{Items: []string{msg}}
	}

	originX, originY := bridge.ContentOrigin()
	availW := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()

	// Minimum legible row pitch from the font; grow it to fill the panel when
	// there are few rows, so a tall strip distributes rows instead of clumping.
	glyphHeight := entry.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = textlayout.RowHeight
	}
	minPitch := glyphHeight + 2
	if minPitch < 8 {
		minPitch = 8
	}
	count := len(fns)
	rowPitch := availH / count
	if rowPitch < minPitch {
		rowPitch = minPitch
	}

	// Drop trailing rows that would not fit, so no row is clipped at the bottom.
	maxRows := availH / rowPitch
	if maxRows < 1 {
		maxRows = 1
	}
	if count > maxRows {
		fns = fns[:maxRows]
		count = maxRows
	}

	// Dot column on the left; text to its right.
	diameter := led.DiameterForRow(rowPitch)
	if diameter < 4 {
		diameter = 4
	}
	// Text sits right of the dot column, except on e-ink where the dot is
	// suppressed and text reclaims the full width from the left margin.
	isEink := ctx.Cap() == style.MonoSlow
	textGap := diameter + 4
	if isEink {
		textGap = 0
	}
	textX := originX + textGap
	textWidth := availW - textGap
	glyphAdvance := entry.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}
	maxChars := 0
	if glyphAdvance > 0 && textWidth > 0 {
		maxChars = textWidth / glyphAdvance
	}

	accent := resolveFGColor(pol.FGColor)
	dimmed := dimFGColor(pol.FGColor)
	colorOn := pol.Color

	// On e-ink the indicator dots are dropped (LEDs are not eink-safe): the
	// text row already carries the full "identity + state" so the dot is
	// redundant there. SuppressOnEink enforces this via the compositor.
	suppCtx := widgets.SuppressionContext{IsEink: isEink, AvailableWidth: availW, AvailableHeight: availH}
	comp := widgets.NewCompositor(suppCtx, widgets.SuppressOnEink())

	for i, fn := range fns {
		rowY := originY + i*rowPitch

		// A function is "lit" when it is live: a bus is always up, a standalone
		// pin is lit when HIGH. Lit → accent/On, unlit → dimmed/Off.
		lit := fn.Kind == source.KindBus || fn.Level
		fg := dimmed
		state := led.Off
		if lit {
			fg = accent
			state = led.On
		}

		// Indicator dot, vertically centered in the row pitch.
		dotY := rowY + (rowPitch-diameter)/2
		if dotY < rowY {
			dotY = rowY
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(originX, dotY)},
			Foreground: fg,
		}))

		// Text label on the SAME row arithmetic, truncated to the width budget.
		if face != nil {
			label := fn.Row()
			if maxChars > 0 && len(label) > maxChars {
				label = label[:maxChars]
			}
			labelFG := fg
			if !colorOn {
				labelFG = accent // mono panels ignore color; keep bright text
			}
			comp.Add(textlabel.New(textlabel.Config{
				Text:       label,
				Bounds:     image.Rect(textX, rowY, originX+availW, rowY+rowPitch),
				Font:       face,
				Alignment:  textlabel.Left,
				Foreground: labelFG,
			}))
		}
	}

	return style.ViewData{
		Sprites: comp.Sprites(),
		
	}
}
