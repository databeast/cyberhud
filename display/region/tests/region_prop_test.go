package tests

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	region2 "github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// =============================================================================
// From: activation_property_test.go
// =============================================================================

// Property 9: Default layout produces correct single region

// For any valid set of screens with no explicit layout, ActivatePanel SHALL generate
// exactly one Region (index 0) whose bounds cover the full panel area, with the
// panel's DefaultMode assigned and input focus enabled.

// mockPropTarget is a minimal DrawTarget for property tests.
type mockPropTarget struct {
	bounds image.Rectangle
}

func (m *mockPropTarget) Bounds() image.Rectangle    { return m.bounds }
func (m *mockPropTarget) DrawImage(draw.Image) error { return nil }

func TestProperty9_DefaultLayoutProducesCorrectSingleRegion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random screen dimensions (single screen).
		w := rapid.IntRange(1, 800).Draw(t, "screenWidth")
		h := rapid.IntRange(1, 600).Draw(t, "screenHeight")

		screen := region2.ScreenPosition{
			Index:  0,
			Name:   "main",
			Bounds: image.Rect(0, 0, w, h),
			Target: &mockPropTarget{bounds: image.Rect(0, 0, w, h)},
		}

		// Generate a random set of available modes (1-5 modes).
		numModes := rapid.IntRange(1, 5).Draw(t, "numModes")
		availModes := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			availModes[i] = rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "mode")
		}

		// Pick one of the available modes as the DefaultMode.
		defaultIdx := rapid.IntRange(0, numModes-1).Draw(t, "defaultModeIdx")
		defaultMode := availModes[defaultIdx]

		// Generate random InputEnabled flag.
		inputEnabled := rapid.Bool().Draw(t, "inputEnabled")

		// Build PanelActivationConfig with Layout: nil (no explicit layout).
		config := region2.PanelActivationConfig{
			Screens:      []region2.ScreenPosition{screen},
			Layout:       nil,
			DefaultMode:  defaultMode,
			InputEnabled: inputEnabled,
			AvailModes:   availModes,
			ModeValidator: func(mode string) bool {
				for _, m := range availModes {
					if m == mode {
						return true
					}
				}
				return false
			},
		}

		activation, err := region2.ActivatePanel(config)
		if err != nil {
			t.Fatalf("ActivatePanel() unexpected error: %v", err)
		}

		// Property: exactly one region (index 0).
		regions := activation.RegionManager.Regions()
		if len(regions) != 1 {
			t.Fatalf("expected exactly 1 region, got %d", len(regions))
		}

		region := regions[0]

		// Property: region bounds cover the full panel area (union of all screen bounds).
		expectedBounds := image.Rect(0, 0, w, h)
		if region.Bounds() != expectedBounds {
			t.Fatalf("region bounds=%v, want %v (full panel area)", region.Bounds(), expectedBounds)
		}

		// Property: the panel's DefaultMode is assigned.
		if region.CurrentMode() != defaultMode {
			t.Fatalf("region mode=%q, want %q (panel's DefaultMode)", region.CurrentMode(), defaultMode)
		}

		// Property: input focus is enabled on the region.
		if !region.HasInputFocus() {
			t.Fatal("region should have input focus enabled")
		}
	})
}

func TestProperty9_DefaultLayoutMultiScreen_CoverFullArea(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 2-4 screens arranged horizontally with random widths.
		numScreens := rapid.IntRange(2, 4).Draw(t, "numScreens")
		h := rapid.IntRange(32, 300).Draw(t, "screenHeight")

		screens := make([]region2.ScreenPosition, numScreens)
		screenModes := make(map[string]string)
		x := 0

		// Generate available modes (at least numScreens modes to assign one per screen).
		numModes := rapid.IntRange(numScreens, numScreens+3).Draw(t, "numModes")
		availModes := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			availModes[i] = rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "mode")
		}

		for i := 0; i < numScreens; i++ {
			w := rapid.IntRange(32, 400).Draw(t, "screenWidth")
			name := fmt.Sprintf("%s_%d", rapid.StringMatching(`[a-z]{3,8}`).Draw(t, "screenName"), i)

			screens[i] = region2.ScreenPosition{
				Index:  i,
				Name:   name,
				Bounds: image.Rect(x, 0, x+w, h),
				Target: &mockPropTarget{bounds: image.Rect(0, 0, w, h)},
			}
			// Assign a mode for each screen from available modes.
			screenModes[name] = availModes[i%numModes]
			x += w
		}

		inputEnabled := rapid.Bool().Draw(t, "inputEnabled")

		config := region2.PanelActivationConfig{
			Screens:      screens,
			Layout:       nil,
			DefaultMode:  "",
			InputEnabled: inputEnabled,
			AvailModes:   availModes,
			ScreenModes:  screenModes,
			ModeValidator: func(mode string) bool {
				for _, m := range availModes {
					if m == mode {
						return true
					}
				}
				return false
			},
		}

		activation, err := region2.ActivatePanel(config)
		if err != nil {
			t.Fatalf("ActivatePanel() unexpected error: %v", err)
		}

		// For multi-screen: regions are created per-screen.
		// The union of all region bounds should cover the full VD area.
		regions := activation.RegionManager.Regions()
		if len(regions) != numScreens {
			t.Fatalf("expected %d regions (one per screen), got %d", numScreens, len(regions))
		}

		// Verify: union of all region bounds covers full VD bounds.
		vdBounds := activation.VirtualDisplay.Bounds()
		covered := image.Rectangle{}
		for i, r := range regions {
			if i == 0 {
				covered = r.Bounds()
			} else {
				covered = covered.Union(r.Bounds())
			}
		}
		if covered != vdBounds {
			t.Fatalf("union of region bounds=%v, want VD bounds=%v", covered, vdBounds)
		}

		// Verify: first region has input focus.
		if !regions[0].HasInputFocus() {
			t.Fatal("first region (index 0) should have input focus")
		}

		// Verify: each region has its screen's mode assigned.
		for i, r := range regions {
			expectedMode := screenModes[screens[i].Name]
			if r.CurrentMode() != expectedMode {
				t.Fatalf("region[%d] mode=%q, want %q", i, r.CurrentMode(), expectedMode)
			}
		}
	})
}

// =============================================================================
// From: flush_property_test.go
// =============================================================================

// captureMockTarget captures the image passed to DrawImage for property verification.
type captureMockTarget struct {
	bounds    image.Rectangle
	captured  *image.RGBA
	callCount int
}

func (m *captureMockTarget) Bounds() image.Rectangle {
	return m.bounds
}

func (m *captureMockTarget) DrawImage(img draw.Image) error {
	m.callCount++
	m.captured = img.(*image.RGBA)
	return nil
}

func TestProperty12_FlushPathExtractsCorrectSubRectangleAtZeroOrigin(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random VD dimensions.
		vdW := rapid.IntRange(10, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 500).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		// Fill the VD framebuffer with random pixel data (A=255 so compositing works).
		fb := vd.FrameBuffer()
		for y := 0; y < vdH; y++ {
			for x := 0; x < vdW; x++ {
				c := color.RGBA{
					R: rapid.Uint8().Draw(t, ""),
					G: rapid.Uint8().Draw(t, ""),
					B: rapid.Uint8().Draw(t, ""),
					A: 255,
				}
				fb.SetRGBA(x, y, c)
			}
		}

		// Generate a random ScreenPosition within the VD bounds.
		screenX0 := rapid.IntRange(0, vdW-1).Draw(t, "screenX0")
		screenY0 := rapid.IntRange(0, vdH-1).Draw(t, "screenY0")
		screenX1 := rapid.IntRange(screenX0+1, vdW).Draw(t, "screenX1")
		screenY1 := rapid.IntRange(screenY0+1, vdH).Draw(t, "screenY1")
		screenBounds := image.Rect(screenX0, screenY0, screenX1, screenY1)

		// Create a mock DrawTarget that captures the image.
		target := &captureMockTarget{bounds: image.Rect(0, 0, screenBounds.Dx(), screenBounds.Dy())}

		screens := []region2.ScreenPosition{
			{Index: 0, Name: "test", Bounds: screenBounds, Target: target},
		}

		fp := region2.NewFlushPath(vd, screens)
		if err := fp.Flush(); err != nil {
			t.Fatalf("Flush() error: %v", err)
		}

		if target.callCount != 1 {
			t.Fatalf("DrawImage called %d times, want 1", target.callCount)
		}

		// Verify: for every pixel (px, py) in the captured image, it equals the VD
		// framebuffer pixel at (Screen.Bounds.Min.X + px, Screen.Bounds.Min.Y + py).
		captured := target.captured
		capturedBounds := captured.Bounds()

		if capturedBounds.Min.X != 0 || capturedBounds.Min.Y != 0 {
			t.Fatalf("captured image not zero-origin: got %v", capturedBounds)
		}

		expectedW := screenBounds.Dx()
		expectedH := screenBounds.Dy()
		if capturedBounds.Dx() != expectedW || capturedBounds.Dy() != expectedH {
			t.Fatalf("captured dimensions = %dx%d, want %dx%d",
				capturedBounds.Dx(), capturedBounds.Dy(), expectedW, expectedH)
		}

		for py := 0; py < expectedH; py++ {
			for px := 0; px < expectedW; px++ {
				got := captured.RGBAAt(px, py)
				want := fb.RGBAAt(screenBounds.Min.X+px, screenBounds.Min.Y+py)
				if got != want {
					t.Fatalf("pixel (%d,%d): got %v, want %v (VD coord (%d,%d))",
						px, py, got, want,
						screenBounds.Min.X+px, screenBounds.Min.Y+py)
				}
			}
		}
	})
}

func TestProperty13_UncoveredPixelsInPhysicalScreenAreaAreBlack(t *testing.T) {
	// Sub-test 1: Fully uncovered screen (VD framebuffer at initial zero state).
	t.Run("fully_uncovered", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate random VD dimensions.
			vdW := rapid.IntRange(10, 300).Draw(t, "vdWidth")
			vdH := rapid.IntRange(10, 300).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("unexpected error creating VD: %v", err)
			}

			// Leave the VD framebuffer at its initial state (all zeros — transparent black).
			// The flush should fill uncovered areas with opaque black.

			// Generate a random ScreenPosition within the VD bounds.
			screenX0 := rapid.IntRange(0, vdW-1).Draw(t, "screenX0")
			screenY0 := rapid.IntRange(0, vdH-1).Draw(t, "screenY0")
			screenX1 := rapid.IntRange(screenX0+1, vdW).Draw(t, "screenX1")
			screenY1 := rapid.IntRange(screenY0+1, vdH).Draw(t, "screenY1")
			screenBounds := image.Rect(screenX0, screenY0, screenX1, screenY1)

			target := &captureMockTarget{bounds: image.Rect(0, 0, screenBounds.Dx(), screenBounds.Dy())}
			screens := []region2.ScreenPosition{
				{Index: 0, Name: "test", Bounds: screenBounds, Target: target},
			}

			fp := region2.NewFlushPath(vd, screens)
			if err := fp.Flush(); err != nil {
				t.Fatalf("Flush() error: %v", err)
			}

			// Verify all pixels in the captured image are RGBA(0, 0, 0, 255).
			captured := target.captured
			opaqueBlack := color.RGBA{0, 0, 0, 255}

			for py := 0; py < screenBounds.Dy(); py++ {
				for px := 0; px < screenBounds.Dx(); px++ {
					got := captured.RGBAAt(px, py)
					if got != opaqueBlack {
						t.Fatalf("pixel (%d,%d) = %v, want opaque black %v", px, py, got, opaqueBlack)
					}
				}
			}
		})
	})

	// Sub-test 2: Partially covered screen — filled pixels (A=255) show through,
	// unfilled pixels (A=0 in VD) come out as opaque black.
	t.Run("partially_covered", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate random VD dimensions.
			vdW := rapid.IntRange(10, 200).Draw(t, "vdWidth")
			vdH := rapid.IntRange(10, 200).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("unexpected error creating VD: %v", err)
			}

			// Screen covers the entire VD for simplicity.
			screenBounds := image.Rect(0, 0, vdW, vdH)

			// Generate a random "covered" sub-rectangle within the VD (simulating a Region).
			covX0 := rapid.IntRange(0, vdW-2).Draw(t, "covX0")
			covY0 := rapid.IntRange(0, vdH-2).Draw(t, "covY0")
			covX1 := rapid.IntRange(covX0+1, vdW).Draw(t, "covX1")
			covY1 := rapid.IntRange(covY0+1, vdH).Draw(t, "covY1")
			coveredRect := image.Rect(covX0, covY0, covX1, covY1)

			// Fill the covered rectangle with non-black, full-alpha pixels.
			fb := vd.FrameBuffer()
			coveredColor := color.RGBA{
				R: rapid.Uint8Range(1, 255).Draw(t, "covR"),
				G: rapid.Uint8Range(1, 255).Draw(t, "covG"),
				B: rapid.Uint8Range(1, 255).Draw(t, "covB"),
				A: 255,
			}
			for y := coveredRect.Min.Y; y < coveredRect.Max.Y; y++ {
				for x := coveredRect.Min.X; x < coveredRect.Max.X; x++ {
					fb.SetRGBA(x, y, coveredColor)
				}
			}

			target := &captureMockTarget{bounds: image.Rect(0, 0, vdW, vdH)}
			screens := []region2.ScreenPosition{
				{Index: 0, Name: "test", Bounds: screenBounds, Target: target},
			}

			fp := region2.NewFlushPath(vd, screens)
			if err := fp.Flush(); err != nil {
				t.Fatalf("Flush() error: %v", err)
			}

			captured := target.captured
			opaqueBlack := color.RGBA{0, 0, 0, 255}

			for py := 0; py < vdH; py++ {
				for px := 0; px < vdW; px++ {
					got := captured.RGBAAt(px, py)
					pt := image.Pt(px, py)
					if pt.In(coveredRect) {
						// Covered pixels should match the written color.
						if got != coveredColor {
							t.Fatalf("covered pixel (%d,%d) = %v, want %v", px, py, got, coveredColor)
						}
					} else {
						// Uncovered pixels should be opaque black.
						if got != opaqueBlack {
							t.Fatalf("uncovered pixel (%d,%d) = %v, want opaque black %v", px, py, got, opaqueBlack)
						}
					}
				}
			}
		})
	})
}

// errorMockTarget is a DrawTarget that optionally returns an error and tracks calls.
type errorMockTarget struct {
	bounds    image.Rectangle
	err       error
	callCount int
}

func (m *errorMockTarget) Bounds() image.Rectangle {
	return m.bounds
}

func (m *errorMockTarget) DrawImage(img draw.Image) error {
	m.callCount++
	return m.err
}

// Property 14: FlushPath error isolation

// For any set of M physical screens where J screens return hardware errors on
// DrawImage, the FlushPath SHALL still call DrawImage on all M screens and return
// a combined error containing all J individual errors.
func TestProperty14_FlushPathErrorIsolation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate M screens (1 to 8).
		m := rapid.IntRange(1, 8).Draw(t, "numScreens")

		// Determine per-screen width/height. Use uniform sizes for simplicity.
		screenW := rapid.IntRange(10, 200).Draw(t, "screenW")
		screenH := rapid.IntRange(10, 200).Draw(t, "screenH")

		// VD bounds: screens laid out horizontally.
		vdW := screenW * m
		vdH := screenH
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		// Fill the VD with some content so DrawImage has something to flush.
		fb := vd.FrameBuffer()
		fillColor := color.RGBA{R: 128, G: 64, B: 32, A: 255}
		for y := 0; y < vdH; y++ {
			for x := 0; x < vdW; x++ {
				fb.SetRGBA(x, y, fillColor)
			}
		}

		// For each screen, randomly decide if it should fail.
		targets := make([]*errorMockTarget, m)
		screens := make([]region2.ScreenPosition, m)
		var expectedErrors []string

		for i := 0; i < m; i++ {
			shouldFail := rapid.Bool().Draw(t, fmt.Sprintf("screen%d_fails", i))
			screenBounds := image.Rect(i*screenW, 0, (i+1)*screenW, screenH)

			var screenErr error
			if shouldFail {
				errMsg := fmt.Sprintf("hw error on screen %d", i)
				screenErr = errors.New(errMsg)
				expectedErrors = append(expectedErrors, errMsg)
			}

			targets[i] = &errorMockTarget{
				bounds: image.Rect(0, 0, screenW, screenH),
				err:    screenErr,
			}
			screens[i] = region2.ScreenPosition{
				Index:  i,
				Name:   fmt.Sprintf("screen-%d", i),
				Bounds: screenBounds,
				Target: targets[i],
			}
		}

		fp := region2.NewFlushPath(vd, screens)
		flushErr := fp.Flush()

		// Verify: ALL M screens had DrawImage called exactly once.
		for i, target := range targets {
			if target.callCount != 1 {
				t.Fatalf("screen %d: DrawImage called %d times, want 1", i, target.callCount)
			}
		}

		j := len(expectedErrors)

		if j == 0 {
			// No errors expected — Flush should return nil.
			if flushErr != nil {
				t.Fatalf("expected nil error when no screens fail, got: %v", flushErr)
			}
		} else {
			// J screens failed — combined error must be non-nil.
			if flushErr == nil {
				t.Fatalf("expected combined error for %d failing screens, got nil", j)
			}

			// Verify the combined error contains all J individual error messages.
			combinedMsg := flushErr.Error()
			for _, expectedMsg := range expectedErrors {
				if !strings.Contains(combinedMsg, expectedMsg) {
					t.Fatalf("combined error missing expected message %q; got: %v", expectedMsg, combinedMsg)
				}
			}
		}
	})
}

