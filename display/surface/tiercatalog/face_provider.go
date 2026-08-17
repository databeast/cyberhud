// Package tiercatalog resolves tier intents into catalog-validated font.Face objects.
// It is the ONLY way sprite-producing code can obtain font faces for widget rendering.
// External callers should use TextHints.Face() or TextHints.AllFaces() rather than
// constructing a faceProvider directly.
package tiercatalog

import (
	"github.com/databeast/cyberhud/display/surface/fonts"
)

// faceResolver is a function that resolves a (family, tier, catalog) request into
// a concrete font.Face. This indirection breaks the import cycle between tiercatalog
// and tierselect — callers inject tierselect.Select at construction time.
type faceResolver func(catalog Catalog, family string, tier Tier) font.Face

// faceProvider resolves tier intents into catalog-validated font faces.
// It is backed by a Catalog and a faceResolver.
type faceProvider struct {
	catalog  Catalog
	family   string
	resolver faceResolver
}

// newFaceProvider creates a faceProvider backed by the given catalog and resolver.
// Mode/style code should use TextHints.Face() instead of constructing a faceProvider
// directly. The family parameter sets the default font family preference (e.g., "spleen",
// "terminus"). The resolver should be tierselect.Select (or equivalent) — it is injected
// to avoid an import cycle between tiercatalog and tierselect.
// All returned faces are guaranteed to satisfy GlyphAdvance ≤ maxAdvance for the catalog's region.
func newFaceProvider(catalog Catalog, family string, resolver faceResolver) *faceProvider {
	return &faceProvider{
		catalog:  catalog,
		family:   family,
		resolver: resolver,
	}
}

// Face resolves the given tier to a catalog-validated font.Face using the
// provider's default family preference.
func (fp *faceProvider) Face(tier Tier) font.Face {
	return fp.resolver(fp.catalog, fp.family, tier)
}

// FaceForFamily resolves the given family and tier to a catalog-validated font.Face.
func (fp *faceProvider) FaceForFamily(family string, tier Tier) font.Face {
	return fp.resolver(fp.catalog, family, tier)
}

// AllFaces returns one catalog-validated face per defined tier, resolved using
// the provider's default family preference. The faces are returned in tier
// order: Small, Normal, Large, Huge, Colossal, Full.
func (fp *faceProvider) AllFaces() []font.Face {
	tiers := fp.catalog.Tiers()
	faces := make([]font.Face, len(tiers))
	for i, tier := range tiers {
		faces[i] = fp.resolver(fp.catalog, fp.family, tier)
	}
	return faces
}
