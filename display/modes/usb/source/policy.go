package source

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
)

const defaultPollInterval = 500 * time.Millisecond

// Policy controls runtime behavior for USB bench detection.
type Policy struct {
	PollMS          int
	HoldUnpluggedMS int
	HideRootHubs    bool
	Style           string // "default" or "compact"
}

// DefaultPolicy returns baseline USB bench behavior.
func DefaultPolicy() Policy {
	return Policy{
		PollMS:          int(defaultPollInterval / time.Millisecond),
		HoldUnpluggedMS: 0,
		HideRootHubs:    true,
		Style:           "",
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "poll_ms", Type: "int", Summary: "Milliseconds between fallback polling scans.", Default: "500"},
		{Key: "hold_unplugged_ms", Type: "int", Summary: "How long to keep a disconnected device visible (0 keeps indefinitely).", Default: "0"},
		{Key: "hide_root_hubs", Type: "bool", Summary: "Hide USB root hubs from bench identification results.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: ""},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("poll=%d|hold=%d|hide=%t|style=%s", p.PollMS, p.HoldUnpluggedMS, p.HideRootHubs, p.Style)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"poll_ms":           p.PollMS,
		"hold_unplugged_ms": p.HoldUnpluggedMS,
		"hide_root_hubs":    p.HideRootHubs,
		"style":             p.Style,
	}
}
