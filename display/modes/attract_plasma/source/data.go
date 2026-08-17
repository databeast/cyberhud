package source

import "github.com/databeast/cyberhud/display/widgets"

// Snapshot is the mode-specific snapshot type consumed by style Build methods.
type Snapshot struct {
	Sprites []widgets.Sprite
	Policy  Policy
	IsEink  bool
}
