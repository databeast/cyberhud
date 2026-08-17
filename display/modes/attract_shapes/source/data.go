package source

import "github.com/databeast/cyberhud/display/widgets"

// Snapshot is the mode-specific snapshot type consumed by style Build methods.
// Since attract_shapes is sprite-based with no text content, the snapshot
// carries pre-rendered sprites and policy state.
type Snapshot struct {
	Sprites []widgets.Sprite
	Policy  Policy
	IsEink  bool
}
