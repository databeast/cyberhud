package waveshare_triple_screen_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"

	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"

	_ "github.com/databeast/cyberhud/hardware/driver/st7735s"
	_ "github.com/databeast/cyberhud/hardware/driver/st7789"
	_ "github.com/databeast/cyberhud/hardware/panels/waveshare_triple_screen"
)

func TestPanelRegistration(t *testing.T) {
	t.Run("Get returns panel without error", func(t *testing.T) {
		_, err := panels.Get("waveshare-triple-screen")
		if err != nil {
			t.Fatalf("panels.Get(\"waveshare-triple-screen\") returned error: %v", err)
		}
	})

	t.Run("name appears in Names()", func(t *testing.T) {
		names := panels.Names()
		found := false
		for _, n := range names {
			if strings.EqualFold(n, "waveshare-triple-screen") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("panels.Names() = %v, want it to contain %q", names, "waveshare-triple-screen")
		}
	})

	t.Run("Controller is virtual", func(t *testing.T) {
		def, err := panels.Get("waveshare-triple-screen")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Controller != "virtual" {
			t.Errorf("Controller = %q, want %q", def.Controller, "virtual")
		}
	})

	t.Run("Virtual slice has length 3", func(t *testing.T) {
		def, err := panels.Get("waveshare-triple-screen")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if len(def.Virtual) != 3 {
			t.Errorf("len(Virtual) = %d, want 3", len(def.Virtual))
		}
	})
}

func TestCenterMainScreen(t *testing.T) {
	def, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if len(def.Virtual) < 1 {
		t.Fatalf("Virtual slice has length %d, need at least 1", len(def.Virtual))
	}
	s := def.Virtual[0]

	t.Run("Name", func(t *testing.T) {
		if s.Name != "main" {
			t.Errorf("Name = %q, want %q", s.Name, "main")
		}
	})

	t.Run("Controller", func(t *testing.T) {
		if s.Controller != "st7789" {
			t.Errorf("Controller = %q, want %q", s.Controller, "st7789")
		}
	})

	t.Run("SPI", func(t *testing.T) {
		if s.SPI != "SPI1.0" {
			t.Errorf("SPI = %q, want %q", s.SPI, "SPI1.0")
		}
	})

	t.Run("Width", func(t *testing.T) {
		if s.Config.Width != 240 {
			t.Errorf("Width = %d, want %d", s.Config.Width, 240)
		}
	})

	t.Run("Height", func(t *testing.T) {
		if s.Config.Height != 240 {
			t.Errorf("Height = %d, want %d", s.Config.Height, 240)
		}
	})

	t.Run("MADCTL", func(t *testing.T) {
		want := driver.MadctlMX
		if s.Config.MADCTL != want {
			t.Errorf("MADCTL = 0x%02X, want 0x%02X", s.Config.MADCTL, want)
		}
	})

	t.Run("DCPin", func(t *testing.T) {
		if s.DCPin != "GPIO22" {
			t.Errorf("DCPin = %q, want %q", s.DCPin, "GPIO22")
		}
	})

	t.Run("RSTPin", func(t *testing.T) {
		if s.RSTPin != "GPIO27" {
			t.Errorf("RSTPin = %q, want %q", s.RSTPin, "GPIO27")
		}
	})

	t.Run("BLPin", func(t *testing.T) {
		if s.BLPin != "GPIO19" {
			t.Errorf("BLPin = %q, want %q", s.BLPin, "GPIO19")
		}
	})

	t.Run("DefaultMode", func(t *testing.T) {
		if s.DefaultMode != "menu" {
			t.Errorf("DefaultMode = %q, want %q", s.DefaultMode, "menu")
		}
	})
}

