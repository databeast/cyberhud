package attract_cube

import (
	"bytes"
	"image"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func resetPolicy(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		if err := SetPolicy(DefaultPolicy()); err != nil {
			t.Error(err)
		}
	})
	if err := SetPolicy(DefaultPolicy()); err != nil {
		t.Fatal(err)
	}
}

func TestPolicyCommandAndPersistence(t *testing.T) {
	resetPolicy(t)
	if got := command.Handle(nil); got != "OK attract_cube base_color=cyan" {
		t.Fatalf("query: %s", got)
	}
	for _, tt := range []struct{ input, want string }{
		{"base_color=AMBER", "amber"}, {"base_color=#Ab12CD", "#ab12cd"},
	} {
		if got := command.Handle([]string{tt.input}); !strings.HasPrefix(got, "OK ") {
			t.Fatalf("command %q: %s", tt.input, got)
		}
		if got := GetPolicy().BaseColor; got != tt.want {
			t.Fatalf("color = %q, want %q", got, tt.want)
		}
	}
	before := GetPolicy()
	if got := command.Handle([]string{"base_color=#oops"}); !strings.HasPrefix(got, "ERR ") {
		t.Fatalf("invalid command: %s", got)
	}
	if GetPolicy() != before {
		t.Fatal("invalid command changed policy")
	}
	s := catalog.Snapshotter("attract_cube")
	if s == nil {
		t.Fatal("snapshotter not registered")
	}
	saved := s.SnapshotPolicy()
	if err := SetPolicy(DefaultPolicy()); err != nil {
		t.Fatal(err)
	}
	if err := s.RestorePolicy(saved); err != nil || GetPolicy() != before {
		t.Fatalf("restore: %v, policy %+v", err, GetPolicy())
	}
	for _, invalid := range []map[string]interface{}{{"base_color": 5}, {"base_color": "#xyzxyz"}} {
		if err := s.RestorePolicy(invalid); err == nil {
			t.Fatalf("invalid snapshot accepted: %v", invalid)
		}
		if GetPolicy() != before {
			t.Fatal("invalid restore changed policy")
		}
	}
	if err := s.RestorePolicy(map[string]interface{}{}); err != nil || GetPolicy() != DefaultPolicy() {
		t.Fatalf("empty restore did not use default: %v", err)
	}
}

func TestInstanceColorAndPanelCapabilities(t *testing.T) {
	resetPolicy(t)
	restore := frameclock.Freeze(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	defer restore()
	i := &instance{}
	i.SetPanelHints(textlayout.TextHints{PixelWidth: 240, PixelHeight: 240, Capability: textlayout.CapColorFast})
	i.Activate()
	firstKey := i.RenderCacheKey()
	first := i.BuildView()
	if first.Static || len(first.Sprites) != 1 {
		t.Fatalf("invalid color-fast view: %+v", first)
	}
	frameclock.Advance(time.Second)
	if i.RenderCacheKey() == firstKey {
		t.Fatal("animated key did not change")
	}
	current := i.BuildView()
	if err := SetPolicy(Policy{BaseColor: "red"}); err != nil {
		t.Fatal(err)
	}
	red := i.BuildView()
	if bytes.Equal(current.Sprites[0].Image.(*image.RGBA).Pix, red.Sprites[0].Image.(*image.RGBA).Pix) {
		t.Fatal("policy change did not affect rendered sprite")
	}
	i.SetPanelHints(textlayout.TextHints{PixelWidth: 240, PixelHeight: 240, Capability: textlayout.CapColorSlow})
	slowKey := i.RenderCacheKey()
	slow := i.BuildView()
	frameclock.Advance(time.Second)
	if !slow.Static || len(slow.Sprites) != 1 || i.RenderCacheKey() != slowKey {
		t.Fatal("slow color view must stay static")
	}
	if err := SetPolicy(Policy{BaseColor: "amber"}); err != nil {
		t.Fatal(err)
	}
	if i.RenderCacheKey() == slowKey {
		t.Fatal("slow cache key did not reflect policy change")
	}
	i.SetPanelHints(textlayout.TextHints{PixelWidth: 128, PixelHeight: 64, Capability: -1})
	unsupported := i.BuildView()
	if !unsupported.Static || len(unsupported.Sprites) != 0 || len(unsupported.Items) == 0 {
		t.Fatal("unsupported panel should receive an explanatory message")
	}
}
