package source

import (
	"sync"
)

// activeReaderState holds the package-level references to the current tail
// reader and line buffer. These are needed because the command handler
// (package-level) must be able to stop/start the reader when the source
// changes at runtime. The instance lifecycle (Activate/Deactivate) also uses
// these to manage the reader.
var activeReaderState struct {
	sync.Mutex
	reader *TailReader
	buffer *LineBuffer
	scroll *ScrollState
	page   *PageState
}

// setActiveReader stores the reader, buffer, and scroll state references for
// use by the command handler's PostApply hook. Called during Activate().
func SetActiveReader(r *TailReader, buf *LineBuffer) {
	activeReaderState.Lock()
	defer activeReaderState.Unlock()
	activeReaderState.reader = r
	activeReaderState.buffer = buf
}

// setActiveScroll stores the scroll state reference for runtime speed updates.
// Called during Activate() when the surface is classified as fast.
func SetActiveScroll(s *ScrollState) {
	activeReaderState.Lock()
	defer activeReaderState.Unlock()
	activeReaderState.scroll = s
}

// setActivePage stores the page state reference for cache key generation.
// Called during Activate() when the surface is classified as slow.
func SetActivePage(p *PageState) {
	activeReaderState.Lock()
	defer activeReaderState.Unlock()
	activeReaderState.page = p
}

// clearActiveReader removes the reader, buffer, scroll, and page references.
// Called during Deactivate().
func ClearActiveReader() {
	activeReaderState.Lock()
	defer activeReaderState.Unlock()
	activeReaderState.reader = nil
	activeReaderState.buffer = nil
	activeReaderState.scroll = nil
	activeReaderState.page = nil
}

// sourceChangePostApply is the PostApply hook for the pager CmdHandler.
// When the "source" key is among the applied keys, it stops the current
// reader and starts a new one with the updated policy. If the new source is
// empty, the reader is left stopped (neutral status — no connection attempt).
// On failure to open the new source, the reader's retry logic handles it:
// it retries at scan_ms without reverting to the previous source.
//
// When the "scroll_speed" key is among the applied keys, it updates the
// active scrollState's base speed so the new speed applies from the next
// frame without altering the current scroll offset.
func ActiveStateSnapshot() (*TailReader, *LineBuffer, *ScrollState, *PageState) {
	activeReaderState.Lock()
	defer activeReaderState.Unlock()
	return activeReaderState.reader, activeReaderState.buffer, activeReaderState.scroll, activeReaderState.page
}

func RestartActiveReader(p Policy) {
	r, buf, _, _ := ActiveStateSnapshot()
	if r == nil {
		return
	}
	r.Stop()
	if p.Source != "" {
		r.Start(p, buf)
	}
}
