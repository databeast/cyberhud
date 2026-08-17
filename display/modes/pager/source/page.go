package source

import (
	"time"

	"github.com/databeast/cyberhud/display/style"
)

// pagePhase represents the current state of a page transition.
type PagePhase int

const (
	// phaseIdle means the current page is being displayed, no transition in progress.
	PhaseIdle PagePhase = iota

	// phaseFadeOut means the current page is fading out before the next page appears.
	PhaseFadeOut

	// phaseFadeIn means the next page is fading in to replace the previous page.
	PhaseFadeIn
)

// pageState tracks the state of page transitions for slow-display rendering.
type PageState struct {
	currentPage []string      // lines currently displayed
	nextPage    []string      // lines queued for next transition
	phase       PagePhase     // idle | fadeOut | fadeIn
	fadeAlpha   float64       // 0.0–1.0 opacity during transition
	cadence     time.Duration // computed from visible rows × per_line_time
	waitTimer   time.Duration // tracks time waiting for a full page
}

// newPageState creates a pageState with the given cadence.
func NewPageState(visibleRows int, perLineTimeMS int) *PageState {
	return &PageState{
		currentPage: []string{},
		nextPage:    []string{},
		phase:       PhaseIdle,
		fadeAlpha:   1.0,
		cadence:     computePageCadence(visibleRows, perLineTimeMS),
		waitTimer:   0,
	}
}

// computePageCadence calculates the page cadence as max(3000ms, visibleRows × perLineTime).
func computePageCadence(visibleRows int, perLineTimeMS int) time.Duration {
	const minCadenceMS = 3000

	cadenceMS := visibleRows * perLineTimeMS
	if cadenceMS < minCadenceMS {
		cadenceMS = minCadenceMS
	}
	return time.Duration(cadenceMS) * time.Millisecond
}

// tick advances the page state by the given delta duration. It handles
// phase transitions and timing logic for page cadence, fade-out, and fade-in.
//
// Parameters:
//   - delta: elapsed time since last tick
//   - bufferLines: current snapshot of lines available in the buffer
//   - visibleRows: number of rows that fit on the display surface
//   - fadeOutMS: configured fade-out duration in milliseconds
//   - fadeInMS: configured fade-in duration in milliseconds
//   - maxWaitS: maximum wait for full page in seconds
func (ps *PageState) Tick(delta time.Duration, bufferLines []string, visibleRows int, fadeOutMS int, fadeInMS int, maxWaitS int) {
	switch ps.phase {
	case PhaseIdle:
		ps.tickIdle(delta, bufferLines, visibleRows, fadeOutMS, maxWaitS)
	case PhaseFadeOut:
		ps.tickFadeOut(delta, fadeOutMS)
	case PhaseFadeIn:
		ps.tickFadeIn(delta, fadeInMS)
	}
}

func (ps *PageState) CurrentPage() []string { return ps.currentPage }
func (ps *PageState) NextPage() []string    { return ps.nextPage }
func (ps *PageState) Phase() PagePhase      { return ps.phase }
func (ps *PageState) FadeAlpha() float64    { return ps.fadeAlpha }

// tickIdle handles the idle phase where the current page is displayed and we
// are waiting for enough content to schedule a transition.
func (ps *PageState) tickIdle(delta time.Duration, bufferLines []string, visibleRows int, fadeOutMS int, maxWaitS int) {
	available := len(bufferLines)

	// Determine how many new lines are available beyond what's currently displayed.
	// The buffer holds all lines; lines already shown on currentPage are consumed.
	// For the page strategy, we look at what's in the buffer as the next content.
	if available >= visibleRows {
		// Full page available — accumulate cadence time.
		ps.waitTimer += delta
		if ps.waitTimer >= ps.cadence {
			// Schedule transition: prepare next page from buffer.
			ps.nextPage = extractPage(bufferLines, visibleRows)
			ps.phase = PhaseFadeOut
			ps.fadeAlpha = 1.0
			ps.waitTimer = 0
		}
	} else if available > 0 {
		// Partial page — wait up to maxWaitS.
		ps.waitTimer += delta
		maxWait := time.Duration(maxWaitS) * time.Second
		if ps.waitTimer >= maxWait {
			// Timeout with partial content: pad with blank lines.
			ps.nextPage = padPage(bufferLines, visibleRows)
			ps.phase = PhaseFadeOut
			ps.fadeAlpha = 1.0
			ps.waitTimer = 0
		}
	} else {
		// No new text available — retain current page, no transition.
		// Accumulate wait time but do not transition.
		ps.waitTimer += delta
		maxWait := time.Duration(maxWaitS) * time.Second
		if ps.waitTimer >= maxWait {
			// Max wait expired with no text: retain current page, reset timer.
			ps.waitTimer = 0
		}
	}
}

