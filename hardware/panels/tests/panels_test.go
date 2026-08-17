package tests_test

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/hardware/panels"
	_ "github.com/databeast/cyberhud/hardware/panels/all"
)

func TestSupportedIncludesKnownPanels(t *testing.T) {
	csv := panels.Supported()
	for _, want := range []string{"waveshare-1.3hat", "waveshare-triple-screen", "adafruit-2.13-ssd1680"} {
		if !strings.Contains(csv, want) {
			t.Fatalf("Supported() missing %q: %q", want, csv)
		}
	}
}

func TestGetPanel(t *testing.T) {
	p, err := panels.Get("WAVESHARE-1.3HAT")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	if p.Name != "waveshare-1.3hat" {
		t.Fatalf("Get().Name=%q", p.Name)
	}
	if p.Controller != "st7789" {
		t.Fatalf("Get().Controller=%q", p.Controller)
	}
}

func TestGetUnknownPanel(t *testing.T) {
	if _, err := panels.Get("definitely-missing"); err == nil {
		t.Fatal("Get() expected error for unknown panel")
	}
}

func TestRegisterPanel(t *testing.T) {
	const name = "test-registry-panel"
	panels.Register(panels.Definition{Name: name, Description: "test panel", Controller: "st7789"})
	p, err := panels.Get(name)
	if err != nil {
		t.Fatalf("Get(%q) error=%v", name, err)
	}
	if p.Name != name {
		t.Fatalf("Get(%q).Name=%q", name, p.Name)
	}
}

func TestNamesSorted(t *testing.T) {
	names := panels.Names()
	if len(names) == 0 {
		t.Fatal("Names() expected non-empty")
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Fatalf("Names() not sorted at %q > %q", names[i-1], names[i])
		}
	}
}
