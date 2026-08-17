package source

import "github.com/databeast/cyberhud/display/style"

// Snapshot is the mode-specific snapshot consumed by style Build methods.
type Snapshot struct {
	Sprites []style.ViewData
	Policy  Policy
	IsEink  bool
}
