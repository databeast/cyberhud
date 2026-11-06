package storage

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("storage", newInstance)
}

// instance implements displaymodes.ModeInstance for the storage mode.
// It manages the filesystem sampler lifecycle via Activate/Deactivate.
type instance struct {
	displaymodes.PanelHints
}

func newInstance() displaymodes.ModeInstance { return &instance{} }

func (i *instance) ID() string { return "storage" }

// Activate starts the background filesystem sampling goroutine.
func (i *instance) Activate() {
	source.Activate(
		func() int { return GetPolicy().RefreshMS },
		func() bool { return GetPolicy().ShowAll },
	)
}

// Deactivate stops the background filesystem sampling goroutine.
func (i *instance) Deactivate() {
	source.Deactivate()
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := i.Hints()
	if !ok {
		hints, ok = modehints.Current()
		if !ok {
			hints = textlayout.TextHints{}
		}
	}
	return buildView(hints)
}

func (i *instance) RenderCacheKey() uint32 {
	snap := source.CurrentSnapshot()
	hints, _ := i.Hints()
	return renderCacheKey(snap, hints.PixelWidth, hints.PixelHeight)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:       "storage",
		Title:    "Storage",
		Category: "System",
		Summary:  "Filesystem usage per mount with capacity bars and low-space alerts.",
		Order:    55,
		Options:  append(source.Policy{}.Options(), catalog.OptionDefinition{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: allowedStyleNames()}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "storage",
		Summary: "Query or set storage display options.",
		Usage:   "storage [style=<mono-fast-128x64-bars|color-fast-240x240-focus|...>] [font=<auto|font-id>] [refresh_ms=<ms>] [unit=<auto|si|iec>] [show_all=<true|false>]",
		Handle:  HandleCommand,
	})
}
