package led

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// ---------------------------------------------------------------------------
// Widget Lifecycle: New, Configure, RenderFrame, Tick
// ---------------------------------------------------------------------------

// TestWidgetLifecycle_NewAndRenderFrame verifies that New returns a non-nil
// Renderable and that RenderFrame produces a valid Sprite.

func TestWidgetLifecycle_NewAndRenderFrame(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Bounds:     image.Rect(5, 10, 15, 20),
	}

	w := New(cfg)
	if w == nil {
		t.Fatal("New returned nil")
	}

	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil for valid config")
	}
	if sprite.Image == nil {
		t.Fatal("Sprite.Image is nil")
	}
	if sprite.Position != image.Pt(5, 10) {
		t.Errorf("Position: got %v, want (5,10)", sprite.Position)
	}
}

// TestWidgetLifecycle_Configure verifies that Configure updates the widget state
// and subsequent RenderFrame uses the new config.

func TestWidgetLifecycle_Configure(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
	}
	w := New(cfg)

	// First render: On state
	sprite1 := w.RenderFrame()
	if sprite1 == nil {
		t.Fatal("first RenderFrame returned nil")
	}
	if sprite1.Label != "led/on" {
		t.Errorf("first render label: got %q, want %q", sprite1.Label, "led/on")
	}

	// Configure to Off state
	newCfg := Config{
		State:      Off,
		Brightness: -1.0,
		Diameter:   10,
	}
	w.(widgets.Configurable).Configure(newCfg)

	sprite2 := w.RenderFrame()
	if sprite2 == nil {
		t.Fatal("second RenderFrame returned nil after Configure")
	}
	if sprite2.Label != "led/off" {
		t.Errorf("second render label: got %q, want %q", sprite2.Label, "led/off")
	}
}

// ---------------------------------------------------------------------------
// Caching Tests
// ---------------------------------------------------------------------------

// TestCaching_SameConfigReturnsCachedSprite verifies that rendering the same
// config twice returns the same pointer (cached).

func TestCaching_SameConfigReturnsCachedSprite(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
	}

	w := New(cfg, widgets.WithCaching())

	sprite1 := w.RenderFrame()
	sprite2 := w.RenderFrame()

	if sprite1 == nil || sprite2 == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// Same config → same pointer (cached)
	if sprite1 != sprite2 {
		t.Error("expected same sprite pointer for unchanged config (cache hit)")
	}
}

// TestCaching_DifferentConfigReRenders verifies that a config change produces
// a new sprite (cache miss).

func TestCaching_DifferentConfigReRenders(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
	}

	w := New(cfg, widgets.WithCaching())

	sprite1 := w.RenderFrame()
	if sprite1 == nil {
		t.Fatal("first RenderFrame returned nil")
	}

	// Change config
	newCfg := Config{
		State:      Off,
		Brightness: -1.0,
		Diameter:   10,
	}
	w.(widgets.Configurable).Configure(newCfg)

	sprite2 := w.RenderFrame()
	if sprite2 == nil {
		t.Fatal("second RenderFrame returned nil")
	}

	// Different config → different pointer (re-rendered)
	if sprite1 == sprite2 {
		t.Error("expected different sprite pointer after config change (cache miss)")
	}
}

// ---------------------------------------------------------------------------
// Animation State Machine (Tick)
// ---------------------------------------------------------------------------

// TestTick_AdvancesAnimElapsed verifies that Tick advances animElapsed when
// animation is configured.

func TestTick_AdvancesAnimElapsed(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Animation: AnimationConfig{
			Type:   Pulse,
			Period: 1000 * time.Millisecond,
		},
	}

	w := New(cfg)
	animated := w.(widgets.Animated)

	// Tick 100ms
	animated.Tick(100 * time.Millisecond)

	// Render after Tick — should not crash and should produce a valid sprite
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil after Tick")
	}

	// The sprite should have a label reflecting the animated brightness
	// (pulse at 100ms into a 1000ms period starts at min brightness)
	if sprite.Label == "" {
		t.Error("expected non-empty label after Tick")
	}
}

// TestTick_NoAnimationSkipsTick verifies that Tick is a no-op when no animation
// is configured.

func TestTick_NoAnimationSkipsTick(t *testing.T) {
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Animation: AnimationConfig{
			Type: NoAnimation,
		},
	}

	w := New(cfg, widgets.WithCaching())
	animated := w.(widgets.Animated)

	sprite1 := w.RenderFrame()

	// Tick should be a no-op
	animated.Tick(100 * time.Millisecond)

	sprite2 := w.RenderFrame()

	// With no animation, Tick doesn't advance animElapsed, so the config
	// signature is unchanged → cache hit → same pointer
	if sprite1 != sprite2 {
		t.Error("expected same sprite pointer when NoAnimation (Tick should be no-op)")
	}
}

