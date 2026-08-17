package attract_bokeh

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
)

// floatRangeValidator returns a KeyValidator that accepts float values in [min, max].
func floatRangeValidator(min, max float64) cmdutil.KeyValidator {
	return func(value string) string {
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return "must be a number"
		}
		if f < min || f > max {
			return fmt.Sprintf("must be in [%.1f, %.1f]", min, max)
		}
		return ""
	}
}

// cmdHandler is the declarative CmdHandler for the "attract_bokeh" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_bokeh",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "size_variance", Validate: floatRangeValidator(0.0, 1.0)},
		{Name: "saturation", Validate: floatRangeValidator(0.0, 1.0)},
	},
	Get:       getPolicy,
	Apply:     applyPolicy,
	PostApply: nil,
}

// getPolicy returns the current value for a given policy key.
func getPolicy(key string) string {
	p := GetPolicy()
	switch key {
	case "speed":
		return fmt.Sprintf("%g", p.Speed)
	case "density":
		return fmt.Sprintf("%g", p.Density)
	case "size_variance":
		return fmt.Sprintf("%g", p.SizeVariance)
	case "saturation":
		return fmt.Sprintf("%g", p.Saturation)
	default:
		return ""
	}
}

// applyPolicy updates a single policy key with its validated value.
func applyPolicy(key, value string) {
	policyState.Lock()
	defer policyState.Unlock()
	switch key {
	case "speed":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Speed = f
		}
	case "density":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Density = f
		}
	case "size_variance":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.SizeVariance = f
		}
	case "saturation":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Saturation = f
		}
	}
	policyState.policy = source.NormalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_bokeh" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_bokeh",
		Title:   "Bokeh",
		Summary: "Soft out-of-focus light circles drifting with gentle color variation.",
		Order:   200,
		Options: []catalog.OptionDefinition{
			{Key: "speed", Type: "float", Summary: "Drift speed multiplier.", Default: "1.0", Allowed: []string{}},
			{Key: "density", Type: "float", Summary: "Circle count scaling from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
			{Key: "size_variance", Type: "float", Summary: "Radius spread, 0.0 uniform to 1.0 maximum.", Default: "0.5", Allowed: []string{}},
			{Key: "saturation", Type: "float", Summary: "Color saturation from 0.0 to 1.0.", Default: "0.7", Allowed: []string{}},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_bokeh",
		Summary: "Query or set bokeh display options.",
		Usage:   "attract_bokeh [speed=<0.1-10.0>] [density=<0.1-1.0>] [size_variance=<0.0-1.0>] [saturation=<0.0-1.0>]",
		Handle:  HandleCommand,
	})
}
