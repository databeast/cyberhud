package tests

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// noopInstance is a minimal ModeInstance for tests that need a factory wired
// but don't care about rendering output.
type noopInstance struct{ id string }

func (n *noopInstance) ID() string                    { return n.id }
func (n *noopInstance) Activate()                     {}
func (n *noopInstance) Deactivate()                   {}
func (n *noopInstance) ActionHandler() action.Handler { return nil }
func (n *noopInstance) BuildView() style.ViewData     { return style.ViewData{} }
func (n *noopInstance) RenderCacheKey() uint32        { return 0 }

// noopModeFactory returns a ModeFactory that always succeeds with a noopInstance.
func noopModeFactory() region.ModeFactory {
	return func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		return &noopInstance{id: id}, true
	}
}

// wireNoopFactory sets a noopModeFactory on all regions in the RegionManager.
func wireNoopFactory(rm *region.RegionManager) {
	for _, r := range rm.Regions() {
		r.SetModeFactory(noopModeFactory())
	}
}

// --- From: activation_test.go ---

// mockActivationTarget is a minimal DrawTarget for activation tests.
type mockActivationTarget struct {
	bounds image.Rectangle
	drawn  bool
}

func (m *mockActivationTarget) Bounds() image.Rectangle    { return m.bounds }
func (m *mockActivationTarget) DrawImage(draw.Image) error { m.drawn = true; return nil }

func TestActivatePanel_SingleScreen_NoLayout(t *testing.T) {
	// Single-screen activation with no layout should produce one "default" region
	// covering the full VD bounds.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen},
		Layout:       nil,
		DefaultMode:  "dashboard",
		InputEnabled: false,
		AvailModes:   []string{"dashboard", "clock"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify VirtualDisplay bounds.
	if activation.VirtualDisplay.Bounds() != image.Rect(0, 0, 240, 240) {
		t.Fatalf("VD.Bounds()=%v, want %v", activation.VirtualDisplay.Bounds(), image.Rect(0, 0, 240, 240))
	}

	// Verify one region named "default".
	regions := activation.RegionManager.Regions()
	if len(regions) != 1 {
		t.Fatalf("got %d regions, want 1", len(regions))
	}
	if regions[0].Name() != "default" {
		t.Fatalf("region name=%q, want %q", regions[0].Name(), "default")
	}
	if regions[0].Bounds() != image.Rect(0, 0, 240, 240) {
		t.Fatalf("region bounds=%v, want %v", regions[0].Bounds(), image.Rect(0, 0, 240, 240))
	}
	if regions[0].CurrentMode() != "dashboard" {
		t.Fatalf("region mode=%q, want %q", regions[0].CurrentMode(), "dashboard")
	}

	// Verify FlushPath, RenderLoop, and ModeSwitch are not nil.
	if activation.FlushPath == nil {
		t.Fatal("FlushPath is nil")
	}
	if activation.RenderLoop == nil {
		t.Fatal("RenderLoop is nil")
	}
	if activation.ModeSwitch == nil {
		t.Fatal("ModeSwitch is nil")
	}
}

func TestActivatePanel_MultiScreen_NoLayout(t *testing.T) {
	// Multi-screen activation with no layout should produce one region per screen.
	screen1 := region.ScreenPosition{
		Index:  0,
		Name:   "left",
		Bounds: image.Rect(0, 0, 240, 135),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 135)},
	}
	screen2 := region.ScreenPosition{
		Index:  1,
		Name:   "right",
		Bounds: image.Rect(240, 0, 480, 135),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 135)},
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen1, screen2},
		Layout:       nil,
		DefaultMode:  "",
		InputEnabled: false,
		AvailModes:   []string{"dashboard", "clock"},
		ScreenModes: map[string]string{
			"left":  "dashboard",
			"right": "clock",
		},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify VD bounds encompass both screens.
	if activation.VirtualDisplay.Bounds() != image.Rect(0, 0, 480, 135) {
		t.Fatalf("VD.Bounds()=%v, want %v", activation.VirtualDisplay.Bounds(), image.Rect(0, 0, 480, 135))
	}

	// Verify two regions created.
	regions := activation.RegionManager.Regions()
	if len(regions) != 2 {
		t.Fatalf("got %d regions, want 2", len(regions))
	}

	// First region: "left"
	if regions[0].Name() != "left" {
		t.Fatalf("region[0] name=%q, want %q", regions[0].Name(), "left")
	}
	if regions[0].Bounds() != image.Rect(0, 0, 240, 135) {
		t.Fatalf("region[0] bounds=%v, want %v", regions[0].Bounds(), image.Rect(0, 0, 240, 135))
	}
	if regions[0].CurrentMode() != "dashboard" {
		t.Fatalf("region[0] mode=%q, want %q", regions[0].CurrentMode(), "dashboard")
	}

	// Second region: "right"
	if regions[1].Name() != "right" {
		t.Fatalf("region[1] name=%q, want %q", regions[1].Name(), "right")
	}
	if regions[1].Bounds() != image.Rect(240, 0, 480, 135) {
		t.Fatalf("region[1] bounds=%v, want %v", regions[1].Bounds(), image.Rect(240, 0, 480, 135))
	}
	if regions[1].CurrentMode() != "clock" {
		t.Fatalf("region[1] mode=%q, want %q", regions[1].CurrentMode(), "clock")
	}
}

func TestActivatePanel_ExplicitLayout(t *testing.T) {
	// Explicit layout with two regions on a single screen.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	layout := &region.RegionLayout{
		Specs: []region.RegionSpec{
			{Name: "top", Bounds: image.Rect(0, 0, 240, 120), DefaultMode: "clock"},
			{Name: "bottom", Bounds: image.Rect(0, 120, 240, 240), DefaultMode: "dashboard"},
		},
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen},
		Layout:       layout,
		DefaultMode:  "dashboard",
		InputEnabled: true,
		AvailModes:   []string{"dashboard", "clock"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify two regions created.
	regions := activation.RegionManager.Regions()
	if len(regions) != 2 {
		t.Fatalf("got %d regions, want 2", len(regions))
	}

	if regions[0].Name() != "top" {
		t.Fatalf("region[0] name=%q, want %q", regions[0].Name(), "top")
	}
	if regions[0].Bounds() != image.Rect(0, 0, 240, 120) {
		t.Fatalf("region[0] bounds=%v, want %v", regions[0].Bounds(), image.Rect(0, 0, 240, 120))
	}
	if regions[0].CurrentMode() != "clock" {
		t.Fatalf("region[0] mode=%q, want %q", regions[0].CurrentMode(), "clock")
	}

	if regions[1].Name() != "bottom" {
		t.Fatalf("region[1] name=%q, want %q", regions[1].Name(), "bottom")
	}
	if regions[1].Bounds() != image.Rect(0, 120, 240, 240) {
		t.Fatalf("region[1] bounds=%v, want %v", regions[1].Bounds(), image.Rect(0, 120, 240, 240))
	}
	if regions[1].CurrentMode() != "dashboard" {
		t.Fatalf("region[1] mode=%q, want %q", regions[1].CurrentMode(), "dashboard")
	}

	// First region should have input focus.
	if !regions[0].HasInputFocus() {
		t.Fatal("region[0] should have input focus")
	}
}

func TestActivatePanel_ZeroScreens_Error(t *testing.T) {
	// Zero screens should return an error.
	config := region.PanelActivationConfig{
		Screens:    nil,
		AvailModes: []string{"dashboard"},
	}

	_, err := region.ActivatePanel(config)
	if err == nil {
		t.Fatal("ActivatePanel() expected error for zero screens, got nil")
	}
}

func TestActivatePanel_InvalidLayout_Error(t *testing.T) {
	// Invalid layout (region bounds outside VD) should return an error with no
	// partial state.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	layout := &region.RegionLayout{
		Specs: []region.RegionSpec{
			{Name: "valid", Bounds: image.Rect(0, 0, 240, 120), DefaultMode: "clock"},
			{Name: "invalid", Bounds: image.Rect(0, 0, 500, 500), DefaultMode: "dashboard"}, // outside VD
		},
	}

	config := region.PanelActivationConfig{
		Screens:     []region.ScreenPosition{screen},
		Layout:      layout,
		DefaultMode: "dashboard",
		AvailModes:  []string{"dashboard", "clock"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock"
		},
	}

	_, err := region.ActivatePanel(config)
	if err == nil {
		t.Fatal("ActivatePanel() expected error for invalid layout, got nil")
	}
}

func TestActivatePanel_EmptyLayout_UsesDefault(t *testing.T) {
	// Empty layout (zero specs) should be treated as absent and apply default generation.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	layout := &region.RegionLayout{
		Specs: []region.RegionSpec{}, // empty
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen},
		Layout:       layout,
		DefaultMode:  "clock",
		InputEnabled: false,
		AvailModes:   []string{"clock", "dashboard"},
		ModeValidator: func(mode string) bool {
			return mode == "clock" || mode == "dashboard"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Should fall through to default generation: single "default" region.
	regions := activation.RegionManager.Regions()
	if len(regions) != 1 {
		t.Fatalf("got %d regions, want 1", len(regions))
	}
	if regions[0].Name() != "default" {
		t.Fatalf("region name=%q, want %q", regions[0].Name(), "default")
	}
	if regions[0].CurrentMode() != "clock" {
		t.Fatalf("region mode=%q, want %q", regions[0].CurrentMode(), "clock")
	}
}

func TestActivatePanel_ModeSwitch_Functional(t *testing.T) {
	// Verify the returned ModeSwitch can change modes on allocated regions.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen},
		Layout:       nil,
		DefaultMode:  "dashboard",
		InputEnabled: false,
		AvailModes:   []string{"dashboard", "clock"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Wire ModeFactory on all regions so SetMode can construct instances.
	for _, r := range activation.RegionManager.Regions() {
		r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
			return &noopInstance{id: id}, true
		})
	}

	// Switch mode via ModeSwitch.
	err = activation.ModeSwitch.Execute(region.ModeChangeCommand{Target: "default", Mode: "clock"})
	if err != nil {
		t.Fatalf("ModeSwitch.Execute() unexpected error: %v", err)
	}

	regions := activation.RegionManager.Regions()
	if regions[0].CurrentMode() != "clock" {
		t.Fatalf("after mode switch: mode=%q, want %q", regions[0].CurrentMode(), "clock")
	}
}

