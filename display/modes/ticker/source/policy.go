package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Policy controls ticker rendering and auto-scroll behavior.
type Policy struct {
	Style        string
	Font         string
	FontTier     string
	LineMode     string
	Direction    string
	AutoScrollMS int
	Accent       string
	ShowBorder   bool
	ShowGlow     bool
}

// DefaultPolicy returns the baseline behavior that matches existing ticker UX.
func DefaultPolicy() Policy {
	return Policy{
		Style:        "",
		Font:         "auto",
		FontTier:     "auto",
		LineMode:     textlayout.LineModeTruncate,
		Direction:    textlayout.TickerDirectionVertical,
		AutoScrollMS: 0,
		Accent:       "cyan",
		ShowBorder:   false,
		ShowGlow:     false,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		{Key: "font_tier", Type: "string", Summary: "Font size tier for panel-appropriate text rendering.", Default: "auto", Allowed: AllowedFontTiers},
		{Key: "line_mode", Type: "string", Summary: "How long lines are constrained to panel width.", Default: textlayout.LineModeTruncate, Allowed: []string{textlayout.LineModeTruncate, textlayout.LineModeClip}},
		{Key: "direction", Type: "string", Summary: "Ticker auto-scroll direction when supported by the panel.", Default: textlayout.TickerDirectionVertical, Allowed: []string{textlayout.TickerDirectionVertical, "horizontal", textlayout.TickerDirectionNone}},
		{Key: "auto_scroll_ms", Type: "int", Summary: "Milliseconds between automatic ticker advances.", Default: "0"},
		{Key: "accent", Type: "string", Summary: "Named accent color from the shared palette.", Default: "cyan", Allowed: color.Names()},
		{Key: "show_border", Type: "bool", Summary: "Show decorative border frame.", Default: "false"},
		{Key: "show_glow", Type: "bool", Summary: "Show accent-colored glow background (ColorFast panels only).", Default: "false"},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%d|%s|%v|%v", p.Style, p.Font, p.FontTier, p.LineMode, p.Direction, p.AutoScrollMS, p.Accent, p.ShowBorder, p.ShowGlow)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":          p.Style,
		"font":           p.Font,
		"font_tier":      p.FontTier,
		"line_mode":      p.LineMode,
		"direction":      p.Direction,
		"auto_scroll_ms": p.AutoScrollMS,
		"accent":         p.Accent,
		"show_border":    p.ShowBorder,
		"show_glow":      p.ShowGlow,
	}
}

func normalizePolicy(policy Policy) Policy {
	d := DefaultPolicy()

	// Normalize style: empty means auto-detect via fitness.
	policy.Style = strings.ToLower(strings.TrimSpace(policy.Style))

	// Reset deprecated style names that are no longer in the registry.
	if policy.Style == "plain" || policy.Style == "bordered" {
		policy.Style = ""
	}

	// Normalize font.
	policy.Font = strings.TrimSpace(policy.Font)
	if policy.Font == "" {
		policy.Font = "auto"
	}

	// Normalize font_tier.
	policy.FontTier = strings.ToLower(strings.TrimSpace(policy.FontTier))
	if !ValidFontTier(policy.FontTier) {
		policy.FontTier = d.FontTier
	}

	// Normalize line_mode.
	policy.LineMode = strings.ToLower(strings.TrimSpace(policy.LineMode))
	if policy.LineMode != textlayout.LineModeClip && policy.LineMode != textlayout.LineModeTruncate {
		policy.LineMode = d.LineMode
	}

	// Normalize direction.
	policy.Direction = strings.ToLower(strings.TrimSpace(policy.Direction))
	if policy.Direction != textlayout.TickerDirectionVertical && policy.Direction != textlayout.TickerDirectionNone && policy.Direction != "horizontal" {
		policy.Direction = d.Direction
	}

	if policy.AutoScrollMS < 0 {
		policy.AutoScrollMS = 0
	}

	// Normalize accent.
	policy.Accent = strings.ToLower(strings.TrimSpace(policy.Accent))
	if !ValidAccent(policy.Accent) {
		policy.Accent = d.Accent
	}

	return policy
}
