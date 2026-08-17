package region

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// ModeFactory is a Region-layer function type that constructs a [ModeInstance]
// for the given mode ID and text hints. It is set on each Region via
// [Region.SetModeFactory] by the runtime layer that can import both the region
// and displaymodes packages (avoiding circular imports).
//
// The id parameter is the normalized mode identifier (lowercase, trimmed).
// The hints parameter provides text layout metrics scoped to the Region's pixel
// dimensions.
//
// Returns (instance, true) on success, or (nil, false) if the mode is unknown.
// The function may panic if the underlying factory panics — callers must recover.
type ModeFactory func(id string, hints textlayout.TextHints) (ModeInstance, bool)

// HintsReceiver is an optional interface a [ModeInstance] may implement to receive
// the hints of the Region hosting it.
//
// # Why injection rather than a global
//
// A mode needs its Region's geometry to lay anything out. Historically it obtained
// that by reading display/region/modehints, a process-wide singleton written by
// whoever activated a mode most recently. That is wrong the moment a second Region
// exists: both modes read the same store, so the last activation wins and the other
// Region lays out for the wrong panel. Per-region data cannot live in process state.
//
// SetMode calls SetPanelHints on the freshly constructed instance, before Activate
// and therefore before any BuildView. An instance that implements this interface
// holds its own Region's hints for its whole lifetime and never needs to consult a
// global.
//
// # Why optional
//
// Making it a required method would change the constructor signature of every
// registered mode at once. As an optional interface, modes migrate individually and
// un-migrated modes keep working through the legacy global. See
// displaymodes.PanelHints for the embeddable implementation, and the clock mode for
// a converted reference. New modes should implement it from the start.
type HintsReceiver interface {
	SetPanelHints(hints textlayout.TextHints)
}

// ModeInstance is the Region-layer interface that the Region uses for its active
// display mode. It defines the lifecycle contract (Activate/Deactivate) and
// rendering contract (BuildView, RenderCacheKey) that each display mode must
// satisfy. This interface is structurally compatible with displaymodes.ModeInstance,
// allowing the region package to avoid importing display/modes (which would
// cause a cycle).
type ModeInstance interface {
	// ID returns the mode's unique identifier string.
	ID() string
	// Activate initializes the mode instance for rendering. Called once after
	// construction when SetMode transitions to this mode.
	Activate()
	// Deactivate tears down the mode instance. Called on the old instance
	// after a successful mode transition, with a 5-second timeout.
	Deactivate()
	// ActionHandler returns the input action handler for this mode.
	ActionHandler() action.Handler
	// BuildView produces the current frame's visual data for the renderer.
	BuildView() style.ViewData
	// RenderCacheKey returns a key used to determine if the frame has changed
	// since the last render, enabling frame-skip optimizations.
	RenderCacheKey() uint32
}

// Region is the Region architectural component: a named, non-overlapping
// rectangular area of the VirtualDisplay that serves as the unit of display mode
// lifecycle, rendering, and input focus. Each Region owns a rendering surface that
// is a zero-copy sub-image of the VirtualDisplay framebuffer, and manages its own
// ModeInstance lifecycle including construction, activation, deactivation, and panic
// recovery.
type Region struct {
	name         string
	bounds       image.Rectangle
	surface      *surface.Surface
	mode         string
	textHints    textlayout.TextHints
	inputFocus   bool
	screens      []ScreenPosition
	panelProduct string       // normalized panel product name for TextHints propagation
	panelPPI     float64      // panel-level PPI for cascading into TextHints
	configPPI    float64      // config-level PPI override; zero means "no override"
	instance     ModeInstance // active mode instance; nil before first SetMode
	modeFactory  ModeFactory  // injected by runtime; nil means legacy (no instance construction)
}

// NewRegion creates a Region with the given name, bounds, and surface.
// The name parameter is the unique identifier for this region. The bounds parameter
// is the rectangle in VirtualDisplay coordinates. The surf parameter is the
// rendering surface (typically a zero-copy sub-image of the VirtualDisplay).
// TextHints are initialized from the surface's bounds using DefaultTextHints.
func NewRegion(name string, bounds image.Rectangle, surf *surface.Surface) *Region {
	return &Region{
		name:      name,
		bounds:    bounds,
		surface:   surf,
		textHints: textlayout.DefaultTextHints(image.Rect(0, 0, bounds.Dx(), bounds.Dy())),
	}
}

