// Package svgrender provides SVG-powered animated rendering primitives for the
// cyberhud widget system using the srwiley/oksvg and srwiley/rasterx libraries.
//
// This package rasterizes SVG content into *widgets.Sprite objects, exposing
// a Canvas abstraction for SVG-to-bitmap conversion, an Animator for
// frame-based SVG animation, and a direct pipeline into the widgets.Sprite
// compositing system.
//
// The package follows the same Render(Config) *Sprite + Sign(Config) uint64
// contract established by ggrender, enabling SVG-based widgets to participate
// in render cache memoization and the panel compositing loop without changes
// to calling code.
//
// Key design decisions:
//
//   - Wrapper, not raw access: Callers interact with Canvas methods rather than
//     oksvg/rasterx types directly, allowing the package to enforce invariants
//     (transparent background, dimension validation, aspect ratio preservation)
//     and keep the SVG dependency contained.
//   - Same contract: Following the existing Render(Config) *Sprite + Sign(Config)
//     uint64 pattern means SVG-based widgets slot into render cache memoization
//     and the panel compositing loop with zero changes to calling code.
//   - Frame-based animation: The Animator struct implements widgets.Animated,
//     driving SVG frame sequences with per-frame durations and optional looping.
//   - Resolution guard: A 16×16 minimum prevents wasted SVG rasterization on
//     tiny or monochrome panels where bitmap widgets are more appropriate.
//   - Defensive parsing: SVG parse failures and panics from the underlying
//     libraries are recovered gracefully, returning nil sprites rather than
//     crashing the compositor.
package svgrender
