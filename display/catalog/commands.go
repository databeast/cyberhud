package catalog

import (
	"sort"
	"strings"
	"sync"
)

// CommandHandler processes mode-specific console sub-commands.
type CommandHandler func(args []string) string

// CommandDefinition describes a registered mode-specific console command verb.
type CommandDefinition struct {
	Verb    string
	Summary string
	Usage   string
	Handle  CommandHandler
}

var (
	commandsMu sync.RWMutex
	commands   = map[string]CommandDefinition{}
)

// RegisterCommand publishes a mode-specific console command handler.
func RegisterCommand(def CommandDefinition) {
	def.Verb = strings.ToLower(strings.TrimSpace(def.Verb))
	if def.Verb == "" || def.Handle == nil {
		return
	}
	def.Summary = strings.TrimSpace(def.Summary)
	def.Usage = strings.TrimSpace(def.Usage)
	commandsMu.Lock()
	defer commandsMu.Unlock()
	commands[def.Verb] = def
}

// Command returns a registered command definition by verb.
func Command(verb string) (CommandDefinition, bool) {
	verb = strings.ToLower(strings.TrimSpace(verb))
	commandsMu.RLock()
	defer commandsMu.RUnlock()
	def, ok := commands[verb]
	return def, ok
}

// Commands returns all registered command definitions ordered by verb.
func Commands() []CommandDefinition {
	commandsMu.RLock()
	defer commandsMu.RUnlock()
	out := make([]CommandDefinition, 0, len(commands))
	for _, def := range commands {
		out = append(out, def)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Verb < out[j].Verb
	})
	return out
}
