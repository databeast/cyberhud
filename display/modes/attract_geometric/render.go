package attract_geometric

import (
	"image"
	"strings"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
	"github.com/databeast/cyberhud/display/modes/attract_geometric/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets"
)

// BuildView produces a complete geometric frame as style.ViewData.
func BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}

	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight

	// Degrade gracefully if hints are not yet set.
	if panelWidth == 0 || panelHeight == 0 {
		panelWidth = 240
		panelHeight = 240
	}

	if strings.Contains(styles.ResolveBestStyleName(hints), "-slow-") {
		return BuildSlowView(hints)
	}

	animState.Lock()
	defer animState.Unlock()

	// Initialize on first call or if panel dimensions changed.
	if !animState.initialized || animState.panelW != panelWidth || animState.panelH != panelHeight {
		initAnimState(panelWidth, panelHeight)
	}

	// Advance frame tick.
	tickFrame()

	// Get policy for glow intensity.
	policy := tunePolicyForUltraLowRes(GetPolicy(), panelWidth, panelHeight)

	// Render frame.
	img := source.RenderFrame(animState.clusters, animState.fragmentState.ActiveFragments, animState.time, panelWidth, panelHeight, animState.scaleFactor, policy.GlowIntensity)

	// Construct bridge and StyleContext per architecture.
	ctx := style.NewStyleContext(hints)

	snap := source.GeometricFrame{Policy: policy}
	s := styles.Registry().Lookup(styles.ResolveBestStyleName(hints))
	vd := s.Build(snap, policy, ctx)

	vd.Static = false
	vd.Sprites = []widgets.Sprite{
		{
			Image:    img,
			Position: image.Point{},
			Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
			Label:    "geometric-clusters",
		},
	}
	vd.StyleReport = style.StyleReport{Name: styles.ResolveBestStyleName(hints), Reason: "fitness"}

	return vd
}

// RenderCacheKey returns a change-detection string for the display runtime.
func RenderCacheKey() uint32 {
	hints, ok := getPanelHints()
	if ok && strings.Contains(styles.ResolveBestStyleName(hints), "-slow-") {
		return SlowRenderCacheKey()
	}
	animState.Lock()
	defer animState.Unlock()
	return region.CalcRegionCacheKey(animState.frameCount, policyFingerprint(GetPolicy()))
}
