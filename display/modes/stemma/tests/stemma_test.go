package tests_test

import (
	"image"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/stemma"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// --- From: scanner_test.go ---

func TestScannerNew_defaults(t *testing.T) {
	s := source.New(nil, 0)
	if s == nil {
		t.Fatal("expected non-nil Scanner")
	}
	// Before any scan, Devices() should return an empty slice.
	devs := s.Devices()
	if len(devs) != 0 {
		t.Fatalf("expected 0 devices before scan, got %d", len(devs))
	}
}

func TestScannerPresentDevices_empty(t *testing.T) {
	s := source.New([]string{"/dev/i2c-99"}, 0)
	// Bus doesn't exist; PresentDevices must still return a slice (possibly empty).
	present := s.PresentDevices()
	if present == nil {
		t.Fatal("PresentDevices must not return nil")
	}
}

func TestDeviceKey(t *testing.T) {
	// scanAll on a non-existent bus must not panic.
	s := source.New([]string{"/dev/i2c-99"}, 0)
	s.Start()
	s.Stop()
	// After Stop, Devices() must still work.
	_ = s.Devices()
}

// --- From: stemma_test.go ---

// fakeIcon returns a 1x1 image for testing icon resolution.
func fakeIcon(name string) (image.Image, bool) {
	switch name {
	case "check", "error":
		return image.NewRGBA(image.Rect(0, 0, 8, 8)), true
	default:
		return nil, false
	}
}

func TestBuildItems(t *testing.T) {
	devs := []*source.Device{
		{Bus: "/dev/i2c-1", Addr: 0x76, Name: "BME280", Present: true},
		{Bus: "/dev/i2c-3", Addr: 0x3C, Name: "OLED", Present: false},
	}
	items := stemma.BuildItems(devs)
	if len(items) != 2 {
		t.Fatalf("BuildItems() len=%d, want 2", len(items))
	}
	if !strings.Contains(items[0], "/dev/i2c-1") || !strings.Contains(items[1], "/dev/i2c-3") {
		t.Fatalf("BuildItems() rows missing bus data: %v", items)
	}
}

func TestBuildItems_empty(t *testing.T) {
	items := stemma.BuildItems(nil)
	if len(items) != 1 || items[0] != "(no devices found)" {
		t.Fatalf("BuildItems(nil) = %v, want [(no devices found)]", items)
	}
}

func TestBuildColors(t *testing.T) {
	devs := []*source.Device{{Present: true}, {Present: false}}
	cols := stemma.BuildColors(devs, stemma.ColorPresent, stemma.ColorAbsent)
	if len(cols) != 2 {
		t.Fatalf("BuildColors() len=%d, want 2", len(cols))
	}
	if cols[0] != stemma.ColorPresent || cols[1] != stemma.ColorAbsent {
		t.Fatalf("BuildColors() unexpected values: %v", cols)
	}
}

func TestBuildColors_empty(t *testing.T) {
	cols := stemma.BuildColors(nil, stemma.ColorPresent, stemma.ColorAbsent)
	if cols != nil {
		t.Fatalf("BuildColors(nil) = %v, want nil", cols)
	}
}

func TestPolicy_default(t *testing.T) {
	stemma.SetPolicy(stemma.DefaultPolicy())
	p := stemma.GetPolicy()
	if p.Style != "" {
		t.Fatalf("DefaultPolicy().Style = %q, want \"\"", p.Style)
	}
}

func TestPolicy_normalization(t *testing.T) {
	stemma.SetPolicy(stemma.Policy{Style: "INVALID"})
	p := stemma.GetPolicy()
	if p.Style != "" {
		t.Fatalf("SetPolicy invalid style -> Style=%q, want \"\"", p.Style)
	}
}

func TestHandleCommand_query(t *testing.T) {
	stemma.SetPolicy(stemma.DefaultPolicy())
	resp := stemma.HandleCommand(nil)
	if !strings.HasPrefix(resp, "OK stemma") {
		t.Fatalf("HandleCommand(nil) = %q, want prefix \"OK stemma\"", resp)
	}
	if !strings.Contains(resp, "style=") {
		t.Fatalf("HandleCommand(nil) = %q, want to contain \"style=\"", resp)
	}
}

func TestHandleCommand_set_valid(t *testing.T) {
	stemma.SetPolicy(stemma.DefaultPolicy())
	resp := stemma.HandleCommand([]string{"style=mono-slow-128x64"})
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("HandleCommand set valid = %q, want prefix \"OK\"", resp)
	}
	p := stemma.GetPolicy()
	if p.Style != "mono-slow-128x64" {
		t.Fatalf("after set style=mono-slow-128x64, Policy.Style = %q", p.Style)
	}
}

