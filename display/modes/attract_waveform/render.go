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

// animPhase tracks the cumulative animation phase in seconds, used to drive the
// signal generators continuously.
var animPhase float64

// phosphor holds the cross-frame persistence buffer. Each frame fades the
// accumulated glow and composites the new beam on top, so the written frame
// carries the ghost of recent sweeps — the phosphor-persistence trail. The
// buffer is keyed by panel size and trace count so a dimension or policy change
// rebuilds it cleanly rather than smearing stale geometry.
var phosphor phosphorBuffer

type phosphorBuffer struct {
	w, h   int
	traces int
	// intensity per pixel in [0,1], row-major (y*w + x). This is the grayscale
	// "charge" of the phosphor; color is applied per trace at composite time.
	r []float64
	g []float64
	b []float64
}

func (pb *phosphorBuffer) reset(w, h, traces int) {
	pb.w, pb.h, pb.traces = w, h, traces
	n := w * h
	pb.r = make([]float64, n)
	pb.g = make([]float64, n)
	pb.b = make([]float64, n)
}

// ensure (re)allocates the buffer when geometry or trace count changes.
func (pb *phosphorBuffer) ensure(w, h, traces int) {
	if pb.w != w || pb.h != h || pb.traces != traces || pb.r == nil {
		pb.reset(w, h, traces)
	}
}

// fade multiplies every pixel's charge by decay (<1), dimming the trail.
func (pb *phosphorBuffer) fade(decay float64) {
	for i := range pb.r {
		pb.r[i] *= decay
		pb.g[i] *= decay
		pb.b[i] *= decay
	}
}

// deposit adds charge at (x,y) in the given color, saturating at 1.0.
func (pb *phosphorBuffer) deposit(x, y int, cr, cg, cb, intensity float64) {
	if x < 0 || x >= pb.w || y < 0 || y >= pb.h {
		return
	}
	i := y*pb.w + x
	pb.r[i] = math.Min(1.0, pb.r[i]+cr*intensity)
	pb.g[i] = math.Min(1.0, pb.g[i]+cg*intensity)
	pb.b[i] = math.Min(1.0, pb.b[i]+cb*intensity)
}

// compose renders the current phosphor charge into an RGBA image on black.
func (pb *phosphorBuffer) compose() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, pb.w, pb.h))
	for y := 0; y < pb.h; y++ {
		for x := 0; x < pb.w; x++ {
			i := y*pb.w + x
			r := clamp01(pb.r[i])
			g := clamp01(pb.g[i])
			b := clamp01(pb.b[i])
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: 255,
			})
		}
	}
	return img
}

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

// ─── Signal model ───────────────────────────────────────────────────────────
//
// The core aesthetic goal: the trace must read as a *live signal from an unknown
// source* — never visibly repeating, within a frame or across time. A pure sine
// (or any single periodic function) tiled across the panel fails this instantly:
// the eye sees the pattern and the illusion of "something out there" collapses
// into "a computer drew a wave."
//
// So each trace is NOT a periodic function. It is the sum of several sine
// partials whose frequencies are mutually incommensurate (irrational ratios):
// such a sum is quasi-periodic — its true period is astronomically long, so it
// never visibly repeats across the panel. Each partial also drifts in phase and
// amplitude over time at its own slow rate, and a wandering macro-envelope plus
// sparse "events" (transient bursts) give the signal dynamics: surges, lulls,
// the odd jolt. The result behaves like a real instrument watching a real,
// unknown process.

// signalPartials describes one trace's additive-synthesis stack.
type signalPartials struct {
	freqs  []float64 // spatial frequency (cycles across sweep) per partial — irrational ratios
	phases []float64 // per-partial phase drift rate
	gains  []float64 // per-partial amplitude weight
}

// irrationalBase supplies incommensurate frequency ratios so partial sums never
// line up into a short period. These are deliberately non-integer and mutually
// irrational (√2, φ, √5, e/2, π/2 ...).
var irrationalRatios = []float64{
	1.0,
	math.Sqrt2,       // 1.4142…
	1.6180339887,     // golden ratio φ
	math.Sqrt(5),     // 2.2360…
	math.E / 2.0,     // 1.3591…
	math.Pi/2 + 0.37, // 1.9408…
}

// buildPartials constructs a trace's partial stack. baseFreq sets how many broad
// undulations span the sweep (kept near 1 so one macro-shape fills the panel,
// never a row of identical humps); the higher partials ride on top at
// irrational multiples to add non-repeating fine structure.
func buildPartials(trace int, baseFreq float64) signalPartials {
	ft := float64(trace)
	// Each trace is seeded differently so traces never look alike.
	seed := 0.73*ft + 0.19
	n := len(irrationalRatios)
	sp := signalPartials{
		freqs:  make([]float64, n),
		phases: make([]float64, n),
		gains:  make([]float64, n),
	}
	for i := 0; i < n; i++ {
		ratio := irrationalRatios[i]
		sp.freqs[i] = baseFreq * ratio
		// Each partial drifts at its own slow, irrational rate so the composite
		// shape is always reorganising — the signal is never the same twice.
		sp.phases[i] = 0.11 + 0.043*float64(i) + 0.017*math.Sin(seed+ft)
		// Higher partials carry progressively less energy (1/f-ish) for a
		// natural, non-buzzy shape, with a per-trace tilt.
		sp.gains[i] = (1.0 / (1.0 + 1.6*float64(i))) * (0.85 + 0.3*math.Sin(seed*1.7+float64(i)))
	}
	return sp
}

