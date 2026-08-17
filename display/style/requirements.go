package style

type SurfaceRequirements struct {
	MinWidth        int // Minimum pixel width for correct rendering (0 = unconstrained)
	MinHeight       int // Minimum pixel height for correct rendering (0 = unconstrained)
	PreferredWidth  int // Width for full-featured rendering (0 = unconstrained)
	PreferredHeight int // Height for full-featured rendering (0 = unconstrained)

	// Capability declares the minimum hardware capability level this style requires.
	// The fitness evaluator returns Unsupported when the panel's capability is below this level.
	Capability Capability

	// MinRows is the minimum number of text rows required for meaningful rendering.
	// 0 means unconstrained; negative values are treated as 0.
	MinRows int

	// MinCharsPerLine is the minimum number of characters per line required for
	// meaningful rendering. 0 means unconstrained; negative values are treated as 0.
	MinCharsPerLine int

	// Minimum PPI required (0 = unconstrained)
	MinPPI float64
	// Maximum PPI allowed (0 = unconstrained)
	MaxPPI float64
}
