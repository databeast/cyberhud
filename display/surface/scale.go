package surface

import (
	"image"

	"github.com/databeast/cyberhud/display/surface/scale"
)

// scaleNearestNeighbor scales src to dstWidth x dstHeight using nearest-neighbor.
// This is a thin wrapper around the exported scale.NearestNeighbor function.
func scaleNearestNeighbor(src image.Image, dstWidth, dstHeight int) *image.RGBA {
	return scale.NearestNeighbor(src, dstWidth, dstHeight)
}
