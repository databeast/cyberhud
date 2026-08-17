package tests

import (
	"bytes"
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/attract_matrix"
	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestCatalogRegistration(t *testing.T) {
	def, ok := catalog.Describe("attract_matrix")
	if !ok {
		t.Fatal("catalog.Describe(\"attract_matrix\") returned false")
	}
	if def.ID != "attract_matrix" {
		t.Errorf("ID = %q, want \"attract_matrix\"", def.ID)
	}
	if def.Title != "Matrix" {
		t.Errorf("Title = %q, want \"Matrix\"", def.Title)
	}
	if def.Summary == "" {
		t.Error("Summary is empty")
	}
	if len(def.Summary) > 120 {
		t.Errorf("Summary too long (%d chars, max 120)", len(def.Summary))
	}
	if def.Order != 200 {
		t.Errorf("Order = %d, want 200", def.Order)
	}
}

func TestCatalogCommand(t *testing.T) {
	cmd, ok := catalog.Command("attract_matrix")
	if !ok {
		t.Fatal("catalog.Command(\"attract_matrix\") returned false")
	}
	if cmd.Verb != "attract_matrix" {
		t.Errorf("Verb = %q, want \"attract_matrix\"", cmd.Verb)
	}
	if cmd.Summary == "" {
		t.Error("Summary is empty")
	}
	if cmd.Usage == "" {
		t.Error("Usage is empty")
	}
	if cmd.Handle == nil {
		t.Error("Handle is nil")
	}
}

func TestHandleCommand_Get(t *testing.T) {
	attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())
	result := attract_matrix.HandleCommand(nil)
	if !strings.HasPrefix(result, "OK attract_matrix") {
		t.Errorf("HandleCommand(nil) = %q, want prefix \"OK attract_matrix\"", result)
	}
	for _, key := range []string{"min_speed", "max_speed", "trail_length", "density", "show_background"} {
		if !strings.Contains(result, key+"=") {
			t.Errorf("response missing key %q: %s", key, result)
		}
	}
}

func TestDefaultPolicy(t *testing.T) {
	p := attract_matrix.DefaultPolicy()
	if p.MinSpeed != 1.5 {
		t.Errorf("MinSpeed = %v, want 1.5", p.MinSpeed)
	}
	if p.MaxSpeed != 6.0 {
		t.Errorf("MaxSpeed = %v, want 6.0", p.MaxSpeed)
	}
	if p.TrailLength != 16 {
		t.Errorf("TrailLength = %d, want 16", p.TrailLength)
	}
	if p.Density != 1.0 {
		t.Errorf("Density = %v, want 1.0", p.Density)
	}
	if p.ShowBackground != false {
		t.Errorf("ShowBackground = %v, want false", p.ShowBackground)
	}
}

func TestBuildView_UltraLowResSparsePixelRain(t *testing.T) {
	attract_matrix.ResetSnapshotState()
	attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

	hints := textlayout.TextHints{
		PixelWidth:   16,
		PixelHeight:  8,
		Capability:   textlayout.CapMonoFast,
		GlyphWidth:   5,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}

	vd := attract_matrix.BuildView(hints)
	if len(vd.Sprites) == 0 {
		t.Fatal("BuildView(16x8) returned no sprites")
	}
	if vd.Sprites[0].Label != "matrix-lowres" {
		t.Fatalf("sprite label = %q, want %q", vd.Sprites[0].Label, "matrix-lowres")
	}
	if vd.StyleReport.Name != "mono-16x8" {
		t.Fatalf("StyleReport.Name = %q, want %q", vd.StyleReport.Name, "mono-16x8")
	}

	img := vd.Sprites[0].Image
	if img == nil {
		t.Fatal("first sprite image is nil")
	}
	if gotW, gotH := img.Bounds().Dx(), img.Bounds().Dy(); gotW != 16 || gotH != 8 {
		t.Fatalf("sprite bounds = %dx%d, want 16x8", gotW, gotH)
	}

	lit := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			if c.R != 0 || c.G != 0 || c.B != 0 {
				lit++
			}
		}
	}

	if lit == 0 {
		t.Fatal("ultra-low-res renderer produced an all-black frame")
	}
	if lit >= 64 {
		t.Fatalf("ultra-low-res renderer is too dense: lit=%d, want <64", lit)
	}
}

