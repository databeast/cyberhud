package pngpanel

import (
	"image"
	"image/color"
	"image/draw"
)

// luminance computes the grayscale luminance of an RGB pixel.
// R, G, B must be in [0, 255]. The result is truncated (floored) to uint8.
func luminance(r, g, b uint8) uint8 {
	return uint8(float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114)
}

// convertToGrayscale converts img to 8-bit grayscale using standard luminance weights.
func convertToGrayscale(img draw.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// RGBA() returns 16-bit pre-multiplied values; convert to 8-bit.
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			gray.SetGray(x, y, color.Gray{Y: luminance(r8, g8, b8)})
		}
	}

	return gray
}

// convertToMonochrome converts img to 1-bit black/white using the given threshold.
// Pixels with luminance >= threshold are white (255); pixels below are black (0).
func convertToMonochrome(img draw.Image, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	mono := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// RGBA() returns 16-bit pre-multiplied values; convert to 8-bit.
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			lum := luminance(r8, g8, b8)

			var val uint8
			if lum >= threshold {
				val = 255
			}
			mono.SetGray(x, y, color.Gray{Y: val})
		}
	}

	return mono
}

// rotateImage rotates img by the given clockwise rotation.
func rotateImage(img image.Image, rot Rotation) image.Image {
	switch src := img.(type) {
	case *image.Gray:
		return rotateGray(src, rot)
	default:
		return rotateGeneric(img, rot)
	}
}

// rotateGray rotates a grayscale image by the given clockwise rotation.
func rotateGray(src *image.Gray, rot Rotation) *image.Gray {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	switch rot {
	case Rotation90:
		dst := image.NewGray(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetGray(h-1-y, x, src.GrayAt(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	case Rotation180:
		dst := image.NewGray(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetGray(w-1-x, h-1-y, src.GrayAt(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	case Rotation270:
		dst := image.NewGray(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetGray(y, w-1-x, src.GrayAt(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	default:
		return src
	}
}

// rotateGeneric rotates any image by the given clockwise rotation.
func rotateGeneric(img image.Image, rot Rotation) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	switch rot {
	case Rotation90:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(h-1-y, x, img.At(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	case Rotation180:
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(w-1-x, h-1-y, img.At(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	case Rotation270:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(y, w-1-x, img.At(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
		return dst
	default:
		return img
	}
}
