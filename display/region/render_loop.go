package region

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/databeast/cyberhud/hardware/input"
)

// Renderer is the RenderLoop-layer interface called once per region per tick
// cycle to produce the current frame's visual content. Implementations
// typically delegate to the Region's active [ModeInstance].
type Renderer interface {
	// Render renders the current frame for the given Region r.
	// Returns an error if rendering fails; the RenderLoop logs the error
	// and skips the Region for that frame.
	Render(r *Region) error
}

// InputDispatcher is the RenderLoop-layer interface that routes hardware input
// events to the active Region. The RenderLoop calls Dispatch after rendering
// each frame for every pending input event.
type InputDispatcher interface {
	// Dispatch delivers the input event ev to Region r, which holds input
	// focus. Implementations typically forward the event to the Region's
	// active ModeInstance action handler.
	Dispatch(r *Region, ev input.Event)
}

// PerRegionTicker is a RenderLoop-layer type that tracks the next-fire time for
// a single Region. It enables independent tick intervals so that fast-refresh
// modes (e.g. animation at 50ms) coexist with slow-refresh modes (e.g. dashboard
// at 1000ms) within the same RenderLoop.
type PerRegionTicker struct {
	region   *Region
	interval time.Duration
	lastFire time.Time
}

// RenderLoopOption is a functional option that configures a [RenderLoop].
// Options may be passed to [NewRenderLoop] at construction time or applied
// later via [RenderLoop.Apply]. The available options are:
//
//   - [WithRenderer] — sets the per-region renderer
//   - [WithInputDispatcher] — sets the input event dispatcher
//   - [WithTickInterval] — sets the global frame tick interval
//   - [WithTickRateResolver] — enables per-region deadline scheduling
type RenderLoopOption func(*RenderLoop)

// WithTickInterval returns a [RenderLoopOption] that sets the global frame
// tick interval to d. The tick interval determines how frequently all regions
// are rendered when no [TickRateResolver] is configured. The default is 1000ms.
func WithTickInterval(d time.Duration) RenderLoopOption {
	return func(rl *RenderLoop) {
		rl.tick = d
	}
}

// WithRenderer returns a [RenderLoopOption] that sets the [Renderer] used to
// draw each region's content every frame. The renderer r is called once per
// region per tick cycle. If no renderer is set, the loop skips rendering.
func WithRenderer(r Renderer) RenderLoopOption {
	return func(rl *RenderLoop) {
		rl.renderer = r
	}
}

// WithInputDispatcher returns a [RenderLoopOption] that sets the
// [InputDispatcher] used to route input events. The dispatcher d receives
// events from the input channel and delivers them to the input-active region.
func WithInputDispatcher(d InputDispatcher) RenderLoopOption {
	return func(rl *RenderLoop) {
		rl.dispatcher = d
	}
}

// WithTickRateResolver returns a [RenderLoopOption] that sets the
// [TickRateResolver] used to determine per-region tick intervals. The resolver
// queries registered [TickRateProvider] values to derive each region's refresh
// rate. When a resolver is configured, [RenderLoop.Run] uses per-region
// deadline scheduling instead of the single global ticker, allowing regions
// with different refresh requirements to coexist efficiently.
func WithTickRateResolver(resolver TickRateResolver) RenderLoopOption {
	return func(rl *RenderLoop) {
		rl.tickRateResolver = resolver
	}
}

// RenderLoop coordinates per-frame rendering across all Regions and flushing
// to Physical Screens. It runs at a fixed tick interval (default 1000ms).
// When a TickRateResolver is configured, it uses per-region deadline scheduling.
type RenderLoop struct {
	rm               *RegionManager
	fp               *FlushPath
	events           <-chan input.Event
	tick             time.Duration
	stopCh           chan struct{}
	renderer         Renderer
	dispatcher       InputDispatcher
	tickRateResolver TickRateResolver
	regionTickers    []*PerRegionTicker
}

// NewRenderLoop creates a new [RenderLoop] that coordinates rendering across
// all regions managed by rm, flushing completed frames through fp.
//
// The loop is not started on construction — call [RenderLoop.Run] to begin the
// render cycle. Run blocks until [RenderLoop.Stop] is called.
//
// After construction, use [RenderLoop.Apply] to configure the loop with
// [RenderLoopOption] values:
//
//   - [WithRenderer] — required for rendering to occur
//   - [WithInputDispatcher] — required for input event delivery
//   - [WithTickInterval] — overrides the default 1000ms tick
//   - [WithTickRateResolver] — enables per-region deadline scheduling
//
// Parameters:
//   - rm: the [RegionManager] whose allocated regions are rendered each frame
//   - fp: the [FlushPath] used to push rendered frames to hardware screens
//   - events: input event channel; may be nil if no input source is available
//   - opts: zero or more [RenderLoopOption] values applied at construction
//
// Returns the configured (but not yet running) RenderLoop.
func NewRenderLoop(rm *RegionManager, fp *FlushPath, events <-chan input.Event, opts ...RenderLoopOption) *RenderLoop {
	rl := &RenderLoop{
		rm:     rm,
		fp:     fp,
		events: events,
		tick:   1000 * time.Millisecond,
		stopCh: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(rl)
	}
	return rl
}