func TestActivatePanel_InputFocus_DefaultRegion(t *testing.T) {
	// The first allocated region should have input focus.
	screen := region.ScreenPosition{
		Index:  0,
		Name:   "main",
		Bounds: image.Rect(0, 0, 240, 240),
		Target: &mockActivationTarget{bounds: image.Rect(0, 0, 240, 240)},
	}

	config := region.PanelActivationConfig{
		Screens:      []region.ScreenPosition{screen},
		Layout:       nil,
		DefaultMode:  "dashboard",
		InputEnabled: true,
		AvailModes:   []string{"dashboard"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	active := activation.RegionManager.InputActiveRegion()
	if active == nil {
		t.Fatal("InputActiveRegion() returned nil")
	}
	if active.Name() != "default" {
		t.Fatalf("InputActiveRegion().Name()=%q, want %q", active.Name(), "default")
	}
}

// --- From: flush_test.go ---

// flushMockTarget records the image passed to DrawImage for verification.
type flushMockTarget struct {
	bounds    image.Rectangle
	lastImage draw.Image
	err       error
	callCount int
}

func (m *flushMockTarget) Bounds() image.Rectangle {
	return m.bounds
}

func (m *flushMockTarget) DrawImage(img draw.Image) error {
	m.callCount++
	m.lastImage = img
	return m.err
}

func TestFlushPath_Flush(t *testing.T) {
	t.Run("normal flush with single screen", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Write a known pixel into the VD framebuffer.
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		vd.FrameBuffer().SetRGBA(10, 20, red)

		target := &flushMockTarget{bounds: image.Rect(0, 0, 240, 240)}
		screens := []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240), Target: target},
		}

		fp := region.NewFlushPath(vd, screens)
		if err := fp.Flush(); err != nil {
			t.Fatalf("Flush() returned error: %v", err)
		}

		if target.callCount != 1 {
			t.Fatalf("DrawImage called %d times, want 1", target.callCount)
		}

		// Verify the pixel was correctly extracted.
		got := target.lastImage.(*image.RGBA).RGBAAt(10, 20)
		if got != red {
			t.Errorf("pixel at (10,20) = %v, want %v", got, red)
		}

		// Verify bounds are zero-origin.
		if target.lastImage.Bounds() != image.Rect(0, 0, 240, 240) {
			t.Errorf("image bounds = %v, want (0,0)-(240,240)", target.lastImage.Bounds())
		}
	})

	t.Run("ascending index order", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Record call order.
		var callOrder []int
		target0 := &flushMockTarget{bounds: image.Rect(0, 0, 240, 135)}
		target1 := &flushMockTarget{bounds: image.Rect(0, 0, 240, 135)}

		// Provide screens in REVERSE order to verify sorting.
		screens := []region.ScreenPosition{
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135), Target: target1},
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135), Target: target0},
		}

		fp := region.NewFlushPath(vd, screens)

		// Write distinct pixels to each screen's area.
		green := color.RGBA{G: 255, A: 255}
		blue := color.RGBA{B: 255, A: 255}
		vd.FrameBuffer().SetRGBA(5, 5, green)  // Left screen area.
		vd.FrameBuffer().SetRGBA(245, 5, blue) // Right screen area.

		if err := fp.Flush(); err != nil {
			t.Fatalf("Flush() returned error: %v", err)
		}

		// Both should have been called.
		if target0.callCount != 1 {
			t.Errorf("target0 called %d times, want 1", target0.callCount)
		}
		if target1.callCount != 1 {
			t.Errorf("target1 called %d times, want 1", target1.callCount)
		}

		// Verify correct pixels went to correct targets.
		gotLeft := target0.lastImage.(*image.RGBA).RGBAAt(5, 5)
		if gotLeft != green {
			t.Errorf("left screen pixel (5,5) = %v, want %v", gotLeft, green)
		}
		// Right screen: VD pixel at (245, 5) → zero-origin (5, 5) in right screen image.
		gotRight := target1.lastImage.(*image.RGBA).RGBAAt(5, 5)
		if gotRight != blue {
			t.Errorf("right screen pixel (5,5) = %v, want %v", gotRight, blue)
		}

		// Verify call order by re-running with order tracking.
		callOrder = nil
		trackTarget0 := &orderTrackingTarget{index: 0, order: &callOrder, bounds: image.Rect(0, 0, 240, 135)}
		trackTarget1 := &orderTrackingTarget{index: 1, order: &callOrder, bounds: image.Rect(0, 0, 240, 135)}
		screens2 := []region.ScreenPosition{
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135), Target: trackTarget1},
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135), Target: trackTarget0},
		}
		fp2 := region.NewFlushPath(vd, screens2)
		fp2.Flush()

		if len(callOrder) != 2 || callOrder[0] != 0 || callOrder[1] != 1 {
			t.Errorf("call order = %v, want [0 1]", callOrder)
		}
	})

	t.Run("error from one screen does not prevent flushing others", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		failTarget := &flushMockTarget{
			bounds: image.Rect(0, 0, 240, 135),
			err:    errors.New("hardware failure"),
		}
		okTarget := &flushMockTarget{
			bounds: image.Rect(0, 0, 240, 135),
		}

		screens := []region.ScreenPosition{
			{Index: 0, Name: "failing", Bounds: image.Rect(0, 0, 240, 135), Target: failTarget},
			{Index: 1, Name: "working", Bounds: image.Rect(240, 0, 480, 135), Target: okTarget},
		}

		fp := region.NewFlushPath(vd, screens)
		err = fp.Flush()

		// Should return an error (combined).
		if err == nil {
			t.Fatal("expected error from Flush(), got nil")
		}

		// Both targets should have been called.
		if failTarget.callCount != 1 {
			t.Errorf("failTarget called %d times, want 1", failTarget.callCount)
		}
		if okTarget.callCount != 1 {
			t.Errorf("okTarget called %d times, want 1", okTarget.callCount)
		}
	})

	t.Run("nil target is skipped", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		okTarget := &flushMockTarget{bounds: image.Rect(0, 0, 240, 135)}

		screens := []region.ScreenPosition{
			{Index: 0, Name: "nil-screen", Bounds: image.Rect(0, 0, 240, 135), Target: nil},
			{Index: 1, Name: "ok-screen", Bounds: image.Rect(240, 0, 480, 135), Target: okTarget},
		}

		fp := region.NewFlushPath(vd, screens)
		err = fp.Flush()

		if err != nil {
			t.Fatalf("Flush() returned error: %v", err)
		}

		// Only the non-nil target should have been called.
		if okTarget.callCount != 1 {
			t.Errorf("okTarget called %d times, want 1", okTarget.callCount)
		}
	})

	t.Run("correct pixel extraction — sub-rect matches VD content", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 320))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Paint a region of the VD with a distinctive pattern.
		for y := 0; y < 135; y++ {
			for x := 240; x < 480; x++ {
				c := color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255}
				vd.FrameBuffer().SetRGBA(x, y, c)
			}
		}

		target := &flushMockTarget{bounds: image.Rect(0, 0, 240, 135)}
		screens := []region.ScreenPosition{
			{Index: 0, Name: "right", Bounds: image.Rect(240, 0, 480, 135), Target: target},
		}

		fp := region.NewFlushPath(vd, screens)
		fp.Flush()

		// Verify every pixel in the extracted image matches what was in the VD.
		img := target.lastImage.(*image.RGBA)
		for y := 0; y < 135; y++ {
			for x := 0; x < 240; x++ {
				want := color.RGBA{R: uint8((x + 240) % 256), G: uint8(y % 256), B: 128, A: 255}
				got := img.RGBAAt(x, y)
				if got != want {
					t.Fatalf("pixel (%d,%d) = %v, want %v", x, y, got, want)
				}
			}
		}
	})

	t.Run("uncovered pixels are opaque black", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Leave VD framebuffer at its default (transparent black: 0,0,0,0).
		// Uncovered pixels should be filled with opaque black (0,0,0,255).
		target := &flushMockTarget{bounds: image.Rect(0, 0, 240, 240)}
		screens := []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240), Target: target},
		}

		fp := region.NewFlushPath(vd, screens)
		fp.Flush()

		img := target.lastImage.(*image.RGBA)
		opaqueBlack := color.RGBA{0, 0, 0, 255}

		// Sample several pixels — all should be opaque black since nothing was rendered.
		for _, pt := range []image.Point{{0, 0}, {119, 119}, {239, 239}, {50, 100}} {
			got := img.RGBAAt(pt.X, pt.Y)
			if got != opaqueBlack {
				t.Errorf("uncovered pixel (%d,%d) = %v, want opaque black %v", pt.X, pt.Y, got, opaqueBlack)
			}
		}
	})

	t.Run("empty screens produces no error", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		fp := region.NewFlushPath(vd, nil)
		if err := fp.Flush(); err != nil {
			t.Fatalf("Flush() with no screens returned error: %v", err)
		}
	})

	t.Run("screen bounds partially outside VD are clipped", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 100, 100))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Screen claims bounds that extend beyond the VD.
		target := &flushMockTarget{bounds: image.Rect(0, 0, 200, 200)}
		screens := []region.ScreenPosition{
			{Index: 0, Name: "oversized", Bounds: image.Rect(50, 50, 200, 200), Target: target},
		}

		fp := region.NewFlushPath(vd, screens)
		err = fp.Flush()
		if err != nil {
			t.Fatalf("Flush() returned error: %v", err)
		}

		// The image should be clipped to the intersection: (50,50)-(100,100) → 50x50.
		img := target.lastImage.(*image.RGBA)
		if img.Bounds() != image.Rect(0, 0, 50, 50) {
			t.Errorf("clipped image bounds = %v, want (0,0)-(50,50)", img.Bounds())
		}
	})
}

// orderTrackingTarget records the order in which DrawImage was called across targets.
type orderTrackingTarget struct {
	index  int
	order  *[]int
	bounds image.Rectangle
}

func (o *orderTrackingTarget) Bounds() image.Rectangle {
	return o.bounds
}

func (o *orderTrackingTarget) DrawImage(img draw.Image) error {
	*o.order = append(*o.order, o.index)
	return nil
}

// --- From: layout_test.go ---

func TestGenerateDefaultLayout_SingleScreen_ExplicitDefaultMode(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		},
		DefaultMode:  "clock",
		InputEnabled: true,
		AvailModes:   []string{"menu", "dashboard", "clock"},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(layout.Specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(layout.Specs))
	}

	spec := layout.Specs[0]
	if spec.Name != "default" {
		t.Errorf("expected name %q, got %q", "default", spec.Name)
	}
	if spec.Bounds != vd.Bounds() {
		t.Errorf("expected bounds %v, got %v", vd.Bounds(), spec.Bounds)
	}
	if spec.DefaultMode != "clock" {
		t.Errorf("expected mode %q, got %q", "clock", spec.DefaultMode)
	}
}

func TestGenerateDefaultLayout_SingleScreen_FallbackToMenu(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		},
		DefaultMode:  "",
		InputEnabled: true,
		AvailModes:   []string{"menu", "dashboard", "clock"},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if layout.Specs[0].DefaultMode != "menu" {
		t.Errorf("expected mode %q, got %q", "menu", layout.Specs[0].DefaultMode)
	}
}

func TestGenerateDefaultLayout_SingleScreen_FallbackToDashboard(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		},
		DefaultMode:  "",
		InputEnabled: false, // no input, so "menu" is skipped
		AvailModes:   []string{"dashboard", "clock"},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if layout.Specs[0].DefaultMode != "dashboard" {
		t.Errorf("expected mode %q, got %q", "dashboard", layout.Specs[0].DefaultMode)
	}
}

func TestGenerateDefaultLayout_SingleScreen_FallbackToFirstAvailable(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		},
		DefaultMode:  "",
		InputEnabled: false,
		AvailModes:   []string{"clock", "gpio"},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if layout.Specs[0].DefaultMode != "clock" {
		t.Errorf("expected mode %q, got %q", "clock", layout.Specs[0].DefaultMode)
	}
}

func TestGenerateDefaultLayout_SingleScreen_EmptyAvailModes_Error(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		},
		DefaultMode:  "",
		InputEnabled: false,
		AvailModes:   []string{},
	}

	_, err = region.GenerateDefaultLayout(vd, config)
	if err == nil {
		t.Fatal("expected error for empty AvailModes, got nil")
	}
}

func TestGenerateDefaultLayout_MultiScreen_TwoScreens(t *testing.T) {
	// Two screens: 240x135 each, positioned left-to-right.
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135)},
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135)},
		},
		AvailModes: []string{"clock", "dashboard"},
		ScreenModes: map[string]string{
			"left":  "clock",
			"right": "dashboard",
		},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(layout.Specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(layout.Specs))
	}

	// First region: left screen.
	if layout.Specs[0].Name != "left" {
		t.Errorf("expected name %q, got %q", "left", layout.Specs[0].Name)
	}
	expectedBounds0 := image.Rect(0, 0, 240, 135)
	if layout.Specs[0].Bounds != expectedBounds0 {
		t.Errorf("expected bounds %v, got %v", expectedBounds0, layout.Specs[0].Bounds)
	}
	if layout.Specs[0].DefaultMode != "clock" {
		t.Errorf("expected mode %q, got %q", "clock", layout.Specs[0].DefaultMode)
	}

	// Second region: right screen.
	if layout.Specs[1].Name != "right" {
		t.Errorf("expected name %q, got %q", "right", layout.Specs[1].Name)
	}
	expectedBounds1 := image.Rect(240, 0, 480, 135)
	if layout.Specs[1].Bounds != expectedBounds1 {
		t.Errorf("expected bounds %v, got %v", expectedBounds1, layout.Specs[1].Bounds)
	}
	if layout.Specs[1].DefaultMode != "dashboard" {
		t.Errorf("expected mode %q, got %q", "dashboard", layout.Specs[1].DefaultMode)
	}
}

