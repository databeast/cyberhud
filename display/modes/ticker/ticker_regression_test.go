package ticker_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/ticker"
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestBuildViewDoesNotPanicWithoutCatalog(t *testing.T) {
	source.SetFeed([]source.LineDirective{{Text: "zero catalog fallback"}})
	t.Cleanup(func() {
		source.SetFeed(nil)
		ticker.SetPolicy(ticker.DefaultPolicy())
	})

	ticker.SetPolicy(ticker.DefaultPolicy())

	hints := textlayout.TextHints{
		PixelWidth:      240,
		PixelHeight:     64,
		GlyphAdvance:    6,
		RowHeight:       10,
		DefaultLineMode: textlayout.LineModeTruncate,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BuildView panicked without catalog: %v", r)
		}
	}()

	v := ticker.BuildView(hints)
	if len(v.Items) != 1 {
		t.Fatalf("BuildView item count=%d, want 1", len(v.Items))
	}
}
