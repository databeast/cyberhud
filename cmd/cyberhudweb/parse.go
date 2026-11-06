package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// This file centralizes ALL coupling to the daemon's verbose "OK" line protocol.
// If the console wording changes, this is the single place to update. The web
// server and templates work only with the typed structs below.

// Region is the web UI's read model for one display region, assembled from
// "display status"/"display modes" responses.
type Region struct {
	Index      int
	Name       string    // region surface name, e.g. "main"
	ID         string    // addressing token, e.g. "main.0"
	Controller string
	Current    string    // current mode ID
	Modes      []ModeRef // available modes, in advertised order
}

// CurrentMode returns the ModeRef for the region's active mode, or nil if the
// current mode is not in the advertised list.
func (r *Region) CurrentMode() *ModeRef {
	for i := range r.Modes {
		if strings.EqualFold(r.Modes[i].ID, r.Current) {
			return &r.Modes[i]
		}
	}
	return nil
}

// IsActive reports whether the given mode ID is the region's current mode.
func (r *Region) IsActive(modeID string) bool {
	return strings.EqualFold(r.Current, modeID)
}

// Categories returns the distinct mode categories present in this region, in
// first-appearance order, for building the carousel's category-jump controls.
func (r *Region) Categories() []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range r.Modes {
		c := m.Category
		if c == "" {
			c = "General"
		}
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out
}

// EditField is one row in the config editor: an option's metadata joined with
// its current value, ready for the template to render a typed input.
type EditField struct {
	Option
	Value string // current value (falls back to Default when unknown)
}

// ModeRef is one available mode for a region, with human-facing metadata when
// the daemon provides it (from "display modes").
type ModeRef struct {
	ID       string
	Title    string
	Category string   // UI grouping, e.g. "Attract", "System", "Data"
	Summary  string
	Options  []Option // configurable policy options, when the mode advertises any
	// Policy holds this mode's current option values, keyed by option key, as
	// display strings. For the active mode these come from the region's live
	// policy; for inactive modes, from the mode's global policy snapshot.
	Policy map[string]string
}

// Cat returns the mode's category, defaulting to "General" when unset. Used as
// a stable grouping/anchor key in the carousel.
func (m *ModeRef) Cat() string {
	if m.Category == "" {
		return "General"
	}
	return m.Category
}

// EditFields returns this mode's editable options joined with their current
// values (falling back to each option's default when a value is unknown).
// Returns nil when the mode advertises no options.
func (m *ModeRef) EditFields() []EditField {
	if len(m.Options) == 0 {
		return nil
	}
	fields := make([]EditField, 0, len(m.Options))
	for _, opt := range m.Options {
		val, ok := m.Policy[opt.Key]
		if !ok {
			val = opt.Default
		}
		fields = append(fields, EditField{Option: opt, Value: val})
	}
	return fields
}

// Option describes one configurable policy option of a mode, from the
// "display modes" per-option metadata line. Type is one of string/int/float/
// bool. Allowed, when non-empty, enumerates the valid values (an enum), which
// the editor renders as a select; otherwise the editor renders a typed input.
type Option struct {
	Key     string
	Type    string
	Default string
	Summary string
	Allowed []string
}

// isEnum reports whether the option should render as a select.
func (o Option) isEnum() bool { return len(o.Allowed) > 0 }

// RegionGeom is one region's placement rectangle in virtual-display coordinates,
// from "display geometry".
type RegionGeom struct {
	Index int
	Name  string
	X     int
	Y     int
	W     int
	H     int
}

// Geometry is the virtual display size plus each region's placement.
type Geometry struct {
	VirtualW int
	VirtualH int
	Regions  []RegionGeom
}

