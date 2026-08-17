package textbox

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
)

// Alignment specifies horizontal text positioning within bounds.
type Alignment int

const (
	Left Alignment = iota // Default
	Center
	Right
)

// VAlign specifies vertical text positioning within bounds.
type VAlign int

const (
	Top VAlign = iota // Default
	Middle
	Bottom
)

// Overflow specifies how text exceeding horizontal bounds is handled.
type Overflow int

const (
	Truncate Overflow = iota // Append ellipsis (U+2026) when text is too wide
	Wrap                     // Break at whitespace boundaries
	Clip                     // Hard cut at bounds edge, no indicator
)

// Config holds all parameters for TextBox rendering.
type Config struct {
	Bounds        image.Rectangle // Pixel region for rendering (required)
	Text          string          // Text content to render (required)
	Font          font.Face       // Bitmap font (nil → fonts.Default())
	Alignment     Alignment       // Horizontal alignment (default: Left)
	VAlign        VAlign          // Vertical alignment (default: Top)
	Overflow      Overflow        // Overflow strategy (default: Truncate)
	Foreground    color.RGBA      // Text color (zero value → opaque white)
	LineSpacing   int             // Extra vertical pixels between lines (default: 0)
	PadX          int             // Horizontal padding in pixels (default: 0)
	PadY          int             // Vertical padding in pixels (default: 0)
	Border        bool            // Draw 1px border at outer edge of Bounds (default: false)
	FontOverrides []font.Face     // Per-line font overrides (nil entries → Config.Font)
	Label         string          // Sprite label, max 128 chars (default: "textbox")
}

// maxLabelLen is the maximum number of characters allowed in a Label.
const maxLabelLen = 128

// Render produces a TextBox bitmap. Returns nil for invalid bounds or when
// padding/border reduces effective area to zero.
func Render(cfg Config) *widgets.Sprite {
	// Validate bounds: zero or negative dimensions → nil.
	width := cfg.Bounds.Dx()
	height := cfg.Bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil
	}

	// Apply defaults.
	face := cfg.Font
	if face == nil {
		face = font.Default()
	}

	fg, _ := widgets.ResolveColors(cfg.Foreground, color.RGBA{})

	padX := cfg.PadX
	if padX < 0 {
		padX = 0
	}

	padY := cfg.PadY
	if padY < 0 {
		padY = 0
	}

	lineSpacing := cfg.LineSpacing
	if lineSpacing < 0 {
		lineSpacing = 0
	}

	// Compute border inset.
	borderInset := 0
	if cfg.Border {
		borderInset = 1
	}

	// Compute effective area.
	effectiveWidth := width - 2*padX - 2*borderInset
	effectiveHeight := height - 2*padY - 2*borderInset

	if effectiveWidth <= 0 || effectiveHeight <= 0 {
		return nil
	}

	// Create the output image (fully transparent by default).
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw 1px border at outer edges if requested.
	if cfg.Border {
		widgets.DrawBorder(img, width, height, fg)
	}

	// Render text if non-empty.
	if cfg.Text != "" {
		textOriginX := padX + borderInset
		textOriginY := padY + borderInset

		// Compute layout.
		lines := computeLayout(cfg.Text, effectiveWidth, effectiveHeight, face, cfg.Overflow, cfg.VAlign, lineSpacing, cfg.FontOverrides)

		// Render each visual line.
		for _, line := range lines {
			renderLine(img, line.face, line.text, textOriginX, textOriginY+line.y, effectiveWidth, cfg.Alignment, cfg.Overflow, fg)
		}
	}

	// Resolve label.
	label := cfg.Label
	if label == "" {
		label = "textbox"
	}
	if len([]rune(label)) > maxLabelLen {
		label = string([]rune(label)[:maxLabelLen])
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    label,
	}
}
