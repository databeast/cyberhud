package source

import (
	"bytes"
	"image/color"
	"math"
	"testing"
)

func TestParseBaseColor(t *testing.T) {
	tests := []struct {
		input, normalized string
		want              color.RGBA
	}{
		{" CYAN ", "cyan", color.RGBA{0, 255, 255, 255}},
		{"emerald", "emerald", color.RGBA{0, 105, 55, 255}},
		{"#Ff6600", "#ff6600", color.RGBA{255, 102, 0, 255}},
		{"#000000", "#000000", color.RGBA{0, 0, 0, 255}},
	}
	for _, tt := range tests {
		normalized, got, err := ParseBaseColor(tt.input)
		if err != nil || normalized != tt.normalized || got != tt.want {
			t.Errorf("ParseBaseColor(%q) = %q, %v, %v", tt.input, normalized, got, err)
		}
	}
	for _, value := range []string{"", "none", "blue", "#fff", "#gg0000", "#12345678", "123456", "#-12345", "#+12345"} {
		if _, _, err := ParseBaseColor(value); err == nil {
			t.Errorf("ParseBaseColor(%q) accepted invalid color", value)
		}
	}
}

func TestCyanPreservesPrototypeColors(t *testing.T) {
	for _, rgb := range [][3]float64{
		{0.16, 0.86, 1}, {0.05, 0.75, 0.92}, {0.24, 0.91, 1}, {0.78, 1, 1},
	} {
		r, g, b := tinted(color.RGBA{0, 255, 255, 255}, rgb[0], rgb[1], rgb[2])
		for i, value := range []float64{r, g, b} {
			if math.Abs(value-rgb[i]) > 1e-12 {
				t.Fatalf("cyan changed %v into %v", rgb, []float64{r, g, b})
			}
		}
	}
}

func TestBaseColorChangesFrame(t *testing.T) {
	cyan := RenderFrame(240, 240, 2, color.RGBA{0, 255, 255, 255})
	red := RenderFrame(240, 240, 2, color.RGBA{255, 0, 0, 255})
	if bytes.Equal(cyan.Pix, red.Pix) {
		t.Fatal("changing the base color did not affect the image")
	}
	if cyan.RGBAAt(0, 0) != red.RGBAAt(0, 0) {
		t.Fatal("base color unexpectedly changed the background")
	}
}
