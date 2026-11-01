package displaymode

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

// GraphicsExampleStyle is the template's worked example of mixing TEXT and
// GRAPHICS in one view. The other 24 styles in this package are text-only
// skeletons; this one exists to teach the single most important and least
// obvious rule about laying out a graphical display mode.
//
// # The two layout planes (read this before drawing a bar next to a label)
//
// A ViewData is rendered through TWO SEPARATE, INDEPENDENT coordinate systems,
// and the classic mistake is trying to make them line up with each other:
//
//  1. The TTY / ROW TEXT plane — ViewData.Items. The RENDERER owns vertical
//     placement here. It flows the item strings top-to-bottom using ITS OWN row
//     height (the surface font's metrics, or the hints in ViewData.RowHeights /
//     Tiers), applies ViewData.OffsetY, and picks a per-row font itself when no
//     Tiers/FontIDs are given. The style does NOT control the pixel Y of any
//     row. Think of it as a terminal: you hand it lines of text, not
//     coordinates.
//
//  2. The PIXEL-ANCHORED SPRITE plane — ViewData.Sprites. These are drawn at the
//     absolute panel pixel coordinates the style sets, verbatim. The style owns
//     every coordinate.
//
// THE TRAP: emitting a caption as an Items row AND a bar as a Sprite, then
// expecting row N of the text to sit at the same Y as a sprite you placed at
// "origin + N * rowHeight". Those two only coincide if the row height you
// assumed equals the one the renderer actually used — and they often do not.
// On panels that report no glyph metrics (e.g. the PNG snapshot panel), a
// style's LayoutBridge falls back to small default metrics (~6x10) while the
// renderer draws the text at the font catalog's real, larger size. The text
// then spreads by one measure and the sprites bunch by another, and the view
// desyncs into an unreadable clump. No amount of tweaking RowHeights/Tiers
// fixes it, because it is a PLANE-CATEGORY mismatch, not an arithmetic error.
//
// THE RULE: a view that mixes text and graphics should commit ENTIRELY to the
// pixel plane. Render the labels as sprites too, using the textlabel widget,
// and leave ViewData.Items empty so the TTY layer is never engaged. Then the
// label and its bar are positioned by the same arithmetic in the same
// coordinate system and cannot drift apart. This style does exactly that.
//
// Two supporting rules this style also demonstrates:
//   - Get the real font face (ctx.Face) and metrics (ctx.Entry) for the tier;
//     do not trust LayoutBridge glyph metrics when the panel reports none.
//   - The font's row height is fixed and does NOT scale with panel size. To
//     FILL a panel with N rows, distribute by availableHeight/N (a per-row
//     pitch), not by the font's row height — otherwise N small rows float as a
//     tiny block in a large panel.
type GraphicsExampleStyle struct{}

func (s GraphicsExampleStyle) Name() string { return "color-240x240-graphics" }

// Requirements: a roomy color panel, where a graphical view makes sense.
func (s GraphicsExampleStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast}
}

func (s GraphicsExampleStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

// Build renders a small "load wall" — labelled bars that fill the panel — using
// only the pixel-anchored sprite plane. Snapshot in this template is empty, so
// the values below are illustrative demo data; a real mode would read them from
// its Snapshot.
func (s GraphicsExampleStyle) Build(snapshot Snapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	// Illustrative demo values (a real mode reads these from its Snapshot).
	values := []float64{0.62, 0.18, 0.91, 0.07}

	bridge := layout.NewLayoutBridge(ctx.Hints(), layout.BridgeConfig{})
	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()

	// Degenerate panel: no room for graphics. Fall back to the TTY plane with a
	// single text row so the view still produces renderable output. (This is the
	// one place a graphics view legitimately uses Items.)
	if width < 8 || availH < 8 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	// Real font face + metrics the labels draw with — NOT the bridge defaults,
	// which are wrong on panels that report no glyph metrics.
	face := ctx.Face("", tiercatalog.TierNormal)
	entry := ctx.Entry(tiercatalog.TierNormal)
	rowH := entry.RowHeight
	glyphAdvance := entry.GlyphAdvance
	if rowH < 1 {
		rowH = bridge.RowHeight()
	}
	if glyphAdvance < 1 {
		glyphAdvance = bridge.GlyphAdvance()
	}

	n := len(values)
	minPitch := rowH + 2
	if availH/minPitch < n {
		n = availH / minPitch
	}
	if n < 1 {
		// Too small for the graphical view: fall back to the TTY plane with a
		// single text row. (This is the one place Items is used — when there is
		// no room for graphics at all.)
		return style.ViewData{Items: []string{"(too small)"}}
	}

	// Distribute rows across the FULL height so the view fills the panel.
	rowPitch := availH / n
	blockH := n * rowPitch
	topPad := 0
	if availH > blockH {
		topPad = (availH - blockH) / 2
	}
	baseY := oy + topPad

	labelW := len("0:100%") * glyphAdvance
	if labelW > width/2 {
		labelW = width / 2
	}
	barX := ox + labelW + 2
	barRight := ox + width
	barH := rowPitch - 4
	if barH < 2 {
		barH = 2
	}

	var vd style.ViewData // Items deliberately left empty — pixel plane only.

	for i := 0; i < n; i++ {
		slotTop := baseY + i*rowPitch
		load := values[i]
		fg := demoLoadColor(load)

		// Label as a SPRITE (textlabel), not an Items row — so it shares the
		// bar's coordinate system.
		labelTop := slotTop + (rowPitch-rowH)/2
		if lbl := textlabel.Render(textlabel.Config{
			Text:       fmt.Sprintf("%d:%d%%", i, int(load*100+0.5)),
			Bounds:     image.Rect(ox, labelTop, ox+labelW, labelTop+rowH),
			Font:       face,
			Alignment:  textlabel.Left,
			Foreground: fg,
		}); lbl != nil {
			vd.Sprites = append(vd.Sprites, *lbl)
		}

		// Bar as a SPRITE on the same row, positioned by the same arithmetic.
		barTop := slotTop + (rowPitch-barH)/2
		if barRight-barX >= 4 {
			if bar := progressbar.Render(progressbar.Config{
				Style:       progressbar.Linear,
				Value:       load,
				Bounds:      image.Rect(barX, barTop, barRight, barTop+barH),
				Foreground:  fg,
				Background:  color.RGBA{R: 40, G: 40, B: 40, A: 255},
				RoundedCaps: true,
			}); bar != nil {
				vd.Sprites = append(vd.Sprites, *bar)
			}
		}
	}

	return vd
}

// demoLoadColor maps load in [0,1] to a green→amber→red ramp, illustrating that
// color can carry meaning (here: how hard something is working).
func demoLoadColor(load float64) color.RGBA {
	if load < 0 {
		load = 0
	}
	if load > 1 {
		load = 1
	}
	if load <= 0.5 {
		t := load / 0.5
		return color.RGBA{R: uint8(0 + 230*t), G: uint8(200 - 30*t), B: 0, A: 255}
	}
	t := (load - 0.5) / 0.5
	return color.RGBA{R: 230, G: uint8(170 - 130*t), B: uint8(0 + 30*t), A: 255}
}
