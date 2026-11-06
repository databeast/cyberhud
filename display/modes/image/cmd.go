package image

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the uniform CmdHandler for the image mode's policy options.
// Policy edits flow through it so that applying key=value pairs fires
// cmdutil.OnApplied, which drives the centralized PolicyStore capture (enabling
// `freeze policy` persistence) — the same path every conformant mode uses.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "image",
	Keys: []cmdutil.KeyDef{
		{Name: "fit", Validate: cmdutil.AllowedValidator([]string{FitTruncate, FitScale, FitStretch})},
		{Name: "style", Validate: cmdutil.AllowedValidator(allowedStyles)},
	},
	Get: func(key string) string {
		p := PolicySnapshot()
		switch key {
		case "fit":
			return p.Fit
		case "style":
			return p.Style
		default:
			return ""
		}
	},
	Apply: func(key, value string) {
		p := PolicySnapshot()
		switch key {
		case "fit":
			p.Fit = strings.ToLower(strings.TrimSpace(value))
		case "style":
			p.Style = strings.ToLower(strings.TrimSpace(value))
		}
		SetPolicy(p)
	},
}
