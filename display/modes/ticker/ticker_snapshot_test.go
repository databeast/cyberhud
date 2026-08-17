package ticker_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/ticker"
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a ticker
// style's name prefix. Ticker styles follow the naming convention:
// mono-*, color-*, grayscale-fast-*, grayscale-slow-*, mono-slow-*, color-slow-*, eink-*.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-slow-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "grayscale-fast-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "grayscale-slow-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	default:
		return testsnapshot.CategoryColor
	}
}

// TestTickerPNGSnapshots enumerates every registered ticker style and
// renders one deterministic PNG per style.
func TestTickerPNGSnapshots(t *testing.T) {
	const showcaseSnapshotDir = "snapshots"
	if err := os.RemoveAll(showcaseSnapshotDir); err != nil {
		t.Fatalf("failed to clean snapshot directory: %v", err)
	}
	if err := os.MkdirAll(showcaseSnapshotDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot directory: %v", err)
	}

	styles := ticker.TickerRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("ticker registry contains zero styles")
	}

	feed := []source.LineDirective{
		{Text: "BTC $67,234", Scroll: "pinned"},
		{Text: "ETH $3,891"},
		{Text: "SOL $142.50"},
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			width, height := reqs.MinWidth, reqs.MinHeight
			if width == 0 || height == 0 {
				t.Fatalf("style %q has zero dimensions", s.Name())
			}

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("ticker"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(categoryFromStyleName(s.Name())),
				testsnapshot.WithOutputDir(showcaseSnapshotDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithFrameCount(3),
				testsnapshot.WithReset(func() {
					source.SetFeed(nil)
					ticker.SetPolicy(ticker.DefaultPolicy())
				}),
				testsnapshot.WithPreRender(func() {
					source.SetFeed(feed)
					ticker.SetPolicy(ticker.Policy{
						Style:        "",
						Font:         "auto",
						FontTier:     "auto",
						LineMode:     textlayout.LineModeTruncate,
						Direction:    textlayout.TickerDirectionVertical,
						AutoScrollMS: 50,
						Accent:       "cyan",
					})
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}

// TestTickerHorizontalMarqueeSnapshot verifies that the ticker in horizontal marquee mode
// produces a valid PNG snapshot with scrolling content rendered through the full production
// pipeline via the snapshottest framework.
//

func TestTickerHorizontalMarqueeSnapshot(t *testing.T) {
	outputDir := filepath.Join("snapshots", "marquee")
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("failed to clean snapshot directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot directory: %v", err)
	}

	feed := []source.LineDirective{
		{Text: "Hello"},
		{Text: "This is a very long line of text that will definitely exceed the panel width and trigger horizontal marquee scrolling behavior"},
	}

	pngPath := testsnapshot.RenderSnapshot(t,
		testsnapshot.WithMode("ticker"),
		testsnapshot.WithDimensions(240, 135),
		testsnapshot.WithDisplayCategory(testsnapshot.CategoryColor),
		testsnapshot.WithOutputDir(outputDir),
		testsnapshot.WithBasename("ticker_marquee"),
		testsnapshot.WithFrameCount(5),
		testsnapshot.WithPreRender(func() {
			source.SetFeed(feed)
			ticker.SetPolicy(ticker.Policy{
				Style:        "",
				Font:         "auto",
				LineMode:     textlayout.LineModeTruncate,
				Direction:    "horizontal",
				AutoScrollMS: 50,
			})
		}),
	)

	testsnapshot.VerifyAll(t, pngPath, 240, 135)
}
