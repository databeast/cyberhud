package waveshare_triple_screen

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
	"pgregory.net/rapid"
)

func TestProperty_OrientationFlipProducesTrue180Rotation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		def, err := panels.Get("waveshare-triple-screen")
		if err != nil {
			t.Fatalf("panels.Get failed: %v", err)
		}

		// Main screen is index 0
		screen := def.Virtual[0]
		config, ok := screen.Orientations[driver.OrientationFlip]
		if !ok {
			t.Fatalf("main screen has no OrientationFlip config")
		}

		// Property: OrientationFlip MADCTL must be MadctlMY (0x80) for hardware row-flip
		if config.MADCTL != driver.MadctlMY {
			t.Fatalf("main screen OrientationFlip: MADCTL = 0x%02X, want 0x%02X (MadctlMY)",
				config.MADCTL, driver.MadctlMY)
		}

		// Property: YOffset must be 80 to compensate for ST7789VW 320-row buffer
		if config.YOffset != 80 {
			t.Fatalf("main screen OrientationFlip: YOffset = %d, want 80", config.YOffset)
		}

		// Property: MirrorX must be true to correct column scan direction change
		if !config.MirrorX {
			t.Fatalf("main screen OrientationFlip: MirrorX = false, want true")
		}

		// Property: No software rotation (hardware handles orientation via MADCTL)
		if config.Rotation != 0 {
			t.Fatalf("main screen OrientationFlip: Rotation = %d, want 0", config.Rotation)
		}
	})
}

// screenOrientationCase holds a (screen-index, orientation) pair along with
// the expected OrientationConfig values observed on the UNFIXED code.
type screenOrientationCase struct {
	ScreenIndex int
	Orientation driver.Orientation
	Expected    driver.OrientationConfig
}

func TestProperty_PreservationNonFlipOrientations(t *testing.T) {
	// All (screen, orientation) pairs EXCEPT main+OrientationFlip.
	// Expected values observed on current (unfixed) code.
	cases := []screenOrientationCase{
		// Main screen (index 0)
		{ScreenIndex: 0, Orientation: driver.OrientationNormal, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMX}},
		{ScreenIndex: 0, Orientation: driver.OrientationCW, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMV | driver.MadctlMX}},
		{ScreenIndex: 0, Orientation: driver.OrientationCCW, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMV | driver.MadctlMY}},
		// Left-aux screen (index 1)
		{ScreenIndex: 1, Orientation: driver.OrientationNormal, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26}},
		{ScreenIndex: 1, Orientation: driver.OrientationFlip, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMX | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26}},
		// Right-aux screen (index 2)
		{ScreenIndex: 2, Orientation: driver.OrientationNormal, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26}},
		{ScreenIndex: 2, Orientation: driver.OrientationFlip, Expected: driver.OrientationConfig{MADCTL: driver.MadctlMX | driver.MadctlMV | driver.MadctlBGR, XOffset: 1, YOffset: 26}},
	}

	rapid.Check(t, func(t *rapid.T) {
		// Use rapid.SampledFrom to pick a random (screen, orientation) pair from the list
		c := rapid.SampledFrom(cases).Draw(t, "screenOrientationCase")

		def, err := panels.Get("waveshare-triple-screen")
		if err != nil {
			t.Fatalf("panels.Get failed: %v", err)
		}

		if c.ScreenIndex >= len(def.Virtual) {
			t.Fatalf("screen index %d out of range (have %d screens)", c.ScreenIndex, len(def.Virtual))
		}

		screen := def.Virtual[c.ScreenIndex]
		config, ok := screen.Orientations[c.Orientation]
		if !ok {
			t.Fatalf("screen[%d] %q has no orientation %q", c.ScreenIndex, screen.Name, c.Orientation)
		}

		// Assert MADCTL matches expected
		if config.MADCTL != c.Expected.MADCTL {
			t.Fatalf("screen[%d] %q orientation %q: MADCTL = 0x%02X, want 0x%02X",
				c.ScreenIndex, screen.Name, c.Orientation, config.MADCTL, c.Expected.MADCTL)
		}

		// Assert XOffset matches expected
		if config.XOffset != c.Expected.XOffset {
			t.Fatalf("screen[%d] %q orientation %q: XOffset = %d, want %d",
				c.ScreenIndex, screen.Name, c.Orientation, config.XOffset, c.Expected.XOffset)
		}

		// Assert YOffset matches expected
		if config.YOffset != c.Expected.YOffset {
			t.Fatalf("screen[%d] %q orientation %q: YOffset = %d, want %d",
				c.ScreenIndex, screen.Name, c.Orientation, config.YOffset, c.Expected.YOffset)
		}

		// Assert Rotation matches expected
		if config.Rotation != c.Expected.Rotation {
			t.Fatalf("screen[%d] %q orientation %q: Rotation = %d, want %d",
				c.ScreenIndex, screen.Name, c.Orientation, config.Rotation, c.Expected.Rotation)
		}
	})
}
