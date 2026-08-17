package font

import (
	"sort"
	"strings"
	"sync"
)

// Metrics describes monospace font dimensions in pixels.
type Metrics struct {
	GlyphWidth   int
	GlyphHeight  int
	GlyphAdvance int
	RowHeight    int
}

// Face describes a bitmap font that can provide one row bitmask per glyph row.
type Face interface {
	ID() string
	Metrics() Metrics
	GlyphRow(ch rune, row int) uint32
}

var (
	facesMu sync.RWMutex
	faces   = map[string]Face{}
)

// Register publishes a font face for runtime selection.
func Register(face Face) {
	if face == nil {
		return
	}
	id := strings.ToLower(strings.TrimSpace(face.ID()))
	if id == "" {
		return
	}
	facesMu.Lock()
	defer facesMu.Unlock()
	faces[id] = face
}

// Get returns a registered font by id.
func Get(id string) (Face, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	facesMu.RLock()
	defer facesMu.RUnlock()
	face, ok := faces[id]
	return face, ok
}

// List returns all registered font faces sorted by ascending glyph height.
func List() []Face {
	facesMu.RLock()
	defer facesMu.RUnlock()
	result := make([]Face, 0, len(faces))
	for _, f := range faces {
		result = append(result, f)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Metrics().GlyphHeight < result[j].Metrics().GlyphHeight
	})
	return result
}

// Default returns the built-in fallback face (Spleen 5×8).
// The generated spleen-5x8 font is always registered via init().
func Default() Face {
	face, _ := Get("spleen-5x8")
	if face == nil {
		// Should never happen — spleen-5x8 is registered at init time.
		// Return the first registered font as emergency fallback.
		all := List()
		if len(all) > 0 {
			return all[0]
		}
		return nil
	}
	return face
}

// familyPriority returns the tie-breaking priority for a font face based on its ID.
// Spleen has highest priority (3), then Terminus (2), then Cozette (1), others (0).
func familyPriority(id string) int {
	lower := strings.ToLower(id)
	switch {
	case strings.HasPrefix(lower, "spleen-"):
		return 3
	case strings.HasPrefix(lower, "terminus-"):
		return 2
	case strings.HasPrefix(lower, "cozette-"):
		return 1
	default:
		return 0
	}
}

// ByHeight returns the best-fit font for the requested pixel height.
// It selects the largest registered variant whose GlyphHeight does not exceed px.
// Priority: Spleen > Terminus > Cozette when heights tie.
// Falls back to smallest variant if px < smallest, largest if px > largest.
func ByHeight(px int) Face {
	all := List()
	if len(all) == 0 {
		return Default()
	}

	// Sort by GlyphHeight descending, then by family priority descending,
	// then by font ID ascending for determinism.
	sort.Slice(all, func(i, j int) bool {
		hi := all[i].Metrics().GlyphHeight
		hj := all[j].Metrics().GlyphHeight
		if hi != hj {
			return hi > hj
		}
		pi := familyPriority(all[i].ID())
		pj := familyPriority(all[j].ID())
		if pi != pj {
			return pi > pj
		}
		return all[i].ID() < all[j].ID()
	})

	// Find the first face where GlyphHeight <= px.
	for _, f := range all {
		if f.Metrics().GlyphHeight <= px {
			return f
		}
	}

	// px is smaller than all registered variants; return the smallest.
	return all[len(all)-1]
}
