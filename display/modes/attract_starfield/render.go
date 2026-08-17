package attract_starfield

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_starfield/source"
	"github.com/databeast/cyberhud/display/modes/attract_starfield/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// maxTickElapsed caps the per-frame advancement to prevent visual jumps.
const maxTickElapsed = 80 * time.Millisecond

// BuildView produces a complete starfield frame as style.ViewData.
func BuildView(hints textlayout.TextHints) style.ViewData {
	p := GetPolicy()
	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight

	// Graceful fallback if no dimensions available.
	if panelWidth == 0 || panelHeight == 0 {
		panelWidth = 128
		panelHeight = 128
	}
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	// E-ink path: produce static starfield scatter.
	if resolveIsEink(hints) {
		s, reason := style.ResolveStyle(starfieldRegistry, hints, "attract_starfield", "")
		vd := styles.BuildEinkView(panelWidth, panelHeight, p)
		if s != nil {
			vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
		}
		return vd
	}

	// Initialize star population if needed.
	if !source.StarfieldInited() {
		source.InitStars(panelWidth, panelHeight, p)
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

	// Advance stars outward from center.
	source.AdvanceStars(elapsed, p, panelWidth, panelHeight)

	// Render stars to sprite.
	sprite := source.RenderStars(panelWidth, panelHeight, p)

	// Build compositor.
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          false,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(source.NewStaticSprite(sprite))

	frameCounter++

	ctx := style.NewStyleContext(hints)
	s, reason := style.ResolveStyle(starfieldRegistry, hints, "attract_starfield", "")
	snap := source.Snapshot{Policy: p, IsEink: false}
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.Sprites = compositor.Sprites()
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}

	return vd
}

// RenderCacheKey returns a change-detection key incorporating the frame counter
// and policy fingerprint.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(fmt.Sprintf("attract_starfield:%d|%s", frameCounter, policyFingerprint(GetPolicy())))
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
	if p.Layers > 2 {
		p.Layers = 2
	}
	if p.Speed > 0.9 {
		p.Speed = 0.9
	}
	return p
}
