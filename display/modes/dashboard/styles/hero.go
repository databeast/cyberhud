package styles

// Large-color "hero" dashboard layout.
//
// This is the flagship layout for big color TFT panels (800x480, 480x800). The
// skeleton/compact builders stack bare centered text rows, which on a large
// color panel reads as a postage stamp marooned in black and leaves the panel's
// color and space unused. This layout earns the hardware instead: a bordered
// frame that makes the space feel authored, a large accent hostname title, and
// an icon-labeled field list (computer / uptime / LAN / wifi / version) with a
// pulsing accent "heartbeat" dot that delivers the uptime-heartbeat the mode's
// description promises.
//
// Everything composites as absolute-positioned sprites — a border, rasterized
// Material-icon glyphs, and text rasterized to labels — so the whole frame lives
// on one plane (no row-text / sprite two-plane drift), exactly as usb and wifi do.

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/style"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	"github.com/databeast/cyberhud/runtime/host"
)

// dashboardFontFamily is the font family the hero layout rasterizes its text
// labels with. Spleen is always registered (it is the system default face), so
// ctx.Face resolves a real face at every tier.
const dashboardFontFamily = "spleen"

// heroField is one icon-labeled status row: a Material icon rune, the value
// text, and whether the value renders in the accent (true) or a quiet tone.
type heroField struct {
	icon   rune
	value  string
	accent bool
}

// buildColorHero is the large-color flagship layout, shared by the landscape
// (800x480) and portrait (480x800) panels. It reads the already-settled
// snapshot and composites sprites; it holds no per-frame state beyond reading
// the frame clock for the heartbeat phase, which is a pure function of the
// current time.
func buildColorHero(data source.DashboardContent, pol source.Policy, ctx style.StyleContext, _ styleDef) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(4)
	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	if cw <= 0 || chh <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: false}
	}

	accent := source.ResolveUptimeAccent(pol.ColorAccent)
	white := color.RGBA{255, 255, 255, 255}
	dimmed := color.RGBA{170, 170, 170, 255}
	ox, oy := bridge.ContentOrigin()

	var sprites []widgets.Sprite

	// --- Authored space: a border frame hugging the panel edge. ---
	if frame := borderframe.Render(borderframe.Config{
		Bounds:    image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:     "sharp",
		ColorTint: accent,
	}); frame != nil {
		sprites = append(sprites, *frame)
	}

	// --- Title band: hostname as a large accent heading, with a computer icon
	// and the live heartbeat dot to its right. ---
	titleFace := ctx.Face(dashboardFontFamily, tiercatalog.TierLarge)
	titleH := tierHeight(ctx, tiercatalog.TierLarge, 24)
	iconSize := titleH

	x := ox
	sprites, x = placeIconInBand(sprites, font.IconComputer, iconSize, accent, x, oy, titleH)
	if lbl := textlabel.Render(textlabel.Config{
		Text:       data.Hostname,
		Bounds:     image.Rect(x, oy, ox+cw, oy+titleH),
		Font:       titleFace,
		Foreground: accent,
	}); lbl != nil {
		sprites = append(sprites, *lbl)
	}
	// Heartbeat: a filled accent dot at the top-right that brightens on a ~1Hz
	// pulse, the literal "uptime heartbeat indicator" the description promises.
	if dot := heartbeatDot(ox+cw-titleH, oy, titleH*2/3, accent); dot != nil {
		sprites = append(sprites, *dot)
	}

	// Two columns below the title: identity fields on the left, the health
	// vitals roll-up on the right. The split fills the panel's width instead of
	// stranding everything in the left gutter. A thin accent rule under the
	// title frames the band and divides it from the header.
	bodyH := tierHeight(ctx, tiercatalog.TierNormal, 14)
	bandTop := oy + titleH + bodyH/2
	bandBottom := oy + chh
	leftW := cw * 55 / 100
	gutter := ox + leftW
	sprites = append(sprites, *ruleSprite(ox, bandTop-bodyH/3, cw, dividerColor(accent), "dash-rule"))

	// --- Left column: identity + liveness fields, icon-labeled. ---
	fields := []heroField{
		{font.IconAvTimer, data.Uptime, true},
		{font.IconLan, data.IPAddress, false},
		{font.IconWifi, data.WifiSSID, false},
		{font.IconInfo, data.Version, false},
	}
	bodyFace := ctx.Face(dashboardFontFamily, tiercatalog.TierNormal)
	fieldIconSize := bodyH

	step := (bandBottom - bandTop) / len(fields)
	if step < bodyH {
		step = bodyH
	}
	for i, f := range fields {
		rowY := bandTop + i*step + (step-bodyH)/2
		tint := accent
		if !f.accent {
			tint = dimmed
		}
		fx := ox
		sprites, fx = placeIconInBand(sprites, f.icon, fieldIconSize, tint, fx, rowY, bodyH)
		valColor := white
		if f.accent {
			valColor = accent
		}
		if lbl := textlabel.Render(textlabel.Config{
			Text:       f.value,
			Bounds:     image.Rect(fx, rowY, gutter-bodyH/2, rowY+bodyH),
			Font:       bodyFace,
			Foreground: valColor,
		}); lbl != nil {
			sprites = append(sprites, *lbl)
		}
	}

	// --- Right column: the health roll-up. ---
	sprites = append(sprites, placeHealth(ctx, data.Health, gutter, bandTop, bandBottom, (ox+cw)-gutter, accent)...)

	return style.ViewData{
		Sprites: sprites,
		Static:  false, // heartbeat animates
		Cursor:  -1,
	}
}

