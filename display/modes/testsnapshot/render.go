package testsnapshot

import (
	"fmt"
	"image"
	"path/filepath"
	"testing"
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	stemmapkg "github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/panels/pngpanel"
	"github.com/databeast/cyberhud/runtime/ui"
)

// RenderSnapshot renders a single frame of the specified display mode through
// the full production pipeline and writes it to a PNG file.
// Returns the absolute path to the output PNG or calls t.Fatal on error.
func RenderSnapshot(t *testing.T, opts ...SnapshotOption) string {
	t.Helper()
	cfg := applyDefaults(opts)
	cfg.t = t
	validate(t, cfg)
	prevScanner := stemmapkg.GlobalScanner()

	// 1. Freeze the frame clock before anything else.
	//
	// This precedes mode construction and activation deliberately. Modes record
	// a "last advanced at" baseline when they are activated, when their content
	// is replaced, and on their first render. A baseline taken from the wall
	// clock while later reads come from the frozen clock yields a hugely
	// negative elapsed time, which every mode treats as "no advancement" — so
	// animation would stall for the whole snapshot.
	restoreClock := frameclock.Freeze(FrameClockStart)
	defer restoreClock()

	// 2. Create PNGPanel with configured dimensions and color mode.
	panel := createPanel(t, cfg)

	// 2. Create Surface matching panel dimensions.
	bounds := image.Rect(0, 0, cfg.width, cfg.height)
	surf := surface.New(bounds)

	// 3. Create Region with PNGPanel's TextHints via ScreenPosition.
	screenPos := region.ScreenPosition{
		Bounds:       bounds,
		HintProvider: panel.TextHints,
	}
	reg := region.NewRegionWithScreens("test", bounds, surf, []region.ScreenPosition{screenPos}, "", 0, 0)

	// 4. Wire ModeFactory so SetMode can resolve mode IDs to instances.
	// Hints are propagated to mode packages via modehints before construction.
	reg.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		modehints.PropagateHints(hints)
		inst, ok := displaymodes.GetInstance(id)
		if !ok {
			return nil, false
		}
		return inst, true
	})

	// 5. Run Reset callback (clears global state before setup).
	if cfg.reset != nil {
		cfg.reset()
	}

	if cfg.scanner != nil {
		stemmapkg.SetGlobalScanner(cfg.scanner)
		defer stemmapkg.SetGlobalScanner(prevScanner)
	}

	// 6. Set the active mode on the region.
	_ = reg.SetMode(cfg.modeID)

	// Set mode-specific globals if configured.
	if cfg.warnings != nil {
		displaymodes.Warnings = cfg.warnings
	}

	// 8. Run PreRender callback (mode-specific state setup).
	if cfg.preRender != nil {
		cfg.preRender()
	}

	// 9. Build RegionRenderer (modes self-source their dependencies via
	// package singletons — no injected scanner/gpiomgr/modeState needed).
	renderer := ui.NewRegionRenderer(cfg.monochrome, nil, nil)

	// 10. Render frames, stepping the clock before each pass so time-gated
	// animation advances through the production tick path.
	//
	// RegionRenderer.Render calls region.Instance().BuildView() directly.
	renderPass := func(n int) {
		frameclock.Advance(cfg.frameInterval)
		if err := renderer.Render(reg); err != nil {
			t.Fatalf("snapshottest: RegionRenderer.Render failed on frame %d: %v", n, err)
		}
	}

	passes := 0
	for i := 0; i < cfg.frameCount; i++ {
		passes++
		renderPass(passes)
	}

	// When a readiness predicate is supplied the configured frame count is a
	// minimum: keep rendering until the mode reports it has something worth
	// capturing.
	if cfg.readyWhen != nil {
		for !cfg.readyWhen() {
			if passes >= MaxReadinessPasses {
				t.Fatalf("snapshottest: mode %q at %dx%d never became ready after %d render passes "+
					"(frame interval %v, elapsed %v)",
					cfg.modeID, cfg.width, cfg.height, passes,
					cfg.frameInterval, time.Duration(passes)*cfg.frameInterval)
			}
			passes++
			renderPass(passes)
		}
	}

	// 13. Write framebuffer to PNGPanel.
	if err := panel.DrawImage(surf.FrameBuffer()); err != nil {
		t.Fatalf("snapshottest: panel.DrawImage failed: %v", err)
	}

	// 14. Deactivate the mode to stop any background goroutines (e.g., thermal
	// sampling loop) that were started by Activate() during SetMode.
	if inst := reg.Instance(); inst != nil {
		inst.Deactivate()
	}

	// 15. Compute and return output path.
	return outputPath(cfg)
}

// createPanel constructs a PNGPanel with the configured dimensions and color mode.
func createPanel(t *testing.T, cfg *snapshotConfig) *pngpanel.PNGPanel {
	t.Helper()

	panel, err := pngpanel.New(
		pngpanel.WithDimensions(cfg.width, cfg.height),
		pngpanel.WithColorMode(cfg.colorMode),
		pngpanel.WithOutputDir(cfg.outputDir),
		pngpanel.WithBasename(cfg.basename),
	)
	if err != nil {
		t.Fatalf("snapshottest: pngpanel.New failed: %v", err)
	}
	return panel
}

// outputPath computes the path to the rendered PNG file.
// The framework calls DrawImage exactly once on a fresh PNGPanel, so the
// panel's internal counter will always be 1. The filename is always _0001.png.
func outputPath(cfg *snapshotConfig) string {
	return filepath.Join(cfg.outputDir, fmt.Sprintf("%s_%04d.png", cfg.basename, 1))
}
