package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
)

// Snapshot is the mode-specific data consumed by style Build methods.
// In the template it is intentionally empty — real modes populate it
// with domain data (time, GPIO state, device list, etc.).
type Snapshot struct{}

// ---------------------------------------------------------------------------
// Compile-time interface compliance checks
// ---------------------------------------------------------------------------

var _ style.Style[Snapshot, Policy] = MonoSmall128x32Style{}
var _ style.Style[Snapshot, Policy] = MonoSmall128x64Style{}
var _ style.Style[Snapshot, Policy] = ColorSmall160x80Style{}
var _ style.Style[Snapshot, Policy] = ColorSmall160x128Style{}
var _ style.Style[Snapshot, Policy] = ColorMedium240x135Style{}
var _ style.Style[Snapshot, Policy] = ColorMedium240x240Style{}
var _ style.Style[Snapshot, Policy] = ColorMedium320x240Style{}
var _ style.Style[Snapshot, Policy] = ColorLarge480x320Style{}
var _ style.Style[Snapshot, Policy] = ColorLarge800x480Style{}
var _ style.Style[Snapshot, Policy] = EinkSmall122x250Style{}
var _ style.Style[Snapshot, Policy] = EinkSmall176x264Style{}
var _ style.Style[Snapshot, Policy] = EinkSmall200x200Style{}
var _ style.Style[Snapshot, Policy] = EinkSmall212x104Style{}
var _ style.Style[Snapshot, Policy] = EinkMedium296x128Style{}
var _ style.Style[Snapshot, Policy] = EinkMedium400x300Style{}
var _ style.Style[Snapshot, Policy] = EinkLarge480x800Style{}
var _ style.Style[Snapshot, Policy] = EinkLarge800x480Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast160x80Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast160x128Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast240x135Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast240x240Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast320x240Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast480x320Style{}
var _ style.Style[Snapshot, Policy] = GrayscaleFast800x480Style{}
var _ style.Style[Snapshot, Policy] = GraphicsExampleStyle{}

// ---------------------------------------------------------------------------
// Style registry
// ---------------------------------------------------------------------------

// templateRegistry holds all 25 styles in registration order: 24 per-resolution
// text-only skeletons plus one graphics-example style (last). The first style
// (mono-128x32) is the default selected by fitness when none is configured.
//
// Each per-resolution style lives in its own file (style_<category>_<WxH>.go)
// so developers can grab the one matching their target panel as a starting
// point. The graphics-example style (style_graphics_example.go) is the
// reference for drawing text alongside graphics correctly.
var templateRegistry = style.NewRegistry[Snapshot, Policy](
	MonoSmall128x32Style{}, // default (first)
	MonoSmall128x64Style{},
	ColorSmall160x80Style{},
	ColorSmall160x128Style{},
	ColorMedium240x135Style{},
	ColorMedium240x240Style{},
	ColorMedium320x240Style{},
	ColorLarge480x320Style{},
	ColorLarge800x480Style{},
	EinkSmall122x250Style{},
	EinkSmall176x264Style{},
	EinkSmall200x200Style{},
	EinkSmall212x104Style{},
	EinkMedium296x128Style{},
	EinkMedium400x300Style{},
	EinkLarge480x800Style{},
	EinkLarge800x480Style{},
	// GrayscaleFast variants for color-capable resolutions
	GrayscaleFast160x80Style{},
	GrayscaleFast160x128Style{},
	GrayscaleFast240x135Style{},
	GrayscaleFast240x240Style{},
	GrayscaleFast320x240Style{},
	GrayscaleFast480x320Style{},
	GrayscaleFast800x480Style{},
	// Worked example of mixing text + graphics in one view (pixel-anchored
	// label + bar sprites). Registered LAST and selected by explicit name
	// (style=graphics-example) so it never displaces the text-only styles by
	// fitness; see style_graphics_example.go for the two-plane layout lesson.
	GraphicsExampleStyle{},
)