// =============================================================================
// From: layout_property_test.go
// =============================================================================

// TestProperty6_SingleScreenDefaultRegionCoversFullVD tests that for any single-screen
// Panel with no explicit Region_Layout, GenerateDefaultLayout creates exactly one Region
// named "default" whose bounds equal the VirtualDisplay bounds.
func TestProperty6_SingleScreenDefaultRegionCoversFullVD(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random single-screen dimensions.
		w := rapid.IntRange(1, 500).Draw(t, "screenWidth")
		h := rapid.IntRange(1, 500).Draw(t, "screenHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, w, h))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		// Generate a random set of available modes (1-5 modes).
		numModes := rapid.IntRange(1, 5).Draw(t, "numModes")
		availModes := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			availModes[i] = rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "mode")
		}

		// Generate random config: DefaultMode sometimes set, InputEnabled random.
		defaultMode := ""
		hasExplicitDefault := rapid.Bool().Draw(t, "hasExplicitDefault")
		if hasExplicitDefault {
			// Pick a mode from availModes as the explicit default.
			idx := rapid.IntRange(0, numModes-1).Draw(t, "defaultModeIdx")
			defaultMode = availModes[idx]
		}

		inputEnabled := rapid.Bool().Draw(t, "inputEnabled")

		config := region2.PanelActivationConfig{
			Screens: []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, w, h)},
			},
			DefaultMode:  defaultMode,
			InputEnabled: inputEnabled,
			AvailModes:   availModes,
		}

		layout, err := region2.GenerateDefaultLayout(vd, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify: exactly 1 spec.
		if len(layout.Specs) != 1 {
			t.Fatalf("expected exactly 1 spec, got %d", len(layout.Specs))
		}

		spec := layout.Specs[0]

		// Verify: name is "default".
		if spec.Name != "default" {
			t.Fatalf("expected spec name %q, got %q", "default", spec.Name)
		}

		// Verify: bounds equal VD bounds.
		if spec.Bounds != vd.Bounds() {
			t.Fatalf("expected spec bounds %v (VD bounds), got %v", vd.Bounds(), spec.Bounds)
		}
	})
}

// TestProperty6_SingleScreenModeResolutionOrder tests the mode resolution order:
// Panel.DefaultMode if registered → "menu" if input enabled → "dashboard" → first available.
func TestProperty6_SingleScreenModeResolutionOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random single-screen dimensions.
		w := rapid.IntRange(1, 500).Draw(t, "screenWidth")
		h := rapid.IntRange(1, 500).Draw(t, "screenHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, w, h))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		// Generate a random set of available modes (1-5 modes).
		numModes := rapid.IntRange(1, 5).Draw(t, "numModes")
		availModes := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			// Use unique mode names that avoid "menu" and "dashboard" to control the test.
			availModes[i] = rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "mode")
		}

		// Decide whether to include "menu" and "dashboard" in the mode list.
		includeMenu := rapid.Bool().Draw(t, "includeMenu")
		includeDashboard := rapid.Bool().Draw(t, "includeDashboard")

		if includeMenu {
			availModes = append(availModes, "menu")
		}
		if includeDashboard {
			availModes = append(availModes, "dashboard")
		}

		// Decide whether to set an explicit DefaultMode.
		hasExplicitDefault := rapid.Bool().Draw(t, "hasExplicitDefault")
		defaultMode := ""
		if hasExplicitDefault && len(availModes) > 0 {
			idx := rapid.IntRange(0, len(availModes)-1).Draw(t, "defaultModeIdx")
			defaultMode = availModes[idx]
		}

		inputEnabled := rapid.Bool().Draw(t, "inputEnabled")

		config := region2.PanelActivationConfig{
			Screens: []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, w, h)},
			},
			DefaultMode:  defaultMode,
			InputEnabled: inputEnabled,
			AvailModes:   availModes,
		}

		layout, err := region2.GenerateDefaultLayout(vd, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		spec := layout.Specs[0]

		// Compute expected mode using the resolution order.
		expectedMode := computeExpectedMode(config)

		if spec.DefaultMode != expectedMode {
			t.Fatalf("mode resolution mismatch: got %q, want %q\n"+
				"config: DefaultMode=%q, InputEnabled=%v, AvailModes=%v",
				spec.DefaultMode, expectedMode,
				config.DefaultMode, config.InputEnabled, config.AvailModes)
		}
	})
}

// TestProperty6_DefaultModeNotInAvailModes tests that when Panel.DefaultMode is set
// but NOT in AvailModes, the resolution falls through to the next option.
func TestProperty6_DefaultModeNotInAvailModes(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 500).Draw(t, "screenWidth")
		h := rapid.IntRange(1, 500).Draw(t, "screenHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, w, h))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		// Generate available modes that do NOT include the default mode.
		numModes := rapid.IntRange(1, 5).Draw(t, "numModes")
		availModes := make([]string, numModes)
		for i := 0; i < numModes; i++ {
			availModes[i] = rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "mode")
		}

		// Set a default mode that is NOT in availModes.
		defaultMode := "nonexistent_mode_xyz"

		inputEnabled := rapid.Bool().Draw(t, "inputEnabled")

		config := region2.PanelActivationConfig{
			Screens: []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, w, h)},
			},
			DefaultMode:  defaultMode,
			InputEnabled: inputEnabled,
			AvailModes:   availModes,
		}

		layout, err := region2.GenerateDefaultLayout(vd, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		spec := layout.Specs[0]

		// The DefaultMode is not in AvailModes, so it should fall through.
		expectedMode := computeExpectedMode(config)

		if spec.DefaultMode != expectedMode {
			t.Fatalf("mode resolution mismatch when DefaultMode not in AvailModes: got %q, want %q\n"+
				"config: DefaultMode=%q, InputEnabled=%v, AvailModes=%v",
				spec.DefaultMode, expectedMode,
				config.DefaultMode, config.InputEnabled, config.AvailModes)
		}

		// The mode should NOT be the unregistered default.
		if spec.DefaultMode == defaultMode {
			t.Fatalf("resolved mode should not be the unregistered DefaultMode %q", defaultMode)
		}
	})
}

// TestProperty6_EmptyAvailModesReturnsError tests that when AvailModes is empty,
// GenerateDefaultLayout returns an error.
func TestProperty6_EmptyAvailModesReturnsError(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 500).Draw(t, "screenWidth")
		h := rapid.IntRange(1, 500).Draw(t, "screenHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, w, h))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		config := region2.PanelActivationConfig{
			Screens: []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, w, h)},
			},
			DefaultMode:  rapid.String().Draw(t, "defaultMode"),
			InputEnabled: rapid.Bool().Draw(t, "inputEnabled"),
			AvailModes:   []string{},
		}

		_, err = region2.GenerateDefaultLayout(vd, config)
		if err == nil {
			t.Fatal("expected error when AvailModes is empty, got nil")
		}
	})
}

// modeInList mirrors the region package's unexported helper for test-side
// expected-mode computation.
func modeInList(mode string, list []string) bool {
	for _, m := range list {
		if m == mode {
			return true
		}
	}
	return false
}

// computeExpectedMode implements the mode resolution order as specified:
// 1. config.DefaultMode if non-empty AND in config.AvailModes
// 2. "menu" if config.InputEnabled AND "menu" in config.AvailModes
// 3. "dashboard" if in config.AvailModes
// 4. First entry in config.AvailModes
func computeExpectedMode(config region2.PanelActivationConfig) string {
	if config.DefaultMode != "" && modeInList(config.DefaultMode, config.AvailModes) {
		return config.DefaultMode
	}
	if config.InputEnabled && modeInList("menu", config.AvailModes) {
		return "menu"
	}
	if modeInList("dashboard", config.AvailModes) {
		return "dashboard"
	}
	if len(config.AvailModes) > 0 {
		return config.AvailModes[0]
	}
	return ""
}

// =============================================================================
// From: manager_property_test.go
// =============================================================================

// lowercaseRunes is used by generators to produce lowercase alphabetic strings.
var lowercaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// TestProperty4_ValidAllocationSucceeds tests that allocation succeeds when all
// constraints are satisfied: unique name (1-64 chars), bounds with width >= 1 and
// height >= 1, bounds within VD, and no overlap with existing regions.
func TestProperty4_ValidAllocationSucceeds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a VirtualDisplay with reasonable dimensions.
		vdW := rapid.IntRange(10, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil) // accept all modes for this test

		// Generate a sequence of non-overlapping regions within VD bounds.
		numRegions := rapid.IntRange(1, 5).Draw(t, "numRegions")
		allocatedNames := make(map[string]bool)
		allocatedBounds := make([]image.Rectangle, 0)

		for i := 0; i < numRegions; i++ {
			// Generate a unique name (1-64 chars).
			name := genUniqueName(t, allocatedNames, i)
			allocatedNames[strings.ToLower(name)] = true

			// Generate bounds that fit within VD and don't overlap existing regions.
			bounds := genNonOverlappingBounds(t, vdW, vdH, allocatedBounds, i)
			if bounds.Empty() {
				// Could not find non-overlapping bounds; skip this iteration.
				continue
			}
			allocatedBounds = append(allocatedBounds, bounds)

			spec := region2.RegionSpec{
				Name:   name,
				Bounds: bounds,
			}

			err := rm.Allocate(spec)
			if err != nil {
				t.Fatalf("expected allocation to succeed for valid spec %+v, got error: %v", spec, err)
			}

			// Verify the region was allocated.
			r, ok := rm.RegionByName(name)
			if !ok {
				t.Fatalf("region %q not found after allocation", name)
			}
			if r.Bounds() != bounds {
				t.Fatalf("region %q bounds mismatch: got %v, want %v", name, r.Bounds(), bounds)
			}
		}
	})
}

// TestProperty4_EmptyNameFails tests that allocation fails when name is empty.
func TestProperty4_EmptyNameFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(10, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		spec := region2.RegionSpec{
			Name:   "",
			Bounds: image.Rect(0, 0, 1, 1),
		}

		err = rm.Allocate(spec)
		if err == nil {
			t.Fatal("expected allocation to fail for empty name")
		}
		if !strings.Contains(err.Error(), "empty string") {
			t.Fatalf("error should mention empty string, got: %v", err)
		}
	})
}

// TestProperty4_NameTooLongFails tests that allocation fails when name exceeds 64 characters.
func TestProperty4_NameTooLongFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(10, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Generate a name longer than 64 characters.
		nameLen := rapid.IntRange(65, 200).Draw(t, "nameLen")
		name := strings.Repeat("a", nameLen)

		spec := region2.RegionSpec{
			Name:   name,
			Bounds: image.Rect(0, 0, 1, 1),
		}

		err = rm.Allocate(spec)
		if err == nil {
			t.Fatalf("expected allocation to fail for name of length %d", nameLen)
		}
		if !strings.Contains(err.Error(), "exceeds 64 character limit") {
			t.Fatalf("error should mention 64 character limit, got: %v", err)
		}
	})
}

// TestProperty4_DuplicateNameFails tests that allocation fails when a duplicate
// name is used (case-insensitive).
func TestProperty4_DuplicateNameFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(20, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(20, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Allocate the first region.
		name := genValidName(t, "firstName")
		spec1 := region2.RegionSpec{
			Name:   name,
			Bounds: image.Rect(0, 0, vdW/2, vdH/2),
		}
		err = rm.Allocate(spec1)
		if err != nil {
			t.Fatal(err)
		}

		// Try to allocate with same name (different case).
		dupName := toggleCase(name)
		spec2 := region2.RegionSpec{
			Name:   dupName,
			Bounds: image.Rect(vdW/2, 0, vdW, vdH/2),
		}
		err = rm.Allocate(spec2)
		if err == nil {
			t.Fatalf("expected allocation to fail for duplicate name %q (original: %q)", dupName, name)
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("error should mention 'already exists', got: %v", err)
		}
	})
}

// TestProperty4_ZeroBoundsFails tests that allocation fails when bounds have
// zero width or height.
func TestProperty4_ZeroBoundsFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(10, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Choose whether to violate width or height or both.
		violation := rapid.IntRange(0, 2).Draw(t, "violation")
		var bounds image.Rectangle
		switch violation {
		case 0: // zero width
			bounds = image.Rect(0, 0, 0, rapid.IntRange(1, vdH).Draw(t, "h"))
		case 1: // zero height
			bounds = image.Rect(0, 0, rapid.IntRange(1, vdW).Draw(t, "w"), 0)
		case 2: // both zero
			bounds = image.Rect(0, 0, 0, 0)
		}

		spec := region2.RegionSpec{
			Name:   "test",
			Bounds: bounds,
		}

		err = rm.Allocate(spec)
		if err == nil {
			t.Fatalf("expected allocation to fail for zero bounds %v", bounds)
		}
		if !strings.Contains(err.Error(), "width >= 1 and height >= 1") {
			t.Fatalf("error should mention bounds constraint, got: %v", err)
		}
	})
}

// TestProperty4_OutOfBoundsFails tests that allocation fails when the region
// bounds extend outside the VirtualDisplay.
func TestProperty4_OutOfBoundsFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(10, 200).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 200).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Generate bounds that extend outside VD in at least one direction.
		bounds := genOutOfBounds(t, vdW, vdH)

		spec := region2.RegionSpec{
			Name:   "oob",
			Bounds: bounds,
		}

		err = rm.Allocate(spec)
		if err == nil {
			t.Fatalf("expected allocation to fail for out-of-bounds %v (VD: %dx%d)", bounds, vdW, vdH)
		}
		if !strings.Contains(err.Error(), "extend outside virtual display") {
			t.Fatalf("error should mention 'extend outside virtual display', got: %v", err)
		}
		// Error should identify the region by name.
		if !strings.Contains(err.Error(), "oob") {
			t.Fatalf("error should identify region by name 'oob', got: %v", err)
		}
	})
}

// TestProperty4_OverlapFails tests that allocation fails when the region bounds
// overlap an existing region.
func TestProperty4_OverlapFails(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(20, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(20, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Allocate a region that takes a portion of the VD.
		w1 := rapid.IntRange(2, vdW-1).Draw(t, "w1")
		h1 := rapid.IntRange(2, vdH-1).Draw(t, "h1")
		spec1 := region2.RegionSpec{
			Name:   "first",
			Bounds: image.Rect(0, 0, w1, h1),
		}
		err = rm.Allocate(spec1)
		if err != nil {
			t.Fatal(err)
		}

		// Generate overlapping bounds.
		x2 := rapid.IntRange(0, w1-1).Draw(t, "x2")
		y2 := rapid.IntRange(0, h1-1).Draw(t, "y2")
		x2End := rapid.IntRange(x2+1, vdW).Draw(t, "x2End")
		y2End := rapid.IntRange(y2+1, vdH).Draw(t, "y2End")

		spec2 := region2.RegionSpec{
			Name:   "second",
			Bounds: image.Rect(x2, y2, x2End, y2End),
		}

		err = rm.Allocate(spec2)
		if err == nil {
			t.Fatalf("expected allocation to fail for overlapping bounds %v (existing: %v)", spec2.Bounds, spec1.Bounds)
		}
		if !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("error should mention 'overlap', got: %v", err)
		}
		if !strings.Contains(err.Error(), "second") {
			t.Fatalf("error should identify region by name 'second', got: %v", err)
		}
	})
}

// TestProperty4_SequentialAllocations tests a sequence of allocations where
// the first few are valid and then an invalid one is attempted.
func TestProperty4_SequentialAllocations(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(40, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(40, 500).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		halfW := vdW / 2
		halfH := vdH / 2
		spec1 := region2.RegionSpec{
			Name:   "topleft",
			Bounds: image.Rect(0, 0, halfW, halfH),
		}
		if err := rm.Allocate(spec1); err != nil {
			t.Fatal(err)
		}

		spec2 := region2.RegionSpec{
			Name:   "topright",
			Bounds: image.Rect(halfW, 0, vdW, halfH),
		}
		if err := rm.Allocate(spec2); err != nil {
			t.Fatal(err)
		}

		// Attempt an invalid allocation (overlapping with topleft).
		specInvalid := region2.RegionSpec{
			Name:   "invalid",
			Bounds: image.Rect(0, 0, halfW, halfH),
		}
		err = rm.Allocate(specInvalid)
		if err == nil {
			t.Fatal("expected invalid allocation to fail")
		}

		// Verify valid regions are still accessible.
		r1, ok := rm.RegionByName("topleft")
		if !ok {
			t.Fatal("topleft region should still exist")
		}
		if r1.Bounds() != spec1.Bounds {
			t.Fatalf("topleft bounds mismatch: got %v, want %v", r1.Bounds(), spec1.Bounds)
		}

		r2, ok := rm.RegionByName("topright")
		if !ok {
			t.Fatal("topright region should still exist")
		}
		if r2.Bounds() != spec2.Bounds {
			t.Fatalf("topright bounds mismatch: got %v, want %v", r2.Bounds(), spec2.Bounds)
		}

		_, ok = rm.RegionByName("invalid")
		if ok {
			t.Fatal("invalid region should not exist after failed allocation")
		}

		if len(rm.Regions()) != 2 {
			t.Fatalf("expected 2 regions, got %d", len(rm.Regions()))
		}
	})
}

