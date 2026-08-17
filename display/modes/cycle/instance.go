package cycle

import (
	"context"
	"log"
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the cycle mode.
const cycleFontFamily = "spleen"

var cycleFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("cycle", newInstance)
}

// cycleInstance implements displaymodes.ModeInstance for the cycle mode.
// Cycle has lifecycle: Activate starts a sequencer goroutine that auto-cycles
// modes, and Deactivate stops it via context cancellation.
type cycleInstance struct {
	cancel      context.CancelFunc
	done        chan struct{}
	activateIdx int // region index that activated this instance (defaults to 0)
}

func newInstance() displaymodes.ModeInstance {
	return &cycleInstance{activateIdx: 0}
}

func (i *cycleInstance) ID() string { return ModeID }

// Activate starts the cycle sequencer goroutine that auto-cycles modes.
func (i *cycleInstance) Activate() {
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel
	i.done = make(chan struct{})

	go i.run(ctx)
}

// Deactivate stops the cycle sequencer goroutine using context cancellation.
// Completes within 5 seconds.
func (i *cycleInstance) Deactivate() {
	if i.cancel != nil {
		i.cancel()
	}
	if i.done != nil {
		select {
		case <-i.done:
		case <-time.After(5 * time.Second):
			log.Printf("cycle.Deactivate: sequencer did not stop within 5s")
		}
	}
}

// run is the background goroutine that cycles display modes on a timer.
func (i *cycleInstance) run(ctx context.Context) {
	defer close(i.done)

	policy := GetPolicy()
	interval := normalizeInterval(policy.Interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state := getGlobalModeState()
			if state == nil {
				continue
			}

			// Re-read policy for hot-reload support.
			p := GetPolicy()

			// Hot-reload interval: if policy interval changed, reset ticker.
			newInterval := normalizeInterval(p.Interval)
			if newInterval != interval {
				interval = newInterval
				ticker.Reset(interval)
			}

			// Determine regions to cycle.
			regions := p.Regions
			if len(regions) == 0 {
				regions = []int{i.activateIdx}
			}

			// Advance each region, skipping "cycle" itself.
			for _, idx := range regions {
				advanceSkippingSelf(idx, state, p.Modes)
			}
		}
	}
}

func (i *cycleInstance) ActionHandler() action.Handler { return Handler{} }

// Handler is the action handler for the cycle mode (no-op).
type Handler struct{}

// HandleAction is a no-op for the cycle mode.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}

// BuildView returns a static placeholder ViewData for the cycle mode.
func (i *cycleInstance) BuildView() style.ViewData {
	state := style.ViewData{
		Items:  []string{"Cycling modes..."},
		Static: true,
	}

	return state
}

// RenderCacheKey returns a stable cache key since cycle mode content doesn't change.
func (i *cycleInstance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey("cycle")
}