func TestHandleCommand_invalid_value(t *testing.T) {
	stemma.SetPolicy(stemma.DefaultPolicy())
	resp := stemma.HandleCommand([]string{"style=badval"})
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("HandleCommand invalid value = %q, want prefix \"ERR\"", resp)
	}
	// Policy should be unchanged.
	p := stemma.GetPolicy()
	if p.Style != "" {
		t.Fatalf("after invalid set, Policy.Style = %q, want \"\"", p.Style)
	}
}

func TestHandleCommand_unknown_key(t *testing.T) {
	resp := stemma.HandleCommand([]string{"bogus=val"})
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("HandleCommand unknown key = %q, want prefix \"ERR\"", resp)
	}
	if !strings.Contains(resp, "bogus") {
		t.Fatalf("HandleCommand unknown key = %q, want to contain \"bogus\"", resp)
	}
}

func TestBuildSprites(t *testing.T) {
	devs := []*source.Device{
		{Addr: 0x76, Name: "BME280", Present: true},
		{Addr: 0x3C, Name: "OLED", Present: false},
	}
	sprites := stemma.BuildSprites(devs, 10, fakeIcon)
	if len(sprites) != 2 {
		t.Fatalf("BuildSprites() len=%d, want 2", len(sprites))
	}
	// First device present -> check icon at column 0.
	if sprites[0].Label == "" {
		t.Errorf("sprites[0].Label is empty, want non-empty")
	}
	if sprites[0].Position.X != 0 {
		t.Errorf("sprites[0].Position.X = %d, want 0", sprites[0].Position.X)
	}
	if sprites[0].Position.Y != 0 {
		t.Errorf("sprites[0].Position.Y = %d, want 0", sprites[0].Position.Y)
	}
	// Second device absent -> error icon at row 1.
	if sprites[1].Label == "" {
		t.Errorf("sprites[1].Label is empty, want non-empty")
	}
	if sprites[1].Position.Y != 10 {
		t.Errorf("sprites[1].Position.Y = %d, want 10", sprites[1].Position.Y)
	}
}

func TestBuildSprites_empty(t *testing.T) {
	sprites := stemma.BuildSprites(nil, 10, fakeIcon)
	if sprites != nil {
		t.Fatalf("BuildSprites(nil) = %v, want nil", sprites)
	}
}

func TestBuildSprites_noIconResolver(t *testing.T) {
	devs := []*source.Device{{Addr: 0x76, Name: "BME280", Present: true}}
	sprites := stemma.BuildSprites(devs, 10, nil)
	if len(sprites) != 0 {
		t.Fatalf("BuildSprites with nil getIcon = %d sprites, want 0", len(sprites))
	}
}