// NewRegionWithScreens creates a Region with the given name, bounds, surface, and
// screen positions. The name parameter is the unique identifier. The bounds parameter
// is the rectangle in VirtualDisplay coordinates. The surf parameter is the rendering
// surface. The screens parameter provides physical screen geometry for hardware-aware
// TextHints resolution:
//   - Region entirely within one Physical Screen with TextHintProvider: use that
//     screen's capability flags + Region's PixelWidth/PixelHeight
//   - Region spans multiple Physical Screens: use default capability flags +
//     Region's PixelWidth/PixelHeight
//   - No TextHintProvider available: use DefaultTextHints(Region.Bounds())
//
// The panelProduct parameter is the normalized panel product name propagated into
// TextHints.PanelProduct so display modes can resolve style aliases.
// The panelPPI parameter is the panel-level PPI used for cascading into TextHints.PPI.
// The configPPI parameter is the user-supplied config file PPI override; it sits
// between panelPPI and the 96.0 default in the PPI cascade. Zero means no override.
func NewRegionWithScreens(name string, bounds image.Rectangle, surf *surface.Surface, screens []ScreenPosition, panelProduct string, panelPPI float64, configPPI float64) *Region {
	r := &Region{
		name:         name,
		bounds:       bounds,
		surface:      surf,
		screens:      screens,
		panelProduct: strings.TrimSpace(panelProduct),
		panelPPI:     panelPPI,
		configPPI:    configPPI,
	}
	r.textHints = r.resolveTextHints()
	return r
}

// resolvePPI determines the PPI for a region based on cascading priority:
// Screen.PPI → Panel.PPI → Config.PPI → 96.0 default. The default ensures
// PPI-aware logic is always active for hardware-backed regions.
func resolvePPI(screenPPI, panelPPI, configPPI float64) float64 {
	if screenPPI > 0 {
		return screenPPI
	}
	if panelPPI > 0 {
		return panelPPI
	}
	if configPPI > 0 {
		return configPPI
	}
	return AssumedPPI
}

// AssumedPPI is the pixel density substituted when no layer of the stack reports a
// real one.
//
// # This number is a fiction, and it matters that you know it
//
// tiercatalog's tier targets are specified in millimetres (DefaultTargetsMM:
// colossal is 18mm) because legibility is a physical property, not a pixel count.
// Converting those to pixels requires real pixel density. This constant is what
// gets used when none is available, and 96 DPI is a desktop-monitor figure. Small
// embedded panels are nothing like it: a 1.3-inch 240x240 display is roughly 260
// PPI, so an 18mm request resolves to 68px here when the physically correct answer
// would be about 184px.
//
// The conversion is therefore internally consistent but not physically meaningful
// on most of this project's hardware. Tier sizes currently behave as "absolute
// pixel sizes wearing millimetre clothing."
//
// # The correct fix, and why it is not simply changing this line
//
// The fix is for panel drivers to report real PPI, which they are positioned to do
// since they know their own diagonal and resolution. The plumbing already exists:
// resolvePPI prefers screen, then panel, then config PPI, and only falls back here.
// Populating ScreenPosition.PPI or the panel-level value makes the millimetre
// intent real, per panel, with no change to the tables.
//
// Setting this to zero instead would route tiercatalog through its pixel-target
// fallback (DefaultTargetsPx). That was considered and rejected: those targets are
// also absolute and equally arbitrary, so it swaps one guess for another while
// visibly shrinking text on every large panel — colossal would drop from 68px to
// 40px on an 800px-tall display. Changing it is a one-line edit if that tradeoff is
// ever wanted, but it should be a deliberate visual decision, not a silent one.
//
// tiercatalog does bound every resolved target to a fraction of the region height
// (see maxTierHeightFraction there), so this fiction can no longer request a glyph
// taller than the panel it is drawn on.
const AssumedPPI = 96.0

