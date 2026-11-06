package styles

// Serial style declarations for 104x212 panels (e-paper, portrait).
//
// Each entry declares the instrument panel at this resolution for one panel
// capability. To hand-tweak a specific panel, change its Params in the stamp
// call — or replace the stamp with a full def literal carrying a bespoke
// Params.BuildFn to swap the engine for just this cell.

var (
	MonoSlow104x212Style      = monoSlowStyle(104, 212)
	GrayscaleSlow104x212Style = grayscaleSlowStyle(104, 212)
	ColorSlow104x212Style     = colorSlowStyle(104, 212)
)
