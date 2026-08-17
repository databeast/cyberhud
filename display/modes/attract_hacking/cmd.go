package attract_hacking

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_hacking/source"
)

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

var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_hacking",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 3.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "glitch", Validate: floatRangeValidator(0.0, 1.0)},
		{Name: "intensity", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "pulse", Validate: floatRangeValidator(0.1, 1.5)},
	},
	Get:   getPolicyValue,
	Apply: applyPolicyValue,
}

func getPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "speed":
		return fmt.Sprintf("%g", p.Speed)
	case "density":
		return fmt.Sprintf("%g", p.Density)
	case "glitch":
		return fmt.Sprintf("%g", p.Glitch)
	case "intensity":
		return fmt.Sprintf("%g", p.Intensity)
	case "pulse":
		return fmt.Sprintf("%g", p.Pulse)
	default:
		return ""
	}
}

func applyPolicyValue(key, value string) {
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
	case "glitch":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Glitch = f
		}
	case "intensity":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Intensity = f
		}
	case "pulse":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Pulse = f
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_hacking",
		Title:   "Hacking",
		Summary: "Cinematic cyberpunk terminal intrusion sequence with neon overlays and fake-system logs.",
		Order:   220,
		Options: source.Policy{}.Options(),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_hacking",
		Summary: "Query or set hacking display options.",
		Usage:   "attract_hacking [speed=<0.1-3.0>] [density=<0.1-1.0>] [glitch=<0.0-1.0>] [intensity=<0.1-1.0>] [pulse=<0.1-1.5>]",
		Handle:  HandleCommand,
	})
}
