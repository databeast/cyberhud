package region

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/surface"
)

// RegionManager is the RegionManager architectural component: it owns Region
// lifecycle including allocation, constraint validation, lookup by name or index,
// mode routing, and input focus management. It operates on a VirtualDisplay to
// carve out non-overlapping Region surfaces.
type RegionManager struct {
	vd            *VirtualDisplay
	screens       []ScreenPosition // physical screens for TextHints resolution
	panelProduct  string           // normalized panel product name for TextHints propagation
	panelPPI      float64          // panel-level PPI for cascading into TextHints
	configPPI     float64          // config-level PPI override; zero means "no override"
	regions       []*Region
	byName        map[string]*Region // keyed by lowercase name
	modeValidator func(string) bool  // if nil, all modes are considered valid
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

// NewRegionManager creates a RegionManager backed by the given VirtualDisplay.
// The vd parameter is the framebuffer from which Region surfaces will be carved as
// zero-copy sub-images. The returned manager has no regions allocated and no mode
// validator configured (all modes are accepted by default).
func NewRegionManager(vd *VirtualDisplay) *RegionManager {
	return &RegionManager{
		vd:     vd,
		byName: make(map[string]*Region),
	}
}

// NewRegionManagerWithScreens creates a RegionManager backed by the given
// VirtualDisplay and physical screen positions. The vd parameter is the framebuffer
// for Region surface allocation. The screens parameter provides hardware screen
// geometry so that Regions allocated through this manager resolve TextHints from
// hardware drivers instead of using degraded defaults.
func NewRegionManagerWithScreens(vd *VirtualDisplay, screens []ScreenPosition) *RegionManager {
	return &RegionManager{
		vd:      vd,
		screens: screens,
		byName:  make(map[string]*Region),
	}
}

// SetModeValidator configures the mode validation function used by Allocate and
// SetMode. Passing nil disables validation and treats all modes as valid.
func (rm *RegionManager) SetModeValidator(validator func(string) bool) {
	rm.modeValidator = validator
}

// Allocate validates a RegionSpec and adds a new Region to the RegionManager.
// The spec parameter defines the region's name, bounds, and optional default mode.
// Returns an error if any constraint is violated (name conflict, out-of-bounds,
// overlap with existing regions, or invalid mode). The first allocated Region
// receives input focus automatically.
func (rm *RegionManager) Allocate(spec RegionSpec) error {
	if err := rm.validateSpec(spec); err != nil {
		return err
	}

	// Create the region's surface as a sub-image view of the VD framebuffer.
	surf := surface.NewFromSubImage(rm.vd.FrameBuffer(), spec.Bounds)

	// Use NewRegionWithScreens when screens are available so that TextHints
	// resolve Capability from the actual hardware driver (e.g., ColorFast for
	// an ST7735S) instead of the DefaultTextHints fallback (MonoFast).
	var r *Region
	if len(rm.screens) > 0 {
		r = NewRegionWithScreens(spec.Name, spec.Bounds, surf, rm.screens, rm.panelProduct, rm.panelPPI, rm.configPPI)
	} else {
		r = NewRegion(spec.Name, spec.Bounds, surf)
	}
	r.controller = strings.ToLower(strings.TrimSpace(spec.Controller))
	r.modes = normalizeModeList(spec.Modes)

	// Set default mode if specified.
	if spec.DefaultMode != "" {
		r.mode = spec.DefaultMode
	}

	// First allocated region gets input focus.
	if len(rm.regions) == 0 {
		r.SetInputFocus(true)
	}

	rm.regions = append(rm.regions, r)
	rm.byName[strings.ToLower(spec.Name)] = r

	return nil
}

func normalizeModeList(modes []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(modes))
	for _, mode := range modes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		if mode == "" {
			continue
		}
		if _, ok := seen[mode]; ok {
			continue
		}
		seen[mode] = struct{}{}
		out = append(out, mode)
	}
	return out
}