// Stop signals the render loop to terminate. The currently executing frame
// completes, and then [RenderLoop.Run] returns. Stop may be called from any
// goroutine. Calling Stop more than once panics (it closes the internal stop
// channel).
func (rl *RenderLoop) Stop() {
	close(rl.stopCh)
}

// Apply applies one or more [RenderLoopOption] values to the RenderLoop after
// construction. This is the primary way to configure the loop when
// [ActivatePanel] returns a pre-built RenderLoop that still needs a renderer
// and dispatcher wired in. The accepted options are [WithRenderer],
// [WithInputDispatcher], [WithTickInterval], and [WithTickRateResolver].
//
// Apply must be called before [RenderLoop.Run]; applying options while the loop
// is running is not safe for concurrent use.
func (rl *RenderLoop) Apply(opts ...RenderLoopOption) {
	for _, opt := range opts {
		opt(rl)
	}
}

// initRegionTickers builds the per-region ticker slice from the currently
// allocated regions and the configured TickRateResolver. Called once at the
// start of Run when a TickRateResolver is configured.
func (rl *RenderLoop) initRegionTickers() {
	regions := rl.rm.Regions()
	now := time.Now()
	rl.regionTickers = make([]*PerRegionTicker, 0, len(regions))
	for _, r := range regions {
		interval := rl.tickRateResolver.TickInterval(r.CurrentMode())
		rl.regionTickers = append(rl.regionTickers, &PerRegionTicker{
			region:   r,
			interval: interval,
			lastFire: now.Add(-interval), // ensure first frame is immediately due
		})
	}
}

// Run blocks the calling goroutine, rendering all regions and flushing each
// frame at the configured tick interval until [RenderLoop.Stop] is called.
// It renders and flushes one complete frame immediately before entering the
// main scheduling loop or processing any input events.
//
// When a [TickRateResolver] is configured (via [WithTickRateResolver]), Run
// initializes per-region tickers from the allocated regions and uses
// deadline-based scheduling instead of a single global ticker. In this mode
// each region has an independent deadline and only regions whose deadline has
// elapsed are rendered each cycle.
//
// Run must be called at most once per RenderLoop. The caller should arrange
// for [RenderLoop.Stop] to be called from another goroutine when shutdown is
// desired.
func (rl *RenderLoop) Run() {
	// Step 1: Initialize per-region tickers before anything else.
	// Tickers must exist before the first render because renderDueRegions
	// iterates the ticker slice to decide which regions to draw. Without
	// initialization here, the first frame would render zero regions and
	// the panel would appear blank until the scheduling loop's first cycle.
	if rl.tickRateResolver != nil {
		rl.initRegionTickers()
	}

	// Step 2: Render and flush one complete frame immediately.
	// This guarantees visible content appears on screen before the loop
	// starts sleeping for tick intervals. Without this, the display stays
	// blank for up to one full tick interval (potentially 1000ms) after
	// activation — unacceptable for the user experience on panel startup.
	if rl.tickRateResolver != nil {
		rl.renderDueRegions(time.Now())
	} else {
		rl.renderFrame()
	}
	rl.flush()

	// Step 3: Check for stop after the first frame.
	// Run honors Stop even if it was called during initialization or first
	// render. Without this check, a race between Stop and the scheduling
	// loop could cause one extra unwanted frame cycle.
	select {
	case <-rl.stopCh:
		return
	default:
	}

	// Step 4: Select the scheduling mode.
	// This decision is deferred until after the first frame because both
	// paths (global ticker and per-region deadlines) block indefinitely.
	// Choosing here lets the first frame render unconditionally regardless
	// of which scheduling strategy is active. If this selection happened
	// before the first render, we would need duplicated first-frame logic
	// inside both runPerRegion and runGlobalTicker.
	//
	// Deadline-based scheduling (runPerRegion) solves the problem of modes
	// with different refresh rates coexisting in the same loop: an animation
	// mode at 50ms and a dashboard mode at 1000ms must not force the slower
	// mode to wake at the faster mode's rate. Per-region deadlines let each
	// region sleep independently, avoiding wasted render cycles and CPU
	// usage. Without this, a single global ticker at the fastest rate would
	// re-render idle regions 20× more often than necessary.
	if rl.tickRateResolver != nil {
		rl.runPerRegion()
		return
	}

	// Fall back to global ticker when no resolver is set.
	rl.runGlobalTicker()
}