// parseField extracts tokens of the form key=value from a protocol line. Values
// may be Go-quoted strings (e.g. name="left-aux") which are unquoted, or bare
// tokens (e.g. mode=clock). It returns a map of key to value. Comma-separated
// list values (e.g. modes=a,b,c) are returned verbatim for the caller to split.
func parseFields(line string) map[string]string {
	fields := make(map[string]string)
	for _, tok := range splitTokens(line) {
		eq := strings.IndexByte(tok, '=')
		if eq < 0 {
			continue
		}
		key := tok[:eq]
		val := tok[eq+1:]
		if len(val) > 0 && val[0] == '"' {
			if unq, err := strconv.Unquote(val); err == nil {
				val = unq
			}
		}
		fields[key] = val
	}
	return fields
}

// splitTokens splits a line into whitespace-separated tokens, but keeps
// double-quoted spans (including spaces) together so that key="a b c" is one
// token. Escaped quotes inside a quoted span are respected.
func splitTokens(line string) []string {
	var tokens []string
	var sb strings.Builder
	inQuote := false
	escaped := false
	flush := func() {
		if sb.Len() > 0 {
			tokens = append(tokens, sb.String())
			sb.Reset()
		}
	}
	for _, r := range line {
		switch {
		case escaped:
			sb.WriteRune(r)
			escaped = false
		case r == '\\':
			sb.WriteRune(r)
			escaped = true
		case r == '"':
			sb.WriteRune(r)
			inQuote = !inQuote
		case (r == ' ' || r == '\t') && !inQuote:
			flush()
		default:
			sb.WriteRune(r)
		}
	}
	flush()
	return tokens
}

// parseStatus parses a "display status" response into regions. Each region line
// looks like:
//
//	region=main.0 name="main" controller=st7789 mode=clock modes=clock,dashboard
//
// The leading bare "OK" line and any non-region lines are ignored.
func parseStatus(lines []string) []Region {
	var regions []Region
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if !strings.HasPrefix(line, "region=") {
			continue
		}
		f := parseFields(line)
		r := Region{
			ID:         f["region"],
			Name:       f["name"],
			Controller: f["controller"],
			Current:    f["mode"],
		}
		r.Index = indexFromID(r.ID)
		if r.Name == "" {
			r.Name = surfaceFromID(r.ID)
		}
		if modes := f["modes"]; modes != "" {
			for _, m := range strings.Split(modes, ",") {
				m = strings.TrimSpace(m)
				if m != "" {
					r.Modes = append(r.Modes, ModeRef{ID: m})
				}
			}
		}
		regions = append(regions, r)
	}
	return regions
}

// parseModes parses a "display modes" response, enriching each region's mode
// list with Title/Summary metadata. The response is region-grouped:
//
//	region=main.0 name="main" controller=st7789 current=clock
//	  mode=clock title="Clock" scope=any summary="A clock."
//	  mode=dashboard title="Dashboard" summary="System dashboard."
//
// It returns a map keyed by region ID (e.g. "main.0") to the ordered mode refs.
// Each mode line may be followed by indented "option=" lines carrying the
// option metadata (key/type/default/summary/allowed) used to build the editor.
func parseModes(lines []string) map[string][]ModeRef {
	out := make(map[string][]ModeRef)
	var currentID string
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "region="):
			f := parseFields(line)
			currentID = f["region"]
			if _, ok := out[currentID]; !ok {
				out[currentID] = nil
			}
		case strings.HasPrefix(line, "mode="):
			if currentID == "" {
				continue
			}
			f := parseFields(line)
			out[currentID] = append(out[currentID], ModeRef{
				ID:       f["mode"],
				Title:    f["title"],
				Category: f["category"],
				Summary:  f["summary"],
			})
		case strings.HasPrefix(line, "option="):
			// Attach to the most recently added mode of the current region.
			if currentID == "" {
				continue
			}
			modes := out[currentID]
			if len(modes) == 0 {
				continue
			}
			f := parseFields(line)
			opt := Option{
				Key:     f["option"],
				Type:    f["type"],
				Default: f["default"],
				Summary: f["summary"],
			}
			if a := strings.TrimSpace(f["allowed"]); a != "" {
				for _, v := range strings.Split(a, ",") {
					if v = strings.TrimSpace(v); v != "" {
						opt.Allowed = append(opt.Allowed, v)
					}
				}
			}
			modes[len(modes)-1].Options = append(modes[len(modes)-1].Options, opt)
			out[currentID] = modes
		}
	}
	return out
}