// placeHealth renders the health vitals in the right column: a STATUS headline
// colored by the overall verdict, and underneath, one named callout line per
// contributor dragging health down. When everything is OK the column shows just
// a calm green STATUS: OK with nothing to call out — the healthy case stays
// quiet, the degraded case points at the culprit.
func placeHealth(ctx style.StyleContext, h host.Health, x, top, bottom, w int, accent color.RGBA) []widgets.Sprite {
	var out []widgets.Sprite
	headFace := ctx.Face(dashboardFontFamily, tiercatalog.TierNormal)
	bodyFace := ctx.Face(dashboardFontFamily, tiercatalog.TierSmall)
	headH := tierHeight(ctx, tiercatalog.TierNormal, 14)
	lineH := tierHeight(ctx, tiercatalog.TierSmall, 10)

	y := top

	// Headline: STATUS: <verdict>, verdict tinted by severity.
	statusColor := severityColor(h.Overall, accent)
	iconX := x
	out, iconX = placeIconInBand(out, severityIcon(h.Overall), headH, statusColor, iconX, y, headH)
	if lbl := textlabel.Render(textlabel.Config{
		Text:       "STATUS: " + h.StatusWord(),
		Bounds:     image.Rect(iconX, y, x+w, y+headH),
		Font:       headFace,
		Foreground: statusColor,
	}); lbl != nil {
		out = append(out, *lbl)
	}
	y += headH + lineH/2

	// Culprit callouts: only the warn/crit contributors, each named with its
	// detail and tinted by its own severity. Healthy contributors stay silent.
	for _, c := range h.Degraded {
		if y+lineH > bottom {
			break
		}
		text := c.Name + " " + c.Detail
		if lbl := textlabel.Render(textlabel.Config{
			Text:       text,
			Bounds:     image.Rect(x, y, x+w, y+lineH),
			Font:       bodyFace,
			Foreground: severityColor(c.Severity, accent),
		}); lbl != nil {
			out = append(out, *lbl)
		}
		y += lineH + lineH/4
	}
	return out
}

// placeIconInBand renders an icon glyph and places it vertically centered within
// a text band of height bandH starting at (x, bandY), so the icon's optical
// center aligns with the row's text center instead of both pinning to the band's
// top-left corner. Returns the updated sprite slice and the x-cursor positioned
// just after the icon plus a gap, ready for the label. When the glyph can't be
// built, x advances by nothing and the label simply starts at x.
func placeIconInBand(sprites []widgets.Sprite, r rune, size int, tint color.RGBA, x, bandY, bandH int) ([]widgets.Sprite, int) {
	g := iconGlyph(r, size, tint)
	if g == nil {
		return sprites, x
	}
	gh := g.Image.Bounds().Dy()
	gw := g.Image.Bounds().Dx()
	g.Position = image.Point{X: x, Y: bandY + (bandH-gh)/2}
	sprites = append(sprites, *g)
	return sprites, x + gw + gw/3
}

