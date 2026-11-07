package attract_cube

import (
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

type instance struct {
	displaymodes.PanelHints
	start time.Time
}

func (*instance) ID() string                    { return "attract_cube" }
func (i *instance) Activate()                   { i.start = frameclock.Now() }
func (*instance) Deactivate()                   {}
func (*instance) ActionHandler() action.Handler { return nil }

func (i *instance) RenderCacheKey() uint32 {
	hints, ok := i.Hints()
	if !ok {
		return region.CalcRegionCacheKey(GetPolicy().Fingerprint())
	}
	switch hints.Capability {
	case textlayout.CapMonoFast, textlayout.CapGrayscaleFast:
		return region.CalcRegionCacheKey(frameclock.Now().Sub(i.start) / (33 * time.Millisecond))
	case textlayout.CapMonoSlow, textlayout.CapGrayscaleSlow:
		return region.CalcRegionCacheKey(hints.Capability)
	case textlayout.CapColorFast:
		return region.CalcRegionCacheKey(frameclock.Now().Sub(i.start)/(33*time.Millisecond), GetPolicy().Fingerprint())
	default:
		return region.CalcRegionCacheKey(GetPolicy().Fingerprint())
	}
}

func (i *instance) BuildView() style.ViewData {
	hints, ok := i.Hints()
	if !ok || hints.PixelWidth <= 0 || hints.PixelHeight <= 0 {
		return style.ViewData{Items: []string{"cube: no panel dimensions"}}
	}
	if hints.Capability < textlayout.CapMonoSlow || hints.Capability > textlayout.CapColorFast {
		return style.ViewData{Items: []string{"cube: unsupported display"}}
	}
	seconds := frameclock.Now().Sub(i.start).Seconds()
	if seconds < 0 {
		seconds = 0
	}
	s, reason := style.ResolveStyle(cubeStyles, hints, "attract_cube", "")
	view := s.Build(frame{seconds: seconds}, GetPolicy(), style.NewStyleContext(hints))
	view.StyleReport = style.ReportFor(s, reason, hints)
	return view
}

type tickProvider struct{}

func (tickProvider) PreferredTickInterval() time.Duration { return 33 * time.Millisecond }

func init() {
	displaymodes.RegisterFactory("attract_cube", func() displaymodes.ModeInstance { return &instance{} })
	region.RegisterTickRate("attract_cube", tickProvider{})
}