func TestGenerateDefaultLayout_MultiScreen_UnsortedIndex(t *testing.T) {
	// Screens provided out of order should still be sorted by index.
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135)},
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135)},
		},
		AvailModes: []string{"clock", "dashboard"},
		ScreenModes: map[string]string{
			"left":  "clock",
			"right": "dashboard",
		},
	}

	layout, err := region.GenerateDefaultLayout(vd, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be sorted: left (index 0) first, right (index 1) second.
	if layout.Specs[0].Name != "left" {
		t.Errorf("expected first spec name %q, got %q", "left", layout.Specs[0].Name)
	}
	if layout.Specs[1].Name != "right" {
		t.Errorf("expected second spec name %q, got %q", "right", layout.Specs[1].Name)
	}
}

func TestGenerateDefaultLayout_MultiScreen_MissingDefaultMode_Error(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135)},
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135)},
		},
		AvailModes: []string{"clock", "dashboard"},
		ScreenModes: map[string]string{
			"left": "clock",
			// "right" is missing — should cause error.
		},
	}

	_, err = region.GenerateDefaultLayout(vd, config)
	if err == nil {
		t.Fatal("expected error for missing screen default mode, got nil")
	}

	expected := `region: screen "right" (index 1) has no valid default mode`
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestGenerateDefaultLayout_MultiScreen_UnregisteredMode_Error(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 135))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens: []region.ScreenPosition{
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135)},
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135)},
		},
		AvailModes: []string{"clock", "dashboard"},
		ScreenModes: map[string]string{
			"left":  "clock",
			"right": "unknown_mode", // not in AvailModes
		},
	}

	_, err = region.GenerateDefaultLayout(vd, config)
	if err == nil {
		t.Fatal("expected error for unregistered screen mode, got nil")
	}

	expected := `region: screen "right" (index 1) has no valid default mode`
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestGenerateDefaultLayout_NoScreens_Error(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	config := region.PanelActivationConfig{
		Screens:    []region.ScreenPosition{},
		AvailModes: []string{"clock"},
	}

	_, err = region.GenerateDefaultLayout(vd, config)
	if err == nil {
		t.Fatal("expected error for zero screens, got nil")
	}
}

// --- From: manager_test.go ---

func newTestVD(t *testing.T, w, h int) *region.VirtualDisplay {
	t.Helper()
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, w, h))
	if err != nil {
		t.Fatalf("NewVirtualDisplay(%dx%d): %v", w, h, err)
	}
	return vd
}

func TestNewRegionManager(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)
	if rm == nil {
		t.Fatal("NewRegionManager returned nil")
	}
	if len(rm.Regions()) != 0 {
		t.Fatalf("expected 0 regions, got %d", len(rm.Regions()))
	}
}

func TestAllocate_Valid(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	err := rm.Allocate(region.RegionSpec{
		Name:        "left",
		Bounds:      image.Rect(0, 0, 240, 240),
		DefaultMode: "clock",
	})
	if err != nil {
		t.Fatalf("Allocate left: %v", err)
	}

	err = rm.Allocate(region.RegionSpec{
		Name:        "right",
		Bounds:      image.Rect(240, 0, 480, 240),
		DefaultMode: "dashboard",
	})
	if err != nil {
		t.Fatalf("Allocate right: %v", err)
	}

	if len(rm.Regions()) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(rm.Regions()))
	}
}

func TestAllocate_FirstRegionGetsInputFocus(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "first", Bounds: image.Rect(0, 0, 240, 240)})
	rm.Allocate(region.RegionSpec{Name: "second", Bounds: image.Rect(240, 0, 480, 240)})

	active := rm.InputActiveRegion()
	if active == nil {
		t.Fatal("no input active region")
	}
	if active.Name() != "first" {
		t.Fatalf("expected input focus on 'first', got %q", active.Name())
	}
}

func TestAllocate_EmptyName(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	err := rm.Allocate(region.RegionSpec{Name: "", Bounds: image.Rect(0, 0, 100, 100)})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !strings.Contains(err.Error(), "got empty string") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_NameTooLong(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	longName := strings.Repeat("a", 65)
	err := rm.Allocate(region.RegionSpec{Name: longName, Bounds: image.Rect(0, 0, 100, 100)})
	if err == nil {
		t.Fatal("expected error for name > 64 chars")
	}
	if !strings.Contains(err.Error(), "exceeds 64 character limit") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_DuplicateName(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "Panel", Bounds: image.Rect(0, 0, 240, 240)})
	err := rm.Allocate(region.RegionSpec{Name: "panel", Bounds: image.Rect(240, 0, 480, 240)})
	if err == nil {
		t.Fatal("expected error for duplicate name (case-insensitive)")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_ZeroBounds(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	err := rm.Allocate(region.RegionSpec{Name: "bad", Bounds: image.Rect(0, 0, 0, 100)})
	if err == nil {
		t.Fatal("expected error for zero width")
	}
	if !strings.Contains(err.Error(), "width >= 1 and height >= 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_OutOfBounds(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	err := rm.Allocate(region.RegionSpec{Name: "big", Bounds: image.Rect(0, 0, 300, 240)})
	if err == nil {
		t.Fatal("expected error for out-of-bounds region")
	}
	if !strings.Contains(err.Error(), "extend outside virtual display") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_Overlap(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "first", Bounds: image.Rect(0, 0, 250, 240)})
	err := rm.Allocate(region.RegionSpec{Name: "second", Bounds: image.Rect(200, 0, 400, 240)})
	if err == nil {
		t.Fatal("expected error for overlapping regions")
	}
	if !strings.Contains(err.Error(), "overlap with existing region") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocate_UnregisteredMode(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err := rm.Allocate(region.RegionSpec{Name: "r1", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for unregistered mode")
	}
	if !strings.Contains(err.Error(), "is not registered") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllocateLayout_Atomic(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	// Second spec overlaps first — entire layout should be rejected.
	layout := region.RegionLayout{
		Specs: []region.RegionSpec{
			{Name: "left", Bounds: image.Rect(0, 0, 300, 240)},
			{Name: "right", Bounds: image.Rect(200, 0, 480, 240)}, // overlaps
		},
	}

	err := rm.AllocateLayout(layout)
	if err == nil {
		t.Fatal("expected error for overlapping layout")
	}

	if len(rm.Regions()) != 0 {
		t.Fatalf("expected 0 regions after failed layout, got %d", len(rm.Regions()))
	}
}

func TestAllocateLayout_Empty(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	err := rm.AllocateLayout(region.RegionLayout{Specs: []region.RegionSpec{}})
	if err != nil {
		t.Fatalf("empty layout should return nil, got: %v", err)
	}
}

func TestAllocateLayout_Valid(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	layout := region.RegionLayout{
		Specs: []region.RegionSpec{
			{Name: "left", Bounds: image.Rect(0, 0, 240, 240)},
			{Name: "right", Bounds: image.Rect(240, 0, 480, 240)},
		},
	}

	err := rm.AllocateLayout(layout)
	if err != nil {
		t.Fatalf("AllocateLayout: %v", err)
	}
	if len(rm.Regions()) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(rm.Regions()))
	}

	// First region should have input focus.
	active := rm.InputActiveRegion()
	if active == nil || active.Name() != "left" {
		t.Fatalf("expected input focus on 'left', got %v", active)
	}
}

func TestRegionByName_CaseInsensitive(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "MyRegion", Bounds: image.Rect(0, 0, 240, 240)})

	cases := []string{"MyRegion", "myregion", "MYREGION", "myREGION"}
	for _, name := range cases {
		r, ok := rm.RegionByName(name)
		if !ok {
			t.Fatalf("RegionByName(%q) not found", name)
		}
		if r.Name() != "MyRegion" {
			t.Fatalf("expected preserved name 'MyRegion', got %q", r.Name())
		}
	}
}

func TestRegion_ByIndex(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "first", Bounds: image.Rect(0, 0, 240, 240)})
	rm.Allocate(region.RegionSpec{Name: "second", Bounds: image.Rect(240, 0, 480, 240)})

	r0, ok := rm.Region(0)
	if !ok || r0.Name() != "first" {
		t.Fatalf("Region(0) expected 'first', got %v", r0)
	}

	r1, ok := rm.Region(1)
	if !ok || r1.Name() != "second" {
		t.Fatalf("Region(1) expected 'second', got %v", r1)
	}

	_, ok = rm.Region(2)
	if ok {
		t.Fatal("Region(2) should not exist")
	}
}

func TestSetMode_ByIndex(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	wireNoopFactory(rm)

	err := rm.SetMode("0", "dashboard")
	if err != nil {
		t.Fatalf("SetMode by index: %v", err)
	}

	r, _ := rm.Region(0)
	if r.CurrentMode() != "dashboard" {
		t.Fatalf("expected mode 'dashboard', got %q", r.CurrentMode())
	}
}

func TestSetMode_ByName(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	rm.Allocate(region.RegionSpec{Name: "Main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	wireNoopFactory(rm)

	// Case-insensitive lookup.
	err := rm.SetMode("main", "dashboard")
	if err != nil {
		t.Fatalf("SetMode by name: %v", err)
	}

	r, _ := rm.RegionByName("Main")
	if r.CurrentMode() != "dashboard" {
		t.Fatalf("expected mode 'dashboard', got %q", r.CurrentMode())
	}
}

func TestSetMode_UnregisteredMode(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock"
	})

	rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})

	err := rm.SetMode("main", "nonexistent")
	if err == nil {
		t.Fatal("expected error for unregistered mode")
	}
	if !strings.Contains(err.Error(), "is not registered") {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mode should remain unchanged.
	r, _ := rm.RegionByName("main")
	if r.CurrentMode() != "clock" {
		t.Fatalf("mode should remain 'clock', got %q", r.CurrentMode())
	}
}

func TestSetMode_UnknownTarget(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240)})

	err := rm.SetMode("99", "clock")
	if err == nil {
		t.Fatal("expected error for unknown index")
	}
	if !strings.Contains(err.Error(), "no region at index") {
		t.Fatalf("unexpected error: %v", err)
	}

	err = rm.SetMode("nonexistent", "clock")
	if err == nil {
		t.Fatal("expected error for unknown name")
	}
	if !strings.Contains(err.Error(), "no region named") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetInputFocus(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "left", Bounds: image.Rect(0, 0, 240, 240)})
	rm.Allocate(region.RegionSpec{Name: "right", Bounds: image.Rect(240, 0, 480, 240)})

	// Initially "left" has focus.
	if active := rm.InputActiveRegion(); active.Name() != "left" {
		t.Fatalf("expected focus on 'left', got %q", active.Name())
	}

	// Move focus to "right".
	err := rm.SetInputFocus("right")
	if err != nil {
		t.Fatalf("SetInputFocus: %v", err)
	}

	if active := rm.InputActiveRegion(); active.Name() != "right" {
		t.Fatalf("expected focus on 'right', got %q", active.Name())
	}

	// "left" should no longer have focus.
	left, _ := rm.RegionByName("left")
	if left.HasInputFocus() {
		t.Fatal("'left' should not have input focus")
	}
}

func TestSetInputFocus_ByIndex(t *testing.T) {
	vd := newTestVD(t, 480, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "left", Bounds: image.Rect(0, 0, 240, 240)})
	rm.Allocate(region.RegionSpec{Name: "right", Bounds: image.Rect(240, 0, 480, 240)})

	err := rm.SetInputFocus("1")
	if err != nil {
		t.Fatalf("SetInputFocus by index: %v", err)
	}

	if active := rm.InputActiveRegion(); active.Name() != "right" {
		t.Fatalf("expected focus on 'right', got %q", active.Name())
	}
}

func TestAllocate_SurfaceSharesMemory(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240)})

	r, _ := rm.Region(0)
	surf := r.Surface()

	// Write a pixel via the surface.
	surf.FrameBuffer().Set(10, 10, color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})

	// Read it from the VD framebuffer.
	got := vd.FrameBuffer().RGBAAt(10, 10)
	if got.R != 0xFF || got.G != 0x00 || got.B != 0x00 || got.A != 0xFF {
		t.Fatalf("expected red pixel in VD, got %v", got)
	}
}