// resolveTextHints determines the appropriate TextHints for this Region based on
// how it relates to the underlying Physical Screens.
// All paths use the full bounds dimensions for PixelWidth/PixelHeight.
// PanelProduct and ScreenName are set on all paths before returning.
func (r *Region) resolveTextHints() textlayout.TextHints {
	width := r.bounds.Dx()
	height := r.bounds.Dy()

	// Branch 1: No screens known — the Region was created without hardware context
	// (e.g., in tests or headless mode). Default hints assume standard glyph metrics
	// and full scrolling capability. Without this fallback, modes would receive a
	// zero-valued TextHints struct and produce no visible text output.
	if len(r.screens) == 0 {
		hints := buildCatalogIfNeeded(textlayout.DefaultTextHints(image.Rect(0, 0, width, height)), 0)
		return r.applyProductIdentity(hints, nil)
	}

	// Determine which Physical Screens geometrically overlap this Region. Only
	// overlapping screens can contribute hardware-specific capabilities (e.g., an
	// e-ink screen disabling smooth scroll). Non-overlapping screens are irrelevant
	// because the Region will never be flushed to them.
	var overlapping []ScreenPosition
	for _, s := range r.screens {
		if !r.bounds.Intersect(s.Bounds).Empty() {
			overlapping = append(overlapping, s)
		}
	}

	if len(overlapping) == 0 {
		// No geometric overlap — same as having no screens. Falling through to the
		// single-screen or multi-screen branches would produce incorrect results
		// because there is no hardware to derive capabilities from.
		hints := buildCatalogIfNeeded(textlayout.DefaultTextHints(image.Rect(0, 0, width, height)), 0)
		return r.applyProductIdentity(hints, nil)
	}

	// Determine if the Region fits entirely within a single Physical Screen. This
	// distinction matters because a single containing screen can provide authoritative
	// hardware capability flags (e.g., an e-ink display may disable horizontal scroll).
	// If we skipped this check and always used defaults, modes running on specialty
	// hardware would ignore that hardware's constraints.
	var containingScreen *ScreenPosition
	for i := range overlapping {
		if r.bounds.In(overlapping[i].Bounds) {
			containingScreen = &overlapping[i]
			break
		}
	}

	// Branch 2: Region is entirely within one Physical Screen — use that screen's
	// HintProvider to obtain hardware-specific capability flags while overriding
	// PixelWidth/PixelHeight with the Region's own dimensions. If we used the
	// screen's pixel dimensions instead, modes would miscalculate text layout for
	// regions smaller than the full screen.
	if containingScreen != nil {
		ppi := resolvePPI(containingScreen.PPI, r.panelPPI, r.configPPI)
		if containingScreen.HintProvider != nil {
			hints := containingScreen.HintProvider()
			hints.PixelWidth = width
			hints.PixelHeight = height
			hints.PPI = ppi
			// Normalize fills any zero-valued glyph metrics from defaults, preventing
			// division-by-zero in text layout calculations.
			hints = textlayout.Normalize(hints, image.Rect(0, 0, width, height))
			return r.applyProductIdentity(buildCatalogIfNeeded(hints, ppi), containingScreen)
		}
		// Screen has no HintProvider — fall back to defaults. This avoids a nil
		// dereference and still produces usable text metrics.
		hints := textlayout.DefaultTextHints(image.Rect(0, 0, width, height))
		hints.PPI = ppi
		hints = buildCatalogIfNeeded(hints, ppi)
		return r.applyProductIdentity(hints, containingScreen)
	}

	// Branch 3: Region spans multiple Physical Screens — no single screen's
	// capability flags are authoritative (e.g., one screen may be e-ink and another
	// LCD). We use conservative defaults that assume full scrolling and standard
	// glyph metrics. Picking one screen's flags arbitrarily would produce incorrect
	// behavior on the other screens the Region spans.
	// For multi-screen spanning, PPI falls back to panel-level with 96.0 default.
	ppi := resolvePPI(0, r.panelPPI, r.configPPI)
	hints := buildCatalogIfNeeded(textlayout.TextHints{
		PixelWidth:               width,
		PixelHeight:              height,
		GlyphWidth:               textlayout.GlyphWidth,
		GlyphHeight:              textlayout.GlyphHeight,
		GlyphAdvance:             textlayout.GlyphAdvance,
		RowHeight:                textlayout.RowHeight,
		SupportsVerticalScroll:   true,
		SupportsHorizontalScroll: true,
		SupportsAutoScroll:       true,
		PreferEventRefresh:       false,
		DefaultTickerDirection:   textlayout.TickerDirectionVertical,
		DefaultLineMode:          textlayout.LineModeTruncate,
		PPI:                      ppi,
	}, ppi)
	return r.applyProductIdentity(hints, nil)
}