func modeAllowed(modes []string, mode string) bool {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		return false
	}
	if len(modes) == 0 {
		return true
	}
	for _, allowed := range modes {
		if allowed == mode {
			return true
		}
	}
	return false
}

func nextMode(r *Region) (string, error) {
	modes := r.Modes()
	if len(modes) == 0 {
		return "", fmt.Errorf("region %q has no modes", r.Name())
	}
	current := strings.ToLower(strings.TrimSpace(r.CurrentMode()))
	for i, mode := range modes {
		if mode == current {
			return modes[(i+1)%len(modes)], nil
		}
	}
	return modes[0], nil
}

func prevMode(r *Region) (string, error) {
	modes := r.Modes()
	if len(modes) == 0 {
		return "", fmt.Errorf("region %q has no modes", r.Name())
	}
	current := strings.ToLower(strings.TrimSpace(r.CurrentMode()))
	for i, mode := range modes {
		if mode == current {
			if i == 0 {
				return modes[len(modes)-1], nil
			}
			return modes[i-1], nil
		}
	}
	return modes[len(modes)-1], nil
}

func regionStatus(index int, r *Region) RegionStatus {
	return RegionStatus{
		Index:      index,
		Name:       r.Name(),
		Controller: r.Controller(),
		Current:    r.CurrentMode(),
		Modes:      r.Modes(),
	}
}

func regionDefinition(index int, r *Region) RegionDefinition {
	return RegionDefinition{
		Index:      index,
		Name:       r.Name(),
		Controller: r.Controller(),
		Current:    r.CurrentMode(),
		Modes:      catalog.DescribeMany(r.Modes()),
	}
}

// pendingRegion tracks a region during layout validation before allocation.
type pendingRegion struct {
	name   string
	bounds image.Rectangle
}

