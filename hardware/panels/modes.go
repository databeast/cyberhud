package panels

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/coordinator"
	_ "github.com/databeast/cyberhud/display/modes"
)

// AllModes returns all registered display mode IDs.
func AllModes() []string {
	return catalog.IDs()
}

func filterExcluded(modes []string, excluded []string) []string {
	if len(modes) == 0 {
		return nil
	}
	ex := make(map[string]struct{}, len(excluded))
	for _, mode := range excluded {
		mode = strings.ToLower(strings.TrimSpace(mode))
		if mode == "" {
			continue
		}
		ex[mode] = struct{}{}
	}
	if len(ex) == 0 {
		return append([]string(nil), modes...)
	}
	out := make([]string, 0, len(modes))
	for _, mode := range modes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		if mode == "" {
			continue
		}
		if _, blocked := ex[mode]; blocked {
			continue
		}
		out = append(out, mode)
	}
	return out
}

// Displays returns mode-state region definitions for one panel product.
func Displays(def Definition, policy catalog.DefaultModePolicy) []coordinator.Region {
	allModes := AllModes()
	if len(def.Virtual) > 0 {
		out := make([]coordinator.Region, 0, len(def.Virtual))
		for _, d := range def.Virtual {
			modes := filterExcluded(allModes, d.ExcludedModes)
			current := policy.ResolveDefault(modes, d.InputEnabled)
			out = append(out, coordinator.Region{
				Index:      d.Index,
				Name:       d.Name,
				Controller: d.Controller,
				Modes:      modes,
				Default:    current,
			})
		}
		return out
	}

	modes := filterExcluded(allModes, def.ExcludedModes)
	current := policy.ResolveDefault(modes, def.InputEnabled)
	return []coordinator.Region{{
		Index:      0,
		Name:       def.Name,
		Controller: def.Controller,
		Modes:      modes,
		Default:    current,
	}}
}
