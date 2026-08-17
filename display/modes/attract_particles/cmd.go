package attract_particles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_particles/source"
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

// cmdHandler is the declarative CmdHandler for the "attract_particles" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_particles",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "drift", Validate: floatRangeValidator(0.0, 1.0)},
		{Name: "glow", Validate: floatRangeValidator(0.1, 1.0)},
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
	case "drift":
		return fmt.Sprintf("%g", p.Drift)
	case "glow":
		return fmt.Sprintf("%g", p.Glow)
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
	case "drift":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Drift = f
		}
	case "glow":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Glow = f
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_particles" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_particles",
		Title:   "Particles",
		Summary: "Drifting firefly-like particles with color cycling and edge-wrapping motion.",
		Order:   200,
		Options: source.Policy{}.Options(),
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_particles",
		Summary: "Query or set particle display options.",
		Usage:   "attract_particles [speed=<0.1-10.0>] [density=<0.1-1.0>] [drift=<0.0-1.0>] [glow=<0.1-1.0>]",
		Handle:  HandleCommand,
	})
}