// applyProductIdentity sets PanelProduct and ScreenName on hints before returning.
// PanelProduct is the normalized panel product name stored on the Region.
// ScreenName is derived from the containing screen's Name field when the panel has
// multiple screens; it is empty for single-screen panels.
func (r *Region) applyProductIdentity(hints textlayout.TextHints, containingScreen *ScreenPosition) textlayout.TextHints {
	hints.PanelProduct = r.panelProduct

	// ScreenName is set only for multi-screen panels. For single-screen panels
	// (exactly one screen in the topology), ScreenName remains empty because
	// there is no disambiguation needed for alias compound key resolution.
	if len(r.screens) > 1 && containingScreen != nil {
		hints.ScreenName = containingScreen.Name
	}

	return hints
}

// Name returns the Region's unique identifier within the RegionManager.
func (r *Region) Name() string {
	return r.name
}

// Bounds returns the Region's bounding rectangle in VirtualDisplay coordinates.
// This is the absolute position within the VirtualDisplay framebuffer, not local
// coordinates.
func (r *Region) Bounds() image.Rectangle {
	return r.bounds
}

// Surface returns the Region's rendering surface with local coordinate origin at
// (0,0). The surface is a zero-copy sub-image of the VirtualDisplay framebuffer,
// so pixel writes are immediately visible in the parent buffer.
func (r *Region) Surface() *surface.Surface {
	return r.surface
}

// TextHints returns text layout hints scoped to this Region's pixel dimensions
// and hardware capabilities. These hints inform display modes how to lay out text
// content (glyph metrics, scrolling support, ticker direction) for this Region.
func (r *Region) TextHints() textlayout.TextHints {
	return r.textHints
}

// SetModeFactory registers factory as the ModeFactory used by SetMode to construct
// ModeInstance values for this Region.
//
// The factory parameter is a function that receives a mode ID and text hints, returning
// a ModeInstance and a boolean indicating whether the mode was recognized. SetModeFactory
// must be called before SetMode — if SetMode is invoked without a factory wired, it
// panics to surface the misconfiguration immediately at startup.
//
// Typically this is called by the runtime layer that can import both the region and
// displaymodes packages, avoiding circular imports between them.
func (r *Region) SetModeFactory(factory ModeFactory) {
	r.modeFactory = factory
}

