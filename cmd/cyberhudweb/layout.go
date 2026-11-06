package main

// Layout computes to-scale SVG placement for the virtual display visualizer.
// It is a pure transform from virtual-display geometry (device pixels) into SVG
// user-space coordinates fitted into a target box, preserving aspect ratio.

const (
	// layoutTargetW/H are the SVG viewport dimensions the layout is fitted into.
	layoutTargetW = 640
	layoutTargetH = 400
	// layoutPad is the inner padding (SVG units) around the scaled display.
	layoutPad = 8
)

// PlacedRegion is one region positioned within the SVG viewport.
type PlacedRegion struct {
	ID      string // region addressing ID, e.g. "main.0"
	Name    string
	Current string // current mode, for labelling
	X       float64
	Y       float64
	W       float64
	H       float64
}

// LayoutView is the full visualizer model: the SVG viewport size, the scaled
// outline of the virtual display, and each region's placed rectangle.
type LayoutView struct {
	ViewW   float64
	ViewH   float64
	OutlineX float64
	OutlineY float64
	OutlineW float64
	OutlineH float64
	Regions  []PlacedRegion
	// Scaled reports whether real geometry was available. When false, the caller
	// should render a plain (non-proportional) region list instead of the SVG.
	Scaled bool
}

// BuildLayout fits the state's virtual-display geometry into a targetW x targetH
// box (minus padding), preserving aspect ratio and centering. It joins geometry
// rectangles to regions (by ID/index) so each placed rectangle carries the
// current mode for labelling.
//
// If geometry is unavailable (daemon not bound, or no regions with positive
// virtual size), Scaled is false and Regions is empty — the template falls back
// to a text listing.
func BuildLayout(st *State, targetW, targetH int) *LayoutView {
	lv := &LayoutView{ViewW: float64(targetW), ViewH: float64(targetH)}
	if st == nil {
		return lv
	}
	vw, vh := st.Geometry.VirtualW, st.Geometry.VirtualH
	if vw <= 0 || vh <= 0 || len(st.Geometry.Regions) == 0 {
		return lv
	}

	availW := float64(targetW) - 2*layoutPad
	availH := float64(targetH) - 2*layoutPad
	if availW <= 0 || availH <= 0 {
		return lv
	}

	// Uniform scale to fit while preserving aspect ratio.
	scale := availW / float64(vw)
	if sh := availH / float64(vh); sh < scale {
		scale = sh
	}

	scaledW := float64(vw) * scale
	scaledH := float64(vh) * scale
	// Center within the viewport.
	offX := (float64(targetW) - scaledW) / 2
	offY := (float64(targetH) - scaledH) / 2

	lv.OutlineX = offX
	lv.OutlineY = offY
	lv.OutlineW = scaledW
	lv.OutlineH = scaledH
	lv.Scaled = true

	// Build a lookup of current mode by region name/index for labelling.
	currentByName := make(map[string]string)
	idByName := make(map[string]string)
	for _, r := range st.Regions {
		currentByName[r.Name] = r.Current
		idByName[r.Name] = r.ID
	}

	for _, g := range st.Geometry.Regions {
		pr := PlacedRegion{
			Name:    g.Name,
			Current: currentByName[g.Name],
			ID:      idByName[g.Name],
			X:       offX + float64(g.X)*scale,
			Y:       offY + float64(g.Y)*scale,
			W:       float64(g.W) * scale,
			H:       float64(g.H) * scale,
		}
		if pr.ID == "" {
			// Fall back to <name>.0 if status didn't list this region.
			pr.ID = g.Name + ".0"
		}
		lv.Regions = append(lv.Regions, pr)
	}
	return lv
}
