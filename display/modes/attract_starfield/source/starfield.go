package source

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// star represents a single star in the starfield.
type star struct {
	x     float64 // current x offset from center
	y     float64 // current y offset from center
	angle float64 // direction in radians
	speed float64 // base speed multiplier (varies by layer)
	layer int     // depth layer [0, layers-1]; higher = faster/brighter
	trail float64 // current trail length in pixels
}

// stars holds the active star population.
var stars []star

// starfieldInited tracks whether the star population has been initialized.
var starfieldInited bool

// initStars creates the initial star population for the given panel dimensions and policy.
func InitStars(panelWidth, panelHeight int, p Policy) {
	area := panelWidth * panelHeight
	count := int(float64(area) * p.Density / 512.0)
	if count < 4 {
		count = 4
	}
	if count > 500 {
		count = 500
	}

	stars = make([]star, count)
	for i := range stars {
		stars[i] = newStar(panelWidth, panelHeight, p, true)
	}
	starfieldInited = true
}

// newStar creates a new star at or near the center with a random direction and layer.
// If randomDist is true, the star starts at a random distance from center (for init).
func newStar(panelWidth, panelHeight int, p Policy, randomDist bool) star {
	angle := rand.Float64() * 2 * math.Pi
	layer := rand.Intn(p.Layers)

	// Layer affects speed: deeper layers (higher index) are faster.
	layerFactor := float64(layer+1) / float64(p.Layers)

	s := star{
		angle: angle,
		speed: 0.5 + layerFactor*1.5, // speed range [0.5, 2.0]
		layer: layer,
	}

	if randomDist {
		// Place at random distance from center for initial fill.
		maxDist := math.Max(float64(panelWidth), float64(panelHeight)) / 2.0
		dist := rand.Float64() * maxDist
		s.x = math.Cos(angle) * dist
		s.y = math.Sin(angle) * dist
		s.trail = dist * 0.15 * layerFactor
	}
	// If not randomDist, star starts at center (x=0, y=0, trail=0).

	return s
}

// advanceStars moves all stars outward from center by the given elapsed duration.
func AdvanceStars(elapsed time.Duration, p Policy, panelWidth, panelHeight int) {
	elapsedSec := elapsed.Seconds() * p.Speed
	halfW := float64(panelWidth) / 2.0
	halfH := float64(panelHeight) / 2.0

	for i := range stars {
		s := &stars[i]

		// Distance from center determines acceleration.
		dist := math.Sqrt(s.x*s.x + s.y*s.y)

		// Speed increases with distance from center (perspective effect).
		accel := 1.0 + dist*0.02
		movement := s.speed * accel * elapsedSec * 60.0

		s.x += math.Cos(s.angle) * movement
		s.y += math.Sin(s.angle) * movement

		// Trail length proportional to distance from center.
		newDist := math.Sqrt(s.x*s.x + s.y*s.y)
		layerFactor := float64(s.layer+1) / float64(p.Layers)
		s.trail = newDist * 0.15 * layerFactor

		// Check if star has exited the panel bounds.
		if math.Abs(s.x) > halfW || math.Abs(s.y) > halfH {
			// Respawn at center in a random direction and layer.
			stars[i] = newStar(panelWidth, panelHeight, p, false)
		}
	}
}

// renderStars draws all stars onto an image.RGBA canvas and returns it as a sprite.
func RenderStars(panelWidth, panelHeight int, p Policy) *widgets.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	centerX := float64(panelWidth) / 2.0
	centerY := float64(panelHeight) / 2.0
	maxDist := math.Sqrt(centerX*centerX + centerY*centerY)

	for _, s := range stars {
		// Star screen position.
		sx := int(centerX + s.x)
		sy := int(centerY + s.y)

		if sx < 0 || sx >= panelWidth || sy < 0 || sy >= panelHeight {
			continue
		}

		// Brightness proportional to distance from center (farther = brighter).
		dist := math.Sqrt(s.x*s.x + s.y*s.y)
		brightness := uint8(55 + int(200.0*(dist/maxDist)))

		c := color.RGBA{R: brightness, G: brightness, B: brightness, A: 255}

		// Draw star point.
		img.SetRGBA(sx, sy, c)

		// Draw trail (line from star back toward center).
		if s.trail > 1 {
			trailLen := int(s.trail)
			if trailLen > 20 {
				trailLen = 20
			}
			for t := 1; t <= trailLen; t++ {
				frac := float64(t) / float64(trailLen)
				tx := int(centerX + s.x - math.Cos(s.angle)*float64(t))
				ty := int(centerY + s.y - math.Sin(s.angle)*float64(t))
				if tx >= 0 && tx < panelWidth && ty >= 0 && ty < panelHeight {
					trailBright := uint8(float64(brightness) * (1.0 - frac*0.8))
					img.SetRGBA(tx, ty, color.RGBA{R: trailBright, G: trailBright, B: trailBright, A: 255})
				}
			}
		}
	}

	return &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Label:    "starfield",
	}
}

// staticSprite wraps a *widgets.Sprite to satisfy the widgets.Renderable interface.
type staticSprite struct {
	s *widgets.Sprite
}

func (ss *staticSprite) RenderFrame() *widgets.Sprite { return ss.s }

// StarfieldInited reports whether the star population has been initialized.
func StarfieldInited() bool { return starfieldInited }

// NewStaticSprite wraps a sprite as a widgets.Renderable.
func NewStaticSprite(s *widgets.Sprite) widgets.Renderable { return &staticSprite{s: s} }
