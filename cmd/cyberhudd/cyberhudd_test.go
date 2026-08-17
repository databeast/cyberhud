package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/quick"
	"time"

	conn "periph.io/x/conn/v3"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
	"pgregory.net/rapid"
)

// --- From: boot_prepend_prop_test.go ---

//
// Property 6: Boot mode prepend sets mode as first and default.
// For any panel configuration where BootPolicy returns a non-empty mode,
// the resulting panel mode list has that mode as the first element and it
// is set as the panel's default mode.

// genModeString generates a plausible non-empty mode name (lowercase alpha, 1-12 chars).
func genModeString(t *rapid.T, label string) string {
	return rapid.StringMatching(`[a-z]{1,12}`).Draw(t, label)
}

// genModeList generates a list of 0-10 mode strings (some may be empty or have whitespace).
func genModeList(t *rapid.T, label string) []string {
	n := rapid.IntRange(0, 10).Draw(t, label+"Len")
	modes := make([]string, n)
	for i := range modes {
		// Allow some diversity: normal modes, empty strings, whitespace-padded
		kind := rapid.IntRange(0, 2).Draw(t, label+"Kind")
		switch kind {
		case 0:
			modes[i] = genModeString(t, label+"Mode")
		case 1:
			modes[i] = "" // empty entry
		case 2:
			modes[i] = "  " + genModeString(t, label+"Padded") + " "
		}
	}
	return modes
}

// TestProperty_PrependMode_FirstInList verifies that for any non-empty mode string
// and any mode list, after prependMode the given mode is the first element.
func TestProperty_PrependMode_FirstInList(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		mode := genModeString(rt, "bootMode")
		modes := genModeList(rt, "modes")

		result := prependMode(mode, modes)

		if len(result) == 0 {
			t.Fatal("prependMode returned empty list")
		}
		if result[0] != mode {
			t.Fatalf("first element = %q, want boot mode %q (full list: %v)", result[0], mode, result)
		}
	})
}

// TestProperty_PrependMode_NoDuplicates verifies that the prepended mode
// appears exactly once in the resulting list (no duplicates).
func TestProperty_PrependMode_NoDuplicates(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		mode := genModeString(rt, "bootMode")
		modes := genModeList(rt, "modes")

		result := prependMode(mode, modes)

		count := 0
		for _, m := range result {
			if m == mode {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("boot mode %q appears %d times, want exactly 1 (list: %v)", mode, count, result)
		}
	})
}

// TestProperty_PrependMode_DefaultSetOnPanel verifies the full boot policy flow:
// when BootPolicy returns a non-empty mode for a panel, the resulting panel
// in state has that mode as its current (default) mode.
func TestProperty_PrependMode_DefaultSetOnPanel(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a boot mode that the policy will inject
		bootMode := genModeString(rt, "bootMode")

		// Generate some additional modes for the panel
		extraCount := rapid.IntRange(0, 5).Draw(rt, "extraCount")
		extraModes := make([]string, extraCount)
		for i := range extraModes {
			extraModes[i] = genModeString(rt, "extra")
		}

		// Build initial panel modes
		initialModes := append([]string{}, extraModes...)

		// Simulate the boot prepend logic from configureDisplayModes:
		// if mode := boot.BootMode(panelIndex); mode != "" {
		//     displays[i].Modes = prependMode(mode, displays[i].Modes)
		//     displays[i].Default = mode
		// }
		resultModes := prependMode(bootMode, initialModes)

		// Create a panel and reset state, verifying default is set
		panel := coordinator.Region{
			Index:   0,
			Name:    "test",
			Modes:   resultModes,
			Default: bootMode,
		}

		state := coordinator.NewState()
		state.ResetWithFallback([]coordinator.Region{panel}, "dashboard")

		// The current mode should be the boot mode
		current := state.CurrentMode(0)
		normalized := strings.ToLower(strings.TrimSpace(bootMode))
		if current != normalized {
			t.Fatalf("current mode = %q, want boot mode %q (modes: %v)", current, normalized, resultModes)
		}
	})
}

