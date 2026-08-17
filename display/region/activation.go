package region

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/hardware/input"
)

// PanelActivationConfig holds all inputs needed to activate a panel's region
// infrastructure via [ActivatePanel]. It is part of the ActivatePanel orchestration
// layer and groups the screen topology, layout strategy, mode configuration, and
// input event source into a single value passed to the activation function.
//
// Fields:
//
//   - Screens: the physical screen positions that define the panel topology.
//     At least one entry is required; ActivatePanel returns an error if empty.
//   - Layout: an optional explicit [RegionLayout] describing how regions map to
//     screen areas. When nil or empty, ActivatePanel generates a default layout
//     from the screen configuration via [GenerateDefaultLayout].
//   - DefaultMode: the mode ID applied to regions that do not have a per-screen
//     override in ScreenModes.
//   - InputEnabled: whether the panel should process hardware input events
//     dispatched through the RenderLoop.
//   - AvailModes: the list of mode IDs that are valid for this panel. Used by
//     [GenerateDefaultLayout] to populate layout defaults.
//   - ScreenModes: per-screen default mode overrides keyed by screen Name.
//     Screens not present in this map fall back to DefaultMode.
//   - ModeValidator: an optional validation function that returns true when a
//     mode ID is recognized. When non-nil, it is assigned via
//     [RegionManager.SetModeValidator] for use during SetMode calls.
//     This field remains exported because cmd/cyberhudd sets it directly
//     when building PanelActivationConfig (see display_runtime.go and tests).
//   - Events: a receive-only channel of hardware input events forwarded to the
//     RenderLoop's input dispatcher on each tick.
type PanelActivationConfig struct {
	Screens       []ScreenPosition
	Layout        *RegionLayout
	DefaultMode   string
	InputEnabled  bool
	AvailModes    []string
	ScreenModes   map[string]string
	ModeValidator func(string) bool
	Events        <-chan input.Event

	// PanelProduct is the normalized panel product name (e.g., "waveshare-1.3hat").
	// Propagated to each Region's TextHints so display modes can resolve style aliases.
	// Empty string if unknown or whitespace-only.
	PanelProduct string

	// PanelPPI is the panel-level pixels-per-inch. Propagated to each Region's
	// TextHints via the cascading resolution: Screen.PPI → PanelPPI → 96.0 default.
	// Zero means "undeclared at panel level."
	PanelPPI float64

	// ConfigPPI is the user-supplied PPI override from the runtime config file.
	// Sits between PanelPPI and the 96.0 default in the cascade.
	// Zero means "no config override."
	ConfigPPI float64
}

// PanelActivation is the ActivatePanel-layer result type that holds the
// initialized region infrastructure for an active Panel. It bundles all
// components constructed by [ActivatePanel] so the caller has a single value
// from which to wire post-activation dependencies.
type PanelActivation struct {
	VirtualDisplay *VirtualDisplay // unified framebuffer backing all screens
	RegionManager  *RegionManager  // region lifecycle owner
	FlushPath      *FlushPath      // per-screen hardware output path
	RenderLoop     *RenderLoop     // tick-driven render coordinator (not yet started)
	ModeSwitch     *ModeSwitch     // command handler for runtime mode changes
}

