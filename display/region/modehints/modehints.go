// Package modehints provides the PanelHintsReceiver registration system.
// Mode packages register receivers during init(), and infrastructure code
// calls PropagateHints to broadcast panel dimensions to all modes.
//
// This package exists separately from display/modes to allow mode
// sub-packages to register without creating import cycles.
package modehints

import (
	"sync"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// PanelHintsReceiver is implemented by mode packages that need panel
// dimensions for fitness-note generation on style changes.
type PanelHintsReceiver interface {
	SetPanelHints(hints textlayout.TextHints)
}

// hintsReceivers holds all registered PanelHintsReceiver implementations.
// Registration happens during init() in each mode package.
var hintsReceivers []PanelHintsReceiver

// RegisterHintsReceiver is called by mode init() functions to register
// a receiver for panel hints propagation.
func RegisterHintsReceiver(r PanelHintsReceiver) {
	hintsReceivers = append(hintsReceivers, r)
}

// current holds the most recently propagated panel hints. All receivers get
// the same hints, so modes can read this central store via Current() instead
// of caching their own copy.
var current struct {
	sync.RWMutex
	hints textlayout.TextHints
	set   bool
}

// Current returns the most recently propagated panel hints and whether any
// have been propagated yet.
func Current() (textlayout.TextHints, bool) {
	current.RLock()
	defer current.RUnlock()
	return current.hints, current.set
}

// PropagateHints stores hints centrally and sends them to all registered
// receivers.
func PropagateHints(hints textlayout.TextHints) {
	current.Lock()
	current.hints = hints
	current.set = true
	current.Unlock()

	for _, r := range hintsReceivers {
		r.SetPanelHints(hints)
	}
}

// ResetHintsReceivers clears the registered receivers slice and the central
// hints store. Intended for use in tests to isolate registration state.
func ResetHintsReceivers() {
	hintsReceivers = nil
	current.Lock()
	current.hints = textlayout.TextHints{}
	current.set = false
	current.Unlock()
}
