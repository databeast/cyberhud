package layout_test

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/style/layout"
)

// A font with advance 9, glyph ink width 8 (spleen-8x16): the final glyph on a
// row occupies 8px, every preceding glyph steps 9px.
func spleen8x16() layout.Metrics { return layout.NewMetrics(9, 8, 18, "spleen-8x16") }

func TestTextWidth_Exact(t *testing.T) {
	m := spleen8x16()
	cases := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"A", 8},         // single glyph: just the ink width, no trailing advance gap
		{"AB", 9 + 8},    // 17: one advance + final ink width
		{"ABC", 9*2 + 8}, // 26
		{"hello", 9*4 + 8},
	}
	for _, c := range cases {
		if got := m.TextWidth(c.s); got != c.want {
			t.Errorf("TextWidth(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestTextWidth_RuneCountedNotByteCounted(t *testing.T) {
	m := spleen8x16()
	// A 3-rune string must measure as 3 glyphs regardless of byte length.
	s := "aé中" // 1 + 2 + 3 bytes, 3 runes
	if got, want := m.TextWidth(s), 9*2+8; got != want {
		t.Errorf("TextWidth(%q) = %d, want %d (3 runes)", s, got, want)
	}
}

func TestTextWidth_UnusableMetricsIsZero(t *testing.T) {
	var zero layout.Metrics // advance 0
	if got := zero.TextWidth("anything"); got != 0 {
		t.Errorf("zero Metrics TextWidth = %d, want 0", got)
	}
}

func TestTextWidth_FallsBackToAdvanceWhenNoGlyphWidth(t *testing.T) {
	// glyphWidth 0 → conservative: every glyph counts a full advance (never
	// under-reports, so a fit test using it can't claim a clipping string fits).
	m := layout.NewMetrics(9, 0, 18, "adv-only")
	if got, want := m.TextWidth("AB"), 9*2; got != want {
		t.Errorf("advance-only TextWidth(%q) = %d, want %d", "AB", got, want)
	}
	// The conservative value is >= the exact value for the same string/font.
	exact := spleen8x16()
	if m.TextWidth("hello") < exact.TextWidth("hello") {
		t.Error("advance-only width should be >= exact width (conservative)")
	}
}

func TestFits(t *testing.T) {
	m := spleen8x16()
	s := "hello" // TextWidth = 44
	if !m.Fits(s, 44) {
		t.Error("Fits should be true at exactly the needed width")
	}
	if m.Fits(s, 43) {
		t.Error("Fits should be false one pixel under")
	}
	if !m.Fits("", 0) {
		t.Error("empty string fits any width including 0")
	}
	if m.Fits("A", 0) {
		t.Error("non-empty string must not fit width 0")
	}
}

// The property that motivated this tool: for a monospace-ish font, TextWidth is
// monotonic in rune count, so a longer string never needs fewer pixels.
func TestTextWidth_MonotonicInLength(t *testing.T) {
	m := spleen8x16()
	prev := -1
	for n := 0; n <= 20; n++ {
		w := m.TextWidth(strings.Repeat("x", n))
		if w < prev {
			t.Fatalf("TextWidth not monotonic: n=%d gave %d < previous %d", n, w, prev)
		}
		prev = w
	}
}
