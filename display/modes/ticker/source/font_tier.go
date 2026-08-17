package source

import (
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ResolveFontTier returns the effective tier for ticker text rendering.
// Auto mode adapts to panel height; explicit tiers are passed through.
func ResolveFontTier(policy Policy, hints textlayout.TextHints) tiercatalog.Tier {
	if policy.FontTier == "auto" || policy.FontTier == "" {
		if hints.PixelHeight >= 200 {
			return tiercatalog.TierLarge
		}
		return tiercatalog.TierNormal
	}
	switch policy.FontTier {
	case "small":
		return tiercatalog.TierSmall
	case "normal":
		return tiercatalog.TierNormal
	case "large":
		return tiercatalog.TierLarge
	case "fullsize":
		return tiercatalog.TierFullsize
	default:
		return tiercatalog.TierNormal
	}
}
