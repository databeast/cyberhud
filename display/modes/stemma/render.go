package stemma

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView constructs a full style.ViewData for the Stemma mode, using the current policy
// and the provided TextHints for adaptive layout.
// The getIcon parameter resolves icon names to images (typically icons.Get).
// The warnings parameter carries hardware pin notices; when non-empty the first
// warning is prepended to the Hint with a "WARN: " prefix.
//
// Dispatch is performed through the stemmaRegistry: lookup the active style
// by name, fall back to the default style if unregistered, and invoke Build.
func BuildView(devs []*source.Device, hints textlayout.TextHints, getIcon func(name string) (image.Image, bool), warnings []string) style.ViewData {
	pol := GetPolicy()

	// Build the snapshot for Style.Build consumption.
	snap := source.StemmaSnapshot{Devices: devs, GetIcon: getIcon}

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(stemmaRegistry, hints, "stemma", pol.Style)

	// Construct StyleContext boundary and invoke the style's Build method.
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, pol, ctx)

	// Report style resolution metadata to the registry layer.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	// Mode decides if warnings matter for its context.
	if len(warnings) > 0 {
		svd.Items = append([]string{"WARN: " + warnings[0]}, svd.Items...)
	}

	return svd
}
