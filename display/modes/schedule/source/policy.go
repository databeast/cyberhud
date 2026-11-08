package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy controls runtime behavior for the schedule mode.
//
// Source is the calendar location: an http(s):// URL or a local filesystem path
// to an iCalendar (.ics) file. RefreshSec is how often the background fetcher
// re-reads it. The remaining fields are presentation knobs owned here so every
// panel size honours them uniformly.
type Policy struct {
	Source      string // http(s) URL or local path to an .ics file; empty disables fetching
	RefreshSec  int    // Seconds between calendar re-fetches
	TimeFormat  string // "24h" or "12h"
	ShowUpcoming bool  // Show the next upcoming event beneath the active one
	Style       string // empty means auto-detect via style fitness
	Font        string // "auto" or a specific font ID
}

const (
	DefaultRefreshSec = 300
	MinRefreshSec     = 30
)

// AllowedTimeFormats lists the accepted time_format values.
var AllowedTimeFormats = []string{"24h", "12h"}

// DefaultPolicy returns baseline schedule behavior.
func DefaultPolicy() Policy {
	return Policy{
		Source:       "",
		RefreshSec:   DefaultRefreshSec,
		TimeFormat:   "24h",
		ShowUpcoming: true,
		Style:        "",
		Font:         "auto",
	}
}

// Options publishes the configurable options for the schedule mode.
func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "source", Type: "string", Summary: "Calendar location: an http(s) URL or a local path to an .ics file.", Default: ""},
		{Key: "refresh_sec", Type: "int", Summary: "Seconds between calendar re-fetches.", Default: "300"},
		{Key: "time_format", Type: "string", Summary: "Event time display format.", Default: "24h", Allowed: AllowedTimeFormats},
		{Key: "show_upcoming", Type: "bool", Summary: "Show the next upcoming event beneath the active one.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
	}
}

// Fingerprint returns a stable token over all policy fields, used for change
// detection in the render cache key.
func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%d|%s|%t|%s|%s", p.Source, p.RefreshSec, p.TimeFormat, p.ShowUpcoming, p.Style, p.Font)
}

// ToMap renders the policy as a JSON-serializable map using snake_case wire keys.
func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"source":        p.Source,
		"refresh_sec":   p.RefreshSec,
		"time_format":   p.TimeFormat,
		"show_upcoming": p.ShowUpcoming,
		"style":         p.Style,
		"font":          p.Font,
	}
}

// NormalizePolicy applies source-owned, registry-independent policy defaults.
func NormalizePolicy(p Policy) Policy {
	d := DefaultPolicy()
	p.Source = strings.TrimSpace(p.Source)
	if p.RefreshSec < MinRefreshSec {
		if p.RefreshSec <= 0 {
			p.RefreshSec = d.RefreshSec
		} else {
			p.RefreshSec = MinRefreshSec
		}
	}
	p.TimeFormat = strings.ToLower(strings.TrimSpace(p.TimeFormat))
	if p.TimeFormat != "12h" {
		p.TimeFormat = "24h"
	}
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.Font = strings.TrimSpace(p.Font)
	if p.Font == "" {
		p.Font = "auto"
	}
	return p
}
