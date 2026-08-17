package pager

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/frameclock"
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/modes/pager/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("pager", newInstance)

	catalog.Register(catalog.Definition{
		ID:      "pager",
		Title:   "Pager",
		Scope:   "any",
		Summary: "Tails a file, pipe, or socket and presents text with smooth scroll or page transitions.",
		Order:   90,
		Options: append(source.Policy{}.Options(), catalog.OptionDefinition{Key: "style", Type: "string", Summary: "Visual style name or empty for default.", Default: "", Allowed: registeredStyleNames()}),
	})
}

// maxTickElapsed caps the per-frame advancement to prevent visual jumps
// after long pauses (e.g., debugger breakpoints, system suspend).
const maxTickElapsed = 80 * time.Millisecond

// instance implements displaymodes.ModeInstance for the pager mode.
type instance struct {
	// PanelHints holds the hosting Region's text hints, injected by
	// Region.SetMode before Activate and before any BuildView call. Embedding
	// it supplies SetPanelHints, which satisfies region.HintsReceiver.
	//
	// Previously this was a private field that nothing ever assigned, so the
	// pager always laid out against zero-valued hints: ComputeLayout returned
	// zero visible rows and columns and BuildView returned an empty static
	// view, meaning the mode rendered nothing on any panel.
	displaymodes.PanelHints

	// lastSurface tracks the surface classification from the previous BuildView
	// call. When it changes we switch the rendering strategy (scroll ↔ page)
	// without discarding buffer content.
	lastSurface styles.SurfaceClass

	// lastTick records when BuildView last advanced animation state.
	// Used to compute frame delta for scroll advancement and page transitions.
	lastTick time.Time
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "pager" }

// Activate creates the buffer, reader, and scroll/page state based on the
// surface classification, then registers them with the package-level active
// state and starts the reader if a source is configured.
func (i *instance) Activate() {
	pol := GetPolicy()

	// Create shared buffer.
	buf := source.NewLineBuffer(pol.MaxLines)

	// Create reader (not started yet).
	reader := &source.TailReader{}

	// Classify surface and create appropriate rendering state.
	hints, _ := i.Hints()
	sc := styles.ClassifySurface(hints)
	i.lastSurface = sc

	// Register the shared reader and buffer.
	source.SetActiveReader(reader, buf)

	switch sc {
	case styles.SurfaceFast:
		// Smooth scroll: register tick rate and create scroll state.
		region.RegisterTickRate("pager", tickProvider{})
		scroll := source.NewScrollState(pol.ScrollSpeed)
		source.SetActiveScroll(scroll)
	case styles.SurfaceSlow:
		// Page transition: create page state with cadence based on layout.
		layout := styles.ComputeLayout(style.NewStyleContext(hints))
		page := source.NewPageState(layout.VisibleRows, pol.LineTimeMS)
		source.SetActivePage(page)
	}

	// Start the reader if source is configured.
	if pol.Source != "" {
		reader.Start(pol, buf)
	}
}

// Deactivate stops the reader, clears active state, and releases resources.
func (i *instance) Deactivate() {
	reader, _, _, _ := source.ActiveStateSnapshot()

	if reader != nil {
		reader.Stop()
	}
	source.ClearActiveReader()
}

func (i *instance) ActionHandler() action.Handler { return nil }

// BuildView returns rendered pager output based on the current strategy.
//
// Strategy dispatch:
//  1. If source is empty/unset → return neutral status message
//  2. Get panel hints (default to fast if unavailable)
//  3. Classify surface → surfaceFast or surfaceSlow
//  4. Compute layout; if zero visible rows/columns → return empty static view
//  5. Get buffer snapshot
//  6. For surfaceFast: call renderSmoothScroll
//  7. For surfaceSlow: call renderPageView
//
// If the surface classification has changed since the last call, the rendering
// strategy is switched without discarding buffer content (the buffer is shared).
func (i *instance) BuildView() style.ViewData {
	pol := GetPolicy()

	// No source configured → neutral status.
	if pol.Source == "" {
		return buildNoSourceView()
	}

	// Determine the effective hints for this frame.
	hints, ok := i.Hints()
	if !ok {
		return style.ViewData{Items: []string{"[pager] no panel hints available"}, Static: true}
	}

	// Classify surface.
	sc := styles.ClassifySurface(hints)

	// Handle strategy switch: if surface classification changed since last
	// BuildView, swap the active rendering state without discarding buffer.
	if sc != i.lastSurface {
		i.switchStrategy(sc, hints, pol)
	}

	// Compute layout from current hints.
	layout := styles.ComputeLayout(style.NewStyleContext(hints))

	// Zero visible area → empty static view.
	if layout.VisibleRows <= 0 || layout.VisibleColumns <= 0 {
		return style.ViewData{
			Items:  []string{},
			Static: true,
		}
	}

	// Get buffer snapshot.
	_, buf, scroll, page := source.ActiveStateSnapshot()

	if buf == nil {
		return style.ViewData{
			Items:  []string{},
			Static: true,
		}
	}

	lines := buf.Snapshot()

	// Compute elapsed time since last frame and advance animation state.
	var elapsed time.Duration
	now := frameclock.Now()
	if !i.lastTick.IsZero() {
		elapsed = now.Sub(i.lastTick)
		if elapsed > maxTickElapsed {
			elapsed = maxTickElapsed
		}
	}
	i.lastTick = now

	// Dispatch to the appropriate rendering strategy.
	switch sc {
	case styles.SurfaceFast:
		if scroll == nil {
			return style.ViewData{Items: []string{}, Static: true}
		}
		// Advance scroll state: adapt velocity then advance offset.
		scroll.AdaptVelocity(len(lines), layout.VisibleRows)
		if elapsed > 0 {
			scroll.Advance(elapsed.Seconds(), layout.RowHeight)
		}
		vd := source.RenderSmoothScroll(scroll, lines, layout)

		return vd
	case styles.SurfaceSlow:
		if page == nil {
			return style.ViewData{Items: []string{}, Static: true}
		}
		// Advance page transition state.
		if elapsed > 0 {
			page.Tick(elapsed, lines, layout.VisibleRows, pol.FadeOutMS, pol.FadeInMS, pol.MaxWaitS)
		}
		// Render via sprite-based slow-refresh path (produces image directly).
		seq := buf.Seq()
		return buildSlowPageView(hints, layout, page.CurrentPage(), page.NextPage(), page.Phase(), page.FadeAlpha(), seq, pol)
	default:
		return style.ViewData{Items: []string{}, Static: true}
	}
}

