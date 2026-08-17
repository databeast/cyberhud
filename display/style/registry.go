package style

import (
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// StyleRegistry is a type-safe ordered collection of Style[S,P] implementations
// for a single display mode. Each mode owns its own registry instance.
// S is the mode-specific snapshot type, providing compile-time type safety.
type StyleRegistry[S any, P catalog.ConfigPolicy] struct {
	styles []Style[S, P]
	byName map[string]Style[S, P]
}

// NewRegistry creates a StyleRegistry[S,P] from the given styles.
// Panics if styles is empty.
// Duplicate names (after normalization) are rejected; the first registration wins.
func NewRegistry[S any, P catalog.ConfigPolicy](styles ...Style[S, P]) *StyleRegistry[S, P] {
	if len(styles) == 0 {
		panic("style: NewRegistry requires at least one Style")
	}

	r := &StyleRegistry[S, P]{
		byName: make(map[string]Style[S, P], len(styles)),
	}

	for _, s := range styles {
		key := normalizeName(s.Name())
		if _, exists := r.byName[key]; exists {
			// Duplicate normalized name: first registration wins, silently reject.
			continue
		}
		r.styles = append(r.styles, s)
		r.byName[key] = s
	}

	return r
}

// Lookup returns the Style with the given name, or nil if not found.
// Matching is case-insensitive and whitespace-trimmed.
// If the name is not found directly, aliases are checked. A deprecation
// warning is logged when a legacy alias resolves successfully.
func (r *StyleRegistry[S, P]) Lookup(name string) Style[S, P] {
	key := normalizeName(name)
	if s, ok := r.byName[key]; ok {
		return s
	}
	return nil
}

// Cycle returns the Style at ((currentIndex + delta) % count + count) % count.
// If current name is not found, resolves via best-fit using the provided hints.
func (r *StyleRegistry[S, P]) Cycle(current string, delta int, hints textlayout.TextHints) Style[S, P] {
	count := len(r.styles)
	if count == 0 {
		return nil
	}

	key := normalizeName(current)
	idx := -1
	for i, s := range r.styles {
		if normalizeName(s.Name()) == key {
			idx = i
			break
		}
	}

	if idx < 0 {
		// Unknown style — resolve via best-fit, then find its index.
		resolved, _ := fitnessFallback[S, P](r, hints)
		for i, s := range r.styles {
			if s.Name() == resolved.Name() {
				idx = i
				break
			}
		}
		if idx < 0 {
			idx = 0
		}
	}

	newIdx := ((idx+delta)%count + count) % count
	return r.styles[newIdx]
}

// Enumerate returns all registered styles in registration order.
// The returned slice is a copy; callers may modify it without affecting the registry.
func (r *StyleRegistry[S, P]) Enumerate() []Style[S, P] {
	out := make([]Style[S, P], len(r.styles))
	copy(out, r.styles)
	return out
}

// Normalize returns the canonical form of a style name (lowercase, trimmed).
// If the name matches a registered style, returns that style's Name().
// If the name matches a legacy alias that resolves to a registered style,
// returns the target style's Name().
// Otherwise returns the default style's Name().
func (r *StyleRegistry[S, P]) Normalize(name string) string {
	key := normalizeName(name)
	if s, ok := r.byName[key]; ok {
		return s.Name()
	}
	return ""
}

// normalizeName applies the canonical normalization: lowercase and trim whitespace.
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// aliasKey is the composite lookup key for alias resolution.
type aliasKey struct {
	product string // normalized panel product name (or "product/screen" compound)
	mode    string // normalized display mode ID
}

// registry is the package-level singleton.
var (
	mu      sync.RWMutex
	aliases = map[aliasKey]string{}
)

// normalize applies lowercase and whitespace trimming to an input string.
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// RegisterAlias registers a single alias mapping from (product, mode) to styleName.
// Empty or whitespace-only product, mode, or styleName silently discards the entry.
// Duplicate keys retain the first registration (first-registration-wins).
func RegisterAlias(product, mode, styleName string) {
	p := normalize(product)
	m := normalize(mode)
	s := normalize(styleName)

	if p == "" || m == "" || s == "" {
		return
	}

	key := aliasKey{product: p, mode: m}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := aliases[key]; !exists {
		aliases[key] = s
	}
}

// RegisterAliases bulk-registers aliases for a single product.
// Each map entry maps a display mode ID to a style name.
// The same validation and normalization rules as RegisterAlias apply to each entry.
func RegisterAliases(product string, modes map[string]string) {
	for mode, styleName := range modes {
		RegisterAlias(product, mode, styleName)
	}
}

// Lookup returns the alias style name for the given product and mode.
// Returns ("", false) if no alias is registered or inputs are invalid.
func Lookup(product, mode string) (string, bool) {
	p := normalize(product)
	m := normalize(mode)

	if p == "" || m == "" {
		return "", false
	}

	key := aliasKey{product: p, mode: m}

	mu.RLock()
	defer mu.RUnlock()

	s, ok := aliases[key]
	return s, ok
}

// LookupCompound resolves aliases for multi-screen panels.
// Tries "product/screen" compound key first, then falls back to product-only.
// Returns ("", false) if neither key has a registered alias.
func LookupCompound(product, screen, mode string) (string, bool) {
	p := normalize(product)
	m := normalize(mode)

	if p == "" || m == "" {
		return "", false
	}

	sc := normalize(screen)

	// If screen is non-empty, try compound key first.
	if sc != "" {
		compound := p + "/" + sc
		key := aliasKey{product: compound, mode: m}

		mu.RLock()
		s, ok := aliases[key]
		mu.RUnlock()

		if ok {
			return s, true
		}
	}

	// Fall back to product-only lookup.
	return Lookup(p, m)
}

// reset clears the alias registry. It is unexported and intended only for
// test isolation — the registry is a package-level singleton and property
// tests mutate global state.
func reset() {
	mu.Lock()
	defer mu.Unlock()
	aliases = map[aliasKey]string{}
}

// AliasEntry is a single (mode, style) pair for enumeration.
type AliasEntry struct {
	Mode  string
	Style string
}

// AliasesForProduct returns all aliases registered for the given product.
// Results are sorted lexicographically by Mode ID (byte-value comparison).
// Returns nil if product is empty/whitespace or has no registrations.
func AliasesForProduct(product string) []AliasEntry {
	p := normalize(product)
	if p == "" {
		return nil
	}

	mu.RLock()
	defer mu.RUnlock()

	var entries []AliasEntry
	for key, styleName := range aliases {
		if key.product == p {
			entries = append(entries, AliasEntry{Mode: key.mode, Style: styleName})
		}
	}

	if len(entries) == 0 {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Mode < entries[j].Mode
	})

	return entries
}

