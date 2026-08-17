package serial

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the serial mode.
// Previously stored in modeFontConfig["serial"] = {Family: "terminus", Tier: TierNormal}.
const serialFontFamily = "terminus"

var serialFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("serial", newInstance)
}

// instance implements displaymodes.ModeInstance for the serial mode.
// It manages the serial monitor lifecycle via Activate/Deactivate and
// owns scroll state (cursor, topRow) internally.
type instance struct {
	hints  textlayout.TextHints
	cursor int
	topRow int
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "serial" }

// Activate starts the serial monitor background goroutine.
func (i *instance) Activate() {
	source.Activate()
}

// Deactivate stops the serial monitor background goroutine.
func (i *instance) Deactivate() {
	source.Deactivate()
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	vd := BuildView(i.hints)

	return vd
}

func (i *instance) RenderCacheKey() uint32 {
	return Signature()
}
