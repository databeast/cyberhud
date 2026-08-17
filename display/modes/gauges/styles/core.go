package styles

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/databeast/cyberhud/display/modes/gauges/source"
	"github.com/databeast/cyberhud/display/style"
	stylecolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

const defaultFontFamily = "spleen"

type Params struct {
	Shape      progressbar.Style
	LabelTier  tiercatalog.Tier
	PaddingPct int
	TileGapPx  int
	Columns    int
	BuildFn    func(data source.GaugeSet, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.GaugeSet, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.GaugeSet, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	return buildGaugeView(data, pol, ctx, d)
}

func buildGaugeView(data source.GaugeSet, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(maxInt(d.p.PaddingPct, pol.PaddingPct))
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Static: true}
	}

	gauges := data.Gauges
	if len(gauges) == 0 {
		return placeholderView(ctx, bridge, "gauges idle")
	}

	columns := d.p.Columns
	if columns <= 0 {
		columns = autoColumns(len(gauges), bridge.AvailableContentWidth(), bridge.AvailableContentHeight())
	}
	if columns < 1 {
		columns = 1
	}
	if columns > len(gauges) {
		columns = len(gauges)
	}
	rows := int(math.Ceil(float64(len(gauges)) / float64(columns)))
	if rows < 1 {
		rows = 1
	}

	originX, originY := bridge.ContentOrigin()
	gap := maxInt(d.p.TileGapPx, pol.TileGapPx)
	tileW := maxInt(1, (bridge.AvailableContentWidth()-gap*(columns-1))/columns)
	tileH := maxInt(1, (bridge.AvailableContentHeight()-gap*(rows-1))/rows)

	labelTier := resolveLabelTier(pol.LabelTier, d.p.LabelTier)
	labelFace := ctx.Face(defaultFontFamily, labelTier)
	labelColor := stylecolor.ResolveAccent(pol.Accent)

	sprites := make([]widgets.Sprite, 0, len(gauges)*2)
	for i, g := range gauges {
		col := i % columns
		row := i / columns
		x := originX + col*(tileW+gap)
		y := originY + row*(tileH+gap)

		labelH := 0
		if pol.ShowLabels && strings.TrimSpace(g.Label) != "" {
			labelH = labelHeight(labelFace, tileH)
			if labelH > 0 {
				label := textlabel.Render(textlabel.Config{
					Text:       g.Label,
					Bounds:     image.Rect(x, y, x+tileW, y+labelH),
					Font:       labelFace,
					Alignment:  textlabel.Center,
					Foreground: labelColor,
				})
				if label != nil {
					sprites = append(sprites, *label)
				}
			}
		}

		barTop := y + labelH
		barHeight := tileH - labelH
		if labelH > 0 {
			barHeight -= 1
		}
		if barHeight < 1 {
			barTop = y
			barHeight = tileH
		}

		shape := chooseShape(pol.Shape, g.Shape, d.p.Shape)
		cfg := progressbar.Config{
			Style:       shape,
			Orientation: chooseOrientation(shape, tileW, barHeight),
			Value:       g.Percent,
			Bounds:      image.Rect(x, barTop, x+tileW, barTop+barHeight),
			Foreground:  stylecolor.ResolveAccent(firstNonEmpty(g.Accent, pol.Accent)),
			Background:  dimBackground(firstNonEmpty(g.Accent, pol.Accent)),
			RoundedCaps: shape == progressbar.Linear,
		}
		if shape == progressbar.Ring || shape == progressbar.Arc {
			cfg.Thickness = maxInt(2, minInt(tileW, barHeight)/6)
		}
		sprite := progressbar.Render(cfg)
		if sprite != nil {
			sprites = append(sprites, *sprite)
		}
	}

	if len(sprites) == 0 {
		return placeholderView(ctx, bridge, "gauges empty")
	}

	return style.ViewData{
		Sprites:    sprites,
		PaddingPct: maxInt(d.p.PaddingPct, pol.PaddingPct),
		Static:     true,
	}
}

func placeholderView(ctx style.StyleContext, bridge layout.LayoutCalculator, text string) style.ViewData {
	x, y := bridge.ContentOrigin()
	h := bridge.RowHeight()
	if h <= 0 {
		h = 10
	}
	label := textlabel.Render(textlabel.Config{
		Text:       text,
		Bounds:     image.Rect(x, y, x+bridge.AvailableContentWidth(), y+h),
		Font:       ctx.Face(defaultFontFamily, tiercatalog.TierNormal),
		Alignment:  textlabel.Center,
		Foreground: stylecolor.Lookup("white"),
	})
	if label == nil {
		return style.ViewData{Items: []string{text}, Static: true}
	}
	return style.ViewData{Sprites: []widgets.Sprite{*label}, Static: true}
}

func autoColumns(count, width, height int) int {
	if count <= 1 {
		return 1
	}
	if width >= 480 && count >= 6 {
		return 3
	}
	if width >= 320 && count >= 4 {
		return 2
	}
	if height < 160 {
		return 1
	}
	if count > 2 {
		return 2
	}
	return 1
}

func labelHeight(face font.Face, tileH int) int {
	if face == nil {
		return maxInt(10, tileH/4)
	}
	m := face.Metrics()
	h := m.RowHeight
	if h <= 0 {
		h = m.GlyphHeight
	}
	if h <= 0 {
		h = 10
	}
	if h > tileH/2 {
		h = tileH / 2
	}
	return h
}

func resolveLabelTier(policyTier string, fallback tiercatalog.Tier) tiercatalog.Tier {
	switch strings.ToLower(strings.TrimSpace(policyTier)) {
	case "", "auto":
		if fallback != "" {
			return fallback
		}
		return tiercatalog.TierNormal
	default:
		return tiercatalog.Tier(policyTier)
	}
}

func chooseShape(policyShape, gaugeShape string, fallback progressbar.Style) progressbar.Style {
	switch strings.ToLower(strings.TrimSpace(firstNonEmpty(gaugeShape, policyShape))) {
	case "linear", "":
		if fallback == 0 {
			return progressbar.Linear
		}
		return fallback
	case "ring":
		return progressbar.Ring
	case "arc":
		return progressbar.Arc
	case "pie":
		return progressbar.Pie
	default:
		return fallback
	}
}

func chooseOrientation(shape progressbar.Style, w, h int) progressbar.Orientation {
	if shape == progressbar.Ring || shape == progressbar.Arc || shape == progressbar.Pie {
		return progressbar.OrientHorizontal
	}
	if h > w {
		return progressbar.OrientVertical
	}
	return progressbar.OrientHorizontal
}

func dimBackground(accent string) color.RGBA {
	c := stylecolor.ResolveAccent(accent)
	return stylecolor.Dim(c)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
