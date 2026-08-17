package source

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

// particle represents a single drifting firefly in the particle system.
type particle struct {
	x, y       float64 // position in pixels
	dir        float64 // direction in radians
	speedScale float64 // speed factor [0.5, 1.0]
	colorPhase float64 // hue phase [0, 1)
}

// particles holds the active particle array, lazily initialized on first frame.
var particles []particle

// ParticleCount computes the particle count for the given panel dimensions and density.
// Formula: clamp(area/512 * density, 4, 200)
func ParticleCount(width, height int, density float64) int {
	area := width * height
	count := int(float64(area) / 512.0 * density)
	if count < 4 {
		count = 4
	}
	if count > 200 {
		count = 200
	}
	return count
}

// InitParticles creates or resizes the particle slice for the given panel dimensions.
func InitParticles(width, height int, density float64) {
	count := ParticleCount(width, height, density)
	if len(particles) == count {
		return
	}
	particles = make([]particle, count)
	for i := range particles {
		particles[i] = particle{
			x:          rand.Float64() * float64(width),
			y:          rand.Float64() * float64(height),
			dir:        rand.Float64() * 2 * math.Pi,
			speedScale: 0.5 + rand.Float64()*0.5,
			colorPhase: rand.Float64(),
		}
	}
}

// AdvanceParticles moves each particle by elapsed time, applying drift and edge wrapping.
func AdvanceParticles(elapsed time.Duration, p Policy, width, height int) {
	dt := elapsed.Seconds()
	fw := float64(width)
	fh := float64(height)

	for i := range particles {
		pt := &particles[i]

		// Apply drift randomness to direction.
		pt.dir += (rand.Float64()*2 - 1) * p.Drift * math.Pi * dt

		// Compute velocity.
		speed := pt.speedScale * p.Speed * 30.0 // base pixels/sec at speed=1
		dx := math.Cos(pt.dir) * speed * dt
		dy := math.Sin(pt.dir) * speed * dt

		pt.x += dx
		pt.y += dy

		// Advance color phase proportional to glow.
		pt.colorPhase += p.Glow * dt * 0.5
		if pt.colorPhase >= 1.0 {
			pt.colorPhase -= 1.0
		}

		// Edge wrapping: exit one edge, re-enter at random position on opposite edge.
		if pt.x < 0 {
			pt.x = fw - 1
			pt.y = rand.Float64() * fh
		} else if pt.x >= fw {
			pt.x = 0
			pt.y = rand.Float64() * fh
		}
		if pt.y < 0 {
			pt.y = fh - 1
			pt.x = rand.Float64() * fw
		} else if pt.y >= fh {
			pt.y = 0
			pt.x = rand.Float64() * fw
		}
	}
}

// HueToRGB converts an HSV hue (0-1) to an RGB color with full saturation and value.
func HueToRGB(hue float64) color.RGBA {
	h := hue * 6.0
	sector := int(h) % 6
	f := h - float64(int(h))
	var r, g, b float64
	switch sector {
	case 0:
		r, g, b = 1, f, 0
	case 1:
		r, g, b = 1-f, 1, 0
	case 2:
		r, g, b = 0, 1, f
	case 3:
		r, g, b = 0, 1-f, 1
	case 4:
		r, g, b = f, 0, 1
	case 5:
		r, g, b = 1, 0, 1-f
	}
	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}

// RenderParticlesColor renders particles as small filled circles on an RGBA image.
func RenderParticlesColor(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	radius := 3
	for _, pt := range particles {
		c := HueToRGB(pt.colorPhase)
		cx := int(pt.x)
		cy := int(pt.y)
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
	return img
}

// RenderParticlesMono renders particles as single-color circles with brightness
// varying by glow intensity.
func RenderParticlesMono(width, height int, glow float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	radius := 3
	brightness := uint8(128 + int(glow*127))
	c := color.RGBA{R: brightness, G: brightness, B: brightness, A: 255}
	for _, pt := range particles {
		cx := int(pt.x)
		cy := int(pt.y)
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
	return img
}
