package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

var (
	// defHostColor accents the hostname line — the single most identifying fact,
	// given prominence through color on capable panels.
	defHostColor = color.RGBA{R: 90, G: 200, B: 255, A: 255}
	// defHealthyColor / defProblemColor color the whole card by health so the
	// verdict reads at a glance (not just via a tiny corner glyph).
	defHealthyColor = color.RGBA{R: 0, G: 200, B: 0, A: 255}
	defProblemColor = color.RGBA{R: 230, G: 70, B: 40, A: 255}
	// defLabelColor dims the "Host:/OS:/IP:" labels so the values stand out.
	defLabelColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	defValueColor = color.RGBA{R: 235, G: 235, B: 235, A: 255}
)

// buildDefaultStyle is the shared "default" variant builder — the identity
// card: hostname (accented), OS/arch, uptime, network addresses, anchored by a
// health verdict the eye lands on first. On color panels the hostname is
// tinted and the whole card leans healthy-green / problem-red; on mono/grayscale
// it falls back to plain bright text plus the status glyph.
func buildDefaultStyle(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	isColor := isColorPanel(ctx)

	allDataOK := snapshot.Hostname != "unknown" && snapshot.Uptime != "n/a" && len(snapshot.IPs) > 0

	// Width budget for truncating values so they never clip off the edge.
	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}
	fit := func(s string) string {
		if maxChars > 0 {
			return textlayout.Truncate(s, maxChars)
		}
		return s
	}

	type row struct {
		text string
		col  color.Color
	}
	healthCol := defHealthyColor
	if !allDataOK {
		healthCol = defProblemColor
	}

	rows := []row{
		{fit(snapshot.Hostname), pick(isColor, defHostColor, nil)},
		{fit("OS: " + snapshot.OSArch), pick(isColor, defValueColor, nil)},
		{fit("Up: " + snapshot.Uptime), pick(isColor, healthCol, nil)},
	}
	if len(snapshot.IPs) == 0 {
		rows = append(rows, row{fit("IP: (none)"), pick(isColor, defProblemColor, nil)})
	} else {
		for _, ip := range snapshot.IPs {
			rows = append(rows, row{fit(ip), pick(isColor, defLabelColor, nil)})
		}
	}

	items := make([]string, len(rows))
	var colors []color.Color
	if isColor {
		colors = make([]color.Color, len(rows))
	}
	for i, r := range rows {
		items[i] = r.text
		if isColor {
			colors[i] = r.col
		}
	}

	// Status glyph, top-right — a secondary cue reinforcing the health verdict
	// (primary cue on color panels is the tinted uptime line; on mono it is the
	// glyph alone).
	var sprites []widgets.Sprite
	if snapshot.GetIcon != nil {
		iconName := "check"
		if !allDataOK {
			iconName = "error"
		}
		if img, ok := snapshot.GetIcon(iconName); ok && img != nil {
			iconSize := 8
			spriteX := bridge.TopRightAnchorX(iconSize)
			_, oy := bridge.ContentOrigin()
			sprites = append(sprites, widgets.Sprite{
				Image:    img,
				Position: image.Point{X: spriteX, Y: oy},
				Label:    "status-" + iconName,
			})
		}
	}

	return style.ViewData{Items: items, Colors: colors, Sprites: sprites}
}

// pick returns a when cond, else b. Used to drop per-row colors on mono panels.
func pick(cond bool, a, b color.Color) color.Color {
	if cond {
		return a
	}
	return b
}