// sample evaluates the trace's normalised value in [-1,1] at sweep position
// u ∈ [0,1] and time t. The sum of irrational-ratio partials is quasi-periodic,
// so scanning u across the panel never produces a repeating pattern, and
// advancing t never reproduces a past frame.
func (sp signalPartials) sample(u, t float64) float64 {
	var sum, norm float64
	for i := range sp.freqs {
		phase := t * sp.phases[i] * 2.0 * math.Pi
		sum += sp.gains[i] * math.Sin((u*sp.freqs[i]*2.0*math.Pi)+phase)
		norm += sp.gains[i]
	}
	if norm == 0 {
		return 0
	}
	return sum / norm
}

// macroEnvelope returns a slow-wandering amplitude multiplier in roughly
// [0.35, 1.0] plus sparse transient "events" (bursts) so the signal surges and
// calms instead of droning at one level. Driven by incommensurate low
// frequencies so the dynamics themselves never repeat.
func macroEnvelope(u, t, traceSeed float64) float64 {
	// Slow breathing swell across the whole trace.
	breath := 0.68 + 0.32*math.Sin(t*0.37+traceSeed*2.1)
	// A travelling swell that moves along the sweep over time, so different
	// regions of the trace are "hot" at different moments.
	travel := 0.5 + 0.5*math.Sin((u*2.3-t*0.21+traceSeed)*2.0*math.Pi*0.18)
	env := breath * (0.6 + 0.4*travel)

	// Sparse transient bursts: a narrow Gaussian bump that sweeps a random-ish
	// position, giving the occasional sharp jolt an unknown signal would show.
	burstPos := math.Mod(t*0.13+traceSeed*0.61, 1.0)
	d := u - burstPos
	burst := 0.9 * math.Exp(-(d*d)/(2*0.015*0.015))
	// Gate bursts so they come and go rather than being always present.
	gate := math.Max(0, math.Sin(t*0.29+traceSeed*3.3))
	env += burst * gate

	return clampf(env, 0.25, 1.6)
}

// buildSprites renders the living-signal traces into the phosphor buffer and
// composes the result. Each frame fades the accumulated charge (the persistence
// trail) and deposits the new sweep on top, the sweep head brightest.
func buildSprites(p source.Policy, panelWidth, panelHeight int, mono bool) []widgets.Sprite {
	if panelWidth <= 0 || panelHeight <= 0 {
		return []widgets.Sprite{{Image: image.NewRGBA(image.Rect(0, 0, 1, 1)), Position: image.Point{}}}
	}

	vertical := resolveVertical(p.Direction, panelWidth, panelHeight)

	sweepLen := panelWidth
	ampLen := panelHeight
	if vertical {
		sweepLen = panelHeight
		ampLen = panelWidth
	}

	traceCount := p.Traces
	if traceCount < 2 {
		traceCount = 2
	}
	if traceCount > 4 {
		traceCount = 4
	}

	phosphor.ensure(panelWidth, panelHeight, traceCount)

	// Persistence → per-frame decay. Higher persistence = slower fade = longer
	// trail. Map [0.1,1.0] → decay [0.55, 0.92]. Below ~0.6 the trail is a short
	// comet; near 0.92 it lingers as a long ghost.
	decay := 0.55 + 0.37*clamp01((p.Persistence-0.1)/0.9)
	phosphor.fade(decay)

	// baseFreq near 1: one broad macro-undulation spans the sweep so the trace
	// never shows a row of identical humps. density nudges the fine-structure
	// richness slightly without ever tiling.
	baseFreq := 0.8 + 0.6*clamp01(p.Density)

	halfAmp := float64(ampLen) / 2.0
	baseAmplitude := clamp01(p.Amplitude) * halfAmp * 0.92

	thickness := 2
	if ampLen < 96 {
		thickness = 1
	}

	denom := sweepLen - 1
	if denom < 1 {
		denom = 1
	}

	for t := 0; t < traceCount; t++ {
		ft := float64(t)
		seed := 0.73*ft + 0.19
		sp := buildPartials(t, baseFreq)

		var cr, cg, cb float64
		if mono {
			cr, cg, cb = 1, 1, 1
		} else {
			cr, cg, cb = traceColor(t, animPhase)
		}

		// Sweep head position: the "live" end of the beam, travelling along the
		// sweep axis over time. Pixels near the head are drawn brightest; the
		// phosphor fade then turns the rest into the fading tail.
		headU := math.Mod(animPhase*0.45+ft*0.27, 1.0)

		prevX, prevY := -1, -1
		for s := 0; s <= denom; s++ {
			u := float64(s) / float64(denom)
			env := macroEnvelope(u, animPhase, seed)
			v := sp.sample(u, animPhase) * env
			amp := v * baseAmplitude
			ay := int(math.Round(float64(ampLen-1)/2.0 + amp))
			if ay < 0 {
				ay = 0
			}
			if ay >= ampLen {
				ay = ampLen - 1
			}
			x, y := sweepToXY(s, ay, vertical)

			// Head brightness: brightest at the live head, tapering with sweep
			// distance so the beam reads as "moving right now". Base level is
			// kept moderate so overlapping passes and cross-frame accumulation
			// build graduated mid-tones (the phosphor glow) rather than slamming
			// every pixel to blown-out white.
			dist := math.Abs(u - headU)
			if dist > 0.5 {
				dist = 1.0 - dist // wrap distance around the sweep
			}
			headGlow := 0.32 + 0.5*math.Exp(-(dist*dist)/(2*0.08*0.08))

			depositTrace(x, y, cr, cg, cb, headGlow, thickness)
			if prevX >= 0 {
				depositLine(prevX, prevY, x, y, cr, cg, cb, headGlow, thickness)
			}
			prevX, prevY = x, y
		}
	}

	img := phosphor.compose()
	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "waveform-traces"},
	}
}

