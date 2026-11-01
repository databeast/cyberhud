package styles

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

// coreTrackBG is the empty-track color behind each core's bar so an idle core
// still reads as "a core at zero" rather than a blank row.
var coreTrackBG = color.RGBA{R: 40, G: 40, B: 40, A: 255}

// buildCoresStyle is the shared "cores" variant builder — the load wall: one
// bar per CPU core, every core always shown, each bar's fill and color carrying
// the load (green→amber→red on color panels).
//
// This view is rendered ENTIRELY in the pixel-anchored sprite plane: each
// core's label is a textlabel widget sprite and its bar is a progressbar
// sprite, both positioned by the same per-row pixel arithmetic. Items is left
// empty so the renderer's TTY/row text layout is never engaged — mixing the
// row-flowed text plane with absolute-pixel sprites is what made earlier
// attempts desync and clump. One plane, one coordinate system, no tofu.
func buildCoresStyle(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	sample := snapshot.CPUSample
	if sample == nil {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	isColor := isColorPanel(ctx)

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()

	// Real font face + metrics the labels will draw with (pixel-anchored, so we
	// own the geometry; no dependence on the renderer's row flow).
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
	coreCount := len(sample)

	// Distribute cores across the FULL available height so the load wall fills
	// the panel instead of floating as a small block (the font's own row height
	// is a fixed ~10px and does not scale with the panel). Each core gets an
	// equal vertical slice; the minimum slice is the font row height so labels
	// stay legible, and if even that will not fit we cap the core count.
	minPitch := rowH + 2
	count := coreCount
	if availH/minPitch < count {
		count = availH / minPitch
	}
	if count < 1 {
		return style.ViewData{Items: []string{"(too small)"}}
	}
	rowPitch := availH / count
	if rowPitch < minPitch {
		rowPitch = minPitch
	}

	// Centre the filled block (rowPitch*count may be slightly under availH due
	// to integer division).
	blockH := count * rowPitch
	topPad := 0
	if availH > blockH {
		topPad = (availH - blockH) / 2
	}
	baseY := oy + topPad

	// Label column width: "N:100%" = up to 6 chars at this resolution's advance.
	labelChars := len(strconv.Itoa(coreCount-1)) + len(":100%")
	labelW := labelChars * glyphAdvance
	if labelW > width/2 {
		labelW = width / 2
	}

	barX := ox + labelW + 2
	barRight := ox + width

	// Bar fills most of its row slice vertically, with a small gap between rows.
	barH := rowPitch - 4
	if barH < 2 {
		barH = 2
	}
	if barH > 28 {
		barH = 28 // keep bars from looking absurdly tall on huge panels
	}

	var vd style.ViewData

	for i := 0; i < count; i++ {
		slotTop := baseY + i*rowPitch
		// Centre the bar (and label) vertically within the slice.
		barTop := slotTop + (rowPitch-barH)/2
		labelTop := slotTop + (rowPitch-rowH)/2
		load := sample[i]
		col := loadColor(load, isColor)

		// Pixel-anchored label sprite.
		pct := int(load*100 + 0.5)
		labelSprite := textlabel.Render(textlabel.Config{
			Text:       fmt.Sprintf("%d:%d%%", i, pct),
			Bounds:     image.Rect(ox, labelTop, ox+labelW, labelTop+rowH),
			Font:       face,
			Alignment:  textlabel.Left,
			Foreground: col,
		})
		if labelSprite != nil {
			labelSprite.Label = "core-label-" + strconv.Itoa(i)
			vd.Sprites = append(vd.Sprites, *labelSprite)
		}

		// Pixel-anchored bar sprite on the same row.
		if barRight-barX >= 4 {
			bar := progressbar.Render(progressbar.Config{
				Style:       progressbar.Linear,
				Value:       load,
				Bounds:      image.Rect(barX, barTop, barRight, barTop+barH),
				Foreground:  col,
				Background:  coreTrackBG,
				RoundedCaps: true,
			})
			if bar != nil {
				bar.Label = "core-bar-" + strconv.Itoa(i)
				vd.Sprites = append(vd.Sprites, *bar)
			}
		}
	}

	return vd
}

// uniformRowHeights returns a slice of n entries each equal to rowHeight.
func uniformRowHeights(n, rowHeight int) []int {
	if n <= 0 || rowHeight <= 0 {
		return nil
	}
	h := make([]int, n)
	for i := range h {
		h[i] = rowHeight
	}
	return h
}
