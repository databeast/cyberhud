package source

import "github.com/databeast/cyberhud/display/widgets"

// Data is the mode-specific snapshot consumed by waveform style Build methods.
type Data struct {
	Sprites []widgets.Sprite
	IsEink  bool
}
