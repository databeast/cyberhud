package gauges

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/gauges/source"
)

// gaugesSnapshotter implements catalog.PolicySnapshotter for the gauges mode so
// that policy edits made via the uniform CmdHandler path are captured by the
// daemon's PolicyStore (and thus persisted by `freeze policy`). Without a
// registered snapshotter the catalog OnPolicyChange bridge no-ops, so this type
// is required for gauges persistence to work.
type gaugesSnapshotter struct{}

// SnapshotPolicy returns the current gauges policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format and the option keys.
func (gaugesSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := sourcePolicySnapshot()
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

// RestorePolicy applies policy values from a JSON map, running them through the
// same normalization as a live edit (SetPolicy). Keys use the protocol wire
// format. Numeric values decode from JSON as float64.
func (gaugesSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"].(string); ok {
		p.Style = v
	}
	if v, ok := data["shape"].(string); ok {
		p.Shape = v
	}
	if v, ok := data["show_labels"].(bool); ok {
		p.ShowLabels = v
	}
	if v, ok := data["label_tier"].(string); ok {
		p.LabelTier = v
	}
	if v, ok := data["accent"].(string); ok {
		p.Accent = v
	}
	if v, ok := cmdutil.ToFloat64(data["default_min"]); ok {
		p.DefaultMin = v
	}
	if v, ok := cmdutil.ToFloat64(data["default_max"]); ok {
		p.DefaultMax = v
	}
	if v, ok := cmdutil.ToInt(data["columns"]); ok {
		p.Columns = v
	}
	if v, ok := cmdutil.ToInt(data["rows"]); ok {
		p.Rows = v
	}
	if v, ok := cmdutil.ToInt(data["tile_gap_px"]); ok {
		p.TileGapPx = v
	}
	if v, ok := cmdutil.ToInt(data["padding_pct"]); ok {
		p.PaddingPct = v
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("gauges", gaugesSnapshotter{})
}
