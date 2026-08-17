package menu

import (
	"image"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/widgets"
)

// Allowed style values for the menu mode.
const (
	StyleFramed = "framed"
	StylePlain  = "plain"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current menu policy (thread-safe).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the menu policy after normalization (thread-safe).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "menu",
		Title:   "Menu",
		Summary: "Interactive main menu for navigating top-level screens.",
		Order:   10,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "menu",
		Summary: "Query or set menu display options.",
		Usage:   "menu [style=<framed|plain>]",
		Handle:  HandleCommand,
	})
}

// Items returns the top-level menu entries for the primary UI mode.
func Items() []string { return source.Items() }

// Destinations maps cursor position to the mode to navigate to on select.
var Destinations = []string{"stemma", "gpio", "gpio-control", "usb", "serial", "system"}

// BorderBuilder is a function that builds border frame sprites for a given pixel bounds.
type BorderBuilder func(bounds image.Rectangle) []widgets.Sprite
