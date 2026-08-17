package testfonts

import (
	"image"
	"image/color"
	"testing"
	"time"

	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestBuildView_NoCatalogStillRendersContent(t *testing.T) {
	hints := textlayout.TextHints{PixelWidth: 128, PixelHeight: 64}
	view := BuildView(hints, time.Unix(0, 0))
	if view.Image == nil {
		t.Fatal("BuildView returned nil image")
	}
	if view.Image.Bounds().Dx() != 128 || view.Image.Bounds().Dy() != 64 {
		t.Fatalf("image bounds = %dx%d, want 128x64", view.Image.Bounds().Dx(), view.Image.Bounds().Dy())
	}
	if countColoredPixels(view.Image, color.RGBA{0x00, 0x00, 0x00, 0xFF}) >= view.Image.Bounds().Dx()*view.Image.Bounds().Dy() {
		t.Fatal("BuildView rendered a fully black image; content should remain visible without a catalog")
	}
}

func TestTierRows_NoCatalogBuildsFallbackCatalog(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	font.Register(testFace{id: "spleen-5x8", metrics: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}})
	font.Register(testFace{id: "spleen-8x16", metrics: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 9, RowHeight: 18}})

	hints := textlayout.TextHints{PixelWidth: 128, PixelHeight: 64}
	rows := tierRows(hints)
	if len(rows) < 2 {
		t.Fatalf("tierRows() without catalog = %d rows, want at least 2 tiers", len(rows))
	}
	if rows[0].label == "small" && rows[1].label == "small" {
		t.Fatal("fallback catalog still collapsed to a single repeated tier; want a true tier stack")
	}
}

type testFace struct {
	id      string
	metrics font.Metrics
}

func (f testFace) ID() string                { return f.id }
func (f testFace) Metrics() font.Metrics     { return f.metrics }
func (f testFace) GlyphRow(rune, int) uint32 { return 0 }

func countColoredPixels(img *image.RGBA, target color.RGBA) int {
	count := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.RGBAAt(x, y) == target {
				count++
			}
		}
	}
	return count
}
