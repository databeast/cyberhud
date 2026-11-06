package main

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 0.01 }

func TestBuildLayoutScaleAndCenter(t *testing.T) {
	// 320x240 virtual display fit into a 640x400 box with 8px padding.
	// avail = 624 x 384; scale = min(624/320, 384/240) = min(1.95, 1.6) = 1.6.
	st := &State{
		Regions: []Region{
			{ID: "main.0", Name: "main", Current: "clock"},
			{ID: "aux.0", Name: "aux", Current: "stemma"},
		},
		Geometry: Geometry{
			VirtualW: 320, VirtualH: 240,
			Regions: []RegionGeom{
				{Name: "main", X: 0, Y: 0, W: 240, H: 240},
				{Name: "aux", X: 240, Y: 0, W: 80, H: 120},
			},
		},
	}
	lv := BuildLayout(st, 640, 400)
	if !lv.Scaled {
		t.Fatal("expected Scaled=true for valid geometry")
	}

	const scale = 1.6
	wantW := 320 * scale // 512
	wantH := 240 * scale // 384
	if !approx(lv.OutlineW, wantW) || !approx(lv.OutlineH, wantH) {
		t.Errorf("outline size = %.2fx%.2f, want %.2fx%.2f", lv.OutlineW, lv.OutlineH, wantW, wantH)
	}
	// Centered: offX = (640-512)/2 = 64; offY = (400-384)/2 = 8.
	if !approx(lv.OutlineX, 64) || !approx(lv.OutlineY, 8) {
		t.Errorf("outline origin = (%.2f,%.2f), want (64,8)", lv.OutlineX, lv.OutlineY)
	}

	if len(lv.Regions) != 2 {
		t.Fatalf("want 2 placed regions, got %d", len(lv.Regions))
	}
	// main at (0,0) 240x240 -> offset (64,8), 384x384.
	m := lv.Regions[0]
	if m.ID != "main.0" || m.Current != "clock" {
		t.Errorf("main metadata mismatch: %+v", m)
	}
	if !approx(m.X, 64) || !approx(m.Y, 8) || !approx(m.W, 384) || !approx(m.H, 384) {
		t.Errorf("main rect = (%.1f,%.1f %.1fx%.1f), want (64,8 384x384)", m.X, m.Y, m.W, m.H)
	}
	// aux at (240,0) 80x120 -> x = 64 + 240*1.6 = 448; w = 128; h = 192.
	a := lv.Regions[1]
	if !approx(a.X, 448) || !approx(a.Y, 8) || !approx(a.W, 128) || !approx(a.H, 192) {
		t.Errorf("aux rect = (%.1f,%.1f %.1fx%.1f), want (448,8 128x192)", a.X, a.Y, a.W, a.H)
	}
}

func TestBuildLayoutNoGeometry(t *testing.T) {
	st := &State{
		Regions:  []Region{{ID: "main.0", Name: "main", Current: "clock"}},
		Geometry: Geometry{VirtualW: 0, VirtualH: 0},
	}
	lv := BuildLayout(st, 640, 400)
	if lv.Scaled {
		t.Error("expected Scaled=false when virtual size is zero")
	}
	if len(lv.Regions) != 0 {
		t.Errorf("expected no placed regions, got %d", len(lv.Regions))
	}
}

func TestBuildLayoutNilState(t *testing.T) {
	lv := BuildLayout(nil, 640, 400)
	if lv.Scaled || len(lv.Regions) != 0 {
		t.Errorf("nil state should yield empty unscaled layout, got %+v", lv)
	}
}
