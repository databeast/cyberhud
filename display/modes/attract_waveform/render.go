package attract_waveform

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
	"github.com/databeast/cyberhud/display/modes/attract_waveform/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// waveformRegistry is the per-mode StyleRegistry for the attract_waveform mode.
var waveformRegistry = style.NewRegistry[source.Data, source.Policy](
	// MonoSlow (e-paper mono)
	styles.MonoSlowStyle{Width: 122, Height: 250},
	styles.MonoSlowStyle{Width: 176, Height: 264},
	styles.MonoSlowStyle{Width: 200, Height: 200},
	styles.MonoSlowStyle{Width: 212, Height: 104},
	styles.MonoSlowStyle{Width: 296, Height: 128},
	styles.MonoSlowStyle{Width: 400, Height: 300},
	styles.MonoSlowStyle{Width: 480, Height: 800},
	styles.MonoSlowStyle{Width: 800, Height: 480},
	styles.MonoSlowStyle{Width: 104, Height: 212},
	styles.MonoSlowStyle{Width: 250, Height: 122},
	styles.MonoSlowStyle{Width: 128, Height: 296},
	styles.MonoSlowStyle{Width: 264, Height: 176},
	styles.MonoSlowStyle{Width: 300, Height: 400},

	// MonoFast (OLED mono)
	styles.MonoFastStyle{Width: 16, Height: 8},
	styles.MonoFastStyle{Width: 8, Height: 16},
	styles.MonoFastStyle{Width: 128, Height: 32},
	styles.MonoFastStyle{Width: 128, Height: 64},
	styles.MonoFastStyle{Width: 128, Height: 128},
	styles.MonoFastStyle{Width: 32, Height: 128},
	styles.MonoFastStyle{Width: 64, Height: 128},

	// GrayscaleSlow (grayscale e-paper)
	styles.GrayscaleSlowStyle{Width: 122, Height: 250},
	styles.GrayscaleSlowStyle{Width: 176, Height: 264},
	styles.GrayscaleSlowStyle{Width: 200, Height: 200},
	styles.GrayscaleSlowStyle{Width: 212, Height: 104},
	styles.GrayscaleSlowStyle{Width: 296, Height: 128},
	styles.GrayscaleSlowStyle{Width: 400, Height: 300},
	styles.GrayscaleSlowStyle{Width: 480, Height: 800},
	styles.GrayscaleSlowStyle{Width: 800, Height: 480},
	styles.GrayscaleSlowStyle{Width: 104, Height: 212},
	styles.GrayscaleSlowStyle{Width: 250, Height: 122},
	styles.GrayscaleSlowStyle{Width: 128, Height: 296},
	styles.GrayscaleSlowStyle{Width: 264, Height: 176},
	styles.GrayscaleSlowStyle{Width: 300, Height: 400},

	// GrayscaleFast (grayscale LED matrix)
	styles.GrayscaleFastStyle{Width: 16, Height: 8},
	styles.GrayscaleFastStyle{Width: 8, Height: 16},
	styles.GrayscaleFastStyle{Width: 160, Height: 80},
	styles.GrayscaleFastStyle{Width: 160, Height: 128},
	styles.GrayscaleFastStyle{Width: 240, Height: 135},
	styles.GrayscaleFastStyle{Width: 240, Height: 240},
	styles.GrayscaleFastStyle{Width: 320, Height: 240},
	styles.GrayscaleFastStyle{Width: 480, Height: 320},
	styles.GrayscaleFastStyle{Width: 800, Height: 480},
	styles.GrayscaleFastStyle{Width: 80, Height: 160},
	styles.GrayscaleFastStyle{Width: 128, Height: 160},
	styles.GrayscaleFastStyle{Width: 135, Height: 240},
	styles.GrayscaleFastStyle{Width: 240, Height: 320},
	styles.GrayscaleFastStyle{Width: 320, Height: 480},
	styles.GrayscaleFastStyle{Width: 480, Height: 800},
	styles.GrayscaleFastStyle{Width: 128, Height: 128},

	// ColorSlow (color e-paper)
	styles.ColorSlowStyle{Width: 122, Height: 250},
	styles.ColorSlowStyle{Width: 176, Height: 264},
	styles.ColorSlowStyle{Width: 200, Height: 200},
	styles.ColorSlowStyle{Width: 212, Height: 104},
	styles.ColorSlowStyle{Width: 296, Height: 128},
	styles.ColorSlowStyle{Width: 400, Height: 300},
	styles.ColorSlowStyle{Width: 480, Height: 800},
	styles.ColorSlowStyle{Width: 800, Height: 480},
	styles.ColorSlowStyle{Width: 104, Height: 212},
	styles.ColorSlowStyle{Width: 250, Height: 122},
	styles.ColorSlowStyle{Width: 128, Height: 296},
	styles.ColorSlowStyle{Width: 264, Height: 176},
	styles.ColorSlowStyle{Width: 300, Height: 400},

	// ColorFast (color TFT)
	styles.ColorFastStyle{Width: 16, Height: 8},
	styles.ColorFastStyle{Width: 8, Height: 16},
	styles.ColorFastStyle{Width: 160, Height: 80},
	styles.ColorFastStyle{Width: 160, Height: 128},
	styles.ColorFastStyle{Width: 240, Height: 135},
	styles.ColorFastStyle{Width: 240, Height: 240},
	styles.ColorFastStyle{Width: 320, Height: 240},
	styles.ColorFastStyle{Width: 480, Height: 320},
	styles.ColorFastStyle{Width: 800, Height: 480},
	styles.ColorFastStyle{Width: 80, Height: 160},
	styles.ColorFastStyle{Width: 128, Height: 160},
	styles.ColorFastStyle{Width: 135, Height: 240},
	styles.ColorFastStyle{Width: 240, Height: 320},
	styles.ColorFastStyle{Width: 320, Height: 480},
	styles.ColorFastStyle{Width: 480, Height: 800},
	styles.ColorFastStyle{Width: 128, Height: 128},
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

// animPhase tracks the cumulative animation phase in seconds, used to morph
// waveform shapes continuously.
var animPhase float64

// BuildView produces a complete waveform frame as style.ViewData.
// It is called by the display runtime each render cycle with the panel's text hints.
func BuildView(hints textlayout.TextHints) style.ViewData {
	// Compute elapsed time since last frame, capped to avoid visual jumps.
	var elapsed time.Duration
	if !lastTick.IsZero() {
		elapsed = time.Since(lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	lastTick = time.Now()

	// Read current policy snapshot.
	p := GetPolicy()

	// Panel dimensions from hints.
	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	// Determine if this is an e-ink panel via best-fit style.
	isEink := resolveIsEink(hints)

	// E-ink path: produce a static frozen frame with no animation.
	if isEink {
		sprites := buildEinkSprites(p, panelWidth, panelHeight)
		return style.ViewData{
			Static:  true,
			Sprites: sprites,
		}
	}

	// Animated path: advance animation phase.
	advance(elapsed, p.Speed)

	// Build sprites for the current frame.
	sprites := buildSprites(p, panelWidth, panelHeight, false)

	frameCounter++

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	s, reason := style.ResolveStyle(waveformRegistry, hints, "attract_waveform", "")
	snap := source.Data{Sprites: sprites, IsEink: false}
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}

	return vd
}

// RenderCacheKey returns a change-detection string incorporating the monotonically
// increasing frame counter and the current policy fingerprint.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(fmt.Sprintf("attract_waveform:%d|%s", frameCounter, GetPolicy().Fingerprint()))
}

// advance updates the animation phase by the elapsed duration scaled by speed.
func advance(elapsed time.Duration, speed float64) {
	animPhase += elapsed.Seconds() * speed
}

// buildSprites renders oscilloscope-style waveform traces into an RGBA sprite.
//
// Each trace is an independent signal with its own frequency, scroll speed, shape,
// and amplitude. All three parameters drift slowly over time via sinusoidal modulation
// so that traces visibly differ from one another and continuously change character.
//
// Line rendering: for each sweep column the waveform is sampled, then every y-value
// the function passes through between the previous column and the current one is filled.
// This is the correct approach for rendering a continuous function on a pixel grid —
// the waveform is present at all intermediate y-values between consecutive samples.
//
// Direction is determined by policy: "horizontal" sweeps left-to-right, "vertical"
// sweeps top-to-bottom, "auto" picks based on aspect ratio (portrait → vertical).
func buildSprites(p source.Policy, panelWidth, panelHeight int, mono bool) []widgets.Sprite {
	if panelWidth <= 0 || panelHeight <= 0 {
		return []widgets.Sprite{{Image: image.NewRGBA(image.Rect(0, 0, 1, 1)), Position: image.Point{}}}
	}

	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	vertical := resolveVertical(p.Direction, panelWidth, panelHeight)

	sweepLen := panelWidth
	ampLen := panelHeight
	if vertical {
		sweepLen = panelHeight
		ampLen = panelWidth
	}

	halfAmp := float64(ampLen) / 2.0
	baseAmplitude := p.Amplitude * halfAmp
	traceCount := p.Traces
	if traceCount < 2 {
		traceCount = 2
	}
	if traceCount > 3 {
		traceCount = 3
	}

	// Map Density [0.1, 1.0] → a tiled cycle count across the sweep axis.
	baseCycles := 3.0 + p.Density*2.5

	for t := 0; t < traceCount; t++ {
		ft := float64(t)

		// Each trace gets a slightly different integer cycle count so the left and
		// right edges line up cleanly instead of tearing at the frame boundary.
		traceCycles := math.Round(baseCycles + ft*0.75)
		if traceCycles < 2 {
			traceCycles = 2
		}
		if traceCycles > 8 {
			traceCycles = 8
		}

		// Keep the waveform exactly tileable: integer cycles ensure both edges line up.
		cycles := traceCycles

		// Each trace scrolls at its own speed. The scroll phase is the product of
		// animPhase and a per-trace multiplier, so traces visibly move at different rates.
		scrollSpeed := p.Speed * (1.0 + ft*0.15)
		tracePhase := animPhase*scrollSpeed*2.0*math.Pi + ft*1.9

		// Shape morphs through sine → sawtooth → harmonic on a per-trace timeline.
		// Rate is slow enough to keep the waveform legible while still moving.
		shapeMix := math.Mod(animPhase*0.12+ft*0.37, 1.0)

		// Amplitude now swings across a much wider envelope so the sweep reads as
		// a living oscilloscope trace instead of a mostly-flat curve.
		ampBreath := 0.7 + 0.5*math.Sin(animPhase*0.13+ft*1.8)
		traceGain := 0.9 + 0.18*ft
		traceAmplitude := baseAmplitude * traceGain * ampBreath
		thickness := 2
		if ampLen < 64 {
			thickness = 1
		}

		var traceColor color.RGBA
		if mono {
			traceColor = color.RGBA{255, 255, 255, 255}
		} else {
			traceColor = waveformTraceColor(cycles, traceAmplitude, tracePhase, baseAmplitude)
		}

		drawOscilloscopeTrace(img, sweepLen, ampLen, vertical, cycles, tracePhase, traceAmplitude, shapeMix, traceColor, thickness)
	}

	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "waveform-traces"},
	}
}

