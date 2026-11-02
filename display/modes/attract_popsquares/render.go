package attract_popsquares

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_popsquares/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
	timePhase    float64 // accumulated animation phase (seconds, scaled by Speed)
	isEinkPanel  bool
)

// ─── The idea ─────────────────────────────────────────────────────────────────
//
// popsquares evokes a SUPERCOMPUTER IDLING: a full-bleed grid of flat colour
// cells — the banks of the machine — sitting at a calm base colour, with a
// sparse, unhurried scatter of cells gently pulsing brighter and easing back.
// The observer knows nothing of what the machine is doing; they just see it is
// alive and at rest. The whole craft is RESTRAINT: there must always be
// *something* lit, but never a frenzy. Calm is the product.
//
// Two properties carry the vibe:
//   - Coherence: the base field is a slow, low-amplitude gradient so neighbours
//     are related — one cohesive system, never random confetti.
//   - Non-repetition: each cell's pulse is driven by a per-cell hash and an
//     irrational-period phase, so the twinkle never falls into a visible
//     pattern and never repeats frame to frame.

// buildView produces a complete popsquares frame as style.ViewData.
func buildView(hints textlayout.TextHints) style.ViewData {
	var elapsed time.Duration
	if !lastTick.IsZero() {
		elapsed = time.Since(lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	lastTick = time.Now()

	p := GetPolicy()

	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight
	if panelWidth <= 0 {
		panelWidth = 64
	}
	if panelHeight <= 0 {
		panelHeight = 64
	}
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	isMono, isEink := resolvePanelType(hints)

	// E-ink path: a single frozen calm frame, no animation.
	if isEink {
		isEinkPanel = true
		return buildEinkFrame(p, panelWidth, panelHeight, hints)
	}
	isEinkPanel = false

	timePhase += elapsed.Seconds() * p.Speed

	sprite := renderGrid(p, panelWidth, panelHeight, timePhase, isMono)

	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          false,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(&staticRenderable{s: sprite})

	frameCounter++

	bestStyleName := resolveBestStyleName(hints)
	ctx := style.NewStyleContext(hints)
	snap := source.Snapshot{Sprites: compositor.Sprites(), Policy: p, IsEink: false}
	s := popsquaresRegistry.Lookup(bestStyleName)
	vd := s.Build(snap, p, ctx)
	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}
	return vd
}

// RenderCacheKey returns a change-detection string. For e-ink panels it is
// stable (changes only on policy mutation); for animated panels it advances
// every frame.
func RenderCacheKey() uint32 {
	if isEinkPanel {
		return region.CalcRegionCacheKey("attract_popsquares:eink", GetPolicy().Fingerprint())
	}
	return region.CalcRegionCacheKey("attract_popsquares", frameCounter, GetPolicy().Fingerprint())
}

// cellSize derives the pixel size of a grid cell from GridScale and panel size.
// The base cell is ~1/16 of the short edge so a default panel shows a readable
// bank of squares; GridScale shrinks or grows that. Clamped so cells never
// become sub-pixel or swallow the whole panel.
func cellSize(gridScale float64, panelWidth, panelHeight int) int {
	short := panelWidth
	if panelHeight < short {
		short = panelHeight
	}
	base := float64(short) / 16.0
	cs := int(math.Round(base / gridScale))
	if cs < 4 {
		cs = 4
	}
	if cs > short {
		cs = short
	}
	return cs
}

// renderGrid draws the full-panel grid of cells into an RGBA sprite.
func renderGrid(p source.Policy, width, height int, t float64, mono bool) *widgets.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	cs := cellSize(p.GridScale, width, height)
	cols := (width + cs - 1) / cs
	rows := (height + cs - 1) / cs

	for gy := 0; gy < rows; gy++ {
		for gx := 0; gx < cols; gx++ {
			c := cellColor(gx, gy, cols, rows, t, p, mono)
			fillCell(img, gx*cs, gy*cs, cs, width, height, c)
		}
	}

	return &widgets.Sprite{Image: img, Position: image.Point{}, Label: "popsquares-frame"}
}