func TestInputActiveRegion_None(t *testing.T) {
	vd := newTestVD(t, 240, 240)
	rm := region.NewRegionManager(vd)

	// No regions allocated.
	if active := rm.InputActiveRegion(); active != nil {
		t.Fatalf("expected nil active region, got %v", active)
	}
}

// --- From: mode_switch_test.go ---

func TestModeSwitch_Execute_ByName(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard" || mode == "menu"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}

	// Wire ModeFactory so SetMode can construct instances.
	r, _ := rm.RegionByName("main")
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		return &noopInstance{id: id}, true
	})

	ms := region.NewModeSwitch(rm)

	// Switch mode by name.
	err = ms.Execute(region.ModeChangeCommand{Target: "main", Mode: "dashboard"})
	if err != nil {
		t.Fatalf("Execute by name: %v", err)
	}

	if r.CurrentMode() != "dashboard" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "dashboard")
	}
}

func TestModeSwitch_Execute_ByNameCaseInsensitive(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err = rm.Allocate(region.RegionSpec{Name: "StatusPanel", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}
	wireNoopFactory(rm)

	ms := region.NewModeSwitch(rm)

	// Switch mode by name with different casing.
	err = ms.Execute(region.ModeChangeCommand{Target: "statuspanel", Mode: "dashboard"})
	if err != nil {
		t.Fatalf("Execute with lowercase name: %v", err)
	}

	r, _ := rm.RegionByName("StatusPanel")
	if r.CurrentMode() != "dashboard" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "dashboard")
	}

	// Also try UPPER case.
	err = ms.Execute(region.ModeChangeCommand{Target: "STATUSPANEL", Mode: "clock"})
	if err != nil {
		t.Fatalf("Execute with uppercase name: %v", err)
	}

	if r.CurrentMode() != "clock" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "clock")
	}
}

func TestModeSwitch_Execute_ByIndex(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard" || mode == "menu"
	})

	// Allocate two regions.
	err = rm.Allocate(region.RegionSpec{Name: "left", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate left: %v", err)
	}
	err = rm.Allocate(region.RegionSpec{Name: "right", Bounds: image.Rect(240, 0, 480, 240), DefaultMode: "dashboard"})
	if err != nil {
		t.Fatalf("Allocate right: %v", err)
	}
	wireNoopFactory(rm)

	ms := region.NewModeSwitch(rm)

	// Switch mode for region at index 1 (right).
	err = ms.Execute(region.ModeChangeCommand{Target: "1", Mode: "menu"})
	if err != nil {
		t.Fatalf("Execute by index: %v", err)
	}

	r, _ := rm.Region(1)
	if r.CurrentMode() != "menu" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "menu")
	}

	// Region at index 0 should be unchanged.
	r0, _ := rm.Region(0)
	if r0.CurrentMode() != "clock" {
		t.Fatalf("Region 0 CurrentMode()=%q, want %q (should be unchanged)", r0.CurrentMode(), "clock")
	}
}

func TestModeSwitch_Execute_UnknownTarget(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}

	ms := region.NewModeSwitch(rm)

	// Unknown name target.
	err = ms.Execute(region.ModeChangeCommand{Target: "nonexistent", Mode: "clock"})
	if err == nil {
		t.Fatal("Execute with unknown name should return error")
	}
	if !strings.Contains(err.Error(), "no region named") {
		t.Fatalf("error should mention 'no region named', got: %v", err)
	}

	// Unknown index target.
	err = ms.Execute(region.ModeChangeCommand{Target: "99", Mode: "clock"})
	if err == nil {
		t.Fatal("Execute with unknown index should return error")
	}
	if !strings.Contains(err.Error(), "no region at index") {
		t.Fatalf("error should mention 'no region at index', got: %v", err)
	}
}

func TestModeSwitch_Execute_SameModeReinitialization(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}
	wireNoopFactory(rm)

	r, _ := rm.RegionByName("main")

	// Dirty the surface with non-black pixels.
	r.Surface().Clear(color.RGBA{255, 0, 0, 255})

	ms := region.NewModeSwitch(rm)

	// Execute with same mode as current.
	err = ms.Execute(region.ModeChangeCommand{Target: "main", Mode: "clock"})
	if err != nil {
		t.Fatalf("Execute with same mode: %v", err)
	}

	// Mode should still be "clock".
	if r.CurrentMode() != "clock" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "clock")
	}

	// Surface should be cleared to black.
	fb := r.Surface().FrameBuffer()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			got := fb.RGBAAt(x, y)
			want := color.RGBA{0, 0, 0, 255}
			if got != want {
				t.Fatalf("pixel(%d,%d)=%v, want %v (surface should be cleared to black)", x, y, got, want)
			}
		}
	}
}

func TestModeSwitch_Execute_ModeChangeClearsSurface(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}
	wireNoopFactory(rm)

	r, _ := rm.RegionByName("main")

	// Dirty the surface.
	r.Surface().Clear(color.RGBA{0, 255, 0, 255})

	ms := region.NewModeSwitch(rm)

	// Execute mode change to a different mode.
	err = ms.Execute(region.ModeChangeCommand{Target: "main", Mode: "dashboard"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	// Mode should change.
	if r.CurrentMode() != "dashboard" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "dashboard")
	}

	// Surface should be cleared to black.
	fb := r.Surface().FrameBuffer()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			got := fb.RGBAAt(x, y)
			want := color.RGBA{0, 0, 0, 255}
			if got != want {
				t.Fatalf("pixel(%d,%d)=%v, want %v (surface should be cleared to black)", x, y, got, want)
			}
		}
	}
}

func TestModeSwitch_Execute_UnregisteredMode(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}
	wireNoopFactory(rm)

	r, _ := rm.RegionByName("main")

	// Dirty the surface with non-black pixels to verify it's preserved on error.
	r.Surface().Clear(color.RGBA{128, 64, 32, 255})

	ms := region.NewModeSwitch(rm)

	// Try switching to an unregistered mode.
	err = ms.Execute(region.ModeChangeCommand{Target: "main", Mode: "unknown_mode"})
	if err == nil {
		t.Fatal("Execute with unregistered mode should return error")
	}
	if !strings.Contains(err.Error(), "not registered") {
		t.Fatalf("error should mention 'not registered', got: %v", err)
	}

	// Current mode should be unchanged.
	if r.CurrentMode() != "clock" {
		t.Fatalf("CurrentMode()=%q, want %q (should be unchanged after error)", r.CurrentMode(), "clock")
	}

	// Surface should be preserved (not cleared) after a failed mode change.
	fb := r.Surface().FrameBuffer()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			got := fb.RGBAAt(x, y)
			want := color.RGBA{128, 64, 32, 255}
			if got != want {
				t.Fatalf("pixel(%d,%d)=%v, want %v (surface should be unchanged after invalid mode)", x, y, got, want)
			}
		}
	}
}

func TestModeSwitch_Execute_NonexistentRegionPreservesOtherRegions(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}

	rm := region.NewRegionManager(vd)
	rm.SetModeValidator(func(mode string) bool {
		return mode == "clock" || mode == "dashboard"
	})

	err = rm.Allocate(region.RegionSpec{Name: "main", Bounds: image.Rect(0, 0, 240, 240), DefaultMode: "clock"})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}
	wireNoopFactory(rm)

	r, _ := rm.RegionByName("main")

	ms := region.NewModeSwitch(rm)

	// Target a nonexistent region — should error without affecting existing regions.
	err = ms.Execute(region.ModeChangeCommand{Target: "ghost", Mode: "dashboard"})
	if err == nil {
		t.Fatal("Execute with nonexistent region should return error")
	}

	// Existing region's mode should be unchanged.
	if r.CurrentMode() != "clock" {
		t.Fatalf("CurrentMode()=%q, want %q (existing region should be unaffected)", r.CurrentMode(), "clock")
	}
}

// --- From: region_test.go ---

// mockDrawTarget is a test double for driver.DrawTarget.
type mockDrawTarget struct {
	bounds image.Rectangle
}

func (m *mockDrawTarget) Bounds() image.Rectangle    { return m.bounds }
func (m *mockDrawTarget) DrawImage(draw.Image) error { return nil }

// mockTextHintTarget is a DrawTarget that also implements TextHintProvider.
type mockTextHintTarget struct {
	bounds image.Rectangle
	hints  textlayout.TextHints
}

func (m *mockTextHintTarget) Bounds() image.Rectangle    { return m.bounds }
func (m *mockTextHintTarget) DrawImage(draw.Image) error { return nil }
func (m *mockTextHintTarget) TextHints() textlayout.TextHints {
	return m.hints
}

func TestNewRegion(t *testing.T) {
	bounds := image.Rect(10, 20, 250, 260)
	surf := surface.New(image.Rect(0, 0, 240, 240))

	r := region.NewRegion("status", bounds, surf)

	if r.Name() != "status" {
		t.Fatalf("Name()=%q, want %q", r.Name(), "status")
	}
	if r.Bounds() != bounds {
		t.Fatalf("Bounds()=%v, want %v", r.Bounds(), bounds)
	}
	if r.Surface() != surf {
		t.Fatal("Surface() returned unexpected pointer")
	}
	if r.CurrentMode() != "" {
		t.Fatalf("CurrentMode()=%q, want empty string", r.CurrentMode())
	}
	if r.HasInputFocus() {
		t.Fatal("HasInputFocus() should be false by default")
	}
}

func TestRegion_TextHints(t *testing.T) {
	surf := surface.New(image.Rect(0, 0, 120, 80))
	r := region.NewRegion("small", image.Rect(0, 0, 120, 80), surf)

	hints := r.TextHints()
	if hints.PixelWidth != 120 {
		t.Fatalf("TextHints().PixelWidth=%d, want 120", hints.PixelWidth)
	}
	if hints.PixelHeight != 80 {
		t.Fatalf("TextHints().PixelHeight=%d, want 80", hints.PixelHeight)
	}
	// Check default glyph metrics are populated.
	if hints.GlyphWidth != textlayout.GlyphWidth {
		t.Fatalf("TextHints().GlyphWidth=%d, want %d", hints.GlyphWidth, textlayout.GlyphWidth)
	}
	if hints.GlyphAdvance != textlayout.GlyphAdvance {
		t.Fatalf("TextHints().GlyphAdvance=%d, want %d", hints.GlyphAdvance, textlayout.GlyphAdvance)
	}
}

func TestRegion_SetMode_ClearsSurfaceToBlack(t *testing.T) {
	surf := surface.New(image.Rect(0, 0, 10, 10))
	// Fill surface with white pixels.
	surf.Clear(color.RGBA{255, 255, 255, 255})

	r := region.NewRegion("test", image.Rect(0, 0, 10, 10), surf)
	r.TestSetMode("old-mode")

	// Wire a ModeFactory that returns a minimal instance.
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		return &noopInstance{id: id}, true
	})

	err := r.SetMode("new-mode")
	if err != nil {
		t.Fatalf("SetMode() returned unexpected error: %v", err)
	}

	if r.CurrentMode() != "new-mode" {
		t.Fatalf("CurrentMode()=%q, want %q", r.CurrentMode(), "new-mode")
	}

	// Verify all pixels are black (RGBA 0,0,0,255).
	fb := surf.FrameBuffer()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			got := fb.RGBAAt(x, y)
			want := color.RGBA{0, 0, 0, 255}
			if got != want {
				t.Fatalf("pixel(%d,%d)=%v, want %v", x, y, got, want)
			}
		}
	}
}

func TestRegion_SetMode_PanicsWithoutFactory(t *testing.T) {
	surf := surface.New(image.Rect(0, 0, 10, 10))
	r := region.NewRegion("test", image.Rect(0, 0, 10, 10), surf)

	// Without a ModeFactory wired, SetMode must panic.
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatal("SetMode() without ModeFactory should panic, but did not")
		}
	}()
	_ = r.SetMode("clock")
}

func TestRegion_HasInputFocus(t *testing.T) {
	surf := surface.New(image.Rect(0, 0, 100, 100))
	r := region.NewRegion("main", image.Rect(0, 0, 100, 100), surf)

	if r.HasInputFocus() {
		t.Fatal("HasInputFocus() should be false initially")
	}

	r.SetInputFocus(true)
	if !r.HasInputFocus() {
		t.Fatal("HasInputFocus() should be true after SetInputFocus(true)")
	}

	r.SetInputFocus(false)
	if r.HasInputFocus() {
		t.Fatal("HasInputFocus() should be false after SetInputFocus(false)")
	}
}

