package coordinator

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
)

// Region describes a logical display region and its advertised modes.
type Region struct {
	Index      int
	Name       string
	Controller string
	Modes      []string
	Default    string
}

// RegionStatus is a read-only snapshot of one region state.
type RegionStatus struct {
	Index      int
	Name       string
	Controller string
	Current    string
	Modes      []string
}

// RegionDefinition is a static view of one region and the modes it advertises.
type RegionDefinition struct {
	Index      int
	Name       string
	Controller string
	Current    string
	Modes      []catalog.Definition
}

type regionState struct {
	name       string
	controller string
	modes      []string
	current    int
}

// State stores per-region display mode state and supports remote control.
type State struct {
	mu      sync.RWMutex
	regions map[int]regionState
}

// NewState creates a mode state and optionally initialises it with regions.
func NewState(regions ...Region) *State {
	s := &State{}
	s.Reset(regions)
	return s
}

// Reset replaces all known regions and mode state.
// It uses "dashboard" as the fallback mode for regions with no configured modes.
func (s *State) Reset(regions []Region) {
	s.ResetWithFallback(regions, "dashboard")
}

// ResetWithFallback replaces all known regions and mode state.
// fallbackMode is used when a region has zero configured modes.
func (s *State) ResetWithFallback(regions []Region, fallbackMode string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := make(map[int]regionState, len(regions))
	for _, p := range regions {
		if p.Index < 0 {
			continue
		}
		modes := normalizeModes(p.Modes)
		if len(modes) == 0 {
			modes = []string{fallbackMode}
		}
		cur := 0
		if idx := findModeIndex(modes, p.Default); idx >= 0 {
			cur = idx
		}
		next[p.Index] = regionState{
			name:       strings.TrimSpace(p.Name),
			controller: strings.ToLower(strings.TrimSpace(p.Controller)),
			modes:      modes,
			current:    cur,
		}
	}
	s.regions = next
}

// CurrentMode returns the current mode for region index.
func (s *State) CurrentMode(index int) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.regions[index]
	if !ok || len(p.modes) == 0 {
		return ""
	}
	if p.current < 0 || p.current >= len(p.modes) {
		return ""
	}
	return p.modes[p.current]
}

// CurrentModeByName returns the current mode for the region matching name
// (case-insensitive). Returns "" if no region has that name.
func (s *State) CurrentModeByName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.regions {
		if strings.EqualFold(p.name, name) {
			if p.current < 0 || p.current >= len(p.modes) {
				return ""
			}
			return p.modes[p.current]
		}
	}
	return ""
}

// Region returns a snapshot for a single region.
func (s *State) Region(index int) (RegionStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.regions[index]
	if !ok {
		return RegionStatus{}, false
	}
	return regionStatus(index, p), true
}

// HasRegion reports whether a region index is configured.
func (s *State) HasRegion(index int) bool {
	_, ok := s.Region(index)
	return ok
}

// Set changes region mode by name and returns the selected mode.
func (s *State) Set(index int, mode string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.regions[index]
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	i := findModeIndex(p.modes, mode)
	if i < 0 {
		return "", fmt.Errorf("unsupported mode %q for region %d", mode, index)
	}
	p.current = i
	s.regions[index] = p
	return p.modes[p.current], nil
}

// Next advances to the next mode and wraps around.
func (s *State) Next(index int) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.regions[index]
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	if len(p.modes) == 0 {
		return "", fmt.Errorf("region %d has no modes", index)
	}
	p.current = (p.current + 1) % len(p.modes)
	s.regions[index] = p
	return p.modes[p.current], nil
}

// Prev moves to the previous mode and wraps around.
func (s *State) Prev(index int) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.regions[index]
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	if len(p.modes) == 0 {
		return "", fmt.Errorf("region %d has no modes", index)
	}
	p.current--
	if p.current < 0 {
		p.current = len(p.modes) - 1
	}
	s.regions[index] = p
	return p.modes[p.current], nil
}

// Status returns region snapshots in ascending index order.
func (s *State) Status() []RegionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	idxs := make([]int, 0, len(s.regions))
	for idx := range s.regions {
		idxs = append(idxs, idx)
	}
	sort.Ints(idxs)
	out := make([]RegionStatus, 0, len(idxs))
	for _, idx := range idxs {
		out = append(out, regionStatus(idx, s.regions[idx]))
	}
	return out
}

// Definitions returns region metadata enriched with built-in mode descriptions.
func (s *State) Definitions() []RegionDefinition {
	statuses := s.Status()
	out := make([]RegionDefinition, 0, len(statuses))
	for _, st := range statuses {
		out = append(out, RegionDefinition{
			Index:      st.Index,
			Name:       st.Name,
			Controller: st.Controller,
			Current:    st.Current,
			Modes:      catalog.DescribeMany(st.Modes),
		})
	}
	return out
}