// SetMode transitions this Region to a new display mode identified by modeID.
// It returns a non-nil error if the mode change fails (factory panic, unknown mode,
// or activation panic); on error the previous mode remains active.
//
// SetModeFactory must be called before SetMode — if no ModeFactory has been wired,
// SetMode panics to surface the misconfiguration immediately at startup.
//
// SetMode performs the following ordered lifecycle steps:
//
//  1. Normalize modeID — lowercase and trim whitespace.
//  2. Construct new instance via ModeFactory with panic recovery — if the factory
//     panics, the error is logged, the previous mode is retained, and the panic
//     is returned as an error.
//  3. Reject unknown mode — if the factory returns found=false, return an error
//     without touching the current state.
//  4. Clear surface — fill the Region's surface with opaque black to avoid visual
//     artifacts from the previous mode.
//  5. Activate new instance with panic recovery — call Activate on the new
//     ModeInstance. If Activate panics, attempt Deactivate on the new instance
//     (also wrapped in recovery), discard the new instance, retain the previous
//     mode, and return the panic as an error.
//  6. Deactivate old instance with a 5-second timeout — run Deactivate on the
//     previous ModeInstance in a goroutine; if it does not complete within 5 seconds,
//     log a timeout warning and proceed. The new instance is stored and modeID is
//     updated only after successful activation.
func (r *Region) SetMode(modeID string) error {
	// Step 1 — Normalize: lowercase and trim the mode ID before any lookup.
	// This must be first so that all subsequent logic uses a canonical key.
	// If skipped, lookups become case-sensitive and whitespace-sensitive, causing
	// spurious "mode not registered" errors for modes that differ only in casing.
	modeID = strings.ToLower(strings.TrimSpace(modeID))

	if r.modeFactory == nil {
		panic(fmt.Sprintf("region %q: SetMode(%q) called but no ModeFactory wired — call SetModeFactory during startup", r.name, modeID))
	}

	// Step 2 — Construct: call the ModeFactory to build the new ModeInstance.
	// Construction must happen before any mutation of Region state (clear, deactivate)
	// so that if the factory fails, the current mode is undisturbed and the user sees
	// no visual glitch. Without this ordering, clearing the surface first would leave
	// the display black with no mode active on factory failure.
	//
	// Panic recovery wraps the factory call because mode constructors are user-supplied
	// code (display-mode packages) that may panic on invalid config or nil dereferences.
	// Without recovery, a single buggy mode would crash the entire process rather than
	// gracefully retaining the previous mode and reporting the error.
	var newInstance ModeInstance
	var found bool
	var constructErr error

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				constructErr = fmt.Errorf("region %q: ModeFactory panicked for mode %q: %v", r.name, modeID, rec)
			}
		}()
		newInstance, found = r.modeFactory(modeID, r.textHints)
	}()

	if constructErr != nil {
		// Factory panicked — retain previous instance, return error.
		return constructErr
	}

	// Step 3 — Reject unknown mode: if the factory does not recognize the mode ID,
	// bail out before any side effects. This gate must follow construction (which
	// determines "found") and precede surface clearing. If this check were skipped or
	// moved after the clear, an unrecognized mode request would wipe the display for
	// no reason, leaving the region in a visually broken state with no replacement mode.
	if !found {
		return fmt.Errorf("region %q: mode %q not registered", r.name, modeID)
	}

	// Step 3b — Inject this Region's hints into the new instance.
	//
	// This must precede Activate, and therefore any BuildView, so the mode never
	// observes a moment where it has no geometry. Injecting per instance is what
	// makes two Regions on two panels independent; see HintsReceiver for why the
	// previous global store could not be.
	//
	// Modes that have not adopted HintsReceiver are skipped silently and continue to
	// read the legacy global. That is not logged: most modes are expected to be
	// mid-migration, and a per-activation warning for each would be noise.
	if hr, ok := newInstance.(HintsReceiver); ok {
		hr.SetPanelHints(r.textHints)
	}

	// Step 4 — Clear surface: fill with opaque black to remove visual artifacts from
	// the previous mode before the new mode renders its first frame. This must happen
	// after construction succeeds (so we know we have a valid replacement) and before
	// activation (so the new mode starts with a clean canvas). Clearing after activation
	// would overwrite the new mode's first rendered frame; clearing before construction
	// would leave a black surface if construction fails.
	r.surface.Clear(color.RGBA{0, 0, 0, 255})

	// Step 5 — Activate: call Activate() on the new ModeInstance to let it set up
	// rendering state (goroutines, timers, GPU resources). Activation must follow
	// the surface clear so the mode starts with a clean framebuffer, and must precede
	// old-instance deactivation so that if activation fails we can fall back to the
	// old mode without having torn it down.
	//
	// Panic recovery wraps Activate because mode init code is user-supplied and may
	// panic (e.g., nil channel, out-of-bounds slice). Without recovery, a crashing
	// Activate would kill the process. With recovery, we discard the broken new
	// instance, attempt cleanup deactivation on it (also wrapped to tolerate panics
	// during teardown), retain the previous mode, and return the error.
	var activateErr error
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				activateErr = fmt.Errorf("region %q: Activate panic for mode %q: %v", r.name, modeID, rec)
				log.Printf("%s", activateErr)
				// Attempt to deactivate the partially-initialized instance
				// in a separate recover block so a secondary panic during
				// cleanup does not mask the original error or crash the process.
				func() {
					defer func() {
						if rec2 := recover(); rec2 != nil {
							log.Printf("region %q: Deactivate panic during cleanup for mode %q: %v", r.name, modeID, rec2)
						}
					}()
					newInstance.Deactivate()
				}()
			}
		}()
		newInstance.Activate()
	}()

	if activateErr != nil {
		// Activate panicked — discard new instance, retain previous, return error.
		return activateErr
	}

	// Step 6 — Deactivate old instance: tear down the previous mode only AFTER the
	// new instance is fully activated. This ordering guarantees that we never leave
	// the region with no active mode — if we deactivated old first and activation
	// then failed, the region would have no instance at all.
	//
	// Deactivation runs in a goroutine with a 5-second timeout because mode teardown
	// may block on I/O, channel drains, or network calls. Without the timeout, a
	// hung Deactivate would block the mode-switch path indefinitely, freezing the
	// display for the user.
	old := r.instance
	if old != nil {
		done := make(chan struct{})
		go func() {
			defer close(done)
			old.Deactivate()
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			log.Printf("region %q: Deactivate timed out for mode %q", r.name, old.ID())
		}
	}

	// Commit: store the new instance and mode ID only after all steps succeed.
	// This is the single point of state mutation, ensuring that any earlier failure
	// leaves r.instance and r.mode unchanged.
	r.instance = newInstance
	r.mode = modeID

	return nil
}