// TestProperty4_AllocationSucceedsIffConstraintsSatisfied is the comprehensive
// property test that verifies allocation succeeds if and only if ALL constraints
// are simultaneously satisfied.
func TestProperty4_AllocationSucceedsIffConstraintsSatisfied(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(10, 300).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 300).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		numExisting := rapid.IntRange(0, 3).Draw(t, "numExisting")
		existingBounds := make([]image.Rectangle, 0, numExisting)
		existingNames := make(map[string]bool)

		for i := 0; i < numExisting; i++ {
			name := genUniqueName(t, existingNames, i)
			bounds := genNonOverlappingBounds(t, vdW, vdH, existingBounds, i)
			if bounds.Empty() {
				break
			}
			existingNames[strings.ToLower(name)] = true
			existingBounds = append(existingBounds, bounds)
			if err := rm.Allocate(region2.RegionSpec{Name: name, Bounds: bounds}); err != nil {
				t.Fatal(err)
			}
		}

		// Generate a random spec that may or may not satisfy constraints.
		name := genArbitraryName(t)
		bounds := genArbitraryBounds(t, vdW, vdH)

		spec := region2.RegionSpec{
			Name:   name,
			Bounds: bounds,
		}

		nameValid := len(name) >= 1 && len(name) <= 64
		nameUnique := !existingNames[strings.ToLower(name)]
		boundsPositive := bounds.Dx() >= 1 && bounds.Dy() >= 1
		boundsWithinVD := bounds.In(image.Rect(0, 0, vdW, vdH))
		noOverlap := true
		for _, eb := range existingBounds {
			if !bounds.Intersect(eb).Empty() {
				noOverlap = false
				break
			}
		}

		allConstraints := nameValid && nameUnique && boundsPositive && boundsWithinVD && noOverlap

		err = rm.Allocate(spec)

		if allConstraints {
			if err != nil {
				t.Fatalf("expected allocation to succeed (all constraints satisfied) for spec %+v, got error: %v", spec, err)
			}
		} else {
			if err == nil {
				t.Fatalf("expected allocation to fail (constraint violated) for spec %+v. "+
					"nameValid=%v, nameUnique=%v, boundsPositive=%v, boundsWithinVD=%v, noOverlap=%v",
					spec, nameValid, nameUnique, boundsPositive, boundsWithinVD, noOverlap)
			}
		}
	})
}

// TestProperty7_LayoutAllocationIsAtomic tests that if any Region in a
// RegionLayout fails validation, no Regions are allocated.
func TestProperty7_LayoutAllocationIsAtomic(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(20, 300).Draw(t, "vdWidth")
		vdH := rapid.IntRange(20, 300).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(nil)

		// Generate a layout with at least one valid spec and one invalid spec.
		numValid := rapid.IntRange(1, 3).Draw(t, "numValid")
		validSpecs := make([]region2.RegionSpec, 0, numValid)
		usedBounds := make([]image.Rectangle, 0)

		halfW := vdW / (numValid + 1)
		for i := 0; i < numValid; i++ {
			bounds := image.Rect(i*halfW, 0, (i+1)*halfW, vdH)
			if !bounds.In(image.Rect(0, 0, vdW, vdH)) || bounds.Dx() < 1 {
				continue
			}
			validSpecs = append(validSpecs, region2.RegionSpec{
				Name:   fmt.Sprintf("region%d", i),
				Bounds: bounds,
			})
			usedBounds = append(usedBounds, bounds)
		}

		if len(validSpecs) == 0 {
			return // skip if we couldn't generate valid specs
		}

		// Generate an invalid spec (overlaps with first valid).
		invalidSpec := region2.RegionSpec{
			Name:   "invalid",
			Bounds: validSpecs[0].Bounds, // same bounds = overlap
		}

		// Place the invalid spec at a random position in the layout.
		position := rapid.IntRange(0, len(validSpecs)).Draw(t, "invalidPos")
		allSpecs := make([]region2.RegionSpec, 0, len(validSpecs)+1)
		for i, s := range validSpecs {
			if i == position {
				allSpecs = append(allSpecs, invalidSpec)
			}
			allSpecs = append(allSpecs, s)
		}
		if position == len(validSpecs) {
			allSpecs = append(allSpecs, invalidSpec)
		}

		layout := region2.RegionLayout{Specs: allSpecs}

		err = rm.AllocateLayout(layout)
		if err == nil {
			t.Fatal("expected layout allocation to fail due to invalid spec")
		}

		// Verify NO regions were allocated (atomicity).
		if len(rm.Regions()) != 0 {
			t.Fatalf("expected 0 regions after failed layout, got %d", len(rm.Regions()))
		}
	})
}

func TestProperty10_ModeChangeClearsRegionSurfaceToBlack(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random VD and region bounds.
		vdW := rapid.IntRange(2, 200).Draw(t, "vdWidth")
		vdH := rapid.IntRange(2, 200).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(func(mode string) bool {
			return mode == "modeA" || mode == "modeB"
		})

		// Generate region bounds within VD.
		x0 := rapid.IntRange(0, vdW-2).Draw(t, "x0")
		y0 := rapid.IntRange(0, vdH-2).Draw(t, "y0")
		x1 := rapid.IntRange(x0+1, vdW).Draw(t, "x1")
		y1 := rapid.IntRange(y0+1, vdH).Draw(t, "y1")
		bounds := image.Rect(x0, y0, x1, y1)

		err = rm.Allocate(region2.RegionSpec{
			Name:        "test",
			Bounds:      bounds,
			DefaultMode: "modeA",
		})
		if err != nil {
			t.Fatalf("Allocate: %v", err)
		}

		r, _ := rm.Region(0)

		// Fill the region's surface with random non-black pixels.
		fb := r.Surface().FrameBuffer()
		w := bounds.Dx()
		h := bounds.Dy()
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				c := color.RGBA{
					R: rapid.Uint8Range(1, 255).Draw(t, "r"),
					G: rapid.Uint8Range(1, 255).Draw(t, "g"),
					B: rapid.Uint8Range(1, 255).Draw(t, "b"),
					A: 255,
				}
				fb.SetRGBA(px, py, c)
			}
		}

		// Verify we actually have non-black pixels.
		sample := fb.RGBAAt(0, 0)
		if sample.R == 0 && sample.G == 0 && sample.B == 0 {
			t.Fatal("setup: expected non-black pixels")
		}

		// Wire ModeFactory so SetMode can construct instances.
		r.SetModeFactory(noopModeFactory())

		// Call SetMode with a valid mode.
		err = rm.SetMode("0", "modeB")
		if err != nil {
			t.Fatalf("SetMode: %v", err)
		}

		// Verify ALL pixels in the Region's surface are RGBA(0, 0, 0, 255).
		black := color.RGBA{0, 0, 0, 255}
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				got := fb.RGBAAt(px, py)
				if got != black {
					t.Fatalf("pixel(%d,%d) = %v after SetMode, want %v", px, py, got, black)
				}
			}
		}
	})
}

func TestProperty11_RegionModeEnginesAreIndependent(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Create a VD large enough for two side-by-side regions.
		regionW := rapid.IntRange(2, 100).Draw(t, "regionWidth")
		regionH := rapid.IntRange(2, 100).Draw(t, "regionHeight")
		vdW := regionW * 2

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, regionH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(func(mode string) bool {
			return mode == "modeA" || mode == "modeB" || mode == "modeC"
		})

		// Allocate two side-by-side regions.
		err = rm.Allocate(region2.RegionSpec{
			Name:        "left",
			Bounds:      image.Rect(0, 0, regionW, regionH),
			DefaultMode: "modeA",
		})
		if err != nil {
			t.Fatalf("Allocate left: %v", err)
		}

		err = rm.Allocate(region2.RegionSpec{
			Name:        "right",
			Bounds:      image.Rect(regionW, 0, vdW, regionH),
			DefaultMode: "modeB",
		})
		if err != nil {
			t.Fatalf("Allocate right: %v", err)
		}

		leftRegion, _ := rm.Region(0)
		rightRegion, _ := rm.Region(1)

		// Wire ModeFactory so SetMode can construct instances.
		wireNoopFactory(rm)

		// Fill both regions with distinct colors.
		leftColor := color.RGBA{
			R: rapid.Uint8Range(1, 255).Draw(t, "leftR"),
			G: rapid.Uint8Range(1, 255).Draw(t, "leftG"),
			B: rapid.Uint8Range(1, 255).Draw(t, "leftB"),
			A: 255,
		}
		rightColor := color.RGBA{
			R: rapid.Uint8Range(1, 255).Draw(t, "rightR"),
			G: rapid.Uint8Range(1, 255).Draw(t, "rightG"),
			B: rapid.Uint8Range(1, 255).Draw(t, "rightB"),
			A: 255,
		}

		leftFB := leftRegion.Surface().FrameBuffer()
		rightFB := rightRegion.Surface().FrameBuffer()

		for py := 0; py < regionH; py++ {
			for px := 0; px < regionW; px++ {
				leftFB.SetRGBA(px, py, leftColor)
				rightFB.SetRGBA(px, py, rightColor)
			}
		}

		// Record the "right" region's state before the change.
		rightModeBefore := rightRegion.CurrentMode()

		// Change mode on the "left" region.
		err = rm.SetMode("left", "modeC")
		if err != nil {
			t.Fatalf("SetMode left: %v", err)
		}

		// Verify the "right" region's mode is unchanged.
		if rightRegion.CurrentMode() != rightModeBefore {
			t.Fatalf("right region mode changed: got %q, want %q",
				rightRegion.CurrentMode(), rightModeBefore)
		}

		// Verify the "right" region's surface content is unchanged.
		for py := 0; py < regionH; py++ {
			for px := 0; px < regionW; px++ {
				got := rightFB.RGBAAt(px, py)
				if got != rightColor {
					t.Fatalf("right region pixel(%d,%d) = %v, want %v (modified by left mode change)",
						px, py, got, rightColor)
				}
			}
		}
	})
}

func TestProperty16_ModeSwitchByIndexOrNameWithValidation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a set of valid mode names.
		numValidModes := rapid.IntRange(1, 5).Draw(t, "numValidModes")
		validModes := make(map[string]bool)
		validModeList := make([]string, 0, numValidModes)
		for i := 0; i < numValidModes; i++ {
			mode := fmt.Sprintf("mode_%d", i)
			validModes[mode] = true
			validModeList = append(validModeList, mode)
		}

		vdW := rapid.IntRange(10, 200).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 200).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		rm := region2.NewRegionManager(vd)
		rm.SetModeValidator(func(mode string) bool {
			return validModes[mode]
		})

		// Allocate a region with the first valid mode.
		initialMode := validModeList[0]
		err = rm.Allocate(region2.RegionSpec{
			Name:        "target",
			Bounds:      image.Rect(0, 0, vdW, vdH),
			DefaultMode: initialMode,
		})
		if err != nil {
			t.Fatalf("Allocate: %v", err)
		}

		r, _ := rm.Region(0)

		// Wire ModeFactory so SetMode can construct instances.
		r.SetModeFactory(func(id string, hints textlayout.TextHints) (region2.ModeInstance, bool) {
			return &noopInstance{id: id}, true
		})

		// Test SetMode with a valid registered mode by index.
		newModeIdx := rapid.IntRange(0, len(validModeList)-1).Draw(t, "newModeIdx")
		newMode := validModeList[newModeIdx]

		err = rm.SetMode("0", newMode)
		if err != nil {
			t.Fatalf("SetMode by index with valid mode %q: %v", newMode, err)
		}
		if r.CurrentMode() != newMode {
			t.Fatalf("after SetMode by index: CurrentMode()=%q, want %q", r.CurrentMode(), newMode)
		}

		// Test SetMode with a valid registered mode by name.
		anotherModeIdx := rapid.IntRange(0, len(validModeList)-1).Draw(t, "anotherModeIdx")
		anotherMode := validModeList[anotherModeIdx]

		err = rm.SetMode("target", anotherMode)
		if err != nil {
			t.Fatalf("SetMode by name with valid mode %q: %v", anotherMode, err)
		}
		if r.CurrentMode() != anotherMode {
			t.Fatalf("after SetMode by name: CurrentMode()=%q, want %q", r.CurrentMode(), anotherMode)
		}

		// Test SetMode with an unregistered mode — should error and leave mode unchanged.
		currentMode := r.CurrentMode()
		invalidMode := "unregistered_mode_xyz"

		err = rm.SetMode("0", invalidMode)
		if err == nil {
			t.Fatal("SetMode with unregistered mode should return error")
		}
		if !strings.Contains(err.Error(), "is not registered") {
			t.Fatalf("unexpected error message: %v", err)
		}
		if r.CurrentMode() != currentMode {
			t.Fatalf("mode changed after invalid SetMode: got %q, want %q", r.CurrentMode(), currentMode)
		}

		// Test SetMode by name with unregistered mode.
		err = rm.SetMode("target", invalidMode)
		if err == nil {
			t.Fatal("SetMode by name with unregistered mode should return error")
		}
		if r.CurrentMode() != currentMode {
			t.Fatalf("mode changed after invalid SetMode by name: got %q, want %q", r.CurrentMode(), currentMode)
		}
	})
}

