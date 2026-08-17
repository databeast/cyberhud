// Package snapshottest provides a shared testing framework for producing PNG
// snapshot output of any registered display mode through the full production
// render pipeline. The framework wires up PNGPanel → Region → ModeInstance →
// RegionRenderer.Render, ensuring test output is bit-for-bit identical to what
// real hardware produces.
//
// Test authors call RenderSnapshot with functional options to configure
// dimensions, color mode, display mode, and runtime dependencies. The framework
// handles all infrastructure wiring so tests only specify what varies.
//
// Font registration is automatic via transitive imports — all fonts registered
// in display/surface/font are available without explicit blank imports.
package testsnapshot

import (
	"image"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	stemmapkg "github.com/databeast/cyberhud/display/modes/stemma/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/hardware/panels/pngpanel"
)

// DisplayCategory represents the logical category of a display panel.
// It determines both the PNGPanel ColorMode and the RegionRenderer monochrome
// flag via the category-to-colorMode mapping.
type DisplayCategory int

const (
	// CategoryColor maps to ColorModeFullColor with monochrome=false.
	// Used for color TFT panels (ST7789, ILI9341).
	CategoryColor DisplayCategory = iota

	// CategoryMono maps to ColorModeMonochrome with monochrome=true.
	// Used for monochrome OLED panels (SH1106, SSD1306).
	CategoryMono

	// CategoryEink maps to ColorModeMonochrome with monochrome=true.
	// Used for e-ink displays (no rapid refresh).
	CategoryEink

	// CategoryGrayscale maps to ColorModeGrayscale with monochrome=false.
	// Used for grayscale displays.
	CategoryGrayscale
)

// Frame clock bounds and the fixed instant every snapshot starts from.
const (
	// DefaultFrameInterval approximates the 30fps cadence the fast-refresh
	// render loop targets on hardware.
	DefaultFrameInterval = 33 * time.Millisecond

	// MinFrameInterval and MaxFrameInterval bound the per-pass advance. A
	// zero or negative interval would leave time-gated modes stalled, which is
	// the bug this mechanism exists to prevent, so it is rejected rather than
	// silently accepted.
	MinFrameInterval = time.Millisecond
	MaxFrameInterval = time.Minute

	// MaxFrameCount caps configured passes. The readiness predicate has its own
	// separate cap, since it may legitimately need more.
	MaxFrameCount = 200

	// MaxReadinessPasses caps how long the framework will keep rendering while
	// waiting for a readiness predicate, so a mode that never becomes ready
	// fails with a diagnosable message instead of looping forever.
	MaxReadinessPasses = 2000
)

// FrameClockStart is the instant every snapshot's frame clock begins at.
//
// Fixed so that any timestamp a mode derives from the clock is identical across
// runs and machines, which is what makes repeated renders byte-identical.
var FrameClockStart = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

// SnapshotOption is a functional option for configuring a RenderSnapshot call.
type SnapshotOption func(*snapshotConfig)

// snapshotConfig holds all configuration for a single RenderSnapshot invocation.
type snapshotConfig struct {
	t             *testing.T
	width         int
	height        int
	colorMode     pngpanel.ColorMode
	monochrome    bool
	modeID        string
	outputDir     string
	basename      string
	padding       int
	reset         func()
	preRender     func()
	frameCount    int
	frameInterval time.Duration
	readyWhen     func() bool
	// colorModeSet tracks whether WithColorMode was explicitly called.
	// Needed because ColorModeFullColor is the zero value.
	colorModeSet bool

	// displayCategory is nil when not set (mutually exclusive with explicit colorMode).
	displayCategory *DisplayCategory

	// Mode-specific globals.
	warnings []string

	// RegionRenderer dependency options.
	scanner   *stemmapkg.Scanner
	gpiomgr   *gpiomgr.Manager
	modeState *coordinator.State
}

// WithDimensions sets the panel width and height in pixels.
// Both dimensions are required and must be in (0, 4096].
func WithDimensions(width, height int) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.width = width
		cfg.height = height
	}
}

// WithMode sets the display mode ID to render.
// The mode must be registered in the display modes registry.
func WithMode(modeID string) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.modeID = modeID
	}
}

// WithOutputDir sets the directory where the output PNG will be written.
func WithOutputDir(dir string) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.outputDir = dir
	}
}

// WithColorMode sets the PNGPanel color mode directly.
// Mutually exclusive with WithDisplayCategory — validate() will fatal if both are set.
// When used without WithDisplayCategory, monochrome is derived automatically:
// ColorModeMonochrome → monochrome=true, all others → monochrome=false.
func WithColorMode(mode pngpanel.ColorMode) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.colorMode = mode
		cfg.colorModeSet = true
	}
}

// WithDisplayCategory sets the logical display category, which derives both
// the PNGPanel ColorMode and RegionRenderer monochrome flag from the mapping.
// Mutually exclusive with WithColorMode.
func WithDisplayCategory(category DisplayCategory) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.displayCategory = &category
	}
}

// WithPadding sets the padding percentage used for LayoutBridge content-area
// computation during tiercatalog pre-validation. The value is stored in
// BridgeConfig.PaddingPct when constructing the LayoutBridge for font
// qualification checks. Valid range is [0, 50]; default is 0.
func WithPadding(pct int) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.padding = pct
	}
}

// WithGPIOManager injects a GPIO manager for modes that need it.
// Passed through to the RegionRenderer as a dependency.
func WithGPIOManager(g *gpiomgr.Manager) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.gpiomgr = g
	}
}

// WithScanner injects a STEMMA scanner for modes that need it.
// Passed through to the RegionRenderer as a dependency.
func WithScanner(s *stemmapkg.Scanner) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.scanner = s
	}
}