func ResolveStyle[S any, P catalog.ConfigPolicy](
	registry *StyleRegistry[S, P],
	hints textlayout.TextHints,
	modeID string,
	configuredStyle string,
) (Style[S, P], string) {
	// 1. Configured style takes precedence.
	if configuredStyle != "" {
		if s := registry.Lookup(configuredStyle); s != nil {
			return s, "configured"
		}
		// Configured style not found in registry — fall through to alias/fitness.
	}

	// 2. Alias resolution via compound key (product/screen fallback to product-only).
	if hints.PanelProduct != "" {
		aliasName, found := LookupCompound(hints.PanelProduct, hints.ScreenName, modeID)
		if found {
			if s := registry.Lookup(aliasName); s != nil {
				return s, "alias"
			}
			// Alias target not in registry — log warning, fall through to fitness.
			log.Printf("style-alias: unresolved alias %q for mode %q on product %q — falling back to fitness", aliasName, modeID, hints.PanelProduct)
		}
	}

	// 3. Fitness-scored fallback.
	return fitnessFallback[S, P](registry, hints)
}

// fitnessFallback iterates all styles in registration order, selecting the one with
// the highest Supports(hints) score. Ties are broken by registration order (first wins).
func fitnessFallback[S any, P catalog.ConfigPolicy](registry *StyleRegistry[S, P], hints textlayout.TextHints) (Style[S, P], string) {
	styles := registry.Enumerate()

	bestStyle := styles[0]
	bestFitness := bestStyle.Supports(hints)

	for _, s := range styles[1:] {
		f := s.Supports(hints)
		if f > bestFitness {
			bestFitness = f
			bestStyle = s
		}
	}

	return bestStyle, "fitness"
}