// sweepToXY converts a (sweep-axis, amplitude-axis) coordinate pair to image (x, y).
func sweepToXY(sweep, amp int, vertical bool) (int, int) {
	if vertical {
		return amp, sweep
	}
	return sweep, amp
}

// resolveVertical determines if traces should run vertically based on the
// direction policy and panel aspect ratio.
func resolveVertical(direction string, w, h int) bool {
	switch direction {
	case "vertical":
		return true
	case "horizontal":
		return false
	default: // "auto"
		return h > w // portrait panels get vertical traces
	}
}

// waveformValue computes the waveform amplitude at a given position.
// It morphs between sine, sawtooth, and complex harmonics based on shapeMix.
func waveformValue(x, shapeMix float64) float64 {
	sine := math.Sin(x)
	sawtooth := 2.0*(math.Mod(x/(2.0*math.Pi)+0.5, 1.0)) - 1.0
	harmonic := math.Sin(x) + 0.5*math.Sin(2.0*x) + 0.25*math.Sin(3.0*x)
	harmonic /= 1.75

	// Morph: 0.0-0.33 = sine→sawtooth, 0.33-0.66 = sawtooth→harmonic, 0.66-1.0 = harmonic→sine
	if shapeMix < 0.333 {
		t := shapeMix / 0.333
		return sine*(1.0-t) + sawtooth*t
	} else if shapeMix < 0.666 {
		t := (shapeMix - 0.333) / 0.333
		return sawtooth*(1.0-t) + harmonic*t
	}

	t := (shapeMix - 0.666) / 0.334
	return harmonic*(1.0-t) + sine*t
}

