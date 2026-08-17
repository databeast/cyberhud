package attract_waveform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
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

// cmdHandler is the declarative CmdHandler for the "attract_waveform" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_waveform",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "amplitude", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "traces", Validate: intRangeValidator(1, 8)},
		{Name: "persistence", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "direction", Validate: cmdutil.AllowedValidator([]string{"auto", "horizontal", "vertical", "h", "v"})},
	},
	Get:       getPolicyValue,
	Apply:     applyPolicyValue,
	PostApply: nil,
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "speed":
		return fmt.Sprintf("%g", p.Speed)
	case "density":
		return fmt.Sprintf("%g", p.Density)
	case "amplitude":
		return fmt.Sprintf("%g", p.Amplitude)
	case "traces":
		return fmt.Sprintf("%d", p.Traces)
	case "persistence":
		return fmt.Sprintf("%g", p.Persistence)
	case "direction":
		return p.Direction
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
	case "density":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Density = f
		}
	case "amplitude":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Amplitude = f
		}
	case "traces":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.Traces = n
		}
	case "persistence":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Persistence = f
		}
	case "direction":
		policyState.policy.Direction = strings.TrimSpace(value)
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_waveform" verb.
// It delegates to the declarative CmdHandler which handles atomic validation
// and application of key=value arguments.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_waveform",
		Title:   "Waveform",
		Summary: "Oscilloscope-style animated waveform traces with phosphor persistence.",
		Order:   200,
		Options: source.Policy{}.Options(),
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_waveform",
		Summary: "Query or set waveform display options.",
		Usage:   "attract_waveform [speed=<0.1-10.0>] [density=<0.1-1.0>] [amplitude=<0.1-1.0>] [traces=<1-8>] [persistence=<0.1-1.0>]",
		Handle:  HandleCommand,
	})
}