func TestBuildView_noDevices(t *testing.T) {
	stemma.SetPolicy(stemma.DefaultPolicy())

	hints := textlayout.TextHints{
		PixelWidth:   128,
		PixelHeight:  64,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	result := stemma.BuildView(nil, hints, fakeIcon, nil)
	if len(result.Items) != 1 || result.Items[0] != "(no devices found)" {
		t.Fatalf("BuildView no devices Items = %v", result.Items)
	}
	if result.Colors != nil {
		t.Fatalf("BuildView no devices Colors = %v, want nil", result.Colors)
	}
	if result.Sprites != nil {
		t.Fatalf("BuildView no devices Sprites = %v, want nil", result.Sprites)
	}
}

func TestBuildView_compact(t *testing.T) {
	stemma.SetPolicy(stemma.Policy{Style: "compact"})

	devs := []*source.Device{
		{Addr: 0x76, Name: "BME280", Present: true},
		{Addr: 0x3C, Name: "OLED", Present: false},
		{Addr: 0x44, Name: "SHT30", Present: true},
	}
	hints := textlayout.TextHints{
		PixelWidth:   128,
		PixelHeight:  64,
		GlyphAdvance: 6,
		RowHeight:    10,
	}
	result := stemma.BuildView(devs, hints, fakeIcon, nil)
	if len(result.Items) != 1 {
		t.Fatalf("BuildView compact Items len=%d, want 1", len(result.Items))
	}
	if result.Items[0] != "2/3 present" {
		t.Fatalf("BuildView compact Items[0] = %q, want \"2/3 present\"", result.Items[0])
	}
	if result.Colors != nil {
		t.Fatalf("BuildView compact Colors = %v, want nil", result.Colors)
	}
	if result.Sprites != nil {
		t.Fatalf("BuildView compact Sprites = %v, want nil", result.Sprites)
	}
}

func TestBuildView_list(t *testing.T) {
	stemma.SetPolicy(stemma.Policy{Style: "list"})

	devs := []*source.Device{
		{Addr: 0x76, Name: "BME280", Present: true},
		{Addr: 0x3C, Name: "OLED", Present: false},
	}
	hints := textlayout.TextHints{
		PixelWidth:   128,
		PixelHeight:  64,
		GlyphAdvance: 6,
		RowHeight:    10,
		GlyphWidth:   5,
		GlyphHeight:  7,
	}
	result := stemma.BuildView(devs, hints, fakeIcon, nil)
	if len(result.Items) != 2 {
		t.Fatalf("BuildView list Items len=%d, want 2", len(result.Items))
	}
	if len(result.Colors) != 2 {
		t.Fatalf("BuildView list Colors len=%d, want 2", len(result.Colors))
	}
	// Present device gets green, absent gets red.
	if result.Colors[0] != stemma.ColorPresent {
		t.Errorf("Colors[0] = %v, want ColorPresent", result.Colors[0])
	}
	if result.Colors[1] != stemma.ColorAbsent {
		t.Errorf("Colors[1] = %v, want ColorAbsent", result.Colors[1])
	}
	// Sprites should be present.
	if len(result.Sprites) != 2 {
		t.Fatalf("BuildView list Sprites len=%d, want 2", len(result.Sprites))
	}
}

func TestBuildView_list_truncation(t *testing.T) {
	stemma.SetPolicy(stemma.Policy{Style: "list"})

	devs := []*source.Device{
		{Addr: 0x76, Name: "A very long device name that exceeds normal width", Present: true},
	}
	// Very narrow display: 48px wide with 6px advance = 8 chars max.
	hints := textlayout.TextHints{
		PixelWidth:   48,
		PixelHeight:  64,
		GlyphAdvance: 6,
		RowHeight:    10,
		GlyphWidth:   5,
		GlyphHeight:  7,
	}
	result := stemma.BuildView(devs, hints, fakeIcon, nil)
	if len(result.Items) != 1 {
		t.Fatalf("BuildView truncation Items len=%d, want 1", len(result.Items))
	}
	// MaxCharsPerRow = 48/6 = 8. Icon width = ceil(8/6) = 2 glyph cells.
	// Name width = 8 - 2 = 6 chars.
	maxLen := 6
	if len([]rune(result.Items[0])) > maxLen {
		t.Errorf("BuildView truncation Items[0] len=%d runes, want <= %d", len([]rune(result.Items[0])), maxLen)
	}
}
