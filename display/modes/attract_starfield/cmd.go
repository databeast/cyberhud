package attract_starfield

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

// intRangeValidator returns a KeyValidator that accepts integer values in [min, max].
func intRangeValidator(min, max int) cmdutil.KeyValidator {
	return func(value string) string {
		v := strings.TrimSpace(value)
		n, err := strconv.Atoi(v)
		if err != nil {
			return "must be an integer"
		}
		if n < min || n > max {
			return fmt.Sprintf("must be in [%d, %d]", min, max)
		}
		return ""
	}
}

// cmdHandler is the declarative CmdHandler for the "attract_starfield" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_starfield",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "layers", Validate: intRangeValidator(1, 8)},
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
	case "layers":
		return fmt.Sprintf("%d", p.Layers)
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
	case "layers":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.Layers = n
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_starfield" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_starfield",
		Title:   "Starfield",
		Summary: "Perspective depth starfield effect with stars emanating from a central vanishing point.",
		Order:   200,
		Options: []catalog.OptionDefinition{
			{Key: "speed", Type: "float", Summary: "Travel speed multiplier (0.1-10.0).", Default: "1.0", Allowed: []string{}},
			{Key: "density", Type: "float", Summary: "Star density scaling (0.1-1.0).", Default: "0.5", Allowed: []string{}},
			{Key: "layers", Type: "int", Summary: "Number of depth layers (1-8).", Default: "4", Allowed: []string{}},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_starfield",
		Summary: "Query or set starfield display options.",
		Usage:   "attract_starfield [speed=<0.1-10.0>] [density=<0.1-1.0>] [layers=<1-8>]",
		Handle:  HandleCommand,
	})
}