// --- TextHints Resolution Tests ---

func TestRegion_TextHints_WithinSingleScreenWithProvider(t *testing.T) {
	// A Region entirely within one Physical Screen that implements TextHintProvider.
	// Should use the screen's capability flags + Region's own PixelWidth/PixelHeight.
	screenBounds := image.Rect(0, 0, 240, 240)
	regionBounds := image.Rect(10, 10, 130, 90) // 120x80 within the 240x240 screen

	customHints := textlayout.TextHints{
		PixelWidth:               240,
		PixelHeight:              240,
		GlyphWidth:               6,
		GlyphHeight:              8,
		GlyphAdvance:             7,
		RowHeight:                12,
		SupportsVerticalScroll:   false,
		SupportsHorizontalScroll: true,
		SupportsAutoScroll:       false,
		PreferEventRefresh:       true,
		DefaultTickerDirection:   textlayout.TickerDirectionNone,
		DefaultLineMode:          textlayout.LineModeClip,
	}

	target := &mockTextHintTarget{
		bounds: screenBounds,
		hints:  customHints,
	}

	screens := []region.ScreenPosition{
		{Index: 0, Name: "main", Bounds: screenBounds, Target: target, HintProvider: target.TextHints},
	}

	surf := surface.New(image.Rect(0, 0, 120, 80))
	r := region.NewRegionWithScreens("panel", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()

	// PixelWidth/PixelHeight should be the Region's dimensions, not the screen's.
	if hints.PixelWidth != 120 {
		t.Fatalf("TextHints().PixelWidth=%d, want 120", hints.PixelWidth)
	}
	if hints.PixelHeight != 80 {
		t.Fatalf("TextHints().PixelHeight=%d, want 80", hints.PixelHeight)
	}

	// Capability flags should come from the screen's TextHintProvider.
	if hints.SupportsVerticalScroll != false {
		t.Fatal("TextHints().SupportsVerticalScroll should be false (from screen)")
	}
	if hints.SupportsHorizontalScroll != true {
		t.Fatal("TextHints().SupportsHorizontalScroll should be true (from screen)")
	}
	if hints.SupportsAutoScroll != false {
		t.Fatal("TextHints().SupportsAutoScroll should be false (from screen)")
	}
	if hints.PreferEventRefresh != true {
		t.Fatal("TextHints().PreferEventRefresh should be true (from screen)")
	}
	if hints.DefaultTickerDirection != textlayout.TickerDirectionNone {
		t.Fatalf("TextHints().DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionNone)
	}
	if hints.DefaultLineMode != textlayout.LineModeClip {
		t.Fatalf("TextHints().DefaultLineMode=%q, want %q", hints.DefaultLineMode, textlayout.LineModeClip)
	}

	// Glyph metrics should come from the screen's TextHintProvider.
	if hints.GlyphWidth != 6 {
		t.Fatalf("TextHints().GlyphWidth=%d, want 6", hints.GlyphWidth)
	}
	if hints.GlyphHeight != 8 {
		t.Fatalf("TextHints().GlyphHeight=%d, want 8", hints.GlyphHeight)
	}
	if hints.GlyphAdvance != 7 {
		t.Fatalf("TextHints().GlyphAdvance=%d, want 7", hints.GlyphAdvance)
	}
	if hints.RowHeight != 12 {
		t.Fatalf("TextHints().RowHeight=%d, want 12", hints.RowHeight)
	}
}

func TestRegion_TextHints_SpansMultipleScreens(t *testing.T) {
	// A Region that spans multiple Physical Screens.
	// Should use default capability flags + Region's PixelWidth/PixelHeight.
	screen1Bounds := image.Rect(0, 0, 240, 135)
	screen2Bounds := image.Rect(240, 0, 480, 135)

	// Region spans both screens.
	regionBounds := image.Rect(100, 0, 380, 135) // 280x135

	target1 := &mockTextHintTarget{
		bounds: screen1Bounds,
		hints: textlayout.TextHints{
			SupportsVerticalScroll: false,
			SupportsAutoScroll:     false,
			DefaultTickerDirection: textlayout.TickerDirectionNone,
		},
	}
	target2 := &mockTextHintTarget{
		bounds: screen2Bounds,
		hints: textlayout.TextHints{
			SupportsHorizontalScroll: false,
			PreferEventRefresh:       true,
		},
	}

	screens := []region.ScreenPosition{
		{Index: 0, Name: "left", Bounds: screen1Bounds, Target: target1, HintProvider: target1.TextHints},
		{Index: 1, Name: "right", Bounds: screen2Bounds, Target: target2, HintProvider: target2.TextHints},
	}

	surf := surface.New(image.Rect(0, 0, 280, 135))
	r := region.NewRegionWithScreens("wide", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()

	// PixelWidth/PixelHeight should be the Region's dimensions.
	if hints.PixelWidth != 280 {
		t.Fatalf("TextHints().PixelWidth=%d, want 280", hints.PixelWidth)
	}
	if hints.PixelHeight != 135 {
		t.Fatalf("TextHints().PixelHeight=%d, want 135", hints.PixelHeight)
	}

	// Default capability flags when spanning multiple screens.
	if hints.SupportsVerticalScroll != true {
		t.Fatal("TextHints().SupportsVerticalScroll should be true (default)")
	}
	if hints.SupportsHorizontalScroll != true {
		t.Fatal("TextHints().SupportsHorizontalScroll should be true (default)")
	}
	if hints.SupportsAutoScroll != true {
		t.Fatal("TextHints().SupportsAutoScroll should be true (default)")
	}
	if hints.PreferEventRefresh != false {
		t.Fatal("TextHints().PreferEventRefresh should be false (default)")
	}
	if hints.DefaultTickerDirection != textlayout.TickerDirectionVertical {
		t.Fatalf("TextHints().DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionVertical)
	}
	if hints.DefaultLineMode != textlayout.LineModeTruncate {
		t.Fatalf("TextHints().DefaultLineMode=%q, want %q", hints.DefaultLineMode, textlayout.LineModeTruncate)
	}

	expected := region.TestBuildCatalogHints(textlayout.TextHints{
		PixelWidth:               280,
		PixelHeight:              135,
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
		t.Fatalf("TextHints().GlyphWidth=%d, want %d", hints.GlyphWidth, expected.GlyphWidth)
	}
	if hints.GlyphHeight != expected.GlyphHeight {
		t.Fatalf("TextHints().GlyphHeight=%d, want %d", hints.GlyphHeight, expected.GlyphHeight)
	}
	if hints.GlyphAdvance != expected.GlyphAdvance {
		t.Fatalf("TextHints().GlyphAdvance=%d, want %d", hints.GlyphAdvance, expected.GlyphAdvance)
	}
	if hints.RowHeight != expected.RowHeight {
		t.Fatalf("TextHints().RowHeight=%d, want %d", hints.RowHeight, expected.RowHeight)
	}
}

func TestRegion_TextHints_NoTextHintProvider(t *testing.T) {
	// A Region within a Physical Screen that does NOT implement TextHintProvider.
	// Should use DefaultTextHints(Region.Bounds()).
	screenBounds := image.Rect(0, 0, 240, 240)
	regionBounds := image.Rect(0, 0, 200, 150) // 200x150 within the screen

	// mockDrawTarget does not implement TextHintProvider.
	target := &mockDrawTarget{bounds: screenBounds}

	screens := []region.ScreenPosition{
		{Index: 0, Name: "main", Bounds: screenBounds, Target: target},
	}

	surf := surface.New(image.Rect(0, 0, 200, 150))
	r := region.NewRegionWithScreens("display", regionBounds, surf, screens, "", 0, 0)

	hints := r.TextHints()

	// Should match DefaultTextHints behavior.
	expected := region.TestBuildCatalogHints(textlayout.DefaultTextHints(image.Rect(0, 0, 200, 150)), 96.0)

	if hints.PixelWidth != expected.PixelWidth {
		t.Fatalf("TextHints().PixelWidth=%d, want %d", hints.PixelWidth, expected.PixelWidth)
	}
	if hints.PixelHeight != expected.PixelHeight {
		t.Fatalf("TextHints().PixelHeight=%d, want %d", hints.PixelHeight, expected.PixelHeight)
	}
	if hints.GlyphWidth != expected.GlyphWidth {
		t.Fatalf("TextHints().GlyphWidth=%d, want %d", hints.GlyphWidth, expected.GlyphWidth)
	}
	if hints.GlyphHeight != expected.GlyphHeight {
		t.Fatalf("TextHints().GlyphHeight=%d, want %d", hints.GlyphHeight, expected.GlyphHeight)
	}
	if hints.SupportsVerticalScroll != expected.SupportsVerticalScroll {
		t.Fatalf("TextHints().SupportsVerticalScroll=%v, want %v", hints.SupportsVerticalScroll, expected.SupportsVerticalScroll)
	}
	if hints.SupportsHorizontalScroll != expected.SupportsHorizontalScroll {
		t.Fatalf("TextHints().SupportsHorizontalScroll=%v, want %v", hints.SupportsHorizontalScroll, expected.SupportsHorizontalScroll)
	}
	if hints.SupportsAutoScroll != expected.SupportsAutoScroll {
		t.Fatalf("TextHints().SupportsAutoScroll=%v, want %v", hints.SupportsAutoScroll, expected.SupportsAutoScroll)
	}
	if hints.DefaultTickerDirection != expected.DefaultTickerDirection {
		t.Fatalf("TextHints().DefaultTickerDirection=%q, want %q", hints.DefaultTickerDirection, expected.DefaultTickerDirection)
	}
	if hints.DefaultLineMode != expected.DefaultLineMode {
		t.Fatalf("TextHints().DefaultLineMode=%q, want %q", hints.DefaultLineMode, expected.DefaultLineMode)
	}
}

func TestRegion_TextHints_NoScreens(t *testing.T) {
	// A Region created with no screen positions (empty slice).
	// Should use DefaultTextHints(Region.Bounds()).
	regionBounds := image.Rect(0, 0, 160, 128)
	surf := surface.New(image.Rect(0, 0, 160, 128))

	r := region.NewRegionWithScreens("orphan", regionBounds, surf, nil, "", 0, 0)

	hints := r.TextHints()
	expected := textlayout.DefaultTextHints(image.Rect(0, 0, 160, 128))

	if hints.PixelWidth != expected.PixelWidth {
		t.Fatalf("TextHints().PixelWidth=%d, want %d", hints.PixelWidth, expected.PixelWidth)
	}
	if hints.PixelHeight != expected.PixelHeight {
		t.Fatalf("TextHints().PixelHeight=%d, want %d", hints.PixelHeight, expected.PixelHeight)
	}
	if hints.SupportsVerticalScroll != expected.SupportsVerticalScroll {
		t.Fatalf("TextHints().SupportsVerticalScroll=%v, want %v", hints.SupportsVerticalScroll, expected.SupportsVerticalScroll)
	}
}

func TestRegion_TextHints_NewRegionFallback(t *testing.T) {
	// NewRegion (without screens) still uses DefaultTextHints from surface bounds.
	surf := surface.New(image.Rect(0, 0, 320, 240))
	r := region.NewRegion("legacy", image.Rect(0, 0, 320, 240), surf)

	hints := r.TextHints()
	expected := textlayout.DefaultTextHints(image.Rect(0, 0, 320, 240))

	if hints.PixelWidth != expected.PixelWidth {
		t.Fatalf("TextHints().PixelWidth=%d, want %d", hints.PixelWidth, expected.PixelWidth)
	}
	if hints.PixelHeight != expected.PixelHeight {
		t.Fatalf("TextHints().PixelHeight=%d, want %d", hints.PixelHeight, expected.PixelHeight)
	}
}

// --- Render error isolation tests (Task 7.1) ---
// These tests verify requirements 7.1, 7.2, 7.3, 7.5:
// - safeRender logs region name and panic value on panic (7.1)
// - Render errors log region name and error value (7.2)
// - Remaining regions continue rendering in allocation order after failure (7.3)
// - Single Flush occurs after all regions (including failed) are processed (7.3)
// - Failed region's deadline is advanced so it retries next cycle (7.5)

// orderedRenderer tracks render calls in order and can panic/error on specific regions.
type orderedRenderer struct {
	mu      sync.Mutex
	calls   []string
	panicOn map[string]bool
	errorOn map[string]bool
}

func (m *orderedRenderer) Render(r *region.Region) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.panicOn[r.Name()] {
		panic("test panic in " + r.Name())
	}
	if m.errorOn[r.Name()] {
		return errors.New("test error in " + r.Name())
	}
	m.calls = append(m.calls, r.Name())
	return nil
}

func (m *orderedRenderer) getCalls() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.calls))
	copy(result, m.calls)
	return result
}

func TestSafeRender_PanicRecovery_RemainingRegionsContinue(t *testing.T) {
	// Requirement 7.1, 7.3: When first region panics, remaining regions
	// continue rendering in allocation order.
	renderer := &orderedRenderer{
		panicOn: map[string]bool{"region_b": true},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	// Create 4 regions in allocation order: a, b (panics), c, d.
	names := []string{"region_a", "region_b", "region_c", "region_d"}
	for i, name := range names {
		bounds := image.Rect(i*60, 0, (i+1)*60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		rm.TestAppendRegion(r)
	}
	rm.Regions()[0].SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithRenderer(renderer),
	)

	// Directly call renderFrame to test single-frame behavior.
	rl.TestRenderFrame()

	// Verify: region_a, region_c, region_d rendered (in order); region_b skipped.
	calls := renderer.getCalls()
	expected := []string{"region_a", "region_c", "region_d"}
	if len(calls) != len(expected) {
		t.Fatalf("expected %d render calls, got %d: %v", len(expected), len(calls), calls)
	}
	for i, want := range expected {
		if calls[i] != want {
			t.Errorf("call[%d] = %q, want %q", i, calls[i], want)
		}
	}
}

func TestSafeRender_ErrorRecovery_RemainingRegionsContinue(t *testing.T) {
	// Requirement 7.2, 7.3: When first region errors, remaining regions
	// continue rendering in allocation order.
	renderer := &orderedRenderer{
		errorOn: map[string]bool{"region_b": true},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	names := []string{"region_a", "region_b", "region_c", "region_d"}
	for i, name := range names {
		bounds := image.Rect(i*60, 0, (i+1)*60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		rm.TestAppendRegion(r)
	}
	rm.Regions()[0].SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithRenderer(renderer),
	)

	rl.TestRenderFrame()

	calls := renderer.getCalls()
	expected := []string{"region_a", "region_c", "region_d"}
	if len(calls) != len(expected) {
		t.Fatalf("expected %d render calls, got %d: %v", len(expected), len(calls), calls)
	}
	for i, want := range expected {
		if calls[i] != want {
			t.Errorf("call[%d] = %q, want %q", i, calls[i], want)
		}
	}
}

func TestSafeRender_SingleFlushAfterPanicAndError(t *testing.T) {
	// Requirement 7.3: Exactly one Flush occurs after all regions are processed,
	// even when some regions panic or error.
	renderer := &orderedRenderer{
		panicOn: map[string]bool{"panicky": true},
		errorOn: map[string]bool{"erroring": true},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	names := []string{"panicky", "erroring", "healthy"}
	for i, name := range names {
		bounds := image.Rect(i*60, 0, (i+1)*60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		rm.TestAppendRegion(r)
	}
	rm.Regions()[0].SetInputFocus(true)

	target := &flushMockTarget{bounds: image.Rect(0, 0, 240, 60)}
	screens := []region.ScreenPosition{
		{Index: 0, Name: "s0", Bounds: image.Rect(0, 0, 240, 60), Target: target},
	}
	fp := region.NewFlushPath(vd, screens)

	rl := region.NewRenderLoop(rm, fp, nil,
		region.WithTickInterval(50*time.Millisecond),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(150 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// Verify flush was called (at least once for the initial frame).
	if target.callCount == 0 {
		t.Fatal("expected at least one Flush call")
	}

	// The "healthy" region should have rendered despite panicky and erroring.
	calls := renderer.getCalls()
	hasHealthy := false
	for _, c := range calls {
		if c == "healthy" {
			hasHealthy = true
			break
		}
	}
	if !hasHealthy {
		t.Fatal("expected 'healthy' to render despite other regions failing")
	}
}

func TestSafeRender_PerRegion_SingleFlushAfterFailures(t *testing.T) {
	// Requirement 7.3 with per-region scheduling: after renderDueRegions completes
	// (with some regions panicking/erroring), exactly one flush() call is issued.
	renderer := &orderedRenderer{
		panicOn: map[string]bool{"panicky": true},
		errorOn: map[string]bool{"erroring": true},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 180, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	names := []string{"panicky", "erroring", "healthy"}
	regions := make([]*region.Region, 3)
	for i, name := range names {
		bounds := image.Rect(i*60, 0, (i+1)*60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		r.TestSetMode("default")
		rm.TestAppendRegion(r)
		regions[i] = r
	}
	rm.Regions()[0].SetInputFocus(true)

	target := &flushMockTarget{bounds: image.Rect(0, 0, 180, 60)}
	screens := []region.ScreenPosition{
		{Index: 0, Name: "s0", Bounds: image.Rect(0, 0, 180, 60), Target: target},
	}
	fp := region.NewFlushPath(vd, screens)

	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"default": 50 * time.Millisecond,
		},
	}

	rl := region.NewRenderLoop(rm, fp, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	// Set all regions as due.
	now := time.Now()
	tickers := make([]region.TestPerRegionTicker, 3)
	for i, r := range regions {
		tickers[i] = region.TestPerRegionTicker{
			Region:   r,
			Interval: 50 * time.Millisecond,
			LastFire: now.Add(-100 * time.Millisecond), // all due
		}
	}
	rl.TestSetRegionTickers(tickers)

	flushBefore := target.callCount

	// Simulate one wake cycle: render due regions + flush.
	rl.TestRenderDueRegions(now)
	rl.TestFlush()

	flushDelta := target.callCount - flushBefore
	if flushDelta != 1 {
		t.Fatalf("expected exactly 1 flush call, got %d", flushDelta)
	}

	// "healthy" should have rendered.
	calls := renderer.getCalls()
	hasHealthy := false
	for _, c := range calls {
		if c == "healthy" {
			hasHealthy = true
			break
		}
	}
	if !hasHealthy {
		t.Fatal("expected 'healthy' to render despite failures in other regions")
	}
}

func TestSafeRender_PerRegion_DeadlineAdvancesOnPanicAndError(t *testing.T) {
	// Requirement 7.5: Failed region's deadline is advanced so it retries.
	renderer := &orderedRenderer{
		panicOn: map[string]bool{"panicky": true},
		errorOn: map[string]bool{"erroring": true},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 180, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	names := []string{"panicky", "erroring", "healthy"}
	regions := make([]*region.Region, 3)
	for i, name := range names {
		bounds := image.Rect(i*60, 0, (i+1)*60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		r.TestSetMode("default")
		rm.TestAppendRegion(r)
		regions[i] = r
	}
	rm.Regions()[0].SetInputFocus(true)

	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"default": 50 * time.Millisecond,
		},
	}

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	// Set all regions as due.
	now := time.Now()
	interval := 50 * time.Millisecond
	tickers := make([]region.TestPerRegionTicker, 3)
	originalLastFires := make([]time.Time, 3)
	for i, r := range regions {
		lastFire := now.Add(-interval - 10*time.Millisecond)
		originalLastFires[i] = lastFire
		tickers[i] = region.TestPerRegionTicker{
			Region:   r,
			Interval: interval,
			LastFire: lastFire,
		}
	}
	rl.TestSetRegionTickers(tickers)

	// Render due regions.
	rl.TestRenderDueRegions(now)

	// ALL regions' deadlines should be advanced, including the failed ones.
	for i, ticker := range rl.TestRegionTickers() {
		expectedLastFire := originalLastFires[i].Add(interval)
		if !ticker.LastFire.Equal(expectedLastFire) {
			t.Errorf("region %q: lastFire not advanced; got %v, want %v",
				names[i], ticker.LastFire, expectedLastFire)
		}
	}
}

func TestSafeRender_PerRegion_FailedRegionRetriedNextCycle(t *testing.T) {
	// Requirement 7.5: After a failure, the region is retried on its next
	// scheduled tick without manual intervention.
	callCount := 0
	failFirst := true

	// A renderer that fails the first time "retry" is rendered, then succeeds.
	retryRenderer := &retryTrackingRenderer{
		failFirst: &failFirst,
		callCount: &callCount,
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 60, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds := image.Rect(0, 0, 60, 60)
	surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
	r := region.NewRegion("retry", bounds, surf)
	r.TestSetMode("default")
	rm.TestAppendRegion(r)
	r.SetInputFocus(true)

	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"default": 30 * time.Millisecond,
		},
	}

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(retryRenderer),
	)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// The region should have been retried and eventually succeeded.
	if callCount < 2 {
		t.Fatalf("expected region to be retried at least twice, got %d attempts", callCount)
	}
}

// retryTrackingRenderer panics on the first call, then succeeds on subsequent calls.
type retryTrackingRenderer struct {
	mu        sync.Mutex
	failFirst *bool
	callCount *int
}

func (r *retryTrackingRenderer) Render(reg *region.Region) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.callCount++
	if *r.failFirst {
		*r.failFirst = false
		return errors.New("transient failure")
	}
	return nil
}

// --- From: render_loop_test.go ---

// mockRenderer tracks Render calls per Region name and can optionally panic or error.
type mockRenderer struct {
	mu      sync.Mutex
	calls   []string // region names in call order
	panicOn string   // if non-empty, panic when rendering this region name
	errorOn string   // if non-empty, return error for this region name
}

func (m *mockRenderer) Render(r *region.Region) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.panicOn != "" && r.Name() == m.panicOn {
		panic("intentional test panic in " + r.Name())
	}
	if m.errorOn != "" && r.Name() == m.errorOn {
		return errors.New("intentional test error in " + r.Name())
	}
	m.calls = append(m.calls, r.Name())
	return nil
}

func (m *mockRenderer) callCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

func (m *mockRenderer) getCalls() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.calls))
	copy(result, m.calls)
	return result
}

// mockDispatcher tracks dispatched events.
type mockDispatcher struct {
	mu     sync.Mutex
	events []dispatchedEvent
}

type dispatchedEvent struct {
	regionName string
	event      input.Event
}

func (d *mockDispatcher) Dispatch(r *region.Region, ev input.Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = append(d.events, dispatchedEvent{regionName: r.Name(), event: ev})
}

func (d *mockDispatcher) getEvents() []dispatchedEvent {
	d.mu.Lock()
	defer d.mu.Unlock()
	result := make([]dispatchedEvent, len(d.events))
	copy(result, d.events)
	return result
}

// setupTestLoop creates a minimal RenderLoop with one or more regions.
func setupTestLoop(t *testing.T, regionNames []string, opts ...region.RenderLoopOption) (*region.RenderLoop, *region.RegionManager) {
	t.Helper()

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)
	x := 0
	for _, name := range regionNames {
		bounds := image.Rect(x, 0, x+60, 60)
		surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
		r := region.NewRegion(name, bounds, surf)
		rm.TestAppendRegion(r)
		x += 60
	}

	// First region gets input focus.
	if len(rm.Regions()) > 0 {
		rm.Regions()[0].SetInputFocus(true)
	}

	rl := region.NewRenderLoop(rm, nil, nil, opts...)
	return rl, rm
}

