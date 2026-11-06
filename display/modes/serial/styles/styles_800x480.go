package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

// Serial style declarations for 800x480 panels (large landscape). Covers both
// slow-refresh e-paper and fast-refresh variants.
//
// Most entries are thin stamps over the shared engine. ColorFast800x480Style is
// the exception: it carries a bespoke Params.BuildFn (posterPanel800x480) to
// show how a single panel takes full control of its layout. When the shared
// engine isn't right for a specific display, drop the stamp helper for a full
// def literal with Params{BuildFn: ...} and write the placement beside it, as
// below.

var (
	MonoSlow800x480Style      = monoSlowStyle(800, 480)
	GrayscaleSlow800x480Style = grayscaleSlowStyle(800, 480)
	ColorSlow800x480Style     = colorSlowStyle(800, 480)

	MonoFast800x480Style      = monoFastStyle(800, 480)
	GrayscaleFast800x480Style = grayscaleFastStyle(800, 480)
)

// ColorFast800x480Style overrides the shared engine with a bespoke poster
// layout for the large 800x480 color TFT. The extra room lets the header read
// as a title at a larger tier with a generous status line, instead of the
// compact header the shared engine uses on every size. This is the worked
// example of per-panel layout override — the rest of the grid delegates to the
// engine; this cell owns its placement.
var ColorFast800x480Style = def{
	name: "color-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast, MinRows: 2, MinCharsPerLine: 8},
	p:    Params{BuildFn: posterPanel800x480},
}

// posterPanel800x480 is the bespoke layout for the large color TFT. It is a
// pure placement function — it reads the already-settled snapshot and composites
// sprites, holding no per-frame state — exactly like the shared engine, but
// hand-tuned for this one panel. Copy this shape to override any other cell.
func posterPanel800x480(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()

	// A larger face for the title band; the body stays at the normal tier.
	titleFace := ctx.Face(serialFontFamily, tiercatalog.TierLarge)
	bodyFace := ctx.Face(serialFontFamily, tiercatalog.TierNormal)
	hints = textlayout.WithFont(hints, bodyFace)

	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 3})
	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	if cw <= 0 || chh <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: false}
	}

	// Raw mode is a policy toggle, honoured here as everywhere: bytes only.
	if pol.Raw {
		return rawView(snap, ctx, bridge, false)
	}

	titleEntry := ctx.Entry(tiercatalog.TierLarge)
	bodyEntry := ctx.Entry(tiercatalog.TierNormal)
	titleH := titleEntry.RowHeight
	bodyH := bodyEntry.RowHeight
	bodyAdvance := bodyEntry.GlyphAdvance
	if titleH <= 0 {
		titleH = 24
	}
	if bodyH <= 0 {
		bodyH = textlayout.RowHeight
	}
	if bodyAdvance <= 0 {
		bodyAdvance = textlayout.GlyphAdvance
	}

	ox, oy := bridge.ContentOrigin()
	comp := widgets.NewCompositor(widgets.SuppressionContext{
		AvailableWidth:  cw,
		AvailableHeight: chh,
	})

	y := oy

	// --- Title band: a large connection LED + port/baud as a heading ---
	ledSize := titleH
	ledState, ledColor := led.Off, ledRed
	if snap.Connected {
		ledState, ledColor = led.On, ledGreen
	}
	comp.Add(led.New(led.Config{
		State:      ledState,
		Brightness: -1.0,
		Diameter:   ledSize,
		Bounds:     image.Rect(ox, y, ox+ledSize, y+ledSize),
		Foreground: ledColor,
	}))
	comp.Add(textlabel.New(textlabel.Config{
		Text:       compactHeader(snap),
		Bounds:     image.Rect(ox+ledSize+6, y, ox+cw, y+titleH),
		Font:       titleFace,
		Foreground: colorSerialHeader,
	}))
	y += titleH + 4

	// --- Status line ---
	stateStr, stateCol := statusText(snap)
	comp.Add(textlabel.New(textlabel.Config{
		Text:       stateStr,
		Bounds:     image.Rect(ox, y, ox+cw, y+bodyH),
		Font:       bodyFace,
		Foreground: stateCol,
	}))
	y += bodyH + 6

	// --- Output body, ANSI-colored, filling the remainder ---
	budget := 0
	if bodyH > 0 {
		budget = (oy + chh - y) / bodyH
	}
	maxChars := cw / bodyAdvance
	placeDataLines(comp, snap, ox, y, cw, bodyH, bodyAdvance, bodyFace, budget, maxChars, true)

	return style.ViewData{
		Sprites: comp.Sprites(),
		Static:  false,
		Cursor:  -1,
	}
}
