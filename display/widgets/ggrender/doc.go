// Package ggrender provides 2D vector graphics rendering primitives for the
// cyberhud widget system using the fogleman/gg library.
//
// This package wraps gg.Context behind a Canvas abstraction, exposing shape
// primitives (rectangle, rounded rectangle, circle, line, arc), TrueType text
// rendering with horizontal alignment, and a direct pipeline into the
// widgets.Sprite compositing system.
//
// The package is additive — existing bitmap-font widgets are unaffected. Display
// mode packages import ggrender when they need anti-aliased vector rendering
// (primarily on color panels 240×135 and above), and use it identically to other
// widget packages: call Render(Config), get a *Sprite.
//
// Key design decisions:
//
//   - Wrapper, not raw access: Callers interact with Canvas methods rather than
//     *gg.Context directly, allowing the package to enforce invariants (transparent
//     background, dimension validation, clipping) and keep the gg dependency
//     contained.
//   - Same contract: Following the existing Render(Config) *Sprite + sign(Config)
//     uint64 pattern means gg-based widgets slot into render cache memoization
//     and the panel compositing loop with zero changes to calling code.
//   - Font caching: TrueType font loading is expensive. A package-level cache
//     keyed on (path, pointSize) avoids repeated file I/O and font.Face creation.
//   - Resolution guard: A 16×16 minimum prevents wasted vector computation on
//     tiny or monochrome panels where bitmap widgets are more appropriate.
package ggrender
