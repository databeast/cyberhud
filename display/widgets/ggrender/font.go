package ggrender

import (
	"fmt"
	"sync"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// Font represents a loaded TrueType font at a specific point size.
type Font struct {
	face font.Face // the rasterized font face
	id   string    // "path:size" identifier
}

// fontKey is used as the cache map key, combining file path and point size.
type fontKey struct {
	path string
	size float64
}

var (
	fontCache map[fontKey]*Font
	fontMu    sync.RWMutex
)

func init() {
	fontCache = make(map[fontKey]*Font)
}

// ListFonts returns the identifiers of all currently cached fonts.
func ListFonts() []string {
	fontMu.RLock()
	defer fontMu.RUnlock()

	ids := make([]string, 0, len(fontCache))
	for _, f := range fontCache {
		ids = append(ids, f.id)
	}
	return ids
}

// LoadFont loads a TrueType font from disk and caches it.
// Subsequent calls with the same path and pointSize return the cached font.
// Returns an error if the file does not exist or is unreadable.
func LoadFont(path string, pointSize float64) (*Font, error) {
	key := fontKey{path: path, size: pointSize}

	// Fast path: check cache with read lock.
	fontMu.RLock()
	if f, ok := fontCache[key]; ok {
		fontMu.RUnlock()
		return f, nil
	}
	fontMu.RUnlock()

	// Slow path: acquire write lock.
	fontMu.Lock()
	defer fontMu.Unlock()

	// Double-check after acquiring write lock.
	if f, ok := fontCache[key]; ok {
		return f, nil
	}

	face, err := gg.LoadFontFace(path, pointSize)
	if err != nil {
		return nil, fmt.Errorf("ggrender: load font %q: %w", path, err)
	}

	f := &Font{
		face: face,
		id:   fmt.Sprintf("%s:%.2f", path, pointSize),
	}
	fontCache[key] = f
	return f, nil
}