func TestLeftAuxScreen(t *testing.T) {
	def, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if len(def.Virtual) < 2 {
		t.Fatalf("Virtual slice has length %d, need at least 2", len(def.Virtual))
	}
	s := def.Virtual[1]

	t.Run("Name", func(t *testing.T) {
		if s.Name != "left-aux" {
			t.Errorf("Name = %q, want %q", s.Name, "left-aux")
		}
	})

	t.Run("Controller", func(t *testing.T) {
		if s.Controller != "st7735s" {
			t.Errorf("Controller = %q, want %q", s.Controller, "st7735s")
		}
	})

	t.Run("SPI", func(t *testing.T) {
		if s.SPI != "SPI0.0" {
			t.Errorf("SPI = %q, want %q", s.SPI, "SPI0.0")
		}
	})

	t.Run("Width", func(t *testing.T) {
		if s.Config.Width != 160 {
			t.Errorf("Width = %d, want %d", s.Config.Width, 160)
		}
	})

	t.Run("Height", func(t *testing.T) {
		if s.Config.Height != 80 {
			t.Errorf("Height = %d, want %d", s.Config.Height, 80)
		}
	})

	t.Run("MADCTL", func(t *testing.T) {
		want := driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR
		if s.Config.MADCTL != want {
			t.Errorf("MADCTL = 0x%02X, want 0x%02X", s.Config.MADCTL, want)
		}
	})

	t.Run("XOffset", func(t *testing.T) {
		if s.Config.XOffset != 1 {
			t.Errorf("XOffset = %d, want %d", s.Config.XOffset, 1)
		}
	})

	t.Run("YOffset", func(t *testing.T) {
		if s.Config.YOffset != 26 {
			t.Errorf("YOffset = %d, want %d", s.Config.YOffset, 26)
		}
	})

	t.Run("DCPin", func(t *testing.T) {
		if s.DCPin != panels.GPIO4 {
			t.Errorf("DCPin = %q, want %q", s.DCPin, panels.GPIO4)
		}
	})

	t.Run("RSTPin", func(t *testing.T) {
		if s.RSTPin != panels.GPIO24 {
			t.Errorf("RSTPin = %q, want %q", s.RSTPin, panels.GPIO24)
		}
	})

	t.Run("BLPin", func(t *testing.T) {
		if s.BLPin != panels.GPIO13 {
			t.Errorf("BLPin = %q, want %q", s.BLPin, panels.GPIO13)
		}
	})

	t.Run("DefaultMode", func(t *testing.T) {
		if s.DefaultMode != "stemma" {
			t.Errorf("DefaultMode = %q, want %q", s.DefaultMode, "stemma")
		}
	})

	t.Run("Index", func(t *testing.T) {
		if s.Index != 1 {
			t.Errorf("Index = %d, want %d", s.Index, 1)
		}
	})
}

func TestRightAuxScreen(t *testing.T) {
	def, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if len(def.Virtual) < 3 {
		t.Fatalf("Virtual slice has length %d, need at least 3", len(def.Virtual))
	}
	s := def.Virtual[2]

	t.Run("Name", func(t *testing.T) {
		if s.Name != "right-aux" {
			t.Errorf("Name = %q, want %q", s.Name, "right-aux")
		}
	})

	t.Run("Controller", func(t *testing.T) {
		if s.Controller != "st7735s" {
			t.Errorf("Controller = %q, want %q", s.Controller, "st7735s")
		}
	})

	t.Run("SPI", func(t *testing.T) {
		if s.SPI != "SPI0.1" {
			t.Errorf("SPI = %q, want %q", s.SPI, "SPI0.1")
		}
	})

	t.Run("Width", func(t *testing.T) {
		if s.Config.Width != 160 {
			t.Errorf("Width = %d, want %d", s.Config.Width, 160)
		}
	})

	t.Run("Height", func(t *testing.T) {
		if s.Config.Height != 80 {
			t.Errorf("Height = %d, want %d", s.Config.Height, 80)
		}
	})

	t.Run("MADCTL", func(t *testing.T) {
		want := driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR
		if s.Config.MADCTL != want {
			t.Errorf("MADCTL = 0x%02X, want 0x%02X", s.Config.MADCTL, want)
		}
	})

	t.Run("XOffset", func(t *testing.T) {
		if s.Config.XOffset != 1 {
			t.Errorf("XOffset = %d, want %d", s.Config.XOffset, 1)
		}
	})

	t.Run("YOffset", func(t *testing.T) {
		if s.Config.YOffset != 26 {
			t.Errorf("YOffset = %d, want %d", s.Config.YOffset, 26)
		}
	})

	t.Run("DCPin", func(t *testing.T) {
		if s.DCPin != panels.GPIO5 {
			t.Errorf("DCPin = %q, want %q", s.DCPin, panels.GPIO5)
		}
	})

	t.Run("RSTPin", func(t *testing.T) {
		if s.RSTPin != panels.GPIO23 {
			t.Errorf("RSTPin = %q, want %q", s.RSTPin, panels.GPIO23)
		}
	})

	t.Run("BLPin", func(t *testing.T) {
		if s.BLPin != panels.GPIO12 {
			t.Errorf("BLPin = %q, want %q", s.BLPin, panels.GPIO12)
		}
	})

	t.Run("DefaultMode", func(t *testing.T) {
		if s.DefaultMode != "gpio" {
			t.Errorf("DefaultMode = %q, want %q", s.DefaultMode, "gpio")
		}
	})

	t.Run("Index", func(t *testing.T) {
		if s.Index != 2 {
			t.Errorf("Index = %d, want %d", s.Index, 2)
		}
	})
}

