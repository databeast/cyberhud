package displaymodes

import (
	"log"
	"sync"

	"github.com/databeast/cyberhud/display/style"
)

// IsKnown reports whether a mode id is registered in the instance-based registry.
func IsKnown(mode string) bool {
	return IsKnownInstance(mode)
}

// styleLogState tracks the last logged mode+style per region for deduplication.
// Key: regionName (string), Value: "mode:style" (string)
var styleLogState sync.Map

// LogStyleSelection emits a style autoselection log entry if the mode+style
// combination differs from the last logged state for this region.
// If report.Name is empty, no log is emitted (mode has no style dispatch).
func LogStyleSelection(regionName string, modeID string, report style.StyleReport) {
	if report.Name == "" {
		return
	}

	current := modeID + ":" + report.Name
	if prev, ok := styleLogState.Load(regionName); ok && prev.(string) == current {
		return // deduplicated
	}
	styleLogState.Store(regionName, current)

	log.Printf("style: region=%s mode=%s style=%s reason=%s",
		regionName, modeID, report.Name, report.Reason)

	// Lever 2: surface the capability-utilization gap. When the chosen style sits
	// below the panel's tier, the panel is being driven with less than it can show
	// — the condition that was previously invisible. Emit a WARN so an
	// under-driven surface (a color panel running a mono-only mode like menu, or a
	// mode with no style at the panel's tier) is observable without a debugger.
	if report.CapGap() > 0 {
		log.Printf("style: region=%s mode=%s UNDER-UTILIZES panel: panel=%s style uses %s (gap=%d tier(s))",
			regionName, modeID, report.PanelCap, report.StyleCap, report.CapGap())
	}
}

// Warnings holds hardware pin notices or other warnings that modes like stemma
// and gpio may incorporate into their Hint field. Set once at startup via WireWarnings.
var Warnings []string

// WireWarnings sets the package-level warnings slice.
// Called from the daemon wiring (cmd/cyberhudd) after computing pin notices.
func WireWarnings(w []string) {
	Warnings = w
}
