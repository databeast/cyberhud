package source

import "github.com/databeast/cyberhud/display/widgets"

// Snapshot is the mode-specific typed snapshot consumed by style Build methods.
type Snapshot struct {
	Sprites []widgets.Sprite
	Policy  Policy
	IsEink  bool
}
