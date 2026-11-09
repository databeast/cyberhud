package styles

import (
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestMatrixStylesRenderWithinPanel(t *testing.T) {
	tests := []struct {
		style def
		w, h  int
		cap   style.Capability
	}{
		{ColorFast13x9Style, 13, 9, style.ColorFast},
		{GrayscaleFast15x7Style, 15, 7, style.GrayscaleFast},
		{GrayscaleFast16x8Style, 16, 8, style.GrayscaleFast},
	}
	for _, tt := range tests {
		t.Run(tt.style.Name(), func(t *testing.T) {
			hints := textlayout.TextHints{PixelWidth: tt.w, PixelHeight: tt.h, Capability: int(tt.cap)}
			if got := tt.style.Supports(hints); got == style.Unsupported {
				t.Fatal("style does not support its panel")
			}
			vd := tt.style.Build(source.Snapshot{Lines: []string{"old", "\x1b[31mABC\x1b[0m"}}, source.Policy{}, style.NewStyleContext(hints))
			if len(vd.Sprites) != 1 || vd.Sprites[0].Label != "serial character" {
				t.Fatalf("expected a character sprite, got %#v", vd.Sprites)
			}
			sprite := vd.Sprites[0]
			bounds := sprite.Image.Bounds().Add(sprite.Position)
			if bounds.Min.X < 0 || bounds.Min.Y < 0 || bounds.Max.X > tt.w || bounds.Max.Y > tt.h {
				t.Fatalf("glyph bounds %v exceed %dx%d panel", bounds, tt.w, tt.h)
			}
			visible := 0
			for y := sprite.Image.Bounds().Min.Y; y < sprite.Image.Bounds().Max.Y; y++ {
				for x := sprite.Image.Bounds().Min.X; x < sprite.Image.Bounds().Max.X; x++ {
					_, _, _, a := sprite.Image.At(x, y).RGBA()
					if a != 0 {
						visible++
					}
				}
			}
			if visible == 0 {
				t.Fatal("character sprite is empty")
			}
		})
	}
}

func TestLastMatrixCharacter(t *testing.T) {
	tests := []struct {
		lines []string
		want  rune
		ok    bool
	}{
		{nil, 0, false},
		{[]string{"hello", "  "}, 'o', true},
		{[]string{"\x1b[31mAB\x1b[0m"}, 'B', true},
		{[]string{"hello\t"}, 'o', true},
		{[]string{"snowman \u2603"}, '?', true},
	}
	for _, tt := range tests {
		got, ok := lastMatrixCharacter(tt.lines)
		if got != tt.want || ok != tt.ok {
			t.Errorf("lastMatrixCharacter(%q) = %q, %t; want %q, %t", tt.lines, got, ok, tt.want, tt.ok)
		}
	}
}

func TestMatrixStatusAndRawMode(t *testing.T) {
	hints := textlayout.TextHints{PixelWidth: 13, PixelHeight: 9, Capability: int(style.ColorFast)}
	ctx := style.NewStyleContext(hints)
	for _, tt := range []struct {
		snap source.Snapshot
		want color.RGBA
	}{
		{source.Snapshot{Connected: true}, colorSerialOK},
		{source.Snapshot{}, colorSerialWarn},
		{source.Snapshot{LastError: "failed"}, colorSerialErr},
	} {
		vd := ColorFast13x9Style.Build(tt.snap, source.Policy{}, ctx)
		if len(vd.Sprites) != 1 {
			t.Fatalf("expected status sprite, got %#v", vd.Sprites)
		}
		got := color.RGBAModel.Convert(vd.Sprites[0].Image.At(0, 0)).(color.RGBA)
		if got != tt.want {
			t.Errorf("status color = %v, want %v", got, tt.want)
		}
	}
	vd := ColorFast13x9Style.Build(source.Snapshot{}, source.Policy{Raw: true}, ctx)
	if len(vd.Sprites) != 0 {
		t.Errorf("raw mode should hide status, got %#v", vd.Sprites)
	}
}
