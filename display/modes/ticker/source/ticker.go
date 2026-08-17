package source

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// Allowed FontTier values.
var AllowedFontTiers = []string{"auto", "small", "normal", "large", "fullsize"}

// ValidFontTier reports whether s is a valid font tier value.
func ValidFontTier(s string) bool {
	for _, t := range AllowedFontTiers {
		if s == t {
			return true
		}
	}
	return false
}

// ValidAccent reports whether s is a recognized accent name.
// Recognized values are those returned by color.Names() plus "none".
func ValidAccent(s string) bool {
	if s == "none" {
		return true
	}
	for _, n := range color.Names() {
		if s == n {
			return true
		}
	}
	return false
}

var feedState = struct {
	sync.RWMutex
	directives      []LineDirective
	policy          Policy
	autoScrollState AutoScrollState
	strips          stripSet
}{
	directives: []LineDirective{{Text: "(ticker idle)"}},
	policy:     DefaultPolicy(),
	autoScrollState: AutoScrollState{
		currentIdx: 0,
		// Left as the zero Time deliberately. Package initialisation happens
		// before a snapshot can freeze the clock, so seeding from the clock
		// here would capture a wall-clock instant that a later frozen clock
		// sits far behind, producing a negative elapsed and permanently
		// stalling auto-scroll. CheckAdvance seeds it on first use instead.
	},
}

// AutoScrollState tracks vertical/horizontal ticker scroll position.
type AutoScrollState struct {
	currentIdx    int
	lastAdvanceAt time.Time
}

// SetFeed atomically replaces the feed buffer with validated directives.
// Resets auto-scroll state on replacement. Defaults to idle placeholder on empty input.
func SetFeed(directives []LineDirective) {
	feedState.Lock()
	defer feedState.Unlock()
	next := make([]LineDirective, 0, len(directives))
	for _, d := range directives {
		d.Text = strings.TrimSpace(d.Text)
		if d.Text == "" {
			continue
		}
		next = append(next, d)
	}
	if len(next) == 0 {
		next = []LineDirective{{Text: "(ticker idle)"}}
	}
	feedState.directives = next
	feedState.autoScrollState = AutoScrollState{
		currentIdx:    0,
		lastAdvanceAt: frameclock.Now(),
	}
	discardStrips()
}

// SetText replaces the ticker feed with externally supplied text lines.
// This is a backward-compatible wrapper that converts []string to []LineDirective.
func SetText(lines []string) {
	directives := make([]LineDirective, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		directives = append(directives, LineDirective{Text: line})
	}
	SetFeed(directives)
}

// Snapshot returns the current ticker feed lines as []string for backward compatibility.
func Snapshot() []string {
	feedState.RLock()
	defer feedState.RUnlock()
	out := make([]string, len(feedState.directives))
	for i, d := range feedState.directives {
		out[i] = d.Text
	}
	return out
}

// FeedSnapshot returns a deep copy of the current feed buffer as []LineDirective.
func FeedSnapshot() []LineDirective {
	feedState.RLock()
	defer feedState.RUnlock()
	out := make([]LineDirective, len(feedState.directives))
	copy(out, feedState.directives)
	return out
}

// SerializeFeed returns the current feed buffer as a pretty-printed JSON array string.
func SerializeFeed() (string, error) {
	snap := FeedSnapshot()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serialize feed: %v", err)
	}
	return string(data), nil
}

func FormatPolicyResponse(p Policy) string {
	return fmt.Sprintf("OK ticker policy style=%s font=%s font_tier=%s line_mode=%s direction=%s auto_scroll_ms=%d accent=%s show_border=%v show_glow=%v", p.Style, p.Font, p.FontTier, p.LineMode, p.Direction, p.AutoScrollMS, p.Accent, p.ShowBorder, p.ShowGlow)
}

