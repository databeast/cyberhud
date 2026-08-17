package displaymodes

import (
	"sync"

	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// PanelHints is embeddable per-instance storage for the hosting Region's text hints.
//
// # What problem this solves
//
// A display mode needs its Region's geometry — pixel dimensions, capability, and the
// tier catalog — to lay anything out. Modes used to fetch that from
// display/region/modehints, a package-level singleton holding whatever was
// propagated most recently. With one Region that works by accident. With two, both
// modes read the same store and the later activation overwrites the earlier, so one
// Region silently lays out for the other Region's panel.
//
// Hints are per-region data and so must be per-instance state. Embedding PanelHints
// gives a mode that state, and the Region fills it during SetMode via the
// region.HintsReceiver interface, before Activate and before any BuildView call.
//
// # How to adopt it in a mode
//
// Embed it in the mode's instance struct and read hints through Hints():
//
//	type instance struct {
//	    displaymodes.PanelHints
//	    // ... mode state
//	}
//
//	func (i *instance) BuildView() style.ViewData {
//	    hints, ok := i.Hints()
//	    if !ok {
//	        return style.ViewData{Items: []string{"error"}}
//	    }
//	    return BuildView(time.Now(), hints)
//	}
//
// Embedding supplies SetPanelHints, which is what satisfies region.HintsReceiver, so
// no further wiring is needed. See display/modes/clock for a converted mode.
//
// # Migration status
//
// Adoption is per-mode and incremental. A mode that has not embedded this continues
// to call its local getPanelHints helper, which reads the legacy global and is
// correct only for single-region setups. Hints() falls back to that same global so a
// partially converted mode cannot end up with nothing; the fallback is the thing to
// delete once every mode is converted, along with the mutable state in modehints.
//
// New modes should embed this rather than adding another getPanelHints.
type PanelHints struct {
	mu    sync.RWMutex
	hints textlayout.TextHints
	set   bool
}

// SetPanelHints records the hosting Region's hints. Called by Region.SetMode; modes
// do not call it themselves.
//
// Guarded by a mutex because activation and rendering can occur on different
// goroutines: SetMode runs on whichever goroutine requested the mode change, while
// BuildView runs on the render loop.
func (p *PanelHints) SetPanelHints(hints textlayout.TextHints) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hints = hints
	p.set = true
}

// Hints returns the hints injected by the hosting Region, and whether any are
// available.
//
// Falls back to the legacy process-global store when the Region has not injected
// any, which happens for instances constructed outside a Region (tests and tooling
// that call GetInstance directly). The fallback keeps those paths working during
// migration; it is not correct for multi-region use, which is the whole reason
// injection exists.
func (p *PanelHints) Hints() (textlayout.TextHints, bool) {
	p.mu.RLock()
	hints, set := p.hints, p.set
	p.mu.RUnlock()
	if set {
		return hints, true
	}
	return legacyGlobalHints()
}

// legacyGlobalHints reads the process-wide hints store.
//
// Isolated in one function so the migration has a single call site to delete. When
// every mode embeds PanelHints and receives injection from its Region, this and the
// mutable state behind it can go.
func legacyGlobalHints() (textlayout.TextHints, bool) {
	return modehints.Current()
}
