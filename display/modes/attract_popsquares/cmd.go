package attract_popsquares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_popsquares/source"
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

// cmdHandler is the declarative CmdHandler for the "attract_popsquares" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_popsquares",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "grid_scale", Validate: floatRangeValidator(0.5, 4.0)},
		{Name: "pop_rate", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "warmth", Validate: floatRangeValidator(0.0, 1.0)},
	},
	Get:       getPolicyValue,
	Apply:     applyPolicyValue,
	PostApply: fitnessNotesPostApply,
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "speed":
		return fmt.Sprintf("%g", p.Speed)
	case "grid_scale":
		return fmt.Sprintf("%g", p.GridScale)
	case "pop_rate":
		return fmt.Sprintf("%g", p.PopRate)
	case "warmth":
		return fmt.Sprintf("%g", p.Warmth)
	default:
		return ""
	}
}

// applyPolicyValue updates a single policy key with its validated value.
func applyPolicyValue(key, value string) {
	policyState.Lock()
	defer policyState.Unlock()
	switch key {
	case "speed":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Speed = f
		}
	case "grid_scale":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.GridScale = f
		}
	case "pop_rate":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.PopRate = f
		}
	case "warmth":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Warmth = f
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_popsquares" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_popsquares",
		Title:   "Pop Squares",
		Summary: "A calm grid of color cells that gently pulse and fade — a supercomputer idling.",
		Order:   200,
		Options: source.Policy{}.Options(),
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_popsquares",
		Summary: "Query or set pop-squares display options.",
		Usage:   "attract_popsquares [speed=<0.1-10.0>] [grid_scale=<0.5-4.0>] [pop_rate=<0.1-1.0>] [warmth=<0.0-1.0>]",
		Handle:  HandleCommand,
	})
}
