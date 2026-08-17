package attract_matrix

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// BuildView produces a complete matrix rain frame as style.ViewData.
// It is called by the display runtime each render cycle with the panel's text hints.
func BuildView(hints textlayout.TextHints) style.ViewData {
	// Set frame-level timestamp for CharAt mutation checks (avoids per-cell syscalls).
	source.SetFrameTime()

	// Compute elapsed time since last frame, capped to avoid visual jumps.
	var elapsed time.Duration
	if !lastTick.IsZero() {
		elapsed = time.Since(lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	lastTick = time.Now()

	// Read current policy snapshot and apply rain intensity cycle.
	p := source.ApplyCycle(GetPolicy())

	// Panel dimensions from hints.
	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight

	// Build StyleContext once for all font resolution in this frame.
	ctx := style.NewStyleContext(hints)

	// Font metrics via StyleContext (catalog-validated).
	// The matrix mode's family preference ("matrix") resolves to matrix-code font
	// through the tier catalog, guaranteeing GlyphAdvance ≤ maxAdvance.
	glyphAdvance, rowHeight, spriteFace := source.ResolveMatrixFont(ctx)

	// Resolve panel type from the best-fitting registered style.
	isMono, isEink := resolvePanelType(hints)

	// Resolve the best-fit style name for StyleReport.
	bestStyleName := resolveBestStyleName(hints)

	// Build snapshot for style dispatch.
	snapshot := source.MatrixSnapshot{
		Policy:       p,
		PanelWidth:   panelWidth,
		PanelHeight:  panelHeight,
		GlyphAdvance: glyphAdvance,
		RowHeight:    rowHeight,
		Mono:         isMono,
		Eink:         isEink,
	}

	// E-ink path: produce a static frozen frame with no animation.
	if isEink {
		strips := source.GetOrRebuildStrips(p, panelWidth, panelHeight, glyphAdvance, rowHeight, isMono, 0, spriteFace)
		compositor := widgets.NewCompositor(widgets.SuppressionContext{
			IsEink:          true,
			AvailableWidth:  panelWidth,
			AvailableHeight: panelHeight,
		})
		for _, strip := range strips {
			compositor.Add(strip)
		}
		frameCounter++

		// Reuse the frame-level StyleContext.
		s := matrixRegistry.Lookup(bestStyleName)
		vd := s.Build(snapshot, p, ctx)

		vd.Static = true
		vd.Sprites = compositor.Sprites()
		vd.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}

		return vd
	}

	// Ultra-low-resolution path (e.g. 16x8 or 8x16 CharliePlex): render a dedicated
	// sparse pixel-rain animation instead of glyph strips, which are too dense
	// at this size and tend to wash out the panel.
	if isUltraLowResPanel(panelWidth, panelHeight) {
		frameCounter++
		return buildUltraLowResView(snapshot, bestStyleName)
	}

	// Animated path (color TFT or monochrome OLED).
	strips := source.GetOrRebuildStrips(p, panelWidth, panelHeight, glyphAdvance, rowHeight, isMono, 0, spriteFace)

	// Tick all strips forward by elapsed time, scaled by cycle speed factor.
	scaledElapsed := time.Duration(float64(elapsed) * source.CycleSpeedFactor())
	for _, strip := range strips {
		strip.Tick(scaledElapsed)
	}

	// Build compositor for this frame.
	compositor := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          false,
		AvailableWidth:  panelWidth,
		AvailableHeight: panelHeight,
	})

	// Add background gradient if enabled and panel is color-capable.
	if p.ShowBackground && !isMono {
		bgSprite := gradient.Render(gradient.Config{
			Style:  gradient.Radial,
			Bounds: image.Rect(0, 0, panelWidth, panelHeight),
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{0, 20, 0, 255}},
				{Position: 1.0, Color: color.RGBA{0, 0, 0, 255}},
			},
		})
		if bgSprite != nil {
			compositor.Add(source.NewStaticSprite(bgSprite))
		}
	}

	// Add each strip to the compositor.
	for _, strip := range strips {
		compositor.Add(strip)
	}

	// Add "CYBERHUD" splash at cycle peak.
	// Splash uses the LARGEST catalog-validated font (spleen/TierFullsize) for dramatic
	// rendering, NOT the rain column's matrix-code font.
	if alpha := source.SplashAlpha(source.CycleProgress()); alpha > 0 {
		splashFace := source.ResolveSplashFont(ctx)
		if splash := source.RenderSplash(panelWidth, panelHeight, alpha, splashFace); splash != nil {
			compositor.Add(source.NewStaticSprite(splash))
		}
	}

	frameCounter++

	// Reuse the frame-level StyleContext.
	sAnim := matrixRegistry.Lookup(bestStyleName)
	vdAnim := sAnim.Build(snapshot, p, ctx)

	vdAnim.Static = false
	vdAnim.Sprites = compositor.Sprites()
	vdAnim.StyleReport = style.StyleReport{Name: bestStyleName, Reason: "fitness"}

	return vdAnim
}

func buildUltraLowResView(snapshot source.MatrixSnapshot, styleName string) style.ViewData {
	w := snapshot.PanelWidth
	h := snapshot.PanelHeight
	if w <= 0 || h <= 0 {
		return style.ViewData{
			Static:      false,
			StyleReport: style.StyleReport{Name: styleName, Reason: "fitness"},
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	columns := clampInt(w/4, 1, 4)
	tailLen := clampInt(snapshot.Policy.TrailLength/8, 2, h)

	sweepX := int(frameCounter % uint64(w))
	sweepY := int((frameCounter / 2) % uint64(h))

	// Moving crosshair anchor makes motion obvious even on tiny 16x8 matrices.
	img.SetRGBA(sweepX, sweepY, color.RGBA{255, 255, 255, 255})
	if sweepX > 0 {
		img.SetRGBA(sweepX-1, sweepY, color.RGBA{120, 120, 120, 255})
	}

	// Rain columns sit at fixed, evenly spaced x positions and their drop heads
	// advance downward one pixel per frame (staggered per column), wrapping at
	// the bottom with a fading trail — coherent falling motion rather than
	// per-frame pseudo-random repositioning.
	colSpacing := w / columns
	for i := 0; i < columns; i++ {
		x := i*colSpacing + colSpacing/2
		head := int((frameCounter + uint64(i)*uint64(h)/uint64(columns)) % uint64(h))

		for t := 0; t < tailLen; t++ {
			y := head - t
			for y < 0 {
				y += h
			}
			alpha := 255 - (t * 190 / tailLen)
			if alpha < 24 {
				alpha = 24
			}
			v := uint8(alpha)
			img.SetRGBA(x, y, color.RGBA{v, v, v, 255})
		}
	}

	return style.ViewData{
		Static: false,
		Sprites: []widgets.Sprite{
			{
				Image:    img,
				Position: image.Point{},
				Label:    "matrix-lowres",
			},
		},
		StyleReport: style.StyleReport{Name: styleName, Reason: "fitness"},
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func isUltraLowResPanel(width, height int) bool {
	if width <= 0 || height <= 0 {
		return false
	}
	longEdge := width
	shortEdge := height
	if height > width {
		longEdge = height
		shortEdge = width
	}
	return longEdge <= 16 && shortEdge <= 8
}

// RenderCacheKey returns a change-detection string incorporating the monotonically
// increasing frame counter and the current policy fingerprint. The display runtime
// compares consecutive keys to determine whether a re-render is needed.
// Because frameCounter increments every BuildView call, animated frames always
// produce a new signature.
func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey("attract_matrix", frameCounter, source.PolicyFingerprint(GetPolicy()))
}
