package attract_particles

import (
	"fmt"
	"image"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_particles/source"
	"github.com/databeast/cyberhud/display/modes/attract_particles/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// buildView produces a complete particles frame as style.ViewData.
func buildView(hints textlayout.TextHints) style.ViewData {
	p := GetPolicy()
	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight
	p = tunePolicyForUltraLowRes(p, panelWidth, panelHeight)

	// Resolve panel type from best-fit style.
	isMono, isEink := resolvePanelType(hints)

	// E-ink path: produce a static frozen frame.
	if isEink {
		return styles.BuildEinkView(panelWidth, panelHeight, p, isMono)
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

	// Initialize or resize particle array.
	source.InitParticles(panelWidth, panelHeight, p.Density)

	// Advance particle positions.
	source.AdvanceParticles(elapsed, p, panelWidth, panelHeight)

	// Render particles to image.
	var img *image.RGBA
	if isMono {
		img = source.RenderParticlesMono(panelWidth, panelHeight, p.Glow)
	} else {
		img = source.RenderParticlesColor(panelWidth, panelHeight)
	}

	// Compose into sprites via Compositor.
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          false,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(&staticSprite{sprite: &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Label:    "particles-frame",
	}})

	frameCounter++

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	bestStyleName := resolveBestStyleName(hints)
	snap := source.Snapshot{Sprites: compositor.Sprites(), Policy: p, IsEink: false}
	s := particlesRegistry.Lookup(bestStyleName)
	vd := s.Build(snap, p, ctx)

	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}

	return vd
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
	if p.Drift > 0.2 {
		p.Drift = 0.2
	}
	if p.Glow > 0.35 {
		p.Glow = 0.35
	}
	if p.Speed > 0.85 {
		p.Speed = 0.85
	}
	return p
}

// renderCacheKey returns a change-detection string for the particles mode.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(fmt.Sprintf("attract_particles:%d|%s", frameCounter, GetPolicy().Fingerprint()))
}

type staticSprite struct {
	sprite *widgets.Sprite
}

func (s *staticSprite) RenderFrame() *widgets.Sprite { return s.sprite }
