package pager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where pager snapshot PNGs are
// written.
//
// Flat by requirement, not by preference: the docsnap collector reads only the
// top-level entries of <mode>/testdata/snapshots and skips directories, so a
// PNG in a subdirectory never reaches the documentation gallery. Test functions
// therefore partition this directory by filename prefix rather than by
// subdirectory.
var snapshotOutputDir = filepath.Join("snapshots")

// snapshotFileName is the file the framework writes for a given basename.
//
// The framework calls DrawImage exactly once per snapshot on a fresh panel, so
// the frame counter in the name is always 1.
func snapshotFileName(basename string) string {
	return basename + "_0001.png"
}

// categoryFromCapability maps a style's declared capability to the panel colour
// category the snapshot renders on.
//
// Derived from the declared capability rather than parsed out of the style name,
// so a style rename cannot silently reclassify a panel.
//
// Note that the resulting panel capability is not the same as the style's
// declared one: a monochrome PNG panel reports CapMonoSlow with
// PreferEventRefresh set, so every mono style — including those declaring
// MonoFast — renders through the slow page-transition path. That is the panel's
// behaviour, and the snapshot reflects it rather than working around it.
func categoryFromCapability(c style.Capability) testsnapshot.DisplayCategory {
	switch c {
	case style.MonoFast:
		return testsnapshot.CategoryMono
	case style.MonoSlow:
		return testsnapshot.CategoryEink
	case style.GrayscaleFast, style.GrayscaleSlow:
		return testsnapshot.CategoryGrayscale
	default:
		return testsnapshot.CategoryColor
	}
}

// sampleLines generates representative text content for snapshot rendering.
func sampleLines(count int) []byte {
	var data []byte
	for i := 0; i < count; i++ {
		line := fmt.Sprintf("Line %02d: The quick brown fox jumps over the lazy dog", i+1)
		data = append(data, []byte(line+"\n")...)
	}
	return data
}

// purgeOwnOutput removes exactly the files this test function is about to
// rewrite, leaving everything else in the directory alone.
//
// Scoping deletion to the owned set is what makes output order-independent. The
// previous implementation removed the whole snapshots tree, which deleted a
// sibling test function's output and left the surviving file set dependent on
// execution order. Matching on the expected names rather than a filename prefix
// also covers styles whose names do not share a common prefix.
func purgeOwnOutput(t *testing.T, dir string, basenames []string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("creating snapshot output dir %q: %v", dir, err)
	}
	for _, name := range basenames {
		path := filepath.Join(dir, snapshotFileName(name))
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			t.Fatalf("removing stale snapshot %q: %v", path, err)
		}
	}
}

// TestPagerPNGSnapshots renders every registered pager style through the full
// production pipeline, one frame per style.
//
// Every style is rendered at exactly the geometry it declares, and no style is
// skipped. Animation state is advanced only by the production tick path, driven
// by the framework's frame clock — the test configures policy and buffer
// content and nothing else.
func TestPagerPNGSnapshots(t *testing.T) {
	styles := pagerRegistry.Enumerate()
	if len(styles) == 0 {
		t.Fatal("pagerRegistry contains zero styles")
	}

	basenames := make([]string, 0, len(styles))
	for _, s := range styles {
		basenames = append(basenames, s.Name())
	}
	purgeOwnOutput(t, snapshotOutputDir, basenames)

	var pngPaths []string

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			width, height := declaredGeometry(t, s.Name(), reqs)

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("pager"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(categoryFromCapability(reqs.Capability)),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				// A handful of passes is enough for the smooth-scroll path to
				// accumulate a visible offset. The page-transition path needs
				// far longer, and is handled by the readiness predicate below.
				testsnapshot.WithFrameCount(4),
				// BuildView caps per-frame delta at 80ms, so stepping the clock
				// by exactly that much advances the most simulated time per
				// pass without any of it being discarded. That matters for the
				// page path, which must cover a multi-second cadence.
				testsnapshot.WithFrameInterval(80*time.Millisecond),
				testsnapshot.WithReset(func() {
					SetPolicy(DefaultPolicy())
				}),
				testsnapshot.WithPreRender(func() {
					pol := DefaultPolicy()
					// A path that cannot exist: the reader retries at scan_ms
					// and never ingests, so buffer content comes only from this
					// test and the snapshot stays deterministic. A non-empty
					// source is still required to get past the no-source view.
					pol.Source = filepath.Join(os.TempDir(), "cyberhud-pager-snapshot-absent")
					// max_lines must exceed the visible row count of the
					// largest panel, or the slow path never sees a full page
					// and falls back to its max_wait_s partial-page timeout.
					pol.MaxLines = 200
					// Page cadence is max(3000ms, visible_rows × line_time_ms).
					// Left at the 1000ms default, a 40-row panel would need 40
					// simulated seconds — hundreds of render passes — to reach
					// its first settled page. Driving line_time_ms to its
					// minimum pins cadence at the 3000ms floor for every panel.
					pol.LineTimeMS = 1
					SetPolicy(pol)

					_, buf, _, _ := source.ActiveStateSnapshot()
					if buf != nil {
						buf.Ingest(sampleLines(120))
					}
				}),
				testsnapshot.WithReadyWhen(pagerHasVisibleContent),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
			pngPaths = append(pngPaths, pngPath)
		})
	}

	if err := testsnapshot.WriteGalleryFragmentFromPaths(snapshotOutputDir, "pager", pngPaths); err != nil {
		t.Fatalf("failed to write gallery fragment: %v", err)
	}
}

// declaredGeometry resolves the pixel size a style is rendered at.
//
// No minimum-size floor is applied. Clamping small dimensions upward produced
// snapshots of geometries no hardware provides, which defeats the purpose of
// capturing per-panel output.
func declaredGeometry(t *testing.T, name string, reqs style.SurfaceRequirements) (int, int) {
	t.Helper()

	width, height := reqs.MinWidth, reqs.MinHeight
	if width == 0 {
		width = reqs.PreferredWidth
	}
	if height == 0 {
		height = reqs.PreferredHeight
	}
	if width <= 0 || height <= 0 {
		t.Fatalf("style %q declares no usable geometry (min %dx%d, preferred %dx%d)",
			name, reqs.MinWidth, reqs.MinHeight, reqs.PreferredWidth, reqs.PreferredHeight)
	}
	return width, height
}

// pagerHasVisibleContent reports whether the pager has reached a state whose
// rendered frame shows content.
//
// Which condition applies depends on the path the panel took. The smooth-scroll
// path is ready once its offset has advanced. The page-transition path is ready
// once it settles back to idle at full opacity holding a non-empty page, since
// a mid-fade frame is a partially transparent snapshot of a transition rather
// than a representative still.
func pagerHasVisibleContent() bool {
	_, _, scroll, page := source.ActiveStateSnapshot()

	switch {
	case scroll != nil:
		return scroll.OffsetPx() > 0
	case page != nil:
		if page.Phase() != source.PhaseIdle || page.FadeAlpha() < 1.0 {
			return false
		}
		for _, line := range page.CurrentPage() {
			if strings.TrimSpace(line) != "" {
				return true
			}
		}
		return false
	default:
		return true
	}
}
