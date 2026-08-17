package zmq

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the zmq mode.
// Previously stored in modeFontConfig["zmq"] = {Family: "spleen", Tier: TierNormal}.
const zmqFontFamily = "spleen"

var zmqFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("zmq", newInstance)
}

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// instance implements displaymodes.ModeInstance for the ZMQ mode.
// It manages the ZMQ subscription lifecycle via Activate/Deactivate.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "zmq" }

// Activate starts the ZMQ receiver using the current policy.
func (i *instance) Activate() {
	content.Activate()
}

// Deactivate stops the ZMQ receiver. The message buffer is preserved.
func (i *instance) Deactivate() {
	content.Deactivate()
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	vd := BuildView(hints)

	return vd
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
