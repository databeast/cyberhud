package region

import (
	"fmt"
	"image"
	"sort"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

// ScreenPosition is a VirtualDisplay-layer type that locates a Physical Screen
// within the Virtual Display coordinate space. It maps a hardware display target
// to a rectangular region of the framebuffer and optionally provides text layout
// hints derived from the hardware driver's capabilities.
type ScreenPosition struct {
	Index        int                         // zero-based ordinal used for ordering and lookup
	Name         string                      // human-readable identifier (e.g. "left", "right")
	Bounds       image.Rectangle             // position and size within Virtual Display coordinates (logical dimensions)
	Target       driver.DrawTarget           // hardware draw target that receives flushed frames
	HintProvider func() textlayout.TextHints // returns driver-specific text hints; nil when unavailable
	Rotation     int                         // clockwise rotation (0, 90, 180, 270) applied in FlushPath before DrawImage
	MirrorX      bool                        // horizontal mirror applied in FlushPath before DrawImage
	PPI          float64                     // optional per-screen PPI; zero means "use panel default"
}

// RegionSpec is a RegionManager-layer type that defines one Region for allocation.
// It carries the region's identity, spatial placement, and initial display mode.
type RegionSpec struct {
	Name        string          // unique identifier (1-64 characters) used for lookup
	Bounds      image.Rectangle // rectangular area within the Virtual Display coordinate space
	DefaultMode string          // display mode ID applied at allocation time; empty means no initial mode
}

// RegionLayout is a RegionManager-layer type that holds an ordered list of
// [RegionSpec] values describing the spatial arrangement of all Regions on a
// Panel. The layout is passed to [RegionManager.AllocateLayout] for atomic
// validation and allocation.
type RegionLayout struct {
	Specs []RegionSpec // ordered region specifications; allocation proceeds in slice order
}

// GenerateDefaultLayout creates the default [RegionLayout] for a Panel that has
// no explicit layout configured. It is part of the RegionManager layer and is
// called by [ActivatePanel] when [PanelActivationConfig].Layout is nil or empty.
//
// The vd parameter is the VirtualDisplay whose bounds determine region sizing.
// The config parameter provides screens, available modes, and per-screen mode
// overrides used to populate the layout.
//
// For a single-screen panel, one "default" region covering the full VirtualDisplay
// is created. For multi-screen panels, one region per screen is created, positioned
// left-to-right by ascending screen index.
//
// Returns the generated RegionLayout and nil error on success, or an empty layout
// and a non-nil error if no screens are present or no valid default mode can be
// resolved.
func GenerateDefaultLayout(vd *VirtualDisplay, config PanelActivationConfig) (RegionLayout, error) {
	if len(config.Screens) == 0 {
		return RegionLayout{}, fmt.Errorf("region: at least one screen is required")
	}

	if len(config.Screens) == 1 {
		return generateSingleScreenLayout(vd, config)
	}

	return generateMultiScreenLayout(vd, config)
}

// generateSingleScreenLayout creates a single "default" region covering the full VD bounds.
func generateSingleScreenLayout(vd *VirtualDisplay, config PanelActivationConfig) (RegionLayout, error) {
	mode, err := resolveDefaultMode(config)
	if err != nil {
		return RegionLayout{}, err
	}

	spec := RegionSpec{
		Name:        "default",
		Bounds:      vd.Bounds(),
		DefaultMode: mode,
	}

	return RegionLayout{Specs: []RegionSpec{spec}}, nil
}

// generateMultiScreenLayout creates one region per Physical Screen, using each
// screen's pre-computed Bounds from buildPositions. Each screen must have a valid
// DefaultMode in config.ScreenModes.
func generateMultiScreenLayout(vd *VirtualDisplay, config PanelActivationConfig) (RegionLayout, error) {
	// Sort screens by ascending index for deterministic ordering.
	screens := make([]ScreenPosition, len(config.Screens))
	copy(screens, config.Screens)
	sort.Slice(screens, func(i, j int) bool {
		return screens[i].Index < screens[j].Index
	})

	specs := make([]RegionSpec, 0, len(screens))
	for _, screen := range screens {
		// Look up the screen's default mode.
		screenMode := ""
		if config.ScreenModes != nil {
			screenMode = config.ScreenModes[screen.Name]
		}

		// Validate: must be non-empty and registered in AvailModes.
		if screenMode == "" || !modeInList(screenMode, config.AvailModes) {
			return RegionLayout{}, fmt.Errorf("region: screen %q (index %d) has no valid default mode", screen.Name, screen.Index)
		}

		// Use the screen's pre-computed Bounds directly. These come from
		// buildPositions which handles both explicit XPosition/YPosition and
		// auto left-to-right accumulation. Using these bounds ensures the
		// region writes to the same VD location that FlushPath reads from.
		spec := RegionSpec{
			Name:        screen.Name,
			Bounds:      screen.Bounds,
			DefaultMode: screenMode,
		}
		specs = append(specs, spec)
	}

	return RegionLayout{Specs: specs}, nil
}

// resolveDefaultMode resolves the default mode for a single-screen panel.
// Resolution order:
//  1. config.DefaultMode if non-empty AND in config.AvailModes → use it
//  2. If config.InputEnabled AND "menu" is in config.AvailModes → use "menu"
//  3. If "dashboard" is in config.AvailModes → use "dashboard"
//  4. If config.AvailModes has at least one entry → use first
//  5. Return error if no modes available
func resolveDefaultMode(config PanelActivationConfig) (string, error) {
	// 1. Explicit Panel default mode if available.
	if config.DefaultMode != "" && modeInList(config.DefaultMode, config.AvailModes) {
		return config.DefaultMode, nil
	}

	// 2. "menu" if input is enabled and available.
	if config.InputEnabled && modeInList("menu", config.AvailModes) {
		return "menu", nil
	}

	// 3. "dashboard" if available.
	if modeInList("dashboard", config.AvailModes) {
		return "dashboard", nil
	}

	// 4. First available mode.
	if len(config.AvailModes) > 0 {
		return config.AvailModes[0], nil
	}

	// 5. No modes available.
	return "", fmt.Errorf("region: no display modes available for default region")
}

// modeInList checks if a mode ID exists in the provided list.
func modeInList(mode string, modes []string) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}
