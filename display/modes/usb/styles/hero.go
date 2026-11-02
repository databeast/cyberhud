package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/usb/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/textbox"
)

// ──────────────────────────────────────────────────────────────────────────────
// buildHero — the "instant device-type recognition" layout for large panels.
//
// The device CLASS is the drama: the moment you plug something into the bench,
// the panel throws up a big, class-colored glyph (disk / keyboard / headset /
// router …) before you have read a single word. The device name is the title,
// a thin class-colored accent bar underlines the identity, and the VID:PID,
// bus/dev, serial and status read as quiet supporting rows beneath.
//
// Everything is composited as absolute-positioned sprites with an empty
// ViewData.Items, committing the whole view to the single sprite plane — the
// renderer's row-text plane and the pixel-sprite plane use different vertical
// metrics, and mixing them desyncs (see the mode docs). One plane, no desync.
//
// Mono safety: on 1-bit / grayscale panels every element renders full white
// (the class fills would quantize to black and vanish), and the accent bar is a
// thin white rule. Color is a color-panel enhancement, never load-bearing.
// ──────────────────────────────────────────────────────────────────────────────

func buildHero(snapshot source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 6})
	ox, oy := bridge.ContentOrigin()
	availW := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()

	colorPanel := isColorPanel(ctx)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Not-yet-seen state: a single centered line, no hero.
	if !snapshot.HasLast {
		sp := centeredLine(ctx, "Waiting for USB...", tiercatalog.TierLarge,
			image.Rect(ox, oy, ox+availW, oy+availH), white)
		vd := style.ViewData{Static: true, PaddingPct: 6}
		if sp != nil {
			vd.Sprites = append(vd.Sprites, *sp)
		}
		return vd
	}

	var sprites []widgets.Sprite

	// ── Hero glyph ──
	// Scale the 24px class glyph by the largest integer factor that fits ~45% of
	// the panel's smaller dimension, clamped to [1,3]. Integer nearest-neighbor
	// keeps a clean bitmap crisp (just bigger pixels), unlike smooth magnification.
	classColor := classHeroColor(snapshot.Device.DeviceClass)
	glyphColor := white
	if colorPanel {
		glyphColor = classColor
	}
	iconRune := ClassIcon(snapshot.Device.DeviceClass)
	smaller := hints.PixelWidth
	if hints.PixelHeight < smaller {
		smaller = hints.PixelHeight
	}
	scale := smaller * 45 / 100 / 24
	if scale < 1 {
		scale = 1
	}
	if scale > 3 {
		scale = 3
	}
	glyphSprite := scaledClassGlyph(iconRune, scale, glyphColor)

	// ── Layout: glyph top-centered, title under it, accent bar, detail rows ──
	// Compute the stack height to vertically center the whole group.
	glyphH := 0
	glyphW := 0
	if glyphSprite != nil {
		b := glyphSprite.Image.Bounds()
		glyphH = b.Dy()
		glyphW = b.Dx()
	}

	titleEntry := ctx.Entry(tiercatalog.TierLarge)
	titleH := titleEntry.RowHeight
	if titleH <= 0 {
		titleH = 18
	}
	const accentGap = 6
	const accentH = 3
	const glyphTitleGap = 8
	const titleDetailGap = 8

	smallEntry := ctx.Entry(tiercatalog.TierSmall)
	detailRowH := smallEntry.RowHeight
	if detailRowH <= 0 {
		detailRowH = 10
	}

	details := detailRows(snapshot)
	detailsH := len(details) * detailRowH

	// Use the ACTUAL title sprite height (face RowHeight), not the catalog Entry
	// estimate, so the centering reflects what is really drawn.
	stackH := glyphH + glyphTitleGap + titleH + accentGap + accentH + titleDetailGap + detailsH
	top := oy
	if stackH < availH {
		top = oy + (availH-stackH)/2
	}

	cx := ox + availW/2
	y := top

	// Glyph.
	if glyphSprite != nil {
		glyphSprite.Position = image.Point{X: cx - glyphW/2, Y: y}
		sprites = append(sprites, *glyphSprite)
		y += glyphH + glyphTitleGap
	}

	// Title (device name), fitted to the largest tier whose whole string fits.
	name := coreDisplayName(snapshot.Device)
	titleSprite, usedTitleH := heroTitle(ctx, name, availW, image.Point{X: ox, Y: y}, availW, white)
	if titleSprite != nil {
		sprites = append(sprites, *titleSprite)
	}
	y += usedTitleH + accentGap

	// Accent bar — class-colored on color panels, white rule on mono.
	accentColor := white
	if colorPanel {
		accentColor = classColor
	}
	// Bar width tracks the title emphasis: a fixed fraction of content width,
	// centered, so it reads as a deliberate underline rather than a full rule.
	barW := availW * 55 / 100
	if barW < 24 {
		barW = 24
	}
	barX := cx - barW/2
	accentSprite := fillRectSprite(image.Rect(barX, y, barX+barW, y+accentH), accentColor, "usb-accent")
	if accentSprite != nil {
		sprites = append(sprites, *accentSprite)
	}
	y += accentH + titleDetailGap

	// Detail rows — quiet supporting cast. Dimmed on color panels for hierarchy;
	// full white on mono (Dim(white) quantizes to black and would vanish).
	detailColor := white
	if colorPanel {
		detailColor = sharedcolor.Dim(white)
	}
	for _, text := range details {
		if y+detailRowH > oy+availH {
			break // Ran out of vertical room; drop remaining detail rows.
		}
		sp := centeredLine(ctx, text, tiercatalog.TierSmall,
			image.Rect(ox, y, ox+availW, y+detailRowH), detailColor)
		if sp != nil {
			sprites = append(sprites, *sp)
		}
		y += detailRowH
	}

	return style.ViewData{
		Sprites:    sprites,
		Static:     true,
		PaddingPct: 6,
	}
}

