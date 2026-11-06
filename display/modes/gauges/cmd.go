package gauges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/gauges/source"
	"github.com/databeast/cyberhud/display/style/color"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:       "gauges",
		Title:    "Gauges",
		Category: "Data",
		Scope:    "any",
		Summary:  "User-supplied progress bars, rings, and dials rendered as panel widgets.",
		Order:    85,
		Options: append(source.Policy{}.Options(), catalog.OptionDefinition{
			Key:     "style",
			Type:    "string",
			Summary: "Visual style name or empty for automatic fitness selection.",
			Default: "",
			Allowed: gaugesRegistry.RegisteredNames(),
		}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "gauges",
		Summary: "Manage user-supplied gauge data and display policy.",
		Usage:   "gauges set <payload> | gauges get | gauges [style=<name>] [shape=<auto|linear|ring|arc|pie>] [show_labels=<true|false>] [label_tier=<auto|small|normal|large|fullsize>] [accent=<name|none>] [default_min=<float>] [default_max=<float>] [columns=<int>] [rows=<int>] [tile_gap_px=<int>] [padding_pct=<int>]",
		Handle:  HandleConsoleCommand,
	})
}

// cmdHandler is the uniform CmdHandler for the gauges mode's policy options.
// Policy edits flow through it so applying key=value pairs fires
// cmdutil.OnApplied, driving centralized PolicyStore capture (freeze policy
// persistence). The data sub-verbs (set/get) are handled by HandleConsoleCommand
// before delegating here.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "gauges",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: func(v string) string {
			if v != "" && gaugesRegistry.Lookup(strings.ToLower(strings.TrimSpace(v))) == nil {
				return fmt.Sprintf("must be one of [%s]", strings.Join(gaugesRegistry.RegisteredNames(), ", "))
			}
			return ""
		}},
		{Name: "shape", Validate: func(v string) string {
			if !source.ValidShape(v) {
				return fmt.Sprintf("must be one of [%s]", strings.Join(source.AllowedShapes, ", "))
			}
			return ""
		}},
		{Name: "show_labels", Validate: cmdutil.BoolValidator()},
		{Name: "label_tier", Validate: func(v string) string {
			if !source.ValidLabelTier(v) {
				return fmt.Sprintf("must be one of [%s]", strings.Join(source.AllowedLabelTiers, ", "))
			}
			return ""
		}},
		{Name: "accent", Validate: func(v string) string {
			if !source.ValidAccent(v) {
				return fmt.Sprintf("must be one of [%s]", strings.Join(append(color.Names(), "none"), ", "))
			}
			return ""
		}},
		{Name: "default_min", Validate: floatValidator},
		{Name: "default_max", Validate: floatValidator},
		{Name: "columns", Validate: cmdutil.IntValidator(0)},
		{Name: "rows", Validate: cmdutil.IntValidator(0)},
		{Name: "tile_gap_px", Validate: cmdutil.IntValidator(0)},
		{Name: "padding_pct", Validate: cmdutil.IntValidator(0)},
	},
	Get:   getPolicyValue,
	Apply: applyPolicyValue,
}

// floatValidator accepts any value parseable as a float64.
func floatValidator(v string) string {
	if _, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err != nil {
		return fmt.Sprintf("must be a number, got %q", v)
	}
	return ""
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := sourcePolicySnapshot()
	switch key {
	case "style":
		return p.Style
	case "shape":
		return p.Shape
	case "show_labels":
		return strconv.FormatBool(p.ShowLabels)
	case "label_tier":
		return p.LabelTier
	case "accent":
		return p.Accent
	case "default_min":
		return strconv.FormatFloat(p.DefaultMin, 'g', -1, 64)
	case "default_max":
		return strconv.FormatFloat(p.DefaultMax, 'g', -1, 64)
	case "columns":
		return strconv.Itoa(p.Columns)
	case "rows":
		return strconv.Itoa(p.Rows)
	case "tile_gap_px":
		return strconv.Itoa(p.TileGapPx)
	case "padding_pct":
		return strconv.Itoa(p.PaddingPct)
	default:
		return ""
	}
}

// applyPolicyValue sets a single policy key on the current policy and persists
// it via SetPolicy (which normalizes). Values have already passed validation.
func applyPolicyValue(key, value string) {
	p := sourcePolicySnapshot()
	value = strings.TrimSpace(value)
	switch key {
	case "style":
		p.Style = strings.ToLower(value)
	case "shape":
		p.Shape = strings.ToLower(value)
	case "show_labels":
		if b, ok := cmdutil.ParseBool(value); ok {
			p.ShowLabels = b
		}
	case "label_tier":
		p.LabelTier = strings.ToLower(value)
	case "accent":
		p.Accent = strings.ToLower(value)
	case "default_min":
		if n, err := strconv.ParseFloat(value, 64); err == nil {
			p.DefaultMin = n
		}
	case "default_max":
		if n, err := strconv.ParseFloat(value, 64); err == nil {
			p.DefaultMax = n
		}
	case "columns":
		if n, err := strconv.Atoi(value); err == nil {
			p.Columns = n
		}
	case "rows":
		if n, err := strconv.Atoi(value); err == nil {
			p.Rows = n
		}
	case "tile_gap_px":
		if n, err := strconv.Atoi(value); err == nil {
			p.TileGapPx = n
		}
	case "padding_pct":
		if n, err := strconv.Atoi(value); err == nil {
			p.PaddingPct = n
		}
	}
	SetPolicy(p)
}

// HandleConsoleCommand handles the top-level gauges command. The data sub-verbs
// "set" and "get" manage the externally-supplied gauge payload; everything else
// is treated as policy key=value edits and delegated to the uniform
// cmdutil.CmdHandler. A bare "gauges" with no args returns the policy query line.
func HandleConsoleCommand(args []string) string {
	if len(args) == 0 {
		return cmdHandler.Handle(args)
	}

	switch strings.ToLower(strings.TrimSpace(args[0])) {
	case "set":
		if len(args) == 1 {
			SetSnapshot(source.GaugeSet{})
			return "OK gauges count=0"
		}
		payload := strings.TrimSpace(strings.Join(args[1:], " "))
		if payload == "" {
			SetSnapshot(source.GaugeSet{})
			return "OK gauges count=0"
		}
		snap, err := source.ParsePayload(payload, sourcePolicySnapshot())
		if err != nil {
			return fmt.Sprintf("ERR %s", err.Error())
		}
		SetSnapshot(snap)
		return fmt.Sprintf("OK gauges count=%d", len(snap.Gauges))
	case "get":
		data, err := source.SerializeSnapshot(sourceSnapshot())
		if err != nil {
			return fmt.Sprintf("ERR %s", err.Error())
		}
		return "OK\n" + data
	default:
		// Policy edits: style=... shape=... etc. — uniform key=value path.
		return cmdHandler.Handle(args)
	}
}
