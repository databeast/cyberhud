package attract_shapes

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_shapes/source"
	"github.com/databeast/cyberhud/display/modes/attract_shapes/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

// maxTickElapsed caps the per-frame advancement to prevent visual jumps when
// frames are occasionally delayed (e.g., GC pause or I/O stall).
const maxTickElapsed = 80 * time.Millisecond

// BuildView produces a complete shapes frame as style.ViewData.
// It is called by the display runtime each render cycle with the panel's text hints.
func BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}

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
		sprites := styles.BuildEinkSprites(p, panelWidth, panelHeight)
		return style.ViewData{
			Static:  true,
			Sprites: sprites,
		}
	}

	// Animated path: advance animation phase.
	elapsedSec := elapsed.Seconds() * p.Speed

	// Ensure shapes are initialized for current dimensions and count.
	source.InitShapes(p, panelWidth, panelHeight)

	// Advance rotation and scale oscillation for each shape.
	source.AdvanceShapes(elapsedSec, p)

	// Build sprites for the current frame.
	sprites := source.BuildSprites(p, panelWidth, panelHeight, false)

	frameCounter++

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	bestStyleName := resolveBestStyleName(hints)
	snap := source.Snapshot{Sprites: sprites, Policy: p, IsEink: false}
	s := shapesRegistry.Lookup(bestStyleName)
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}

	return vd
}

// RenderCacheKey returns a change-detection string incorporating the monotonically
// increasing frame counter and the current policy fingerprint.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(fmt.Sprintf("attract_shapes:%d|%s", frameCounter, GetPolicy().Fingerprint()))
}

// resolveIsEink determines if the current panel is a slow-refresh (e-ink) panel
// based on the best-fit style from the registry.
func resolveIsEink(hints textlayout.TextHints) bool {
	s, _ := style.ResolveStyle(shapesRegistry, hints, "attract_shapes", "")
	reqs := s.Requirements()
	return reqs.Capability == style.MonoSlow ||
		reqs.Capability == style.GrayscaleSlow ||
		reqs.Capability == style.ColorSlow
}

// resolveBestStyleName returns the Name() of the best-fit style for the given hints.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(shapesRegistry, hints, "attract_shapes", "")
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
	if p.ShapeCount > 2 {
		p.ShapeCount = 2
	}
	if p.Complexity > 4 {
		p.Complexity = 4
	}
	if p.PulseRate > 0.8 {
		p.PulseRate = 0.8
	}
	if p.Speed > 0.85 {
		p.Speed = 0.85
	}
	return p
}