func TestInputPins(t *testing.T) {
	def, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	t.Run("Key1 is GPIO25", func(t *testing.T) {
		if def.Inputs.Key1 != panels.GPIO25 {
			t.Errorf("Key1 = %q, want %q", def.Inputs.Key1, panels.GPIO25)
		}
	})

	t.Run("Key2 is GPIO26", func(t *testing.T) {
		if def.Inputs.Key2 != panels.GPIO26 {
			t.Errorf("Key2 = %q, want %q", def.Inputs.Key2, panels.GPIO26)
		}
	})

	t.Run("Key3 is empty", func(t *testing.T) {
		if def.Inputs.Key3 != "" {
			t.Errorf("Key3 = %q, want empty", def.Inputs.Key3)
		}
	})

	t.Run("JoyUp is empty", func(t *testing.T) {
		if def.Inputs.JoyUp != "" {
			t.Errorf("JoyUp = %q, want empty", def.Inputs.JoyUp)
		}
	})

	t.Run("JoyDown is empty", func(t *testing.T) {
		if def.Inputs.JoyDown != "" {
			t.Errorf("JoyDown = %q, want empty", def.Inputs.JoyDown)
		}
	})

	t.Run("JoyLeft is empty", func(t *testing.T) {
		if def.Inputs.JoyLeft != "" {
			t.Errorf("JoyLeft = %q, want empty", def.Inputs.JoyLeft)
		}
	})

	t.Run("JoyRight is empty", func(t *testing.T) {
		if def.Inputs.JoyRight != "" {
			t.Errorf("JoyRight = %q, want empty", def.Inputs.JoyRight)
		}
	})

	t.Run("JoyPressed is empty", func(t *testing.T) {
		if def.Inputs.JoyPressed != "" {
			t.Errorf("JoyPressed = %q, want empty", def.Inputs.JoyPressed)
		}
	})

	t.Run("Any returns true", func(t *testing.T) {
		if !def.Inputs.Any() {
			t.Errorf("Inputs.Any() = false, want true")
		}
	})
}

func TestPinConflicts(t *testing.T) {
	def, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	t.Run("PinNotices returns exactly 2 entries", func(t *testing.T) {
		notices := panels.PinNotices(def)
		if len(notices) != 2 {
			t.Fatalf("len(PinNotices) = %d, want 2; notices: %v", len(notices), notices)
		}
	})

	t.Run("one notice contains GPIO13 and display1_bl", func(t *testing.T) {
		notices := panels.PinNotices(def)
		found := false
		for _, n := range notices {
			if strings.Contains(n, "GPIO13") && strings.Contains(n, "display1_bl") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no notice contains both \"GPIO13\" and \"display1_bl\"; notices: %v", notices)
		}
	})

	t.Run("one notice contains GPIO18 and display0_cs", func(t *testing.T) {
		notices := panels.PinNotices(def)
		found := false
		for _, n := range notices {
			if strings.Contains(n, "GPIO18") && strings.Contains(n, "display0_cs") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no notice contains both \"GPIO18\" and \"display0_cs\"; notices: %v", notices)
		}
	})

	t.Run("BuildPinReport has conflict on 3-pin connector GPIO13", func(t *testing.T) {
		report := panels.BuildPinReport(def)
		found := false
		for _, line := range strings.Split(report, "\n") {
			if strings.Contains(line, "3-pin connector GPIO13") && strings.Contains(line, "status=conflict") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BuildPinReport missing line with \"3-pin connector GPIO13\" and \"status=conflict\";\nreport:\n%s", report)
		}
	})

	t.Run("BuildPinReport has conflict on 3-pin connector GPIO18", func(t *testing.T) {
		report := panels.BuildPinReport(def)
		found := false
		for _, line := range strings.Split(report, "\n") {
			if strings.Contains(line, "3-pin connector GPIO18") && strings.Contains(line, "status=conflict") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BuildPinReport missing line with \"3-pin connector GPIO18\" and \"status=conflict\";\nreport:\n%s", report)
		}
	})
}