// buildCompactBadge renders a compact "badge": a scaled class glyph stacked
// above the device name (and optional sub-line), vertically centered as one
// group. It is the shared body for the narrow/minimal paths — which land on
// tall, width-starved portrait strips where a lone centered name floats in a
// mostly-empty panel. The glyph both fills that void and delivers the mode's
// device-type recognition promise even on the small panels.
//
// The glyph is included only when the panel has vertical headroom for it
// (short landscape strips skip it and fall back to name/sub-line only). Returns
// a sprite-plane ViewData (empty Items) so it composes on one plane.
func buildCompactBadge(snapshot source.Snapshot, ctx style.StyleContext, name, subLine string) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 6})
	ox, oy := bridge.ContentOrigin()
	availW := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()

	colorPanel := isColorPanel(ctx)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Fit the name to the largest tier that fits the (narrow) width.
	nameTier := tiercatalog.TierSmall
	for _, t := range nameTierLadder {
		ga := ctx.Entry(t).GlyphAdvance
		if ga <= 0 {
			continue
		}
		if len([]rune(name))*ga <= availW {
			nameTier = t
			break
		}
	}
	nameEntry := ctx.Entry(nameTier)
	nameH := nameEntry.RowHeight
	if nameH <= 0 {
		nameH = 12
	}

	subEntry := ctx.Entry(tiercatalog.TierSmall)
	subH := subEntry.RowHeight
	if subH <= 0 {
		subH = 10
	}

	const glyphGap = 6
	const subGap = 3

	// Glyph scale: fit the width comfortably (glyph ~70% of content width) but
	// never so tall it crowds out the name. Clamp to [1,3]; drop the glyph
	// entirely when there is not enough vertical room for glyph + name together.
	glyphNative := 24
	scale := availW * 70 / 100 / glyphNative
	if scale < 1 {
		scale = 1
	}
	if scale > 3 {
		scale = 3
	}
	glyphH := glyphNative * scale
	// Need room for glyph + gap + name (+ sub). If not, drop the glyph.
	needed := glyphH + glyphGap + nameH
	if subLine != "" {
		needed += subGap + subH
	}
	showGlyph := needed <= availH

	glyphColor := white
	if colorPanel {
		glyphColor = classHeroColor(snapshot.Device.DeviceClass)
	}

	var glyphSprite *widgets.Sprite
	if showGlyph {
		glyphSprite = scaledClassGlyph(ClassIcon(snapshot.Device.DeviceClass), scale, glyphColor)
	} else {
		glyphH = 0
	}

	// Vertically center the group (glyph + name + optional sub).
	stackH := nameH
	if showGlyph {
		stackH += glyphH + glyphGap
	}
	if subLine != "" {
		stackH += subGap + subH
	}
	top := oy
	if stackH < availH {
		top = oy + (availH-stackH)/2
	}

	cx := ox + availW/2
	y := top
	var sprites []widgets.Sprite

	if glyphSprite != nil {
		gw := glyphSprite.Image.Bounds().Dx()
		glyphSprite.Position = image.Point{X: cx - gw/2, Y: y}
		sprites = append(sprites, *glyphSprite)
		y += glyphH + glyphGap
	}

	// Name — full white (never dimmed; this is the headline on a small panel).
	if sp := centeredLine(ctx, name, nameTier,
		image.Rect(ox, y, ox+availW, y+nameH), white); sp != nil {
		sprites = append(sprites, *sp)
	}
	y += nameH

	if subLine != "" {
		y += subGap
		subColor := white
		if colorPanel {
			subColor = sharedcolor.Dim(white)
		}
		if sp := centeredLine(ctx, subLine, tiercatalog.TierSmall,
			image.Rect(ox, y, ox+availW, y+subH), subColor); sp != nil {
			sprites = append(sprites, *sp)
		}
	}

	return style.ViewData{Sprites: sprites, Static: true, PaddingPct: 6}
}

