package pngpanel

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
)

// Compile-time interface assertions.
var _ driver.DrawTarget = (*PNGPanel)(nil)
var _ textlayout.TextHintProvider = (*PNGPanel)(nil)

// PNGPanel implements driver.DrawTarget and textlayout.TextHintProvider,
// writing rendered frames as PNG files to disk.
type PNGPanel struct {
	width     int
	height    int
	colorMode ColorMode
	threshold uint8
	rotation  Rotation
	outputDir string
	basename  string
	counter   uint64
	bounds    image.Rectangle
}

// New creates a PNGPanel with the given options.
// Returns an error if:
//   - width or height is outside [1, 4096]
//   - outputDir is empty
//   - an unsupported ColorMode is provided
func New(opts ...Option) (*PNGPanel, error) {
	cfg := &config{}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	// Validate dimensions.
	if cfg.width < 1 || cfg.width > 4096 || cfg.height < 1 || cfg.height > 4096 {
		return nil, fmt.Errorf("pngpanel: invalid dimensions: width and height must be in [1, 4096]")
	}

	// Validate output directory.
	if cfg.outputDir == "" {
		return nil, fmt.Errorf("pngpanel: output path must not be empty")
	}

	// Validate color mode.
	switch cfg.colorMode {
	case ColorModeFullColor, ColorModeGrayscale, ColorModeMonochrome:
		// valid
	default:
		return nil, fmt.Errorf("pngpanel: unsupported color mode %d (supported: full-color, monochrome, grayscale)", int(cfg.colorMode))
	}

	// Apply defaults.
	threshold := uint8(128)
	if cfg.threshold != nil {
		threshold = *cfg.threshold
	}

	// Validate rotation.
	rotation := Rotation0
	if cfg.rotation != nil {
		switch *cfg.rotation {
		case Rotation0, Rotation90, Rotation180, Rotation270:
			rotation = *cfg.rotation
		default:
			return nil, fmt.Errorf("pngpanel: invalid rotation %d (valid values: 0, 90, 180, 270)", int(*cfg.rotation))
		}
	}

	basename := "frame"
	if cfg.basename != "" {
		basename = cfg.basename
	}

	return &PNGPanel{
		width:     cfg.width,
		height:    cfg.height,
		colorMode: cfg.colorMode,
		threshold: threshold,
		rotation:  rotation,
		outputDir: cfg.outputDir,
		basename:  basename,
		counter:   0,
		bounds:    image.Rect(0, 0, cfg.width, cfg.height),
	}, nil
}

// Bounds returns the configured rectangle (origin 0,0).
// Implements driver.DrawTarget.
func (p *PNGPanel) Bounds() image.Rectangle {
	return p.bounds
}

// DrawImage converts the frame per color mode and writes it as a PNG file.
// Implements driver.DrawTarget.
func (p *PNGPanel) DrawImage(img draw.Image) error {
	// Validate nil image.
	if img == nil {
		return fmt.Errorf("pngpanel: image must not be nil")
	}

	// Validate dimension mismatch.
	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()
	if imgW != p.width || imgH != p.height {
		return fmt.Errorf("pngpanel: dimension mismatch: expected %dx%d, got %dx%d", p.width, p.height, imgW, imgH)
	}

	// Apply color conversion pipeline.
	var output image.Image
	switch p.colorMode {
	case ColorModeGrayscale:
		output = convertToGrayscale(img)
	case ColorModeMonochrome:
		output = convertToMonochrome(img, p.threshold)
	default:
		output = img
	}

	// Apply rotation if configured.
	if p.rotation != Rotation0 {
		output = rotateImage(output, p.rotation)
	}

	// Increment frame counter before writing (save old for rollback on error).
	oldCounter := p.counter
	p.counter++

	// Generate filename with appropriate zero-padding.
	var filename string
	if p.counter <= 9999 {
		filename = fmt.Sprintf("%s_%04d.png", p.basename, p.counter)
	} else {
		filename = fmt.Sprintf("%s_%d.png", p.basename, p.counter)
	}

	// Create output directory (including intermediate dirs) if it does not exist.
	if err := os.MkdirAll(p.outputDir, 0o755); err != nil {
		p.counter = oldCounter
		return fmt.Errorf("pngpanel: cannot create output directory %q: %w", p.outputDir, err)
	}

	// Atomic write: temp file → encode PNG → close → rename to final path.
	finalPath := filepath.Join(p.outputDir, filename)
	tmpPath := finalPath + ".tmp"

	f, err := os.Create(tmpPath)
	if err != nil {
		p.counter = oldCounter
		return fmt.Errorf("pngpanel: write failed for %q: %w", finalPath, err)
	}

	if err := png.Encode(f, output); err != nil {
		f.Close()
		os.Remove(tmpPath)
		p.counter = oldCounter
		return fmt.Errorf("pngpanel: write failed for %q: %w", finalPath, err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		p.counter = oldCounter
		return fmt.Errorf("pngpanel: write failed for %q: %w", finalPath, err)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		p.counter = oldCounter
		return fmt.Errorf("pngpanel: rename failed for %q: %w", finalPath, err)
	}

	return nil
}

// TextHints returns layout hints based on configured dimensions and color mode.
// Implements textlayout.TextHintProvider.
//
// Glyph metrics are intentionally absent. A panel does not own a font; the Region
// derives glyph metrics from the tier catalog it builds for the region's real
// dimensions (see region.applyBaselineGlyphMetrics). This panel previously reported
// the textlayout 5x8/6/10 constants, which happened to match the surface's default
// face and so masked the fact that nothing was choosing a font deliberately.
func (p *PNGPanel) TextHints() textlayout.TextHints {
	h := textlayout.TextHints{
		PixelWidth:      p.width,
		PixelHeight:     p.height,
		DefaultLineMode: textlayout.LineModeTruncate,
	}
	if p.colorMode == ColorModeMonochrome {
		h.PreferEventRefresh = true
		h.Capability = textlayout.CapMonoSlow
		h.DefaultTickerDirection = textlayout.TickerDirectionNone
	} else if p.colorMode == ColorModeGrayscale {
		h.SupportsVerticalScroll = true
		h.SupportsHorizontalScroll = true
		h.SupportsAutoScroll = true
		h.Capability = textlayout.CapGrayscaleFast
		h.DefaultTickerDirection = textlayout.TickerDirectionVertical
	} else {
		h.SupportsVerticalScroll = true
		h.SupportsHorizontalScroll = true
		h.SupportsAutoScroll = true
		h.Capability = textlayout.CapColorFast
		h.DefaultTickerDirection = textlayout.TickerDirectionVertical
	}
	return h
}

// ResetCounter sets the frame counter back to 0.
func (p *PNGPanel) ResetCounter() {
	p.counter = 0
}
