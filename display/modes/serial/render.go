package serial

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView resolves the serial style for the panel and lets it place the
// current snapshot. All logic — the monitor goroutine, ANSI parsing, snapshot
// assembly — lives in the source package; this only dispatches through the
// registry, exactly like clock/wifi/pager.
func BuildView(hints textlayout.TextHints) style.ViewData {
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		if current, ok := modehints.Current(); ok {
			hints = current
		}
	}
	snap := source.SnapshotNow()
	p := normalizePolicy(PolicySnapshot())

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(serialRegistry, hints, "serial", p.Style)
	ctx := style.NewStyleContext(hints)
	vd := s.Build(snap, p, ctx)
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	return vd
}

// Signature returns a stable change token for UI refresh, incorporating the
// buffer sequence, connection state, port/baud, error, and line count.
func Signature() uint32 {
	return region.CalcRegionCacheKey(signatureFromSnapshot(source.SnapshotNow(), PolicySnapshot()))
}

func signatureFromSnapshot(snap Snapshot, p Policy) string {
	conn := "0"
	if snap.Connected {
		conn = "1"
	}
	p = normalizePolicy(p)
	return fmt.Sprintf("serial:%d:%s:%s:%d:%t:%s:%d", snap.Sequence, conn, snap.Port, snap.Baud, snap.AutoSelect, snap.LastError, len(snap.Lines))
}