// CurrentMode returns the Region's active display mode ID as a lowercase string.
// Returns an empty string if no mode has been set.
func (r *Region) CurrentMode() string {
	return r.mode
}

// Instance returns the active ModeInstance for this Region, or nil if no mode has
// been set via SetMode. Used by the RenderLoop to access the mode's BuildView and
// RenderCacheKey methods during rendering.
func (r *Region) Instance() ModeInstance {
	return r.instance
}

// UpdateMode updates the Region's mode ID field without clearing the surface or
// stopping any instance. The modeID parameter is the new mode identifier to record.
// Used by the renderer to keep the Region's mode in sync with the active
// ModeInstance after SyncMode picks up external mode changes.
func (r *Region) UpdateMode(modeID string) {
	r.mode = modeID
}

// HasInputFocus reports whether this Region currently receives input events.
// Only one Region at a time has input focus within a RegionManager.
func (r *Region) HasInputFocus() bool {
	return r.inputFocus
}

// SetInputFocus sets whether this Region receives input events. The focus parameter
// controls the input state: true grants focus, false removes it. The RegionManager
// coordinates focus so that only one Region is active at a time.
func (r *Region) SetInputFocus(focus bool) {
	r.inputFocus = focus
}

// LegibilityMinChars is the number of monospace characters the display system
// requires to fit across a region before it considers text legible at a given tier.
//
// This was previously an unexplained literal 10 buried inside
// buildCatalogIfNeeded. It is a legibility policy, not an implementation detail:
// it is the single knob that decides how large a font any region may use, since
// tiercatalog filters candidate fonts by PixelWidth/MinChars. Naming it here makes
// it findable and reviewable.
//
// It is deliberately not configurable per panel yet. If that becomes necessary,
// thread it through ScreenPosition rather than reintroducing a literal at the call
// site.
const LegibilityMinChars = 10

