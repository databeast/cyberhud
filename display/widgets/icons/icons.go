package icons

import (
	"image"
	"sort"
	"strings"
	"sync"
)

var (
	mu       sync.RWMutex
	registry = make(map[string]image.Image)
)

// Register stores an icon by name. Empty name (after normalization) or nil image is a no-op.
func Register(name string, img image.Image) {
	name = normalize(name)
	if name == "" || img == nil {
		return
	}
	mu.Lock()
	registry[name] = img
	mu.Unlock()
}

// Get retrieves an icon by name. Returns the image and true if found,
// or nil and false if no icon is registered under the normalized name.
func Get(name string) (image.Image, bool) {
	name = normalize(name)
	mu.RLock()
	img, ok := registry[name]
	mu.RUnlock()
	return img, ok
}

// Reset clears all entries from the icon registry, enabling test isolation.
// After Reset(), Get(name) returns (nil, false) for all previously registered names.
func Reset() {
	mu.Lock()
	registry = make(map[string]image.Image)
	mu.Unlock()
}

// Names returns all registered icon names sorted in lexicographic ascending order.
// Returns a non-nil empty slice if no icons are registered.
func Names() []string {
	mu.RLock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	mu.RUnlock()
	sort.Strings(names)
	return names
}

func normalize(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
