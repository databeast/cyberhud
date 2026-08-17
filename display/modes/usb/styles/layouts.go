package styles

// Shared layout helpers for the USB core layouts in core.go and any
// bespoke Params.BuildFn implementations.

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// selectTier picks the primary font tier for a region of the given pixel height.
//
// Delegates to style.TierForHeight so that the USB mode inherits centralized tier
// thresholds without maintaining its own copy.
func selectTier(pixelHeight int) tiercatalog.Tier {
	return style.TierForHeight(pixelHeight)
}

// secondaryTier returns the tier one step below the given primary tier.
// Used to render detail rows at a smaller size than the primary heading.
//
// Delegates to style.SecondaryTier for the same reason as selectTier.
func secondaryTier(primary tiercatalog.Tier) tiercatalog.Tier {
	return style.SecondaryTier(primary)
}

// guardEmptyItems ensures vd.Items is never empty. If Items has length 0,
// it is replaced with a single-element slice containing one empty string.
// A non-empty Items slice is left unchanged.
func guardEmptyItems(vd *style.ViewData) {
	if len(vd.Items) == 0 {
		vd.Items = []string{""}
	}
}