// ActivatePanel is the top-level orchestrator for the region infrastructure. It
// constructs a [VirtualDisplay], a [RegionManager], a [FlushPath], a [RenderLoop],
// and a [ModeSwitch] for the given panel configuration and returns them bundled in
// a [PanelActivation].
//
// The config parameter supplies the screen topology, layout strategy, mode
// defaults, and input event source (see [PanelActivationConfig] for field
// details).
//
// The returned [PanelActivation] contains pointers to the fully-wired
// components. The RenderLoop is created but not started — the caller is
// responsible for the following post-activation steps:
//
//  1. Wire a [ModeFactory] on each [Region] via [Region.SetModeFactory] so that
//     mode changes can construct new mode instances.
//  2. Configure the [RenderLoop] via [RenderLoop.Apply], supplying a [Renderer]
//     (via [WithRenderer]), an [InputDispatcher] (via [WithInputDispatcher]),
//     and a [TickRateResolver] (via [WithTickRateResolver]).
//  3. Call [RenderLoop.Run] to start the render/input loop. Run blocks until
//     [RenderLoop.Stop] is called.
//
// ActivatePanel returns an error if:
//   - config.Screens is empty (at least one screen is required)
//   - [VirtualDisplay] construction fails
//   - Layout resolution or allocation fails
func ActivatePanel(config PanelActivationConfig) (*PanelActivation, error) {
	if len(config.Screens) == 0 {
		return nil, fmt.Errorf("region: at least one screen is required for panel activation")
	}

	// Step 1: Construct VirtualDisplay first because it provides the unified
	// framebuffer that all downstream components (RegionManager, FlushPath,
	// RenderLoop) reference. If this is deferred, SubImage-based region surfaces
	// cannot be allocated and FlushPath has no source buffer to extract from.
	vd, err := NewVirtualDisplayFromScreens(config.Screens)
	if err != nil {
		return nil, fmt.Errorf("region: virtual display construction failed: %w", err)
	}

	// Step 2: Create RegionManager immediately after VirtualDisplay because it
	// wraps the VD to provide region lifecycle operations (Allocate, SetMode).
	// Screen positions are passed here so that allocated regions receive correct
	// TextHints reflecting physical screen geometry. If created before VD, the
	// manager would hold a nil framebuffer and all SubImage calls would panic.
	rm := NewRegionManagerWithScreens(vd, config.Screens)
	rm.panelProduct = strings.TrimSpace(config.PanelProduct)
	rm.panelPPI = config.PanelPPI
	rm.configPPI = config.ConfigPPI
	if config.ModeValidator != nil {
		rm.SetModeValidator(config.ModeValidator)
	}

	// Step 3: Resolve layout and allocate regions before creating FlushPath or
	// RenderLoop. Allocation carves zero-copy SubImage surfaces from the VD for
	// each region. FlushPath needs the VD populated with region pixel data, and
	// RenderLoop iterates allocated regions for per-frame rendering. If allocation
	// is deferred, FlushPath flushes an empty framebuffer and RenderLoop has no
	// regions to render.
	layout, err := resolveLayout(vd, config)
	if err != nil {
		return nil, fmt.Errorf("region: layout resolution failed: %w", err)
	}

	if err := rm.AllocateLayout(layout); err != nil {
		return nil, fmt.Errorf("region: layout allocation failed: %w", err)
	}

	// Step 4: Create FlushPath after VD and allocation because it reads
	// per-screen rectangles from the VD framebuffer. It depends on the VD being
	// fully constructed and regions allocated so that flushed pixels reflect
	// actual rendered content. If created before allocation, the screen rectangles
	// would contain only zeroed memory on first flush.
	fp := NewFlushPath(vd, config.Screens)

	// Step 5: Create ModeSwitch after RegionManager because it routes mode-change
	// commands through rm.SetMode. If created before the manager exists, the
	// ModeSwitch would hold a nil RegionManager reference and panic on Execute.
	ms := NewModeSwitch(rm)

	// Step 6: Create RenderLoop last because it depends on RegionManager (to
	// iterate regions), FlushPath (to push rendered frames to hardware), and the
	// Events channel. If created earlier, it would capture incomplete or nil
	// dependencies and either panic or silently skip rendering on first tick.
	rl := NewRenderLoop(rm, fp, config.Events)

	return &PanelActivation{
		VirtualDisplay: vd,
		RegionManager:  rm,
		FlushPath:      fp,
		RenderLoop:     rl,
		ModeSwitch:     ms,
	}, nil
}

// resolveLayout determines the RegionLayout to use for allocation.
// If an explicit layout is provided with specs, it is returned directly.
// Otherwise, GenerateDefaultLayout is called to produce the default layout.
func resolveLayout(vd *VirtualDisplay, config PanelActivationConfig) (RegionLayout, error) {
	// If an explicit layout is provided and has specs, use it.
	if config.Layout != nil && len(config.Layout.Specs) > 0 {
		return *config.Layout, nil
	}

	// Otherwise, generate the default layout.
	return GenerateDefaultLayout(vd, config)
}
