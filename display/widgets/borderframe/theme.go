package borderframe

import (
	"image/color"
	"sync"
	"time"
)

// BorderTheme defines a complete border visual style.
//
// Theme field precedence: Config fields override theme defaults when non-zero.
// Specifically:
//   - Config.ColorTint overrides theme ForegroundColor (when non-zero RGBA)
//   - Config.GlowRadius overrides theme GlowRadius (when > 0)
//   - Config.PulseCycle overrides theme AnimInterval (when > 0)
//
// Zero-value Config fields fall through to theme values.
type BorderTheme struct {
	TileSetPrefix   string        // Icon registry prefix (e.g., "circuit").
	ForegroundColor color.RGBA    // Default foreground RGBA for tile compositing.
	GlowRadius      int           // Default glow radius in pixels (0 = none).
	AnimInterval    time.Duration // Default animation interval (0 = static).
}

var (
	themeMu       sync.RWMutex
	themeRegistry = map[string]BorderTheme{}
)

// RegisterTheme adds a named theme to the registry.
// The name must be between 1 and 64 characters (inclusive); otherwise the
// registration is silently ignored.
func RegisterTheme(name string, theme BorderTheme) {
	if len(name) < 1 || len(name) > 64 {
		return
	}
	themeMu.Lock()
	themeRegistry[name] = theme
	themeMu.Unlock()
}

// LookupTheme returns the BorderTheme registered under name.
// If name is empty or not found in the registry, it returns the "sharp" theme.
func LookupTheme(name string) BorderTheme {
	if name == "" {
		return defaultSharpTheme()
	}
	themeMu.RLock()
	theme, ok := themeRegistry[name]
	themeMu.RUnlock()
	if !ok {
		return defaultSharpTheme()
	}
	return theme
}

// defaultSharpTheme returns the built-in sharp theme definition.
func defaultSharpTheme() BorderTheme {
	return BorderTheme{
		TileSetPrefix:   "border",
		ForegroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	}
}

// init registers the five built-in border themes.
// Tile graphics for each theme are registered separately in the icons package
// (see icons/borders.go and icons/borders_enhanced.go).
func init() {
	// "sharp" — single-pixel-line edges and right-angle corners.
	// Tiles: border/h, border/v, border/corner-{tl,tr,bl,br}
	// Accents: border/accent-{tl,tr,bl,br}
	RegisterTheme("sharp", BorderTheme{
		TileSetPrefix:   "border",
		ForegroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	})

	// "rounded" — quarter-arc corner tiles with the same edge tiles as sharp.
	// Tiles: border/round/h, border/round/v, border/round/corner-{tl,tr,bl,br}
	// Accents: border/round/accent-{tl,tr,bl,br}
	RegisterTheme("rounded", BorderTheme{
		TileSetPrefix:   "border/round",
		ForegroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	})

	// "double-line" — parallel 1px lines spaced 2px apart.
	// Tiles: doubleline/h, doubleline/v, doubleline/corner-{tl,tr,bl,br}
	// Accents: doubleline/accent-{tl,tr,bl,br}
	RegisterTheme("double-line", BorderTheme{
		TileSetPrefix:   "doubleline",
		ForegroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	})

	// "circuit" — PCB-trace-style tiles with 2px-wide traces and via stubs.
	// Tiles: circuit/h, circuit/v, circuit/corner-{tl,tr,bl,br}
	// Accents: circuit/accent-{tl,tr,bl,br}
	RegisterTheme("circuit", BorderTheme{
		TileSetPrefix:   "circuit",
		ForegroundColor: color.RGBA{R: 0, G: 255, B: 128, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	})

	// "hex" — hexagonal-segment tiles with angled edges and 120° corners.
	// Tiles: hex/h, hex/v, hex/corner-{tl,tr,bl,br}
	// Accents: hex/accent-{tl,tr,bl,br}
	RegisterTheme("hex", BorderTheme{
		TileSetPrefix:   "hex",
		ForegroundColor: color.RGBA{R: 0, G: 200, B: 255, A: 255},
		GlowRadius:      0,
		AnimInterval:    0,
	})
}