func TestProperty_InputPinsAny(t *testing.T) {
	property := func(pins panels.InputPins) bool {
		expected := pins.Key1 != "" || pins.Key2 != "" || pins.Key3 != "" ||
			pins.JoyUp != "" || pins.JoyDown != "" || pins.JoyLeft != "" ||
			pins.JoyRight != "" || pins.JoyPressed != ""
		return pins.Any() == expected
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("property failed: %v", err)
	}
}

func TestProperty_UnknownPanelReturnsError(t *testing.T) {
	registeredNames := panels.Names()
	known := make(map[string]bool, len(registeredNames))
	for _, n := range registeredNames {
		known[n] = true
	}

	prop := func(name string) bool {
		// If the random string happens to match a registered panel, skip it
		if known[strings.ToLower(strings.TrimSpace(name))] {
			return true
		}
		_, err := panels.Get(name)
		return err != nil
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(prop, cfg); err != nil {
		t.Errorf("property failed: %v", err)
	}
}

func TestProperty_VirtualPanelInvariants(t *testing.T) {
	controllers := []string{"st7789", "st7735s", "ssd1306", "sh1106", "ssd1351"}

	// generateVirtualPanel builds a well-formed virtual panel definition
	// with the given screen count using randomized but valid field values.
	generateVirtualPanel := func(rng *rand.Rand, screenCount int) panels.Definition {
		screens := make([]panels.Screen, screenCount)
		for i := 0; i < screenCount; i++ {
			screens[i] = panels.Screen{
				Index:      i,
				Name:       fmt.Sprintf("screen-%d-%d", i, rng.Intn(10000)),
				Controller: controllers[rng.Intn(len(controllers))],
				Config: driver.DriverConfig{
					Width:  rng.Intn(1920) + 1, // 1..1920
					Height: rng.Intn(1080) + 1, // 1..1080
				},
			}
		}
		return panels.Definition{
			Name:       fmt.Sprintf("test-panel-%d", rng.Intn(100000)),
			Controller: "virtual",
			Virtual:    screens,
		}
	}

	prop := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))
		screenCount := rng.Intn(10) + 1 // 1..10
		def := generateVirtualPanel(rng, screenCount)

		// Invariant 1: All screen names are unique
		names := make(map[string]bool, len(def.Virtual))
		for _, s := range def.Virtual {
			if names[s.Name] {
				t.Logf("duplicate screen name: %q", s.Name)
				return false
			}
			names[s.Name] = true
		}

		// Invariant 2: Indices form a contiguous sequence starting from 0
		for i, s := range def.Virtual {
			if s.Index != i {
				t.Logf("non-contiguous index: screen %d has Index=%d", i, s.Index)
				return false
			}
		}

		// Invariant 3: Each screen has non-empty Controller
		for i, s := range def.Virtual {
			if s.Controller == "" {
				t.Logf("screen %d has empty Controller", i)
				return false
			}
		}

		// Invariant 4: Each screen Config has Width > 0 and Height > 0
		for i, s := range def.Virtual {
			if s.Config.Width <= 0 {
				t.Logf("screen %d has Width=%d (want > 0)", i, s.Config.Width)
				return false
			}
			if s.Config.Height <= 0 {
				t.Logf("screen %d has Height=%d (want > 0)", i, s.Config.Height)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(prop, cfg); err != nil {
		t.Errorf("property failed: %v", err)
	}
}

func TestProperty_PinConflictCount(t *testing.T) {
	// Pool of GPIO pin names to randomly assign to panel fields.
	pinPool := []string{
		"", // empty means unassigned
		panels.GPIO4, panels.GPIO5, panels.GPIO6,
		panels.GPIO12, panels.GPIO13, panels.GPIO17,
		panels.GPIO18, panels.GPIO19, panels.GPIO22,
		panels.GPIO23, panels.GPIO24, panels.GPIO25,
		panels.GPIO27,
	}

	// SPI options - SPI1.0 automatically assigns GPIO18 as CS.
	spiOptions := []string{"", "SPI0.0", "SPI0.1", "SPI1.0", "SPI1.1"}

	// pinsContain checks if a given GPIO name appears in any of the provided pin strings.
	pinsContain := func(pins []string, gpio string) bool {
		for _, p := range pins {
			if p == gpio {
				return true
			}
		}
		return false
	}

	// spiPins returns the CS pin for a given SPI device (the one that could conflict).
	spiCSPin := func(spi string) string {
		switch strings.ToUpper(strings.TrimSpace(spi)) {
		case "SPI0.0":
			return panels.GPIO8
		case "SPI0.1":
			return panels.GPIO7
		case "SPI1.0":
			return panels.GPIO18
		case "SPI1.1":
			return panels.GPIO17
		default:
			return ""
		}
	}

	// spiAllPins returns all GPIO pins assigned by a SPI device.
	spiAllPins := func(spi string) []string {
		switch strings.ToUpper(strings.TrimSpace(spi)) {
		case "SPI0.0":
			return []string{panels.GPIO10, panels.GPIO11, panels.GPIO8}
		case "SPI0.1":
			return []string{panels.GPIO10, panels.GPIO11, panels.GPIO7}
		case "SPI1.0":
			return []string{panels.GPIO20, panels.GPIO21, panels.GPIO18}
		case "SPI1.1":
			return []string{panels.GPIO20, panels.GPIO21, panels.GPIO17}
		default:
			return nil
		}
	}

	prop := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))

		// Randomly decide how many virtual screens (0-3)
		numScreens := rng.Intn(4)

		// Pick random pins for the top-level panel definition
		topDC := pinPool[rng.Intn(len(pinPool))]
		topRST := pinPool[rng.Intn(len(pinPool))]
		topBL := pinPool[rng.Intn(len(pinPool))]

		def := panels.Definition{
			Name:       fmt.Sprintf("test-panel-%d", seed),
			Controller: "virtual",
			DCPin:      topDC,
			RSTPin:     topRST,
			BLPin:      topBL,
		}

		// Build virtual screens with random pin assignments
		for i := 0; i < numScreens; i++ {
			screen := panels.Screen{
				Index:      i,
				Name:       fmt.Sprintf("screen-%d", i),
				Controller: "st7789",
				SPI:        spiOptions[rng.Intn(len(spiOptions))],
				DCPin:      pinPool[rng.Intn(len(pinPool))],
				RSTPin:     pinPool[rng.Intn(len(pinPool))],
				BLPin:      pinPool[rng.Intn(len(pinPool))],
			}
			def.Virtual = append(def.Virtual, screen)
		}

		// Compute the expected count: how many of {GPIO13, GPIO18} appear
		// in the panel's effective pin set.
		var allPins []string
		if topDC != "" {
			allPins = append(allPins, topDC)
		}
		if topRST != "" {
			allPins = append(allPins, topRST)
		}
		if topBL != "" {
			allPins = append(allPins, topBL)
		}
		for _, s := range def.Virtual {
			allPins = append(allPins, spiAllPins(s.SPI)...)
			if s.DCPin != "" {
				allPins = append(allPins, s.DCPin)
			}
			if s.RSTPin != "" {
				allPins = append(allPins, s.RSTPin)
			}
			if s.BLPin != "" {
				allPins = append(allPins, s.BLPin)
			}
		}

		// Also include input pins (they are empty for our generated panel)
		// but the real function checks them, so we leave Inputs at zero value.

		_ = spiCSPin // used indirectly via spiAllPins

		expectedCount := 0
		if pinsContain(allPins, panels.GPIO13) {
			expectedCount++
		}
		if pinsContain(allPins, panels.GPIO18) {
			expectedCount++
		}

		notices := panels.PinNotices(def)
		if len(notices) != expectedCount {
			t.Logf("seed=%d numScreens=%d topDC=%q topRST=%q topBL=%q",
				seed, numScreens, topDC, topRST, topBL)
			for i, s := range def.Virtual {
				t.Logf("  screen[%d]: SPI=%q DC=%q RST=%q BL=%q",
					i, s.SPI, s.DCPin, s.RSTPin, s.BLPin)
			}
			t.Logf("  allPins=%v expectedCount=%d got=%d notices=%v",
				allPins, expectedCount, len(notices), notices)
			return false
		}
		return true
	}

	cfg := &quick.Config{MaxCount: 100}
	if err := quick.Check(prop, cfg); err != nil {
		t.Errorf("property failed: %v", err)
	}
}