// CheckAdvance updates auto-scroll state if sufficient time has passed.
// Returns the effective scroll offset (line index for vertical) constrained by panel capabilities.
// For horizontal direction, advances marquee pixel offsets (side effect) and returns 0.
func CheckAdvance(hints textlayout.TextHints) int {
	feedState.Lock()
	defer feedState.Unlock()

	effective := EffectivePolicy(feedState.policy, hints)
	if effective.AutoScrollMS <= 0 || effective.Direction == textlayout.TickerDirectionNone {
		feedState.autoScrollState.currentIdx = 0
		return 0
	}

	now := frameclock.Now()

	// Seed on first use rather than at package initialisation, so the baseline
	// comes from whichever clock is active when the mode first renders.
	if feedState.autoScrollState.lastAdvanceAt.IsZero() {
		feedState.autoScrollState.lastAdvanceAt = now
		return feedState.autoScrollState.currentIdx
	}

	elapsed := now.Sub(feedState.autoScrollState.lastAdvanceAt)
	if elapsed < time.Duration(effective.AutoScrollMS)*time.Millisecond {
		return feedState.autoScrollState.currentIdx
	}

	feedState.autoScrollState.lastAdvanceAt = now

	switch effective.Direction {
	case textlayout.TickerDirectionVertical:
		feedState.autoScrollState.currentIdx++
		if feedState.autoScrollState.currentIdx >= len(feedState.directives) {
			feedState.autoScrollState.currentIdx = 0
		}
	case "horizontal":
		borderInset := 0
		ensureStrips(hints, feedState.directives, effective, borderInset)
		now := frameclock.Now()
		elapsed := now.Sub(feedState.strips.lastTickAt)
		if feedState.strips.lastTickAt.IsZero() || elapsed <= 0 {
			elapsed = time.Duration(effective.AutoScrollMS) * time.Millisecond
		}
		tickStrips(elapsed)
		feedState.strips.lastTickAt = now
	}

	return feedState.autoScrollState.currentIdx
}

// PolicySnapshot returns the current ticker policy.
func PolicySnapshot() Policy {
	feedState.RLock()
	defer feedState.RUnlock()
	return feedState.policy
}

// ReplacePolicy installs an already-normalized policy and clears stale horizontal strips as needed.
func ReplacePolicy(policy Policy) {
	feedState.Lock()
	defer feedState.Unlock()
	prev := feedState.policy
	feedState.policy = policy
	next := feedState.policy

	if (prev.Direction == "horizontal" && next.Direction != "horizontal") ||
		(prev.AutoScrollMS > 0 && next.AutoScrollMS <= 0) {
		feedState.strips = stripSet{}
	}
}

// RenderStripSprites renders active horizontal marquee strips as sprites.
func RenderStripSprites() []widgets.Sprite {
	feedState.RLock()
	defer feedState.RUnlock()
	return renderStripSprites()
}

// FormatLines applies panel-specific width constraints to ticker rows.
// leftPad and rightPad are pixel paddings used by the caller's layout.
func FormatLines(lines []string, hints textlayout.TextHints, policy Policy, leftPad, rightPad int) []string {
	out := make([]string, len(lines))
	maxChars := textlayout.MaxCharsPerRow(hints, leftPad+rightPad)
	effective := EffectivePolicy(policy, hints)
	for i, line := range lines {
		if maxChars <= 0 {
			out[i] = ""
			continue
		}
		switch effective.LineMode {
		case textlayout.LineModeClip:
			out[i] = clipToChars(line, maxChars)
		default:
			out[i] = textlayout.Truncate(line, maxChars)
		}
	}
	return out
}

// EffectivePolicy returns policy normalized for provided panel capabilities.
func EffectivePolicy(policy Policy, hints textlayout.TextHints) Policy {
	policy = normalizePolicy(policy)
	if !hints.SupportsAutoScroll || isSlowRefresh(hints.Capability) {
		policy.AutoScrollMS = 0
	}
	switch policy.Direction {
	case textlayout.TickerDirectionVertical:
		if !hints.SupportsVerticalScroll {
			policy.Direction = fallbackDirection(hints)
		}
	case "horizontal":
		if !hints.SupportsHorizontalScroll {
			policy.Direction = fallbackDirection(hints)
		}
	case textlayout.TickerDirectionNone:
		// Explicit static mode.
	default:
		policy.Direction = fallbackDirection(hints)
	}
	return policy
}

// isSlowRefresh returns true when the capability indicates a slow-refresh
// (e-ink) panel that cannot support animation or scrolling.
func isSlowRefresh(capability int) bool {
	return capability == textlayout.CapMonoSlow ||
		capability == textlayout.CapGrayscaleSlow ||
		capability == textlayout.CapColorSlow
}

func fallbackDirection(hints textlayout.TextHints) string {
	if strings.TrimSpace(hints.DefaultTickerDirection) != "" {
		return strings.ToLower(strings.TrimSpace(hints.DefaultTickerDirection))
	}
	if hints.SupportsVerticalScroll {
		return textlayout.TickerDirectionVertical
	}
	if hints.SupportsHorizontalScroll {
		return "horizontal"
	}
	return textlayout.TickerDirectionNone
}

func clipToChars(s string, maxChars int) string {
	r := []rune(s)
	if len(r) <= maxChars {
		return s
	}
	return string(r[:maxChars])
}
