package styles

// Serial style declarations for 480x800 panels (large portrait). Covers both
// slow-refresh e-paper and fast-refresh variants.

var (
	MonoSlow480x800Style      = monoSlowStyle(480, 800)
	GrayscaleSlow480x800Style = grayscaleSlowStyle(480, 800)
	ColorSlow480x800Style     = colorSlowStyle(480, 800)

	GrayscaleFast480x800Style = grayscaleFastStyle(480, 800)
	ColorFast480x800Style     = colorFastStyle(480, 800)
)
