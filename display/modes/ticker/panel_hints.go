package ticker

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/modes/ticker/styles"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// fitnessNotesPostApply generates fitness notes after policy keys are applied
// in the ticker policy command. It checks if the "style", "font_tier", or
// "accent" keys were among the applied keys and generates appropriate notes
// based on the current panel hints.
func fitnessNotesPostApply(appliedKeys []string) []string {
	styleChanged := false
	fontTierChanged := false
	accentChanged := false
	for _, k := range appliedKeys {
		switch k {
		case "style":
			styleChanged = true
		case "font_tier":
			fontTierChanged = true
		case "accent":
			accentChanged = true
		}
	}
	if !styleChanged && !fontTierChanged && !accentChanged {
		return nil
	}

	hints, ok := getPanelHints()
	if !ok {
		// No panel hints available (e.g., testing or headless mode).
		return nil
	}

	p := PolicySnapshot()
	var notes []string

	// Style fitness notes: evaluate the applied style against the panel.
	if styleChanged {
		s := tickerRegistry.Lookup(p.Style)
		if s != nil {
			notes = append(notes, style.FitnessNotes(s, hints)...)
		}
	}

	// Font tier notes: warn when an explicit large tier is set on a small panel.
	if fontTierChanged && p.FontTier != "auto" {
		tier := styles.ResolveFontTier(p, hints)
		autoTier := styles.ResolveFontTier(source.Policy{FontTier: "auto"}, hints)
		if tierOrdinal(tier) > tierOrdinal(autoTier) {
			notes = append(notes, fmt.Sprintf(
				"note: font_tier %q is larger than auto-selected tier %q for %d×%d panel",
				p.FontTier, string(autoTier), hints.PixelWidth, hints.PixelHeight))
		}
	}

	// Accent notes: warn when accent is set on a mono panel where color has no effect.
	if accentChanged && p.Accent != "none" {
		cap := style.Capability(hints.Capability)
		if cap <= style.MonoFast {
			notes = append(notes, fmt.Sprintf(
				"note: accent %q has no visible effect on %s panel",
				p.Accent, cap))
		}
	}

	return notes
}

// tierOrdinal returns a numeric ordering for tier comparison.
// Larger tiers return higher values.
func tierOrdinal(t tiercatalog.Tier) int {
	switch t {
	case tiercatalog.TierSmall:
		return 0
	case tiercatalog.TierNormal:
		return 1
	case tiercatalog.TierLarge:
		return 2
	case tiercatalog.TierFullsize:
		return 3
	default:
		return 1
	}
}
