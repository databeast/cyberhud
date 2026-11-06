package pager

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestPageRenderUsesElapsedTimeAndHoldsWhenCaughtUp(t *testing.T) {
	restoreClock := frameclock.Freeze(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	defer restoreClock()
	oldPolicy := GetPolicy()
	defer SetPolicy(oldPolicy)

	pol := DefaultPolicy()
	pol.LineTimeMS = 1
	pol.FadeOutMS = 0
	pol.FadeInMS = 0
	pol.MaxWaitS = 2
	SetPolicy(pol)

	hints := textlayout.TextHints{PixelWidth: 128, PixelHeight: 32, Capability: textlayout.CapMonoSlow}
	layout := source.ComputeLayout(style.NewStyleContext(hints))
	if layout.VisibleRows <= 0 || layout.VisibleColumns <= 0 {
		t.Fatal("test panel has no usable text area")
	}

	i := &instance{}
	i.SetPanelHints(hints)
	i.Activate()
	defer i.Deactivate()
	pol.Source = "/unused"
	SetPolicy(pol)

	_, buf, _, page := source.ActiveStateSnapshot()
	buf.Ingest([]byte(strings.Repeat("first\n", layout.VisibleRows)))
	i.BuildView()
	for n := 0; n < 3; n++ {
		frameclock.Advance(time.Second)
		i.BuildView()
	}
	if page.Phase() != source.PhaseFadeOut {
		t.Fatalf("phase after 3s = %v, want fade out", page.Phase())
	}

	for n := 0; n < 2; n++ {
		frameclock.Advance(time.Second)
		i.BuildView()
	}
	if page.Phase() != source.PhaseIdle || len(page.CurrentPage()) != layout.VisibleRows {
		t.Fatalf("page did not settle: phase %v, lines %q", page.Phase(), page.CurrentPage())
	}
	frameclock.Advance(10 * time.Second)
	i.BuildView()
	if page.Phase() != source.PhaseIdle {
		t.Fatalf("phase with no unread lines = %v, want idle", page.Phase())
	}

	buf.Ingest([]byte(strings.Repeat("second\n", layout.VisibleRows)))
	for n := 0; n < 5; n++ {
		frameclock.Advance(time.Second)
		i.BuildView()
	}
	if got, want := page.CurrentPage(), strings.Fields(strings.Repeat("second ", layout.VisibleRows)); !reflect.DeepEqual(got, want) {
		t.Fatalf("second page = %q, want %q", got, want)
	}
}
