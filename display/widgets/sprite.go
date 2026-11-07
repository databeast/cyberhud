package widgets

import "image"

// Sprite is a positioned bitmap element used for rendering.
// Modes include Sprites in their view state for compositing onto surfaces.
type Sprite struct {
	Image    image.Image     // Pixel data; nil means skip.
	Position image.Point     // Top-left draw coordinate (used when Bounds is zero).
	Bounds   image.Rectangle // Destination rect; non-zero triggers scaled drawing.
	Label    string          // Accessibility/debug identifier (max 128 chars).

	// Z is the sprite's compositing depth relative to the text/cell plane.
	//
	// The text plane is the fixed datum: it is the first thing drawn at depth 0.
	// Sprites are drawn in ascending Z order around it, declaration order breaking
	// ties within a depth. The default (0) therefore draws a sprite *after* the
	// text — on top of it — which is the overwhelmingly common case (icons, bars,
	// labels decorating a row). Set a negative Z to sink a sprite *beneath* the
	// text: a full-panel gradient or backdrop is Z:-1 so the text reads on top of
	// it. 99% of styles never set this, because their elements do not overlap.
	Z int
}
