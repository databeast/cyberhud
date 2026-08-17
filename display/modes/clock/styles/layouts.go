package styles

// Shared layout helpers for the clock core layouts in core.go and any
// bespoke Params.BuildFn implementations.

import (
	"image/color"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// selectTier picks the primary font tier for a region of the given pixel height.
//
// This now delegates to style.TierForHeight. The thresholds originated here, and
// were lifted into the style package because "choose a tier from region height" is
// a question every tiered mode must answer, not a clock concern. Keeping a private
// copy would have guaranteed the two drifted apart — the same duplication that
// produced this system's layout-math bugs.
//
// It is kept as a named local function rather than inlining the call so that a
// future clock-specific tier policy has an obvious single place to live: change the
// body here and the whole mode follows.
func selectTier(minHeight int) tiercatalog.Tier {
	return style.TierForHeight(minHeight)
}

// secondaryTier returns the tier one step below the given primary tier.
// Used to render date/weekday rows at a smaller size than the time row.
//
// Delegates to style.SecondaryTier for the same reason as selectTier.
func secondaryTier(primary tiercatalog.Tier) tiercatalog.Tier {
	return style.SecondaryTier(primary)
}

func guardEmptyItems(vd *style.ViewData) {
	if len(vd.Items) == 0 {
		vd.Items = []string{""}
	}
}

func resolveFGColor(name string) color.RGBA {
	return sharedcolor.ResolveAccent(name)
}
