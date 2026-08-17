package displaymode

import (
	"sync"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// panelHintsState stores the current panel's TextHints for fitness evaluation.
// These are set by the runtime during display initialization and read by the
// PostApply hook to generate fitness notes on style changes.
var panelHintsState struct {
	sync.RWMutex
	hints textlayout.TextHints
	set   bool
}

// SetPanelHints stores the active panel's TextHints for fitness note generation.
// Called by the display runtime when the panel is initialized.
func SetPanelHints(hints textlayout.TextHints) {
	panelHintsState.Lock()
	defer panelHintsState.Unlock()
	panelHintsState.hints = hints
	panelHintsState.set = true
}

// getPanelHints returns the stored panel hints and whether they've been set.
func getPanelHints() (textlayout.TextHints, bool) {
	panelHintsState.RLock()
	defer panelHintsState.RUnlock()
	return panelHintsState.hints, panelHintsState.set
}

// fitnessNotesPostApply is the PostApply hook for the template CmdHandler.
// It checks if the "style" key was applied and, if so, generates fitness notes
// for the applied style against the current panel hints.
func fitnessNotesPostApply(appliedKeys []string) []string {
	// Only generate notes when the style key was changed.
	styleChanged := false
	for _, k := range appliedKeys {
		if k == "style" {
			styleChanged = true
			break
		}
	}
	if !styleChanged {
		return nil
	}

	hints, ok := getPanelHints()
	if !ok {
		// No panel hints available (e.g., testing or headless mode).
		return nil
	}

	// Look up the applied style from the registry.
	p := GetPolicy()
	s := templateRegistry.Lookup(p.Style)
	if s == nil {
		return nil
	}

	return style.FitnessNotes(s, hints)
}
