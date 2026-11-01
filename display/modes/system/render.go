package system

import (
	"fmt"
	"image"
	"strings"

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

	// For cores/top: if no stored sample, try sampling fresh. The active
	// variant is the trailing suffix of the resolution-qualified style name
	// (e.g. "color-240x240-cores" → "cores"), not the whole name.
	variant := styleVariant(p.Style)
	if cpuSample == nil && variant == StyleCores {
		cpuSample = source.SampleCPU()
	}
	if processes == nil && variant == StyleTop {
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

// styleVariant returns the content variant (default|compact|cores|top) encoded
// as the trailing suffix of a resolution-qualified style name such as
// "color-240x240-cores". A bare legacy name ("cores") or a skeleton name with
// no variant suffix ("color-240x240") returns the whole-name / empty case:
// for a skeleton the variant is "" (adaptiveBuild decides), and for a legacy
// bare name the name itself is the variant.
func styleVariant(styleName string) string {
	switch {
	case strings.HasSuffix(styleName, "-"+StyleCores):
		return StyleCores
	case strings.HasSuffix(styleName, "-"+StyleTop):
		return StyleTop
	case strings.HasSuffix(styleName, "-"+StyleDefault):
		return StyleDefault
	case strings.HasSuffix(styleName, "-"+StyleCompact):
		return StyleCompact
	case styleName == StyleCores, styleName == StyleTop,
		styleName == StyleDefault, styleName == StyleCompact:
		return styleName // legacy bare name
	default:
		return "" // skeleton (no variant suffix)
	}
}

// RenderCacheKey returns a change-detection signature incorporating sampled data.
// When the active variant is "cores", it samples CPU data and includes the
// values. When "top", it samples process data. For other variants, it returns
// just the style fingerprint.
func RenderCacheKey() uint32 {
	p := GetPolicy()
	switch styleVariant(p.Style) {
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
