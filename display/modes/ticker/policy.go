package ticker

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// tickerSnapshotter implements catalog.PolicySnapshotter for the ticker mode.
type tickerSnapshotter struct{}

// SnapshotPolicy returns the current ticker policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (tickerSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := PolicySnapshot()
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

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (tickerSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["font"]; ok {
		if s, ok := v.(string); ok {
			p.Font = s
		}
	}
	if v, ok := data["font_tier"]; ok {
		if s, ok := v.(string); ok {
			p.FontTier = s
		}
	}
	if v, ok := data["line_mode"]; ok {
		if s, ok := v.(string); ok {
			p.LineMode = s
		}
	}
	if v, ok := data["direction"]; ok {
		if s, ok := v.(string); ok {
			p.Direction = s
		}
	}
	if v, ok := data["auto_scroll_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.AutoScrollMS = n
		}
	}
	if v, ok := data["accent"]; ok {
		if s, ok := v.(string); ok {
			p.Accent = s
		}
	}
	if v, ok := data["show_border"]; ok {
		if b, ok := v.(bool); ok {
			p.ShowBorder = b
		}
	}
	if v, ok := data["show_glow"]; ok {
		if b, ok := v.(bool); ok {
			p.ShowGlow = b
		}
	}

	SetPolicy(p)
	return nil
}

// toInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

func init() {
	catalog.RegisterSnapshotter("ticker", tickerSnapshotter{})
}

// PolicySnapshot returns the current ticker policy.
func PolicySnapshot() Policy {
	return source.PolicySnapshot()
}

// SetPolicy replaces ticker policy using safe normalization.
// Discards strips if the effective direction moves away from horizontal
// or if AutoScrollMS becomes zero or negative (scroll requires a positive interval).
func SetPolicy(policy Policy) {
	source.ReplacePolicy(normalizePolicy(policy))
}

func normalizePolicy(policy Policy) Policy {
	d := source.DefaultPolicy()

	policy.Style = strings.ToLower(strings.TrimSpace(policy.Style))
	if policy.Style != "" && tickerRegistry.Lookup(policy.Style) == nil {
		policy.Style = ""
	}

	policy.Font = strings.TrimSpace(policy.Font)
	if policy.Font == "" {
		policy.Font = "auto"
	}

	policy.FontTier = strings.ToLower(strings.TrimSpace(policy.FontTier))
	if !source.ValidFontTier(policy.FontTier) {
		policy.FontTier = d.FontTier
	}

	policy.LineMode = strings.ToLower(strings.TrimSpace(policy.LineMode))
	if policy.LineMode != textlayout.LineModeClip && policy.LineMode != textlayout.LineModeTruncate {
		policy.LineMode = d.LineMode
	}

	policy.Direction = strings.ToLower(strings.TrimSpace(policy.Direction))
	if policy.Direction != textlayout.TickerDirectionVertical && policy.Direction != textlayout.TickerDirectionNone && policy.Direction != "horizontal" {
		policy.Direction = d.Direction
	}

	if policy.AutoScrollMS < 0 {
		policy.AutoScrollMS = 0
	}

	policy.Accent = strings.ToLower(strings.TrimSpace(policy.Accent))
	if !source.ValidAccent(policy.Accent) {
		policy.Accent = d.Accent
	}

	return policy
}
