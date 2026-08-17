package ticker

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

// defaultFontFamily is the preferred font family when policy.Font is "auto".
// The tier catalog priority system (spleen > terminus > cozette) determines
// the actual face; this default ensures backward compatibility.
const defaultFontFamily = "spleen"

func init() {
	displaymodes.RegisterFactory("ticker", newInstance)
}

// instance implements displaymodes.ModeInstance for the ticker mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "ticker" }
func (i *instance) Activate()                     {} // ticker has no background work
func (i *instance) Deactivate()                   {} // ticker has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
