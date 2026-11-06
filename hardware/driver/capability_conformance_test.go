package driver_test

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/driver"

	// Register every real hardware driver so the registry is complete. This is
	// lever 1 of the capability-utilization work: a single gate that proves every
	// driver honestly declares what its DrawImage actually drives. A new driver
	// that forgets to declare ColorDepth/Capability — the exact mistake that left
	// is31fl3731 silently understated as mono — fails here at build time instead
	// of shipping a panel that renders below its hardware.
	_ "github.com/databeast/cyberhud/hardware/driver/epd4in26"
	_ "github.com/databeast/cyberhud/hardware/driver/ht16k33"
	_ "github.com/databeast/cyberhud/hardware/driver/ili9486"
	_ "github.com/databeast/cyberhud/hardware/driver/is31fl3731"
	_ "github.com/databeast/cyberhud/hardware/driver/max7219"
	_ "github.com/databeast/cyberhud/hardware/driver/sh1106"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1305"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1306"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1351"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1680"
	_ "github.com/databeast/cyberhud/hardware/driver/st7735s"
	_ "github.com/databeast/cyberhud/hardware/driver/st7789"
	_ "github.com/databeast/cyberhud/hardware/driver/uc8179"
)

// TestEveryDriverDeclaresColorDepth asserts that no driver leaves ColorDepth at
// its zero value. The zero value is invalid precisely so that omission is a
// failure rather than a silent default to mono — the trap that hid the
// is31fl3731 under-declaration.
func TestEveryDriverDeclaresColorDepth(t *testing.T) {
	for _, id := range driver.Names() {
		def, ok := driver.Get(id)
		if !ok {
			t.Fatalf("driver %q registered but Get failed", id)
		}
		if def.ColorDepth == driver.DepthUnset {
			t.Errorf("driver %q leaves ColorDepth unset; it must declare mono1/gray/color so its advertised capability can be checked against hardware reality", id)
		}
	}
}

// TestDriverCapabilityMatchesColorDepth asserts that the capability each driver
// advertises to style selection is consistent with the per-pixel depth its
// DrawImage actually drives. A driver that renders color but advertises a mono
// capability (ssd1351, uc8179 before this work) under-drives the panel: style
// selection treats it as mono and the hue/luminance the hardware can show is
// stripped upstream. A driver that advertises more than it renders would send
// color a panel cannot show. Both directions fail here.
func TestDriverCapabilityMatchesColorDepth(t *testing.T) {
	for _, id := range driver.Names() {
		def, ok := driver.Get(id)
		if !ok {
			t.Fatalf("driver %q registered but Get failed", id)
		}
		if def.ColorDepth == driver.DepthUnset {
			continue // reported by TestEveryDriverDeclaresColorDepth
		}
		advertised := driver.DepthForCapability(def.Capability)
		if advertised != def.ColorDepth {
			t.Errorf("driver %q advertises capability %d (depth %s) but its hardware ColorDepth is %s: the panel is driven below (or above) what it can show",
				id, def.Capability, advertised, def.ColorDepth)
		}
	}
}

// TestMonochromeFlagAgreesWithColorDepth keeps the legacy Monochrome bool — which
// panels.normalize() still force-copies onto panel definitions — in sync with the
// richer ColorDepth, so the two cannot drift.
func TestMonochromeFlagAgreesWithColorDepth(t *testing.T) {
	for _, id := range driver.Names() {
		def, ok := driver.Get(id)
		if !ok {
			t.Fatalf("driver %q registered but Get failed", id)
		}
		if def.ColorDepth == driver.DepthUnset {
			continue
		}
		wantMono := def.ColorDepth == driver.DepthMono1
		if def.Monochrome != wantMono {
			t.Errorf("driver %q has Monochrome=%v but ColorDepth=%s (expected Monochrome=%v)",
				id, def.Monochrome, def.ColorDepth, wantMono)
		}
	}
}
