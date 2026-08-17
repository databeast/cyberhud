package system

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the system mode view data with sprites, adaptive font, and style support.
// The getIcon parameter resolves icon names to images (typically icons.Get).
func BuildView(hints textlayout.TextHints, getIcon func(name string) (image.Image, bool)) style.ViewData {
	p := GetPolicy()

	// Consume stored sampler data (primed by RenderCacheKey).
	cpuSample := source.ConsumeCPUSample()
	processes := source.ConsumeTopProcesses()

	// For cores/top: if no stored sample, try sampling fresh.
	if cpuSample == nil && p.Style == StyleCores {
		cpuSample = source.SampleCPU()
	}
	if processes == nil && p.Style == StyleTop {
		processes = source.TopProcesses()
	}

	// Build snapshot for style dispatch.
	snap := source.BuildSnapshot(cpuSample, processes, getIcon)

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(systemRegistry, hints, "system", p.Style)

	// Construct StyleContext boundary and invoke the style's Build method.
	ctx := style.NewStyleContext(hints)
	result := s.Build(snap, p, ctx)

	return style.ViewData{
		Static:      true,
		Items:       result.Items,
		Sprites:     result.Sprites,
		Colors:      result.Colors,
		Tiers:       result.Tiers,
		FontIDs:     result.FontIDs,
		LineOffsets: result.LineOffsets,
		OffsetY:     result.OffsetY,
		StyleReport: style.StyleReport{Name: s.Name(), Reason: reason},
	}
}

// RenderCacheKey returns a change-detection signature incorporating sampled data.
// When style is "cores", it samples CPU data and includes the values.
// When style is "top", it samples process data and includes a representation.
// For other styles, it returns just the style name.
func RenderCacheKey() uint32 {
	p := GetPolicy()
	switch p.Style {
	case StyleCores:
		sample := source.PrimeCPUSample()
		if sample == nil {
			return region.CalcRegionCacheKey("cores:nil", p.Fingerprint())
		}
		return region.CalcRegionCacheKey(fmt.Sprintf("cores:%v", sample), p.Fingerprint())
	case StyleTop:
		procs := source.PrimeTopProcesses()
		if procs == nil {
			return region.CalcRegionCacheKey("top:nil", p.Fingerprint())
		}
		return region.CalcRegionCacheKey(fmt.Sprintf("top:%v", procs), p.Fingerprint())
	default:
		return region.CalcRegionCacheKey(p.Fingerprint())
	}
}
