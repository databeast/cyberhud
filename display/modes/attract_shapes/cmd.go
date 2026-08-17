package attract_shapes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_shapes/source"
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

// cmdHandler is the declarative CmdHandler for the "attract_shapes" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_shapes",
	Keys: []cmdutil.KeyDef{
		{Name: "speed", Validate: floatRangeValidator(0.1, 10.0)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "shape_count", Validate: intRangeValidator(1, 50)},
		{Name: "pulse_rate", Validate: floatRangeValidator(0.1, 5.0)},
		{Name: "complexity", Validate: intRangeValidator(3, 8)},
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
	case "shape_count":
		return fmt.Sprintf("%d", p.ShapeCount)
	case "pulse_rate":
		return fmt.Sprintf("%g", p.PulseRate)
	case "complexity":
		return fmt.Sprintf("%d", p.Complexity)
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
	case "shape_count":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.ShapeCount = n
		}
	case "pulse_rate":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.PulseRate = f
		}
	case "complexity":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.Complexity = n
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}

// HandleCommand is the catalog command handler for the "attract_shapes" verb.
// It delegates to the declarative CmdHandler which handles atomic validation
// and application of key=value arguments.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_shapes",
		Title:   "Shapes",
		Summary: "Pulsing and rotating geometric shapes with configurable complexity and speed.",
		Order:   200,
		Options: source.Policy{}.Options(),
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_shapes",
		Summary: "Query or set shapes display options.",
		Usage:   "attract_shapes [speed=<0.1-10.0>] [density=<0.1-1.0>] [shape_count=<1-50>] [pulse_rate=<0.1-5.0>] [complexity=<3-8>]",
		Handle:  HandleCommand,
	})
}
