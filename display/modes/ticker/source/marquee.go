package source

import (
	"image"
	"time"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/marquee"
)

// stripSet holds the active marquee.Strip instances and their creation context.
type stripSet struct {
	strips      []*marquee.Strip // One per directive; nil for non-scrolling lines
	panelWidth  int              // hints.PixelWidth when strips were created (for invalidation)
	pixelHeight int              // hints.PixelHeight when strips were created
	lastTickAt  time.Time        // Wall-clock time of last Tick invocation
}

// textPixelWidth returns the rendered pixel width of text for a given font.
// It computes the width as the number of runes multiplied by the font's GlyphAdvance.
// Returns 0 for empty text.
func textPixelWidth(text string, face font.Face) int {
	if face == nil {
		face = font.Default()
	}
	return utf8.RuneCountInString(text) * face.Metrics().GlyphAdvance
}

// ensureStrips creates or recreates Strip instances if needed.
// Called from CheckAdvance when direction=horizontal and AutoScrollMS > 0.
// Recreates when: strips are nil, directive count changed, or panel dimensions changed.
func ensureStrips(hints textlayout.TextHints, directives []LineDirective, policy Policy, borderInset int) {
	// Check if recreation is needed.
	current := feedState.strips
	if current.strips != nil &&
		len(directives) == len(current.strips) &&
		hints.PixelWidth == current.panelWidth &&
		hints.PixelHeight == current.pixelHeight {
		return
	}

	usableWidth := hints.PixelWidth - 2*borderInset
	fontTier := ResolveFontTier(policy, hints)
	strips := make([]*marquee.Strip, len(directives))

	// Guard: if usable width is zero or negative, all lines are static.
	if usableWidth <= 0 {
		feedState.strips = stripSet{
			strips:      strips,
			panelWidth:  hints.PixelWidth,
			pixelHeight: hints.PixelHeight,
			lastTickAt:  feedState.strips.lastTickAt,
		}
		return
	}

	for i, d := range directives {
		fontID := d.Font
		if fontID == "" {
			fontID = policy.Font
		}
		if fontID == "" || fontID == "auto" {
			fontID = "spleen"
		}
		face := ResolveFace(hints, fontID, fontTier)
		metrics := face.Metrics()

		textWidth := utf8.RuneCountInString(d.Text) * metrics.GlyphAdvance

		// Skip strip creation for pinned lines, zero-advance fonts, or text that fits.
		if d.Scroll == "pinned" || metrics.GlyphAdvance == 0 || textWidth <= usableWidth {
			strips[i] = nil
			continue
		}

		lineY := borderInset + i*metrics.RowHeight
		speed := 1000.0 / (float64(policy.AutoScrollMS) * float64(metrics.GlyphAdvance))

		strips[i] = marquee.New(marquee.Config{
			Bounds:    image.Rect(borderInset, lineY, hints.PixelWidth-borderInset, lineY+metrics.RowHeight),
			Direction: marquee.Horizontal,
			Font:      face,
			Source:    &marquee.FixedStringSource{Runes: []rune(d.Text)},
			Speed:     speed,
			Phase:     0,
		})
	}

	feedState.strips = stripSet{
		strips:      strips,
		panelWidth:  hints.PixelWidth,
		pixelHeight: hints.PixelHeight,
		lastTickAt:  feedState.strips.lastTickAt,
	}
}

// tickStrips advances all active strips by the given elapsed duration.
// Called from CheckAdvance after ensureStrips to animate scrolling.
func tickStrips(elapsed time.Duration) {
	for _, strip := range feedState.strips.strips {
		if strip != nil {
			strip.Tick(elapsed)
		}
	}
}

// discardStrips clears all Strip instances by resetting feedState.strips to its zero value.
// Called by SetFeed, SetText, and policy changes that leave horizontal mode.
func discardStrips() {
	feedState.strips = stripSet{}
}

// renderStripSprites calls RenderFrame on each active strip and returns
// the collected non-nil sprites. Called by BuildView to populate TickerSnapshot.StripSprites.
func renderStripSprites() []widgets.Sprite {
	var sprites []widgets.Sprite
	for _, strip := range feedState.strips.strips {
		if strip == nil {
			continue
		}
		sp := strip.RenderFrame()
		if sp != nil {
			sprites = append(sprites, *sp)
		}
	}
	return sprites
}
