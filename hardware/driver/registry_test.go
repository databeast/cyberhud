package driver_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/hardware/driver"
	_ "github.com/databeast/cyberhud/hardware/driver/sh1106"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1680"
	_ "github.com/databeast/cyberhud/hardware/driver/st7735s"
	_ "github.com/databeast/cyberhud/hardware/driver/st7789"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/spi"
	"pgregory.net/rapid"
)

func TestRegistryIncludesBuiltIns(t *testing.T) {
	cases := []struct {
		id         string
		monochrome bool
	}{
		{id: "st7789", monochrome: false},
		{id: "st7735s", monochrome: false},
		{id: "sh1106", monochrome: true},
		{id: "ssd1680", monochrome: true},
	}
	for _, tc := range cases {
		def, ok := driver.Get(tc.id)
		if !ok {
			t.Fatalf("Get(%q) missing", tc.id)
		}
		if def.ID != tc.id {
			t.Fatalf("Get(%q).ID=%q", tc.id, def.ID)
		}
		if def.Monochrome != tc.monochrome {
			t.Fatalf("Get(%q).Monochrome=%v, want %v", tc.id, def.Monochrome, tc.monochrome)
		}
		if len(def.OptionDefs) == 0 {
			t.Fatalf("Get(%q).OptionDefs expected published options", tc.id)
		}
	}
}

// dummySPIFactory is a stub SPI factory for testing registration acceptance.
func dummySPIFactory(_ spi.Port, _, _, _ gpio.PinOut, _ gpio.PinIn, _ driver.DriverConfig) (driver.DrawTarget, error) {
	return nil, nil
}

// dummyI2CFactory is a stub I2C factory for testing registration acceptance.
func dummyI2CFactory(_ i2c.Bus, _ driver.DriverConfig) (driver.DrawTarget, error) {
	return nil, nil
}

// genValidID generates a non-empty lowercase alphanumeric ID with a unique prefix
// to avoid collisions with other tests and built-in drivers.
func genValidID(t *rapid.T, prefix string) string {
	// Generate a suffix of 4-16 lowercase letters
	suffix := rapid.StringMatching(`[a-z][a-z0-9]{3,15}`).Draw(t, "id_suffix")
	return fmt.Sprintf("_prop_%s_%s", prefix, suffix)
}

func TestProperty_RegistrationAcceptanceWithFactory(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		id := genValidID(t, "accept")

		// Pick a factory combination: I2C-only, SPI-only, or both
		combo := rapid.SampledFrom([]string{"i2c-only", "spi-only", "both"}).Draw(t, "factory_combo")

		def := driver.Definition{
			ID:         id,
			Monochrome: rapid.Bool().Draw(t, "monochrome"),
		}

		switch combo {
		case "i2c-only":
			def.NewI2C = dummyI2CFactory
		case "spi-only":
			def.NewSPI = dummySPIFactory
		case "both":
			def.NewSPI = dummySPIFactory
			def.NewI2C = dummyI2CFactory
		}

		driver.Register(def)

		got, ok := driver.Get(id)
		if !ok {
			t.Fatalf("Get(%q) returned false after Register with combo=%s", id, combo)
		}
		// Verify the stored definition matches what was registered
		if got.ID != strings.ToLower(strings.TrimSpace(id)) {
			t.Fatalf("Get(%q).ID = %q, want %q", id, got.ID, id)
		}
		if got.Monochrome != def.Monochrome {
			t.Fatalf("Get(%q).Monochrome = %v, want %v", id, got.Monochrome, def.Monochrome)
		}

		// Verify factory presence matches the combo
		switch combo {
		case "i2c-only":
			if got.NewI2C == nil {
				t.Fatalf("Get(%q).NewI2C is nil for i2c-only combo", id)
			}
			if got.NewSPI != nil {
				t.Fatalf("Get(%q).NewSPI is non-nil for i2c-only combo", id)
			}
		case "spi-only":
			if got.NewSPI == nil {
				t.Fatalf("Get(%q).NewSPI is nil for spi-only combo", id)
			}
			if got.NewI2C != nil {
				t.Fatalf("Get(%q).NewI2C is non-nil for spi-only combo", id)
			}
		case "both":
			if got.NewSPI == nil {
				t.Fatalf("Get(%q).NewSPI is nil for both combo", id)
			}
			if got.NewI2C == nil {
				t.Fatalf("Get(%q).NewI2C is nil for both combo", id)
			}
		}
	})
}

func TestProperty_RegistrationRejectionNoFactories(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		id := genValidID(t, "reject")

		def := driver.Definition{
			ID:         id,
			Monochrome: rapid.Bool().Draw(t, "monochrome"),
			// Both factories are nil — registration should be discarded
			NewSPI: nil,
			NewI2C: nil,
		}

		driver.Register(def)

		_, ok := driver.Get(id)
		if ok {
			t.Fatalf("Get(%q) returned true after Register with both factories nil", id)
		}
	})
}
