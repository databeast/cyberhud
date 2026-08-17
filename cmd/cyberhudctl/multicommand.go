package main

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/regionid"
)

// SplitCommands splits a CLI argument list by standalone semicolons into
// individual commands. Only standalone ";" args are treated as separators.
// No quoted-string multi-command form is supported.
func SplitCommands(args []string) [][]string {
	if len(args) == 0 {
		return nil
	}

	var result [][]string
	var current []string

	for _, arg := range args {
		if arg == ";" {
			if len(current) > 0 {
				result = append(result, current)
				current = nil
			}
		} else {
			current = append(current, arg)
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

// RegionContext tracks the active region for scoped commands.
type RegionContext struct {
	Active *regionid.ID
}

// ResolveCommand expands a scoped command (mode, config, next, prev, status)
// into a fully-qualified protocol command using the active region context.
// If the command is "region <id>", it sets the active context and returns an
// empty string (no protocol command to send).
func (ctx *RegionContext) ResolveCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("empty command")
	}

	verb := strings.ToLower(args[0])

	switch verb {
	case "region":
		if len(args) < 2 {
			return "", fmt.Errorf("usage: region <region_id>")
		}
		raw := args[1]
		// Try parsing as a fully-qualified region ID first.
		id, err := regionid.Parse(raw)
		if err != nil {
			// Check if it's a bare integer — still store it as a region ID
			// with the integer as the surface (bare int resolution happens
			// daemon-side, but for context tracking we store the parsed form).
			if n, ok := regionid.ParseBareInt(raw); ok {
				// For bare integers, we can't fully resolve without coordinator
				// state. Store a placeholder that uses the integer as the string
				// representation. The daemon will resolve it.
				// We store the raw string so ResolveCommand can pass it through.
				ctx.Active = &regionid.ID{Surface: fmt.Sprintf("%d", n), Index: -1}
				return "", nil
			}
			return "", fmt.Errorf("invalid region identifier: %v", err)
		}
		ctx.Active = &id
		return "", nil

	case "mode":
		if ctx.Active == nil {
			return "", fmt.Errorf("no region context active; use 'region <id>' first")
		}
		if len(args) < 2 {
			return "", fmt.Errorf("usage: mode <mode_name>")
		}
		regionStr := ctx.regionString()
		mode := strings.ToLower(args[1])
		// Include any additional key=value args after mode name.
		parts := []string{"display set", regionStr, mode}
		for _, kv := range args[2:] {
			parts = append(parts, strings.TrimSpace(kv))
		}
		return strings.Join(parts, " "), nil

	case "config":
		if ctx.Active == nil {
			return "", fmt.Errorf("no region context active; use 'region <id>' first")
		}
		regionStr := ctx.regionString()
		parts := []string{"display config", regionStr}
		for _, kv := range args[1:] {
			parts = append(parts, strings.TrimSpace(kv))
		}
		return strings.Join(parts, " "), nil

	case "next":
		if ctx.Active == nil {
			return "", fmt.Errorf("no region context active; use 'region <id>' first")
		}
		return fmt.Sprintf("display next %s", ctx.regionString()), nil

	case "prev":
		if ctx.Active == nil {
			return "", fmt.Errorf("no region context active; use 'region <id>' first")
		}
		return fmt.Sprintf("display prev %s", ctx.regionString()), nil

	case "status":
		if ctx.Active == nil {
			return "", fmt.Errorf("no region context active; use 'region <id>' first")
		}
		return fmt.Sprintf("display policy %s", ctx.regionString()), nil

	default:
		return "", fmt.Errorf("unknown scoped command %q", verb)
	}
}

// regionString returns the string representation of the active region.
// For bare integers (Index == -1), it returns just the surface string (the integer).
// For fully-qualified IDs, it returns the canonical <surface>.<index> form.
func (ctx *RegionContext) regionString() string {
	if ctx.Active == nil {
		return ""
	}
	if ctx.Active.Index == -1 {
		// Bare integer stored as surface name.
		return ctx.Active.Surface
	}
	return ctx.Active.String()
}
