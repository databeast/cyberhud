package schedule

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// showcaseSnapshot builds a representative schedule state anchored to the render
// instant: one event in progress right now (so the card shows its live accent
// title and time window) plus two upcoming entries (so ShowUpcoming has a "next:"
// line to preview). Anchoring to now at seed time is safe because the render path
// reads time.Now() microseconds later, within the active window.
//
// Without this the gallery rendered every panel as "No calendar set", which
// showcases none of the mode's actual layout.
func showcaseSnapshot() (source.Snapshot, source.Policy) {
	now := time.Now()
	events := []source.Event{
		{
			Summary:  "Standup",
			Location: "Room 2",
			Start:    now.Add(-15 * time.Minute),
			End:      now.Add(45 * time.Minute),
		},
		{
			Summary: "Design review",
			Start:   now.Add(90 * time.Minute),
			End:     now.Add(150 * time.Minute),
		},
		{
			Summary: "1:1 with Sam",
			Start:   now.Add(4 * time.Hour),
			End:     now.Add(4*time.Hour + 30*time.Minute),
		},
	}
	snap := source.Snapshot{
		Sequence:     1,
		Loaded:       true,
		Source:       "showcase.ics",
		FetchedAt:    now,
		Events:       events,
		CalendarName: "Work",
	}
	p := source.DefaultPolicy()
	p.Source = "showcase.ics"
	p.ShowUpcoming = true
	return snap, p
}

// snapshotOutputDir is the persistent directory where schedule snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// categoryFromCapability maps a style's declared hardware Capability to the
// appropriate testsnapshot.DisplayCategory, derived from the style's own
// SurfaceRequirements rather than a name heuristic, so styles added to the
// registry are handled automatically.
func categoryFromCapability(cap style.Capability) testsnapshot.DisplayCategory {
	switch cap {
	case style.MonoSlow:
		return testsnapshot.CategoryEink
	case style.MonoFast:
		return testsnapshot.CategoryMono
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	default:
		return testsnapshot.CategoryColor
	}
}

// TestSchedulePNGSnapshots enumerates every registered schedule style and
// renders a snapshot PNG for each through the full production render pipeline.
// The background calendar fetcher is never activated here, so each style renders
// its deterministic empty-calendar ("nothing scheduled") card.
func TestSchedulePNGSnapshots(t *testing.T) {
	styles := scheduleRegistry.Enumerate()
	if len(styles) == 0 {
		t.Fatal("scheduleRegistry contains zero styles")
	}

	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			category := categoryFromCapability(reqs.Capability)

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("schedule"),
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					source.ResetForTest()
				}),
				testsnapshot.WithPreRender(func() {
					// Seed a representative calendar, then point the policy at the
					// style under test. SetSnapshotForTest installs both snapshot
					// and policy atomically; set the target style on that policy.
					snap, p := showcaseSnapshot()
					p.Style = s.Name()
					source.SetSnapshotForTest(snap, p)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