// runGlobalTicker is the legacy loop using a single global ticker interval.
func (rl *RenderLoop) runGlobalTicker() {
	ticker := time.NewTicker(rl.tick)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopCh:
			return
		case <-ticker.C:
			frameStart := time.Now()

			rl.renderFrame()
			rl.flush()
			rl.processInput()

			// If frame time exceeded tick interval, reset the ticker to avoid
			// drift accumulation. The next frame begins immediately on the next
			// ticker fire (which will already be ready if we're late).
			elapsed := time.Since(frameStart)
			if elapsed >= rl.tick {
				ticker.Reset(rl.tick)
			}
		}
	}
}

// runPerRegion implements per-region deadline-based scheduling.
// It sleeps until the earliest region deadline, renders all due regions,
// flushes once, and processes input.
func (rl *RenderLoop) runPerRegion() {
	for {
		// Compute minimum sleep duration across all region deadlines.
		sleepDur := rl.minSleepDuration()

		if sleepDur > 0 {
			timer := time.NewTimer(sleepDur)
			select {
			case <-rl.stopCh:
				timer.Stop()
				return
			case <-timer.C:
			}
		} else {
			// Check for stop without sleeping.
			select {
			case <-rl.stopCh:
				return
			default:
			}
		}

		now := time.Now()

		// Re-derive tick intervals for regions whose mode has changed.
		rl.rederiveTickIntervals()

		// Render only regions whose deadline has elapsed.
		rl.renderDueRegions(now)
		rl.flush()
		rl.processInput()
	}
}

// minSleepDuration computes the minimum time until the next region deadline fires.
// Returns 0 if any deadline is already past.
func (rl *RenderLoop) minSleepDuration() time.Duration {
	now := time.Now()
	minDur := time.Duration(1<<63 - 1) // max duration

	for _, rt := range rl.regionTickers {
		deadline := rt.lastFire.Add(rt.interval)
		remaining := deadline.Sub(now)
		if remaining <= 0 {
			return 0
		}
		if remaining < minDur {
			minDur = remaining
		}
	}

	if len(rl.regionTickers) == 0 {
		return rl.tick
	}

	return minDur
}

// renderDueRegions renders only regions whose deadline has elapsed (deadline ≤ now),
// then advances each rendered region's deadline by its tick interval.
// On render error or panic, the region's deadline is still advanced.
func (rl *RenderLoop) renderDueRegions(now time.Time) {
	if rl.renderer == nil {
		return
	}

	for _, rt := range rl.regionTickers {
		deadline := rt.lastFire.Add(rt.interval)
		if deadline.After(now) {
			continue // not due yet
		}

		// Render this region (safe from panics).
		rl.safeRender(rt.region)

		// Advance deadline regardless of render success/failure.
		rt.lastFire = rt.lastFire.Add(rt.interval)
	}
}

// rederiveTickIntervals checks each region ticker's mode against the resolver
// and updates the interval if the mode has changed since the ticker was created.
func (rl *RenderLoop) rederiveTickIntervals() {
	for _, rt := range rl.regionTickers {
		newInterval := rl.tickRateResolver.TickInterval(rt.region.CurrentMode())
		if newInterval != rt.interval {
			rt.interval = newInterval
		}
	}
}

// renderFrame renders all Regions in allocation order.
// On render error or panic, the Region is skipped and its surface left unchanged.
func (rl *RenderLoop) renderFrame() {
	if rl.renderer == nil {
		return
	}

	regions := rl.rm.Regions()
	for _, r := range regions {
		rl.safeRender(r)
	}
}

// safeRender calls the renderer for a single Region, recovering from panics.
func (rl *RenderLoop) safeRender(r *Region) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("render loop: panic rendering region %q: %v\n%s", r.Name(), rec, debug.Stack())
		}
	}()

	if err := rl.renderer.Render(r); err != nil {
		log.Printf("render loop: error rendering region %q: %v", r.Name(), err)
	}
}

// flush calls FlushPath.Flush to push all Physical Screens.
func (rl *RenderLoop) flush() {
	if rl.fp == nil {
		return
	}
	if err := rl.fp.Flush(); err != nil {
		log.Printf("render loop: flush error: %v", err)
	}
}

// processInput drains all pending input events and dispatches them to the
// input-active Region.
func (rl *RenderLoop) processInput() {
	if rl.events == nil {
		return
	}

	activeRegion := rl.rm.InputActiveRegion()
	if activeRegion == nil {
		// No input-active region; drain events without dispatching.
		for {
			select {
			case _, ok := <-rl.events:
				if !ok {
					return
				}
			default:
				return
			}
		}
	}

	// Drain all pending events.
	for {
		select {
		case ev, ok := <-rl.events:
			if !ok {
				return
			}
			if rl.dispatcher != nil {
				rl.dispatcher.Dispatch(activeRegion, ev)
			}
		default:
			return
		}
	}
}
