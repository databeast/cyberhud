package tests

import (
	"testing"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	_ "github.com/databeast/cyberhud/display/modes/attract_geometric"
	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
	_ "github.com/databeast/cyberhud/display/modes/attract_particles"
	_ "github.com/databeast/cyberhud/display/modes/attract_plasma"
	_ "github.com/databeast/cyberhud/display/modes/attract_shapes"
	_ "github.com/databeast/cyberhud/display/modes/attract_starfield"
	_ "github.com/databeast/cyberhud/display/modes/attract_waveform"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type styleCase struct {
	name       string
	modeID     string
	capability int
	wantStyle  string
}

func TestAttractModes_16x8StyleResolution(t *testing.T) {
	runUltraLowResStyleResolution(t, 16, 8, []styleCase{
		{name: "particles mono", modeID: "attract_particles", capability: textlayout.CapMonoFast, wantStyle: "mono-16x8"},
		{name: "particles grayscale", modeID: "attract_particles", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "particles color", modeID: "attract_particles", capability: textlayout.CapColorFast, wantStyle: "color-16x8"},
		{name: "plasma mono", modeID: "attract_plasma", capability: textlayout.CapMonoFast, wantStyle: "mono-16x8"},
		{name: "plasma grayscale", modeID: "attract_plasma", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "plasma color", modeID: "attract_plasma", capability: textlayout.CapColorFast, wantStyle: "color-16x8"},
		{name: "shapes mono", modeID: "attract_shapes", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-16x8"},
		{name: "shapes grayscale", modeID: "attract_shapes", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "shapes color", modeID: "attract_shapes", capability: textlayout.CapColorFast, wantStyle: "color-fast-16x8"},
		{name: "starfield mono", modeID: "attract_starfield", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-16x8"},
		{name: "starfield grayscale", modeID: "attract_starfield", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "starfield color", modeID: "attract_starfield", capability: textlayout.CapColorFast, wantStyle: "color-fast-16x8"},
		{name: "waveform mono", modeID: "attract_waveform", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-16x8"},
		{name: "waveform grayscale", modeID: "attract_waveform", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "waveform color", modeID: "attract_waveform", capability: textlayout.CapColorFast, wantStyle: "color-fast-16x8"},
		{name: "matrix mono", modeID: "attract_matrix", capability: textlayout.CapMonoFast, wantStyle: "mono-16x8"},
		{name: "matrix grayscale", modeID: "attract_matrix", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "matrix color", modeID: "attract_matrix", capability: textlayout.CapColorFast, wantStyle: "color-16x8"},
		{name: "geometric mono", modeID: "attract_geometric", capability: textlayout.CapMonoFast, wantStyle: "mono-16x8"},
		{name: "geometric grayscale", modeID: "attract_geometric", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-16x8"},
		{name: "geometric color", modeID: "attract_geometric", capability: textlayout.CapColorFast, wantStyle: "color-16x8"},
	})
}

func TestAttractModes_8x16StyleResolution(t *testing.T) {
	runUltraLowResStyleResolution(t, 8, 16, []styleCase{
		{name: "particles mono", modeID: "attract_particles", capability: textlayout.CapMonoFast, wantStyle: "mono-8x16"},
		{name: "particles grayscale", modeID: "attract_particles", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "particles color", modeID: "attract_particles", capability: textlayout.CapColorFast, wantStyle: "color-8x16"},
		{name: "plasma mono", modeID: "attract_plasma", capability: textlayout.CapMonoFast, wantStyle: "mono-8x16"},
		{name: "plasma grayscale", modeID: "attract_plasma", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "plasma color", modeID: "attract_plasma", capability: textlayout.CapColorFast, wantStyle: "color-8x16"},
		{name: "shapes mono", modeID: "attract_shapes", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-8x16"},
		{name: "shapes grayscale", modeID: "attract_shapes", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "shapes color", modeID: "attract_shapes", capability: textlayout.CapColorFast, wantStyle: "color-fast-8x16"},
		{name: "starfield mono", modeID: "attract_starfield", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-8x16"},
		{name: "starfield grayscale", modeID: "attract_starfield", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "starfield color", modeID: "attract_starfield", capability: textlayout.CapColorFast, wantStyle: "color-fast-8x16"},
		{name: "waveform mono", modeID: "attract_waveform", capability: textlayout.CapMonoFast, wantStyle: "mono-fast-8x16"},
		{name: "waveform grayscale", modeID: "attract_waveform", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "waveform color", modeID: "attract_waveform", capability: textlayout.CapColorFast, wantStyle: "color-fast-8x16"},
		{name: "matrix mono", modeID: "attract_matrix", capability: textlayout.CapMonoFast, wantStyle: "mono-8x16"},
		{name: "matrix grayscale", modeID: "attract_matrix", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "matrix color", modeID: "attract_matrix", capability: textlayout.CapColorFast, wantStyle: "color-8x16"},
		{name: "geometric mono", modeID: "attract_geometric", capability: textlayout.CapMonoFast, wantStyle: "mono-8x16"},
		{name: "geometric grayscale", modeID: "attract_geometric", capability: textlayout.CapGrayscaleFast, wantStyle: "grayscale-fast-8x16"},
		{name: "geometric color", modeID: "attract_geometric", capability: textlayout.CapColorFast, wantStyle: "color-8x16"},
	})
}

func runUltraLowResStyleResolution(t *testing.T, w, h int, cases []styleCase) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			hints := textlayout.TextHints{
				PixelWidth:   w,
				PixelHeight:  h,
				GlyphWidth:   5,
				GlyphHeight:  7,
				GlyphAdvance: 6,
				RowHeight:    10,
				Capability:   tc.capability,
			}
			modehints.PropagateHints(hints)

			inst, ok := displaymodes.GetInstance(tc.modeID)
			if !ok {
				t.Fatalf("GetInstance(%q) returned ok=false", tc.modeID)
			}

			got := inst.BuildView().StyleReport.Name
			if got != tc.wantStyle {
				t.Fatalf("%s resolved style %q, want %q", tc.modeID, got, tc.wantStyle)
			}
		})
	}
}