// tickFadeOut advances the fade-out animation. fadeAlpha decreases from 1.0 to 0.0
// over fadeOutMS milliseconds.
func (ps *PageState) tickFadeOut(delta time.Duration, fadeOutMS int) {
	if fadeOutMS <= 0 {
		// Instant fade-out.
		ps.fadeAlpha = 0.0
		ps.transitionToFadeIn()
		return
	}

	fadeOutDuration := time.Duration(fadeOutMS) * time.Millisecond
	decrement := float64(delta) / float64(fadeOutDuration)
	ps.fadeAlpha -= decrement

	if ps.fadeAlpha <= 0.0 {
		ps.fadeAlpha = 0.0
		ps.transitionToFadeIn()
	}
}

// tickFadeIn advances the fade-in animation. fadeAlpha increases from 0.0 to 1.0
// over fadeInMS milliseconds.
func (ps *PageState) tickFadeIn(delta time.Duration, fadeInMS int) {
	if fadeInMS <= 0 {
		// Instant fade-in.
		ps.fadeAlpha = 1.0
		ps.transitionToIdle()
		return
	}

	fadeInDuration := time.Duration(fadeInMS) * time.Millisecond
	increment := float64(delta) / float64(fadeInDuration)
	ps.fadeAlpha += increment

	if ps.fadeAlpha >= 1.0 {
		ps.fadeAlpha = 1.0
		ps.transitionToIdle()
	}
}

// transitionToFadeIn moves from fade-out to fade-in, promoting nextPage to
// currentPage.
func (ps *PageState) transitionToFadeIn() {
	ps.currentPage = ps.nextPage
	ps.nextPage = nil
	ps.phase = PhaseFadeIn
	ps.fadeAlpha = 0.0
}

// transitionToIdle completes the fade-in and returns to idle state.
func (ps *PageState) transitionToIdle() {
	ps.phase = PhaseIdle
	ps.fadeAlpha = 1.0
	ps.waitTimer = 0
}

// extractPage takes the first visibleRows lines from the buffer.
func extractPage(bufferLines []string, visibleRows int) []string {
	if len(bufferLines) >= visibleRows {
		page := make([]string, visibleRows)
		copy(page, bufferLines[:visibleRows])
		return page
	}
	return padPage(bufferLines, visibleRows)
}

// padPage pads the given lines with blank lines to fill visibleRows.
func padPage(lines []string, visibleRows int) []string {
	page := make([]string, visibleRows)
	copy(page, lines)
	// Remaining entries are zero-value empty strings (blank lines).
	return page
}

// renderPageView produces a style.ViewData representing the current page state.
// This is a standalone rendering function that does NOT modify BuildView directly;
// callers are responsible for wiring it into the instance's BuildView method.
//
// Parameters:
//   - ps: the current page state
//   - bufferSnapshot: current lines from the line buffer
//   - layout: computed layout information (visible rows, columns, dimensions)
//
// During idle phase, it renders the current page at full opacity.
// During fadeOut, it renders the current page at decreasing opacity.
// During fadeIn, it renders the new (current) page at increasing opacity.
func RenderPageView(ps *PageState, bufferSnapshot []string, layout Layout) style.ViewData {
	if layout.VisibleRows <= 0 || layout.VisibleColumns <= 0 {
		return style.ViewData{
			Items:  []string{},
			Static: true,
		}
	}

	var displayLines []string
	switch ps.phase {
	case PhaseIdle:
		displayLines = ps.currentPage
	case PhaseFadeOut:
		// During fade-out, show the current (outgoing) page.
		displayLines = ps.currentPage
	case PhaseFadeIn:
		// During fade-in, show the new (incoming) page (already promoted to currentPage).
		displayLines = ps.currentPage
	}

	// Ensure we have exactly visibleRows lines for display.
	items := make([]string, layout.VisibleRows)
	for i := 0; i < layout.VisibleRows && i < len(displayLines); i++ {
		line := displayLines[i]
		// Truncate to visible columns if needed.
		if len(line) > layout.VisibleColumns {
			line = line[:layout.VisibleColumns]
		}
		items[i] = line
	}

	return style.ViewData{
		Items:  items,
		Static: ps.phase == PhaseIdle,
	}
}
