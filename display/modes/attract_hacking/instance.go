package attract_hacking

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_hacking", newInstance)
}

type instance struct{}

var _ displaymodes.ModeInstance = (*instance)(nil)

func newInstance() displaymodes.ModeInstance { return &instance{} }

func (i *instance) ID() string { return "attract_hacking" }

func (i *instance) Activate() {}

func (i *instance) Deactivate() {}

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