// buildCatalogIfNeeded populates hints.Catalog via tiercatalog.Build if not already set.
// The ppi parameter is passed through to tiercatalog.Params.PPI for glyph height scaling.
//
// Best-effort construction is requested via AllowRelaxedMinChars. The reason is
// worth stating plainly, because the previous behaviour looked like graceful
// degradation and was not:
//
// When no registered font is narrow enough to fit LegibilityMinChars characters —
// a 32px-wide region, for instance, where even the narrowest 6px-advance font
// manages only 5 — Build used to fail, this function logged and returned an
// unpopulated catalog, and every style then substituted textlayout's 6x10 default
// metrics. Those metrics belong to no registered font, so each style performed its
// centring and row-fitting arithmetic against a font the renderer would never draw
// with. The visible result was a few pixels of text jammed off-centre.
//
// Relaxing the character floor instead yields a catalog whose entries carry the
// true metrics of a real, narrow font. The text is small because the region is
// small, which is honest, and every downstream calculation stays consistent with
// what is actually drawn. Callers that need to know this happened can consult
// Catalog.Relaxed.
func buildCatalogIfNeeded(hints textlayout.TextHints, ppi float64) textlayout.TextHints {
	if hints.Catalog.PixelWidth() > 0 {
		return hints // already populated (manual Build or HintProvider)
	}
	if hints.PixelWidth <= 0 || hints.PixelHeight <= 0 {
		return hints // invalid dimensions — nothing to build
	}
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:           hints.PixelWidth,
		PixelHeight:          hints.PixelHeight,
		MinChars:             LegibilityMinChars,
		PPI:                  ppi,
		AllowRelaxedMinChars: true,
	})
	if err != nil {
		// With AllowRelaxedMinChars set, the only remaining failure is an empty font
		// registry, which indicates a broken build rather than a small panel. Leave
		// the catalog zero; Catalog.Entry still yields real default-face metrics for
		// any style that asks.
		log.Printf("region: tiercatalog build failed for %dx%d @ %.1f PPI: %v (styles will fall back to the default face)",
			hints.PixelWidth, hints.PixelHeight, ppi, err)
		return hints
	}
	if cat.Relaxed() {
		log.Printf("region: %dx%d is too narrow for %d characters; relaxed to %d (text will be small but correctly measured)",
			hints.PixelWidth, hints.PixelHeight, cat.RequestedMinChars(), cat.MinChars())
	}
	hints.Catalog = cat
	return applyBaselineGlyphMetrics(hints, cat)
}

// applyBaselineGlyphMetrics overwrites the glyph metric fields of hints with those
// of the catalog's normal tier.
//
// # Who owns font metrics
//
// A panel has pixels, a physical size and a colour capability. It does not have a
// font. Yet every driver in hardware/driver populated TextHints.GlyphWidth,
// GlyphHeight, GlyphAdvance and RowHeight, all of them with the same hardcoded
// 5x8/6/10 constants, as though the panel came with a typeface. Font choice is a
// Surface/Region concern: it depends on how much room the region has and which faces
// are registered, neither of which hardware knows.
//
// Those four fields cannot simply be deleted — around sixty call sites read them, and
// they are the defaults LayoutCalculator falls back to. So ownership moves instead:
// drivers no longer assert them, and the Region fills them here from the tier catalog
// it just built for this region's real dimensions.
//
// # Why the normal tier
//
// TierNormal is the baseline "body text" size for the region, so it is the right
// answer to "what are the glyph metrics here" for any caller that does not ask for a
// specific tier. It also names a real registered face, which the previous constants
// did not necessarily correspond to on any given panel.
//
// # The invariant this maintains
//
// RegionRenderer sets the surface's baseline face from this same catalog entry before
// drawing (see setBaselineFont). That is what makes these metrics true rather than
// merely plausible: a mode measuring with hints.GlyphAdvance and a renderer drawing
// without per-row font IDs now agree on one real face. Breaking that pairing
// reintroduces the measure-with-one-font, draw-with-another defect that this system
// has already suffered from twice.
func applyBaselineGlyphMetrics(hints textlayout.TextHints, cat tiercatalog.Catalog) textlayout.TextHints {
	entry := cat.Entry(tiercatalog.TierNormal)
	if entry.GlyphAdvance <= 0 || entry.RowHeight <= 0 {
		return hints // nothing usable; leave whatever Normalize supplied
	}
	hints.GlyphWidth = entry.GlyphWidth
	hints.GlyphHeight = entry.GlyphHeight
	hints.GlyphAdvance = entry.GlyphAdvance
	hints.RowHeight = entry.RowHeight
	return hints
}