// drawOscilloscopeTrace renders one phase-shifted trace as connected line segments
// so the waveform reads as a continuous oscilloscope sweep.
func drawOscilloscopeTrace(img *image.RGBA, sweepLen, ampLen int, vertical bool, cycles, phase, amplitude, shapeMix float64, traceColor color.RGBA, thickness int) {
	if sweepLen <= 0 || ampLen <= 0 {
		return
	}

	denom := sweepLen - 1
	if denom < 1 {
		denom = 1
	}

	type pt struct{ x, y int }
	points := make([]pt, sweepLen)
	for s := 0; s < sweepLen; s++ {
		sNorm := float64(s) / float64(denom) * cycles * 2.0 * math.Pi
		wv := waveformValue(sNorm+phase, shapeMix) * amplitude
		py := int(math.Round(float64(ampLen-1)/2.0 + wv))
		if py < 0 {
			py = 0
		}
		if py >= ampLen {
			py = ampLen - 1
		}
		px, py := sweepToXY(s, py, vertical)
		points[s] = pt{x: px, y: py}
	}

	for i := 0; i < len(points); i++ {
		curr := points[i]
		drawTraceDot(img, curr.x, curr.y, traceColor, thickness)
		if i > 0 {
			prev := points[i-1]
			drawTraceLine(img, prev.x, prev.y, curr.x, curr.y, traceColor, thickness)
		}
	}
}

