package cyberhudd

import (
	"slices"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/thermal"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func TestWaveshare144ThermalStyleAlias(t *testing.T) {
	hints := textlayout.TextHints{PanelProduct: "waveshare-1.44"}
	selected, reason := style.ResolveStyle(thermal.ThermalRegistry, hints, "thermal", "")
	if selected.Name() != "color-fast-128x128-minimal" || reason != "alias" {
		t.Fatalf("resolved style = %q (%s), want color-fast-128x128-minimal (alias)", selected.Name(), reason)
	}
}

func TestPanelStyleAliases(t *testing.T) {
	for _, panel := range panels.All() {
		products := []string{panel.Name}
		for _, screen := range panel.Virtual {
			products = append(products, panel.Name+"/"+screen.Name)
		}
		for _, product := range products {
			for _, alias := range style.AliasesForProduct(product) {
				def, ok := catalog.Describe(alias.Mode)
				if !ok {
					t.Errorf("%s: alias references unknown mode %q", product, alias.Mode)
					continue
				}
				var allowed []string
				for _, option := range def.Options {
					if option.Key == "style" && len(option.Allowed) > 0 {
						allowed = option.Allowed
						break
					}
				}
				if len(allowed) == 0 {
					// Attract modes resolve styles but do not expose a style option.
					if alias.Mode != "attract_plasma" && alias.Mode != "attract_popsquares" && alias.Mode != "attract_waveform" {
						t.Errorf("%s: mode %q does not publish registered styles", product, alias.Mode)
					}
					continue
				}
				if !slices.Contains(allowed, alias.Style) {
					t.Errorf("%s: mode %q has unresolved alias %q", product, alias.Mode, alias.Style)
				}
			}
		}
	}
}

func TestThermalAliasesFitDefaultPanel(t *testing.T) {
	check := func(product, screen, controller string, config driver.DriverConfig) {
		t.Helper()
		if _, ok := style.LookupCompound(product, screen, "thermal"); !ok {
			return
		}

		d, ok := driver.Get(controller)
		if !ok {
			t.Errorf("%s/%s: unknown controller %q", product, screen, controller)
			return
		}
		hints := textlayout.TextHints{
			PanelProduct: product, ScreenName: screen, PixelWidth: config.Width,
			PixelHeight: config.Height, Capability: d.Capability,
		}
		selected, reason := style.ResolveStyle(thermal.ThermalRegistry, hints, "thermal", "")
		if reason != "alias" || selected.Supports(hints) == style.Unsupported {
			t.Errorf("%s/%s: thermal style %q (%s) does not fit %dx%d capability %d",
				product, screen, selected.Name(), reason, config.Width, config.Height, d.Capability)
		}
	}
	for _, panel := range panels.All() {
		if len(panel.Virtual) == 0 {
			check(panel.Name, "", panel.Controller, panel.Config)
		}
		for _, screen := range panel.Virtual {
			check(panel.Name, screen.Name, screen.Controller, screen.Config)
		}
	}
}