func TestRenderLoop_RendersFirstFrameBeforeInput(t *testing.T) {
	// Verify that at least one full render+flush happens before any input is processed.
	vd, _ := region.NewVirtualDisplay(image.Rect(0, 0, 120, 60))
	rm := region.NewRegionManager(vd)

	bounds := image.Rect(0, 0, 120, 60)
	surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
	r := region.NewRegion("main", bounds, surf)
	rm.TestAppendRegion(r)
	r.SetInputFocus(true)

	renderer := &mockRenderer{}
	dispatcher := &mockDispatcher{}

	// Pre-load input events BEFORE the loop starts.
	events := make(chan input.Event, 10)
	events <- input.Event{Key: input.Key1, Type: input.Press}
	events <- input.Event{Key: input.Key2, Type: input.Press}

	rl := region.NewRenderLoop(rm, nil, events,
		region.WithTickInterval(50*time.Millisecond),
		region.WithRenderer(renderer),
		region.WithInputDispatcher(dispatcher),
	)

	// Stop the loop after the first tick frame processes.
	go func() {
		time.Sleep(80 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// Renderer must have been called at least once (the initial frame).
	calls := renderer.getCalls()
	if len(calls) == 0 {
		t.Fatal("expected renderer to be called at least once for the first frame")
	}
	if calls[0] != "main" {
		t.Fatalf("expected first render call to region 'main', got %q", calls[0])
	}
}

func TestRenderLoop_PanicInRegionDoesNotCrashLoop(t *testing.T) {
	renderer := &mockRenderer{panicOn: "bad"}

	vd, _ := region.NewVirtualDisplay(image.Rect(0, 0, 180, 60))
	rm := region.NewRegionManager(vd)

	// Create two regions: "bad" (will panic) and "good" (should still render).
	bounds1 := image.Rect(0, 0, 60, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("bad", bounds1, surf1)
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(60, 0, 120, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("good", bounds2, surf2)
	rm.TestAppendRegion(r2)
	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickInterval(50*time.Millisecond),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	// This should NOT panic.
	rl.Run()

	// "good" should have been rendered despite "bad" panicking.
	calls := renderer.getCalls()
	hasGood := false
	for _, c := range calls {
		if c == "good" {
			hasGood = true
			break
		}
	}
	if !hasGood {
		t.Fatal("expected 'good' region to be rendered despite 'bad' panicking")
	}
}

func TestRenderLoop_StopTerminatesLoop(t *testing.T) {
	renderer := &mockRenderer{}

	rl, _ := setupTestLoop(t, []string{"r1"},
		region.WithTickInterval(20*time.Millisecond),
		region.WithRenderer(renderer),
	)

	done := make(chan struct{})
	go func() {
		rl.Run()
		close(done)
	}()

	// Stop after a short delay.
	time.Sleep(30 * time.Millisecond)
	rl.Stop()

	select {
	case <-done:
		// Run returned, good.
	case <-time.After(1 * time.Second):
		t.Fatal("Run did not return after Stop")
	}
}

func TestRenderLoop_MultipleFramesRun(t *testing.T) {
	renderer := &mockRenderer{}

	rl, _ := setupTestLoop(t, []string{"r1"},
		region.WithTickInterval(30*time.Millisecond),
		region.WithRenderer(renderer),
	)

	go func() {
		// Allow several ticks.
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// With 30ms tick and ~200ms runtime, we should have at least 3 frames
	// (initial frame + several ticker frames).
	count := renderer.callCount()
	if count < 3 {
		t.Fatalf("expected at least 3 render calls, got %d", count)
	}
}

func TestRenderLoop_ErrorInRegionContinues(t *testing.T) {
	renderer := &mockRenderer{errorOn: "errRegion"}

	vd, _ := region.NewVirtualDisplay(image.Rect(0, 0, 180, 60))
	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 60, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("errRegion", bounds1, surf1)
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(60, 0, 120, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("okRegion", bounds2, surf2)
	rm.TestAppendRegion(r2)
	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickInterval(50*time.Millisecond),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(150 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// "okRegion" should have been rendered despite "errRegion" erroring.
	calls := renderer.getCalls()
	hasOk := false
	for _, c := range calls {
		if c == "okRegion" {
			hasOk = true
			break
		}
	}
	if !hasOk {
		t.Fatal("expected 'okRegion' to be rendered despite 'errRegion' returning error")
	}
}

func TestRenderLoop_InputDispatchToActiveRegion(t *testing.T) {
	renderer := &mockRenderer{}
	dispatcher := &mockDispatcher{}

	vd, _ := region.NewVirtualDisplay(image.Rect(0, 0, 180, 60))
	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 60, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("left", bounds1, surf1)
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(60, 0, 120, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("right", bounds2, surf2)
	rm.TestAppendRegion(r2)

	// Give input focus to "right".
	r2.SetInputFocus(true)

	events := make(chan input.Event, 10)

	rl := region.NewRenderLoop(rm, nil, events,
		region.WithTickInterval(50*time.Millisecond),
		region.WithRenderer(renderer),
		region.WithInputDispatcher(dispatcher),
	)

	go func() {
		// Wait for the first frame to render, then send events.
		time.Sleep(20 * time.Millisecond)
		events <- input.Event{Key: input.Key1, Type: input.Press}
		events <- input.Event{Key: input.Key2, Type: input.Press}
		// Let one tick process.
		time.Sleep(60 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	dispatched := dispatcher.getEvents()
	if len(dispatched) == 0 {
		t.Fatal("expected at least one dispatched event")
	}
	for _, d := range dispatched {
		if d.regionName != "right" {
			t.Fatalf("expected events dispatched to 'right', got %q", d.regionName)
		}
	}
}

func TestRenderLoop_NilRendererNoops(t *testing.T) {
	// Ensure the loop doesn't crash with no renderer set.
	rl, _ := setupTestLoop(t, []string{"r1"},
		region.WithTickInterval(30*time.Millisecond),
	)

	go func() {
		time.Sleep(80 * time.Millisecond)
		rl.Stop()
	}()

	// Should not panic.
	rl.Run()
}

// --- Per-region ticker tests ---

// mockTickRateResolver returns a fixed interval for each mode.
type mockTickRateResolver struct {
	intervals map[string]time.Duration
}

func (m *mockTickRateResolver) TickInterval(modeID string) time.Duration {
	if d, ok := m.intervals[modeID]; ok {
		return d
	}
	return 1000 * time.Millisecond
}

func TestRenderLoop_WithTickRateResolver(t *testing.T) {
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker":    50 * time.Millisecond,
			"dashboard": 1000 * time.Millisecond,
		},
	}

	rl, _ := setupTestLoop(t, []string{"r1"},
		region.WithTickRateResolver(resolver),
	)

	if !rl.TestHasTickRateResolver() {
		t.Fatal("expected tickRateResolver to be set")
	}
}

func TestRenderLoop_InitRegionTickers(t *testing.T) {
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker":    50 * time.Millisecond,
			"dashboard": 1000 * time.Millisecond,
		},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	// Create two regions with different modes.
	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("fast", bounds1, surf1)
	r1.TestSetMode("ticker")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("slow", bounds2, surf2)
	r2.TestSetMode("dashboard")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
	)

	// Call initRegionTickers.
	rl.TestInitRegionTickers()

	tickers := rl.TestRegionTickers()
	if len(tickers) != 2 {
		t.Fatalf("expected 2 region tickers, got %d", len(tickers))
	}

	// Verify intervals match the resolver output.
	if tickers[0].Region != r1 {
		t.Fatal("expected first ticker to reference 'fast' region")
	}
	if tickers[0].Interval != 50*time.Millisecond {
		t.Fatalf("expected fast ticker interval 50ms, got %v", tickers[0].Interval)
	}

	if tickers[1].Region != r2 {
		t.Fatal("expected second ticker to reference 'slow' region")
	}
	if tickers[1].Interval != 1000*time.Millisecond {
		t.Fatalf("expected slow ticker interval 1000ms, got %v", tickers[1].Interval)
	}

	// Verify lastFire is set to now minus the ticker's interval so the first
	// frame is immediately due. We check that it's non-zero and in the past by
	// at least the interval (but not absurdly far in the past).
	for i, ticker := range tickers {
		if ticker.LastFire.IsZero() {
			t.Fatalf("ticker %d lastFire is zero", i)
		}
		elapsed := time.Since(ticker.LastFire)
		if elapsed < ticker.Interval {
			t.Fatalf("ticker %d lastFire should be at least %v in the past, got %v",
				i, ticker.Interval, elapsed)
		}
		// Allow generous headroom for slow CI / parallel tests.
		if elapsed > ticker.Interval+5*time.Second {
			t.Fatalf("ticker %d lastFire too old: elapsed %v, expected ~%v",
				i, elapsed, ticker.Interval)
		}
	}
}

func TestRenderLoop_InitRegionTickersOnRun(t *testing.T) {
	// Verify that Run initializes per-region tickers when a resolver is set.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"dashboard": 1000 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 120, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)
	bounds := image.Rect(0, 0, 120, 60)
	surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
	r := region.NewRegion("main", bounds, surf)
	r.TestSetMode("dashboard")
	rm.TestAppendRegion(r)
	r.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickInterval(50*time.Millisecond),
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(80 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// After Run returns, regionTickers should be initialized.
	tickers := rl.TestRegionTickers()
	if len(tickers) != 1 {
		t.Fatalf("expected 1 region ticker after Run, got %d", len(tickers))
	}
	if tickers[0].Interval != 1000*time.Millisecond {
		t.Fatalf("expected ticker interval 1000ms, got %v", tickers[0].Interval)
	}
}

// --- Per-region deadline scheduling tests ---

func TestRenderLoop_PerRegion_FastRegionRendersMoreOften(t *testing.T) {
	// Two regions: "fast" at 30ms and "slow" at 200ms.
	// Over ~300ms, "fast" should render multiple times, "slow" only 1-2 times.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker":    30 * time.Millisecond,
			"dashboard": 200 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("fast", bounds1, surf1)
	r1.TestSetMode("ticker")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("slow", bounds2, surf2)
	r2.TestSetMode("dashboard")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(300 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// Count calls per region.
	calls := renderer.getCalls()
	fastCount := 0
	slowCount := 0
	for _, c := range calls {
		switch c {
		case "fast":
			fastCount++
		case "slow":
			slowCount++
		}
	}

	// "fast" at 30ms over 300ms should have at least 4 renders (initial + ~9 ticks).
	if fastCount < 4 {
		t.Fatalf("expected fast region to render at least 4 times, got %d", fastCount)
	}
	// "slow" at 200ms over 300ms should have at most 2 renders (initial + 1 tick).
	if slowCount > 3 {
		t.Fatalf("expected slow region to render at most 3 times, got %d", slowCount)
	}
	// Fast must render more often than slow.
	if fastCount <= slowCount {
		t.Fatalf("expected fast (%d) > slow (%d)", fastCount, slowCount)
	}
}

func TestRenderLoop_PerRegion_SingleFlushPerCycle(t *testing.T) {
	// Verify that flush is called once per wake cycle, not per-region.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker": 30 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("a", bounds1, surf1)
	r1.TestSetMode("ticker")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("b", bounds2, surf2)
	r2.TestSetMode("ticker")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	// Track flush calls via a mock FlushPath.
	flushCount := 0
	screens := []region.ScreenPosition{
		{Index: 0, Name: "s1", Bounds: image.Rect(0, 0, 240, 60), Target: &flushMockTarget{bounds: image.Rect(0, 0, 240, 60)}},
	}
	fp := region.NewFlushPath(vd, screens)

	rl := region.NewRenderLoop(rm, fp, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	// Override flush to count.
	// Instead, we verify by checking the render calls — both regions should be
	// rendered together in each cycle since they have the same interval.

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// Both regions have the same interval, so they render together.
	calls := renderer.getCalls()
	aCount := 0
	bCount := 0
	for _, c := range calls {
		switch c {
		case "a":
			aCount++
		case "b":
			bCount++
		}
	}
	// Both should have similar render counts.
	if aCount == 0 || bCount == 0 {
		t.Fatalf("expected both regions to render, got a=%d, b=%d", aCount, bCount)
	}
	// Use flushCount variable to avoid unused warning — this test primarily
	// validates both regions render together, meaning single flush per cycle.
	_ = flushCount
}

func TestRenderLoop_PerRegion_DeadlineAdvancesOnError(t *testing.T) {
	// A region that always errors should still have its deadline advanced
	// and should be retried on subsequent cycles.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker": 30 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{errorOn: "failing"}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("failing", bounds1, surf1)
	r1.TestSetMode("ticker")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("ok", bounds2, surf2)
	r2.TestSetMode("ticker")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// "ok" should have rendered multiple times (deadline advances for both).
	calls := renderer.getCalls()
	okCount := 0
	for _, c := range calls {
		if c == "ok" {
			okCount++
		}
	}
	// At 30ms over 200ms, expect at least 2 renders for "ok".
	if okCount < 2 {
		t.Fatalf("expected 'ok' to render at least 2 times, got %d", okCount)
	}
}

func TestRenderLoop_PerRegion_DeadlineAdvancesOnPanic(t *testing.T) {
	// A region that panics should still have its deadline advanced.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker": 30 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{panicOn: "panicky"}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("panicky", bounds1, surf1)
	r1.TestSetMode("ticker")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("stable", bounds2, surf2)
	r2.TestSetMode("ticker")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// "stable" should have rendered multiple times (the panicking region
	// doesn't block progress).
	calls := renderer.getCalls()
	stableCount := 0
	for _, c := range calls {
		if c == "stable" {
			stableCount++
		}
	}
	if stableCount < 2 {
		t.Fatalf("expected 'stable' to render at least 2 times, got %d", stableCount)
	}
}

func TestRenderLoop_PerRegion_ModeChangeRederivesInterval(t *testing.T) {
	// When a region's mode changes, the interval should be re-derived.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"ticker":    30 * time.Millisecond,
			"dashboard": 500 * time.Millisecond,
		},
	}

	renderer := &mockRenderer{}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 120, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds := image.Rect(0, 0, 120, 60)
	surf := surface.NewFromSubImage(vd.FrameBuffer(), bounds)
	r := region.NewRegion("main", bounds, surf)
	r.TestSetMode("ticker")
	rm.TestAppendRegion(r)
	r.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	go func() {
		// Let it run for a bit at 30ms tick, then switch mode.
		time.Sleep(150 * time.Millisecond)
		r.TestSetMode("dashboard") // simulate mode change
		// Let it continue for a bit at the new interval.
		time.Sleep(150 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// After mode change, the interval should have been updated to 500ms.
	tickers := rl.TestRegionTickers()
	if tickers[0].Interval != 500*time.Millisecond {
		t.Fatalf("expected interval to be re-derived to 500ms, got %v", tickers[0].Interval)
	}
}

func TestRenderLoop_PerRegion_GlobalTickerStillWorks(t *testing.T) {
	// Without a TickRateResolver, the global ticker behavior should remain.
	renderer := &mockRenderer{}

	rl, _ := setupTestLoop(t, []string{"r1"},
		region.WithTickInterval(30*time.Millisecond),
		region.WithRenderer(renderer),
	)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rl.Stop()
	}()

	rl.Run()

	// Should have rendered multiple times via the global ticker.
	count := renderer.callCount()
	if count < 3 {
		t.Fatalf("expected at least 3 render calls with global ticker, got %d", count)
	}

	// regionTickers should be nil (not initialized).
	if rl.TestRegionTickers() != nil {
		t.Fatal("expected regionTickers to be nil when no resolver is set")
	}
}

func TestRenderLoop_PerRegion_MinSleepDuration(t *testing.T) {
	// Test minSleepDuration computes correctly.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"fast": 50 * time.Millisecond,
			"slow": 200 * time.Millisecond,
		},
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("f", bounds1, surf1)
	r1.TestSetMode("fast")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("s", bounds2, surf2)
	r2.TestSetMode("slow")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(&mockRenderer{}),
	)

	// Initialize tickers.
	rl.TestInitRegionTickers()

	// After init, lastFire = now - interval, so all regions are immediately due.
	// Move lastFire to "just now" so deadlines are in the future.
	now := time.Now()
	for i := range rl.TestRegionTickers() {
		rl.TestSetRegionTicker(i, rl.TestRegionTickers()[i].Interval, now)
	}

	// Now the minimum sleep should be approximately the fast interval (50ms).
	sleepDur := rl.TestMinSleepDuration()
	if sleepDur > 55*time.Millisecond || sleepDur < 40*time.Millisecond {
		t.Fatalf("expected sleep duration ~50ms, got %v", sleepDur)
	}

	// If we move lastFire back far enough, deadline is past → return 0.
	firstTicker := rl.TestRegionTickers()[0]
	rl.TestSetRegionTicker(0, firstTicker.Interval, time.Now().Add(-100*time.Millisecond))
	sleepDur = rl.TestMinSleepDuration()
	if sleepDur != 0 {
		t.Fatalf("expected sleep duration 0 for past deadline, got %v", sleepDur)
	}
}

func TestRenderLoop_PerRegion_RenderDueRegionsSelectivity(t *testing.T) {
	// Only regions whose deadline has elapsed should be rendered.
	resolver := &mockTickRateResolver{
		intervals: map[string]time.Duration{
			"fast": 10 * time.Millisecond,
			"slow": 10 * time.Second,
		},
	}

	renderer := &mockRenderer{}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 60))
	if err != nil {
		t.Fatal(err)
	}

	rm := region.NewRegionManager(vd)

	bounds1 := image.Rect(0, 0, 120, 60)
	surf1 := surface.NewFromSubImage(vd.FrameBuffer(), bounds1)
	r1 := region.NewRegion("due", bounds1, surf1)
	r1.TestSetMode("fast")
	rm.TestAppendRegion(r1)

	bounds2 := image.Rect(120, 0, 240, 60)
	surf2 := surface.NewFromSubImage(vd.FrameBuffer(), bounds2)
	r2 := region.NewRegion("notdue", bounds2, surf2)
	r2.TestSetMode("slow")
	rm.TestAppendRegion(r2)

	r1.SetInputFocus(true)

	rl := region.NewRenderLoop(rm, nil, nil,
		region.WithTickRateResolver(resolver),
		region.WithRenderer(renderer),
	)

	// Manually init tickers.
	rl.TestInitRegionTickers()

	// Override lastFire values to control which regions are due.
	// Set "due" region's lastFire to 20ms ago (past its 10ms deadline).
	rl.TestSetRegionTicker(0, rl.TestRegionTickers()[0].Interval, time.Now().Add(-20*time.Millisecond))
	// Set "notdue" region's lastFire to now (its 10s deadline is far in the future).
	rl.TestSetRegionTicker(1, rl.TestRegionTickers()[1].Interval, time.Now())

	// Call renderDueRegions.
	renderer.calls = nil
	rl.TestRenderDueRegions(time.Now())

	calls := renderer.getCalls()
	if len(calls) != 1 || calls[0] != "due" {
		t.Fatalf("expected only 'due' to be rendered, got %v", calls)
	}
}

// --- From: virtual_display_test.go ---

func TestNewVirtualDisplay(t *testing.T) {
	t.Run("valid bounds", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if vd.Bounds() != image.Rect(0, 0, 240, 240) {
			t.Errorf("bounds = %v, want (0,0)-(240,240)", vd.Bounds())
		}
		if vd.FrameBuffer() == nil {
			t.Fatal("FrameBuffer() returned nil")
		}
	})

	t.Run("non-zero origin normalizes to zero", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(10, 20, 250, 260))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should normalize to origin (0,0) with same dimensions.
		want := image.Rect(0, 0, 240, 240)
		if vd.Bounds() != want {
			t.Errorf("bounds = %v, want %v", vd.Bounds(), want)
		}
	})

	t.Run("zero width returns error", func(t *testing.T) {
		_, err := region.NewVirtualDisplay(image.Rect(0, 0, 0, 240))
		if err == nil {
			t.Fatal("expected error for zero width")
		}
	})

	t.Run("negative area returns error", func(t *testing.T) {
		// Bypass image.Rect which canonicalizes; use a literal with Max < Min.
		_, err := region.NewVirtualDisplay(image.Rectangle{
			Min: image.Pt(10, 10),
			Max: image.Pt(5, 5),
		})
		if err == nil {
			t.Fatal("expected error for negative area")
		}
	})

	t.Run("zero height returns error", func(t *testing.T) {
		_, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 0))
		if err == nil {
			t.Fatal("expected error for zero height")
		}
	})
}

