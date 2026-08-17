package source

import "testing"

func TestNormalizePolicy(t *testing.T) {
	p := NormalizePolicy(Policy{
		Shape:      "bogus",
		LabelTier:  "bogus",
		Accent:     "bogus",
		DefaultMin: 10,
		DefaultMax: 10,
		Columns:    -1,
		Rows:       -2,
		TileGapPx:  -3,
		PaddingPct: 40,
	})

	if p.Shape != "auto" {
		t.Fatalf("Shape=%q, want auto", p.Shape)
	}
	if p.LabelTier != "normal" {
		t.Fatalf("LabelTier=%q, want normal", p.LabelTier)
	}
	if p.Accent != "cyan" {
		t.Fatalf("Accent=%q, want cyan", p.Accent)
	}
	if p.DefaultMax <= p.DefaultMin {
		t.Fatalf("DefaultMax=%v must be greater than DefaultMin=%v", p.DefaultMax, p.DefaultMin)
	}
	if p.Columns != 0 || p.Rows != 0 || p.TileGapPx != 0 || p.PaddingPct != 25 {
		t.Fatalf("unexpected clamping result: %#v", p)
	}
}

func TestParsePayloadUsesPolicyDefaults(t *testing.T) {
	pol := Policy{DefaultMin: 25, DefaultMax: 75, Shape: "linear", Accent: "amber", LabelTier: "normal"}
	snap, err := ParsePayload("50", pol)
	if err != nil {
		t.Fatalf("ParsePayload returned error: %v", err)
	}
	if len(snap.Gauges) != 1 {
		t.Fatalf("len(Gauges)=%d, want 1", len(snap.Gauges))
	}
	g := snap.Gauges[0]
	if g.Min != 25 || g.Max != 75 {
		t.Fatalf("Min/Max = %v/%v, want 25/75", g.Min, g.Max)
	}
	if g.Percent != 0.5 {
		t.Fatalf("Percent=%v, want 0.5", g.Percent)
	}
	if g.Shape != "linear" {
		t.Fatalf("Shape=%q, want linear", g.Shape)
	}
}

func TestParsePayloadArrayAndObject(t *testing.T) {
	pol := DefaultPolicy()
	arraySnap, err := ParsePayload(`[{"label":"cpu","value":42,"min":0,"max":100},{"label":"mem","value":88}]`, pol)
	if err != nil {
		t.Fatalf("ParsePayload(array) returned error: %v", err)
	}
	if len(arraySnap.Gauges) != 2 {
		t.Fatalf("len(arraySnap.Gauges)=%d, want 2", len(arraySnap.Gauges))
	}

	objectSnap, err := ParsePayload(`{"label":"disk","value":12,"min":0,"max":24}`, pol)
	if err != nil {
		t.Fatalf("ParsePayload(object) returned error: %v", err)
	}
	if len(objectSnap.Gauges) != 1 {
		t.Fatalf("len(objectSnap.Gauges)=%d, want 1", len(objectSnap.Gauges))
	}
	if objectSnap.Gauges[0].Label != "disk" {
		t.Fatalf("Label=%q, want disk", objectSnap.Gauges[0].Label)
	}
}