// iconGlyph rasterizes a Material-icons-24 glyph at the given pixel size, tinted
// fg. It integer-scales the 24px face to the nearest size that fits, mirroring
// the usb mode's class-glyph technique. Returns nil if the icon face is
// unavailable or size is non-positive.
func iconGlyph(r rune, size int, fg color.RGBA) *widgets.Sprite {
	if size <= 0 {
		return nil
	}
	face, ok := font.Get("material-icons-24")
	if !ok {
		return nil
	}
	m := face.Metrics()
	gw, gh := m.GlyphWidth, m.GlyphHeight
	if gw <= 0 || gh <= 0 {
		return nil
	}
	scale := size / gh
	if scale < 1 {
		scale = 1
	}
	img := image.NewRGBA(image.Rect(0, 0, gw*scale, gh*scale))
	for row := 0; row < gh; row++ {
		bits := face.GlyphRow(r, row)
		for col := 0; col < gw; col++ {
			if bits&(1<<uint(31-col)) == 0 {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetRGBA(col*scale+dx, row*scale+dy, fg)
				}
			}
		}
	}
	return &widgets.Sprite{Image: img, Label: "dashboard-icon"}
}

// heartbeatDot renders a filled circle whose brightness pulses on a ~1Hz cycle
// derived from the frame clock, giving a live "still alive" indicator. The pulse
// is a pure function of the current time, so the style stays stateless.
func heartbeatDot(x, y, size int, accent color.RGBA) *widgets.Sprite {
	if size <= 0 {
		return nil
	}
	// Phase in [0,1) over a 1s period; brightness eases between a dim floor and
	// full via a triangle wave so it reads as a steady pulse, not a blink.
	ns := frameclock.Now().UnixNano() % int64(1e9)
	phase := float64(ns) / 1e9
	tri := phase * 2
	if tri > 1 {
		tri = 2 - tri
	}
	const floor = 0.35
	b := floor + (1.0-floor)*tri
	c := color.RGBA{
		R: uint8(float64(accent.R) * b),
		G: uint8(float64(accent.G) * b),
		B: uint8(float64(accent.B) * b),
		A: 255,
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	cx, cy := float64(size-1)/2, float64(size-1)/2
	rad := float64(size) / 2
	for py := 0; py < size; py++ {
		for px := 0; px < size; px++ {
			dx, dy := float64(px)-cx, float64(py)-cy
			if dx*dx+dy*dy <= rad*rad {
				img.SetRGBA(px, py, c)
			}
		}
	}
	return &widgets.Sprite{Image: img, Position: image.Point{X: x, Y: y}, Label: "dashboard-heartbeat"}
}

// tierHeight returns the row height for a tier, falling back to def when the
// catalog has no entry (degraded/no-font regions).
func tierHeight(ctx style.StyleContext, tier tiercatalog.Tier, def int) int {
	if h := ctx.Entry(tier).RowHeight; h > 0 {
		return h
	}
	return def
}

// ruleSprite renders a 1px horizontal accent rule spanning w pixels at (x, y) —
// a touch of structure that frames the content band so the panel reads as
// authored rather than a text block floating in black.
func ruleSprite(x, y, w int, c color.RGBA, label string) *widgets.Sprite {
	if w < 1 {
		w = 1
	}
	img := image.NewRGBA(image.Rect(0, 0, w, 1))
	for px := 0; px < w; px++ {
		img.SetRGBA(px, 0, c)
	}
	return &widgets.Sprite{Image: img, Position: image.Point{X: x, Y: y}, Label: label}
}

// dividerColor returns a dimmed version of the accent for the structural rule,
// so it frames without competing with the data for attention.
func dividerColor(accent color.RGBA) color.RGBA {
	return color.RGBA{R: accent.R / 2, G: accent.G / 2, B: accent.B / 2, A: 255}
}

// severityColor maps a health severity to its display color: accent for OK (so
// a healthy box reads in the panel's own accent), amber for warn, red for crit.
func severityColor(s host.Severity, accent color.RGBA) color.RGBA {
	switch s {
	case host.SeverityWarn:
		return color.RGBA{R: 255, G: 200, B: 64, A: 255}
	case host.SeverityCritical:
		return color.RGBA{R: 255, G: 96, B: 96, A: 255}
	default:
		return accent
	}
}

// severityIcon maps a health severity to a Material glyph: check for OK, error
// for warn/crit, so the status line carries its meaning at a glance.
func severityIcon(s host.Severity) rune {
	if s > host.SeverityOK {
		return font.IconError
	}
	return font.IconCheckCircle
}
