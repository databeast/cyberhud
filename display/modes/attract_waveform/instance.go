package attract_waveform

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_waveform", newInstance)
}

// instance implements displaymodes.ModeInstance for the attract_waveform mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_waveform" }

func (i *instance) Activate() {} // waveform is frame-driven, no background work

func (i *instance) Deactivate() {} // waveform is frame-driven, no background work

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