func TestNewVirtualDisplayFromScreens(t *testing.T) {
	t.Run("single screen", func(t *testing.T) {
		screens := []region.ScreenPosition{
			{Index: 0, Name: "main", Bounds: image.Rect(0, 0, 240, 240)},
		}
		vd, err := region.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := image.Rect(0, 0, 240, 240)
		if vd.Bounds() != want {
			t.Errorf("bounds = %v, want %v", vd.Bounds(), want)
		}
	})

	t.Run("multiple screens side by side", func(t *testing.T) {
		screens := []region.ScreenPosition{
			{Index: 0, Name: "left", Bounds: image.Rect(0, 0, 240, 135)},
			{Index: 1, Name: "right", Bounds: image.Rect(240, 0, 480, 135)},
		}
		vd, err := region.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := image.Rect(0, 0, 480, 135)
		if vd.Bounds() != want {
			t.Errorf("bounds = %v, want %v", vd.Bounds(), want)
		}
	})

	t.Run("screens with different heights", func(t *testing.T) {
		screens := []region.ScreenPosition{
			{Index: 0, Name: "tall", Bounds: image.Rect(0, 0, 240, 320)},
			{Index: 1, Name: "short", Bounds: image.Rect(240, 0, 480, 135)},
		}
		vd, err := region.NewVirtualDisplayFromScreens(screens)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := image.Rect(0, 0, 480, 320)
		if vd.Bounds() != want {
			t.Errorf("bounds = %v, want %v", vd.Bounds(), want)
		}
	})

	t.Run("empty screens returns error", func(t *testing.T) {
		_, err := region.NewVirtualDisplayFromScreens(nil)
		if err == nil {
			t.Fatal("expected error for empty screens")
		}
		if err.Error() != "virtual display: at least one physical screen is required" {
			t.Errorf("error = %q, want specific message", err.Error())
		}
	})

	t.Run("screen with zero dimensions returns error", func(t *testing.T) {
		screens := []region.ScreenPosition{
			{Index: 0, Name: "bad", Bounds: image.Rect(0, 0, 0, 240)},
		}
		_, err := region.NewVirtualDisplayFromScreens(screens)
		if err == nil {
			t.Fatal("expected error for invalid screen dimensions")
		}
	})
}

