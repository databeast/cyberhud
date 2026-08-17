package styles

import (
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// resolveFontTier determines the effective tier based on policy and panel dimensions.
// When FontTier is "auto" or empty, it uses panel height to select an appropriate tier.
// An explicit tier setting bypasses the auto logic entirely.
func ResolveFontTier(policy source.Policy, hints textlayout.TextHints) tiercatalog.Tier {
	return source.ResolveFontTier(policy, hints)
}
