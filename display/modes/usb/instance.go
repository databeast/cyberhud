package usb

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	styles2 "github.com/databeast/cyberhud/display/modes/usb/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets/icons"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("usb", newInstance)
}

// instance implements displaymodes.ModeInstance for the usb mode.
// It embeds PanelHints so that per-instance panel geometry is available
// without relying on the process-global modehints store.
type instance struct {
	displaymodes.PanelHints
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "usb" }
func (i *instance) Activate()                     {} // usb has no background work
func (i *instance) Deactivate()                   {} // usb has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	// Hints come from this instance's embedded PanelHints, injected by the
	// hosting Region before Activate. Falls back to the process-wide store
	// (getPanelHints) for instances constructed outside a Region (tests, tooling).
	hints, ok := i.Hints()
	if !ok {
		hints, ok = getPanelHints()
		if !ok {
			return style.ViewData{Items: []string{"error"}}
		}
	}

	snap := SnapshotNow()
	p := PolicySnapshot()
	p = normalizePolicy(p)

	s, reason := style.ResolveStyle(usbRegistry, hints, "usb", p.Style)
	ctx := style.NewStyleContext(hints)
	if consumer, ok := any(s).(styles2.IconConsumer); ok {
		consumer.SetIconGetter(icons.Get)
	}
	vd := s.Build(snap, p, ctx)
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	vd.Static = true
	return vd
}

func (i *instance) RenderCacheKey() uint32 {
	return Signature()
}
