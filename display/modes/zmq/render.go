package zmq

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// RenderCacheKey returns a deterministic change-detection string (≤512 bytes)
// incorporating the buffer sequence counter and all policy fields. The display
// runtime compares consecutive keys to skip unnecessary re-renders.
func RenderCacheKey() uint32 {
	p := content.GetPolicy()
	seq := content.BufferSeq()
	return region.CalcRegionCacheKey(fmt.Sprintf("%d|%s|%s|%s|%d|%s|%s|%s",
		seq, p.Endpoint, p.SocketType, p.Topic,
		p.MaxLines, p.Style, p.Font, p.JSONFields))
}

// BuildView returns the ZMQ mode view data for the display runtime.
// It populates a style.ViewData with the current buffer snapshot, font
// resolution, and cursor/topRow clamping based on visible row count.
func BuildView(hints textlayout.TextHints) style.ViewData {
	items := content.SnapshotLines()

	// When the buffer is empty, show a placeholder item.
	if len(items) == 0 {
		items = []string{"Waiting for messages..."}
	}

	p := GetPolicy()

	s, reason := style.ResolveStyle(zmqRegistry, hints, "zmq", p.Style)

	// Build the snapshot for style dispatch.
	snap := content.ZMQData{Lines: items}

	// Construct StyleContext boundary and invoke the style's Build method.
	ctx := style.NewStyleContext(hints)
	svd := s.Build(snap, p, ctx)

	// Report style resolution metadata to the registry layer.
	svd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	return svd
}
