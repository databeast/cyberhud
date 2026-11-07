package attract_cube

import (
	"bytes"
	"image"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestMonoCapabilities(t *testing.T) {
	resetPolicy(t)
	restore := frameclock.Freeze(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	defer restore()
	i := &instance{}
	i.SetPanelHints(textlayout.TextHints{PixelWidth: 128, PixelHeight: 64, Capability: textlayout.CapMonoFast})
	i.Activate()
	firstKey := i.RenderCacheKey()
	first := i.BuildView()
	if len(first.Sprites) != 1 {
		t.Fatal("mono-fast must produce an animated sprite")
	}
	if err := SetPolicy(Policy{BaseColor: "#000000"}); err != nil {
		t.Fatal(err)
	}
	if firstKey != i.RenderCacheKey() || !bytes.Equal(first.Sprites[0].Image.(*image.RGBA).Pix, i.BuildView().Sprites[0].Image.(*image.RGBA).Pix) {
		t.Fatal("mono rendering must not depend on base_color")
	}
	frameclock.Advance(4 * time.Second)
	next := i.BuildView()
	if firstKey == i.RenderCacheKey() || bytes.Equal(first.Sprites[0].Image.(*image.RGBA).Pix, next.Sprites[0].Image.(*image.RGBA).Pix) {
		t.Fatal("mono-fast frame and cache key must advance")
	}
	i.SetPanelHints(textlayout.TextHints{PixelWidth: 128, PixelHeight: 64, Capability: textlayout.CapMonoSlow})
	slowKey := i.RenderCacheKey()
	slow := i.BuildView()
	frameclock.Advance(4 * time.Second)
	if len(slow.Sprites) != 1 || slowKey != i.RenderCacheKey() ||
		!bytes.Equal(slow.Sprites[0].Image.(*image.RGBA).Pix, i.BuildView().Sprites[0].Image.(*image.RGBA).Pix) {
		t.Fatal("mono-slow frame and cache key must remain static")
	}
	if err := SetPolicy(Policy{BaseColor: "red"}); err != nil {
		t.Fatal(err)
	}
	if slowKey != i.RenderCacheKey() {
		t.Fatal("mono-slow cache key must ignore color changes")
	}
}
