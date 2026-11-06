package pager

import (
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestTyperCommandAndPolicyPersistence(t *testing.T) {
	old := GetPolicy()
	defer SetPolicy(old)
	SetPolicy(DefaultPolicy())

	response := HandleCommand([]string{"mode=typer", "type_speed=42"})
	if !strings.HasPrefix(response, "OK ") || !strings.Contains(response, "mode=typer") || !strings.Contains(response, "type_speed=42") {
		t.Fatalf("unexpected command response: %s", response)
	}
	snapshot := (pagerSnapshotter{}).SnapshotPolicy()
	SetPolicy(DefaultPolicy())
	if err := (pagerSnapshotter{}).RestorePolicy(snapshot); err != nil {
		t.Fatal(err)
	}
	if p := GetPolicy(); p.Mode != "typer" || p.TypeSpeed != 42 {
		t.Fatalf("restored policy: %+v", p)
	}
	if response := HandleCommand([]string{"mode=scroll", "type_speed=0"}); !strings.HasPrefix(response, "ERR ") {
		t.Fatalf("invalid speed accepted: %s", response)
	}
	if p := GetPolicy(); p.Mode != "typer" || p.TypeSpeed != 42 {
		t.Fatalf("invalid command changed policy: %+v", p)
	}
	if response := HandleCommand([]string{"mode=unknown"}); !strings.HasPrefix(response, "ERR ") {
		t.Fatalf("invalid mode accepted: %s", response)
	}
	if err := (pagerSnapshotter{}).RestorePolicy(map[string]interface{}{}); err != nil {
		t.Fatal(err)
	}
	if p := GetPolicy(); p.Mode != "scroll" || p.TypeSpeed != 30 {
		t.Fatalf("legacy policy defaults: %+v", p)
	}
	if err := (pagerSnapshotter{}).RestorePolicy(map[string]interface{}{"mode": "typer", "type_speed": 55}); err != nil {
		t.Fatal(err)
	}
	if p := GetPolicy(); p.Mode != "typer" || p.TypeSpeed != 55 {
		t.Fatalf("partial policy restore: %+v", p)
	}
}

func TestTyperFastRenderAndSlowFallback(t *testing.T) {
	restore := frameclock.Freeze(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	defer restore()
	old := GetPolicy()
	defer SetPolicy(old)
	p := DefaultPolicy()
	p.Mode = "typer"
	p.TypeSpeed = 10
	SetPolicy(p)
	i := &instance{}
	hints := textlayout.TextHints{PixelWidth: 128, PixelHeight: 64, Capability: textlayout.CapMonoFast}
	i.SetPanelHints(hints)
	i.Activate()
	defer i.Deactivate()
	p.Source = "/unused"
	SetPolicy(p)

	_, buf, scroll, page := source.ActiveStateSnapshot()
	if scroll != nil || page != nil || i.typer == nil {
		t.Fatal("typer strategy not installed on fast display")
	}
	buf.Ingest([]byte("abcd"))
	i.BuildView()
	key := i.RenderCacheKey()
	frameclock.Advance(100 * time.Millisecond)
	view := i.BuildView()
	if len(view.Items) != 1 || view.Items[0] != "a" || len(view.LineOffsets) != 1 || view.LineOffsets[0] != 0 {
		t.Fatalf("unexpected first typer frame: %+v", view)
	}
	if i.RenderCacheKey() == key {
		t.Fatal("cache key did not change when a character was revealed")
	}
	frameclock.Advance(200 * time.Millisecond)
	if view = i.BuildView(); len(view.Items) != 1 || view.Items[0] != "abc" {
		t.Fatalf("unexpected later typer frame: %+v", view)
	}

	p.Mode = "scroll"
	SetPolicy(p)
	i.BuildView()
	_, _, scroll, page = source.ActiveStateSnapshot()
	if i.typer != nil || scroll == nil || page != nil {
		t.Fatal("switch to scroll did not install scroll strategy")
	}

	p.Mode = "typer"
	SetPolicy(p)
	hints.Capability = textlayout.CapMonoSlow
	i.SetPanelHints(hints)
	i.BuildView()
	_, _, scroll, page = source.ActiveStateSnapshot()
	if i.typer != nil || scroll != nil || page == nil {
		t.Fatal("slow display did not retain page strategy")
	}
}
