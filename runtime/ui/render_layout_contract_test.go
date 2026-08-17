package ui

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// These tests guard the ViewData layout contract described on style.ViewData: a
// style solves layout and the renderer draws what it was given.
//
// They exist because that contract has been broken three separate ways, each time
// silently and each time only visible as mis-positioned text on hardware:
// per-row font heights collapsed to one global height, inter-row spacing dropped,
// and per-row horizontal offsets ignored on the non-static path. Unit tests that
// only assert "renders without error" cannot catch any of those, so these assert
// on actual pixel positions.
//
// Both draw paths are covered. renderItems regressed independently of
// renderStaticItems precisely because nothing pinned its behaviour.

// inkRows returns the inclusive y ranges of rows containing any lit pixel, so a
// test can assert where text actually landed rather than trusting the call
// completed.
func inkRows(img *image.RGBA) [][2]int {
	b := img.Bounds()
	var bands [][2]int
	inBand := false
	start := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		lit := false
		for x := b.Min.X; x < b.Max.X; x++ {
			if r, g, bl, _ := img.At(x, y).RGBA(); r > 0x2000 || g > 0x2000 || bl > 0x2000 {
				lit = true
				break
			}
		}
		if lit && !inBand {
			inBand, start = true, y
		} else if !lit && inBand {
			bands = append(bands, [2]int{start, y - 1})
			inBand = false
		}
	}
	if inBand {
		bands = append(bands, [2]int{start, b.Max.Y - 1})
	}
	return bands
}

// firstLitX returns the leftmost lit pixel within the given y range, or -1.
func firstLitX(img *image.RGBA, yMin, yMax int) int {
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := yMin; y <= yMax && y < b.Max.Y; y++ {
			if r, g, bl, _ := img.At(x, y).RGBA(); r > 0x2000 || g > 0x2000 || bl > 0x2000 {
				return x
			}
		}
	}
	return -1
}

// contractTestFonts registers two fonts with deliberately different row heights so
// that a renderer collapsing rows to a single height produces a detectably wrong
// layout. Returns a restore function.
func contractTestFonts() func() {
	restore := font.SnapshotAndClear()
	font.Register(contractFace{id: "tall", m: font.Metrics{GlyphWidth: 8, GlyphHeight: 16, GlyphAdvance: 8, RowHeight: 20}})
	font.Register(contractFace{id: "short", m: font.Metrics{GlyphWidth: 4, GlyphHeight: 6, GlyphAdvance: 4, RowHeight: 8}})
	return restore
}

type contractFace struct {
	id string
	m  font.Metrics
}

func (f contractFace) ID() string            { return f.id }
func (f contractFace) Metrics() font.Metrics { return f.m }

// GlyphRow lights every pixel of every glyph row so ink detection is unambiguous.
func (f contractFace) GlyphRow(_ rune, row int) uint32 {
	if row < 0 || row >= f.m.GlyphHeight {
		return 0
	}
	return 0xFFFFFFFF
}

// contractView is a two-row view whose rows use different fonts, with an explicit
// spacing and per-row horizontal offsets.
//
// Cursor is -1 so no row is drawn highlighted. That matters for these tests: the
// scrolling path fills the cursor row with a bright highlight rectangle, which the
// ink scan would otherwise report as text spanning the full row width.
func contractView(static bool) style.ViewData {
	return style.ViewData{
		Items:        []string{"AA", "BB"},
		FontIDs:      []string{"tall", "short"},
		RowHeights:   []int{20, 8},
		Spacing:      5,
		LineOffsets:  []int{7, 30},
		VisibleCount: 2,
		OffsetY:      10,
		Cursor:       -1,
		Static:       static,
	}
}

func contractBridge() layout.LayoutCalculator {
	return layout.NewLayoutBridge(textlayout.TextHints{
		PixelWidth:   120,
		PixelHeight:  100,
		GlyphAdvance: 6,
		RowHeight:    10,
	}, layout.BridgeConfig{PaddingPct: 0})
}

// TestLayoutContract_RowHeightsAndSpacingHonoured asserts that both draw paths
// advance by each row's own height plus the style's spacing.
//
// With rows of 20 and 8 pixels, spacing 5 and OffsetY 10, row 0 occupies y=10..29
// and row 1 begins at y=35. A renderer using one global row height would place row
// 1 at y=30 (or wherever its single height led), which this detects.
func TestLayoutContract_RowHeightsAndSpacingHonoured(t *testing.T) {
	restore := contractTestFonts()
	defer restore()

	for _, static := range []bool{true, false} {
		name := "scrolling"
		if static {
			name = "static"
		}
		t.Run(name, func(t *testing.T) {
			surf := surface.New(image.Rect(0, 0, 120, 100))
			surf.Clear(colBackground)
			rr := NewRegionRenderer(false, nil, nil)
			v := contractView(static)

			if static {
				rr.renderStaticItems(surf, v, contractBridge())
			} else {
				rr.renderItems(surf, v, contractBridge())
			}

			bands := inkRows(surf.FrameBuffer())
			if len(bands) != 2 {
				t.Fatalf("expected 2 ink bands (one per row), got %d: %v", len(bands), bands)
			}

			// Row 0 is the 16px-tall glyph inside a 20px row starting at y=10.
			if bands[0][0] != 12 {
				t.Errorf("row 0 ink starts at y=%d, want 12 (OffsetY 10 + centring within a 20px row)", bands[0][0])
			}
			// Row 1 must start after row 0's full height plus spacing: 10+20+5 = 35.
			// A single-global-row-height renderer would start it earlier.
			if bands[1][0] < 35 || bands[1][0] > 36 {
				t.Errorf("row 1 ink starts at y=%d, want 35-36 (row0 height 20 + spacing 5); "+
					"an earlier value means RowHeights or Spacing was ignored", bands[1][0])
			}
		})
	}
}

