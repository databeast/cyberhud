package dashboard

import (
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// fitnessNotesPostApply generates fitness notes for the applied style
// against the current panel hints when the "style" key changes.
var fitnessNotesPostApply = dashboardRegistry.FitnessPostApply(modehints.Current, func() string { return GetPolicy().Style })

// BuildView returns the full dashboard view result using registry-based dispatch.
// It resolves system state internally via buildSnapshot, normalizes the policy,
// dispatches to the active style, and injects Title and Hint chrome.
func BuildView(hints textlayout.TextHints) style.ViewData {
	pol := normalizePolicy(GetPolicy())

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(dashboardRegistry, hints, "dashboard", pol.Style)

	// Build the snapshot from OS queries.
	data := source.BuildDashboardContent()

	// Construct StyleContext for the style boundary.
	ctx := style.NewStyleContext(hints)

	// validate it's even possible to render anything on the attached surface region
	if ctx.Layout(0).AvailableContentWidth() == 0 || ctx.Layout(0).AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	// Call the style's Build method with the snapshot and context.
	svd := s.Build(data, policy, ctx)

	// Report style resolution metadata to the registry layer.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	return svd
}

func RenderCacheKey() uint32 {
	dat := source.BuildDashboardContent()
	pol := GetPolicy()

	return region.CalcRegionCacheKey(
		dat.Hostname, dat.IPAddress, dat.Uptime,
		pol.Style, pol.ColorAccent,
		dat.WifiSSID, dat.Version,
	)
}
