package gauges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/gauges/source"
	"github.com/databeast/cyberhud/display/style/color"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "gauges",
		Title:   "Gauges",
		Scope:   "any",
		Summary: "User-supplied progress bars, rings, and dials rendered as panel widgets.",
		Order:   85,
		Options: append(source.Policy{}.Options(), catalog.OptionDefinition{
			Key:     "style",
			Type:    "string",
			Summary: "Visual style name or empty for automatic fitness selection.",
			Default: "",
			Allowed: registeredStyleNames(),
		}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "gauges",
		Summary: "Manage user-supplied gauge data and display policy.",
		Usage:   "gauges <set|get|policy> ...",
		Handle:  HandleConsoleCommand,
	})
}

// HandleConsoleCommand handles the top-level gauges command.
func HandleConsoleCommand(args []string) string {
	if len(args) == 0 {
		return "ERR usage: gauges <set|get|policy> ..."
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
	case "policy":
		if len(args) == 1 {
			return source.FormatPolicyResponse(sourcePolicySnapshot())
		}
		next := sourcePolicySnapshot()
		for _, token := range args[1:] {
			kv := strings.SplitN(token, "=", 2)
			if len(kv) != 2 {
				return "ERR usage: gauges policy [style=<name>] [shape=<auto|linear|ring|arc|pie>] [show_labels=<true|false>] [label_tier=<auto|small|normal|large|fullsize>] [accent=<name|none>] [default_min=<float>] [default_max=<float>] [columns=<int>] [rows=<int>] [tile_gap_px=<int>] [padding_pct=<int>]"
			}
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			val := strings.TrimSpace(kv[1])
			switch key {
			case "style":
				if val != "" && gaugesRegistry.Lookup(val) == nil {
					return fmt.Sprintf("ERR style: must be one of [%s]", strings.Join(registeredStyleNames(), ", "))
				}
				next.Style = strings.ToLower(val)
			case "shape":
				if !source.ValidShape(val) {
					return fmt.Sprintf("ERR shape: must be one of [%s]", strings.Join(source.AllowedShapes, ", "))
				}
				next.Shape = strings.ToLower(val)
			case "show_labels":
				b, ok := parseBool(val)
				if !ok {
					return "ERR show_labels: must be true or false"
				}
				next.ShowLabels = b
			case "label_tier":
				if !source.ValidLabelTier(val) {
					return fmt.Sprintf("ERR label_tier: must be one of [%s]", strings.Join(source.AllowedLabelTiers, ", "))
				}
				next.LabelTier = strings.ToLower(val)
			case "accent":
				if !source.ValidAccent(val) {
					return fmt.Sprintf("ERR accent: must be one of [%s]", strings.Join(append(color.Names(), "none"), ", "))
				}
				next.Accent = strings.ToLower(val)
			case "default_min":
				n, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Sprintf("ERR default_min must be a number, got %q", val)
				}
				next.DefaultMin = n
			case "default_max":
				n, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Sprintf("ERR default_max must be a number, got %q", val)
				}
				next.DefaultMax = n
			case "columns":
				n, err := strconv.Atoi(val)
				if err != nil || n < 0 {
					return fmt.Sprintf("ERR columns must be an integer >= 0, got %q", val)
				}
				next.Columns = n
			case "rows":
				n, err := strconv.Atoi(val)
				if err != nil || n < 0 {
					return fmt.Sprintf("ERR rows must be an integer >= 0, got %q", val)
				}
				next.Rows = n
			case "tile_gap_px":
				n, err := strconv.Atoi(val)
				if err != nil || n < 0 {
					return fmt.Sprintf("ERR tile_gap_px must be an integer >= 0, got %q", val)
				}
				next.TileGapPx = n
			case "padding_pct":
				n, err := strconv.Atoi(val)
				if err != nil || n < 0 {
					return fmt.Sprintf("ERR padding_pct must be an integer >= 0, got %q", val)
				}
				next.PaddingPct = n
			default:
				return fmt.Sprintf("ERR unknown gauges policy key %q", key)
			}
		}
		SetPolicy(next)
		return source.FormatPolicyResponse(sourcePolicySnapshot())
	default:
		return "ERR usage: gauges <set|get|policy> ..."
	}
}

func parseBool(s string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "on", "1":
		return true, true
	case "false", "no", "off", "0":
		return false, true
	default:
		return false, false
	}
}
