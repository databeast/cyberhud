package driver

// MonoThreshold is the BT.601 luminance midpoint used by 1-bit panels to decide
// whether a pixel is "on". Values are in the same 16-bit-ish scale produced by
// Luminance (0..0xFFFF-ish), so the comparison is Luminance(...) vs MonoThreshold.
const MonoThreshold uint32 = 0x3000

// Luminance returns the BT.601 perceptual luminance of an RGB triple using the
// integer coefficients (299, 587, 114)/1000. Inputs are the 32-bit channel
// values as returned by image/color's RGBA(). Callers compare the result
// against MonoThreshold to drive 1-bit (black/white) panels; note that e-paper
// drivers invert the comparison because dark ink is the "on" state.
func Luminance(r, g, b uint32) uint32 {
	return (299*r + 587*g + 114*b) / 1000
}
