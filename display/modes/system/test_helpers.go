package system

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/system/sampler"
	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func NormalizePolicy(p Policy) Policy { return normalizePolicy(p) }
func RegisteredStyleNames() []string  { return registeredStyleNames() }
func BuildCoresView(hints textlayout.TextHints, getIcon func(string) (image.Image, bool)) style.ViewData {
	old := GetPolicy()
	p := old
	p.Style = StyleCores
	SetPolicy(p)
	defer SetPolicy(old)
	return BuildView(hints, getIcon)
}
func BuildTopView(hints textlayout.TextHints) style.ViewData {
	old := GetPolicy()
	p := old
	p.Style = StyleTop
	SetPolicy(p)
	defer SetPolicy(old)
	return BuildView(hints, nil)
}
func SetSamplerForTest(s sampler.Sampler) { source.SetSampler(s) }
