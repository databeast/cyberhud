package displaymodes

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

// ModeInstance is the per-region runtime object for an active display mode.
// Constructed once per mode activation; destroyed when the mode is deactivated.
// The Region hosts the instance and the Renderer calls it directly.
type ModeInstance interface {
	// ID returns the mode identifier (e.g., "clock", "dashboard").
	// The returned string is non-empty, lowercase, and matches the ID used
	// at factory registration.
	ID() string

	// Activate is called once when the mode becomes active on a region.
	// Modes with background goroutines (serial, thermal, demo) start them here.
	// Must be called before any BuildView invocation.
	Activate()

	// Deactivate is called once when the mode is being replaced or the region
	// shuts down. Modes with background goroutines stop them here.
	// No BuildView calls will occur after Deactivate returns.
	Deactivate()

	// ActionHandler returns the input handler for this mode instance.
	// Returns nil if the mode does not respond to input.
	ActionHandler() action.Handler

	// BuildView produces the rendering output for the current frame.
	//
	// The mode is responsible for layout; the renderer only draws. Everything the
	// renderer needs about geometry must therefore be present in the returned
	// ViewData — see the layout contract documented on style.ViewData for why, and
	// for the class of bug that arises when a field is omitted and the renderer
	// re-derives it.
	//
	// Row budget flows from the mode to the renderer, not the reverse. BuildView
	// takes no arguments deliberately: a style solves its own fit against real
	// per-row font metrics and reports the answer as ViewData.VisibleCount. An
	// earlier revision of this comment described a maxVisible parameter, which the
	// signature never had.
	//
	// Font resolution is not the mode's job. A mode declares per-row intent in
	// ViewData.Tiers and leaves ViewData.FontIDs nil; the renderer resolves tiers
	// against the region's catalog. A mode may set FontIDs to override that, and
	// the renderer will not overwrite a non-nil value. (This comment previously
	// claimed the opposite — that fonts were fully resolved before returning and
	// nothing post-processed them. That expectation is what left Tiers and FontIDs
	// permanently disconnected, since every mode was supposed to bridge them
	// privately and none did.)
	BuildView() style.ViewData

	// RenderCacheKey returns a change-detection fingerprint.
	// Identical return values across consecutive calls indicate no visual
	// change has occurred, allowing the renderer to skip a redraw.
	RenderCacheKey() uint32
}

// ModeFactory constructs a ModeInstance.
//
// It takes no parameters. Modes that need hardware access (scanner, GPIO manager)
// import those singletons directly.
//
// Region geometry is not a construction parameter. It arrives immediately after
// construction: Region.SetMode calls SetPanelHints on the new instance, before
// Activate and before any BuildView, for instances that implement
// region.HintsReceiver. Embed [PanelHints] to get that method and read the values
// back through its Hints accessor.
//
// (An earlier version of this comment claimed "TextHints is the sole construction
// parameter", which the signature never supported. Modes consequently reached for
// the process-wide modehints store instead, which is incorrect for more than one
// Region — see PanelHints for why, and for the migration state.)
type ModeFactory func() ModeInstance
