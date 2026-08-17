package attract_plasma

import (
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/modes/attract_plasma/source"
)

func TestTunePolicyForUltraLowResCapsRates(t *testing.T) {
	p := tunePolicyForUltraLowRes(source.Policy{Speed: 1, BlobScale: 1.4, CycleRate: 1.7}, 16, 8)
	if p.Speed > 0.8 || p.BlobScale > 0.8 || p.CycleRate > 0.9 {
		t.Fatalf("unexpected tuned policy: %+v", p)
	}
}

func TestRenderPlasmaSpriteUltraLowResLeavesDarkPixels(t *testing.T) {
	sprite := renderPlasmaSprite(source.DefaultPolicy(), 16, 8, 1.2, false)
	dark := 0
	bounds := sprite.Image.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if sprite.Image.At(x, y) == (color.RGBA{0, 0, 0, 255}) {
				dark++
			}
		}
	}
	if dark == 0 {
		t.Fatalf("expected some dark pixels for ultra-low-res plasma output")
	}
}