func drawTraceDot(img *image.RGBA, x, y int, c color.RGBA, thickness int) {
	if thickness < 1 {
		thickness = 1
	}
	radius := thickness / 2
	bounds := img.Bounds()
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			px := x + dx
			py := y + dy
			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			img.SetRGBA(px, py, c)
		}
	}
}

func drawTraceLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA, thickness int) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		drawTraceDot(img, x0, y0, c, thickness)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// waveformTraceColor maps a trace's live frequency, amplitude, and phase to color.
func waveformTraceColor(cycles, traceAmplitude, tracePhase, baseAmplitude float64) color.RGBA {
	freq := clamp01((cycles - 2.0) / 6.0)
	phase := math.Mod(tracePhase/(2.0*math.Pi), 1.0)
	if phase < 0 {
		phase += 1.0
	}
	amp := 0.0
	if baseAmplitude > 0 {
		amp = clamp01(traceAmplitude / baseAmplitude)
	}

	hue := math.Mod(0.06+0.58*freq+0.27*phase+0.09*amp, 1.0)
	sat := clamp01(0.42 + 0.42*amp + 0.16*math.Sin(tracePhase*0.5+cycles))
	val := clamp01(0.68 + 0.22*amp)

	r, g, b := hsvToRGB(hue, sat, val)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	h = math.Mod(h, 1.0)
	if h < 0 {
		h += 1.0
	}
	h *= 6.0
	i := int(h)
	f := h - float64(i)
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
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

// resolveIsEink determines if the current panel is a slow-refresh (e-ink) panel
// based on the best-fit style from the registry.
func resolveIsEink(hints textlayout.TextHints) bool {
	s, _ := style.ResolveStyle(waveformRegistry, hints, "attract_waveform", "")
	reqs := s.Requirements()
	return reqs.Capability == style.MonoSlow ||
		reqs.Capability == style.GrayscaleSlow ||
		reqs.Capability == style.ColorSlow
}

// resolveBestStyleName returns the Name() of the best-fit style for the given hints.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(waveformRegistry, hints, "attract_waveform", "")
	return s.Name()
}

func tunePolicyForUltraLowRes(p source.Policy, panelWidth, panelHeight int) source.Policy {
	longEdge, shortEdge := panelWidth, panelHeight
	if panelHeight > panelWidth {
		longEdge, shortEdge = panelHeight, panelWidth
	}
	if longEdge > 16 || shortEdge > 8 {
		return p
	}
	if p.Traces > 1 {
		p.Traces = 1
	}
	if p.Persistence > 0.25 {
		p.Persistence = 0.25
	}
	if p.Amplitude > 0.65 {
		p.Amplitude = 0.65
	}
	if p.Speed > 0.85 {
		p.Speed = 0.85
	}
	return p
}
