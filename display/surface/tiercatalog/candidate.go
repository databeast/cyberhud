package tiercatalog

import "github.com/databeast/cyberhud/display/surface/fonts"

// Candidate abstracts a font that can participate in tier selection.
// Bitmap fonts return their fixed metrics; scalable fonts compute metrics
// for a requested pixel height.
type Candidate interface {
	ID() string
	MetricsAt(pixelHeight int) font.Metrics
	IsScalable() bool
}

// bitmapCandidate wraps a font.Face for the existing bitmap fonts.
type bitmapCandidate struct {
	face font.Face
}

func (b bitmapCandidate) ID() string                   { return b.face.ID() }
func (b bitmapCandidate) MetricsAt(_ int) font.Metrics { return b.face.Metrics() }
func (b bitmapCandidate) IsScalable() bool             { return false }
