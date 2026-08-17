package displaymodes

import (
	"strings"
	"sync"
)

// registryMu protects concurrent access to the factories map.
var registryMu sync.RWMutex

// factories holds instance-based mode registrations keyed by normalized mode ID.
var factories = map[string]ModeFactory{}

// RegisterFactory is called by mode packages in their init() to self-register
// a factory for the given mode ID. The ID is normalized (lowercased, trimmed).
//
// Panics if id is empty/whitespace-only after normalization or if factory is nil.
// If a factory already exists for the normalized ID it is silently overwritten.
func RegisterFactory(id string, factory ModeFactory) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		panic("displaymodes.RegisterFactory: id must not be empty or whitespace-only")
	}
	if factory == nil {
		panic("displaymodes.RegisterFactory: factory must not be nil")
	}

	registryMu.Lock()
	defer registryMu.Unlock()
	factories[id] = factory
}

// GetInstance returns a freshly constructed ModeInstance for the given mode ID.
// The ID is normalized before lookup. Returns (nil, false) if the mode is unknown.
//
// Each call constructs a new instance — no caching is performed.
func GetInstance(id string) (ModeInstance, bool) {
	id = strings.ToLower(strings.TrimSpace(id))

	registryMu.RLock()
	defer registryMu.RUnlock()

	if factory, ok := factories[id]; ok {
		return factory(), true
	}
	return nil, false
}

// IsKnownInstance reports whether a mode ID has a factory in the instance-based registry.
func IsKnownInstance(id string) bool {
	id = strings.ToLower(strings.TrimSpace(id))

	registryMu.RLock()
	defer registryMu.RUnlock()

	_, ok := factories[id]
	return ok
}