// switchStrategy handles a change in surface classification between consecutive
// BuildView calls. It swaps the active scroll/page state without discarding
// buffer content — the buffer is shared, only the rendering strategy changes.
func (i *instance) switchStrategy(newClass styles.SurfaceClass, hints textlayout.TextHints, pol source.Policy) {
	i.lastSurface = newClass

	switch newClass {
	case styles.SurfaceFast:
		// Switching to smooth scroll: clear page state, create scroll state.
		source.SetActivePage(nil)
		scroll := source.NewScrollState(pol.ScrollSpeed)
		source.SetActiveScroll(scroll)
		region.RegisterTickRate("pager", tickProvider{})
	case styles.SurfaceSlow:
		// Switching to page transition: clear scroll state, create page state.
		source.SetActiveScroll(nil)
		layout := styles.ComputeLayout(style.NewStyleContext(hints))
		page := source.NewPageState(layout.VisibleRows, pol.LineTimeMS)
		source.SetActivePage(page)
	}
}

// buildNoSourceView returns a ViewData with a neutral status message indicating
// no data source is configured. No tier or font is declared; the renderer's
// text-fit fallback will select the largest font that fits the message on the
// panel without truncation.
func buildNoSourceView() style.ViewData {
	return style.ViewData{
		Items:  []string{"[pager] no source"},
		Static: true,
	}
}

// RenderCacheKey returns a deterministic string that changes whenever the
// visual output of the pager would change. It incorporates:
//   - Buffer sequence number (changes on every ingest)
//   - Scroll offset (smooth scroll) or page phase + quantized fade alpha (page transition)
//   - Policy fingerprint (changes on any policy field update)
//
// Format for scroll strategy: pager:<seq>:<scrollOffset>:<policyHash>
// Format for page strategy:   pager:<seq>:<phase>:<quantizedAlpha>:<policyHash>
//
// Maximum length is 256 bytes; the policy hash portion is truncated if exceeded.
func (i *instance) RenderCacheKey() uint32 {
	_, buf, scroll, page := source.ActiveStateSnapshot()

	// When the mode is not yet fully wired (no buffer), return a safe stable key.
	if buf == nil {
		key := "pager:0:0:" + truncatePolicyHash(policyFingerprint(GetPolicy()), "pager:0:0:")
		return region.CalcRegionCacheKey(key)
	}

	seq := buf.Seq()
	pHash := policyFingerprint(GetPolicy())

	var key string
	if scroll != nil {
		// Smooth scroll strategy: use scroll offset.
		prefix := fmt.Sprintf("pager:%d:%d:", seq, scroll.OffsetPx())
		key = prefix + truncatePolicyHash(pHash, prefix)
	} else if page != nil {
		// Page transition strategy: use phase + quantized fade alpha.
		// fadeAlpha quantization: multiply by 100 and truncate to int.
		quantizedAlpha := int(page.FadeAlpha() * 100)
		prefix := fmt.Sprintf("pager:%d:%d:%d:", seq, page.Phase(), quantizedAlpha)
		key = prefix + truncatePolicyHash(pHash, prefix)
	} else {
		// Neither scroll nor page state active — fallback with just seq and policy.
		prefix := fmt.Sprintf("pager:%d:0:", seq)
		key = prefix + truncatePolicyHash(pHash, prefix)
	}

	return region.CalcRegionCacheKey(key)
}

// maxCacheKeyBytes is the maximum allowed length for the render cache key.
const maxCacheKeyBytes = 256

// truncatePolicyHash truncates the policy hash so that prefix + hash ≤ 256 bytes.
func truncatePolicyHash(policyHash string, prefix string) string {
	available := maxCacheKeyBytes - len(prefix)
	if available <= 0 {
		return ""
	}
	if len(policyHash) <= available {
		return policyHash
	}
	return policyHash[:available]
}
