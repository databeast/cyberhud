package pager

import (
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("pager", newInstance)

	catalog.Register(catalog.Definition{
		ID:       "pager",
		Title:    "Pager",
		Category: "Data",
		Scope:    "any",
		Summary:  "Tails a file, pipe, or socket and presents text with smooth scroll or page transitions.",
		Order:    90,
		Options:  append(source.Policy{}.Options(), catalog.OptionDefinition{Key: "style", Type: "string", Summary: "Visual style name or empty for default.", Default: "", Allowed: pagerRegistry.RegisteredNames()}),
	})
}

// maxTickElapsed caps per-frame advancement so a long pause (debugger, suspend)
// does not jump the scroll or page animation forward in one step.
const maxTickElapsed = 80 * time.Millisecond

// instance implements displaymodes.ModeInstance for the pager mode.
//
// It owns the per-frame clock (lastTick) and the surface classification from the
// previous pass (lastSurface); both are mode state, not layout. The embedded
// PanelHints supplies this instance's own Region hints, injected by
// Region.SetMode before Activate.
type instance struct {
	displaymodes.PanelHints

	lastSurface source.SurfaceClass
	lastTick    time.Time
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "pager" }

// Activate creates the shared buffer and reader, classifies the surface, and
// sets up the matching animation state, then starts the reader if a source is
// configured.
func (i *instance) Activate() {
	pol := GetPolicy()

	buf := source.NewLineBuffer(pol.MaxLines)
	reader := &source.TailReader{}

	hints, _ := i.Hints()
	sc := source.ClassifySurface(hints)
	i.lastSurface = sc

	source.SetActiveReader(reader, buf)
	i.installStrategy(sc, hints, pol)

	if pol.Source != "" {
		reader.Start(pol, buf)
	}
}

// Deactivate stops the reader and releases the active state.
func (i *instance) Deactivate() {
	reader, _, _, _ := source.ActiveStateSnapshot()
	if reader != nil {
		reader.Stop()
	}
	source.ClearActiveReader()
}

func (i *instance) ActionHandler() action.Handler { return nil }

// BuildView advances the per-frame clock to compute elapsed time, then delegates
// to buildView, which owns all pager logic and dispatches through the registry.
func (i *instance) BuildView() style.ViewData {
	hints, ok := i.Hints()
	if !ok {
		return style.ViewData{Items: []string{"[pager] no panel hints available"}, Static: true}
	}

	// Swap animation strategy if the panel's classification changed since the
	// last pass, without discarding buffered content.
	sc := source.ClassifySurface(hints)
	if sc != i.lastSurface {
		i.lastSurface = sc
		i.installStrategy(sc, hints, GetPolicy())
	}

	// Compute capped elapsed time for this pass.
	var elapsed time.Duration
	now := frameclock.Now()
	if !i.lastTick.IsZero() {
		elapsed = now.Sub(i.lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	i.lastTick = now

	return buildView(hints, elapsed)
}

func (i *instance) RenderCacheKey() uint32 {
	return renderCacheKey()
}

// installStrategy sets up the animation state for a surface class, clearing the
// other strategy's state. The shared buffer is left intact — only the animation
// changes.
func (i *instance) installStrategy(sc source.SurfaceClass, hints textlayout.TextHints, pol source.Policy) {
	switch sc {
	case source.SurfaceFast:
		source.SetActivePage(nil)
		source.SetActiveScroll(source.NewScrollState(pol.ScrollSpeed))
		region.RegisterTickRate("pager", tickProvider{})
	case source.SurfaceSlow:
		source.SetActiveScroll(nil)
		layout := source.ComputeLayout(style.NewStyleContext(hints))
		source.SetActivePage(source.NewPageState(layout.VisibleRows, pol.LineTimeMS))
	}
}