// detailRows returns the supporting identity lines beneath the title.
func detailRows(snapshot source.Snapshot) []string {
	var out []string
	out = append(out, fmt.Sprintf("VID:PID %s:%s", snapshot.Device.VendorID, snapshot.Device.ProductID))
	if snapshot.Device.BusNum != "" || snapshot.Device.DevNum != "" {
		out = append(out, fmt.Sprintf("Bus %s Dev %s",
			coreSafeValue(snapshot.Device.BusNum, "?"), coreSafeValue(snapshot.Device.DevNum, "?")))
	}
	if snapshot.Device.Serial != "" {
		out = append(out, "SN "+snapshot.Device.Serial)
	}
	status := "connected"
	if !snapshot.Connected {
		status = "unplugged"
	}
	out = append(out, "Status "+status)
	return out
}

// heroTitle renders the device name at the largest name tier whose whole string
// fits availW, as a sprite positioned at pos. Returns the sprite and the row
// height consumed. Uses Clip overflow (no ellipsis tofu) as a last resort.
func heroTitle(ctx style.StyleContext, name string, availW int, pos image.Point, boundsW int, fg color.RGBA) (*widgets.Sprite, int) {
	tier := tiercatalog.TierSmall
	for _, t := range nameTierLadder {
		ga := ctx.Entry(t).GlyphAdvance
		if ga <= 0 {
			continue
		}
		if len([]rune(name))*ga <= availW {
			tier = t
			break
		}
	}
	face := ctx.Face("", tier)
	if face == nil {
		return nil, 0
	}
	m := face.Metrics()
	rowH := m.RowHeight
	if rowH <= m.GlyphHeight {
		rowH = m.GlyphHeight + 2
	}
	sp := textbox.Render(textbox.Config{
		Bounds:     image.Rect(0, 0, boundsW, rowH),
		Text:       name,
		Font:       face,
		Alignment:  textbox.Center,
		VAlign:     textbox.Top,
		Overflow:   textbox.Clip,
		Foreground: fg,
		Label:      "usb-title",
	})
	if sp != nil {
		sp.Position = pos
	}
	return sp, rowH
}