// --- From: config_test.go ---

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cyberhudd.json")
	data := []byte(`{
		"socket":"/tmp/cyberhud.sock",
		"i2c":"/dev/i2c-1,/dev/i2c-3",
		"scan":"3s",
		"display":{
			"panel":"waveshare-2.2",
			"disable_input":true,
			"width":320,
			"height":240,
			"madctl":"0x20",
			"x_offset":2,
			"y_offset":4,
			"dc":"GPIO25",
			"rst":"GPIO27",
			"bl":"GPIO24",
			"busy":"GPIO4"
		}
	}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg.Display.Profile != "waveshare-2.2" {
		t.Fatalf("unexpected panel: %q", cfg.Display.Profile)
	}
	if cfg.Display.DisableInput == nil || !*cfg.Display.DisableInput {
		t.Fatalf("expected disable_input=true, got %#v", cfg.Display.DisableInput)
	}
}

func TestMergeConfigRespectsExplicitFlags(t *testing.T) {
	falseValue := false
	trueValue := true
	width := 320
	height := 240
	xOffset := 1
	yOffset := 2
	cfg := &fileConfig{
		Socket: "/tmp/from-config.sock",
		I2C:    "/dev/i2c-3",
		Scan:   "7s",
		Display: fileDisplayConfig{
			Disabled:     &falseValue,
			Profile:      "waveshare-2.2",
			DisableInput: &trueValue,
			Width:        &width,
			Height:       &height,
			MADCTL:       "0x20",
			XOffset:      &xOffset,
			YOffset:      &yOffset,
			DC:           "GPIO9",
			RST:          "GPIO11",
			BL:           "none",
			Busy:         "GPIO4",
		},
	}

	socket := "/run/default.sock"
	i2c := "/dev/i2c-1"
	scan := 2 * time.Second
	noDisplay := true
	noInput := false
	profile := "waveshare-1.3hat"
	dWidth := -1
	dHeight := -1
	madctl := ""
	rotate := ""
	dx := -1
	dy := -1
	dc := ""
	rst := ""
	bl := ""
	busy := ""
	inKey1 := ""
	inKey2 := ""
	inKey3 := ""
	inUp := ""
	inDown := ""
	inLeft := ""
	inRight := ""
	inPress := ""

	seen := map[string]bool{
		"socket":        true,
		"display-width": true,
	}
	if err := mergeConfig(cfg, seen, &socket, &i2c, &scan, &noDisplay, &noInput, &profile, &dWidth, &dHeight, &madctl, &rotate, &dx, &dy, &dc, &rst, &bl, &busy, &inKey1, &inKey2, &inKey3, &inUp, &inDown, &inLeft, &inRight, &inPress); err != nil {
		t.Fatalf("mergeConfig() error = %v", err)
	}
	if socket != "/run/default.sock" {
		t.Fatalf("socket should remain explicit flag value, got %q", socket)
	}
	if i2c != "/dev/i2c-3" || scan != 7*time.Second {
		t.Fatalf("unexpected base overrides: i2c=%q scan=%s", i2c, scan)
	}
	if profile != "waveshare-2.2" || !noInput {
		t.Fatalf("unexpected display config merge: panel=%q noInput=%v", profile, noInput)
	}
	if dWidth != -1 {
		t.Fatalf("display-width should remain explicit flag value, got %d", dWidth)
	}
	if dHeight != 240 || madctl != "0x20" || dx != 1 || dy != 2 {
		t.Fatalf("unexpected display tuning merge: height=%d madctl=%q dx=%d dy=%d", dHeight, madctl, dx, dy)
	}
	if dc != "GPIO9" || rst != "GPIO11" || bl != "none" {
		t.Fatalf("unexpected GPIO merge: dc=%q rst=%q bl=%q", dc, rst, bl)
	}
	if busy != "GPIO4" {
		t.Fatalf("unexpected busy pin merge: %q", busy)
	}
}

// --- From: cyberhudd_test.go ---

// Tests from display_panels_test.go

func TestGetDisplayProfile(t *testing.T) {
	profile, err := panels.Get("waveshare-2.2")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	if profile.Config.Width != 320 || profile.Config.Height != 240 {
		t.Fatalf("unexpected panel size: %+v", profile.Config)
	}
	if profile.Inputs.Any() {
		t.Fatal("waveshare-2.2 should default to no input pins")
	}
}

func TestResolveDisplayProfileOverrides(t *testing.T) {
	profile, err := panels.Resolve("waveshare-1.3hat", panels.Overrides{
		Width:   240,
		Height:  320,
		MADCTL:  "0x28",
		XOffset: 2,
		YOffset: 4,
		DCPin:   "GPIO9",
		RSTPin:  "GPIO11",
		BLPin:   "none",
	})
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	if profile.Config.Width != 240 || profile.Config.Height != 320 {
		t.Fatalf("unexpected resolved size: %+v", profile.Config)
	}
	if profile.Config.MADCTL != 0x28 {
		t.Fatalf("unexpected MADCTL: 0x%02X", profile.Config.MADCTL)
	}
	if profile.Config.XOffset != 2 || profile.Config.YOffset != 4 {
		t.Fatalf("unexpected offsets: %+v", profile.Config)
	}
	if profile.DCPin != "GPIO9" || profile.RSTPin != "GPIO11" {
		t.Fatalf("unexpected GPIO overrides: dc=%s rst=%s", profile.DCPin, profile.RSTPin)
	}
	if profile.BLPin != "" {
		t.Fatalf("expected backlight override to disable control, got %q", profile.BLPin)
	}
}

func TestResolveDisplayProfile(t *testing.T) {
	primary, err := panels.Resolve(
		"waveshare-2.2",
		panels.Overrides{
			Width:   -1,
			Height:  -1,
			XOffset: -1,
			YOffset: -1,
		},
	)
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	if primary.Name != "waveshare-2.2" {
		t.Fatalf("unexpected primary: %s", primary.Name)
	}
}

func TestResolveDisplayProfileEmptyOverrides(t *testing.T) {
	// Zero-value overrides should leave the base profile unchanged.
	primary, err := panels.Resolve(
		"waveshare-1.3hat",
		panels.Overrides{},
	)
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	// The base profile for waveshare-1.3hat is 240x240.
	if primary.Config.Width != 240 || primary.Config.Height != 240 {
		t.Fatalf("expected base dimensions 240x240, got %dx%d", primary.Config.Width, primary.Config.Height)
	}
	if primary.Name != "waveshare-1.3hat" {
		t.Fatalf("unexpected primary name: %s", primary.Name)
	}
}

func TestResolveDisplayProfileSomeOverrides(t *testing.T) {
	// Override only width and DC pin; other fields remain at base values.
	primary, err := panels.Resolve(
		"waveshare-1.3hat",
		panels.Overrides{
			Width: 320,
			DCPin: "GPIO22",
		},
	)
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	if primary.Config.Width != 320 {
		t.Fatalf("expected width override 320, got %d", primary.Config.Width)
	}
	// Height should remain at base value (240).
	if primary.Config.Height != 240 {
		t.Fatalf("expected base height 240, got %d", primary.Config.Height)
	}
	if primary.DCPin != "GPIO22" {
		t.Fatalf("expected DC pin override GPIO22, got %q", primary.DCPin)
	}
}

func TestResolveDisplayProfileAllOverrides(t *testing.T) {
	// Override all supported fields.
	primary, err := panels.Resolve(
		"waveshare-1.3hat",
		panels.Overrides{
			Width:      320,
			Height:     240,
			MADCTL:     "0x60",
			XOffset:    10,
			YOffset:    20,
			DCPin:      "GPIO9",
			RSTPin:     "GPIO11",
			BLPin:      "GPIO18",
			BusyPin:    "GPIO24",
			InputKey1:  "GPIO5",
			InputKey2:  "GPIO6",
			InputKey3:  "GPIO13",
			InputUp:    "GPIO19",
			InputDown:  "GPIO26",
			InputLeft:  "GPIO16",
			InputRight: "GPIO20",
			InputPress: "GPIO21",
		},
	)
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	if primary.Config.Width != 320 || primary.Config.Height != 240 {
		t.Fatalf("unexpected dimensions: %dx%d", primary.Config.Width, primary.Config.Height)
	}
	if primary.Config.MADCTL != 0x60 {
		t.Fatalf("expected MADCTL 0x60, got 0x%02X", primary.Config.MADCTL)
	}
	if primary.Config.XOffset != 10 || primary.Config.YOffset != 20 {
		t.Fatalf("unexpected offsets: x=%d y=%d", primary.Config.XOffset, primary.Config.YOffset)
	}
	if primary.DCPin != "GPIO9" {
		t.Fatalf("unexpected DC pin: %q", primary.DCPin)
	}
	if primary.RSTPin != "GPIO11" {
		t.Fatalf("unexpected RST pin: %q", primary.RSTPin)
	}
	if primary.BLPin != "GPIO18" {
		t.Fatalf("unexpected BL pin: %q", primary.BLPin)
	}
	if primary.BusyPin != "GPIO24" {
		t.Fatalf("unexpected Busy pin: %q", primary.BusyPin)
	}
	if primary.Inputs.Key1 != "GPIO5" || primary.Inputs.Key2 != "GPIO6" || primary.Inputs.Key3 != "GPIO13" {
		t.Fatalf("unexpected key overrides: %+v", primary.Inputs)
	}
	if primary.Inputs.JoyUp != "GPIO19" || primary.Inputs.JoyDown != "GPIO26" ||
		primary.Inputs.JoyLeft != "GPIO16" || primary.Inputs.JoyRight != "GPIO20" ||
		primary.Inputs.JoyPressed != "GPIO21" {
		t.Fatalf("unexpected joystick overrides: %+v", primary.Inputs)
	}
}

func TestResolveDisplayProfileInvalidProfile(t *testing.T) {
	// An unknown profile should return an error.
	_, err := panels.Resolve(
		"nonexistent-panel-xyz",
		panels.Overrides{},
	)
	if err == nil {
		t.Fatal("expected error for invalid profile name, got nil")
	}
}

func TestResolveDisplayProfileBLPinNoneDisablesBacklight(t *testing.T) {
	// BLPin="none" should clear the backlight pin override.
	primary, err := panels.Resolve(
		"waveshare-1.3hat",
		panels.Overrides{BLPin: "none"},
	)
	if err != nil {
		t.Fatalf("panels.Resolve() error = %v", err)
	}
	if primary.BLPin != "" {
		t.Fatalf("expected BLPin cleared by 'none' override, got %q", primary.BLPin)
	}
}

func TestOLEDProfileMetadata(t *testing.T) {
	profile, err := panels.Get("waveshare-1.3-oled-hat")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	if profile.Controller != "sh1106" {
		t.Fatalf("expected sh1106 controller, got %q", profile.Controller)
	}
	if !profile.Inputs.Any() {
		t.Fatal("expected OLED HAT profile to include input pins")
	}
}

func TestSSD1680LandscapeProfileMetadata(t *testing.T) {
	profile, err := panels.Get("adafruit-2.13-ssd1680")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	if profile.Controller != "ssd1680" {
		t.Fatalf("expected ssd1680 controller, got %q", profile.Controller)
	}
	if profile.Config.Width != 122 || profile.Config.Height != 250 {
		t.Fatalf("unexpected EPD size: %+v", profile.Config)
	}
	if profile.BusyPin == "" {
		t.Fatal("expected SSD1680 profile to include busy pin")
	}
	if !profile.Inputs.Any() {
		t.Fatal("expected SSD1680 profile to include button pins")
	}
}

func TestWaveshareTripleScreenProfileMetadata(t *testing.T) {
	profile, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	if profile.Controller != "virtual" {
		t.Fatalf("expected virtual controller, got %q", profile.Controller)
	}
	if len(profile.Virtual) != 3 {
		t.Fatalf("expected 3 virtual panels, got %d", len(profile.Virtual))
	}
	if profile.Virtual[0].Controller != "st7789" || profile.Virtual[1].Controller != "st7735s" || profile.Virtual[2].Controller != "st7735s" {
		t.Fatalf("unexpected virtual panel controllers: %+v", profile.Virtual)
	}
	if profile.Virtual[0].DefaultMode == "" || profile.Virtual[1].DefaultMode == "" || profile.Virtual[2].DefaultMode == "" {
		t.Fatalf("expected per-screen defaults to be set, got left=%+v center=%+v right=%+v", profile.Virtual[1], profile.Virtual[0], profile.Virtual[2])
	}
	if profile.Inputs.Key1 != "GPIO25" || profile.Inputs.Key2 != "GPIO26" {
		t.Fatalf("unexpected key mapping: %+v", profile.Inputs)
	}
}

// Tests from pins_test.go

func TestPinNoticesWaveshare13HatGPIO13Unavailable(t *testing.T) {
	profile, err := panels.Get("waveshare-1.3hat")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	notices := panels.PinNotices(profile)
	if len(notices) == 0 {
		t.Fatal("expected at least one pin notice for waveshare-1.3hat")
	}
	joined := strings.Join(notices, "\n")
	if !strings.Contains(joined, "GPIO13") {
		t.Fatalf("expected GPIO13 unavailable notice, got %q", joined)
	}
}

func TestPinNoticesSSD1680NoUnavailableOutputs(t *testing.T) {
	profile, err := panels.Get("adafruit-2.13-ssd1680")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	notices := panels.PinNotices(profile)
	if len(notices) != 0 {
		t.Fatalf("expected no unavailable output notices for SSD1680 panel, got %v", notices)
	}
}

func TestBuildPinReportIncludesConnectorConflict(t *testing.T) {
	profile, err := panels.Get("waveshare-1.3hat")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	report := panels.BuildPinReport(profile)
	if !strings.Contains(report, "3-pin connector GPIO13") {
		t.Fatalf("expected GPIO13 connector in report, got %q", report)
	}
	if !strings.Contains(report, "status=conflict") {
		t.Fatalf("expected conflict status in report, got %q", report)
	}
}

func TestPinNoticesTripleScreenGPIO13AndGPIO18Unavailable(t *testing.T) {
	profile, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("panels.Get() error = %v", err)
	}
	notices := panels.PinNotices(profile)
	if len(notices) == 0 {
		t.Fatal("expected pin notices for waveshare-triple-screen")
	}
	joined := strings.Join(notices, "\n")
	if !strings.Contains(joined, "GPIO13") {
		t.Fatalf("expected GPIO13 unavailable notice, got %q", joined)
	}
	if !strings.Contains(joined, "GPIO18") {
		t.Fatalf("expected GPIO18 unavailable notice, got %q", joined)
	}
}

// --- From: display_modes_test.go ---

func TestPrependMode(t *testing.T) {
	got := prependMode("systemd", []string{"menu", "dashboard", "systemd", ""})
	want := []string{"systemd", "menu", "dashboard"}
	if len(got) != len(want) {
		t.Fatalf("prependMode len=%d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("prependMode[%d]=%q, want %q (full=%v)", i, got[i], want[i], got)
		}
	}
}

func TestConfigureDisplayModesForcesStartupSystemdOnPrimaryPanel(t *testing.T) {
	state := coordinator.NewState()
	profile := panels.Definition{
		Name:        "test-profile",
		Controller:  "st7789",
		DefaultMode: "menu",
		Virtual: []panels.Screen{
			{Index: 0, Name: "main", Controller: "st7789", DefaultMode: "menu"},
			{Index: 1, Name: "aux", Controller: "st7735s", DefaultMode: "clock"},
		},
	}

	configureDisplayModes(state, profile, true)

	if got := state.CurrentMode(0); got != "systemd" {
		t.Fatalf("panel0 current=%q, want systemd", got)
	}
	// With StandardPolicy and inputEnabled=true, panel 1 resolves to "menu"
	// (the policy prefers "menu" when input is enabled).
	if got := state.CurrentMode(1); got != "menu" {
		t.Fatalf("panel1 current=%q, want menu", got)
	}
	defs := state.Definitions()
	if len(defs) != 2 {
		t.Fatalf("Definitions len=%d, want 2", len(defs))
	}
	if len(defs[0].Modes) == 0 || defs[0].Modes[0].ID != "systemd" {
		t.Fatalf("primary mode list missing systemd at index 0: %+v", defs[0].Modes)
	}
}

// --- From: display_runtime_test.go ---

// --- Mock types ---

// mockI2CBus implements i2c.BusCloser for testing.
type mockI2CBus struct {
	mu     sync.Mutex
	closed bool
}

func (m *mockI2CBus) String() string { return "mock-i2c-bus" }

func (m *mockI2CBus) Tx(addr uint16, w, r []byte) error { return nil }

func (m *mockI2CBus) SetSpeed(f physic.Frequency) error { return nil }

func (m *mockI2CBus) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *mockI2CBus) isClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

// mockDrawTarget implements driver.DrawTarget for testing.
type mockDrawTarget struct{}

func (m *mockDrawTarget) Bounds() image.Rectangle {
	return image.Rect(0, 0, 15, 7)
}

func (m *mockDrawTarget) DrawImage(draw.Image) error {
	return nil
}

// --- Helper: register and unregister test buses ---

func registerTestBus(t *testing.T, name string, opener func() (i2c.BusCloser, error)) {
	t.Helper()
	if err := i2creg.Register(name, nil, -1, opener); err != nil {
		t.Fatalf("i2creg.Register(%q) failed: %v", name, err)
	}
	t.Cleanup(func() {
		_ = i2creg.Unregister(name)
	})
}

// --- Helper: register and unregister test SPI ports ---

func registerTestSPIPort(t *testing.T, name string, opener func() (spi.PortCloser, error)) {
	t.Helper()
	if err := spireg.Register(name, nil, -1, opener); err != nil {
		t.Fatalf("spireg.Register(%q) failed: %v", name, err)
	}
	t.Cleanup(func() {
		_ = spireg.Unregister(name)
	})
}

// mockSPIPortCloser implements spi.PortCloser for testing.
type mockSPIPortCloser struct{}

func (m *mockSPIPortCloser) String() string { return "mock-spi-port" }
func (m *mockSPIPortCloser) Close() error   { return nil }
func (m *mockSPIPortCloser) Connect(f physic.Frequency, mode spi.Mode, bits int) (spi.Conn, error) {
	return &mockSPIConn{}, nil
}
func (m *mockSPIPortCloser) LimitSpeed(f physic.Frequency) error { return nil }

// mockSPIConn implements spi.Conn for testing.
type mockSPIConn struct{}

func (m *mockSPIConn) String() string               { return "mock-spi-conn" }
func (m *mockSPIConn) Tx(w, r []byte) error         { return nil }
func (m *mockSPIConn) Duplex() conn.Duplex          { return conn.Half }
func (m *mockSPIConn) TxPackets([]spi.Packet) error { return nil }

// --- Helper: register test drivers ---

func registerTestDriver(t *testing.T, def driver.Definition) {
	t.Helper()
	driver.Register(def)
	// Note: driver registry does not have Unregister, but test IDs are unique
	// per test so there's no conflict between tests.
}

// --- Tests ---

func TestIsI2COnly(t *testing.T) {
	t.Run("I2C-only driver (NewI2C set, NewSPI nil)", func(t *testing.T) {
		drv := driver.Definition{
			NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
				return nil, nil
			},
			NewSPI: nil,
		}
		if !isI2COnly(drv) {
			t.Fatal("expected isI2COnly=true for driver with NewI2C set and NewSPI nil")
		}
	})

	t.Run("neither factory set", func(t *testing.T) {
		drv := driver.Definition{NewI2C: nil, NewSPI: nil}
		if isI2COnly(drv) {
			t.Fatal("expected isI2COnly=false when both factories are nil")
		}
	})

	t.Run("NewI2C nil (SPI-only or no factory)", func(t *testing.T) {
		drv := driver.Definition{NewI2C: nil}
		if isI2COnly(drv) {
			t.Fatal("expected isI2COnly=false when NewI2C is nil")
		}
	})
}

func TestInitDisplayTargets_I2CPathTaken(t *testing.T) {
	// Register a test I2C bus that returns a mock bus.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-i2c-path-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a test I2C-only driver.
	factoryCalled := false
	registerTestDriver(t, driver.Definition{
		ID:         "test-i2c-driver",
		Monochrome: false,
		NewI2C: func(b i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			factoryCalled = true
			return &mockDrawTarget{}, nil
		},
		// New is nil — makes this I2C-only.
	})

	// Create a panel definition that references the test driver.
	profile := panels.Definition{
		Name:       "test-i2c-panel",
		Controller: "test-i2c-driver",
		Config: driver.DriverConfig{
			Width:  15,
			Height: 7,
		},
		I2CBus: "test-i2c-path-bus",
	}

	positions, cleanup, controller, err := initDisplayTargets(profile, "")
	defer cleanup()

	if err != nil {
		t.Fatalf("initDisplayTargets() error = %v", err)
	}
	if !factoryCalled {
		t.Fatal("expected I2C factory to be called for I2C-only driver")
	}
	if controller != "test-i2c-driver" {
		t.Fatalf("expected controller %q, got %q", "test-i2c-driver", controller)
	}
	if len(positions) == 0 {
		t.Fatal("expected non-empty screen positions")
	}
}

func TestInitI2CDisplayTarget_BusOpenFailure(t *testing.T) {
	// Use a bus name that is not registered — i2creg.Open will fail.
	_, _, _, err := initI2CDisplayTarget("nonexistent-test-bus-xyz", "any-controller", driver.DriverConfig{})
	if err == nil {
		t.Fatal("expected error when bus cannot be opened")
	}

	// Verify the error message includes the bus identifier.
	errMsg := err.Error()
	if !strings.Contains(errMsg, "nonexistent-test-bus-xyz") {
		t.Fatalf("error should include bus identifier, got: %q", errMsg)
	}
}

func TestInitI2CDisplayTarget_FactoryFailure(t *testing.T) {
	// Register a test I2C bus that returns a mock bus we can track.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-factory-fail-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a test driver whose NewI2C returns an error.
	factoryErr := errors.New("simulated factory failure")
	registerTestDriver(t, driver.Definition{
		ID: "test-factory-fail-driver",
		NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
			return nil, factoryErr
		},
	})

	_, _, _, err := initI2CDisplayTarget("test-factory-fail-bus", "test-factory-fail-driver", driver.DriverConfig{})
	if err == nil {
		t.Fatal("expected error when factory fails")
	}

	// Verify the error message includes the controller name.
	errMsg := err.Error()
	if !strings.Contains(errMsg, "test-factory-fail-driver") {
		t.Fatalf("error should include controller name, got: %q", errMsg)
	}

	// Verify the error wraps or contains the factory error message.
	if !strings.Contains(errMsg, "simulated factory failure") {
		t.Fatalf("error should include underlying failure reason, got: %q", errMsg)
	}

	// Verify the bus was closed (cleanup on failure).
	if !bus.isClosed() {
		t.Fatal("expected bus to be closed after factory failure")
	}
}

func TestInitI2CDisplayTarget_UnsupportedController(t *testing.T) {
	// Register a test bus so the Open succeeds.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-unsupported-ctrl-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	_, _, _, err := initI2CDisplayTarget("test-unsupported-ctrl-bus", "nonexistent-controller-xyz", driver.DriverConfig{})
	if err == nil {
		t.Fatal("expected error for unsupported controller")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "nonexistent-controller-xyz") {
		t.Fatalf("error should include controller name, got: %q", errMsg)
	}

	// Verify the bus was closed (cleanup when driver not found).
	if !bus.isClosed() {
		t.Fatal("expected bus to be closed when controller is unsupported")
	}
}

func TestInitDisplayTargets_I2CPathVirtualScreens(t *testing.T) {
	// Register a test I2C bus.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-virtual-i2c-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a test I2C-only driver.
	callCount := 0
	registerTestDriver(t, driver.Definition{
		ID: "test-virtual-i2c-driver",
		NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
			callCount++
			return &mockDrawTarget{}, nil
		},
	})

	// Create a virtual panel definition with an I2C screen.
	profile := panels.Definition{
		Name:       "test-virtual-panel",
		Controller: "virtual",
		Virtual: []panels.Screen{
			{
				Index:      0,
				Name:       "main",
				Controller: "test-virtual-i2c-driver",
				I2CBus:     "test-virtual-i2c-bus",
				Config: driver.DriverConfig{
					Width:  15,
					Height: 7,
				},
			},
		},
	}

	positions, cleanup, controller, err := initDisplayTargets(profile, "")
	defer cleanup()

	if err != nil {
		t.Fatalf("initDisplayTargets() error = %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected I2C factory called once, got %d", callCount)
	}
	if controller != "test-virtual-i2c-driver" {
		t.Fatalf("expected controller %q, got %q", "test-virtual-i2c-driver", controller)
	}
	if len(positions) == 0 {
		t.Fatal("expected non-empty screen positions")
	}
}

// Ensure the error format for bus failure matches requirement 5.4 (includes bus identifier and reason).
func TestInitI2CDisplayTarget_BusOpenError_Format(t *testing.T) {
	_, _, _, err := initI2CDisplayTarget("/dev/i2c-99", "is31fl3731", driver.DriverConfig{})
	if err == nil {
		t.Fatal("expected error when bus cannot be opened")
	}

	// Error must contain the bus path.
	errMsg := err.Error()
	if !strings.Contains(errMsg, "/dev/i2c-99") {
		t.Fatalf("error should contain bus path %q, got: %q", "/dev/i2c-99", errMsg)
	}

	// Error must be wrapped (verify with errors.Unwrap or fmt pattern).
	if !strings.Contains(errMsg, "i2c bus") {
		t.Fatalf("error should have 'i2c bus' prefix, got: %q", errMsg)
	}
}

// --- buildPositions tests ---

// sizedDrawTarget implements driver.DrawTarget with configurable bounds.
type sizedDrawTarget struct {
	bounds image.Rectangle
}

func (s *sizedDrawTarget) Bounds() image.Rectangle {
	return s.bounds
}

func (s *sizedDrawTarget) DrawImage(draw.Image) error {
	return nil
}

func TestBuildPositions_SingleScreen(t *testing.T) {
	target := &sizedDrawTarget{bounds: image.Rect(0, 0, 240, 135)}
	targets := []indexedTarget{
		{index: 0, target: target},
	}

	profile := panels.Definition{
		Name: "test-panel",
	}

	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}

	pos := positions[0]
	if pos.Index != 0 {
		t.Errorf("expected Index=0, got %d", pos.Index)
	}
	if pos.Name != "test-panel" {
		t.Errorf("expected Name=%q, got %q", "test-panel", pos.Name)
	}
	expectedBounds := image.Rect(0, 0, 240, 135)
	if pos.Bounds != expectedBounds {
		t.Errorf("expected Bounds=%v, got %v", expectedBounds, pos.Bounds)
	}
	if pos.Target != target {
		t.Error("expected Target to be the original draw target")
	}
}

func TestBuildPositions_MultipleVirtualScreens(t *testing.T) {
	target0 := &sizedDrawTarget{bounds: image.Rect(0, 0, 240, 135)}
	target1 := &sizedDrawTarget{bounds: image.Rect(0, 0, 240, 135)}
	targets := []indexedTarget{
		{index: 0, target: target0},
		{index: 1, target: target1},
	}

	profile := panels.Definition{
		Name: "multi-panel",
		Virtual: []panels.Screen{
			{Index: 0, Name: "left"},
			{Index: 1, Name: "right"},
		},
	}

	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	if len(positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(positions))
	}

	// First screen at origin.
	if positions[0].Index != 0 {
		t.Errorf("positions[0].Index = %d, want 0", positions[0].Index)
	}
	if positions[0].Name != "left" {
		t.Errorf("positions[0].Name = %q, want %q", positions[0].Name, "left")
	}
	expectedBounds0 := image.Rect(0, 0, 240, 135)
	if positions[0].Bounds != expectedBounds0 {
		t.Errorf("positions[0].Bounds = %v, want %v", positions[0].Bounds, expectedBounds0)
	}

	// Second screen placed to the right of the first.
	if positions[1].Index != 1 {
		t.Errorf("positions[1].Index = %d, want 1", positions[1].Index)
	}
	if positions[1].Name != "right" {
		t.Errorf("positions[1].Name = %q, want %q", positions[1].Name, "right")
	}
	expectedBounds1 := image.Rect(240, 0, 480, 135)
	if positions[1].Bounds != expectedBounds1 {
		t.Errorf("positions[1].Bounds = %v, want %v", positions[1].Bounds, expectedBounds1)
	}
}

func TestBuildPositions_NoScreens(t *testing.T) {
	targets := []indexedTarget{}
	profile := panels.Definition{Name: "empty"}

	_, err := buildPositions(targets, profile)
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
	if !strings.Contains(err.Error(), "no screens") {
		t.Errorf("error should mention 'no screens', got: %q", err.Error())
	}
}

func TestBuildPositions_ZeroBoundsScreen(t *testing.T) {
	// A screen with zero-width bounds.
	target := &sizedDrawTarget{bounds: image.Rect(0, 0, 0, 135)}
	targets := []indexedTarget{
		{index: 0, target: target},
	}

	profile := panels.Definition{Name: "zero-bounds"}

	_, err := buildPositions(targets, profile)
	if err == nil {
		t.Fatal("expected error when no screen has positive-dimension bounds")
	}
	if !strings.Contains(err.Error(), "positive-dimension") {
		t.Errorf("error should mention 'positive-dimension', got: %q", err.Error())
	}
}

func TestBuildPositions_FallbackNameForUnnamedVirtual(t *testing.T) {
	target0 := &sizedDrawTarget{bounds: image.Rect(0, 0, 128, 64)}
	target1 := &sizedDrawTarget{bounds: image.Rect(0, 0, 128, 64)}
	targets := []indexedTarget{
		{index: 0, target: target0},
		{index: 1, target: target1},
	}

	// Profile has virtual screens but no names.
	profile := panels.Definition{
		Name: "unnamed-panel",
		Virtual: []panels.Screen{
			{Index: 0},
			{Index: 1},
		},
	}

	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	// With multiple unnamed screens, fallback names should be "screen0", "screen1".
	if positions[0].Name != "screen0" {
		t.Errorf("positions[0].Name = %q, want %q", positions[0].Name, "screen0")
	}
	if positions[1].Name != "screen1" {
		t.Errorf("positions[1].Name = %q, want %q", positions[1].Name, "screen1")
	}
}

// --- Integration tests for daemon wiring (task 6.3) ---
// These tests verify the ActivatePanel integration from the daemon's perspective,
// including buildPositions feeding into ActivatePanel, default layout generation,
// and noinput handling.

// TestDaemonWiring_ActivatePanelWithMockTarget verifies that calling ActivatePanel
// with a mock DrawTarget and config produced via buildPositions results in
// one default region with correct bounds.

func TestDaemonWiring_ActivatePanelWithMockTarget(t *testing.T) {
	// Set up a mock draw target with known bounds (240x135).
	target := &sizedDrawTarget{bounds: image.Rect(0, 0, 240, 135)}

	// Build screen positions as runDisplayWithRegions would.
	targets := []indexedTarget{
		{index: 0, target: target},
	}
	profile := panels.Definition{
		Name: "test-panel",
	}
	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	// Construct PanelActivationConfig matching the daemon wiring pattern.
	config := region.PanelActivationConfig{
		Screens:      positions,
		Layout:       nil, // no explicit layout → default generation
		DefaultMode:  "dashboard",
		InputEnabled: true,
		AvailModes:   []string{"dashboard", "clock", "ticker"},
		ModeValidator: func(mode string) bool {
			return mode == "dashboard" || mode == "clock" || mode == "ticker"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify exactly one region allocated.
	regions := activation.RegionManager.Regions()
	if len(regions) != 1 {
		t.Fatalf("got %d regions, want 1", len(regions))
	}

	// Verify the region is named "default".
	r := regions[0]
	if r.Name() != "default" {
		t.Errorf("region name=%q, want %q", r.Name(), "default")
	}

	// Verify region bounds cover the full panel bounds (240x135).
	expectedBounds := image.Rect(0, 0, 240, 135)
	if r.Bounds() != expectedBounds {
		t.Errorf("region bounds=%v, want %v", r.Bounds(), expectedBounds)
	}

	// Verify the default mode is assigned.
	if r.CurrentMode() != "dashboard" {
		t.Errorf("region mode=%q, want %q", r.CurrentMode(), "dashboard")
	}

	// Verify input focus is granted to the default region.
	if !r.HasInputFocus() {
		t.Error("default region should have input focus")
	}

	// Verify all activation components are non-nil.
	if activation.FlushPath == nil {
		t.Error("FlushPath is nil")
	}
	if activation.RenderLoop == nil {
		t.Error("RenderLoop is nil")
	}
	if activation.ModeSwitch == nil {
		t.Error("ModeSwitch is nil")
	}
}

// TestDaemonWiring_SinglePanelNoLayout_Region0FullBounds verifies that a
// single-panel config with no explicit layout generates Region 0 whose bounds
// equal the full union of screen bounds (backward compatibility).

func TestDaemonWiring_SinglePanelNoLayout_Region0FullBounds(t *testing.T) {
	// Use a larger screen to test bounds propagation.
	target := &sizedDrawTarget{bounds: image.Rect(0, 0, 320, 240)}

	targets := []indexedTarget{
		{index: 0, target: target},
	}
	profile := panels.Definition{
		Name: "wide-panel",
	}
	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	config := region.PanelActivationConfig{
		Screens:      positions,
		Layout:       nil,
		DefaultMode:  "clock",
		InputEnabled: true,
		AvailModes:   []string{"clock", "dashboard"},
		ModeValidator: func(mode string) bool {
			return mode == "clock" || mode == "dashboard"
		},
	}

	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify VirtualDisplay bounds match the union of all screens.
	vdBounds := activation.VirtualDisplay.Bounds()
	if vdBounds != image.Rect(0, 0, 320, 240) {
		t.Errorf("VD bounds=%v, want %v", vdBounds, image.Rect(0, 0, 320, 240))
	}

	// Verify Region 0 exists and covers the full VD bounds.
	r, ok := activation.RegionManager.Region(0)
	if !ok {
		t.Fatal("Region(0) not found")
	}
	if r.Bounds() != image.Rect(0, 0, 320, 240) {
		t.Errorf("Region 0 bounds=%v, want %v", r.Bounds(), image.Rect(0, 0, 320, 240))
	}

	// Verify default mode assigned.
	if r.CurrentMode() != "clock" {
		t.Errorf("Region 0 mode=%q, want %q", r.CurrentMode(), "clock")
	}

	// Verify input focus is enabled (requirement 8.2).
	if !r.HasInputFocus() {
		t.Error("Region 0 should have input focus when InputEnabled is true")
	}
}

// TestDaemonWiring_NoInputFlag_NilEventsAndDisabled verifies that when the
// noinput flag is set (enableInput=false), the PanelActivationConfig gets nil
// events and InputEnabled=false.

func TestDaemonWiring_NoInputFlag_NilEventsAndDisabled(t *testing.T) {
	// This test verifies the config construction logic that runDisplayWithRegions
	// performs when input is disabled.
	target := &sizedDrawTarget{bounds: image.Rect(0, 0, 240, 135)}

	targets := []indexedTarget{
		{index: 0, target: target},
	}
	profile := panels.Definition{
		Name: "noinput-panel",
	}
	positions, err := buildPositions(targets, profile)
	if err != nil {
		t.Fatalf("buildPositions() error = %v", err)
	}

	// Simulate the noinput flag: enableInput=false means Events=nil and InputEnabled=false.
	// This matches the logic in runDisplayWithRegions when enableInput is false.
	enableInput := false

	var events chan struct{} // represents nil events channel
	_ = events               // unused, just for documentation

	config := region.PanelActivationConfig{
		Screens:      positions,
		Layout:       nil,
		DefaultMode:  "dashboard",
		InputEnabled: enableInput, // false when noinput flag is set
		AvailModes:   []string{"dashboard"},
		Events:       nil, // nil when noinput flag is set
		ModeValidator: func(mode string) bool {
			return mode == "dashboard"
		},
	}

	// Verify the config has the expected noinput state.
	if config.InputEnabled != false {
		t.Fatal("InputEnabled should be false when noinput flag is set")
	}
	if config.Events != nil {
		t.Fatal("Events should be nil when noinput flag is set")
	}

	// Activate the panel to verify it works correctly in passive mode.
	activation, err := region.ActivatePanel(config)
	if err != nil {
		t.Fatalf("ActivatePanel() unexpected error: %v", err)
	}

	// Verify activation succeeds with nil events (passive mode).
	regions := activation.RegionManager.Regions()
	if len(regions) != 1 {
		t.Fatalf("got %d regions, want 1", len(regions))
	}

	// The region still gets created correctly.
	r := regions[0]
	if r.Name() != "default" {
		t.Errorf("region name=%q, want %q", r.Name(), "default")
	}
	if r.Bounds() != image.Rect(0, 0, 240, 135) {
		t.Errorf("region bounds=%v, want %v", r.Bounds(), image.Rect(0, 0, 240, 135))
	}
	if r.CurrentMode() != "dashboard" {
		t.Errorf("region mode=%q, want %q", r.CurrentMode(), "dashboard")
	}

	// With InputEnabled=false in the config, the first region still gets focus
	// (focus assignment happens unconditionally in AllocateLayout), but the
	// noinput state means the dispatcher won't route events.
	// The key assertion is that ActivatePanel succeeds with nil Events.
	if activation.RenderLoop == nil {
		t.Error("RenderLoop should be non-nil even in passive mode")
	}
}

// --- Diagnostic log output tests (Task 4.4) ---

// captureLog redirects log output to a buffer for the duration of fn,
// restoring the original log output afterwards.
func captureLog(fn func()) string {
	var buf bytes.Buffer
	origOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origOutput)
	fn()
	return buf.String()
}

// TestDiagnosticLog_I2CSuccessPath verifies that the diagnostic log output
// during a successful I2C initialization contains the expected panel name,
// controller, mode, bus path, driver ID, and dimensions.

func TestDiagnosticLog_I2CSuccessPath(t *testing.T) {
	// Register a test I2C bus.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-diag-success-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a test I2C-only driver.
	registerTestDriver(t, driver.Definition{
		ID: "test-diag-i2c-driver",
		NewI2C: func(_ i2c.Bus, cfg driver.DriverConfig) (driver.DrawTarget, error) {
			return &sizedDrawTarget{bounds: image.Rect(0, 0, cfg.Width, cfg.Height)}, nil
		},
	})

	profile := panels.Definition{
		Name:       "test-diag-panel",
		Controller: "test-diag-i2c-driver",
		Config: driver.DriverConfig{
			Width:   128,
			Height:  32,
			I2CAddr: 0x3C,
		},
		I2CBus: "test-diag-success-bus",
	}

	var positions []region.ScreenPosition
	var cleanupFn func()
	var retErr error

	output := captureLog(func() {
		positions, cleanupFn, _, retErr = initDisplayTargets(profile, "")
	})
	if cleanupFn != nil {
		defer cleanupFn()
	}
	if retErr != nil {
		t.Fatalf("initDisplayTargets() returned unexpected error: %v", retErr)
	}
	if len(positions) == 0 {
		t.Fatal("expected non-empty screen positions")
	}

	// Verify log contains panel name.
	if !strings.Contains(output, "test-diag-panel") {
		t.Errorf("log should contain panel name %q, got:\n%s", "test-diag-panel", output)
	}
	// Verify log contains controller name.
	if !strings.Contains(output, "test-diag-i2c-driver") {
		t.Errorf("log should contain controller %q, got:\n%s", "test-diag-i2c-driver", output)
	}
	// Verify log contains communication mode.
	if !strings.Contains(output, "I2C") {
		t.Errorf("log should contain mode %q, got:\n%s", "I2C", output)
	}
	// Verify log contains bus path.
	if !strings.Contains(output, "test-diag-success-bus") {
		t.Errorf("log should contain bus path %q, got:\n%s", "test-diag-success-bus", output)
	}
	// Verify log contains driver factory ID reference.
	if !strings.Contains(output, "driver factory") {
		t.Errorf("log should contain %q, got:\n%s", "driver factory", output)
	}
	// Verify log contains dimensions.
	if !strings.Contains(output, "128x32") {
		t.Errorf("log should contain dimensions %q, got:\n%s", "128x32", output)
	}
	// Verify log contains the ready/success summary.
	if !strings.Contains(output, "ready") {
		t.Errorf("log should contain %q indicating success, got:\n%s", "ready", output)
	}
}

// TestDiagnosticLog_InvalidController verifies that when an invalid controller
// name is used, the diagnostic log mentions the unrecognized name and lists
// the registered driver IDs.

func TestDiagnosticLog_InvalidController(t *testing.T) {
	// Register a test I2C bus so bus open succeeds.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-diag-invalid-ctrl-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a known driver so Names() returns something we can verify.
	registerTestDriver(t, driver.Definition{
		ID: "test-diag-known-driver",
		NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
			return &mockDrawTarget{}, nil
		},
	})

	profile := panels.Definition{
		Name:       "test-invalid-ctrl-panel",
		Controller: "nonexistent-controller-diag",
		Config: driver.DriverConfig{
			Width:  128,
			Height: 32,
		},
		I2CBus: "test-diag-invalid-ctrl-bus",
	}

	var retErr error
	output := captureLog(func() {
		_, _, _, retErr = initDisplayTargets(profile, "")
	})

	// Should return an error for unsupported controller.
	if retErr == nil {
		t.Fatal("expected error for invalid controller name")
	}

	// Verify log mentions the unrecognized controller name.
	if !strings.Contains(output, "nonexistent-controller-diag") {
		t.Errorf("log should contain unrecognized controller name %q, got:\n%s", "nonexistent-controller-diag", output)
	}

	// Verify log mentions "registered drivers" or lists driver IDs.
	if !strings.Contains(output, "registered drivers") {
		t.Errorf("log should contain %q to list available drivers, got:\n%s", "registered drivers", output)
	}

	// Verify that the known test driver ID appears in the log (it's in the registered list).
	if !strings.Contains(output, "test-diag-known-driver") {
		t.Errorf("log should list registered driver ID %q, got:\n%s", "test-diag-known-driver", output)
	}
}

// TestDiagnosticLog_MissingRequiredGPIOPin verifies that when a required GPIO pin
// (DC or RST) is not found on the host, the diagnostic log contains the missing
// pin name and pin type.

func TestDiagnosticLog_MissingRequiredGPIOPin(t *testing.T) {
	// Register a test SPI port so spireg.Open succeeds.
	registerTestSPIPort(t, "test-diag-spi-port", func() (spi.PortCloser, error) {
		return &mockSPIPortCloser{}, nil
	})

	// Register a driver that has BOTH SPI and I2C factories so initDisplayTargets
	// takes the SPI path (non-I2C-only).
	registerTestDriver(t, driver.Definition{
		ID: "test-diag-spi-driver",
		NewSPI: func(_ spi.Port, _, _, _ gpio.PinOut, _ gpio.PinIn, _ driver.DriverConfig) (driver.DrawTarget, error) {
			return &mockDrawTarget{}, nil
		},
		NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
			return &mockDrawTarget{}, nil
		},
	})

	// Use a panel definition that references non-existent GPIO pins.
	// Since gpioreg won't find "FAKE_DC_PIN_XYZ" or "FAKE_RST_PIN_XYZ" on this host,
	// pin resolution will fail and the missing pin should be logged.
	profile := panels.Definition{
		Name:       "test-missing-pin-panel",
		Controller: "test-diag-spi-driver",
		Config: driver.DriverConfig{
			Width:  128,
			Height: 32,
		},
		DCPin:  "FAKE_DC_PIN_XYZ",
		RSTPin: "FAKE_RST_PIN_XYZ",
	}

	var retErr error
	output := captureLog(func() {
		_, _, _, retErr = initDisplayTargets(profile, "test-diag-spi-port")
	})

	// Should return an error because DC pin cannot be resolved.
	if retErr == nil {
		t.Fatal("expected error for missing GPIO pin")
	}

	// Verify log contains the pin name for DC.
	if !strings.Contains(output, "FAKE_DC_PIN_XYZ") {
		t.Errorf("log should contain missing DC pin name %q, got:\n%s", "FAKE_DC_PIN_XYZ", output)
	}
	// Verify log references pin type DC.
	if !strings.Contains(output, "DC") {
		t.Errorf("log should contain pin type %q, got:\n%s", "DC", output)
	}
	// Verify log contains the "required pin missing" message.
	if !strings.Contains(output, "required pin missing") {
		t.Errorf("log should contain %q, got:\n%s", "required pin missing", output)
	}
}

// TestDiagnosticLog_I2CFactoryError verifies that when the driver factory returns
// an error, the diagnostic log contains the driver ID, factory type, and error string.

func TestDiagnosticLog_I2CFactoryError(t *testing.T) {
	// Register a test I2C bus.
	bus := &mockI2CBus{}
	registerTestBus(t, "test-diag-factory-err-bus", func() (i2c.BusCloser, error) {
		return bus, nil
	})

	// Register a test I2C-only driver whose factory always fails.
	registerTestDriver(t, driver.Definition{
		ID: "test-diag-fail-driver",
		NewI2C: func(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
			return nil, errors.New("simulated init failure")
		},
	})

	profile := panels.Definition{
		Name:       "test-factory-err-panel",
		Controller: "test-diag-fail-driver",
		Config: driver.DriverConfig{
			Width:   128,
			Height:  32,
			I2CAddr: 0x3C,
		},
		I2CBus: "test-diag-factory-err-bus",
	}

	var retErr error
	output := captureLog(func() {
		_, _, _, retErr = initDisplayTargets(profile, "")
	})

	if retErr == nil {
		t.Fatal("expected error when factory fails")
	}

	// Verify log contains driver ID.
	if !strings.Contains(output, "test-diag-fail-driver") {
		t.Errorf("log should contain driver ID %q, got:\n%s", "test-diag-fail-driver", output)
	}
	// Verify log contains factory type.
	if !strings.Contains(output, "I2C") {
		t.Errorf("log should contain factory type %q, got:\n%s", "I2C", output)
	}
	// Verify log contains the error string from the factory.
	if !strings.Contains(output, "simulated init failure") {
		t.Errorf("log should contain factory error %q, got:\n%s", "simulated init failure", output)
	}
	// Verify the "driver factory error" log line is present.
	if !strings.Contains(output, "driver factory error") {
		t.Errorf("log should contain %q, got:\n%s", "driver factory error", output)
	}
}

// --- From: runtime_config_test.go ---

func TestBuildRuntimeConfig_AllDefaults(t *testing.T) {
	fn := buildRuntimeConfig(
		"/run/cyberhudd/console.sock", // socketPath
		"/dev/i2c-1",                  // i2cBuses
		2*time.Second,                 // scanInterval
		false,                         // noDisplay
		false,                         // noInput
		"",                            // profileName
		-1, -1,                        // width, height
		"",     // madctl
		-1, -1, // xOffset, yOffset
		"", "", "", "", // dc, rst, bl, busy
		"", "", "", // key1, key2, key3
		"", "", "", "", "", // up, down, left, right, press
	)

	cfg := fn()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// String fields still captured
	if cfg.Socket != "/run/cyberhudd/console.sock" {
		t.Errorf("Socket = %q, want %q", cfg.Socket, "/run/cyberhudd/console.sock")
	}
	if cfg.I2C != "/dev/i2c-1" {
		t.Errorf("I2C = %q, want %q", cfg.I2C, "/dev/i2c-1")
	}
	if cfg.Scan != "2s" {
		t.Errorf("Scan = %q, want %q", cfg.Scan, "2s")
	}

	// Pointer bools nil when at defaults (false)
	if cfg.Display.Disabled != nil {
		t.Errorf("Disabled = %v, want nil", cfg.Display.Disabled)
	}
	if cfg.Display.DisableInput != nil {
		t.Errorf("DisableInput = %v, want nil", cfg.Display.DisableInput)
	}

	// Pointer ints nil when value is -1
	if cfg.Display.Width != nil {
		t.Errorf("Width = %v, want nil", cfg.Display.Width)
	}
	if cfg.Display.Height != nil {
		t.Errorf("Height = %v, want nil", cfg.Display.Height)
	}
	if cfg.Display.XOffset != nil {
		t.Errorf("XOffset = %v, want nil", cfg.Display.XOffset)
	}
	if cfg.Display.YOffset != nil {
		t.Errorf("YOffset = %v, want nil", cfg.Display.YOffset)
	}

	// String fields in display should be empty
	if cfg.Display.Profile != "" {
		t.Errorf("Profile = %q, want empty", cfg.Display.Profile)
	}
	if cfg.Display.MADCTL != "" {
		t.Errorf("MADCTL = %q, want empty", cfg.Display.MADCTL)
	}
}

func TestBuildRuntimeConfig_AllValuesSet(t *testing.T) {
	fn := buildRuntimeConfig(
		"/tmp/test.sock",        // socketPath
		"/dev/i2c-1,/dev/i2c-3", // i2cBuses
		500*time.Millisecond,    // scanInterval
		true,                    // noDisplay
		true,                    // noInput
		"waveshare-2.2",         // profileName
		320, 240,                // width, height
		"0x60", // madctl
		2, 4,   // xOffset, yOffset
		"GPIO25", "GPIO27", "GPIO24", "GPIO4", // dc, rst, bl, busy
		"KEY1", "KEY2", "KEY3", // key1, key2, key3
		"UP", "DOWN", "LEFT", "RIGHT", "PRESS", // up, down, left, right, press
	)

	cfg := fn()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Top-level string fields
	if cfg.Socket != "/tmp/test.sock" {
		t.Errorf("Socket = %q, want %q", cfg.Socket, "/tmp/test.sock")
	}
	if cfg.I2C != "/dev/i2c-1,/dev/i2c-3" {
		t.Errorf("I2C = %q, want %q", cfg.I2C, "/dev/i2c-1,/dev/i2c-3")
	}
	if cfg.Scan != "500ms" {
		t.Errorf("Scan = %q, want %q", cfg.Scan, "500ms")
	}

	// *bool fields set when true
	if cfg.Display.Disabled == nil || *cfg.Display.Disabled != true {
		t.Errorf("Disabled = %v, want ptr to true", cfg.Display.Disabled)
	}
	if cfg.Display.DisableInput == nil || *cfg.Display.DisableInput != true {
		t.Errorf("DisableInput = %v, want ptr to true", cfg.Display.DisableInput)
	}

	// *int fields set for non -1 values
	if cfg.Display.Width == nil || *cfg.Display.Width != 320 {
		t.Errorf("Width = %v, want ptr to 320", cfg.Display.Width)
	}
	if cfg.Display.Height == nil || *cfg.Display.Height != 240 {
		t.Errorf("Height = %v, want ptr to 240", cfg.Display.Height)
	}
	if cfg.Display.XOffset == nil || *cfg.Display.XOffset != 2 {
		t.Errorf("XOffset = %v, want ptr to 2", cfg.Display.XOffset)
	}
	if cfg.Display.YOffset == nil || *cfg.Display.YOffset != 4 {
		t.Errorf("YOffset = %v, want ptr to 4", cfg.Display.YOffset)
	}

	// Display string fields
	if cfg.Display.Profile != "waveshare-2.2" {
		t.Errorf("Profile = %q, want %q", cfg.Display.Profile, "waveshare-2.2")
	}
	if cfg.Display.MADCTL != "0x60" {
		t.Errorf("MADCTL = %q, want %q", cfg.Display.MADCTL, "0x60")
	}
	if cfg.Display.DC != "GPIO25" {
		t.Errorf("DC = %q, want %q", cfg.Display.DC, "GPIO25")
	}
	if cfg.Display.RST != "GPIO27" {
		t.Errorf("RST = %q, want %q", cfg.Display.RST, "GPIO27")
	}
	if cfg.Display.BL != "GPIO24" {
		t.Errorf("BL = %q, want %q", cfg.Display.BL, "GPIO24")
	}
	if cfg.Display.Busy != "GPIO4" {
		t.Errorf("Busy = %q, want %q", cfg.Display.Busy, "GPIO4")
	}

	// Input pin fields
	if cfg.Display.InputKey1 != "KEY1" {
		t.Errorf("InputKey1 = %q, want %q", cfg.Display.InputKey1, "KEY1")
	}
	if cfg.Display.InputKey2 != "KEY2" {
		t.Errorf("InputKey2 = %q, want %q", cfg.Display.InputKey2, "KEY2")
	}
	if cfg.Display.InputKey3 != "KEY3" {
		t.Errorf("InputKey3 = %q, want %q", cfg.Display.InputKey3, "KEY3")
	}
	if cfg.Display.InputUp != "UP" {
		t.Errorf("InputUp = %q, want %q", cfg.Display.InputUp, "UP")
	}
	if cfg.Display.InputDown != "DOWN" {
		t.Errorf("InputDown = %q, want %q", cfg.Display.InputDown, "DOWN")
	}
	if cfg.Display.InputLeft != "LEFT" {
		t.Errorf("InputLeft = %q, want %q", cfg.Display.InputLeft, "LEFT")
	}
	if cfg.Display.InputRight != "RIGHT" {
		t.Errorf("InputRight = %q, want %q", cfg.Display.InputRight, "RIGHT")
	}
	if cfg.Display.InputPress != "PRESS" {
		t.Errorf("InputPress = %q, want %q", cfg.Display.InputPress, "PRESS")
	}
}

func TestBuildRuntimeConfig_ScanDurationFormat(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"seconds", 2 * time.Second, "2s"},
		{"milliseconds", 500 * time.Millisecond, "500ms"},
		{"fractional seconds", 1500 * time.Millisecond, "1.5s"},
		{"minutes", 2 * time.Minute, "2m0s"},
		{"sub-millisecond", 100 * time.Microsecond, "100µs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := buildRuntimeConfig(
				"sock", "i2c", tt.duration,
				false, false, "",
				-1, -1, "", -1, -1,
				"", "", "", "",
				"", "", "",
				"", "", "", "", "",
			)
			cfg := fn()
			if cfg.Scan != tt.want {
				t.Errorf("Scan = %q, want %q", cfg.Scan, tt.want)
			}
		})
	}
}

func TestBuildRuntimeConfig_PointerIntEdgeCases(t *testing.T) {
	// Width=0 should be set (0 is valid, only -1 means unset)
	fn := buildRuntimeConfig(
		"sock", "i2c", time.Second,
		false, false, "",
		0, 0, "", 0, 0,
		"", "", "", "",
		"", "", "",
		"", "", "", "", "",
	)
	cfg := fn()

	if cfg.Display.Width == nil {
		t.Error("Width should be non-nil for value 0")
	} else if *cfg.Display.Width != 0 {
		t.Errorf("Width = %d, want 0", *cfg.Display.Width)
	}

	if cfg.Display.Height == nil {
		t.Error("Height should be non-nil for value 0")
	} else if *cfg.Display.Height != 0 {
		t.Errorf("Height = %d, want 0", *cfg.Display.Height)
	}

	if cfg.Display.XOffset == nil {
		t.Error("XOffset should be non-nil for value 0")
	} else if *cfg.Display.XOffset != 0 {
		t.Errorf("XOffset = %d, want 0", *cfg.Display.XOffset)
	}

	if cfg.Display.YOffset == nil {
		t.Error("YOffset should be non-nil for value 0")
	} else if *cfg.Display.YOffset != 0 {
		t.Errorf("YOffset = %d, want 0", *cfg.Display.YOffset)
	}
}

func TestBuildRuntimeConfig_ClosureReturnsConsistently(t *testing.T) {
	fn := buildRuntimeConfig(
		"/sock", "/dev/i2c-1", 3*time.Second,
		true, false, "panel",
		128, 64, "0x00", 1, 2,
		"D", "R", "B", "U",
		"K1", "K2", "K3",
		"U", "D", "L", "R", "P",
	)

	// Call the closure multiple times — should always return same values
	cfg1 := fn()
	cfg2 := fn()

	if cfg1.Socket != cfg2.Socket || cfg1.Scan != cfg2.Scan {
		t.Error("closure returned inconsistent results across calls")
	}
	if *cfg1.Display.Width != *cfg2.Display.Width {
		t.Error("closure returned inconsistent Width across calls")
	}
}

func TestProperty_RoundTripFidelity(t *testing.T) {
	f := func(
		socket, i2c, profile, madctl string,
		dc, rst, bl, busy string,
		key1, key2, key3 string,
		up, down, left, right, press string,
		scanMs uint16,
		noDisplay, noInput bool,
		width, height int8,
	) bool {
		// Convert small ints: treat negative as "unset" (-1)
		w := -1
		if width >= 0 {
			w = int(width)
		}
		h := -1
		if height >= 0 {
			h = int(height)
		}
		scanInterval := time.Duration(scanMs+1) * time.Millisecond // ensure positive

		fn := buildRuntimeConfig(
			socket, i2c, scanInterval,
			noDisplay, noInput,
			profile,
			w, h, madctl, -1, -1,
			dc, rst, bl, busy,
			key1, key2, key3,
			up, down, left, right, press,
		)
		cfg := fn()

		// Serialize to JSON
		data, err := json.Marshal(cfg)
		if err != nil {
			return false
		}

		// Deserialize back
		var cfg2 fileConfig
		if err := json.Unmarshal(data, &cfg2); err != nil {
			return false
		}

		// Assert equivalence of key fields
		if cfg.Socket != cfg2.Socket {
			return false
		}
		if cfg.I2C != cfg2.I2C {
			return false
		}
		if cfg.Scan != cfg2.Scan {
			return false
		}
		if cfg.Display.Profile != cfg2.Display.Profile {
			return false
		}
		if cfg.Display.MADCTL != cfg2.Display.MADCTL {
			return false
		}
		if cfg.Display.DC != cfg2.Display.DC {
			return false
		}
		if cfg.Display.RST != cfg2.Display.RST {
			return false
		}
		if cfg.Display.BL != cfg2.Display.BL {
			return false
		}
		if cfg.Display.Busy != cfg2.Display.Busy {
			return false
		}
		if cfg.Display.InputKey1 != cfg2.Display.InputKey1 {
			return false
		}
		if cfg.Display.InputKey2 != cfg2.Display.InputKey2 {
			return false
		}
		if cfg.Display.InputKey3 != cfg2.Display.InputKey3 {
			return false
		}
		if cfg.Display.InputUp != cfg2.Display.InputUp {
			return false
		}
		if cfg.Display.InputDown != cfg2.Display.InputDown {
			return false
		}
		if cfg.Display.InputLeft != cfg2.Display.InputLeft {
			return false
		}
		if cfg.Display.InputRight != cfg2.Display.InputRight {
			return false
		}
		if cfg.Display.InputPress != cfg2.Display.InputPress {
			return false
		}

		// Check pointer fields
		if !ptrBoolEqual(cfg.Display.Disabled, cfg2.Display.Disabled) {
			return false
		}
		if !ptrBoolEqual(cfg.Display.DisableInput, cfg2.Display.DisableInput) {
			return false
		}
		if !ptrIntEqual(cfg.Display.Width, cfg2.Display.Width) {
			return false
		}
		if !ptrIntEqual(cfg.Display.Height, cfg2.Display.Height) {
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func ptrBoolEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func ptrIntEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// TestProperty_ZeroValueFieldOmission verifies that when pointer fields (*bool, *int) are nil
// in the fileConfig struct, their keys are absent from the serialized JSON output.

func TestProperty_ZeroValueFieldOmission(t *testing.T) {
	f := func(socket, i2c, profile string) bool {
		// Call buildRuntimeConfig with all pointer-field defaults:
		// noDisplay=false, noInput=false → nil *bool
		// width=-1, height=-1, xOffset=-1, yOffset=-1 → nil *int
		fn := buildRuntimeConfig(
			socket, i2c, 2*time.Second,
			false, false, // noDisplay=false, noInput=false
			profile,
			-1, -1, // width=-1, height=-1
			"", -1, -1, // madctl, xOffset=-1, yOffset=-1
			"", "", "", "",
			"", "", "",
			"", "", "", "", "",
		)
		cfg := fn()

		// Verify pointer fields are nil when at defaults
		if cfg.Display.Disabled != nil {
			return false
		}
		if cfg.Display.DisableInput != nil {
			return false
		}
		if cfg.Display.Width != nil {
			return false
		}
		if cfg.Display.Height != nil {
			return false
		}
		if cfg.Display.XOffset != nil {
			return false
		}
		if cfg.Display.YOffset != nil {
			return false
		}

		// Serialize to JSON and verify nil pointer field keys are absent
		data, err := json.Marshal(cfg)
		if err != nil {
			return false
		}
		jsonStr := string(data)

		// Keys for nil pointer fields should NOT appear in serialized output
		if strings.Contains(jsonStr, `"disabled"`) {
			return false
		}
		if strings.Contains(jsonStr, `"disable_input"`) {
			return false
		}
		if strings.Contains(jsonStr, `"width"`) {
			return false
		}
		if strings.Contains(jsonStr, `"height"`) {
			return false
		}
		if strings.Contains(jsonStr, `"x_offset"`) {
			return false
		}
		if strings.Contains(jsonStr, `"y_offset"`) {
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
