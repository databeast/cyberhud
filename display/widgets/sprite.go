package widgets

import "image"

// Sprite is a positioned bitmap element used for rendering.
// Modes include Sprites in their view state for compositing onto surfaces.
type Sprite struct {
	Image    image.Image     // Pixel data; nil means skip.
	Position image.Point     // Top-left draw coordinate (used when Bounds is zero).
	Bounds   image.Rectangle // Destination rect; non-zero triggers scaled drawing.
	Label    string          // Accessibility/debug identifier (max 128 chars).
}
