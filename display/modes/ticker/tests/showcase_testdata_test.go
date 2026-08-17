package tests

import (
	"github.com/databeast/cyberhud/display/modes/ticker/source"
)

// showcaseCase defines a single test configuration for the ticker snapshot showcase.
// Shared between ticker_snapshot_test.go (external) and ticker_prop_test.go (internal).
type showcaseCase struct {
	name       string // Subtest name AND PNGPanel basename
	width      int    // Panel pixel width
	height     int    // Panel pixel height
	mono       bool   // Whether this is a monochrome panel
	style      string // Style name (empty = fitness auto-selection)
	showBorder bool   // Whether to render decorative border frame
	direction  string // "horizontal", "vertical", or "none"
	autoScroll int    // AutoScrollMS value
	feed       []source.LineDirective
	ticks      int    // Number of render frames for animation advancement
	lineMode   string // "truncate" (default) or "clip"
	font       string // "auto" (default) or a specific font ID
}

var showcaseCases = []showcaseCase{
	// 1: Tiny OLED, single scroll
	{
		name:       "128x32_plain_horizontal",
		width:      128,
		height:     32,
		mono:       true,
		style:      "",
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"},
		},
		ticks: 5,
	},
	// 2: Small OLED rotation
	{
		name:       "128x64_plain_vertical",
		width:      128,
		height:     64,
		mono:       true,
		style:      "",
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234"},
			{Text: "ETH $3,891"},
		},
		ticks: 3,
	},
	// 3: Static mono baseline
	{
		name:       "128x64_plain_none",
		width:      128,
		height:     64,
		mono:       true,
		style:      "",
		direction:  "none",
		autoScroll: 0,
		feed: []source.LineDirective{
			{Text: "BTC $67,234"},
			{Text: "ETH $3,891"},
		},
		ticks: 0,
	},
	// 4: Small TFT, border + marquee
	{
		name:       "160x80_bordered_horizontal",
		width:      160,
		height:     80,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025"},
		},
		ticks: 5,
	},
	// 5: Wide e-ink, rotation
	{
		name:       "296x128_plain_vertical",
		width:      296,
		height:     128,
		mono:       true,
		style:      "",
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234"},
			{Text: "ETH $3,891"},
			{Text: "SOL $142.50"},
		},
		ticks: 4,
	},
	// 6: Large e-ink, static dense
	{
		name:       "400x300_bordered_none",
		width:      400,
		height:     300,
		mono:       true,
		style:      "",
		showBorder: true,
		direction:  "none",
		autoScroll: 0,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "SOL $142.50", Scroll: "pinned"},
			{Text: "DOGE $0.124", Scroll: "pinned"},
		},
		ticks: 0,
	},
	// 7: Tiny portrait
	{
		name:       "64x128_plain_vertical",
		width:      64,
		height:     128,
		mono:       true,
		style:      "",
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67K"},
			{Text: "ETH $3.8K"},
		},
		ticks: 3,
	},
	// 8: Narrow portrait marquee
	{
		name:       "80x160_plain_horizontal",
		width:      80,
		height:     160,
		mono:       false,
		style:      "",
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BREAKING: Markets surge on Federal Reserve announcement of rate cuts for Q3 2025"},
			{Text: "ALERT: Crypto markets rally — Bitcoin breaks $70K resistance as institutional buying accelerates"},
		},
		ticks: 5,
	},
	// 9: Mid TFT, dense mixed marquee
	{
		name:       "240x135_plain_horizontal",
		width:      240,
		height:     135,
		mono:       false,
		style:      "",
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end"},
		},
		ticks: 5,
	},
	// 10: Mid TFT bordered, deep scroll
	{
		name:       "240x135_bordered_horizontal",
		width:      240,
		height:     135,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025"},
			{Text: "ALERT: European Central Bank follows Fed guidance — euro weakens against dollar as divergence narrows"},
			{Text: "UPDATE: Oil prices surge past $95 as OPEC+ extends production cuts through end of year"},
		},
		ticks: 5,
	},
	// 11: Mid TFT bordered rotation
	{
		name:       "240x135_bordered_vertical",
		width:      240,
		height:     135,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891"},
			{Text: "SOL $142.50"},
			{Text: "DOGE $0.124", Scroll: "pinned"},
			{Text: "ADA $0.45"},
		},
		ticks: 5,
	},
	// 12: Portrait bordered, max density
	{
		name:       "135x240_bordered_horizontal",
		width:      135,
		height:     240,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle"},
			{Text: "ALERT: Crypto markets rally — Bitcoin breaks $70K resistance as institutional buying accelerates"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end"},
			{Text: "LIVE: Senate committee hearing on stablecoin regulation draws heated debate from industry leaders"},
		},
		ticks: 5,
	},
	// 13: Medium TFT, font + dense marquee
	{
		name:       "320x240_plain_horizontal",
		width:      320,
		height:     240,
		mono:       false,
		style:      "",
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices worldwide"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end amid softening labor data"},
			{Text: "LIVE: Senate committee hearing on cryptocurrency regulation draws heated debate from industry leaders and consumer advocates"},
		},
		ticks: 5,
	},
	// 14: Medium bordered, 6-line rotation
	{
		name:       "320x240_bordered_vertical",
		width:      320,
		height:     240,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234"},
			{Text: "ETH $3,891"},
			{Text: "SOL $142.50"},
			{Text: "DOGE $0.124"},
			{Text: "ADA $0.45"},
			{Text: "DOT $7.23"},
		},
		ticks: 5,
	},
	// 15: Clip mode static
	{
		name:       "320x240_plain_none",
		width:      320,
		height:     240,
		mono:       false,
		style:      "",
		direction:  "none",
		autoScroll: 0,
		feed: []source.LineDirective{
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in rate cuts"},
			{Text: "LIVE: Senate hearing on crypto regulation draws heated debate"},
		},
		ticks:    0,
		lineMode: "clip",
	},
	// 16: Tall portrait, max line count
	{
		name:       "320x480_bordered_horizontal",
		width:      320,
		height:     480,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "SOL $142.50", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices worldwide"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end"},
			{Text: "LIVE: Senate committee hearing on cryptocurrency regulation draws heated debate from industry leaders"},
			{Text: "ANALYSIS: Housing market shows signs of cooling as mortgage rates remain elevated above 7% threshold"},
		},
		ticks: 5,
	},
	// 17: Large TFT bordered, deep marquee
	{
		name:       "480x320_bordered_horizontal",
		width:      480,
		height:     320,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices worldwide in historic session"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end amid softening labor market data"},
			{Text: "LIVE: Senate committee hearing on cryptocurrency regulation draws heated debate from industry leaders and consumer advocates alike"},
		},
		ticks: 5,
	},
	// 18: Large TFT, 7-line rotation
	{
		name:       "480x320_plain_vertical",
		width:      480,
		height:     320,
		mono:       false,
		style:      "",
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234"},
			{Text: "ETH $3,891"},
			{Text: "SOL $142.50"},
			{Text: "DOGE $0.124"},
			{Text: "ADA $0.45"},
			{Text: "DOT $7.23"},
			{Text: "AVAX $35.67"},
		},
		ticks: 5,
	},
	// 19: Large static bordered density
	{
		name:       "480x320_bordered_none",
		width:      480,
		height:     320,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "none",
		autoScroll: 0,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "SOL $142.50", Scroll: "pinned"},
			{Text: "DOGE $0.124", Scroll: "pinned"},
			{Text: "ADA $0.45", Scroll: "pinned"},
		},
		ticks: 0,
	},
	// 20: Largest TFT, max density + deep scroll
	{
		name:       "800x480_plain_horizontal",
		width:      800,
		height:     480,
		mono:       false,
		style:      "",
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data and rising unemployment claims"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices worldwide in historic trading session with record volume"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end amid softening labor market data and declining consumer confidence"},
			{Text: "LIVE: Senate committee hearing on cryptocurrency regulation draws heated debate from industry leaders and consumer advocates over stablecoin oversight framework"},
			{Text: "ANALYSIS: Housing market shows signs of cooling as mortgage rates remain elevated above 7% threshold — new construction permits decline for third consecutive month"},
			{Text: "GLOBAL: European Central Bank announces emergency meeting as sovereign debt concerns resurface in peripheral eurozone economies amid political instability in key member states"},
		},
		ticks: 5,
	},
	// 21: Largest bordered, complex frame
	{
		name:       "800x480_bordered_horizontal",
		width:      800,
		height:     480,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data and rising unemployment claims"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices worldwide in historic trading session with record volume"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end amid softening labor market data and declining consumer confidence"},
			{Text: "LIVE: Senate committee hearing on cryptocurrency regulation draws heated debate from industry leaders and consumer advocates over stablecoin oversight framework"},
			{Text: "ANALYSIS: Housing market shows signs of cooling as mortgage rates remain elevated above 7% — new construction permits decline for third consecutive month nationwide"},
		},
		ticks: 5,
	},
	// 22: Largest bordered rotation
	{
		name:       "800x480_bordered_vertical",
		width:      800,
		height:     480,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "vertical",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891"},
			{Text: "SOL $142.50"},
			{Text: "DOGE $0.124", Scroll: "pinned"},
			{Text: "ADA $0.45"},
			{Text: "DOT $7.23"},
			{Text: "AVAX $35.67"},
			{Text: "LINK $14.82"},
		},
		ticks: 5,
	},
	// 23: Large static max density
	{
		name:       "800x480_plain_none",
		width:      800,
		height:     480,
		mono:       false,
		style:      "",
		direction:  "none",
		autoScroll: 0,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "ETH $3,891", Scroll: "pinned"},
			{Text: "SOL $142.50", Scroll: "pinned"},
			{Text: "DOGE $0.124", Scroll: "pinned"},
			{Text: "ADA $0.45", Scroll: "pinned"},
			{Text: "DOT $7.23", Scroll: "pinned"},
		},
		ticks: 0,
	},
	// 24: Square bordered marquee
	{
		name:       "240x240_bordered_horizontal",
		width:      240,
		height:     240,
		mono:       false,
		style:      "",
		showBorder: true,
		direction:  "horizontal",
		autoScroll: 50,
		feed: []source.LineDirective{
			{Text: "BTC $67,234", Scroll: "pinned"},
			{Text: "BREAKING: Federal Reserve signals rate cuts — markets rally as investors anticipate easing cycle through Q3 2025"},
			{Text: "ALERT: S&P 500 hits all-time high as tech sector leads broad market gains across all major indices"},
			{Text: "UPDATE: Treasury yields fall sharply as bond traders price in multiple rate cuts before year end"},
		},
		ticks: 5,
	},
}
