package storage

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// buildView constructs a full style.ViewData for the storage mode, using the
// current policy and the provided TextHints for adaptive style resolution.
func buildView(hints textlayout.TextHints) style.ViewData {
	pol := GetPolicy()
	snap := source.CurrentSnapshot()

	// If no sample has landed yet, try one synchronously (the loop primes the
	// first sample on activation, but tooling/test contexts may not have run it).
	if len(snap.Mounts) == 0 {
		if fresh, ok := source.SampleActive(pol.ShowAll); ok {
			snap = fresh
		} else {
			return style.ViewData{Items: []string{"Scanning..."}, Static: true}
		}
	}

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(storageRegistry, hints, "storage", pol.Style)

	ctx := style.NewStyleContext(hints)
	vd := s.Build(snap, pol, ctx)

	vd.Static = true
	vd.StyleReport = style.ReportFor(s, reason, hints)
	if vd.OffsetY < 0 {
		vd.OffsetY = 0
	}
	return vd
}

// renderCacheKey returns a change-detection signature combining the mount usage
// values and a policy fingerprint, plus the instance pixel dimensions so two
// panels of different sizes produce distinct keys.
func renderCacheKey(snap source.StorageSnapshot, pixelWidth, pixelHeight int) uint32 {
	pol := GetPolicy()

	var b strings.Builder
	for i, m := range snap.Mounts {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(fmt.Sprintf("%s:%d/%d", m.Mountpoint, m.UsedBytes, m.TotalBytes))
	}
	b.WriteByte('|')
	b.WriteString(pol.Fingerprint())
	b.WriteString(fmt.Sprintf("|%dx%d", pixelWidth, pixelHeight))

	return region.CalcRegionCacheKey(b.String())
}
