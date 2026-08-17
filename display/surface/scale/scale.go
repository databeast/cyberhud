package scale

import "image"

// NearestNeighbor scales src to dstWidth × dstHeight using nearest-neighbor
// interpolation. Returns nil if src is nil or dst dimensions are non-positive.
func NearestNeighbor(src image.Image, dstWidth, dstHeight int) *image.RGBA {
	if src == nil {
		return nil
	}
	if dstWidth <= 0 || dstHeight <= 0 {
		return nil
	}

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	if srcW <= 0 || srcH <= 0 {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	for y := 0; y < dstHeight; y++ {
		srcY := srcBounds.Min.Y + y*srcH/dstHeight
		for x := 0; x < dstWidth; x++ {
			srcX := srcBounds.Min.X + x*srcW/dstWidth
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}