// WithModeState injects mode state for cross-panel coordination.
// Passed through to the RegionRenderer as a dependency.
func WithModeState(ms *coordinator.State) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.modeState = ms
	}
}

// WithReset sets a callback that runs BEFORE PreRender.
// Use this to clear global state (e.g., clock.ResetPolicy(), gpio.ResetPolicy()).
func WithReset(fn func()) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.reset = fn
	}
}

// WithPreRender sets a callback that runs after reset but before Render().
// Use this to configure mode-specific state (e.g., clock.SetPolicy(...)).
func WithPreRender(fn func()) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.preRender = fn
	}
}

// WithIconGetter is a no-op retained for legacy test compatibility.
//
// Deprecated: modes resolve icons directly via display/widgets/icons.
// Register stub icons with icons.Register / icons.Reset instead.
func WithIconGetter(fn func(name string) (image.Image, bool)) SnapshotOption {
	return func(cfg *snapshotConfig) {}
}

// WithWarnings sets the warnings slice for modes that display hardware notices
// (gpio, stemma).
func WithWarnings(w []string) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.warnings = w
	}
}

// WithFrameCount sets the number of render cycles.
// The framework renders N frames but writes only the final surface state to PNG.
// Default is 1.
func WithFrameCount(n int) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.frameCount = n
	}
}

// WithFrameInterval sets how far the frame clock advances before each render
// pass.
//
// Modes gate animation on elapsed time, so without this the whole render loop
// occupies a single instant and nothing advances. Default is 33ms, roughly the
// 30fps the fast-refresh render loop targets on hardware.
func WithFrameInterval(d time.Duration) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.frameInterval = d
	}
}

// WithReadyWhen supplies a predicate the framework evaluates after each render
// pass, continuing to render until it reports true.
//
// This exists for modes whose displayable state is reached after a duration
// rather than a fixed number of frames. The pager's page-transition path is the
// motivating case: it holds a blank page until its cadence elapses, so the
// number of passes needed depends on layout and policy and is awkward for a
// caller to compute. The frame count becomes a minimum in that case.
//
// The predicate must not mutate mode state; it is only asked whether the frame
// just rendered is worth capturing.
func WithReadyWhen(fn func() bool) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.readyWhen = fn
	}
}

// WithBasename sets the filename prefix for the output PNG.
// Default is the mode ID.
func WithBasename(name string) SnapshotOption {
	return func(cfg *snapshotConfig) {
		cfg.basename = name
	}
}

// applyDefaults creates a snapshotConfig with all options applied and
// sensible defaults for unset fields.
func applyDefaults(opts []SnapshotOption) *snapshotConfig {
	cfg := &snapshotConfig{
		frameCount:    1,
		frameInterval: DefaultFrameInterval,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	// Default basename to modeID if not explicitly set.
	if cfg.basename == "" {
		cfg.basename = cfg.modeID
	}
	return cfg
}

// validate checks the snapshotConfig for invalid or missing configuration.
// It calls t.Fatal on any error condition.
func validate(t *testing.T, cfg *snapshotConfig) {
	t.Helper()

	// Check required options.
	if cfg.width == 0 && cfg.height == 0 {
		t.Fatalf("snapshottest: missing required option: dimensions")
	}
	if cfg.modeID == "" {
		t.Fatalf("snapshottest: missing required option: mode")
	}
	if cfg.outputDir == "" {
		t.Fatalf("snapshottest: missing required option: outputDir")
	}

	// Check dimension validity.
	if cfg.width <= 0 || cfg.height <= 0 || cfg.width > 4096 || cfg.height > 4096 {
		t.Fatalf("snapshottest: invalid dimensions %dx%d", cfg.width, cfg.height)
	}

	if cfg.frameCount < 1 || cfg.frameCount > MaxFrameCount {
		t.Fatalf("snapshottest: frame count %d out of range [1, %d]", cfg.frameCount, MaxFrameCount)
	}

	if cfg.frameInterval < MinFrameInterval || cfg.frameInterval > MaxFrameInterval {
		t.Fatalf("snapshottest: frame interval %v out of range [%v, %v]",
			cfg.frameInterval, MinFrameInterval, MaxFrameInterval)
	}

	// Check mutual exclusivity of WithDisplayCategory and WithColorMode.
	if cfg.displayCategory != nil && cfg.colorModeSet {
		t.Fatalf("snapshottest: WithDisplayCategory and WithColorMode are mutually exclusive")
	}

	// Resolve color mode from display category or explicit color mode.
	resolveColorMode(cfg)

	// Check mode is registered.
	if !displaymodes.IsKnown(cfg.modeID) {
		t.Fatalf("snapshottest: mode %q not found in registry", cfg.modeID)
	}
}

// resolveColorMode derives colorMode and monochrome from the display category
// if set, or derives monochrome from explicit colorMode otherwise.
func resolveColorMode(cfg *snapshotConfig) {
	if cfg.displayCategory != nil {
		switch *cfg.displayCategory {
		case CategoryColor:
			cfg.colorMode = pngpanel.ColorModeFullColor
			cfg.monochrome = false
		case CategoryMono:
			cfg.colorMode = pngpanel.ColorModeMonochrome
			cfg.monochrome = true
		case CategoryEink:
			cfg.colorMode = pngpanel.ColorModeMonochrome
			cfg.monochrome = true
		case CategoryGrayscale:
			cfg.colorMode = pngpanel.ColorModeGrayscale
			cfg.monochrome = false
		default:
			cfg.colorMode = pngpanel.ColorModeFullColor
			cfg.monochrome = false
		}
	} else {
		// Derive monochrome from explicit colorMode.
		cfg.monochrome = cfg.colorMode == pngpanel.ColorModeMonochrome
	}
}
