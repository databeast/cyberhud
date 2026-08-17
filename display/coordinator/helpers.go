package coordinator

import "strings"

func regionStatus(index int, p regionState) RegionStatus {
	cur := ""
	if p.current >= 0 && p.current < len(p.modes) {
		cur = p.modes[p.current]
	}
	modes := make([]string, len(p.modes))
	copy(modes, p.modes)
	return RegionStatus{
		Index:      index,
		Name:       p.name,
		Controller: p.controller,
		Current:    cur,
		Modes:      modes,
	}
}

func normalizeModes(modes []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(modes))
	for _, mode := range modes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		if mode == "" || seen[mode] {
			continue
		}
		seen[mode] = true
		out = append(out, mode)
	}
	return out
}

func findModeIndex(modes []string, mode string) int {
	mode = strings.ToLower(strings.TrimSpace(mode))
	for i, m := range modes {
		if m == mode {
			return i
		}
	}
	return -1
}
