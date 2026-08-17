package attract_geometric

import (
	"image"
	"math"
	"sync"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// slowRefreshState caches the static frame for slow-refresh panels.
var slowRefreshState = struct {
	sync.Mutex
	img         *image.RGBA
	fingerprint string
	panelW      int
	panelH      int
}{}

// BuildSlowView produces a static frozen frame for slow-refresh panels.
// It renders one frame at time=0 with all clusters at fadeInFactor=1.0,
// omits all pseudocode text fragments, and returns Static=true.
// The cached frame is only re-rendered when the policy fingerprint or
// panel dimensions change.
func BuildSlowView(hints textlayout.TextHints) style.ViewData {
	panelWidth := hints.PixelWidth
	panelHeight := hints.PixelHeight

	if panelWidth == 0 || panelHeight == 0 {
		panelWidth = 240
		panelHeight = 240
	}

	policy := tunePolicyForUltraLowRes(GetPolicy(), panelWidth, panelHeight)
	fp := policyFingerprint(policy)

	slowRefreshState.Lock()
	defer slowRefreshState.Unlock()

	// Return cached frame if policy and panel dimensions haven't changed.
	if slowRefreshState.img != nil &&
		slowRefreshState.fingerprint == fp &&
		slowRefreshState.panelW == panelWidth &&
		slowRefreshState.panelH == panelHeight {
		return style.ViewData{
			Static: true,
			Sprites: []widgets.Sprite{
				{
					Image:    slowRefreshState.img,
					Position: image.Point{},
					Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
					Label:    "geometric-static",
				},
			},
		}
	}

	// Generate clusters deterministically for the static frame.
	rng := newRNG()
	baseCount := int(math.Round(float64(10) * policy.Density))
	clusters := source.InitializeClusters(panelWidth, panelHeight, baseCount, rng)

	// Force all clusters to fadeInFactor=1.0 at time=0 by setting SpawnTime
	// to negative FadeInDuration. This ensures:
	//   elapsed = 0 - (-FadeInDuration) = FadeInDuration
	//   fadeInFactor = FadeInDuration / FadeInDuration = 1.0
	for i := range clusters {
		clusters[i].SpawnTime = -clusters[i].FadeInDuration
	}

	scaleFactor := computeScaleFactor(panelWidth, panelHeight)

	// Render at time=0 with no fragments (nil).
	img := source.RenderFrame(clusters, nil, 0, panelWidth, panelHeight, scaleFactor, policy.GlowIntensity)

	// Cache the result.
	slowRefreshState.img = img
	slowRefreshState.fingerprint = fp
	slowRefreshState.panelW = panelWidth
	slowRefreshState.panelH = panelHeight

	return style.ViewData{
		Static: true,
		Sprites: []widgets.Sprite{
			{
				Image:    img,
				Position: image.Point{},
				Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
				Label:    "geometric-static",
			},
		},
	}
}

// SlowRenderCacheKey returns a cache key that only changes on policy fingerprint change.
func SlowRenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(0, policyFingerprint(GetPolicy()))
}