// parseGeometry parses a "display geometry" response:
//
//	OK
//	  virtual w=320 h=240
//	  region=main.0 index=0 name="main" x=0 y=0 w=240 h=240
func parseGeometry(lines []string) Geometry {
	var geo Geometry
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "virtual "):
			f := parseFields(line)
			geo.VirtualW = atoiDefault(f["w"], 0)
			geo.VirtualH = atoiDefault(f["h"], 0)
		case strings.HasPrefix(line, "region="):
			f := parseFields(line)
			g := RegionGeom{
				Index: atoiDefault(f["index"], indexFromID(f["region"])),
				Name:  f["name"],
				X:     atoiDefault(f["x"], 0),
				Y:     atoiDefault(f["y"], 0),
				W:     atoiDefault(f["w"], 0),
				H:     atoiDefault(f["h"], 0),
			}
			if g.Name == "" {
				g.Name = surfaceFromID(f["region"])
			}
			geo.Regions = append(geo.Regions, g)
		}
	}
	return geo
}

// mergeModes enriches regions (from parseStatus) with the richer mode metadata
// from parseModes, matching by region ID. Modes present in status but missing
// from the modes map keep their bare ID. Order follows the modes map when
// available, otherwise the status order.
func mergeModes(regions []Region, modesByID map[string][]ModeRef) []Region {
	for i := range regions {
		rich, ok := modesByID[regions[i].ID]
		if !ok || len(rich) == 0 {
			continue
		}
		regions[i].Modes = rich
	}
	return regions
}

// indexFromID extracts the trailing index from a region ID like "main.0" -> 0.
// Returns 0 if not parseable.
func indexFromID(id string) int {
	dot := strings.LastIndexByte(id, '.')
	if dot < 0 {
		return 0
	}
	return atoiDefault(id[dot+1:], 0)
}

// surfaceFromID extracts the surface name from a region ID like "main.0" -> "main".
func surfaceFromID(id string) string {
	dot := strings.LastIndexByte(id, '.')
	if dot < 0 {
		return id
	}
	return id[:dot]
}

func atoiDefault(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// parsePolicy parses a "display policy <region>" response into a flat map of
// option key to its current value as a display string. The response is a single
// line in one of two forms:
//
//	OK mode=<id> {"fit":"scale","style":""}
//	OK {"fit":"scale"}
//
// Values decode from JSON (strings, numbers, bools) and are normalized to the
// string forms the editor and the daemon's key=value apply path expect: bools as
// "true"/"false", integral numbers without a decimal point. Returns an empty map
// if no JSON object is present.
func parsePolicy(lines []string) map[string]string {
	out := make(map[string]string)
	if len(lines) == 0 {
		return out
	}
	line := strings.TrimSpace(lines[0])
	// Locate the embedded JSON object.
	brace := strings.IndexByte(line, '{')
	if brace < 0 {
		return out
	}
	raw := line[brace:]

	var parsed map[string]any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return out
	}
	for k, v := range parsed {
		out[k] = jsonValueToString(v)
	}
	return out
}

// jsonValueToString renders a decoded JSON value as the string form used by the
// editor inputs and the daemon's key=value wire format.
func jsonValueToString(v any) string {
	switch n := v.(type) {
	case nil:
		return ""
	case string:
		return n
	case bool:
		if n {
			return "true"
		}
		return "false"
	case float64:
		// JSON numbers decode as float64. Render integral values without a
		// trailing ".0" so they match int option expectations.
		if n == math.Trunc(n) && !math.IsInf(n, 0) {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", n)
	}
}
