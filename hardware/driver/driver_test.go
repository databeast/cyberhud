package driver_test

import (
	"image"
	"image/draw"
	"testing"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type hintsTarget struct {
	bounds image.Rectangle
	hints  textlayout.TextHints
}

func (t *hintsTarget) Bounds() image.Rectangle { return t.bounds }

func (t *hintsTarget) DrawImage(img draw.Image) error { return nil }

func (t *hintsTarget) TextHints() textlayout.TextHints { return t.hints }

func TestResolveTextHintsProvider(t *testing.T) {
	target := &hintsTarget{
		bounds: image.Rect(0, 0, 250, 122),
		hints: textlayout.TextHints{
			SupportsVerticalScroll:   false,
			SupportsHorizontalScroll: false,
			SupportsAutoScroll:       false,
			PreferEventRefresh:       true,
			DefaultTickerDirection:   textlayout.TickerDirectionNone,
			DefaultLineMode:          textlayout.LineModeTruncate,
		},
	}
	h := textlayout.ResolveTextHints(target, target.TextHints)
	if h.PixelWidth != 250 || h.PixelHeight != 122 {
		t.Fatalf("TextHints size=%dx%d, want 250x122", h.PixelWidth, h.PixelHeight)
	}
	if h.SupportsVerticalScroll || h.SupportsHorizontalScroll || h.SupportsAutoScroll {
		t.Fatalf("TextHints scrolling flags not preserved: %+v", h)
	}
	if !h.PreferEventRefresh {
		t.Fatalf("PreferEventRefresh=false, want true")
	}
}
