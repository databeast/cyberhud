package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the pager display mode.
type Policy struct {
	Source      string // Absolute path to file, pipe, or socket; empty = no source configured.
	Mode        string // Fast-display presentation: scroll or typer.
	ScrollSpeed int    // Pixels per second for smooth scroll [1..1000], default 60.
	TypeSpeed   int    // Characters per second for typer [1..1000], default 30.
	MaxLines    int    // Buffer capacity [1..1000], default 24.
	ScanMS      int    // Reconnect/retry interval in ms [100..30000], default 500.
	Font        string // Font family identifier; empty = auto.
	Style       string // Visual style name; empty = default.
	FadeOutMS   int    // Fade-out duration in ms for page transitions, default 300.
	FadeInMS    int    // Fade-in duration in ms for page transitions, default 300.
	LineTimeMS  int    // Per-line reading time in ms for page cadence, default 1000.
	MaxWaitS    int    // Max wait for full page in seconds, default 30.
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%d|%d|%d|%d|%s|%s|%d|%d|%d|%d",
		p.Source, p.Mode, p.ScrollSpeed, p.TypeSpeed, p.MaxLines, p.ScanMS,
		p.Font, p.Style,
		p.FadeOutMS, p.FadeInMS, p.LineTimeMS, p.MaxWaitS)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"source":       p.Source,
		"mode":         p.Mode,
		"scroll_speed": p.ScrollSpeed,
		"type_speed":   p.TypeSpeed,
		"max_lines":    p.MaxLines,
		"scan_ms":      p.ScanMS,
		"font":         p.Font,
		"style":        p.Style,
		"fade_out_ms":  p.FadeOutMS,
		"fade_in_ms":   p.FadeInMS,
		"line_time_ms": p.LineTimeMS,
		"max_wait_s":   p.MaxWaitS,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "source", Type: "string", Summary: "Absolute path to file, pipe, or socket to tail.", Default: ""},
		{Key: "mode", Type: "string", Summary: "Fast-display presentation strategy.", Default: "scroll", Allowed: []string{"scroll", "typer"}},
		{Key: "scroll_speed", Type: "int", Summary: "Pixels per second for smooth scroll (1-1000).", Default: "60"},
		{Key: "type_speed", Type: "int", Summary: "Characters per second for typer (1-1000).", Default: "30"},
		{Key: "max_lines", Type: "int", Summary: "Maximum buffered lines (1-1000).", Default: "24"},
		{Key: "scan_ms", Type: "int", Summary: "Reconnect interval in milliseconds (100-30000).", Default: "500"},
		{Key: "font", Type: "string", Summary: "Font family identifier or empty for auto.", Default: ""},
		{Key: "fade_out_ms", Type: "int", Summary: "Fade-out duration in ms for page transitions.", Default: "300"},
		{Key: "fade_in_ms", Type: "int", Summary: "Fade-in duration in ms for page transitions.", Default: "300"},
		{Key: "line_time_ms", Type: "int", Summary: "Per-line reading time in ms for page cadence.", Default: "1000"},
		{Key: "max_wait_s", Type: "int", Summary: "Max wait seconds for a full page before partial display.", Default: "30"},
	}
}

// DefaultPolicy returns the default pager policy with all fields initialized.
func DefaultPolicy() Policy {
	return Policy{
		Source:      "",
		Mode:        "scroll",
		ScrollSpeed: 60,
		TypeSpeed:   30,
		MaxLines:    24,
		ScanMS:      500,
		Font:        "",
		Style:       "",
		FadeOutMS:   300,
		FadeInMS:    300,
		LineTimeMS:  1000,
		MaxWaitS:    30,
	}
}

// normalizePolicy ensures all numeric policy fields are within their valid
// ranges, clamping out-of-range entries to their nearest valid bound.
func NormalizePolicy(p Policy) Policy {
	p.Mode = strings.ToLower(strings.TrimSpace(p.Mode))
	if p.Mode != "scroll" && p.Mode != "typer" {
		p.Mode = "scroll"
	}
	if p.TypeSpeed < 1 {
		p.TypeSpeed = 1
	}
	if p.TypeSpeed > 1000 {
		p.TypeSpeed = 1000
	}
	// ScrollSpeed: [1, 1000]
	if p.ScrollSpeed < 1 {
		p.ScrollSpeed = 1
	}
	if p.ScrollSpeed > 1000 {
		p.ScrollSpeed = 1000
	}

	// MaxLines: [1, 1000]
	if p.MaxLines < 1 {
		p.MaxLines = 1
	}
	if p.MaxLines > 1000 {
		p.MaxLines = 1000
	}

	// ScanMS: [100, 30000]
	if p.ScanMS < 100 {
		p.ScanMS = 100
	}
	if p.ScanMS > 30000 {
		p.ScanMS = 30000
	}

	// FadeOutMS: clamp to positive
	if p.FadeOutMS < 0 {
		p.FadeOutMS = 0
	}

	// FadeInMS: clamp to positive
	if p.FadeInMS < 0 {
		p.FadeInMS = 0
	}

	// LineTimeMS: clamp to positive
	if p.LineTimeMS < 1 {
		p.LineTimeMS = 1
	}

	// MaxWaitS: clamp to positive
	if p.MaxWaitS < 1 {
		p.MaxWaitS = 1
	}

	return p
}
