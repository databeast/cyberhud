package attract_bokeh

import (
	"image"
	"log"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

// BuildView produces a complete bokeh frame as style.ViewData.
func BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		log.Fatal("Could not find hints")
		return style.ViewData{}
	}
	p := GetPolicy()

	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	// Resolve panel type from best-fit style.
	isMono, isEink := resolvePanelType(hints)
	// E-ink path: static frozen frame.
	if isEink {
		return source.RenderEInkBokehPicture(panelWidth, panelHeight, p)
	}

	// Compute elapsed time.
	var elapsed time.Duration
	if !lastTick.IsZero() {
		elapsed = time.Since(lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	lastTick = time.Now()

	// Initialize or reinitialize circles if panel dimensions changed.
	source.EnsureCircles(panelWidth, panelHeight, p)

	// Advance circles.
	source.AdvanceCircles(elapsed, p.Speed, panelWidth, panelHeight)

	// Render circles to image.
	img := source.RenderCircles(panelWidth, panelHeight, p, isMono)

	frameCounter++

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	snap := source.BokehFrame{IsEink: isEink}
	s := bokehRegistry.Lookup(resolveBestStyleName(hints))
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.Sprites = []widgets.Sprite{
		{
			Image:    img,
			Position: image.Point{},
			Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
			Label:    "bokeh-circles",
		},
	}

	vd.StyleReport = style.StyleReport{Name: resolveBestStyleName(hints), Reason: "fitness"}

	return vd
}

// RenderCacheKey returns a change-detection string for the display runtime.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(frameCounter, GetPolicy().Fingerprint())
}

func tunePolicyForUltraLowRes(p source.Policy, panelWidth, panelHeight int) source.Policy {
	longEdge, shortEdge := panelWidth, panelHeight
	if panelHeight > panelWidth {
		longEdge, shortEdge = panelHeight, panelWidth
	}
	if longEdge > 16 || shortEdge > 8 {
		return p
	}
	if p.Density > 0.2 {
		p.Density = 0.2
	}
	if p.SizeVariance > 0.35 {
		p.SizeVariance = 0.35
	}
	if p.Saturation > 0.6 {
		p.Saturation = 0.6
	}
	if p.Speed > 0.85 {
		p.Speed = 0.85
	}
	return p
}
