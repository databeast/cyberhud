package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the storage mode.
type Policy struct {
	Style     string // resolution-qualified style name, or "" for fitness/alias resolution
	Font      string // "auto" or a registered font ID
	RefreshMS int    // sampling interval in milliseconds [500, 120000]
	Unit      string // "auto", "si" (1000-based) or "iec" (1024-based)
	ShowAll   bool   // include pseudo/overlay filesystems (tmpfs, overlay, ...)
}

// AllowedUnits lists valid capacity-unit display modes.
//
//	"auto" — size the unit per value (B/KB/MB/GB/TB), 1024-based
//	"si"   — 1000-based units (KB = 1000 bytes)
//	"iec"  — 1024-based units (KiB = 1024 bytes)
var AllowedUnits = []string{"auto", "si", "iec"}

// DefaultRefreshMS is the default sampling interval. Storage changes slowly, so
// the default is deliberately relaxed relative to CPU/thermal modes.
const DefaultRefreshMS = 5000

// DefaultPolicy returns the default storage policy.
func DefaultPolicy() Policy {
	return Policy{
		Style:     "",
		Font:      "auto",
		RefreshMS: DefaultRefreshMS,
		Unit:      "auto",
		ShowAll:   false,
	}
}

// Options returns catalog option definitions for the storage mode policy.
func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		{Key: "refresh_ms", Type: "int", Summary: "Sampling interval in milliseconds.", Default: "5000"},
		{Key: "unit", Type: "string", Summary: "Capacity unit style.", Default: "auto", Allowed: AllowedUnits},
		{Key: "show_all", Type: "bool", Summary: "Include pseudo/virtual filesystems (tmpfs, overlay).", Default: "false", Allowed: []string{"true", "false"}},
	}
}

// Fingerprint encodes all policy fields into a deterministic string.
func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%d|%s|%v", p.Style, p.Font, p.RefreshMS, p.Unit, p.ShowAll)
}

// ToMap returns the policy as a JSON-serializable map using snake_case keys.
func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":      p.Style,
		"font":       p.Font,
		"refresh_ms": p.RefreshMS,
		"unit":       p.Unit,
		"show_all":   p.ShowAll,
	}
}

// NormalizePolicy performs registry-independent policy normalization.
func NormalizePolicy(p Policy) Policy {
	d := DefaultPolicy()

	p.Style = strings.ToLower(strings.TrimSpace(p.Style))

	p.Font = strings.TrimSpace(p.Font)
	if p.Font == "" {
		p.Font = d.Font
	}

	p.RefreshMS = ClampRefreshMS(p.RefreshMS)

	p.Unit = strings.ToLower(strings.TrimSpace(p.Unit))
	if !isAllowed(p.Unit, AllowedUnits) {
		p.Unit = d.Unit
	}

	return p
}

// ClampRefreshMS clamps a RefreshMS value to [500, 120000], defaulting to
// DefaultRefreshMS when zero or negative.
func ClampRefreshMS(ms int) int {
	if ms <= 0 {
		ms = DefaultRefreshMS
	}
	if ms < 500 {
		ms = 500
	}
	if ms > 120000 {
		ms = 120000
	}
	return ms
}

func isAllowed(value string, allowed []string) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}
