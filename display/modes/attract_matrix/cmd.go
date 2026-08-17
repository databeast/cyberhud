package attract_matrix

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// cmdHandler is the declarative CmdHandler for the "attract_matrix" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "attract_matrix",
	Keys: []cmdutil.KeyDef{
		{Name: "min_speed", Validate: floatPositiveValidator()},
		{Name: "max_speed", Validate: floatPositiveValidator()},
		{Name: "trail_length", Validate: intRangeValidator(1, 128)},
		{Name: "density", Validate: floatRangeValidator(0.1, 1.0)},
		{Name: "show_background", Validate: cmdutil.BoolValidator()},
	},
	Get:       getPolicy,
	Apply:     applyPolicy,
	PostApply: nil,
}

// HandleCommand is the catalog command handler for the "matrix" verb.
// It delegates to the declarative CmdHandler which handles atomic validation
// and application of key=value arguments.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func init() {
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "attract_matrix",
		Summary: "Query or set matrix rain display options.",
		Usage:   "attract_matrix [min_speed=<float>] [max_speed=<float>] [trail_length=<int>] [density=<float>] [show_background=<bool>]",
		Handle:  HandleCommand,
	})
}
