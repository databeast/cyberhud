package main

import (
	"reflect"
	"testing"
)

func TestParseStatus(t *testing.T) {
	lines := []string{
		"OK",
		`  region=main.0 name="main" controller=st7789 mode=clock modes=clock,dashboard`,
		`  region=left-aux.0 name="left-aux" controller=st7735s mode=stemma modes=stemma,gpio`,
	}
	got := parseStatus(lines)
	if len(got) != 2 {
		t.Fatalf("want 2 regions, got %d: %+v", len(got), got)
	}

	m := got[0]
	if m.ID != "main.0" || m.Name != "main" || m.Index != 0 || m.Controller != "st7789" || m.Current != "clock" {
		t.Errorf("region 0 mismatch: %+v", m)
	}
	if len(m.Modes) != 2 || m.Modes[0].ID != "clock" || m.Modes[1].ID != "dashboard" {
		t.Errorf("region 0 modes mismatch: %+v", m.Modes)
	}

	a := got[1]
	if a.ID != "left-aux.0" || a.Name != "left-aux" || a.Current != "stemma" {
		t.Errorf("region 1 mismatch: %+v", a)
	}
}

func TestParseStatusQuotedNameWithSpace(t *testing.T) {
	lines := []string{
		"OK",
		`  region=main.0 name="main panel" controller=st7789 mode=clock modes=clock`,
	}
	got := parseStatus(lines)
	if len(got) != 1 {
		t.Fatalf("want 1 region, got %d", len(got))
	}
	if got[0].Name != "main panel" {
		t.Errorf("quoted name with space: want %q, got %q", "main panel", got[0].Name)
	}
}

func TestParseModes(t *testing.T) {
	lines := []string{
		"OK",
		`  region=main.0 name="main" controller=st7789 current=clock`,
		`    mode=clock title="Clock" scope=any summary="A simple clock."`,
		`    mode=dashboard title="Dashboard" summary="System dashboard with stats."`,
		`  region=aux.0 name="aux" current=stemma`,
		`    mode=stemma title="STEMMA" summary="Detected I2C devices."`,
	}
	got := parseModes(lines)

	main, ok := got["main.0"]
	if !ok {
		t.Fatalf("missing main.0 in %+v", got)
	}
	if len(main) != 2 {
		t.Fatalf("main.0: want 2 modes, got %d", len(main))
	}
	if main[0].ID != "clock" || main[0].Title != "Clock" || main[0].Summary != "A simple clock." {
		t.Errorf("main.0 mode 0 mismatch: %+v", main[0])
	}
	if main[1].ID != "dashboard" || main[1].Summary != "System dashboard with stats." {
		t.Errorf("main.0 mode 1 mismatch: %+v", main[1])
	}

	aux, ok := got["aux.0"]
	if !ok || len(aux) != 1 || aux[0].ID != "stemma" {
		t.Errorf("aux.0 mismatch: %+v", aux)
	}
}

func TestParseModesWithOptions(t *testing.T) {
	lines := []string{
		"OK",
		`  region=main.0 name="main" controller=st7789 current=image`,
		`    mode=image title="Image" scope=any summary="Full-panel image."`,
		`      option=fit type=string default="scale" summary="How the image fits." allowed="truncate,scale,stretch"`,
		`      option=auto_scroll_ms type=int default="0" summary="Scroll interval."`,
		`      option=show_border type=bool default="false" summary="Border frame."`,
		`    mode=clock title="Clock" summary="A clock."`,
	}
	got := parseModes(lines)
	modes := got["main.0"]
	if len(modes) != 2 {
		t.Fatalf("want 2 modes, got %d: %+v", len(modes), modes)
	}

	img := modes[0]
	if img.ID != "image" || len(img.Options) != 3 {
		t.Fatalf("image options mismatch: %+v", img)
	}

	fit := img.Options[0]
	if fit.Key != "fit" || fit.Type != "string" || fit.Default != "scale" {
		t.Errorf("fit option mismatch: %+v", fit)
	}
	if !fit.isEnum() || len(fit.Allowed) != 3 || fit.Allowed[0] != "truncate" || fit.Allowed[2] != "stretch" {
		t.Errorf("fit allowed mismatch: %+v", fit.Allowed)
	}

	if img.Options[1].Key != "auto_scroll_ms" || img.Options[1].Type != "int" || img.Options[1].isEnum() {
		t.Errorf("int option mismatch: %+v", img.Options[1])
	}
	if img.Options[2].Key != "show_border" || img.Options[2].Type != "bool" {
		t.Errorf("bool option mismatch: %+v", img.Options[2])
	}

	// The second mode (clock) must NOT inherit image's options.
	if len(modes[1].Options) != 0 {
		t.Errorf("clock should have no options, got %+v", modes[1].Options)
	}
}

