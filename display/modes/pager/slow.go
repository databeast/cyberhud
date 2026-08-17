package pager

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sync"

	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/modes/pager/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// slowRenderState holds the cached rendering state for slow-refresh panels to
// produce stable cache keys when page content is unchanged.
var slowRenderState struct {
	mu        sync.Mutex
	cacheKey  string
	policyFP  string
	lastSeq   uint64
	lastPhase source.PagePhase
	lastAlpha int // quantized fade alpha (0-100)
}

// buildSlowPageView produces a ViewData for slow-refresh panels using the
// page-transition rendering strategy. It coordinates with the page phase and
// fade alpha to render the correct frame:
//   - phaseIdle: render current page at full opacity
//   - phaseFadeOut: render current page with decreasing opacity
//   - phaseFadeIn: render next page with increasing opacity
//
// Text is rendered with dithering-friendly contrast: full black on white for
// mono slow-refresh panels, and grayscale-appropriate values for grayscale panels.
// A deterministic (seeded) render ensures stable output for unchanged content.
func buildSlowPageView(
	hints textlayout.TextHints,
	layout styles.Layout,
	lines []string,
	nextLines []string,
	phase source.PagePhase,
	fadeAlpha float64,
	seq uint64,
	pol source.Policy,
) style.ViewData {
	panelWidth := layout.PixelWidth
	panelHeight := layout.PixelHeight

	if panelWidth <= 0 || panelHeight <= 0 {
		return style.ViewData{
			Static: true,
			Items:  []string{},
		}
	}

	// Determine panel type from capability.
	isMono := hints.Capability == textlayout.CapMonoSlow

	// Select which lines to render and at what alpha based on phase.
	var renderLines []string
	var alpha float64

	switch phase {
	case source.PhaseIdle:
		renderLines = lines
		alpha = 1.0
	case source.PhaseFadeOut:
		renderLines = lines
		alpha = 1.0 - fadeAlpha // fadeAlpha goes 0→1, so opacity decreases
	case source.PhaseFadeIn:
		renderLines = nextLines
		alpha = fadeAlpha // fadeAlpha goes 0→1, so opacity increases
	default:
		renderLines = lines
		alpha = 1.0
	}

	// Render the text page to an image.
	img := renderSlowTextPage(panelWidth, panelHeight, layout, renderLines, alpha, isMono, seq)

	// Update cache state.
	quantizedAlpha := int(fadeAlpha * 100)
	fp := policyFingerprint(pol)
	slowRenderState.mu.Lock()
	slowRenderState.cacheKey = fmt.Sprintf("pager:slow:%d:%d:%d:%s", seq, phase, quantizedAlpha, fp)
	slowRenderState.policyFP = fp
	slowRenderState.lastSeq = seq
	slowRenderState.lastPhase = phase
	slowRenderState.lastAlpha = quantizedAlpha
	slowRenderState.mu.Unlock()

	// Use compositor with slow-refresh suppression context.
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})
	compositor.Add(&slowPageRenderable{sprite: &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
		Label:    "pager-slow-page",
	}})

	return style.ViewData{
		Static:      true,
		Sprites:     compositor.Sprites(),
		StyleReport: style.StyleReport{Name: "slow-pager", Reason: "fitness"},
	}
}

// renderSlowTextPage renders lines of text onto an image optimized for
// slow-refresh display characteristics. Uses a deterministic seed for stable
// output and high-contrast rendering (black text on white background).
func renderSlowTextPage(
	panelWidth, panelHeight int,
	lay styles.Layout,
	lines []string,
	alpha float64,
	isMono bool,
	seq uint64,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	// Background is black (zero-value RGBA). Modes render with standard
	// "light on dark" convention — the panel driver handles any inversion
	// needed for the physical display technology.

	if len(lines) == 0 || lay.VisibleRows == 0 || lay.VisibleColumns == 0 {
		return img
	}

	// Use a deterministic RNG seeded from the sequence number for stable frames.
	_ = rand.New(rand.NewSource(int64(seq)))

	// Resolve the font face for text rendering.
	face := resolveSlowFont(lay.RowHeight)
	if face == nil {
		return img
	}

	// Compute text foreground color based on panel type and alpha.
	fgColor := slowForegroundColor(alpha, isMono)

	// Render each visible line of text.
	metrics := face.Metrics()
	maxX := panelWidth
	visibleRows := lay.VisibleRows
	if visibleRows > len(lines) {
		visibleRows = len(lines)
	}

	for row := 0; row < visibleRows; row++ {
		x := 0
		y := row * lay.RowHeight

		// Center the glyph vertically within the row if row height exceeds glyph height.
		yOffset := 0
		if lay.RowHeight > metrics.GlyphHeight {
			yOffset = (lay.RowHeight - metrics.GlyphHeight) / 2
		}

		line := lines[row]
		font.DrawText(img, face, line, x, y+yOffset, fgColor, maxX)
	}

	return img
}

// slowForegroundColor computes the text foreground color for slow-refresh rendering.
// Uses white/light text on dark background — standard mode convention.
// The panel driver handles any inversion for the physical display.
func slowForegroundColor(alpha float64, isMono bool) color.RGBA {
	a := uint8(alpha * 255)
	return color.RGBA{R: 255, G: 255, B: 255, A: a}
}

// resolveSlowFont returns an appropriate font face for slow-refresh text rendering.
// Falls back to the default system font.
func resolveSlowFont(rowHeight int) font.Face {
	if rowHeight <= 0 {
		return font.Default()
	}
	return font.ByHeight(rowHeight)
}

// slowRenderCacheKey returns the stable cache key for slow-refresh pager rendering.
func slowRenderCacheKey(seq uint64, phase source.PagePhase, fadeAlpha float64, pol source.Policy) string {
	quantizedAlpha := int(fadeAlpha * 100)
	fp := policyFingerprint(pol)
	return fmt.Sprintf("pager:slow:%d:%d:%d:%s", seq, phase, quantizedAlpha, fp)
}

// slowPageRenderable wraps a Sprite for use with the widget Compositor.
type slowPageRenderable struct {
	sprite *widgets.Sprite
}

func (r *slowPageRenderable) RenderFrame() *widgets.Sprite {
	return r.sprite
}
