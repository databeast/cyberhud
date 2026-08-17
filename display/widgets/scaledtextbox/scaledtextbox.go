package scaledtextbox

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/scale"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/textbox"
)

// Config holds parameters for ScaledTextBox rendering.
// Text is rendered at LogicalSize resolution, then scaled to TargetSize.
type Config struct {
	// LogicalSize is the resolution at which text is rendered internally.
	// The TextBox layout engine operates at this size.
	LogicalSize image.Point

	// TargetSize is the final Sprite output size after scaling.
	TargetSize image.Point

	// Position is the top-left coordinate for the output Sprite.
	Position image.Point

	// Text content to render.
	Text string

	// Font is the bitmap font for rendering. Nil falls back to fonts.Default().
	Font font.Face

	// Alignment specifies horizontal text positioning within bounds.
	Alignment textbox.Alignment

	// VAlign specifies vertical text positioning within bounds.
	VAlign textbox.VAlign

	// Overflow specifies how text exceeding horizontal bounds is handled.
	Overflow textbox.Overflow

	// Foreground is the text color. Zero value defaults to opaque white.
	Foreground color.RGBA

	// LineSpacing is extra vertical pixels between lines.
	LineSpacing int

	// PadX is horizontal padding in pixels.
	PadX int

	// PadY is vertical padding in pixels.
	PadY int

	// Border draws a 1px border at TargetSize (NOT at LogicalSize).
	Border bool

	// FontOverrides provides per-line font overrides. Nil entries use Config.Font.
	FontOverrides []font.Face

	// Label is the Sprite label, max 128 chars. Defaults to "scaledtextbox".
	Label string
}

// maxLabelLen is the maximum number of characters allowed in a Label.
const maxLabelLen = 128

// Render produces a ScaledTextBox bitmap. Internally renders text at
// LogicalSize using textbox.Render, then scales to TargetSize using
// nearest-neighbor interpolation. If Border is true, draws a 1px border
// at TargetSize AFTER scaling (for crisp output). Returns nil for invalid
// LogicalSize or TargetSize (zero/negative dimensions).
func Render(cfg Config) *widgets.Sprite {
	// 1. Validate dimensions.
	if cfg.LogicalSize.X <= 0 || cfg.LogicalSize.Y <= 0 {
		return nil
	}
	if cfg.TargetSize.X <= 0 || cfg.TargetSize.Y <= 0 {
		return nil
	}

	// Border edge case: TargetSize 2×2 or smaller with Border=true → nil.
	if cfg.Border && (cfg.TargetSize.X <= 2 || cfg.TargetSize.Y <= 2) {
		return nil
	}

	// 2. Construct internal textbox.Config.
	innerConfig := textbox.Config{
		Bounds:        image.Rect(0, 0, cfg.LogicalSize.X, cfg.LogicalSize.Y),
		Text:          cfg.Text,
		Font:          cfg.Font,
		Alignment:     cfg.Alignment,
		VAlign:        cfg.VAlign,
		Overflow:      cfg.Overflow,
		Foreground:    cfg.Foreground,
		LineSpacing:   cfg.LineSpacing,
		PadX:          cfg.PadX,
		PadY:          cfg.PadY,
		Border:        false, // Border is NOT applied at logical resolution.
		FontOverrides: cfg.FontOverrides,
		Label:         "inner", // Not used externally.
	}

	// 3. Call textbox.Render — if nil, return nil.
	innerResult := textbox.Render(innerConfig)
	if innerResult == nil {
		return nil
	}

	// 4. Scale result using nearest-neighbor.
	scaled := scale.NearestNeighbor(innerResult.Image, cfg.TargetSize.X, cfg.TargetSize.Y)
	if scaled == nil {
		return nil
	}

	// 5. If Border=true, draw 1px border on scaled image at TargetSize edges.
	if cfg.Border {
		fg := cfg.Foreground
		if fg == (color.RGBA{}) {
			fg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		widgets.DrawBorder(scaled, cfg.TargetSize.X, cfg.TargetSize.Y, fg)
	}

	// 6. Resolve label.
	label := cfg.Label
	if label == "" {
		label = "scaledtextbox"
	}
	if len([]rune(label)) > maxLabelLen {
		label = string([]rune(label)[:maxLabelLen])
	}

	// 7. Return result.
	return &widgets.Sprite{
		Image:    scaled,
		Position: cfg.Position,
		Label:    label,
	}
}
