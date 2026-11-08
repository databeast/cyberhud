package styles

// Schedule style declarations across the common panel sizes for each capability
// class. Schedule renders one adaptive layout; these cells exist only so the
// registry can resolve the best-fitting capability/size for a panel. Each cell
// is a one-liner — the layout itself lives in layout.go.

var (
	// MonoSlow — e-paper mono.
	MonoSlow128x64Style   = monoSlowStyle(128, 64)
	MonoSlow128x128Style  = monoSlowStyle(128, 128)
	MonoSlow200x200Style  = monoSlowStyle(200, 200)
	MonoSlow250x122Style  = monoSlowStyle(250, 122)
	MonoSlow296x128Style  = monoSlowStyle(296, 128)
	MonoSlow400x300Style  = monoSlowStyle(400, 300)
	MonoSlow800x480Style  = monoSlowStyle(800, 480)

	// MonoFast — OLED mono.
	MonoFast128x32Style  = monoFastStyle(128, 32)
	MonoFast128x64Style  = monoFastStyle(128, 64)
	MonoFast128x128Style = monoFastStyle(128, 128)
	MonoFast240x240Style = monoFastStyle(240, 240)
	MonoFast320x240Style = monoFastStyle(320, 240)

	// GrayscaleSlow — grayscale e-paper.
	GrayscaleSlow200x200Style = grayscaleSlowStyle(200, 200)
	GrayscaleSlow296x128Style = grayscaleSlowStyle(296, 128)
	GrayscaleSlow400x300Style = grayscaleSlowStyle(400, 300)
	GrayscaleSlow480x800Style = grayscaleSlowStyle(480, 800)
	GrayscaleSlow800x480Style = grayscaleSlowStyle(800, 480)

	// GrayscaleFast — grayscale LED matrix / fast panels.
	GrayscaleFast128x128Style = grayscaleFastStyle(128, 128)
	GrayscaleFast240x240Style = grayscaleFastStyle(240, 240)
	GrayscaleFast320x240Style = grayscaleFastStyle(320, 240)
	GrayscaleFast480x320Style = grayscaleFastStyle(480, 320)

	// ColorSlow — color e-paper.
	ColorSlow200x200Style = colorSlowStyle(200, 200)
	ColorSlow296x128Style = colorSlowStyle(296, 128)
	ColorSlow400x300Style = colorSlowStyle(400, 300)
	ColorSlow480x800Style = colorSlowStyle(480, 800)
	ColorSlow800x480Style = colorSlowStyle(800, 480)

	// ColorFast — color TFT.
	ColorFast128x128Style = colorFastStyle(128, 128)
	ColorFast160x80Style  = colorFastStyle(160, 80)
	ColorFast240x135Style = colorFastStyle(240, 135)
	ColorFast240x240Style = colorFastStyle(240, 240)
	ColorFast240x320Style = colorFastStyle(240, 320)
	ColorFast320x240Style = colorFastStyle(320, 240)
	ColorFast320x480Style = colorFastStyle(320, 480)
	ColorFast480x320Style = colorFastStyle(480, 320)
	ColorFast480x800Style = colorFastStyle(480, 800)
	ColorFast800x480Style = colorFastStyle(800, 480)
)
