package source

import (
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ResolveFace returns a catalog face when available and falls back to the default
// registered face when the catalog is unavailable or cannot resolve the request.
func ResolveFace(hints textlayout.TextHints, family string, tier tiercatalog.Tier) font.Face {
	if face := hints.Face(family, tier); face != nil {
		return face
	}
	return font.Default()
}
