package attract_plasma

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

// cmdHandler is the declarative CmdHandler for the "attract_plasma" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_plasma",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "cycle_rate", Validate: floatRangeValidator(0.1, 5.0)},
		{Name: "blob_scale", Validate: floatRangeValidator(0.5, 4.0)},
	},
	Get:       getPolicy,
	Apply:     applyPolicy,
	PostApply: fitnessNotesPostApply,
}

// getPolicy returns the current value for a given policy key.
func getPolicy(key string) string {
	p := GetPolicy()
	switch key {
	case "speed":
		return fmt.Sprintf("%g", p.Speed)
	case "density":
		return fmt.Sprintf("%g", p.Density)
	case "cycle_rate":
		return fmt.Sprintf("%g", p.CycleRate)
	case "blob_scale":
		return fmt.Sprintf("%g", p.BlobScale)
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
	case "cycle_rate":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.CycleRate = f
		}
	case "blob_scale":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.BlobScale = f
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_plasma" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_plasma",
		Title:   "Plasma",
		Summary: "Lava-lamp plasma blob effect with smooth morphing color gradients and configurable spatial frequency.",
		Order:   200,
		Options: []catalog.OptionDefinition{
			{Key: "speed", Type: "float", Summary: "Morph speed multiplier.", Default: "1.0", Allowed: []string{}},
			{Key: "density", Type: "float", Summary: "Unused, kept for uniformity.", Default: "0.5", Allowed: []string{}},
			{Key: "cycle_rate", Type: "float", Summary: "Color palette cycles per second.", Default: "1.0", Allowed: []string{}},
			{Key: "blob_scale", Type: "float", Summary: "Spatial frequency multiplier for blob size.", Default: "1.0", Allowed: []string{}},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_plasma",
		Summary: "Query or set plasma display options.",
		Usage:   "attract_plasma [speed=<0.1-10.0>] [density=<0.1-1.0>] [cycle_rate=<0.1-5.0>] [blob_scale=<0.5-4.0>]",
		Handle:  HandleCommand,
	})
}
