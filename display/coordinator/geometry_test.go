package coordinator_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/region"
)

// TestGeometryUnbound verifies that Geometry() returns a zero-value Geometry when
// the State has not been bound to a live RegionManager (the pre-bind compatibility
// snapshot carries no bounds).
func TestGeometryUnbound(t *testing.T) {
	s := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu"}, Default: "menu"},
	)

	geo := s.Geometry()
	if geo.VirtualW != 0 || geo.VirtualH != 0 {
		t.Fatalf("unbound Geometry: want virtual 0x0, got %dx%d", geo.VirtualW, geo.VirtualH)
	}
	if len(geo.Regions) != 0 {
		t.Fatalf("unbound Geometry: want 0 regions, got %d", len(geo.Regions))
	}
}

// TestGeometryBound verifies that once a live RegionManager is bound, Geometry()
// reports the virtual-display size and each region's rectangle in virtual-display
// coordinates, in allocation order.
func TestGeometryBound(t *testing.T) {
	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 320, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}
	rm := region.NewRegionManager(vd)

	specs := []region.RegionSpec{
		{Name: "main", Controller: "st7789", Bounds: image.Rect(0, 0, 240, 240), Modes: []string{"clock"}, DefaultMode: "clock"},
		{Name: "aux", Controller: "st7735s", Bounds: image.Rect(240, 0, 320, 120), Modes: []string{"clock"}, DefaultMode: "clock"},
	}
	for _, spec := range specs {
		if err := rm.Allocate(spec); err != nil {
			t.Fatalf("Allocate %q: %v", spec.Name, err)
		}
	}

	s := coordinator.NewState()
	s.Bind(rm)

	geo := s.Geometry()
	if geo.VirtualW != 320 || geo.VirtualH != 240 {
		t.Fatalf("bound Geometry: want virtual 320x240, got %dx%d", geo.VirtualW, geo.VirtualH)
	}
	if len(geo.Regions) != 2 {
		t.Fatalf("bound Geometry: want 2 regions, got %d", len(geo.Regions))
	}

	want := []coordinator.RegionGeometry{
		{Index: 0, Name: "main", X: 0, Y: 0, W: 240, H: 240},
		{Index: 1, Name: "aux", X: 240, Y: 0, W: 80, H: 120},
	}
	for i, w := range want {
		got := geo.Regions[i]
		if got != w {
			t.Errorf("region[%d]: want %+v, got %+v", i, w, got)
		}
	}
}