// depositTrace lays down a glowing dab of phosphor charge: a bright core
// surrounded by a wide, soft bloom halo so the beam reads as a lit CRT trace
// rather than a hard vector line. The halo extends well beyond the core and
// falls off gently, so a rendered frame carries noticeably more soft glow than
// hard core — the visual signature of phosphor.
func depositTrace(x, y int, cr, cg, cb, intensity float64, thickness int) {
	if thickness < 1 {
		thickness = 1
	}
	core := thickness / 2
	// The bloom halo reaches two pixels past the core, giving a soft glowing
	// edge. sigma controls the falloff; a larger sigma spreads the glow wider.
	halo := core + 3
	sigma := float64(core) + 1.6
	twoSigmaSq := 2.0 * sigma * sigma
	for dy := -halo; dy <= halo; dy++ {
		for dx := -halo; dx <= halo; dx++ {
			px, py := x+dx, y+dy
			dd := float64(dx*dx + dy*dy)
			var falloff float64
			if dd <= float64(core*core) {
				// Core: full strength.
				falloff = 1.0
			} else {
				// Halo: gentle Gaussian bloom. Scaled down so the halo is
				// clearly dimmer than the core but still deposits real glow.
				falloff = 0.5 * math.Exp(-dd/twoSigmaSq)
			}
			if falloff < 0.02 {
				continue
			}
			phosphor.deposit(px, py, cr, cg, cb, intensity*falloff)
		}
	}
}

// depositLine fills phosphor charge along the segment between two samples so the
// trace is continuous (Bresenham).
func depositLine(x0, y0, x1, y1 int, cr, cg, cb, intensity float64, thickness int) {
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
		depositTrace(x0, y0, cr, cg, cb, intensity, thickness)
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

// traceColor assigns each trace a punchy phosphor hue — the classic scope
// palette (green / cyan / amber / magenta) that reads as a glowing CRT rather
// than the dull olive a narrow warm band produced before. Hue breathes slowly
// over time so the colour is alive without ever becoming muddy.
func traceColor(trace int, t float64) (float64, float64, float64) {
	// Anchor hues spaced around the phosphor-friendly part of the wheel.
	anchors := []float64{
		0.33, // green
		0.50, // cyan
		0.12, // amber
		0.85, // magenta
	}
	h := anchors[trace%len(anchors)]
	// Gentle hue drift keeps it living but stays near the anchor so it never
	// drifts into muddy territory.
	h = math.Mod(h+0.015*math.Sin(t*0.21+float64(trace)), 1.0)
	if h < 0 {
		h += 1
	}
	// High saturation and value for that bright-phosphor pop.
	s := 0.9
	v := 1.0
	r, g, b := hsvToRGBf(h, s, v)
	return r, g, b
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

// hsvToRGBf converts HSV in [0,1] to RGB in [0,1].
func hsvToRGBf(h, s, v float64) (float64, float64, float64) {
	h = math.Mod(h, 1.0)
	if h < 0 {
		h += 1.0
	}
	h *= 6.0
	i := int(h)
	f := h - float64(i)
	p := v * (1 - s)
	q := v * (1 - s*f)
	tt := v * (1 - s*(1-f))

	switch i % 6 {
	case 0:
		return v, tt, p
	case 1:
		return q, v, p
	case 2:
		return p, v, tt
	case 3:
		return p, q, v
	case 4:
		return tt, p, v
	default:
		return v, p, q
	}
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

func clampf(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
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
