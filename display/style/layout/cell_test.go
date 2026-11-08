package layout_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestCellPixelRoundTrip(t *testing.T) {
	lb := layout.NewLayoutBridge(textlayout.TextHints{PixelWidth: 400, PixelHeight: 200}, layout.BridgeConfig{PaddingPct: 10})
	m := layout.NewMetrics(8, 8, 16, "test")

	for row := 0; row < 5; row++ {
		for col := 0; col < 10; col++ {
			r := lb.CellRect(row, col, m)
			if r.Min != lb.CellToPixel(row, col, m) || r.Dx() != 8 || r.Dy() != 16 {
				t.Fatalf("CellRect(%d,%d) = %v", row, col, r)
			}
			for _, p := range []image.Point{r.Min, r.Max.Sub(image.Pt(1, 1))} {
				gr, gc, ok := lb.PixelToCell(p, m)
				if !ok || gr != row || gc != col {
					t.Fatalf("PixelToCell(%v) = %d,%d,%v want %d,%d", p, gr, gc, ok, row, col)
				}
			}
		}
	}
}

func TestPixelToCellOutside(t *testing.T) {
	lb := layout.NewLayoutBridge(textlayout.TextHints{PixelWidth: 400, PixelHeight: 200}, layout.BridgeConfig{PaddingPct: 10})
	m := layout.NewMetrics(8, 8, 16, "test")
	ox, oy := lb.ContentOrigin()

	for _, p := range []image.Point{{ox - 1, oy}, {ox, oy - 1}, {ox + lb.AvailableContentWidth(), oy}} {
		if _, _, ok := lb.PixelToCell(p, m); ok {
			t.Fatalf("PixelToCell(%v) ok, want false", p)
		}
	}
	if _, _, ok := lb.PixelToCell(image.Pt(ox, oy), layout.Metrics{}); ok {
		t.Fatal("unusable metrics should return ok=false")
	}
}