// TestLayoutContract_LineOffsetsHonoured asserts per-row horizontal offsets are
// applied on both paths.
//
// renderItems ignored LineOffsets entirely and drew every row at the content
// origin. Live styles depend on these offsets for centring: wifi's 240x240 view is
// Static:false and supplies its own per-row offsets.
func TestLayoutContract_LineOffsetsHonoured(t *testing.T) {
	restore := contractTestFonts()
	defer restore()

	for _, static := range []bool{true, false} {
		name := "scrolling"
		if static {
			name = "static"
		}
		t.Run(name, func(t *testing.T) {
			surf := surface.New(image.Rect(0, 0, 120, 100))
			surf.Clear(colBackground)
			rr := NewRegionRenderer(false, nil, nil)
			v := contractView(static)

			if static {
				rr.renderStaticItems(surf, v, contractBridge())
			} else {
				rr.renderItems(surf, v, contractBridge())
			}

			bands := inkRows(surf.FrameBuffer())
			if len(bands) != 2 {
				t.Fatalf("expected 2 ink bands, got %d: %v", len(bands), bands)
			}

			if x := firstLitX(surf.FrameBuffer(), bands[0][0], bands[0][1]); x != 7 {
				t.Errorf("row 0 ink starts at x=%d, want 7 (LineOffsets[0]); 0 means offsets were ignored", x)
			}
			if x := firstLitX(surf.FrameBuffer(), bands[1][0], bands[1][1]); x != 30 {
				t.Errorf("row 1 ink starts at x=%d, want 30 (LineOffsets[1]); 0 means offsets were ignored", x)
			}
		})
	}
}

// TestLayoutContract_VisibleCountHonoured asserts the style's row budget wins over
// the renderer's own estimate.
//
// The renderer's fallback divides region height by the current face's row height,
// which both depends on draw order and cannot be right for mixed tiers. A style
// that has already solved its fit reports VisibleCount, and that must be respected
// even when more rows would physically fit.
func TestLayoutContract_OffsetY_RespectsContentOrigin(t *testing.T) {
	restore := contractTestFonts()
	defer restore()

	for _, static := range []bool{true, false} {
		name := "scrolling"
		if static {
			name = "static"
		}
		t.Run(name, func(t *testing.T) {
			surf := surface.New(image.Rect(0, 0, 120, 100))
			surf.Clear(colBackground)
			rr := NewRegionRenderer(false, nil, nil)
			v := contractView(static)
			v.OffsetY = 10
			bridge := layout.NewLayoutBridge(textlayout.TextHints{
				PixelWidth:   120,
				PixelHeight:  100,
				GlyphAdvance: 6,
				RowHeight:    10,
			}, layout.BridgeConfig{PaddingPct: 10})

			if static {
				rr.renderStaticItems(surf, v, bridge)
			} else {
				rr.renderItems(surf, v, bridge)
			}

			bands := inkRows(surf.FrameBuffer())
			if len(bands) < 2 {
				t.Fatalf("expected at least 2 ink bands, got %d: %v", len(bands), bands)
			}
			if got, want := bands[0][0], 22; got != want {
				t.Fatalf("row 0 ink starts at y=%d, want %d (content origin Y=10 + OffsetY=10 + centring)", got, want)
			}
		})
	}
}

func TestLayoutContract_VisibleCountHonoured(t *testing.T) {
	restore := contractTestFonts()
	defer restore()

	for _, static := range []bool{true, false} {
		name := "scrolling"
		if static {
			name = "static"
		}
		t.Run(name, func(t *testing.T) {
			surf := surface.New(image.Rect(0, 0, 120, 100))
			surf.Clear(colBackground)
			rr := NewRegionRenderer(false, nil, nil)

			v := style.ViewData{
				Items:        []string{"AA", "BB", "CC"},
				FontIDs:      []string{"short", "short", "short"},
				RowHeights:   []int{8, 8, 8},
				Spacing:      4,
				VisibleCount: 2, // third row must not be drawn
				OffsetY:      10,
				Cursor:       -1, // no highlight rect; see contractView
				Static:       static,
			}

			if static {
				rr.renderStaticItems(surf, v, contractBridge())
			} else {
				rr.renderItems(surf, v, contractBridge())
			}

			if bands := inkRows(surf.FrameBuffer()); len(bands) != 2 {
				t.Fatalf("expected exactly 2 drawn rows (VisibleCount=2), got %d bands: %v", len(bands), bands)
			}
		})
	}
}
