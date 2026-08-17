package cycle

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// durationValidator validates that a string is a valid duration >= MinInterval.
// Values > MaxInterval are accepted by the validator (clamped in applyKey).
func durationValidator(value string) string {
	d, err := time.ParseDuration(value)
	if err != nil {
		return "must be a valid duration (e.g., 10s, 1m)"
	}
	if d < MinInterval {
		return fmt.Sprintf("must be >= %v", MinInterval)
	}
	return ""
}

// modesValidator validates a mode list string. Always accepts any input
// (no catalog validation at command time per Requirement 8.4).
// An empty string means "clear the mode list".
func modesValidator(value string) string {
	return ""
}

// regionsValidator validates a region index string. Each comma-separated
// token must parse as a non-negative integer.
func regionsValidator(value string) string {
	if value == "" {
		return ""
	}
	tokens := strings.Split(value, ",")
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		n, err := strconv.Atoi(tok)
		if err != nil || n < 0 {
			return "must be non-negative integers"
		}
	}
	return ""
}

var cmdHandler = &cmdutil.CmdHandler{
	Verb: "cycle",
	Keys: []cmdutil.KeyDef{
		{Name: "interval", Validate: durationValidator},
		{Name: "modes", Validate: modesValidator},
		{Name: "regions", Validate: regionsValidator},
	},
	Get:   getKey,
	Apply: applyKey,
}

// getKey returns the current value for a given key name.
func getKey(key string) string {
	p := GetPolicy()
	switch key {
	case "interval":
		return p.Interval.String()
	case "modes":
		return strings.Join(p.Modes, ",")
	case "regions":
		if len(p.Regions) == 0 {
			return ""
		}
		parts := make([]string, len(p.Regions))
		for i, r := range p.Regions {
			parts[i] = strconv.Itoa(r)
		}
		return strings.Join(parts, ",")
	default:
		return ""
	}
}

// applyKey sets the value for a given key. Called only after all
// key=value pairs have passed validation.
func applyKey(key, value string) {
	switch key {
	case "interval":
		d, err := time.ParseDuration(value)
		if err != nil {
			return
		}
		p := GetPolicy()
		p.Interval = normalizeInterval(d)
		SetPolicy(p)
	case "modes":
		p := GetPolicy()
		if value == "" {
			p.Modes = nil
		} else {
			parts := strings.Split(value, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			p.Modes = parts
		}
		SetPolicy(p)
	case "regions":
		p := GetPolicy()
		if value == "" {
			p.Regions = nil
		} else {
			tokens := strings.Split(value, ",")
			ints := make([]int, 0, len(tokens))
			for _, tok := range tokens {
				n, err := strconv.Atoi(strings.TrimSpace(tok))
				if err != nil {
					continue
				}
				ints = append(ints, n)
			}
			p.Regions = ints
		}
		SetPolicy(p)
	}
}

// HandleCommand is the catalog command handler for the "cycle" verb.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}
