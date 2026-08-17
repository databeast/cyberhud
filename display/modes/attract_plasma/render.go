package attract_plasma

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_plasma/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
	timePhase    float64 // accumulated time phase for plasma animation
	isEinkPanel  bool    // tracks whether the current panel is e-ink
)

// buildView produces a complete plasma frame as style.ViewData.
func buildView(hints textlayout.TextHints) style.ViewData {
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

	// Graceful degradation: if hints are missing, use minimal defaults.
	if panelWidth <= 0 {
		panelWidth = 64
	}
	if panelHeight <= 0 {
		panelHeight = 64
	}
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	// Resolve panel type from the best-fitting registered style.
	isMono, isEink := resolvePanelType(hints)

	// E-ink path: produce a static frozen frame with no animation.
	// Frame counter is NOT incremented so RenderCacheKey remains stable.
	if isEink {
		isEinkPanel = true
		return buildEinkFrame(p, panelWidth, panelHeight, hints)
	}

	isEinkPanel = false

	// Advance time phase by elapsed scaled by Speed.
	elapsedSec := elapsed.Seconds()
	timePhase += elapsedSec * p.Speed

	// Render the plasma sprite.
	sprite := renderPlasmaSprite(p, panelWidth, panelHeight, timePhase, isMono)

	// Build compositor for this frame.
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          false,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(&staticRenderable{s: sprite})

	frameCounter++

	// Resolve best style name for StyleReport.
	bestStyleName := resolveBestStyleName(hints)

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	snap := source.Snapshot{Sprites: compositor.Sprites(), Policy: p, IsEink: false}
	s := plasmaRegistry.Lookup(bestStyleName)
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}

	return vd
}

// RenderCacheKey returns a change-detection string incorporating the monotonically
// increasing frame counter and the current policy fingerprint. For e-ink panels,
// returns a stable key that changes only on policy mutation.
func RenderCacheKey() uint32 {
	if isEinkPanel {
		return region.CalcRegionCacheKey("attract_plasma:eink", GetPolicy().Fingerprint())
	}
	return region.CalcRegionCacheKey("attract_plasma", frameCounter, GetPolicy().Fingerprint())
}

// renderPlasmaSprite generates the full-panel plasma pattern as an image.RGBA.
// The pattern is computed as the sum of 3+ overlapping sine functions with
// varying frequencies and phase offsets. Color mapping cycles through a 256-step
// gradient palette at CycleRate. BlobScale controls spatial frequency.
func renderPlasmaSprite(p source.Policy, width, height int, t float64, mono bool) *widgets.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Precompute spatial frequency multiplier (higher BlobScale = finer detail).
	sf := p.BlobScale * 0.04 // base spatial frequency

	// Color palette offset cycles at CycleRate.
	paletteOffset := t * p.CycleRate
	longEdge, shortEdge := width, height
	if height > width {
		longEdge, shortEdge = height, width
	}
	ultraLowRes := longEdge <= 16 && shortEdge <= 8

	for y := 0; y < height; y++ {
		fy := float64(y)
		for x := 0; x < width; x++ {
			fx := float64(x)

			// Sum of overlapping sine functions with varying frequencies and phases.
			// Function 1: diagonal wave
			v1 := math.Sin(fx*sf + t*0.7)
			// Function 2: circular wave from center
			cx := fx - float64(width)/2
			cy := fy - float64(height)/2
			v2 := math.Sin(math.Sqrt(cx*cx+cy*cy)*sf*1.5 - t*0.5)
			// Function 3: horizontal wave with vertical modulation
			v3 := math.Sin(fy*sf*0.8 + math.Sin(fx*sf*0.5+t*0.3)*2.0)
			// Function 4: radial wave with phase offset
			v4 := math.Sin((fx*sf*0.6+fy*sf*0.4)*1.2 + t*0.9)

			// Combine and normalize to [0, 1].
			combined := (v1 + v2 + v3 + v4) / 4.0 // range [-1, 1]
			normalized := (combined + 1.0) / 2.0  // range [0, 1]

			// Map to palette index with cycling offset.
			paletteIdx := math.Mod(normalized+paletteOffset, 1.0)
			if paletteIdx < 0 {
				paletteIdx += 1.0
			}

			if ultraLowRes {
				isActive := normalized > 0.72 && ((x+y+int(t*8))%2 == 0)
				if !isActive {
					img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
					continue
				}
				if mono {
					lum := uint8(150 + paletteIdx*105)
					img.SetRGBA(x, y, color.RGBA{R: lum, G: lum, B: lum, A: 255})
				} else {
					r, g, b := hsvToRGB(paletteIdx, 0.7, 0.85)
					img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
				}
				continue
			}

			if mono {
				// Monochrome: luminance-only mapping.
				lum := uint8(paletteIdx * 255)
				img.SetRGBA(x, y, color.RGBA{R: lum, G: lum, B: lum, A: 255})
			} else {
				// Color: map through HSV-like gradient palette.
				r, g, b := hsvToRGB(paletteIdx, 0.85, 0.9)
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}

	return &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Label:    "plasma-frame",
	}
}

// buildEinkFrame renders a single static plasma frame at time=0 for e-ink panels.
func buildEinkFrame(p source.Policy, panelWidth, panelHeight int, hints textlayout.TextHints) style.ViewData {
	sprite := renderPlasmaSprite(p, panelWidth, panelHeight, 0, true)
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(&staticRenderable{s: sprite})

	styleName := resolveBestStyleName(hints)
	ctx := style.NewStyleContext(hints)
	snap := source.Snapshot{Sprites: compositor.Sprites(), Policy: p, IsEink: true}
	s := plasmaRegistry.Lookup(styleName)
	vd := s.Build(snap, p, ctx)
	vd.Static = true
	vd.StyleReport = style.StyleReport{Name: styleName, Reason: "fitness"}
	return vd
}

// hsvToRGB converts HSV (h in [0,1], s in [0,1], v in [0,1]) to RGB bytes.
func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	h = h * 6.0
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
	case 5:
		r, g, b = v, p, q
	}
	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

// staticRenderable wraps a *widgets.Sprite to satisfy the widgets.Renderable interface.
type staticRenderable struct {
	s *widgets.Sprite
}

func (sr *staticRenderable) RenderFrame() *widgets.Sprite { return sr.s }

// resolveBestStyleName returns the Name() of the best-fit style for the given hints.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(plasmaRegistry, hints, "attract_plasma", "")
	return s.Name()
}

// resolvePanelType infers mono/eink flags from the best-fit style's SurfaceRequirements.
func resolvePanelType(hints textlayout.TextHints) (mono bool, eink bool) {
	s, _ := style.ResolveStyle(plasmaRegistry, hints, "attract_plasma", "")
	reqs := s.Requirements()
	mono = reqs.Capability == style.MonoFast
	eink = reqs.Capability == style.MonoSlow || reqs.Capability == style.GrayscaleSlow || reqs.Capability == style.ColorSlow
	return
}

func tunePolicyForUltraLowRes(p source.Policy, panelWidth, panelHeight int) source.Policy {
	longEdge, shortEdge := panelWidth, panelHeight
	if panelHeight > panelWidth {
		longEdge, shortEdge = panelHeight, panelWidth
	}
	if longEdge > 16 || shortEdge > 8 {
		return p
	}
	if p.BlobScale > 0.8 {
		p.BlobScale = 0.8
	}
	if p.CycleRate > 0.9 {
		p.CycleRate = 0.9
	}
	if p.Speed > 0.8 {
		p.Speed = 0.8
	}
	return p
}