// cellColor computes a single cell's colour: a calm base hue from a slow
// drifting gradient, brightened toward white by that cell's current pulse.
func cellColor(gx, gy, cols, rows int, t float64, p source.Policy, mono bool) color.RGBA {
	// Normalised grid position in [0,1].
	ux := 0.0
	uy := 0.0
	if cols > 1 {
		ux = float64(gx) / float64(cols-1)
	}
	if rows > 1 {
		uy = float64(gy) / float64(rows-1)
	}

	// Base field: a gentle 2D gradient that drifts slowly, so neighbours are
	// related and the whole grid reads as one cohesive, breathing system.
	grad := 0.5 + 0.5*math.Sin(ux*1.7+uy*1.3+t*0.15)
	grad2 := 0.5 + 0.5*math.Sin(uy*2.1-ux*0.9-t*0.11)
	field := 0.5*grad + 0.5*grad2 // [0,1]

	// Pulse: each cell has its own phase (from a stable hash) and an
	// irrational-period cycle, so pulses are scattered and never repeat.
	seed := cellHash(gx, gy)
	phase := seed * 2.0 * math.Pi
	// Two incommensurate periods multiplied give a sparse, non-repeating "on"
	// window — the cell is mostly dim, occasionally swelling bright.
	pulseRaw := math.Sin(t*0.9*popIrrA+phase) * math.Sin(t*0.31*popIrrB+phase*1.7)
	// Shape it: only the upper band lights up, so most cells rest dark and a few
	// glow — the sparse idle twinkle. PopRate widens the lit band.
	threshold := 1.0 - 0.9*p.PopRate
	pulse := 0.0
	if pulseRaw > threshold {
		pulse = (pulseRaw - threshold) / (1.0 - threshold) // [0,1] ease-in
		pulse = pulse * pulse                              // soften toward the low end
	}

	if mono {
		// Base cells sit at a low grey; pulses ease toward bright white.
		base := 0.12 + 0.16*field
		lum := base + (1.0-base)*pulse
		v := uint8(clamp01(lum) * 255)
		return color.RGBA{R: v, G: v, B: v, A: 255}
	}

	// Colour: a cool "nominal/idle" palette by default (deep blue → teal →
	// green), nudged warmer by Warmth. Hue varies gently with the field so the
	// bank is cohesive, not a rainbow.
	hueCool := 0.52 + 0.12*field // ~teal/cyan band
	hueWarm := 0.08 + 0.10*field // ~amber/orange band
	hue := hueCool + (hueWarm-hueCool)*p.Warmth

	// Base cells are dim and fairly saturated (deep colour); a pulse lifts value
	// toward full and desaturates slightly toward a hot white core.
	baseVal := 0.22 + 0.20*field
	val := baseVal + (0.95-baseVal)*pulse
	sat := 0.85 - 0.55*pulse

	r, g, b := hsvToRGB(hue, clamp01(sat), clamp01(val))
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// popIrr* are irrational multipliers so the two pulse oscillators never share a
// short common period — the twinkle never falls into a visible pattern.
const (
	popIrrA = 1.0
	popIrrB = 1.6180339887 // golden ratio
)

// cellHash maps a cell's grid coordinates to a stable pseudo-random value in
// [0,1). Deterministic so a cell's phase is consistent across frames, but
// scattered enough that neighbouring cells pulse independently.
func cellHash(gx, gy int) float64 {
	h := uint32(gx)*73856093 ^ uint32(gy)*19349663
	h ^= h >> 13
	h *= 0x85ebca6b
	h ^= h >> 16
	return float64(h&0xffffff) / float64(0x1000000)
}

// fillCell paints a single cell block, clipped to the panel bounds so edge
// cells (which may be partial) stay inside the image.
func fillCell(img *image.RGBA, x0, y0, cs, width, height int, c color.RGBA) {
	x1 := x0 + cs
	y1 := y0 + cs
	if x1 > width {
		x1 = width
	}
	if y1 > height {
		y1 = height
	}
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

// buildEinkFrame renders a single static calm frame for e-ink panels: the base
// field with a few cells frozen mid-pulse, no animation.
func buildEinkFrame(p source.Policy, width, height int, hints textlayout.TextHints) style.ViewData {
	// A fixed instant chosen to show a handful of cells lit.
	const frozenT = 2.3
	sprite := renderGrid(p, width, height, frozenT, true)

	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  width,
		AvailableHeight: height,
	})
	compositor.Add(&staticRenderable{s: sprite})

	styleName := resolveBestStyleName(hints)
	ctx := style.NewStyleContext(hints)
	snap := source.Snapshot{Sprites: compositor.Sprites(), Policy: p, IsEink: true}
	s := popsquaresRegistry.Lookup(styleName)
	vd := s.Build(snap, p, ctx)
	vd.Static = true
	vd.StyleReport = style.StyleReport{Name: styleName, Reason: "fitness"}
	return vd
}

// hsvToRGB converts HSV (all in [0,1]) to RGB bytes.
func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	h = math.Mod(h, 1.0)
	if h < 0 {
		h += 1.0
	}
	h *= 6.0
	i := int(h)
	f := h - float64(i)
	pp := v * (1 - s)
	q := v * (1 - s*f)
	tt := v * (1 - s*(1-f))

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, tt, pp
	case 1:
		r, g, b = q, v, pp
	case 2:
		r, g, b = pp, v, tt
	case 3:
		r, g, b = pp, q, v
	case 4:
		r, g, b = tt, pp, v
	default:
		r, g, b = v, pp, q
	}
	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// staticRenderable wraps a *widgets.Sprite to satisfy widgets.Renderable.
type staticRenderable struct {
	s *widgets.Sprite
}

func (sr *staticRenderable) RenderFrame() *widgets.Sprite { return sr.s }

// resolveBestStyleName returns the Name() of the best-fit style for the hints.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(popsquaresRegistry, hints, "attract_popsquares", "")
	return s.Name()
}

// resolvePanelType infers mono/eink flags from the best-fit style's requirements.
func resolvePanelType(hints textlayout.TextHints) (mono bool, eink bool) {
	s, _ := style.ResolveStyle(popsquaresRegistry, hints, "attract_popsquares", "")
	reqs := s.Requirements()
	mono = reqs.Capability == style.MonoFast
	eink = reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow
	return
}

// tunePolicyForUltraLowRes keeps tiny panels calm: on a 16x8 OLED a grid of
// pops would be noise, so we slow it down and thin the pops out.
func tunePolicyForUltraLowRes(p source.Policy, panelWidth, panelHeight int) source.Policy {
	longEdge, shortEdge := panelWidth, panelHeight
	if panelHeight > panelWidth {
		longEdge, shortEdge = panelHeight, panelWidth
	}
	if longEdge > 16 || shortEdge > 8 {
		return p
	}
	if p.PopRate > 0.3 {
		p.PopRate = 0.3
	}
	if p.Speed > 0.8 {
		p.Speed = 0.8
	}
	if p.GridScale < 1.0 {
		p.GridScale = 1.0
	}
	return p
}
