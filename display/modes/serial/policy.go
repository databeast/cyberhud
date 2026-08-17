package serial

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/modes/serial/styles"
	"github.com/databeast/cyberhud/display/style"
)

type Policy = source.Policy
type Snapshot = source.Snapshot

const StyleDefault = source.StyleDefault

// serialRegistry is the per-mode StyleRegistry for the serial display mode.
var serialRegistry = style.NewRegistry[source.Snapshot, source.Policy](
	styles.DefaultStyle,
	styles.RawStyle,
	styles.DashboardStyle,
	styles.CompactStyle,
	styles.FramedStyle,
	styles.MonoSlow64x128Style,
	styles.MonoSlow128x64Style,
	styles.MonoSlow800x480Style,
	styles.MonoSlow104x212Style,
	styles.MonoSlow122x250Style,
	styles.MonoSlow176x264Style,
	styles.MonoSlow200x200Style,
	styles.MonoSlow212x104Style,
	styles.MonoSlow250x122Style,
	styles.MonoSlow264x176Style,
	styles.MonoSlow296x128Style,
	styles.MonoSlow300x400Style,
	styles.MonoSlow400x300Style,
	styles.MonoSlow480x800Style,
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast128x128Style,
	styles.MonoFast160x80Style,
	styles.MonoFast240x135Style,
	styles.MonoFast240x240Style,
	styles.MonoFast320x240Style,
	styles.MonoFast480x320Style,
	styles.MonoFast800x480Style,
	styles.GrayscaleSlow160x80Style,
	styles.GrayscaleSlow240x135Style,
	styles.GrayscaleSlow240x240Style,
	styles.GrayscaleSlow320x240Style,
	styles.GrayscaleSlow480x320Style,
	styles.GrayscaleSlow800x480Style,
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast800x480Style,
	styles.ColorSlow122x250Style,
	styles.ColorSlow176x264Style,
	styles.ColorSlow200x200Style,
	styles.ColorSlow400x300Style,
	styles.ColorSlow800x480Style,
	styles.ColorFast128x128Style,
	styles.ColorFast160x80Style,
	styles.ColorFast240x135Style,
	styles.ColorFast240x240Style,
	styles.ColorFast240x320Style,
	styles.ColorFast320x240Style,
	styles.ColorFast320x480Style,
	styles.ColorFast480x320Style,
	styles.ColorFast480x800Style,
	styles.ColorFast800x480Style,
)

// registeredStyleNames returns the list of style names from the registry.
func registeredStyleNames() []string {
	styles := serialRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// normalizePolicy ensures policy fields contain valid values and registered styles.
func normalizePolicy(p Policy) Policy {
	p = source.NormalizePolicy(p)
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && serialRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// serialSnapshotter implements catalog.PolicySnapshotter for the serial mode.
type serialSnapshotter struct{}

// SnapshotPolicy returns the current serial policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (serialSnapshotter) SnapshotPolicy() map[string]interface{} {
	return PolicySnapshot().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (serialSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["port"]; ok {
		if s, ok := v.(string); ok {
			p.Port = s
		}
	}
	if v, ok := data["baud"]; ok {
		if n, ok := toInt(v); ok {
			p.Baud = n
		}
	}
	if v, ok := data["lines"]; ok {
		if n, ok := toInt(v); ok {
			p.MaxLines = n
		}
	}
	if v, ok := data["autoselect"]; ok {
		if b, ok := toBool(v); ok {
			p.AutoSelect = b
		}
	}
	if v, ok := data["scan_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.ScanMS = n
		}
	}
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

// toBool extracts a bool from an interface value.
func toBool(v interface{}) (bool, bool) {
	b, ok := v.(bool)
	return b, ok
}

func init() {
	catalog.RegisterSnapshotter("serial", serialSnapshotter{})
}

// PolicySnapshot returns the current serial monitor policy.
func PolicySnapshot() Policy {
	return source.PolicySnapshot()
}

// SetPolicy updates the serial monitor policy.
func SetPolicy(policy Policy) {
	source.SetPolicy(normalizePolicy(policy))
}
