package pager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/pager/source"
)

// sourceValidator validates that the source path is absolute.
// Empty string is accepted (means "no source configured").
func sourceValidator() cmdutil.KeyValidator {
	return func(value string) string {
		v := strings.TrimSpace(value)
		if v == "" {
			return "" // empty is valid (unset source)
		}
		if err := source.ValidateSource(v); err != nil {
			return err.Error()
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

// intAcceptValidator returns a KeyValidator that accepts any valid integer.
// It rejects non-integer values but does NOT reject out-of-range integers.
// Clamping to the valid range is handled by normalizePolicy during SetPolicy.
func intAcceptValidator() cmdutil.KeyValidator {
	return func(value string) string {
		v := strings.TrimSpace(value)
		_, err := strconv.Atoi(v)
		if err != nil {
			return "must be an integer"
		}
		return ""
	}
}

// intMinValidator returns a KeyValidator that accepts integer values >= min.
func intMinValidator(min int) cmdutil.KeyValidator {
	return func(value string) string {
		v := strings.TrimSpace(value)
		n, err := strconv.Atoi(v)
		if err != nil {
			return "must be an integer"
		}
		if n < min {
			return fmt.Sprintf("must be >= %d", min)
		}
		return ""
	}
}

// stringValidator accepts any string value (no validation).
func stringValidator() cmdutil.KeyValidator {
	return func(value string) string {
		return ""
	}
}

// cmdHandler is the declarative CmdHandler for the "pager" command verb.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "pager",
	Keys: []cmdutil.KeyDef{
		{Name: "source", Validate: sourceValidator()},
		{Name: "scroll_speed", Validate: intAcceptValidator()},
		{Name: "max_lines", Validate: intRangeValidator(1, 1000)},
		{Name: "scan_ms", Validate: intRangeValidator(100, 30000)},
		{Name: "font", Validate: stringValidator()},
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
		{Name: "fade_out_ms", Validate: intMinValidator(0)},
		{Name: "fade_in_ms", Validate: intMinValidator(0)},
		{Name: "line_time_ms", Validate: intMinValidator(1)},
		{Name: "max_wait_s", Validate: intMinValidator(1)},
	},
	Get:       getPolicyValue,
	Apply:     applyPolicyValue,
	PostApply: sourceChangePostApply,
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "source":
		return p.Source
	case "scroll_speed":
		return strconv.Itoa(p.ScrollSpeed)
	case "max_lines":
		return strconv.Itoa(p.MaxLines)
	case "scan_ms":
		return strconv.Itoa(p.ScanMS)
	case "font":
		return p.Font
	case "style":
		return p.Style
	case "fade_out_ms":
		return strconv.Itoa(p.FadeOutMS)
	case "fade_in_ms":
		return strconv.Itoa(p.FadeInMS)
	case "line_time_ms":
		return strconv.Itoa(p.LineTimeMS)
	case "max_wait_s":
		return strconv.Itoa(p.MaxWaitS)
	default:
		return ""
	}
}

// applyPolicyValue updates a single policy key with its validated value.
// This is called only after all key=value pairs have passed validation,
// ensuring atomic all-or-nothing semantics.
func applyPolicyValue(key, value string) {
	p := GetPolicy()
	v := strings.TrimSpace(value)
	switch key {
	case "source":
		p.Source = v
	case "scroll_speed":
		if n, err := strconv.Atoi(v); err == nil {
			p.ScrollSpeed = n
		}
	case "max_lines":
		if n, err := strconv.Atoi(v); err == nil {
			p.MaxLines = n
		}
	case "scan_ms":
		if n, err := strconv.Atoi(v); err == nil {
			p.ScanMS = n
		}
	case "font":
		p.Font = v
	case "style":
		p.Style = v
	case "fade_out_ms":
		if n, err := strconv.Atoi(v); err == nil {
			p.FadeOutMS = n
		}
	case "fade_in_ms":
		if n, err := strconv.Atoi(v); err == nil {
			p.FadeInMS = n
		}
	case "line_time_ms":
		if n, err := strconv.Atoi(v); err == nil {
			p.LineTimeMS = n
		}
	case "max_wait_s":
		if n, err := strconv.Atoi(v); err == nil {
			p.MaxWaitS = n
		}
	}
	SetPolicy(p)
}

// HandleCommand is the catalog command handler for the "pager" verb.
// It delegates to the declarative CmdHandler which handles atomic validation
// and application of key=value arguments.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

func sourceChangePostApply(appliedKeys []string) []string {
	sourceChanged := false
	scrollSpeedChanged := false
	for _, k := range appliedKeys {
		switch k {
		case "source":
			sourceChanged = true
		case "scroll_speed":
			scrollSpeedChanged = true
		}
	}

	if scrollSpeedChanged {
		_, _, scroll, _ := source.ActiveStateSnapshot()
		if scroll != nil {
			scroll.SetBaseSpeed(GetPolicy().ScrollSpeed)
		}
	}

	if sourceChanged {
		source.RestartActiveReader(GetPolicy())
	}
	return nil
}

func init() {
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "pager",
		Summary: "Query or set pager display mode options.",
		Usage:   "pager [source=<path>] [scroll_speed=<1-1000>] [max_lines=<1-1000>] [scan_ms=<100-30000>] [font=<name>] [style=<name>] [fade_out_ms=<int>] [fade_in_ms=<int>] [line_time_ms=<int>] [max_wait_s=<int>]",
		Handle:  HandleCommand,
	})
}