func TestParsePolicy(t *testing.T) {
	cases := []struct {
		name string
		line string
		want map[string]string
	}{
		{
			name: "mode-prefixed with mixed types",
			line: `OK mode=ticker {"style":"bordered","auto_scroll_ms":250,"show_border":true,"accent":"cyan"}`,
			want: map[string]string{"style": "bordered", "auto_scroll_ms": "250", "show_border": "true", "accent": "cyan"},
		},
		{
			name: "bare json",
			line: `OK {"fit":"scale","style":""}`,
			want: map[string]string{"fit": "scale", "style": ""},
		},
		{
			name: "float renders without trailing zero",
			line: `OK mode=attract {"speed":2,"density":0.75}`,
			want: map[string]string{"speed": "2", "density": "0.75"},
		},
		{
			name: "empty object",
			line: `OK mode=clock {}`,
			want: map[string]string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parsePolicy([]string{tc.line})
			if len(got) != len(tc.want) {
				t.Fatalf("len mismatch: got %v want %v", got, tc.want)
			}
			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("key %q: got %q want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestEditFields(t *testing.T) {
	m := &ModeRef{
		ID: "image",
		Options: []Option{
			{Key: "fit", Type: "string", Default: "scale", Allowed: []string{"truncate", "scale", "stretch"}},
			{Key: "show_border", Type: "bool", Default: "false"},
		},
		Policy: map[string]string{"fit": "stretch"}, // show_border not in policy → falls back to default
	}
	fields := m.EditFields()
	if len(fields) != 2 {
		t.Fatalf("want 2 fields, got %d", len(fields))
	}
	if fields[0].Key != "fit" || fields[0].Value != "stretch" {
		t.Errorf("fit field should use live value: %+v", fields[0])
	}
	if fields[1].Key != "show_border" || fields[1].Value != "false" {
		t.Errorf("show_border should fall back to default: %+v", fields[1])
	}

	// A mode with no options yields no fields.
	m2 := &ModeRef{ID: "clock"}
	if f := m2.EditFields(); f != nil {
		t.Errorf("clock should have no edit fields, got %+v", f)
	}
}

func TestParseGeometry(t *testing.T) {
	lines := []string{
		"OK",
		"  virtual w=320 h=240",
		`  region=main.0 index=0 name="main" x=0 y=0 w=240 h=240`,
		`  region=aux.0 index=1 name="aux" x=240 y=0 w=80 h=120`,
	}
	got := parseGeometry(lines)
	if got.VirtualW != 320 || got.VirtualH != 240 {
		t.Fatalf("virtual size: want 320x240, got %dx%d", got.VirtualW, got.VirtualH)
	}
	want := []RegionGeom{
		{Index: 0, Name: "main", X: 0, Y: 0, W: 240, H: 240},
		{Index: 1, Name: "aux", X: 240, Y: 0, W: 80, H: 120},
	}
	if !reflect.DeepEqual(got.Regions, want) {
		t.Errorf("geometry regions mismatch:\n got  %+v\n want %+v", got.Regions, want)
	}
}

func TestParseGeometryEmpty(t *testing.T) {
	// Unbound daemon: just OK + virtual 0x0.
	lines := []string{"OK", "  virtual w=0 h=0"}
	got := parseGeometry(lines)
	if got.VirtualW != 0 || got.VirtualH != 0 || len(got.Regions) != 0 {
		t.Errorf("empty geometry mismatch: %+v", got)
	}
}

func TestMergeModes(t *testing.T) {
	regions := []Region{
		{ID: "main.0", Name: "main", Current: "clock", Modes: []ModeRef{{ID: "clock"}, {ID: "dashboard"}}},
		{ID: "aux.0", Name: "aux", Current: "stemma", Modes: []ModeRef{{ID: "stemma"}}},
	}
	modesByID := map[string][]ModeRef{
		"main.0": {{ID: "clock", Title: "Clock"}, {ID: "dashboard", Title: "Dashboard"}},
	}
	got := mergeModes(regions, modesByID)

	if got[0].Modes[0].Title != "Clock" || got[0].Modes[1].Title != "Dashboard" {
		t.Errorf("main.0 not enriched: %+v", got[0].Modes)
	}
	// aux.0 had no rich entry; should keep its bare mode.
	if len(got[1].Modes) != 1 || got[1].Modes[0].ID != "stemma" || got[1].Modes[0].Title != "" {
		t.Errorf("aux.0 should keep bare mode: %+v", got[1].Modes)
	}
}

func TestSplitTokensQuoted(t *testing.T) {
	toks := splitTokens(`region=main.0 name="a b" x=1`)
	want := []string{"region=main.0", `name="a b"`, "x=1"}
	if !reflect.DeepEqual(toks, want) {
		t.Errorf("splitTokens mismatch:\n got  %q\n want %q", toks, want)
	}
}