// ---------------------------------------------------------------------------
// Label Strings
// ---------------------------------------------------------------------------

// TestLabelStrings verifies that the correct label is assigned for each state.

func TestLabelStrings(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected string
	}{
		{"On state", On, "led/on"},
		{"Off state", Off, "led/off"},
		{"Warning state", Warning, "led/warning"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				State:      tt.state,
				Brightness: -1.0,
				Diameter:   10,
			}
			w := New(cfg)
			sprite := w.RenderFrame()
			if sprite == nil {
				t.Fatal("RenderFrame returned nil")
			}
			if sprite.Label != tt.expected {
				t.Errorf("label: got %q, want %q", sprite.Label, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Warning State Default Amber Color
// ---------------------------------------------------------------------------

// TestWarningStateDefaultAmber verifies that Warning state with zero WarningColor
// uses default amber (255, 191, 0, 255) for the center pixel.

func TestWarningStateDefaultAmber(t *testing.T) {
	cfg := Config{
		State:      Warning,
		Brightness: -1.0,
		Diameter:   10,
		// Zero WarningColor → default amber
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// Check center pixel — should be amber (255, 191, 0, 255)
	expectedAmber := color.RGBA{R: 255, G: 191, B: 0, A: 255}
	cx, cy := 5, 5
	r, g, b, a := sprite.Image.At(cx, cy).RGBA()
	got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}

	if got != expectedAmber {
		t.Errorf("center pixel (%d,%d): got %v, want %v (default amber)", cx, cy, got, expectedAmber)
	}
}

// ---------------------------------------------------------------------------
// Shine Style Tests
// ---------------------------------------------------------------------------

// TestShineStyle_DotAtDiameter10 verifies that ShineDot at diameter 10 produces
// a white pixel at the expected shine dot position.

func TestShineStyle_DotAtDiameter10(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Foreground: fg,
		ShineStyle: ShineDot,
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// For diameter 10, bodyRadius = 5.
	// Dot center offset = floor(5 * 0.25) = 1 pixel from center.
	// Body center = (5, 5). Dot center = (5-1, 5-1) = (4, 4).
	// The dot pixel at (4, 4) should be white (255, 255, 255, 255).
	r, g, b, a := sprite.Image.At(4, 4).RGBA()
	got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if got != white {
		t.Errorf("shine dot at (4,4): got %v, want %v (white)", got, white)
	}
}

// TestShineStyle_CrescentAtDiameter20 verifies that ShineCrescent at diameter 20
// produces white pixels in the upper-left arc region.

func TestShineStyle_CrescentAtDiameter20(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   20,
		Foreground: fg,
		ShineStyle: ShineCrescent,
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil")
	}

	// For diameter 20, bodyRadius = 10.
	// Crescent: inner = 10*0.8 = 8, outer = 10*0.9 = 9
	// Body center = (10, 10).
	// A pixel in the upper-left quadrant at correct distance should be white.
	// Check pixel at (2, 10) — dx = 2.5-10 = -7.5, dy = 10.5-10 = 0.5
	// dist = sqrt(56.25 + 0.25) = 7.52... Not in [8, 9] range.
	// Try pixel at (1, 10) — dx = 1.5-10 = -8.5, dy = 10.5-10 = 0.5
	// dist = sqrt(72.25 + 0.25) = 8.52. dx < 0, dy > 0 → doesn't meet dy <= 0.
	// Try pixel at (1, 9) — dx = 1.5-10 = -8.5, dy = 9.5-10 = -0.5
	// dist = sqrt(72.25 + 0.25) = 8.52. dx < 0, dy < 0 → in crescent!
	// This should be white.

	foundCrescentPixel := false
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Scan the upper-left quadrant for a white pixel
	for py := 0; py < 10; py++ {
		for px := 0; px < 10; px++ {
			r, g, b, a := sprite.Image.At(px, py).RGBA()
			got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			if got == white {
				foundCrescentPixel = true
				break
			}
		}
		if foundCrescentPixel {
			break
		}
	}

	if !foundCrescentPixel {
		t.Error("expected at least one white crescent pixel in upper-left quadrant for diameter 20")
	}
}

// TestShineStyle_NoneProducesNoShine verifies that ShineNone does not produce
// any white pixels in the body.

func TestShineStyle_NoneProducesNoShine(t *testing.T) {
	fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}
	cfg := Config{
		State:      On,
		Brightness: -1.0,
		Diameter:   10,
		Foreground: fg,
		ShineStyle: ShineNone,
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil")
	}

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	bounds := sprite.Image.Bounds()
	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			r, g, b, a := sprite.Image.At(px, py).RGBA()
			got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			if got == white {
				t.Fatalf("found white pixel at (%d,%d) with ShineNone", px, py)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Group Layout Tests
// ---------------------------------------------------------------------------

// TestGroupLayout_HorizontalDimensions verifies that a 3-entry horizontal group
// with diameter=10 and spacing=2 produces an output of width=34, height=10.

func TestGroupLayout_HorizontalDimensions(t *testing.T) {
	cfg := Config{
		State:       On,
		Brightness:  -1.0,
		Diameter:    10,
		Spacing:     2,
		Orientation: Horizontal,
		Group: []GroupEntry{
			{State: On},
			{State: Off},
			{State: Warning},
		},
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil for group config")
	}

	bounds := sprite.Image.Bounds()
	expectedWidth := 3*10 + 2*2 // 34
	expectedHeight := 10

	if bounds.Dx() != expectedWidth {
		t.Errorf("group width: got %d, want %d", bounds.Dx(), expectedWidth)
	}
	if bounds.Dy() != expectedHeight {
		t.Errorf("group height: got %d, want %d", bounds.Dy(), expectedHeight)
	}
}

// TestGroupLayout_VerticalDimensions verifies that a 3-entry vertical group
// with diameter=10 and spacing=2 produces an output of width=10, height=34.

func TestGroupLayout_VerticalDimensions(t *testing.T) {
	cfg := Config{
		State:       On,
		Brightness:  -1.0,
		Diameter:    10,
		Spacing:     2,
		Orientation: Vertical,
		Group: []GroupEntry{
			{State: On},
			{State: Off},
			{State: Warning},
		},
	}

	w := New(cfg)
	sprite := w.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame returned nil for vertical group config")
	}

	bounds := sprite.Image.Bounds()
	expectedWidth := 10
	expectedHeight := 3*10 + 2*2 // 34

	if bounds.Dx() != expectedWidth {
		t.Errorf("group width: got %d, want %d", bounds.Dx(), expectedWidth)
	}
	if bounds.Dy() != expectedHeight {
		t.Errorf("group height: got %d, want %d", bounds.Dy(), expectedHeight)
	}
}

// ---------------------------------------------------------------------------
// Off-State Dimming Formula
// ---------------------------------------------------------------------------

// TestOffStateDimmingFormula verifies the dimming formula at specific foreground
// values: each channel becomes floor(channel * 0.3).

func TestOffStateDimmingFormula(t *testing.T) {
	tests := []struct {
		name     string
		fg       color.RGBA
		expected color.RGBA // dimmed outline color
	}{
		{
			"red 100",
			color.RGBA{R: 100, G: 0, B: 0, A: 255},
			color.RGBA{R: 30, G: 0, B: 0, A: 255},
		},
		{
			"green 200",
			color.RGBA{R: 0, G: 200, B: 0, A: 255},
			color.RGBA{R: 0, G: 60, B: 0, A: 255},
		},
		{
			"white 255",
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
			color.RGBA{R: 76, G: 76, B: 76, A: 255},
		},
		{
			"mixed",
			color.RGBA{R: 100, G: 200, B: 50, A: 200},
			color.RGBA{R: 30, G: 60, B: 15, A: 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				State:      Off,
				Brightness: -1.0,
				Diameter:   10,
				Foreground: tt.fg,
				Background: color.RGBA{R: 10, G: 10, B: 10, A: 255},
			}

			w := New(cfg)
			sprite := w.RenderFrame()
			if sprite == nil {
				t.Fatal("RenderFrame returned nil")
			}

			// Find an outline pixel. For a circle with diameter 10, center=5.0,
			// radius=5.0. Pixel (5, 0) is an edge pixel.
			r, g, b, a := sprite.Image.At(5, 0).RGBA()
			got := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}

			if got != tt.expected {
				t.Errorf("outline pixel (5,0): got %v, want %v", got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Registry Registration
// ---------------------------------------------------------------------------

// TestRegistryRegistration verifies that the LED widget is registered as "led"
// in the global widget registry during init().

func TestRegistryRegistration(t *testing.T) {
	names := widgets.Registered()

	found := false
	for _, name := range names {
		if name == "led" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected 'led' to be registered in widget registry, got names: %v", names)
	}
}