// centeredLine renders text centered within bounds as a single sprite, using
// the given tier's face and Clip overflow (no ellipsis tofu).
func centeredLine(ctx style.StyleContext, text string, tier tiercatalog.Tier, bounds image.Rectangle, fg color.RGBA) *widgets.Sprite {
	face := ctx.Face("", tier)
	if face == nil {
		return nil
	}
	sp := textbox.Render(textbox.Config{
		Bounds:     image.Rect(0, 0, bounds.Dx(), bounds.Dy()),
		Text:       text,
		Font:       face,
		Alignment:  textbox.Center,
		VAlign:     textbox.Middle,
		Overflow:   textbox.Clip,
		Foreground: fg,
		Label:      "usb-line",
	})
	// textbox renders into a 0-origin image; anchor it at the bounds' panel
	// position so it lands where the caller intended.
	if sp != nil {
		sp.Position = bounds.Min
	}
	return sp
}

// scaledClassGlyph rasterizes the material-icons-24 glyph for iconRune and
// integer-upscales it by `scale` using nearest-neighbor, producing a crisp
// blocky glyph. Returns nil if the icon face is unavailable.
func scaledClassGlyph(iconRune rune, scale int, fg color.RGBA) *widgets.Sprite {
	face, ok := font.Get("material-icons-24")
	if !ok {
		return nil
	}
	m := face.Metrics()
	gw, gh := m.GlyphWidth, m.GlyphHeight
	if gw <= 0 || gh <= 0 {
		return nil
	}
	if scale < 1 {
		scale = 1
	}
	img := image.NewRGBA(image.Rect(0, 0, gw*scale, gh*scale))
	for row := 0; row < gh; row++ {
		bits := face.GlyphRow(iconRune, row)
		for col := 0; col < gw; col++ {
			if bits&(1<<uint(31-col)) == 0 {
				continue
			}
			// Fill the scale×scale block for this source pixel.
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetRGBA(col*scale+dx, row*scale+dy, fg)
				}
			}
		}
	}
	return &widgets.Sprite{Image: img, Label: "usb-class-glyph"}
}

// fillRectSprite returns a sprite of a solid filled rectangle at r (panel
// coordinates), used for the class accent bar.
func fillRectSprite(r image.Rectangle, c color.RGBA, label string) *widgets.Sprite {
	w, h := r.Dx(), r.Dy()
	if w <= 0 || h <= 0 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return &widgets.Sprite{Image: img, Position: r.Min, Label: label}
}

// classHeroColor maps a device class to its hero/accent color. These are
// deliberately saturated so each class reads as a distinct identity at a glance
// on color panels. On mono panels the caller substitutes white.
func classHeroColor(class string) color.RGBA {
	switch class {
	case "Storage":
		return color.RGBA{R: 0, G: 200, B: 0, A: 255} // green — disks
	case "HID":
		return color.RGBA{R: 0, G: 255, B: 255, A: 255} // cyan — input
	case "Audio":
		return color.RGBA{R: 255, G: 128, B: 0, A: 255} // orange — audio
	case "Network":
		return color.RGBA{R: 80, G: 160, B: 255, A: 255} // blue — network
	case "Printer":
		return color.RGBA{R: 200, G: 120, B: 255, A: 255} // violet — printers
	case "Video":
		return color.RGBA{R: 255, G: 80, B: 160, A: 255} // magenta — video
	case "Wireless":
		return color.RGBA{R: 0, G: 220, B: 180, A: 255} // teal — wireless
	case "Hub":
		return color.RGBA{R: 255, G: 191, B: 0, A: 255} // amber — hubs
	default:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // neutral — unknown/Device
	}
}

// isColorPanel reports whether the panel renders in color (vs mono/grayscale),
// so hero fills use class color only where they will not quantize to black.
func isColorPanel(ctx style.StyleContext) bool {
	switch ctx.Cap() {
	case style.ColorSlow, style.ColorFast:
		return true
	default:
		return false
	}
}
