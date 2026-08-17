package attract_geometric

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
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

// cmdHandler is the declarative CmdHandler for the "attract_geometric" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_geometric",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "glow_intensity", Validate: floatRangeValidator(0.0, 1.0)},
		{Name: "fragment_rate", Validate: floatRangeValidator(0.0, 2.0)},
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
	case "glow_intensity":
		return fmt.Sprintf("%g", p.GlowIntensity)
	case "fragment_rate":
		return fmt.Sprintf("%g", p.FragmentRate)
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
	case "glow_intensity":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.GlowIntensity = f
		}
	case "fragment_rate":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.FragmentRate = f
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_geometric" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_geometric",
		Title:   "Geometric",
		Summary: "Rotating semi-transparent rectangles with sinusoidal fades and drifting pseudocode fragments.",
		Order:   200,
		Options: []catalog.OptionDefinition{
			{Key: "speed", Type: "float", Summary: "Animation speed multiplier.", Default: "1.0", Allowed: []string{}},
			{Key: "density", Type: "float", Summary: "Cluster count scaling from 0.1 to 1.0.", Default: "0.5", Allowed: []string{}},
			{Key: "glow_intensity", Type: "float", Summary: "Glow effect strength from 0.0 to 1.0.", Default: "1.0", Allowed: []string{}},
			{Key: "fragment_rate", Type: "float", Summary: "Pseudocode fragment spawn rate from 0.0 to 2.0.", Default: "1.0", Allowed: []string{}},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_geometric",
		Summary: "Query or set geometric display options.",
		Usage:   "attract_geometric [speed=<0.1-10.0>] [density=<0.1-1.0>] [glow_intensity=<0.0-1.0>] [fragment_rate=<0.0-2.0>]",
		Handle:  HandleCommand,
	})
}
