package image

import (
	"github.com/databeast/cyberhud/display/catalog"
)

// imageSnapshotter implements catalog.PolicySnapshotter for the image mode.
type imageSnapshotter struct{}

// SnapshotPolicy returns the current image policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (imageSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := PolicySnapshot()
	return map[string]interface{}{
		"fit":   p.Fit,
		"style": p.Style,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (imageSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["fit"]; ok {
		if s, ok := v.(string); ok {
			p.Fit = s
		}
	}
	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("image", imageSnapshotter{})
}
