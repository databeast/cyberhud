package ticker

import (
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/style"
)

type ShowcaseCase struct {
	Name       string
	Width      int
	Height     int
	Mono       bool
	Style      string
	ShowBorder bool
	Direction  string
	AutoScroll int
	Feed       []source.LineDirective
	Ticks      int
	LineMode   string
	Font       string
}

func TickerRegistryEnumerate() []style.Style[source.TickerSnapshot, source.Policy] {
	return tickerRegistry.Enumerate()
}
func AllowedFontTiers() []string      { return source.AllowedFontTiers }
func NormalizePolicy(p Policy) Policy { return normalizePolicy(p) }
func NewSnapshotterForTest() interface {
	SnapshotPolicy() map[string]interface{}
	RestorePolicy(map[string]interface{}) error
} {
	return tickerSnapshotter{}
}

func ExportedShowcaseCases() []ShowcaseCase {
	feed1 := []source.LineDirective{{Text: "BREAKING: Federal Reserve signals rate cuts - markets rally as investors anticipate easing cycle through Q3 2025 amid cooling inflation data"}}
	feed2 := []source.LineDirective{{Text: "BTC $67,234"}, {Text: "ETH $3,891"}, {Text: "SOL $142.50"}}
	return []ShowcaseCase{
		{Name: "128x32_plain_horizontal", Width: 128, Height: 32, Mono: true, Style: "", Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "128x64_plain_vertical", Width: 128, Height: 64, Mono: true, Style: "", Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 3},
		{Name: "128x64_plain_none", Width: 128, Height: 64, Mono: true, Style: "", Direction: "none", Feed: feed2},
		{Name: "160x80_bordered_horizontal", Width: 160, Height: 80, Style: "", ShowBorder: true, Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "296x128_plain_vertical", Width: 296, Height: 128, Mono: true, Style: "", Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 4},
		{Name: "400x300_bordered_none", Width: 400, Height: 300, Mono: true, Style: "", ShowBorder: true, Direction: "none", Feed: feed2},
		{Name: "64x128_plain_vertical", Width: 64, Height: 128, Mono: true, Style: "", Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 3},
		{Name: "80x160_plain_horizontal", Width: 80, Height: 160, Style: "", Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "240x135_plain_horizontal", Width: 240, Height: 135, Style: "", Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "240x135_bordered_horizontal", Width: 240, Height: 135, Style: "", ShowBorder: true, Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "240x135_bordered_vertical", Width: 240, Height: 135, Style: "", ShowBorder: true, Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 5},
		{Name: "135x240_bordered_horizontal", Width: 135, Height: 240, Style: "", ShowBorder: true, Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "320x240_plain_horizontal", Width: 320, Height: 240, Style: "", Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "320x240_bordered_vertical", Width: 320, Height: 240, Style: "", ShowBorder: true, Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 5},
		{Name: "320x240_plain_none", Width: 320, Height: 240, Style: "", Direction: "none", Feed: feed1, LineMode: "clip"},
		{Name: "320x480_bordered_horizontal", Width: 320, Height: 480, Style: "", ShowBorder: true, Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "480x320_bordered_horizontal", Width: 480, Height: 320, Style: "", ShowBorder: true, Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
		{Name: "480x320_plain_vertical", Width: 480, Height: 320, Style: "", Direction: "vertical", AutoScroll: 50, Feed: feed2, Ticks: 5},
		{Name: "480x320_bordered_none", Width: 480, Height: 320, Style: "", ShowBorder: true, Direction: "none", Feed: feed2},
		{Name: "800x480_plain_horizontal", Width: 800, Height: 480, Style: "", Direction: "horizontal", AutoScroll: 50, Feed: feed1, Ticks: 5},
	}
}
