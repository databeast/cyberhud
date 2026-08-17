package displaymode

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func TestRegisteredStyleNamesCount(t *testing.T) {
	names := registeredStyleNames()
	if len(names) != 24 {
		t.Fatalf("registeredStyleNames() returned %d entries, want 24", len(names))
	}
}

func TestCatalogDescribe(t *testing.T) {
	def, ok := catalog.Describe("displaymode")
	if !ok {
		t.Fatal("catalog.Describe(\"displaymode\") returned false")
	}
	if def.ID != "displaymode" {
		t.Fatalf("ID=%q, want \"displaymode\"", def.ID)
	}
	if def.Title != "Display Mode Template" {
		t.Fatalf("Title=%q, want \"Display Mode Template\"", def.Title)
	}
	if def.Order != 999 {
		t.Fatalf("Order=%d, want 999", def.Order)
	}
	if len(def.Options) == 0 {
		t.Fatal("no options registered")
	}
	if def.Options[0].Key != "style" {
		t.Fatalf("Options[0].Key=%q, want \"style\"", def.Options[0].Key)
	}
}

func TestHandleCommandQuery(t *testing.T) {
	SetPolicy(DefaultPolicy())
	result := HandleCommand(nil)
	want := "OK displaymode style=mono-128x32"
	if result != want {
		t.Fatalf("HandleCommand(nil)=%q, want %q", result, want)
	}
}

func TestHandleCommandSetColorStyle(t *testing.T) {
	SetPolicy(DefaultPolicy())
	result := HandleCommand([]string{"style=color-240x240"})
	if strings.HasPrefix(result, "ERR") {
		t.Fatalf("HandleCommand rejected valid style: %s", result)
	}
	if GetPolicy().Style != "color-240x240" {
		t.Fatalf("policy style=%q after setting color-240x240", GetPolicy().Style)
	}
}

func TestHandleCommandSetGrayscaleFastStyle(t *testing.T) {
	SetPolicy(DefaultPolicy())
	result := HandleCommand([]string{"style=grayscale-fast-240x240"})
	if strings.HasPrefix(result, "ERR") {
		t.Fatalf("HandleCommand rejected valid grayscale-fast style: %s", result)
	}
	if GetPolicy().Style != "grayscale-fast-240x240" {
		t.Fatalf("policy style=%q after setting grayscale-fast-240x240", GetPolicy().Style)
	}
}



func TestRenderCacheKeyNonEmpty(t *testing.T) {
	SetPolicy(DefaultPolicy())
	sig := RenderCacheKey()
	if sig == "" {
		t.Fatal("RenderCacheKey() returned empty string")
	}
}

func TestHandlerLeftAction(t *testing.T) {
	result := Handler{}.HandleAction(action.Left, 0, 0)
	if result != (action.Result{Dirty: true}) {
		t.Fatalf("HandleAction(Left)=%+v, want {Dirty:true}", result)
	}
}

func TestTwo800x480StylesDistinct(t *testing.T) {
	color800 := templateRegistry.Lookup("color-800x480")
	eink800 := templateRegistry.Lookup("eink-800x480")
	if color800 == nil || eink800 == nil {
		t.Fatal("one of the 800x480 styles is nil")
	}
	if color800.Name() == eink800.Name() {
		t.Fatal("800x480 styles have same name")
	}
	colorReqs := color800.Requirements()
	einkReqs := eink800.Requirements()
	if colorReqs.Capability == einkReqs.Capability {
		t.Fatal("800x480 styles have same Capability")
	}
}

func TestGrayscaleFastStyleCount(t *testing.T) {
	names := registeredStyleNames()
	count := 0
	for _, n := range names {
		if strings.HasPrefix(n, "grayscale-fast-") {
			count++
		}
	}
	if count != 7 {
		t.Fatalf("grayscale-fast style count=%d, want 7", count)
	}
}

func TestGrayscaleFastStyleDimensionsMatchColor(t *testing.T) {
	styles := templateRegistry.Enumerate()
	for _, s := range styles {
		name := s.Name()
		if !strings.HasPrefix(name, "grayscale-fast-") {
			continue
		}
		// Find corresponding color style
		colorName := "color-" + strings.TrimPrefix(name, "grayscale-fast-")
		colorStyle := templateRegistry.Lookup(colorName)
		if colorStyle == nil {
			t.Fatalf("grayscale-fast style %q has no matching color style %q", name, colorName)
		}
		dreqs := s.Requirements()
		creqs := colorStyle.Requirements()
		if dreqs.MinWidth != creqs.MinWidth || dreqs.MinHeight != creqs.MinHeight {
			t.Fatalf("grayscale-fast %q dimensions (%dx%d) != color %q (%dx%d)",
				name, dreqs.MinWidth, dreqs.MinHeight,
				colorName, creqs.MinWidth, creqs.MinHeight)
		}
	}
}

func TestGrayscaleFastStylesNoColorNoRapid(t *testing.T) {
	styles := templateRegistry.Enumerate()
	for _, s := range styles {
		if !strings.HasPrefix(s.Name(), "grayscale-fast-") {
			continue
		}
		reqs := s.Requirements()
		if reqs.Capability >= style.ColorSlow {
			t.Fatalf("grayscale-fast style %q has Capability >= ColorSlow (%v)", s.Name(), reqs.Capability)
		}
	}
}
