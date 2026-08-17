package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/style/color"
)

var AllowedShapes = []string{"auto", "linear", "ring", "arc", "pie"}
var AllowedLabelTiers = []string{"auto", "small", "normal", "large", "fullsize"}

func ValidShape(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, allowed := range AllowedShapes {
		if allowed == s {
			return true
		}
	}
	return false
}

func ValidLabelTier(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, allowed := range AllowedLabelTiers {
		if allowed == s {
			return true
		}
	}
	return false
}

func ValidAccent(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "none" {
		return true
	}
	for _, allowed := range color.Names() {
		if allowed == s {
			return true
		}
	}
	return false
}

// Policy captures all runtime-configurable parameters for gauges mode.
type Policy struct {
	Style      string
	Shape      string
	ShowLabels bool
	LabelTier  string
	Accent     string

	DefaultMin float64
	DefaultMax float64

	Columns    int
	Rows       int
	TileGapPx  int
	PaddingPct int
}

func DefaultPolicy() Policy {
	return Policy{
		Style:      "",
		Shape:      "auto",
		ShowLabels: true,
		LabelTier:  "normal",
		Accent:     "cyan",
		DefaultMin: 0,
		DefaultMax: 100,
		Columns:    0,
		Rows:       0,
		TileGapPx:  1,
		PaddingPct: 0,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "style", Type: "string", Summary: "Visual style name or empty for automatic selection.", Default: "", Allowed: nil},
		{Key: "shape", Type: "string", Summary: "Default gauge shape when the payload does not override it.", Default: "auto", Allowed: AllowedShapes},
		{Key: "show_labels", Type: "bool", Summary: "Render labels above each gauge tile.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "label_tier", Type: "string", Summary: "Font tier used for gauge labels.", Default: "normal", Allowed: AllowedLabelTiers},
		{Key: "accent", Type: "string", Summary: "Named accent color for gauge fills.", Default: "cyan", Allowed: append(color.Names(), "none")},
		{Key: "default_min", Type: "float", Summary: "Default minimum value when a payload omits min.", Default: "0"},
		{Key: "default_max", Type: "float", Summary: "Default maximum value when a payload omits max.", Default: "100"},
		{Key: "columns", Type: "int", Summary: "Explicit gauge grid columns, or 0 for auto.", Default: "0"},
		{Key: "rows", Type: "int", Summary: "Explicit gauge grid rows, or 0 for auto.", Default: "0"},
		{Key: "tile_gap_px", Type: "int", Summary: "Gap between gauge tiles in pixels.", Default: "1"},
		{Key: "padding_pct", Type: "int", Summary: "Layout inset percentage applied to the panel.", Default: "0"},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%t|%s|%s|%g|%g|%d|%d|%d|%d",
		p.Style, p.Shape, p.ShowLabels, p.LabelTier, p.Accent,
		p.DefaultMin, p.DefaultMax, p.Columns, p.Rows, p.TileGapPx, p.PaddingPct)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":       p.Style,
		"shape":       p.Shape,
		"show_labels": p.ShowLabels,
		"label_tier":  p.LabelTier,
		"accent":      p.Accent,
		"default_min": p.DefaultMin,
		"default_max": p.DefaultMax,
		"columns":     p.Columns,
		"rows":        p.Rows,
		"tile_gap_px": p.TileGapPx,
		"padding_pct": p.PaddingPct,
	}
}

func NormalizePolicy(p Policy) Policy {
	d := DefaultPolicy()

	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.Shape = strings.ToLower(strings.TrimSpace(p.Shape))
	if !ValidShape(p.Shape) {
		p.Shape = d.Shape
	}

	p.LabelTier = strings.ToLower(strings.TrimSpace(p.LabelTier))
	if !ValidLabelTier(p.LabelTier) {
		p.LabelTier = d.LabelTier
	}

	p.Accent = strings.ToLower(strings.TrimSpace(p.Accent))
	if !ValidAccent(p.Accent) {
		p.Accent = d.Accent
	}

	if p.DefaultMax <= p.DefaultMin {
		p.DefaultMax = p.DefaultMin + 1
	}
	if p.Columns < 0 {
		p.Columns = 0
	}
	if p.Rows < 0 {
		p.Rows = 0
	}
	if p.TileGapPx < 0 {
		p.TileGapPx = 0
	}
	if p.TileGapPx > 32 {
		p.TileGapPx = 32
	}
	if p.PaddingPct < 0 {
		p.PaddingPct = 0
	}
	if p.PaddingPct > 25 {
		p.PaddingPct = 25
	}

	return p
}