// AllocateLayout validates and allocates an entire RegionLayout atomically.
// The layout parameter contains the set of RegionSpec values to allocate together.
// If any spec fails validation, no Regions are allocated and a combined error is
// returned. An empty layout (no specs) is treated as absent — the caller handles
// default generation.
func (rm *RegionManager) AllocateLayout(layout RegionLayout) error {
	// Empty layout is treated as absent — caller handles default generation.
	if len(layout.Specs) == 0 {
		return nil
	}

	// Validate ALL specs first, collecting errors.
	var errs []string
	// Track specs validated so far to check inter-spec overlap and name uniqueness.
	var pending []pendingRegion

	for _, spec := range layout.Specs {
		if err := rm.validateSpecWithPending(spec, pending); err != nil {
			errs = append(errs, err.Error())
		} else {
			pending = append(pending, pendingRegion{name: spec.Name, bounds: spec.Bounds})
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}

	// All passed — allocate in order.
	for _, spec := range layout.Specs {
		if err := rm.Allocate(spec); err != nil {
			// This should not happen since we pre-validated, but handle gracefully.
			return err
		}
	}

	return nil
}

// Region returns the Region at the given zero-based index in allocation order.
// The index parameter identifies the Region positionally. Returns the Region and
// true if found, or nil and false if the index is out of range.
func (rm *RegionManager) Region(index int) (*Region, bool) {
	if index < 0 || index >= len(rm.regions) {
		return nil, false
	}
	return rm.regions[index], true
}

// RegionByName returns the Region matching the given name (case-insensitive lookup).
// The name parameter is the region identifier to search for. Returns the Region and
// true if found, or nil and false if no region has that name.
func (rm *RegionManager) RegionByName(name string) (*Region, bool) {
	r, ok := rm.byName[strings.ToLower(name)]
	return r, ok
}

// Regions returns a copy of all allocated Regions in allocation order. The returned
// slice is a shallow copy so callers cannot mutate the manager's internal list.
func (rm *RegionManager) Regions() []*Region {
	result := make([]*Region, len(rm.regions))
	copy(result, rm.regions)
	return result
}

// SetMode changes the active display mode for a Region identified by target.
// The target parameter is either a zero-based integer index (as a string) or a
// region name (case-insensitive). The modeID parameter is the display mode
// identifier to activate. Returns an error if the target cannot be resolved, the
// mode is not registered (when a modeValidator is configured), or the Region's
// SetMode fails.
func (rm *RegionManager) SetMode(target string, modeID string) error {
	r, err := rm.resolveTarget(target)
	if err != nil {
		return err
	}

	// Validate mode exists.
	if rm.modeValidator != nil && !rm.modeValidator(modeID) {
		return fmt.Errorf("region %q: display mode %q is not registered", r.Name(), modeID)
	}
	if !modeAllowed(r.modes, modeID) {
		return fmt.Errorf("region %q: display mode %q is not allowed", r.Name(), modeID)
	}

	return r.SetMode(modeID)
}

// Set switches a region by index to the named mode and returns the selected mode.
func (rm *RegionManager) Set(index int, modeID string) (string, error) {
	r, ok := rm.Region(index)
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	if err := rm.SetMode(strconv.Itoa(index), modeID); err != nil {
		return "", err
	}
	return r.CurrentMode(), nil
}

// Next advances a region by index to the next allowed mode.
func (rm *RegionManager) Next(index int) (string, error) {
	r, ok := rm.Region(index)
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	next, err := nextMode(r)
	if err != nil {
		return "", err
	}
	if err := r.SetMode(next); err != nil {
		return "", err
	}
	return r.CurrentMode(), nil
}

// Prev moves a region by index to the previous allowed mode.
func (rm *RegionManager) Prev(index int) (string, error) {
	r, ok := rm.Region(index)
	if !ok {
		return "", fmt.Errorf("region %d is not configured", index)
	}
	prev, err := prevMode(r)
	if err != nil {
		return "", err
	}
	if err := r.SetMode(prev); err != nil {
		return "", err
	}
	return r.CurrentMode(), nil
}

// CurrentMode returns the current mode for the region index.
func (rm *RegionManager) CurrentMode(index int) string {
	r, ok := rm.Region(index)
	if !ok {
		return ""
	}
	return r.CurrentMode()
}

// CurrentModeByName returns the current mode for the region matching name.
func (rm *RegionManager) CurrentModeByName(name string) string {
	r, ok := rm.RegionByName(name)
	if !ok {
		return ""
	}
	return r.CurrentMode()
}

// HasRegion reports whether a region index is configured.
func (rm *RegionManager) HasRegion(index int) bool {
	_, ok := rm.Region(index)
	return ok
}

// Status returns region snapshots in allocation order.
func (rm *RegionManager) Status() []RegionStatus {
	out := make([]RegionStatus, 0, len(rm.regions))
	for i, r := range rm.regions {
		out = append(out, regionStatus(i, r))
	}
	return out
}

// Definitions returns region metadata enriched with built-in mode descriptions.
func (rm *RegionManager) Definitions() []RegionDefinition {
	out := make([]RegionDefinition, 0, len(rm.regions))
	for i, r := range rm.regions {
		out = append(out, regionDefinition(i, r))
	}
	return out
}

// InputActiveRegion returns the Region currently receiving input events, or nil
// if no region has input focus. Part of the RegionManager's input focus management.
func (rm *RegionManager) InputActiveRegion() *Region {
	for _, r := range rm.regions {
		if r.HasInputFocus() {
			return r
		}
	}
	return nil
}

// SetInputFocus moves input focus to the Region identified by target, removing
// focus from all other regions. The target parameter is either a zero-based integer
// index (as a string) or a region name (case-insensitive). Returns an error if the
// target cannot be resolved.
func (rm *RegionManager) SetInputFocus(target string) error {
	r, err := rm.resolveTarget(target)
	if err != nil {
		return err
	}

	// Remove focus from all regions.
	for _, existing := range rm.regions {
		existing.SetInputFocus(false)
	}

	// Set focus on target.
	r.SetInputFocus(true)
	return nil
}

// resolveTarget resolves a region target string: tries integer index first, then name.
func (rm *RegionManager) resolveTarget(target string) (*Region, error) {
	// Try to parse as integer index.
	if idx, err := strconv.Atoi(target); err == nil {
		r, ok := rm.Region(idx)
		if !ok {
			return nil, fmt.Errorf("region: no region at index %d", idx)
		}
		return r, nil
	}

	// Look up by name (case-insensitive).
	r, ok := rm.RegionByName(target)
	if !ok {
		return nil, fmt.Errorf("region: no region named %q", target)
	}
	return r, nil
}

// validateSpecCore performs shared validation checks for a RegionSpec: name validation,
// name uniqueness vs existing regions, bounds dimensions, bounds containment,
// overlap vs existing regions, and mode validation.
func (rm *RegionManager) validateSpecCore(spec RegionSpec) error {
	// Validate name.
	if err := rm.validateName(spec.Name); err != nil {
		return err
	}

	// Check name uniqueness against existing regions.
	if _, exists := rm.byName[strings.ToLower(spec.Name)]; exists {
		return fmt.Errorf("region: name %q already exists", spec.Name)
	}

	// Validate bounds dimensions.
	if spec.Bounds.Dx() < 1 || spec.Bounds.Dy() < 1 {
		return fmt.Errorf("region %q: bounds must have width >= 1 and height >= 1", spec.Name)
	}

	// Validate bounds within VD.
	if !spec.Bounds.In(rm.vd.Bounds()) {
		return fmt.Errorf("region %q: bounds %v extend outside virtual display %v", spec.Name, spec.Bounds, rm.vd.Bounds())
	}

	// Validate no overlap with existing regions.
	for _, existing := range rm.regions {
		if !spec.Bounds.Intersect(existing.Bounds()).Empty() {
			return fmt.Errorf("region %q: bounds %v overlap with existing region %q at %v", spec.Name, spec.Bounds, existing.Name(), existing.Bounds())
		}
	}

	// Validate mode if specified and validator is set.
	if spec.DefaultMode != "" && rm.modeValidator != nil && !rm.modeValidator(spec.DefaultMode) {
		return fmt.Errorf("region %q: display mode %q is not registered", spec.Name, spec.DefaultMode)
	}
	if spec.DefaultMode != "" && !modeAllowed(spec.Modes, spec.DefaultMode) {
		return fmt.Errorf("region %q: display mode %q is not allowed", spec.Name, spec.DefaultMode)
	}

	return nil
}

// validateSpec validates a RegionSpec against constraints and existing regions.
func (rm *RegionManager) validateSpec(spec RegionSpec) error {
	return rm.validateSpecCore(spec)
}

// validateSpecWithPending validates a spec against both existing regions AND pending specs
// (used during atomic layout validation).
func (rm *RegionManager) validateSpecWithPending(spec RegionSpec, pending []pendingRegion) error {
	if err := rm.validateSpecCore(spec); err != nil {
		return err
	}

	// Check name uniqueness against pending specs.
	for _, p := range pending {
		if strings.EqualFold(p.name, spec.Name) {
			return fmt.Errorf("region: name %q already exists", spec.Name)
		}
	}

	// Validate no overlap with pending specs.
	for _, p := range pending {
		if !spec.Bounds.Intersect(p.bounds).Empty() {
			return fmt.Errorf("region %q: bounds %v overlap with existing region %q at %v", spec.Name, spec.Bounds, p.name, p.bounds)
		}
	}

	return nil
}

// validateName checks the name constraints (1-64 chars, non-empty).
func (rm *RegionManager) validateName(name string) error {
	if name == "" {
		return fmt.Errorf("region: name must be 1-64 characters, got empty string")
	}
	if len(name) > 64 {
		return fmt.Errorf("region: name %q exceeds 64 character limit", name)
	}
	return nil
}
