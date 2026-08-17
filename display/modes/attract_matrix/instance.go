package attract_matrix

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_matrix", newInstance)
}

// instance implements displaymodes.ModeInstance for the matrix mode.
type instance struct {
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_matrix" }

func (i *instance) Activate() {} // matrix has no background work (animation is frame-driven)

func (i *instance) Deactivate() {} // matrix has no background work

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

func init() {
	catalog.Register(catalog.Definition{
		ID:      "attract_matrix",
		Title:   "Matrix",
		Summary: "Animated matrix digital rain effect with scrolling columns of pseudorandom glyphs.",
		Order:   200,
		Options: source.Policy{}.Options(),
	})
}
