package styles

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/databeast/cyberhud/display/modes/attract_particles/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets"
)

// einkCacheKey holds the stable RenderCacheKey for e-ink panels.
// It changes only when the policy is modified.
var einkCacheKey string

// einkPolicyFP tracks the last policy fingerprint used for e-ink rendering.
var einkPolicyFP string

// BuildEinkView produces a static particle scatter frame for e-ink panels.
// It renders a frozen decorative frame without advancing animation state.
func BuildEinkView(width, height int, p source.Policy, isMono bool) style.ViewData {
	fp := p.Fingerprint()

	// Only regenerate when policy changes.
	if einkCacheKey == "" || einkPolicyFP != fp {
		einkPolicyFP = fp
		einkCacheKey = fmt.Sprintf("attract_particles:eink|%s", fp)
	}

	// Generate a static scatter of particles.
	count := source.ParticleCount(width, height, p.Density)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	radius := 3

	for i := 0; i < count; i++ {
		cx := rand.Intn(width)
		cy := rand.Intn(height)

		var c color.RGBA
		if isMono {
			brightness := uint8(128 + int(p.Glow*127))
			c = color.RGBA{R: brightness, G: brightness, B: brightness, A: 255}
		} else {
			// Static hue distribution evenly spaced.
			hue := float64(i) / float64(count)
			c = source.HueToRGB(hue)
		}

		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				if dx*dx+dy*dy <= radius*radius {
					px, py := cx+dx, cy+dy
					if px >= 0 && px < width && py >= 0 && py < height {
						img.SetRGBA(px, py, c)
					}
				}
			}
		}
	}

	return style.ViewData{
		Static:      true,
		Sprites:     []widgets.Sprite{{Image: img, Position: image.Point{}, Label: "particles-eink"}},
		StyleReport: style.StyleReport{Name: "eink-fallback", Reason: "fitness"},
	}
}
