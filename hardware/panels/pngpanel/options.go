package pngpanel

// Rotation specifies the clockwise rotation angle for output images.
type Rotation int

const (
	Rotation0   Rotation = 0   // No rotation
	Rotation90  Rotation = 90  // 90° clockwise
	Rotation180 Rotation = 180 // 180°
	Rotation270 Rotation = 270 // 270° clockwise
)

// ColorMode specifies the color depth for output encoding.
type ColorMode int

const (
	ColorModeFullColor  ColorMode = iota // RGBA, no conversion
	ColorModeGrayscale                   // 8-bit grayscale via luminance
	ColorModeMonochrome                  // 1-bit B/W via threshold
)

// Option is a functional option for configuring a PNGPanel.
type Option func(*config) error

// config holds the internal construction parameters for a PNGPanel.
type config struct {
	width     int
	height    int
	colorMode ColorMode
	threshold *uint8    // nil = use default (128)
	rotation  *Rotation // nil = use default (0)
	outputDir string
	basename  string // default: "frame"
}

// WithDimensions sets the panel width and height in pixels.
func WithDimensions(width, height int) Option {
	return func(c *config) error {
		c.width = width
		c.height = height
		return nil
	}
}

// WithColorMode sets the color depth mode for output encoding.
func WithColorMode(mode ColorMode) Option {
	return func(c *config) error {
		c.colorMode = mode
		return nil
	}
}

// WithThreshold sets the luminance threshold for monochrome conversion.
// Values are always valid since uint8 is bounded to [0, 255].
func WithThreshold(t uint8) Option {
	return func(c *config) error {
		c.threshold = &t
		return nil
	}
}

// WithOutputDir sets the directory where PNG files are written.
func WithOutputDir(dir string) Option {
	return func(c *config) error {
		c.outputDir = dir
		return nil
	}
}

// WithBasename sets the filename prefix for output PNG files.
// The default basename is "frame", producing files like "frame_0001.png".
func WithBasename(name string) Option {
	return func(c *config) error {
		c.basename = name
		return nil
	}
}

// WithRotation sets the clockwise rotation angle for output images.
// Valid values are 0, 90, 180, and 270.
func WithRotation(degrees int) Option {
	return func(c *config) error {
		r := Rotation(degrees)
		c.rotation = &r
		return nil
	}
}