func TestProperty17_RegionNameLookupIsCaseInsensitive(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random region name (1-20 alphanumeric chars for reasonable names).
		name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9]{0,19}`).Draw(t, "regionName")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, 100, 100))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		rm := region2.NewRegionManager(vd)

		err = rm.Allocate(region2.RegionSpec{
			Name:   name,
			Bounds: image.Rect(0, 0, 100, 100),
		})
		if err != nil {
			t.Fatalf("Allocate(%q): %v", name, err)
		}

		// Generate random case variations of the name.
		numVariations := rapid.IntRange(1, 5).Draw(t, "numVariations")
		for i := 0; i < numVariations; i++ {
			variation := randomCaseVariation(t, name, i)

			r, ok := rm.RegionByName(variation)
			if !ok {
				t.Fatalf("RegionByName(%q) not found (original name: %q)", variation, name)
			}
			if r.Name() != name {
				t.Fatalf("RegionByName(%q).Name() = %q, want preserved original %q",
					variation, r.Name(), name)
			}
		}

		// Also test that the original name works.
		r, ok := rm.RegionByName(name)
		if !ok {
			t.Fatalf("RegionByName(%q) not found with original name", name)
		}
		if r.Name() != name {
			t.Fatalf("RegionByName original: got %q, want %q", r.Name(), name)
		}

		// Test all-uppercase.
		r, ok = rm.RegionByName(strings.ToUpper(name))
		if !ok {
			t.Fatalf("RegionByName(upper %q) not found", strings.ToUpper(name))
		}
		if r.Name() != name {
			t.Fatalf("RegionByName upper: got %q, want %q", r.Name(), name)
		}

		// Test all-lowercase.
		r, ok = rm.RegionByName(strings.ToLower(name))
		if !ok {
			t.Fatalf("RegionByName(lower %q) not found", strings.ToLower(name))
		}
		if r.Name() != name {
			t.Fatalf("RegionByName lower: got %q, want %q", r.Name(), name)
		}
	})
}

// TestProperty2_RegionValidationBehavioralEquivalence verifies that the unified
// validateSpecCore produces the same accept/reject decision and error message
// as validateSpec for any RegionSpec and existing region state.
// Since validateSpec delegates to validateSpecCore, this test confirms the
// unified function correctly handles all validation scenarios:
// - Empty names, names exceeding 64 characters
// - Zero-dimension bounds (width=0 or height=0)
// - Bounds outside virtual display
// - Overlapping bounds with existing regions
// - Valid specs that should be accepted
// - Name uniqueness (case-insensitive)
func TestProperty2_RegionValidationBehavioralEquivalence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a VirtualDisplay with reasonable dimensions.
		vdW := rapid.IntRange(10, 300).Draw(t, "vdWidth")
		vdH := rapid.IntRange(10, 300).Draw(t, "vdHeight")
		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatal(err)
		}

		rm := region2.NewRegionManager(vd)

		// Optionally set a mode validator.
		useModeValidator := rapid.Bool().Draw(t, "useModeValidator")
		validModes := map[string]bool{"modeA": true, "modeB": true, "modeC": true}
		if useModeValidator {
			rm.SetModeValidator(func(mode string) bool {
				return validModes[mode]
			})
		}

		// Pre-allocate some existing regions to test overlap/name-uniqueness checks.
		numExisting := rapid.IntRange(0, 3).Draw(t, "numExisting")
		existingBounds := make([]image.Rectangle, 0, numExisting)
		existingNames := make(map[string]bool)

		for i := 0; i < numExisting; i++ {
			name := genUniqueName(t, existingNames, i)
			bounds := genNonOverlappingBounds(t, vdW, vdH, existingBounds, i)
			if bounds.Empty() {
				break
			}
			existingNames[strings.ToLower(name)] = true
			existingBounds = append(existingBounds, bounds)
			if allocErr := rm.Allocate(region2.RegionSpec{Name: name, Bounds: bounds}); allocErr != nil {
				t.Fatal(allocErr)
			}
		}

		// Generate a spec using various scenarios that may be valid or invalid.
		spec := genValidationTestSpec(t, vdW, vdH, existingNames, existingBounds, useModeValidator, validModes)

		// Call both validateSpec and validateSpecCore and compare results.
		errSpec := rm.TestValidateSpec(spec)
		errCore := rm.TestValidateSpecCore(spec)

		// Both must agree on accept/reject.
		if (errSpec == nil) != (errCore == nil) {
			t.Fatalf("validateSpec and validateSpecCore disagree on spec %+v:\n  validateSpec: %v\n  validateSpecCore: %v",
				spec, errSpec, errCore)
		}

		// If both reject, error messages must be identical.
		if errSpec != nil && errCore != nil {
			if errSpec.Error() != errCore.Error() {
				t.Fatalf("validateSpec and validateSpecCore produce different error messages for spec %+v:\n  validateSpec: %q\n  validateSpecCore: %q",
					spec, errSpec.Error(), errCore.Error())
			}
		}

		// Additionally verify the unified function produces correct decisions:
		// Manually compute expected outcome.
		nameValid := len(spec.Name) >= 1 && len(spec.Name) <= 64
		nameUnique := !existingNames[strings.ToLower(spec.Name)]
		boundsPositive := spec.Bounds.Dx() >= 1 && spec.Bounds.Dy() >= 1
		boundsWithinVD := spec.Bounds.In(image.Rect(0, 0, vdW, vdH))
		noOverlap := true
		for _, eb := range existingBounds {
			if !spec.Bounds.Intersect(eb).Empty() {
				noOverlap = false
				break
			}
		}
		modeOK := true
		if spec.DefaultMode != "" && useModeValidator {
			modeOK = validModes[spec.DefaultMode]
		}

		allConstraints := nameValid && nameUnique && boundsPositive && boundsWithinVD && noOverlap && modeOK

		if allConstraints && errCore != nil {
			t.Fatalf("expected validateSpecCore to accept spec %+v (all constraints satisfied), got error: %v\n"+
				"  nameValid=%v, nameUnique=%v, boundsPositive=%v, boundsWithinVD=%v, noOverlap=%v, modeOK=%v",
				spec, errCore, nameValid, nameUnique, boundsPositive, boundsWithinVD, noOverlap, modeOK)
		}

		if !allConstraints && errCore == nil {
			t.Fatalf("expected validateSpecCore to reject spec %+v (constraint violated), but it accepted.\n"+
				"  nameValid=%v, nameUnique=%v, boundsPositive=%v, boundsWithinVD=%v, noOverlap=%v, modeOK=%v",
				spec, nameValid, nameUnique, boundsPositive, boundsWithinVD, noOverlap, modeOK)
		}
	})
}

// genValidationTestSpec generates a RegionSpec designed to exercise various
// validation scenarios with a mix of valid and invalid inputs.
func genValidationTestSpec(t *rapid.T, vdW, vdH int, existingNames map[string]bool, existingBounds []image.Rectangle, useModeValidator bool, validModes map[string]bool) region2.RegionSpec {
	scenario := rapid.IntRange(0, 9).Draw(t, "scenario")

	var name string
	var bounds image.Rectangle
	var defaultMode string

	switch scenario {
	case 0: // Empty name
		name = ""
		bounds = image.Rect(0, 0, 1, 1)
	case 1: // Name exceeding 64 characters
		nameLen := rapid.IntRange(65, 200).Draw(t, "longNameLen")
		name = strings.Repeat("n", nameLen)
		bounds = image.Rect(0, 0, 1, 1)
	case 2: // Zero-width bounds
		name = genFreshName(t, existingNames, "zeroW")
		h := rapid.IntRange(1, vdH).Draw(t, "zeroWH")
		x := rapid.IntRange(0, vdW-1).Draw(t, "zeroWX")
		bounds = image.Rect(x, 0, x, h)
	case 3: // Zero-height bounds
		name = genFreshName(t, existingNames, "zeroH")
		w := rapid.IntRange(1, vdW).Draw(t, "zeroHW")
		y := rapid.IntRange(0, vdH-1).Draw(t, "zeroHY")
		bounds = image.Rect(0, y, w, y)
	case 4: // Bounds outside virtual display
		name = genFreshName(t, existingNames, "oob")
		bounds = genOutOfBounds(t, vdW, vdH)
	case 5: // Overlapping bounds with existing region
		name = genFreshName(t, existingNames, "overlap")
		if len(existingBounds) > 0 {
			idx := rapid.IntRange(0, len(existingBounds)-1).Draw(t, "overlapTarget")
			eb := existingBounds[idx]
			x := rapid.IntRange(eb.Min.X, max(eb.Min.X, eb.Max.X-1)).Draw(t, "overlapX")
			y := rapid.IntRange(eb.Min.Y, max(eb.Min.Y, eb.Max.Y-1)).Draw(t, "overlapY")
			x2 := rapid.IntRange(x+1, min(vdW, x+50)).Draw(t, "overlapX2")
			y2 := rapid.IntRange(y+1, min(vdH, y+50)).Draw(t, "overlapY2")
			bounds = image.Rect(x, y, x2, y2)
		} else {
			w := rapid.IntRange(1, vdW).Draw(t, "overlapW")
			h := rapid.IntRange(1, vdH).Draw(t, "overlapH")
			bounds = image.Rect(0, 0, w, h)
		}
	case 6: // Duplicate name (case-insensitive)
		if len(existingNames) > 0 {
			var existingName string
			for n := range existingNames {
				existingName = n
				break
			}
			name = strings.ToUpper(existingName)
			bounds = genNonOverlappingBounds(t, vdW, vdH, existingBounds, 99)
			if bounds.Empty() {
				bounds = image.Rect(0, 0, 1, 1)
			}
		} else {
			name = genFreshName(t, existingNames, "dup")
			w := rapid.IntRange(1, vdW).Draw(t, "dupW")
			h := rapid.IntRange(1, vdH).Draw(t, "dupH")
			bounds = image.Rect(0, 0, w, h)
		}
	case 7: // Valid spec that should be accepted
		name = genFreshName(t, existingNames, "valid")
		bounds = genNonOverlappingBounds(t, vdW, vdH, existingBounds, 88)
		if bounds.Empty() {
			bounds = image.Rect(0, 0, 1, 1)
		}
	case 8: // Invalid default mode
		name = genFreshName(t, existingNames, "badmode")
		bounds = genNonOverlappingBounds(t, vdW, vdH, existingBounds, 77)
		if bounds.Empty() {
			bounds = image.Rect(0, 0, 1, 1)
		}
		if useModeValidator {
			defaultMode = "unregistered_mode"
		}
	case 9: // Valid spec with valid default mode
		name = genFreshName(t, existingNames, "goodmode")
		bounds = genNonOverlappingBounds(t, vdW, vdH, existingBounds, 66)
		if bounds.Empty() {
			bounds = image.Rect(0, 0, 1, 1)
		}
		if useModeValidator {
			defaultMode = "modeA"
		}
	}

	return region2.RegionSpec{
		Name:        name,
		Bounds:      bounds,
		DefaultMode: defaultMode,
	}
}

// genFreshName generates a unique name that isn't in the existing set.
func genFreshName(t *rapid.T, existing map[string]bool, prefix string) string {
	base := rapid.StringMatching(`[a-z]{1,40}`).Draw(t, prefix+"Name")
	name := prefix + base
	for existing[strings.ToLower(name)] {
		name = name + "x"
	}
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

// --- Helper generators ---

// genValidName generates a valid name between 1-64 characters.
func genValidName(t *rapid.T, label string) string {
	return rapid.StringMatching(`[a-z]{1,64}`).Draw(t, label)
}

// genUniqueName generates a unique valid name not in the existing set.
func genUniqueName(t *rapid.T, existing map[string]bool, idx int) string {
	base := rapid.StringMatching(`[a-z]{1,50}`).Draw(t, fmt.Sprintf("name%d", idx))
	name := base
	for existing[strings.ToLower(name)] {
		name = name + "x"
		if len(name) > 64 {
			name = name[:60] + "uniq"
		}
	}
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

// genNonOverlappingBounds generates bounds within VD that don't overlap with existing.
func genNonOverlappingBounds(t *rapid.T, vdW, vdH int, existing []image.Rectangle, idx int) image.Rectangle {
	// Try to generate random non-overlapping bounds.
	for attempts := 0; attempts < 10; attempts++ {
		label := fmt.Sprintf("b%d_%d", idx, attempts)
		w := rapid.IntRange(1, max(1, vdW/4)).Draw(t, label+"w")
		h := rapid.IntRange(1, max(1, vdH/4)).Draw(t, label+"h")
		x := rapid.IntRange(0, max(0, vdW-w)).Draw(t, label+"x")
		y := rapid.IntRange(0, max(0, vdH-h)).Draw(t, label+"y")

		bounds := image.Rect(x, y, x+w, y+h)

		if !bounds.In(image.Rect(0, 0, vdW, vdH)) {
			continue
		}

		overlaps := false
		for _, eb := range existing {
			if !bounds.Intersect(eb).Empty() {
				overlaps = true
				break
			}
		}
		if !overlaps {
			return bounds
		}
	}
	return image.Rectangle{}
}

// genOutOfBounds generates bounds that extend outside the VD.
func genOutOfBounds(t *rapid.T, vdW, vdH int) image.Rectangle {
	direction := rapid.IntRange(0, 3).Draw(t, "direction")
	switch direction {
	case 0: // extends past right edge
		x := rapid.IntRange(0, vdW-1).Draw(t, "oobX")
		w := rapid.IntRange(vdW-x+1, vdW-x+100).Draw(t, "oobW")
		h := rapid.IntRange(1, vdH).Draw(t, "oobH")
		return image.Rect(x, 0, x+w, h)
	case 1: // extends past bottom edge
		y := rapid.IntRange(0, vdH-1).Draw(t, "oobY")
		w := rapid.IntRange(1, vdW).Draw(t, "oobW")
		h := rapid.IntRange(vdH-y+1, vdH-y+100).Draw(t, "oobH")
		return image.Rect(0, y, w, y+h)
	case 2: // extends past left edge (negative X)
		w := rapid.IntRange(1, vdW).Draw(t, "oobW")
		h := rapid.IntRange(1, vdH).Draw(t, "oobH")
		return image.Rect(-rapid.IntRange(1, 50).Draw(t, "negX"), 0, w, h)
	default: // extends past top edge (negative Y)
		w := rapid.IntRange(1, vdW).Draw(t, "oobW")
		h := rapid.IntRange(1, vdH).Draw(t, "oobH")
		return image.Rect(0, -rapid.IntRange(1, 50).Draw(t, "negY"), w, h)
	}
}

// genArbitraryName generates a name that may or may not be valid.
func genArbitraryName(t *rapid.T) string {
	choice := rapid.IntRange(0, 9).Draw(t, "nameChoice")
	switch {
	case choice < 2: // empty
		return ""
	case choice == 2: // too long
		length := rapid.IntRange(65, 128).Draw(t, "longNameLen")
		return strings.Repeat("x", length)
	default: // valid
		return rapid.StringMatching(`[a-z]{1,64}`).Draw(t, "validName")
	}
}

// genArbitraryBounds generates bounds that may or may not satisfy constraints.
func genArbitraryBounds(t *rapid.T, vdW, vdH int) image.Rectangle {
	choice := rapid.IntRange(0, 9).Draw(t, "boundsChoice")
	switch {
	case choice < 1: // zero-size bounds
		x := rapid.IntRange(0, vdW).Draw(t, "zX")
		y := rapid.IntRange(0, vdH).Draw(t, "zY")
		return image.Rect(x, y, x, y)
	case choice < 2: // out of bounds
		return genOutOfBounds(t, vdW, vdH)
	default: // within VD with positive size
		w := rapid.IntRange(1, vdW).Draw(t, "arbW")
		h := rapid.IntRange(1, vdH).Draw(t, "arbH")
		x := rapid.IntRange(0, vdW-1).Draw(t, "arbX")
		y := rapid.IntRange(0, vdH-1).Draw(t, "arbY")
		return image.Rect(x, y, x+w, y+h)
	}
}

// toggleCase toggles the case of at least one character in the string.
func toggleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'a' && r <= 'z' {
			runes[i] = r - 32
			return string(runes)
		}
		if r >= 'A' && r <= 'Z' {
			runes[i] = r + 32
			return string(runes)
		}
	}
	return s
}

// randomCaseVariation generates a random case variation of a string.
func randomCaseVariation(t *rapid.T, s string, idx int) string {
	runes := []rune(s)
	result := make([]rune, len(runes))
	for i, r := range runes {
		if rapid.Bool().Draw(t, fmt.Sprintf("case_%d_%d", idx, i)) {
			result[i] = unicode.ToUpper(r)
		} else {
			result[i] = unicode.ToLower(r)
		}
	}
	return string(result)
}

// =============================================================================
// From: mode_switch_property_test.go
// =============================================================================

// For any Region with current mode A, when SetMode is called with a valid mode B
// (B ≠ A), the Region SHALL: stop the current ModeEngine, clear the surface to
// all-black pixels, and set CurrentMode() to B.
func TestProperty10_ModeChangeOnValidModeClearsAndTransitions(t *testing.T) {
	t.Run("valid_mode_change_transitions_and_clears", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a set of valid modes (at least 2 so we can change from A to B).
			numModes := rapid.IntRange(2, 10).Draw(t, "numModes")
			validModes := make([]string, numModes)
			validModeSet := make(map[string]bool)
			for i := 0; i < numModes; i++ {
				validModes[i] = fmt.Sprintf("mode_%d", i)
				validModeSet[validModes[i]] = true
			}

			// Generate VD and region dimensions.
			vdW := rapid.IntRange(2, 200).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 200).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return validModeSet[mode]
			})

			// Pick initial mode A and target mode B (A ≠ B).
			idxA := rapid.IntRange(0, numModes-1).Draw(t, "modeAIdx")
			idxB := rapid.IntRange(0, numModes-2).Draw(t, "modeBIdx")
			if idxB >= idxA {
				idxB++ // ensure B ≠ A
			}
			modeA := validModes[idxA]
			modeB := validModes[idxB]

			// Allocate a region with mode A.
			err = rm.Allocate(region2.RegionSpec{
				Name:        "target",
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: modeA,
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.Region(0)

			// Wire ModeFactory so SetMode can construct instances.
			r.SetModeFactory(noopModeFactory())

			// Fill the region's surface with random non-black pixels.
			fb := r.Surface().FrameBuffer()
			w := r.Bounds().Dx()
			h := r.Bounds().Dy()
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					c := color.RGBA{
						R: rapid.Uint8Range(1, 255).Draw(t, "r"),
						G: rapid.Uint8Range(0, 255).Draw(t, "g"),
						B: rapid.Uint8Range(0, 255).Draw(t, "b"),
						A: 255,
					}
					fb.SetRGBA(px, py, c)
				}
			}

			// Execute mode change from A to B.
			err = rm.SetMode("0", modeB)
			if err != nil {
				t.Fatalf("SetMode(%q → %q): %v", modeA, modeB, err)
			}

			// Verify: CurrentMode is now B.
			if r.CurrentMode() != modeB {
				t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), modeB)
			}

			// Verify: all surface pixels are black (RGBA 0,0,0,255).
			black := color.RGBA{0, 0, 0, 255}
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					got := fb.RGBAAt(px, py)
					if got != black {
						t.Fatalf("pixel(%d,%d)=%v after SetMode, want %v (surface should be cleared to black)",
							px, py, got, black)
					}
				}
			}
		})
	})

	t.Run("same_mode_reinitializes_fully", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a valid mode.
			modeA := fmt.Sprintf("mode_%d", rapid.IntRange(0, 99).Draw(t, "modeIdx"))

			vdW := rapid.IntRange(2, 100).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 100).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return strings.HasPrefix(mode, "mode_")
			})

			err = rm.Allocate(region2.RegionSpec{
				Name:        "reinit",
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: modeA,
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.Region(0)

			// Wire ModeFactory so SetMode can construct instances.
			r.SetModeFactory(noopModeFactory())

			// Dirty the surface.
			fb := r.Surface().FrameBuffer()
			w := r.Bounds().Dx()
			h := r.Bounds().Dy()
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					c := color.RGBA{
						R: rapid.Uint8Range(1, 255).Draw(t, "r"),
						G: rapid.Uint8Range(0, 255).Draw(t, "g"),
						B: rapid.Uint8Range(0, 255).Draw(t, "b"),
						A: 255,
					}
					fb.SetRGBA(px, py, c)
				}
			}

			// Set mode to the SAME mode (full re-initialization per Req 5.6).
			err = rm.SetMode("0", modeA)
			if err != nil {
				t.Fatalf("SetMode same mode: %v", err)
			}

			// Verify: CurrentMode is still A.
			if r.CurrentMode() != modeA {
				t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), modeA)
			}

			// Verify: surface cleared to black.
			black := color.RGBA{0, 0, 0, 255}
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					got := fb.RGBAAt(px, py)
					if got != black {
						t.Fatalf("pixel(%d,%d)=%v after same-mode SetMode, want %v",
							px, py, got, black)
					}
				}
			}
		})
	})

	t.Run("mode_change_by_name_clears_and_transitions", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			modeA := "alpha"
			modeB := "beta"

			vdW := rapid.IntRange(2, 150).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 150).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return mode == modeA || mode == modeB
			})

			// Generate a valid region name.
			regionName := rapid.StringMatching(`[a-z]{1,20}`).Draw(t, "regionName")

			err = rm.Allocate(region2.RegionSpec{
				Name:        regionName,
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: modeA,
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.RegionByName(regionName)

			// Wire ModeFactory so SetMode can construct instances.
			r.SetModeFactory(noopModeFactory())

			// Dirty surface.
			r.Surface().Clear(color.RGBA{255, 128, 64, 255})

			// Set mode by name.
			err = rm.SetMode(regionName, modeB)
			if err != nil {
				t.Fatalf("SetMode by name: %v", err)
			}

			// Verify transition.
			if r.CurrentMode() != modeB {
				t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), modeB)
			}

			// Verify surface is black.
			black := color.RGBA{0, 0, 0, 255}
			fb := r.Surface().FrameBuffer()
			w := r.Bounds().Dx()
			h := r.Bounds().Dy()
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					got := fb.RGBAAt(px, py)
					if got != black {
						t.Fatalf("pixel(%d,%d)=%v, want %v", px, py, got, black)
					}
				}
			}
		})
	})
}

// For any Region with current mode A, when SetMode is called with a mode string
// not in the mode registry, the Region SHALL retain mode A unchanged, its surface
// SHALL remain unmodified, and an error SHALL be returned.
func TestProperty11_InvalidModeChangePreservesState(t *testing.T) {
	t.Run("unregistered_mode_preserves_mode_and_surface", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate a set of valid modes.
			numValidModes := rapid.IntRange(1, 5).Draw(t, "numValidModes")
			validModes := make([]string, numValidModes)
			validModeSet := make(map[string]bool)
			for i := 0; i < numValidModes; i++ {
				validModes[i] = fmt.Sprintf("valid_%d", i)
				validModeSet[validModes[i]] = true
			}

			// Generate an invalid mode name that is NOT in the valid set.
			invalidMode := rapid.StringMatching(`invalid_[a-z]{1,20}`).Draw(t, "invalidMode")

			// Generate VD and region dimensions.
			vdW := rapid.IntRange(2, 150).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 150).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return validModeSet[mode]
			})

			// Allocate a region with a valid initial mode.
			initialMode := validModes[rapid.IntRange(0, numValidModes-1).Draw(t, "initialModeIdx")]

			err = rm.Allocate(region2.RegionSpec{
				Name:        "preserve",
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: initialMode,
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.Region(0)

			// Fill the surface with random colored pixels and record their state.
			fb := r.Surface().FrameBuffer()
			w := r.Bounds().Dx()
			h := r.Bounds().Dy()
			originalPixels := make([]color.RGBA, w*h)
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					c := color.RGBA{
						R: rapid.Uint8Range(0, 255).Draw(t, "r"),
						G: rapid.Uint8Range(0, 255).Draw(t, "g"),
						B: rapid.Uint8Range(0, 255).Draw(t, "b"),
						A: 255,
					}
					fb.SetRGBA(px, py, c)
					originalPixels[py*w+px] = c
				}
			}

			// Attempt to set an invalid mode.
			err = rm.SetMode("0", invalidMode)

			// Verify: error is returned.
			if err == nil {
				t.Fatal("SetMode with unregistered mode should return error")
			}
			if !strings.Contains(err.Error(), "not registered") {
				t.Fatalf("error should mention 'not registered', got: %v", err)
			}

			// Verify: current mode is unchanged.
			if r.CurrentMode() != initialMode {
				t.Fatalf("CurrentMode()=%q, want %q (should be unchanged)", r.CurrentMode(), initialMode)
			}

			// Verify: surface pixels are unchanged.
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					got := fb.RGBAAt(px, py)
					want := originalPixels[py*w+px]
					if got != want {
						t.Fatalf("pixel(%d,%d)=%v, want %v (surface should be unmodified after invalid mode change)",
							px, py, got, want)
					}
				}
			}
		})
	})

	t.Run("invalid_mode_by_name_preserves_state", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			validMode := "dashboard"

			vdW := rapid.IntRange(2, 100).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 100).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return mode == validMode
			})

			regionName := rapid.StringMatching(`[a-z]{1,20}`).Draw(t, "regionName")

			err = rm.Allocate(region2.RegionSpec{
				Name:        regionName,
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: validMode,
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.RegionByName(regionName)

			// Fill surface with a consistent color and record it.
			fillColor := color.RGBA{
				R: rapid.Uint8Range(0, 255).Draw(t, "fillR"),
				G: rapid.Uint8Range(0, 255).Draw(t, "fillG"),
				B: rapid.Uint8Range(0, 255).Draw(t, "fillB"),
				A: 255,
			}
			r.Surface().Clear(fillColor)

			// Generate an invalid mode.
			invalidMode := rapid.StringMatching(`bad_[a-z]{1,15}`).Draw(t, "badMode")

			// Attempt SetMode by name with invalid mode.
			err = rm.SetMode(regionName, invalidMode)

			// Verify error.
			if err == nil {
				t.Fatal("SetMode with unregistered mode should return error")
			}

			// Verify mode unchanged.
			if r.CurrentMode() != validMode {
				t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), validMode)
			}

			// Verify surface unchanged.
			fb := r.Surface().FrameBuffer()
			w := r.Bounds().Dx()
			h := r.Bounds().Dy()
			for py := 0; py < h; py++ {
				for px := 0; px < w; px++ {
					got := fb.RGBAAt(px, py)
					if got != fillColor {
						t.Fatalf("pixel(%d,%d)=%v, want %v (surface unchanged)", px, py, got, fillColor)
					}
				}
			}
		})
	})

	t.Run("nonexistent_region_target_returns_error", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			vdW := rapid.IntRange(2, 200).Draw(t, "vdWidth")
			vdH := rapid.IntRange(2, 200).Draw(t, "vdHeight")

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
			if err != nil {
				t.Fatalf("NewVirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			rm.SetModeValidator(func(mode string) bool {
				return mode == "clock"
			})

			err = rm.Allocate(region2.RegionSpec{
				Name:        "existing",
				Bounds:      image.Rect(0, 0, vdW, vdH),
				DefaultMode: "clock",
			})
			if err != nil {
				t.Fatalf("Allocate: %v", err)
			}

			r, _ := rm.Region(0)

			// Record state before bad target attempt.
			modeBefore := r.CurrentMode()

			// Generate a nonexistent target (either an out-of-range index or unknown name).
			choice := rapid.IntRange(0, 1).Draw(t, "targetChoice")
			var badTarget string
			if choice == 0 {
				// Out-of-range index.
				badIdx := rapid.IntRange(1, 100).Draw(t, "badIdx")
				badTarget = fmt.Sprintf("%d", badIdx)
			} else {
				// Unknown name.
				badTarget = rapid.StringMatching(`nonexist_[a-z]{1,10}`).Draw(t, "badName")
			}

			err = rm.SetMode(badTarget, "clock")

			// Verify error is returned.
			if err == nil {
				t.Fatalf("SetMode with nonexistent target %q should return error", badTarget)
			}

			// Verify original region state is unchanged.
			if r.CurrentMode() != modeBefore {
				t.Fatalf("mode changed after failed SetMode: got %q, want %q", r.CurrentMode(), modeBefore)
			}
		})
	})
}

// =============================================================================
// From: render_failure_isolation_property_test.go
// =============================================================================

// --- Property 13: Render failure isolation ---

// For any set of N regions where K regions fail to render (via panic or error
// return), the remaining N-K regions SHALL render successfully, the failed
// regions' surfaces SHALL remain unchanged from their prior state, and exactly
// one Flush call SHALL occur.

// failureIsolationRenderer renders healthy regions by writing a known pixel,
// and panics or errors for failed regions. It tracks which regions were
// successfully rendered.
type failureIsolationRenderer struct {
	mu          sync.Mutex
	rendered    []string
	panicNames  map[string]bool
	errorNames  map[string]bool
	renderColor color.RGBA
}

func (r *failureIsolationRenderer) Render(reg *region2.Region) error {
	name := reg.Name()

	if r.panicNames[name] {
		panic("intentional panic for " + name)
	}
	if r.errorNames[name] {
		return errors.New("intentional error for " + name)
	}

	// Successful render: write a known pixel to prove the surface was modified.
	reg.Surface().FrameBuffer().SetRGBA(0, 0, r.renderColor)

	r.mu.Lock()
	r.rendered = append(r.rendered, name)
	r.mu.Unlock()
	return nil
}

func (r *failureIsolationRenderer) getRendered() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, len(r.rendered))
	copy(result, r.rendered)
	return result
}

func TestProperty13_RenderFailureIsolation(t *testing.T) {
	t.Run("failed_regions_surface_unchanged_healthy_render_flush_once", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// Generate between 2 and 8 regions.
			numRegions := rapid.IntRange(2, 8).Draw(t, "numRegions")
			numFailing := rapid.IntRange(1, numRegions-1).Draw(t, "numFailing")

			failIndices := make(map[int]bool)
			for len(failIndices) < numFailing {
				idx := rapid.IntRange(0, numRegions-1).Draw(t, "failIdx")
				failIndices[idx] = true
			}

			panicNames := make(map[string]bool)
			errorNames := make(map[string]bool)
			for idx := range failIndices {
				usePanic := rapid.Bool().Draw(t, "usePanic")
				name := regionName(idx)
				if usePanic {
					panicNames[name] = true
				} else {
					errorNames[name] = true
				}
			}

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)

			priorColor := color.RGBA{R: 42, G: 84, B: 126, A: 255}
			renderColor := color.RGBA{R: 200, G: 100, B: 50, A: 255}

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := regionName(i)
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
				surf.FrameBuffer().SetRGBA(0, 0, priorColor)
			}

			if regions := rm.Regions(); len(regions) > 0 {
				regions[0].SetInputFocus(true)
			}

			target := &flushMockTarget{bounds: image.Rect(0, 0, totalWidth, 60)}
			screens := []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, totalWidth, 60), Target: target},
			}
			fp := region2.NewFlushPath(vd, screens)

			renderer := &failureIsolationRenderer{
				panicNames:  panicNames,
				errorNames:  errorNames,
				renderColor: renderColor,
			}

			mockResolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 50 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, fp, nil,
				region2.WithTickRateResolver(mockResolver),
				region2.WithRenderer(renderer),
			)

			now := time.Now()
			tickers := make([]region2.TestPerRegionTicker, numRegions)
			for i, r := range regions {
				tickers[i] = region2.TestPerRegionTicker{
					Region:   r,
					Interval: 50 * time.Millisecond,
					LastFire: now.Add(-100 * time.Millisecond),
				}
			}
			rl.TestSetRegionTickers(tickers)

			flushBefore := target.callCount
			rl.TestRenderDueRegions(now)
			rl.TestFlush()

			flushDelta := target.callCount - flushBefore
			if flushDelta != 1 {
				t.Fatalf("expected exactly 1 flush call, got %d", flushDelta)
			}

			rendered := renderer.getRendered()
			renderedSet := make(map[string]bool)
			for _, name := range rendered {
				renderedSet[name] = true
			}

			for i := 0; i < numRegions; i++ {
				name := regionName(i)
				if failIndices[i] {
					if renderedSet[name] {
						t.Fatalf("failed region %q should not have been rendered successfully", name)
					}
				} else {
					if !renderedSet[name] {
						t.Fatalf("healthy region %q was not rendered", name)
					}
				}
			}

			for i := 0; i < numRegions; i++ {
				pixel := regions[i].Surface().FrameBuffer().RGBAAt(0, 0)
				if failIndices[i] {
					if pixel != priorColor {
						t.Fatalf("failed region %q surface was modified: got %v, want %v",
							regionName(i), pixel, priorColor)
					}
				} else {
					if pixel != renderColor {
						t.Fatalf("healthy region %q surface was not updated: got %v, want %v",
							regionName(i), pixel, renderColor)
					}
				}
			}
		})
	})

	t.Run("global_ticker_failure_isolation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(2, 8).Draw(t, "numRegions")
			numFailing := rapid.IntRange(1, numRegions-1).Draw(t, "numFailing")

			failIndices := make(map[int]bool)
			for len(failIndices) < numFailing {
				idx := rapid.IntRange(0, numRegions-1).Draw(t, "failIdx")
				failIndices[idx] = true
			}

			panicNames := make(map[string]bool)
			errorNames := make(map[string]bool)
			for idx := range failIndices {
				usePanic := rapid.Bool().Draw(t, "usePanic")
				name := regionName(idx)
				if usePanic {
					panicNames[name] = true
				} else {
					errorNames[name] = true
				}
			}

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)

			priorColor := color.RGBA{R: 10, G: 20, B: 30, A: 255}
			renderColor := color.RGBA{R: 250, G: 240, B: 230, A: 255}

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := regionName(i)
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
				surf.FrameBuffer().SetRGBA(0, 0, priorColor)
			}

			if regions := rm.Regions(); len(regions) > 0 {
				regions[0].SetInputFocus(true)
			}

			target := &flushMockTarget{bounds: image.Rect(0, 0, totalWidth, 60)}
			screens := []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, totalWidth, 60), Target: target},
			}
			fp := region2.NewFlushPath(vd, screens)

			renderer := &failureIsolationRenderer{
				panicNames:  panicNames,
				errorNames:  errorNames,
				renderColor: renderColor,
			}

			rl := region2.NewRenderLoop(rm, fp, nil,
				region2.WithRenderer(renderer),
			)

			flushBefore := target.callCount
			rl.TestRenderFrame()
			rl.TestFlush()

			flushDelta := target.callCount - flushBefore
			if flushDelta != 1 {
				t.Fatalf("expected exactly 1 flush call, got %d", flushDelta)
			}

			rendered := renderer.getRendered()
			renderedSet := make(map[string]bool)
			for _, name := range rendered {
				renderedSet[name] = true
			}

			for i := 0; i < numRegions; i++ {
				name := regionName(i)
				if failIndices[i] {
					if renderedSet[name] {
						t.Fatalf("failed region %q should not have been rendered", name)
					}
				} else {
					if !renderedSet[name] {
						t.Fatalf("healthy region %q was not rendered", name)
					}
				}
			}

			for i := 0; i < numRegions; i++ {
				pixel := regions[i].Surface().FrameBuffer().RGBAAt(0, 0)
				if failIndices[i] {
					if pixel != priorColor {
						t.Fatalf("failed region %q surface modified: got %v, want %v",
							regionName(i), pixel, priorColor)
					}
				} else {
					if pixel != renderColor {
						t.Fatalf("healthy region %q not updated: got %v, want %v",
							regionName(i), pixel, renderColor)
					}
				}
			}
		})
	})
}

// regionName returns a deterministic name for a region at the given index.
func regionName(idx int) string {
	return string(rune('a'+idx)) + "_iso"
}

// =============================================================================
// From: render_loop_property_test.go
// =============================================================================

// --- Property 6 test helpers ---

// trackingRenderer records which regions were rendered per call to allow assertions.
type trackingRenderer struct {
	mu    sync.Mutex
	calls []string
}

func (tr *trackingRenderer) Render(r *region2.Region) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.calls = append(tr.calls, r.Name())
	return nil
}

func (tr *trackingRenderer) reset() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.calls = nil
}

func (tr *trackingRenderer) getCalls() []string {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	result := make([]string, len(tr.calls))
	copy(result, tr.calls)
	return result
}

// For any set of N regions with varying tick intervals, at a given wake time T,
// the PerRegionRenderLoop SHALL render exactly those regions whose deadline has
// elapsed (deadline ≤ T) and skip all others. A single Flush call SHALL follow
// the rendering of all due regions.
func TestProperty6_SelectiveRegionRenderingByDeadline(t *testing.T) {
	t.Run("only_due_regions_are_rendered", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(2, 8).Draw(t, "numRegions")

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			renderer := &trackingRenderer{}

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := string(rune('a'+i)) + "_prop6"
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
			}

			mockResolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 1000 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(mockResolver),
				region2.WithRenderer(renderer),
			)

			now := time.Now()
			tickers := make([]region2.TestPerRegionTicker, numRegions)

			expectedDue := make(map[string]bool)
			for i := 0; i < numRegions; i++ {
				intervalMs := rapid.IntRange(10, 2000).Draw(t, "intervalMs")
				interval := time.Duration(intervalMs) * time.Millisecond
				elapsedMs := rapid.IntRange(1, 3000).Draw(t, "elapsedMs")
				elapsed := time.Duration(elapsedMs) * time.Millisecond

				tickers[i] = region2.TestPerRegionTicker{
					Region:   regions[i],
					Interval: interval,
					LastFire: now.Add(-elapsed),
				}

				if elapsedMs >= intervalMs {
					expectedDue[regions[i].Name()] = true
				}
			}
			rl.TestSetRegionTickers(tickers)

			renderer.reset()
			rl.TestRenderDueRegions(now)

			calls := renderer.getCalls()
			actualRendered := make(map[string]bool)
			for _, name := range calls {
				actualRendered[name] = true
			}

			for name := range expectedDue {
				if !actualRendered[name] {
					t.Fatalf("region %q was due (elapsed >= interval) but was NOT rendered", name)
				}
			}

			for name := range actualRendered {
				if !expectedDue[name] {
					t.Fatalf("region %q was NOT due but WAS rendered", name)
				}
			}
		})
	})

	t.Run("single_flush_follows_all_due_renders", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(1, 6).Draw(t, "numRegions")

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			renderer := &trackingRenderer{}

			target := &flushMockTarget{bounds: image.Rect(0, 0, totalWidth, 60)}
			screens := []region2.ScreenPosition{
				{Index: 0, Name: "screen0", Bounds: image.Rect(0, 0, totalWidth, 60), Target: target},
			}
			fp := region2.NewFlushPath(vd, screens)

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := string(rune('a'+i)) + "_flush"
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
			}

			mockResolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 50 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, fp, nil,
				region2.WithTickRateResolver(mockResolver),
				region2.WithRenderer(renderer),
			)

			now := time.Now()
			tickers := make([]region2.TestPerRegionTicker, numRegions)
			for i, r := range regions {
				tickers[i] = region2.TestPerRegionTicker{
					Region:   r,
					Interval: 50 * time.Millisecond,
					LastFire: now.Add(-100 * time.Millisecond),
				}
			}
			rl.TestSetRegionTickers(tickers)

			flushBefore := target.callCount
			renderer.reset()
			rl.TestRenderDueRegions(now)
			rl.TestFlush()

			calls := renderer.getCalls()
			if len(calls) != numRegions {
				t.Fatalf("expected %d regions rendered, got %d", numRegions, len(calls))
			}

			flushAfter := target.callCount
			flushDelta := flushAfter - flushBefore
			if flushDelta != 1 {
				t.Fatalf("expected exactly 1 flush call after rendering %d due regions, got %d",
					numRegions, flushDelta)
			}
		})
	})

	t.Run("non_due_regions_are_skipped", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(2, 6).Draw(t, "numRegions")

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			renderer := &trackingRenderer{}

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := string(rune('a'+i)) + "_skip"
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
			}

			mockResolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 1000 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(mockResolver),
				region2.WithRenderer(renderer),
			)

			now := time.Now()
			tickers := make([]region2.TestPerRegionTicker, numRegions)
			for i, r := range regions {
				intervalMs := rapid.IntRange(5000, 10000).Draw(t, "intervalMs")
				tickers[i] = region2.TestPerRegionTicker{
					Region:   r,
					Interval: time.Duration(intervalMs) * time.Millisecond,
					LastFire: now.Add(-1 * time.Millisecond),
				}
			}
			rl.TestSetRegionTickers(tickers)

			renderer.reset()
			rl.TestRenderDueRegions(now)

			calls := renderer.getCalls()
			if len(calls) != 0 {
				t.Fatalf("expected 0 regions rendered when none are due, got %d: %v",
					len(calls), calls)
			}
		})
	})

	t.Run("deadline_advances_after_render", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(1, 5).Draw(t, "numRegions")
			intervalMs := rapid.IntRange(10, 2000).Draw(t, "intervalMs")

			totalWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, totalWidth, 60))
			if err != nil {
				t.Fatalf("failed to create VirtualDisplay: %v", err)
			}

			rm := region2.NewRegionManager(vd)
			renderer := &trackingRenderer{}

			regions := make([]*region2.Region, numRegions)
			for i := 0; i < numRegions; i++ {
				name := string(rune('a'+i)) + "_adv"
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
				regions[i] = r
			}

			mockResolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": time.Duration(intervalMs) * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(mockResolver),
				region2.WithRenderer(renderer),
			)

			interval := time.Duration(intervalMs) * time.Millisecond
			now := time.Now()

			tickers := make([]region2.TestPerRegionTicker, numRegions)
			originalLastFires := make([]time.Time, numRegions)
			for i, r := range regions {
				lastFire := now.Add(-interval - 10*time.Millisecond)
				originalLastFires[i] = lastFire
				tickers[i] = region2.TestPerRegionTicker{
					Region:   r,
					Interval: interval,
					LastFire: lastFire,
				}
			}
			rl.TestSetRegionTickers(tickers)

			rl.TestRenderDueRegions(now)

			for i, ticker := range rl.TestRegionTickers() {
				expectedLastFire := originalLastFires[i].Add(interval)
				if !ticker.LastFire.Equal(expectedLastFire) {
					t.Fatalf("region %d: expected lastFire to advance by interval to %v, got %v",
						i, expectedLastFire, ticker.LastFire)
				}
			}
		})
	})
}

// For any set of region deadlines, the PerRegionRenderLoop's sleep duration SHALL
// equal the minimum (deadline - now) across all regions.
func TestProperty7_SleepDurationIsMinimumDeadlineOffset(t *testing.T) {
	t.Run("sleep_duration_equals_minimum_remaining", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(1, 8).Draw(t, "numRegions")

			vdWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdWidth, 60))
			if err != nil {
				t.Fatal(err)
			}

			rm := region2.NewRegionManager(vd)

			for i := 0; i < numRegions; i++ {
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				name := string(rune('A' + i))
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
			}

			resolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 1000 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(resolver),
			)

			rl.TestInitRegionTickers()

			now := time.Now()
			var expectedMin time.Duration
			expectedMin = time.Duration(1<<63 - 1)

			for i := 0; i < numRegions; i++ {
				intervalMs := rapid.IntRange(50, 5000).Draw(t, "intervalMs")
				interval := time.Duration(intervalMs) * time.Millisecond
				// Ensure at least 20ms remaining to avoid flakiness from clock drift
				// between test setup and minSleepDuration()'s internal time.Now() call.
				maxElapsed := intervalMs - 20
				if maxElapsed < 0 {
					maxElapsed = 0
				}
				elapsedMs := rapid.IntRange(0, maxElapsed).Draw(t, "elapsedMs")
				elapsed := time.Duration(elapsedMs) * time.Millisecond

				rl.TestSetRegionTicker(i, interval, now.Add(-elapsed))

				remaining := interval - elapsed
				if remaining < expectedMin {
					expectedMin = remaining
				}
			}

			sleepDur := rl.TestMinSleepDuration()

			tolerance := 15 * time.Millisecond
			diff := sleepDur - expectedMin
			if diff < 0 {
				diff = -diff
			}
			if diff > tolerance {
				t.Fatalf("minSleepDuration() = %v, expected ~%v (diff=%v exceeds tolerance %v)",
					sleepDur, expectedMin, diff, tolerance)
			}

			if sleepDur <= 0 {
				t.Fatalf("minSleepDuration() = %v, expected positive duration", sleepDur)
			}
		})
	})

	t.Run("sleep_duration_zero_when_deadline_past", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(1, 8).Draw(t, "numRegions")

			vdWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdWidth, 60))
			if err != nil {
				t.Fatal(err)
			}

			rm := region2.NewRegionManager(vd)

			for i := 0; i < numRegions; i++ {
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				name := string(rune('A' + i))
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
			}

			resolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 1000 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(resolver),
			)

			rl.TestInitRegionTickers()

			now := time.Now()
			pastIdx := rapid.IntRange(0, numRegions-1).Draw(t, "pastIdx")

			for i := 0; i < numRegions; i++ {
				intervalMs := rapid.IntRange(10, 5000).Draw(t, "intervalMs")
				interval := time.Duration(intervalMs) * time.Millisecond
				if i == pastIdx {
					pastOffsetMs := rapid.IntRange(intervalMs, intervalMs+5000).Draw(t, "pastOffsetMs")
					rl.TestSetRegionTicker(i, interval, now.Add(-time.Duration(pastOffsetMs)*time.Millisecond))
				} else {
					elapsedMs := rapid.IntRange(0, intervalMs-1).Draw(t, "elapsedMs")
					rl.TestSetRegionTicker(i, interval, now.Add(-time.Duration(elapsedMs)*time.Millisecond))
				}
			}

			sleepDur := rl.TestMinSleepDuration()
			if sleepDur != 0 {
				t.Fatalf("minSleepDuration() = %v, expected 0 when at least one deadline is past", sleepDur)
			}
		})
	})

	t.Run("sleep_duration_is_minimum_not_maximum_or_average", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			numRegions := rapid.IntRange(2, 8).Draw(t, "numRegions")

			vdWidth := numRegions * 60
			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdWidth, 60))
			if err != nil {
				t.Fatal(err)
			}

			rm := region2.NewRegionManager(vd)

			for i := 0; i < numRegions; i++ {
				bounds := image.Rect(i*60, 0, (i+1)*60, 60)
				surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
				name := string(rune('A' + i))
				r := region2.NewRegion(name, bounds, surf)
				r.TestSetMode("default")
				rm.TestAppendRegion(r)
			}

			resolver := &mockTickRateResolver{
				intervals: map[string]time.Duration{
					"default": 1000 * time.Millisecond,
				},
			}

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickRateResolver(resolver),
			)

			rl.TestInitRegionTickers()

			now := time.Now()

			for i := 0; i < numRegions; i++ {
				intervalMs := rapid.IntRange(100, 5000).Draw(t, "intervalMs")
				interval := time.Duration(intervalMs) * time.Millisecond
				maxElapsed := intervalMs - 20
				if maxElapsed < 0 {
					maxElapsed = 0
				}
				elapsedMs := rapid.IntRange(0, maxElapsed).Draw(t, "elapsedMs")
				elapsed := time.Duration(elapsedMs) * time.Millisecond

				rl.TestSetRegionTicker(i, interval, now.Add(-elapsed))
			}

			sleepDur := rl.TestMinSleepDuration()

			tolerance := 5 * time.Millisecond
			for i, rt := range rl.TestRegionTickers() {
				deadline := rt.LastFire.Add(rt.Interval)
				remaining := time.Until(deadline)
				if sleepDur > remaining+tolerance {
					t.Fatalf("minSleepDuration() = %v exceeds region %d remaining %v",
						sleepDur, i, remaining)
				}
			}

			if sleepDur <= 0 {
				t.Fatalf("minSleepDuration() = %v, expected positive", sleepDur)
			}
		})
	})

	t.Run("empty_tickers_returns_default_tick", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			tickMs := rapid.IntRange(10, 5000).Draw(t, "tickMs")
			tick := time.Duration(tickMs) * time.Millisecond

			vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, 60, 60))
			if err != nil {
				t.Fatal(err)
			}

			rm := region2.NewRegionManager(vd)

			rl := region2.NewRenderLoop(rm, nil, nil,
				region2.WithTickInterval(tick),
			)

			sleepDur := rl.TestMinSleepDuration()
			if sleepDur != tick {
				t.Fatalf("minSleepDuration() with no tickers = %v, expected default tick %v",
					sleepDur, tick)
			}
		})
	})
}

// =============================================================================
// From: surface_property_test.go
// =============================================================================

func TestProperty5_RegionSubImageSharesVirtualDisplayFramebufferMemory(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(2, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(2, 500).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		x0 := rapid.IntRange(0, vdW-2).Draw(t, "x0")
		y0 := rapid.IntRange(0, vdH-2).Draw(t, "y0")
		x1 := rapid.IntRange(x0+1, vdW).Draw(t, "x1")
		y1 := rapid.IntRange(y0+1, vdH).Draw(t, "y1")
		subRect := image.Rect(x0, y0, x1, y1)

		surf := surface.NewFromSubImage(vd.FrameBuffer(), subRect)

		localX := rapid.IntRange(0, subRect.Dx()-1).Draw(t, "localX")
		localY := rapid.IntRange(0, subRect.Dy()-1).Draw(t, "localY")

		r := rapid.Uint8().Draw(t, "r")
		g := rapid.Uint8().Draw(t, "g")
		b := rapid.Uint8().Draw(t, "b")
		a := rapid.Uint8().Draw(t, "a")
		testColor := color.RGBA{R: r, G: g, B: b, A: a}

		surf.FrameBuffer().Set(localX, localY, testColor)

		globalX := subRect.Min.X + localX
		globalY := subRect.Min.Y + localY
		vdColor := vd.FrameBuffer().RGBAAt(globalX, globalY)

		if vdColor != testColor {
			t.Fatalf("Surface→VD mismatch: wrote %v at local (%d,%d), read %v at VD (%d,%d)",
				testColor, localX, localY, vdColor, globalX, globalY)
		}

		r2 := rapid.Uint8().Draw(t, "r2")
		g2 := rapid.Uint8().Draw(t, "g2")
		b2 := rapid.Uint8().Draw(t, "b2")
		a2 := rapid.Uint8().Draw(t, "a2")
		testColor2 := color.RGBA{R: r2, G: g2, B: b2, A: a2}

		vd.FrameBuffer().SetRGBA(globalX, globalY, testColor2)

		surfColor := surf.FrameBuffer().RGBAAt(localX, localY)
		if surfColor != testColor2 {
			t.Fatalf("VD→Surface mismatch: wrote %v at VD (%d,%d), read %v at local (%d,%d)",
				testColor2, globalX, globalY, surfColor, localX, localY)
		}
	})
}

func TestProperty8_RegionSurfaceUsesLocalCoordinateSystemWithOriginZeroZero(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(2, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(2, 500).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		x0 := rapid.IntRange(0, vdW-2).Draw(t, "x0")
		y0 := rapid.IntRange(0, vdH-2).Draw(t, "y0")
		x1 := rapid.IntRange(x0+1, vdW).Draw(t, "x1")
		y1 := rapid.IntRange(y0+1, vdH).Draw(t, "y1")
		subRect := image.Rect(x0, y0, x1, y1)

		surf := surface.NewFromSubImage(vd.FrameBuffer(), subRect)

		expectedBounds := image.Rect(0, 0, x1-x0, y1-y0)
		actualBounds := surf.Bounds()

		if actualBounds != expectedBounds {
			t.Fatalf("Surface bounds = %v, want %v (subRect=%v)", actualBounds, expectedBounds, subRect)
		}
	})
}

func TestProperty9_SurfaceDrawingIsClippedToRegionBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		vdW := rapid.IntRange(4, 500).Draw(t, "vdWidth")
		vdH := rapid.IntRange(4, 500).Draw(t, "vdHeight")

		vd, err := region2.NewVirtualDisplay(image.Rect(0, 0, vdW, vdH))
		if err != nil {
			t.Fatalf("unexpected error creating VD: %v", err)
		}

		x0 := rapid.IntRange(1, vdW-3).Draw(t, "x0")
		y0 := rapid.IntRange(1, vdH-3).Draw(t, "y0")
		x1 := rapid.IntRange(x0+1, vdW-1).Draw(t, "x1")
		y1 := rapid.IntRange(y0+1, vdH-1).Draw(t, "y1")
		subRect := image.Rect(x0, y0, x1, y1)

		surf := surface.NewFromSubImage(vd.FrameBuffer(), subRect)

		bgColor := color.RGBA{R: 42, G: 84, B: 126, A: 255}
		for y := 0; y < vdH; y++ {
			for x := 0; x < vdW; x++ {
				vd.FrameBuffer().SetRGBA(x, y, bgColor)
			}
		}

		testColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		surfW := subRect.Dx()
		surfH := subRect.Dy()

		outOfBoundsCoords := []image.Point{
			{X: -rapid.IntRange(1, 100).Draw(t, "negX"), Y: rapid.IntRange(0, surfH-1).Draw(t, "validY1")},
			{X: rapid.IntRange(0, surfW-1).Draw(t, "validX1"), Y: -rapid.IntRange(1, 100).Draw(t, "negY")},
			{X: surfW + rapid.IntRange(0, 100).Draw(t, "overX"), Y: rapid.IntRange(0, surfH-1).Draw(t, "validY2")},
			{X: rapid.IntRange(0, surfW-1).Draw(t, "validX2"), Y: surfH + rapid.IntRange(0, 100).Draw(t, "overY")},
		}

		for _, pt := range outOfBoundsCoords {
			surf.FrameBuffer().Set(pt.X, pt.Y, testColor)
		}

		for y := 0; y < vdH; y++ {
			for x := 0; x < vdW; x++ {
				pt := image.Pt(x, y)
				if pt.In(subRect) {
					continue
				}
				pixel := vd.FrameBuffer().RGBAAt(x, y)
				if pixel != bgColor {
					t.Fatalf("pixel at VD (%d,%d) outside region %v was modified: got %v, want %v",
						x, y, subRect, pixel, bgColor)
				}
			}
		}
	})
}

// =============================================================================
// From: texthints_property_test.go
// =============================================================================

// propertyTextHintTarget is a DrawTarget that implements TextHintProvider with configurable hints.
type propertyTextHintTarget struct {
	bounds image.Rectangle
	hints  textlayout.TextHints
}

func (pt *propertyTextHintTarget) Bounds() image.Rectangle         { return pt.bounds }
func (pt *propertyTextHintTarget) DrawImage(draw.Image) error      { return nil }
func (pt *propertyTextHintTarget) TextHints() textlayout.TextHints { return pt.hints }

// propertyPlainTarget is a DrawTarget that does NOT implement TextHintProvider.
type propertyPlainTarget struct {
	bounds image.Rectangle
}

func (pt *propertyPlainTarget) Bounds() image.Rectangle    { return pt.bounds }
func (pt *propertyPlainTarget) DrawImage(draw.Image) error { return nil }

func TestProperty14_SingleScreenWithTextHintProvider(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		screenW := rapid.IntRange(10, 500).Draw(t, "screenW")
		screenH := rapid.IntRange(10, 500).Draw(t, "screenH")
		screenBounds := image.Rect(0, 0, screenW, screenH)

		rx0 := rapid.IntRange(0, screenW-2).Draw(t, "rx0")
		ry0 := rapid.IntRange(0, screenH-2).Draw(t, "ry0")
		rx1 := rapid.IntRange(rx0+1, screenW).Draw(t, "rx1")
		ry1 := rapid.IntRange(ry0+1, screenH).Draw(t, "ry1")
		regionBounds := image.Rect(rx0, ry0, rx1, ry1)

		glyphW := rapid.IntRange(1, 20).Draw(t, "glyphW")
		glyphH := rapid.IntRange(1, 20).Draw(t, "glyphH")
		glyphAdv := rapid.IntRange(1, 25).Draw(t, "glyphAdv")
		rowH := rapid.IntRange(1, 30).Draw(t, "rowH")
		supportsVScroll := rapid.Bool().Draw(t, "supportsVScroll")
		supportsHScroll := rapid.Bool().Draw(t, "supportsHScroll")
		supportsAutoScroll := rapid.Bool().Draw(t, "supportsAutoScroll")
		preferEventRefresh := rapid.Bool().Draw(t, "preferEventRefresh")
		tickerDir := rapid.SampledFrom([]string{
			textlayout.TickerDirectionVertical,
			textlayout.TickerDirectionNone,
		}).Draw(t, "tickerDir")
		lineMode := rapid.SampledFrom([]string{
			textlayout.LineModeTruncate,
			textlayout.LineModeClip,
		}).Draw(t, "lineMode")

		providerHints := textlayout.TextHints{
			PixelWidth:               screenW,
			PixelHeight:              screenH,
			GlyphWidth:               glyphW,
			GlyphHeight:              glyphH,
			GlyphAdvance:             glyphAdv,
			RowHeight:                rowH,
			SupportsVerticalScroll:   supportsVScroll,
			SupportsHorizontalScroll: supportsHScroll,
			SupportsAutoScroll:       supportsAutoScroll,
			PreferEventRefresh:       preferEventRefresh,
			DefaultTickerDirection:   tickerDir,
			DefaultLineMode:          lineMode,
		}

		target := &propertyTextHintTarget{
			bounds: screenBounds,
			hints:  providerHints,
		}

		screens := []region2.ScreenPosition{
			{Index: 0, Name: "screen0", Bounds: screenBounds, Target: target, HintProvider: target.TextHints},
		}

		regionW := regionBounds.Dx()
		regionH := regionBounds.Dy()
		surf := surface.New(image.Rect(0, 0, regionW, regionH))
		r := region2.NewRegionWithScreens("test", regionBounds, surf, screens, "", 0, 0)

		hints := r.TextHints()

		if hints.PixelWidth != regionW {
			t.Fatalf("PixelWidth=%d, want %d (region width)", hints.PixelWidth, regionW)
		}
		if hints.PixelHeight != regionH {
			t.Fatalf("PixelHeight=%d, want %d (region height)", hints.PixelHeight, regionH)
		}
		if hints.GlyphWidth != glyphW {
			t.Fatalf("GlyphWidth=%d, want %d (from provider)", hints.GlyphWidth, glyphW)
		}
		if hints.GlyphHeight != glyphH {
			t.Fatalf("GlyphHeight=%d, want %d (from provider)", hints.GlyphHeight, glyphH)
		}
		if hints.GlyphAdvance != glyphAdv {
			t.Fatalf("GlyphAdvance=%d, want %d (from provider)", hints.GlyphAdvance, glyphAdv)
		}
		if hints.RowHeight != rowH {
			t.Fatalf("RowHeight=%d, want %d (from provider)", hints.RowHeight, rowH)
		}
		if hints.SupportsVerticalScroll != supportsVScroll {
			t.Fatalf("SupportsVerticalScroll=%v, want %v", hints.SupportsVerticalScroll, supportsVScroll)
		}
		if hints.SupportsHorizontalScroll != supportsHScroll {
			t.Fatalf("SupportsHorizontalScroll=%v, want %v", hints.SupportsHorizontalScroll, supportsHScroll)
		}
		if hints.SupportsAutoScroll != supportsAutoScroll {
			t.Fatalf("SupportsAutoScroll=%v, want %v", hints.SupportsAutoScroll, supportsAutoScroll)
		}
		if hints.PreferEventRefresh != preferEventRefresh {
			t.Fatalf("PreferEventRefresh=%v, want %v", hints.PreferEventRefresh, preferEventRefresh)
		}
		if hints.DefaultTickerDirection != tickerDir {
			t.Fatalf("DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, tickerDir)
		}
		if hints.DefaultLineMode != lineMode {
			t.Fatalf("DefaultLineMode=%q, want %q", hints.DefaultLineMode, lineMode)
		}
	})
}

func TestProperty14_SpansMultipleScreens(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		screen1W := rapid.IntRange(10, 300).Draw(t, "screen1W")
		screen1H := rapid.IntRange(10, 300).Draw(t, "screen1H")
		screen2W := rapid.IntRange(10, 300).Draw(t, "screen2W")
		screen2H := rapid.IntRange(10, 300).Draw(t, "screen2H")

		screen1Bounds := image.Rect(0, 0, screen1W, screen1H)
		screen2Bounds := image.Rect(screen1W, 0, screen1W+screen2W, screen2H)

		maxH := screen1H
		if screen2H < maxH {
			maxH = screen2H
		}
		rx0 := rapid.IntRange(0, screen1W-1).Draw(t, "rx0")
		rx1 := rapid.IntRange(screen1W+1, screen1W+screen2W).Draw(t, "rx1")
		ry0 := rapid.IntRange(0, maxH-2).Draw(t, "ry0")
		ry1 := rapid.IntRange(ry0+1, maxH).Draw(t, "ry1")
		regionBounds := image.Rect(rx0, ry0, rx1, ry1)

		target1 := &propertyTextHintTarget{
			bounds: screen1Bounds,
			hints: textlayout.TextHints{
				GlyphWidth:             8,
				GlyphHeight:            10,
				SupportsVerticalScroll: false,
			},
		}
		target2 := &propertyTextHintTarget{
			bounds: screen2Bounds,
			hints: textlayout.TextHints{
				GlyphWidth:               9,
				SupportsHorizontalScroll: false,
			},
		}

		screens := []region2.ScreenPosition{
			{Index: 0, Name: "left", Bounds: screen1Bounds, Target: target1, HintProvider: target1.TextHints},
			{Index: 1, Name: "right", Bounds: screen2Bounds, Target: target2, HintProvider: target2.TextHints},
		}

		regionW := regionBounds.Dx()
		regionH := regionBounds.Dy()
		surf := surface.New(image.Rect(0, 0, regionW, regionH))
		r := region2.NewRegionWithScreens("spanning", regionBounds, surf, screens, "", 0, 0)

		hints := r.TextHints()

		if hints.PixelWidth != regionW {
			t.Fatalf("PixelWidth=%d, want %d", hints.PixelWidth, regionW)
		}
		if hints.PixelHeight != regionH {
			t.Fatalf("PixelHeight=%d, want %d", hints.PixelHeight, regionH)
		}
		if hints.SupportsVerticalScroll != true {
			t.Fatal("SupportsVerticalScroll should be true (default for multi-screen)")
		}
		if hints.SupportsHorizontalScroll != true {
			t.Fatal("SupportsHorizontalScroll should be true (default for multi-screen)")
		}
		if hints.SupportsAutoScroll != true {
			t.Fatal("SupportsAutoScroll should be true (default for multi-screen)")
		}
		if hints.PreferEventRefresh != false {
			t.Fatal("PreferEventRefresh should be false (default for multi-screen)")
		}
		if hints.DefaultTickerDirection != textlayout.TickerDirectionVertical {
			t.Fatalf("DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionVertical)
		}
		if hints.DefaultLineMode != textlayout.LineModeTruncate {
			t.Fatalf("DefaultLineMode=%q, want %q", hints.DefaultLineMode, textlayout.LineModeTruncate)
		}
		expected := region2.TestBuildCatalogHints(textlayout.TextHints{
			PixelWidth:               regionW,
			PixelHeight:              regionH,
			GlyphWidth:               textlayout.GlyphWidth,
			GlyphHeight:              textlayout.GlyphHeight,
			GlyphAdvance:             textlayout.GlyphAdvance,
			RowHeight:                textlayout.RowHeight,
			SupportsVerticalScroll:   true,
			SupportsHorizontalScroll: true,
			SupportsAutoScroll:       true,
			PreferEventRefresh:       false,
			DefaultTickerDirection:   textlayout.TickerDirectionVertical,
			DefaultLineMode:          textlayout.LineModeTruncate,
		}, 96.0)
		if hints.GlyphWidth != expected.GlyphWidth {
			t.Fatalf("GlyphWidth=%d, want %d", hints.GlyphWidth, expected.GlyphWidth)
		}
		if hints.GlyphHeight != expected.GlyphHeight {
			t.Fatalf("GlyphHeight=%d, want %d", hints.GlyphHeight, expected.GlyphHeight)
		}
		if hints.GlyphAdvance != expected.GlyphAdvance {
			t.Fatalf("GlyphAdvance=%d, want %d", hints.GlyphAdvance, expected.GlyphAdvance)
		}
		if hints.RowHeight != expected.RowHeight {
			t.Fatalf("RowHeight=%d, want %d", hints.RowHeight, expected.RowHeight)
		}
	})
}

func TestProperty14_NoTextHintProvider(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		screenW := rapid.IntRange(10, 500).Draw(t, "screenW")
		screenH := rapid.IntRange(10, 500).Draw(t, "screenH")
		screenBounds := image.Rect(0, 0, screenW, screenH)

		rx0 := rapid.IntRange(0, screenW-2).Draw(t, "rx0")
		ry0 := rapid.IntRange(0, screenH-2).Draw(t, "ry0")
		rx1 := rapid.IntRange(rx0+1, screenW).Draw(t, "rx1")
		ry1 := rapid.IntRange(ry0+1, screenH).Draw(t, "ry1")
		regionBounds := image.Rect(rx0, ry0, rx1, ry1)

		target := &propertyPlainTarget{bounds: screenBounds}
		screens := []region2.ScreenPosition{
			{Index: 0, Name: "plain", Bounds: screenBounds, Target: target},
		}

		regionW := regionBounds.Dx()
		regionH := regionBounds.Dy()
		surf := surface.New(image.Rect(0, 0, regionW, regionH))
		r := region2.NewRegionWithScreens("noProvider", regionBounds, surf, screens, "", 0, 0)

		hints := r.TextHints()
		expected := region2.TestBuildCatalogHints(textlayout.DefaultTextHints(image.Rect(0, 0, regionW, regionH)), 96.0)

		if hints.PixelWidth != expected.PixelWidth {
			t.Fatalf("PixelWidth=%d, want %d", hints.PixelWidth, expected.PixelWidth)
		}
		if hints.PixelHeight != expected.PixelHeight {
			t.Fatalf("PixelHeight=%d, want %d", hints.PixelHeight, expected.PixelHeight)
		}
		if hints.GlyphWidth != expected.GlyphWidth {
			t.Fatalf("GlyphWidth=%d, want %d", hints.GlyphWidth, expected.GlyphWidth)
		}
		if hints.GlyphHeight != expected.GlyphHeight {
			t.Fatalf("GlyphHeight=%d, want %d", hints.GlyphHeight, expected.GlyphHeight)
		}
		if hints.GlyphAdvance != expected.GlyphAdvance {
			t.Fatalf("GlyphAdvance=%d, want %d", hints.GlyphAdvance, expected.GlyphAdvance)
		}
		if hints.RowHeight != expected.RowHeight {
			t.Fatalf("RowHeight=%d, want %d", hints.RowHeight, expected.RowHeight)
		}
		if hints.SupportsVerticalScroll != expected.SupportsVerticalScroll {
			t.Fatalf("SupportsVerticalScroll=%v, want %v", hints.SupportsVerticalScroll, expected.SupportsVerticalScroll)
		}
		if hints.SupportsHorizontalScroll != expected.SupportsHorizontalScroll {
			t.Fatalf("SupportsHorizontalScroll=%v, want %v", hints.SupportsHorizontalScroll, expected.SupportsHorizontalScroll)
		}
		if hints.SupportsAutoScroll != expected.SupportsAutoScroll {
			t.Fatalf("SupportsAutoScroll=%v, want %v", hints.SupportsAutoScroll, expected.SupportsAutoScroll)
		}
		if hints.PreferEventRefresh != expected.PreferEventRefresh {
			t.Fatalf("PreferEventRefresh=%v, want %v", hints.PreferEventRefresh, expected.PreferEventRefresh)
		}
		if hints.DefaultTickerDirection != expected.DefaultTickerDirection {
			t.Fatalf("DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, expected.DefaultTickerDirection)
		}
		if hints.DefaultLineMode != expected.DefaultLineMode {
			t.Fatalf("DefaultLineMode=%q, want %q", hints.DefaultLineMode, expected.DefaultLineMode)
		}
	})
}

func TestProperty14_PixelDimensionsAlwaysCorrect(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		scenario := rapid.IntRange(0, 2).Draw(t, "scenario")
		screenW := rapid.IntRange(10, 500).Draw(t, "screenW")
		screenH := rapid.IntRange(10, 500).Draw(t, "screenH")

		var regionBounds image.Rectangle
		var screens []region2.ScreenPosition

		switch scenario {
		case 0:
			screenBounds := image.Rect(0, 0, screenW, screenH)
			rx0 := rapid.IntRange(0, screenW-2).Draw(t, "rx0")
			ry0 := rapid.IntRange(0, screenH-2).Draw(t, "ry0")
			rx1 := rapid.IntRange(rx0+1, screenW).Draw(t, "rx1")
			ry1 := rapid.IntRange(ry0+1, screenH).Draw(t, "ry1")
			regionBounds = image.Rect(rx0, ry0, rx1, ry1)

			target := &propertyTextHintTarget{
				bounds: screenBounds,
				hints: textlayout.TextHints{
					GlyphWidth:             8,
					GlyphHeight:            10,
					GlyphAdvance:           9,
					RowHeight:              14,
					SupportsVerticalScroll: false,
				},
			}
			screens = []region2.ScreenPosition{
				{Index: 0, Name: "s0", Bounds: screenBounds, Target: target, HintProvider: target.TextHints},
			}

		case 1:
			screen2W := rapid.IntRange(10, 300).Draw(t, "screen2W")
			screen1Bounds := image.Rect(0, 0, screenW, screenH)
			screen2Bounds := image.Rect(screenW, 0, screenW+screen2W, screenH)

			rx0 := rapid.IntRange(0, screenW-1).Draw(t, "rx0")
			rx1 := rapid.IntRange(screenW+1, screenW+screen2W).Draw(t, "rx1")
			ry0 := rapid.IntRange(0, screenH-2).Draw(t, "ry0")
			ry1 := rapid.IntRange(ry0+1, screenH).Draw(t, "ry1")
			regionBounds = image.Rect(rx0, ry0, rx1, ry1)

			screens = []region2.ScreenPosition{
				{Index: 0, Name: "s0", Bounds: screen1Bounds, Target: &propertyPlainTarget{bounds: screen1Bounds}},
				{Index: 1, Name: "s1", Bounds: screen2Bounds, Target: &propertyPlainTarget{bounds: screen2Bounds}},
			}

		case 2:
			screenBounds := image.Rect(0, 0, screenW, screenH)
			rx0 := rapid.IntRange(0, screenW-2).Draw(t, "rx0")
			ry0 := rapid.IntRange(0, screenH-2).Draw(t, "ry0")
			rx1 := rapid.IntRange(rx0+1, screenW).Draw(t, "rx1")
			ry1 := rapid.IntRange(ry0+1, screenH).Draw(t, "ry1")
			regionBounds = image.Rect(rx0, ry0, rx1, ry1)

			screens = []region2.ScreenPosition{
				{Index: 0, Name: "s0", Bounds: screenBounds, Target: &propertyPlainTarget{bounds: screenBounds}},
			}
		}

		regionW := regionBounds.Dx()
		regionH := regionBounds.Dy()
		surf := surface.New(image.Rect(0, 0, regionW, regionH))
		r := region2.NewRegionWithScreens("dim-check", regionBounds, surf, screens, "", 0, 0)

		hints := r.TextHints()

		if hints.PixelWidth != regionBounds.Dx() {
			t.Fatalf("scenario %d: PixelWidth=%d, want %d",
				scenario, hints.PixelWidth, regionBounds.Dx())
		}
		if hints.PixelHeight != regionBounds.Dy() {
			t.Fatalf("scenario %d: PixelHeight=%d, want %d",
				scenario, hints.PixelHeight, regionBounds.Dy())
		}
	})
}

// =============================================================================
// From: tick_rate_property_test.go
// =============================================================================

func TestProperty4_TickIntervalDerivationCorrectness(t *testing.T) {
	resolver := &region2.DefaultTickRateResolver{}

	t.Run("registered_mode_with_positive_interval", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			intervalMS := rapid.IntRange(1, 10000).Draw(t, "intervalMS")
			modeID := fmt.Sprintf("__prop4_pos_%d", intervalMS)
			region2.RegisterTickRate(modeID, &propTestProvider{time.Duration(intervalMS) * time.Millisecond})

			interval := resolver.TickInterval(modeID)
			expected := time.Duration(intervalMS) * time.Millisecond

			if interval != expected {
				t.Fatalf("TickInterval(%q) with provider returning %dms: got %v, want %v",
					modeID, intervalMS, interval, expected)
			}
		})
	})

	t.Run("registered_mode_with_zero_or_negative_interval_clamped_to_min", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			intervalMS := rapid.IntRange(-1000, 0).Draw(t, "intervalMS")
			modeID := fmt.Sprintf("__prop4_neg_%d", intervalMS)
			region2.RegisterTickRate(modeID, &propTestProvider{time.Duration(intervalMS) * time.Millisecond})

			interval := resolver.TickInterval(modeID)

			if interval != 1*time.Millisecond {
				t.Fatalf("TickInterval(%q) with provider returning %dms: got %v, want 1ms (clamped min)",
					modeID, intervalMS, interval)
			}
		})
	})

	t.Run("non_registered_modes_return_default", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			modeID := rapid.OneOf(
				rapid.Just("dashboard"),
				rapid.Just("system"),
				rapid.Just("gpio"),
				rapid.Just("stemma"),
				rapid.Just("menu"),
				rapid.Just("clock"),
				rapid.Just("image"),
				rapid.Just("serial"),
				rapid.Just("usb"),
				rapid.Just(""),
				rapid.StringMatching(`[a-z]{1,20}`),
			).Filter(func(s string) bool {
				return s != "ticker"
			}).Draw(t, "modeID")

			// Only check modes that are not registered
			if region2.TestHasTickRateProvider(modeID) {
				return // skip — mode happens to be registered
			}

			interval := resolver.TickInterval(modeID)

			if interval != 1000*time.Millisecond {
				t.Fatalf("TickInterval(%q): got %v, want 1000ms (default for unregistered modes)",
					modeID, interval)
			}
		})
	})
}

// propTestProvider is a test helper implementing TickRateProvider with a fixed interval.
type propTestProvider struct {
	interval time.Duration
}

func (p *propTestProvider) PreferredTickInterval() time.Duration {
	return p.interval
}

func TestProperty5_TickIntervalBoundsInvariant(t *testing.T) {
	resolver := &region2.DefaultTickRateResolver{}

	rapid.Check(t, func(t *rapid.T) {
		modeID := rapid.String().Draw(t, "modeID")

		interval := resolver.TickInterval(modeID)

		if interval < 1*time.Millisecond {
			t.Fatalf("TickInterval(%q) = %v, want >= 1ms", modeID, interval)
		}
		if interval > 10000*time.Millisecond {
			t.Fatalf("TickInterval(%q) = %v, want <= 10000ms", modeID, interval)
		}
	})
}

// =============================================================================
// From: virtual_display_property_test.go
// =============================================================================

func TestProperty1_VirtualDisplayBoundsFromLayoutEqualsMinBoundingRect(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		numScreens := rapid.IntRange(1, 8).Draw(t, "numScreens")
		screens := make([]region2.ScreenPosition, numScreens)

		var expectedMaxX, expectedMaxY int
		for i := 0; i < numScreens; i++ {
			originX := rapid.IntRange(0, 1000).Draw(t, "originX")
			originY := rapid.IntRange(0, 1000).Draw(t, "originY")
			w := rapid.IntRange(1, 500).Draw(t, "width")
			h := rapid.IntRange(1, 500).Draw(t, "height")

			bounds := image.Rect(originX, originY, originX+w, originY+h)
			screens[i] = region2.ScreenPosition{
				Index:  i,
				Name:   rapid.StringMatching(`[a-zA-Z0-9]{1,16}`).Draw(t, "name"),
				Bounds: bounds,
			}

			if bounds.Max.X > expectedMaxX {
				expectedMaxX = bounds.Max.X
			}
			if bounds.Max.Y > expectedMaxY {
				expectedMaxY = bounds.Max.Y
			}
		}

		vd, err := region2.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedBounds := image.Rect(0, 0, expectedMaxX, expectedMaxY)
		if vd.Bounds() != expectedBounds {
			t.Fatalf("VD bounds = %v, want %v", vd.Bounds(), expectedBounds)
		}
	})
}

func TestProperty2_VirtualDisplaySingleScreenDefaultBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 500).Draw(t, "width")
		h := rapid.IntRange(1, 500).Draw(t, "height")

		screens := []region2.ScreenPosition{
			{
				Index:  0,
				Name:   rapid.StringMatching(`[a-zA-Z0-9]{1,16}`).Draw(t, "name"),
				Bounds: image.Rect(0, 0, w, h),
			},
		}

		vd, err := region2.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedBounds := image.Rect(0, 0, w, h)
		if vd.Bounds() != expectedBounds {
			t.Fatalf("VD bounds = %v, want %v (W=%d, H=%d)", vd.Bounds(), expectedBounds, w, h)
		}
	})
}

func TestProperty3_VirtualDisplayMultiScreenDefaultBounds(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		numScreens := rapid.IntRange(2, 8).Draw(t, "numScreens")

		type screenDef struct {
			w, h int
		}
		defs := make([]screenDef, numScreens)
		for i := 0; i < numScreens; i++ {
			defs[i] = screenDef{
				w: rapid.IntRange(1, 500).Draw(t, "width"),
				h: rapid.IntRange(1, 500).Draw(t, "height"),
			}
		}

		screens := make([]region2.ScreenPosition, numScreens)
		cumulativeX := 0
		totalWidth := 0
		maxHeight := 0

		for i := 0; i < numScreens; i++ {
			screens[i] = region2.ScreenPosition{
				Index:  i,
				Name:   rapid.StringMatching(`[a-zA-Z0-9]{1,16}`).Draw(t, "name"),
				Bounds: image.Rect(cumulativeX, 0, cumulativeX+defs[i].w, defs[i].h),
			}

			if screens[i].Bounds.Min.X != cumulativeX {
				t.Fatalf("screen %d X origin = %d, want %d", i, screens[i].Bounds.Min.X, cumulativeX)
			}

			cumulativeX += defs[i].w
			totalWidth += defs[i].w
			if defs[i].h > maxHeight {
				maxHeight = defs[i].h
			}
		}

		vd, err := region2.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if vd.Bounds().Dx() != totalWidth {
			t.Fatalf("VD width = %d, want %d (sum of screen widths)", vd.Bounds().Dx(), totalWidth)
		}

		if vd.Bounds().Dy() != maxHeight {
			t.Fatalf("VD height = %d, want %d (max screen height)", vd.Bounds().Dy(), maxHeight)
		}

		if vd.Bounds().Min != image.ZP {
			t.Fatalf("VD origin = %v, want (0,0)", vd.Bounds().Min)
		}

		expectedBounds := image.Rect(0, 0, totalWidth, maxHeight)
		if vd.Bounds() != expectedBounds {
			t.Fatalf("VD bounds = %v, want %v", vd.Bounds(), expectedBounds)
		}
	})
}

// =============================================================================
// From: tick_rate_prop_test.go
// =============================================================================

// genModeID generates random mode ID strings with a proptest prefix to avoid
// collisions with real registered modes.
func genModeID(t *rapid.T, label string) string {
	suffix := rapid.StringMatching(`[a-z][a-z0-9\-]{0,14}`).Draw(t, label)
	return "__proptest_" + suffix
}

// genTickInterval generates random time.Duration values including edge cases:
// 0, negative, very small, in-range, and very large.
func genTickInterval(t *rapid.T, label string) time.Duration {
	// Draw from a wide range including negative, zero, and very large values.
	// Range: -1s to 20s in nanoseconds covers all interesting cases.
	ns := rapid.Int64Range(-int64(time.Second), int64(20*time.Second)).Draw(t, label)
	return time.Duration(ns)
}

func TestProp_TickRateResolverDispatch(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		modeID := genModeID(t, "modeID")
		interval := genTickInterval(t, "interval")

		// Register a provider for this mode.
		provider := &propTestProvider{interval: interval}
		region2.RegisterTickRate(modeID, provider)

		resolver := &region2.DefaultTickRateResolver{}

		// Property 1: Registered providers return clamped interval.
		got := resolver.TickInterval(modeID)

		// Compute expected clamped value.
		expected := interval
		if expected < 1*time.Millisecond {
			expected = 1 * time.Millisecond
		}
		if expected > 10000*time.Millisecond {
			expected = 10000 * time.Millisecond
		}

		if got != expected {
			t.Fatalf("registered mode %q with interval %v: got %v, want %v (clamped)",
				modeID, interval, got, expected)
		}

		// Property 2: Unregistered mode returns default 1000ms.
		unregistered := genModeID(t, "unregisteredModeID")
		gotDefault := resolver.TickInterval(unregistered + "_unregistered")
		wantDefault := 1000 * time.Millisecond
		if gotDefault != wantDefault {
			t.Fatalf("unregistered mode %q: got %v, want %v",
				unregistered+"_unregistered", gotDefault, wantDefault)
		}

		// Property 3: Clamping boundary checks.
		// Values below 1ms → 1ms
		if interval < 1*time.Millisecond && got != 1*time.Millisecond {
			t.Fatalf("interval %v below min: got %v, want 1ms", interval, got)
		}
		// Values above 10000ms → 10000ms
		if interval > 10000*time.Millisecond && got != 10000*time.Millisecond {
			t.Fatalf("interval %v above max: got %v, want 10000ms", interval, got)
		}
		// Values in range → unchanged
		if interval >= 1*time.Millisecond && interval <= 10000*time.Millisecond && got != interval {
			t.Fatalf("interval %v in range: got %v, want %v", interval, got, interval)
		}
	})
}

// =============================================================================
// Properties 1-4 removed: padding no longer lives on Region (moved to LayoutBridge).
// See layout-padding-refactor spec for the new property tests.
// =============================================================================

// For any Region constructed with a screen providing specific TextHints capability flags
// and glyph metrics, the Region's TextHints SHALL preserve GlyphWidth, GlyphHeight,
// GlyphAdvance, RowHeight, SupportsVerticalScroll, SupportsHorizontalScroll,
// SupportsAutoScroll, PreferEventRefresh, SupportsRapidRefresh, SupportsColor,
// DefaultTickerDirection, and DefaultLineMode unchanged regardless of padding percentage.

func TestPaddingProperty5_TextHintsNonDimensionFieldsPreserved(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(10, 800).Draw(t, "width")
		h := rapid.IntRange(10, 600).Draw(t, "height")

		// Generate random capability flags and glyph metrics.
		glyphW := rapid.IntRange(3, 12).Draw(t, "glyphWidth")
		glyphH := rapid.IntRange(5, 16).Draw(t, "glyphHeight")
		glyphAdv := rapid.IntRange(4, 14).Draw(t, "glyphAdvance")
		rowH := rapid.IntRange(8, 20).Draw(t, "rowHeight")

		supportsVScroll := rapid.Bool().Draw(t, "supportsVerticalScroll")
		supportsHScroll := rapid.Bool().Draw(t, "supportsHorizontalScroll")
		supportsAutoScroll := rapid.Bool().Draw(t, "supportsAutoScroll")
		preferEventRefresh := rapid.Bool().Draw(t, "preferEventRefresh")
		capability := rapid.IntRange(0, 5).Draw(t, "capability")

		tickerDir := rapid.SampledFrom([]string{
			textlayout.TickerDirectionVertical,
			textlayout.TickerDirectionNone,
		}).Draw(t, "tickerDirection")
		lineMode := rapid.SampledFrom([]string{
			textlayout.LineModeTruncate,
			textlayout.LineModeClip,
		}).Draw(t, "lineMode")

		// Create a HintProvider that returns these specific hints.
		screenHints := textlayout.TextHints{
			PixelWidth:               w, // will be overridden by post-padding content area
			PixelHeight:              h, // will be overridden by post-padding content area
			GlyphWidth:               glyphW,
			GlyphHeight:              glyphH,
			GlyphAdvance:             glyphAdv,
			RowHeight:                rowH,
			SupportsVerticalScroll:   supportsVScroll,
			SupportsHorizontalScroll: supportsHScroll,
			SupportsAutoScroll:       supportsAutoScroll,
			PreferEventRefresh:       preferEventRefresh,
			Capability:               capability,
			DefaultTickerDirection:   tickerDir,
			DefaultLineMode:          lineMode,
		}

		bounds := image.Rect(0, 0, w, h)
		surf := surface.New(bounds)

		// Create a ScreenPosition whose bounds contain the region and has a HintProvider.
		screen := region2.ScreenPosition{
			Index:  0,
			Name:   "test-screen",
			Bounds: bounds,
			HintProvider: func() textlayout.TextHints {
				return screenHints
			},
		}

		r := region2.NewRegionWithScreens("test", bounds, surf, []region2.ScreenPosition{screen}, "", 0, 0)
		hints := r.TextHints()

		// Verify non-dimension fields are preserved.
		if hints.GlyphWidth != glyphW {
			t.Fatalf("GlyphWidth=%d, want %d", hints.GlyphWidth, glyphW)
		}
		if hints.GlyphHeight != glyphH {
			t.Fatalf("GlyphHeight=%d, want %d", hints.GlyphHeight, glyphH)
		}
		if hints.GlyphAdvance != glyphAdv {
			t.Fatalf("GlyphAdvance=%d, want %d", hints.GlyphAdvance, glyphAdv)
		}
		if hints.RowHeight != rowH {
			t.Fatalf("RowHeight=%d, want %d", hints.RowHeight, rowH)
		}
		if hints.SupportsVerticalScroll != supportsVScroll {
			t.Fatalf("SupportsVerticalScroll=%v, want %v", hints.SupportsVerticalScroll, supportsVScroll)
		}
		if hints.SupportsHorizontalScroll != supportsHScroll {
			t.Fatalf("SupportsHorizontalScroll=%v, want %v", hints.SupportsHorizontalScroll, supportsHScroll)
		}
		if hints.SupportsAutoScroll != supportsAutoScroll {
			t.Fatalf("SupportsAutoScroll=%v, want %v", hints.SupportsAutoScroll, supportsAutoScroll)
		}
		if hints.PreferEventRefresh != preferEventRefresh {
			t.Fatalf("PreferEventRefresh=%v, want %v", hints.PreferEventRefresh, preferEventRefresh)
		}
		if hints.Capability != capability {
			t.Fatalf("Capability=%v, want %v", hints.Capability, capability)
		}
		if hints.DefaultTickerDirection != tickerDir {
			t.Fatalf("DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, tickerDir)
		}
		if hints.DefaultLineMode != lineMode {
			t.Fatalf("DefaultLineMode=%q, want %q", hints.DefaultLineMode, lineMode)
		}
	})
}
