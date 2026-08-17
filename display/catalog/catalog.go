package catalog

import (
	"sort"
	"strings"
	"sync"
)

// OptionDefinition describes one published configuration option for a display mode.
type OptionDefinition struct {
	Key     string
	Type    string
	Summary string
	Default string
	Allowed []string
}

// Definition describes a known display mode.
type Definition struct {
	ID      string
	Title   string
	Scope   string
	Summary string
	Order   int
	Options []OptionDefinition
}

var (
	definitionsMu sync.RWMutex
	definitions   = map[string]Definition{}
)

// Register publishes a display mode definition from the mode package that owns it.
func Register(def Definition) {
	def.ID = strings.ToLower(strings.TrimSpace(def.ID))
	if def.ID == "" {
		return
	}
	def.Title = strings.TrimSpace(def.Title)
	if def.Title == "" {
		def.Title = strings.ToUpper(def.ID)
	}
	def.Scope = strings.TrimSpace(def.Scope)
	def.Summary = strings.TrimSpace(def.Summary)
	def.Options = append([]OptionDefinition(nil), def.Options...)

	definitionsMu.Lock()
	defer definitionsMu.Unlock()
	definitions[def.ID] = def
}

// Describe returns metadata for a known display mode.
func Describe(mode string) (Definition, bool) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	definitionsMu.RLock()
	defer definitionsMu.RUnlock()
	d, ok := definitions[mode]
	return d, ok
}

// Definitions returns all registered definitions ordered by priority then ID.
func Definitions() []Definition {
	definitionsMu.RLock()
	defer definitionsMu.RUnlock()
	out := make([]Definition, 0, len(definitions))
	for _, def := range definitions {
		out = append(out, def)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// IDs returns all registered mode IDs in priority order.
func IDs() []string {
	defs := Definitions()
	out := make([]string, 0, len(defs))
	for _, def := range defs {
		out = append(out, def.ID)
	}
	return out
}

// DescribeMany returns descriptions for the provided mode IDs in order.
func DescribeMany(modes []string) []Definition {
	out := make([]Definition, 0, len(modes))
	for _, mode := range modes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		if mode == "" {
			continue
		}
		if d, ok := Describe(mode); ok {
			out = append(out, d)
			continue
		}
		out = append(out, Definition{ID: mode, Title: strings.ToUpper(mode), Scope: "display mode", Summary: "User-defined mode without built-in catalog metadata."})
	}
	return out
}