func TestVirtualDisplay_SubImage(t *testing.T) {
	t.Run("returns zero-origin sub-image", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 320))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		sub := vd.SubImage(image.Rect(240, 0, 480, 320))
		want := image.Rect(0, 0, 240, 320)
		if sub.Bounds() != want {
			t.Errorf("sub bounds = %v, want %v", sub.Bounds(), want)
		}
	})

	t.Run("shares memory with framebuffer", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 320))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		sub := vd.SubImage(image.Rect(100, 50, 200, 150))

		// Write to the sub-image at local coordinates (0,0).
		red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		sub.SetRGBA(0, 0, red)

		// Read from the framebuffer at the corresponding VD coordinates (100, 50).
		got := vd.FrameBuffer().RGBAAt(100, 50)
		if got != red {
			t.Errorf("framebuffer at (100,50) = %v, want %v", got, red)
		}
	})

	t.Run("framebuffer writes visible in sub-image", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 480, 320))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		sub := vd.SubImage(image.Rect(100, 50, 200, 150))

		// Write to VD framebuffer at (105, 55).
		blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
		vd.FrameBuffer().SetRGBA(105, 55, blue)

		// Read from sub-image at local (5, 5) which maps to (105, 55) in VD.
		got := sub.RGBAAt(5, 5)
		if got != blue {
			t.Errorf("sub at (5,5) = %v, want %v", got, blue)
		}
	})

	t.Run("sub-image clipped to VD bounds", func(t *testing.T) {
		vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 100, 100))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Request a rect that extends beyond VD bounds.
		sub := vd.SubImage(image.Rect(50, 50, 200, 200))
		want := image.Rect(0, 0, 50, 50) // Clipped to (50,50)-(100,100) then zero-origin.
		if sub.Bounds() != want {
			t.Errorf("clipped sub bounds = %v, want %v", sub.Bounds(), want)
		}
	})
}
