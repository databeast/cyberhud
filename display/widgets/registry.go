package widgets

import (
	"fmt"
	"sort"
)

// registry holds widget type factories keyed by unique name.
// Populated during init() via Register calls; read-only after init completes.
var registry = map[string]func() Described{}

// Register adds a widget type factory to the global registry.
// Must be called from init(). Panics on duplicate name.
func Register(name string, factory func() Described) {
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("widgets: duplicate registration for %q", name))
	}
	registry[name] = factory
}

// Registered returns all registered widget names in sorted order.
func Registered() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Factory returns the widget factory function for the given name and true,
// or nil and false if no widget is registered under that name.
func Factory(name string) (func() Described, bool) {
	f, ok := registry[name]
	return f, ok
}