func TestBuildView_UltraLowResPortraitPanelAnimates(t *testing.T) {
	attract_matrix.ResetSnapshotState()
	attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

	hints := textlayout.TextHints{
		PixelWidth:   8,
		PixelHeight:  16,
		Capability:   textlayout.CapMonoFast,
		GlyphWidth:   5,
		GlyphHeight:  7,
		GlyphAdvance: 6,
		RowHeight:    10,
	}

	first := attract_matrix.BuildView(hints)
	second := attract_matrix.BuildView(hints)

	if len(first.Sprites) == 0 || len(second.Sprites) == 0 {
		t.Fatal("expected sprites for ultra-low-res portrait panel")
	}
	if first.Sprites[0].Label != "matrix-lowres" || second.Sprites[0].Label != "matrix-lowres" {
		t.Fatalf("expected matrix-lowres sprite labels, got %q and %q", first.Sprites[0].Label, second.Sprites[0].Label)
	}

	firstRGBA, ok := first.Sprites[0].Image.(*image.RGBA)
	if !ok {
		t.Fatalf("expected first sprite image type *image.RGBA, got %T", first.Sprites[0].Image)
	}
	secondRGBA, ok := second.Sprites[0].Image.(*image.RGBA)
	if !ok {
		t.Fatalf("expected second sprite image type *image.RGBA, got %T", second.Sprites[0].Image)
	}
	if bytes.Equal(firstRGBA.Pix, secondRGBA.Pix) {
		t.Fatal("expected consecutive ultra-low-res portrait frames to animate")
	}
}

// headY returns the y of the unique full-brightness pixel in column x, or -1
// if not exactly one is found. Rain heads render at 255; trails are dimmer.
func headY(img image.Image, x int) int {
	found := -1
	for y := 0; y < img.Bounds().Dy(); y++ {
		c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
		if c.R == 255 && c.G == 255 && c.B == 255 {
			if found != -1 {
				return -1
			}
			found = y
		}
	}
	return found
}

func TestBuildView_UltraLowResRainFallsCoherently(t *testing.T) {
	for _, dims := range []struct{ w, h int }{{16, 8}, {8, 16}} {
		attract_matrix.ResetSnapshotState()
		attract_matrix.SetPolicy(attract_matrix.DefaultPolicy())

		hints := textlayout.TextHints{
			PixelWidth:   dims.w,
			PixelHeight:  dims.h,
			Capability:   textlayout.CapMonoFast,
			GlyphWidth:   5,
			GlyphHeight:  7,
			GlyphAdvance: 6,
			RowHeight:    10,
		}

		// Column 0 sits at fixed x = spacing/2 and its head advances one pixel
		// per frame: head = frameCounter % h (frameCounter starts at 1 here).
		const colX = 2
		first := attract_matrix.BuildView(hints)
		second := attract_matrix.BuildView(hints)

		y1 := headY(first.Sprites[0].Image, colX)
		y2 := headY(second.Sprites[0].Image, colX)
		if y1 != 1%dims.h || y2 != 2%dims.h {
			t.Fatalf("%dx%d: column head at y=%d then y=%d, want %d then %d (coherent downward motion)",
				dims.w, dims.h, y1, y2, 1%dims.h, 2%dims.h)
		}
	}
}

func TestResolveMatrixFont_SmallPanelUsesCompactMatrixFace(t *testing.T) {
	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64))
	ctx := style.NewStyleContext(hints)

	_, _, face := source.ResolveMatrixFont(ctx)
	if face == nil {
		t.Fatal("ResolveMatrixFont(128x64) returned nil face")
	}
	if got := face.ID(); got != "matrix-10x10" {
		t.Fatalf("ResolveMatrixFont(128x64) = %q, want %q", got, "matrix-10x10")
	}
}

func TestResolveMatrixFont_LargerPanelKeepsMatrixCodeFace(t *testing.T) {
	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240))
	ctx := style.NewStyleContext(hints)

	_, _, face := source.ResolveMatrixFont(ctx)
	if face == nil {
		t.Fatal("ResolveMatrixFont(240x240) returned nil face")
	}
	if got := face.ID(); got != "matrix-code" {
		t.Fatalf("ResolveMatrixFont(240x240) = %q, want %q", got, "matrix-code")
	}
}
