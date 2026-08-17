package source

// HSLColor represents a color in HSL space (hue 0-360, saturation 0-100, lightness 0-100).
type HSLColor struct {
	H float64
	S float64
	L float64
}

// SquareConfig defines a single rectangle within a cluster.
type SquareConfig struct {
	OffsetX       float64
	OffsetY       float64
	Size          float64
	Aspect        float64
	Rotation      float64
	Color         HSLColor
	PhaseOffset   float64
	CycleDuration float64
	PeakOpacity   float64
}

// ClusterConfig defines a cluster of rectangles.
type ClusterConfig struct {
	CenterXPct     float64
	CenterYPct     float64
	Squares        []SquareConfig
	BoundingRadius float64
	SpawnTime      float64
	FadeInDuration float64
}

// ActiveFragment represents a visible pseudocode text snippet.
type ActiveFragment struct {
	Text            string
	X, Y            float64
	StartTime       float64
	FadeInDuration  float64
	HoldDuration    float64
	FadeOutDuration float64
	FontSize        float64
	Color           HSLColor
	PeakOpacity     float64
}

// PerfState tracks frame performance for adaptive scaling.
type PerfState struct {
	FrameTimes          []float64
	CurrentSquareCount  int
	OriginalSquareCount int
	HasReduced          bool
}

// FragmentState manages fragment scheduling.
type FragmentState struct {
	ActiveFragments []ActiveFragment
	LastSpawnTime   float64
	LastSpawnedText string
}

// GeometricFrame is the mode-specific snapshot type for style dispatch.
type GeometricFrame struct {
	Policy Policy
}
