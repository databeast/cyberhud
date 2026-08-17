package borderframe

import (
	"encoding/binary"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/gradient"
	"github.com/databeast/cyberhud/display/widgets/icons"
)

const tileSize = 8

// Corner identifies a frame corner for animation origin purposes.
type Corner int

const (
	TopLeft     Corner = iota // Top-left corner (default).
	TopRight                  // Top-right corner.
	BottomLeft                // Bottom-left corner.
	BottomRight               // Bottom-right corner.
)

// Config holds all parameters needed to render a border frame.
type Config struct {
	// Existing fields (backward compatible).
	Bounds  image.Rectangle // Pixel bounds for the border frame.
	TileSet string          // Identifies which icon set to use (e.g., "border").

	// Theme system.
	Theme string // Theme name for Theme_Registry lookup (empty = "sharp").

	// Color tinting.
	ColorTint  color.RGBA // Foreground color for tile alpha compositing.
	Background color.RGBA // Background fill for transparent tile pixels.

	// Gradient sweep.
	GradientStops []gradient.ColorStop // 2-64 stops for perimeter gradient.

	// Glow effect.
	GlowRadius int        // 0-32 pixels outward bloom (0 = disabled).
	GlowColor  color.RGBA // RGB + starting alpha for glow.

	// Pulse animation.
	PulseCycle time.Duration // Cycle period for glow intensity oscillation (0 = static).

	// Scan line animation.
	ScanSpeed  int // Pixels per second (0 = disabled, 1-1000).
	ScanLength int // Highlight length in pixels (default 16, min 1, max perimeter).

	// Corner accents.
	CornerAccent  bool          // Enable corner accent tiles.
	CornerFlash   bool          // Enable corner alpha flash animation.
	FlashDuration time.Duration // Peak hold duration (default 150ms, range 50-1000ms).
	FlashInterval time.Duration // Full cycle interval (default 2000ms, range 200-10000ms).

	// Inner border.
	InnerBorder  bool       // Enable secondary inset border.
	InnerOffset  int        // Inset pixels (default 4, clamped 1-8).
	InnerColor   color.RGBA // Inner border foreground (default opaque white).
	InnerTileSet string     // Inner border tiles (empty = same as primary).

	// Segment reveal animation.
	SegmentReveal bool   // Enable boot-sequence tile reveal.
	RevealSpeed   int    // Tiles per second (1-240, default 30; 0 = disabled).
	RevealOrigin  Corner // Starting corner (default TopLeft).

	// Ticker notches.
	NotchInterval int // Pixels between notch marks (0 = disabled).
	NotchLength   int // Notch height in pixels (default 2, max 8).

	// Opacity.
	// nil = unset (full opacity); non-nil 0.0 = fully transparent;
	// non-nil 1.0 = full opacity; clamped to [0.0, 1.0].
	Opacity *float64

	// ShowBorder flag.
	// nil = border shown (default); explicit false suppresses border rendering.
	ShowBorder *bool
}

// Render produces a single composited border frame image as a Sprite.
// Returns nil when bounds are smaller than 16×16 pixels.
func Render(cfg Config) *widgets.Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	if w < 16 || h < 16 {
		return nil
	}

	cols := w / tileSize
	rows := h / tileSize

	// Create the composited RGBA image covering the full bounds.
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Resolve tile prefix via theme system.
	// Precedence: Config.TileSet (if non-empty) > theme.TileSetPrefix > "border"
	prefix := resolvePrefix(cfg)

	// Draw corner tiles.
	drawTile(img, prefix+"/corner-tl", 0, 0)
	drawTile(img, prefix+"/corner-tr", (cols-1)*tileSize, 0)
	drawTile(img, prefix+"/corner-bl", 0, (rows-1)*tileSize)
	drawTile(img, prefix+"/corner-br", (cols-1)*tileSize, (rows-1)*tileSize)

	// Draw horizontal edge tiles (top and bottom rows, excluding corners).
	for col := 1; col < cols-1; col++ {
		x := col * tileSize
		drawTile(img, prefix+"/h", x, 0)
		drawTile(img, prefix+"/h", x, (rows-1)*tileSize)
	}

	// Draw vertical edge tiles (left and right columns, excluding corners).
	for row := 1; row < rows-1; row++ {
		y := row * tileSize
		drawTile(img, prefix+"/v", 0, y)
		drawTile(img, prefix+"/v", (cols-1)*tileSize, y)
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    "borderframe",
	}
}

// resolvePrefix determines the tile set prefix for icon lookups.
//
// Precedence:
//  1. Config.TileSet if non-empty (backward compatibility override)
//  2. Theme's TileSetPrefix (resolved via LookupTheme)
//  3. Falls back to "border" (sharp theme default via LookupTheme(""))
func resolvePrefix(cfg Config) string {
	if cfg.TileSet != "" {
		return cfg.TileSet
	}
	theme := LookupTheme(cfg.Theme)
	return theme.TileSetPrefix
}

// drawTile composites a named icon tile onto dst at position (x, y).
func drawTile(dst *image.RGBA, name string, x, y int) {
	tile, ok := icons.Get(name)
	if !ok {
		return
	}
	draw.Draw(dst, image.Rect(x, y, x+tileSize, y+tileSize), tile, image.Point{}, draw.Over)
}

// Sign produces a deterministic uint64 hash of a Config for render cache memoization.
// When all new fields are at Go zero-value, this produces the same hash as the
// pre-enhancement Sign (only Bounds + TileSet are hashed). Non-zero fields are
// written to the hasher to maintain backward compatibility.
func Sign(cfg Config) uint64 {
	h := fnv.New64a()
	// Original fields — always hashed for backward compatibility.
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.Y))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.Y))
	h.Write([]byte(cfg.TileSet))

	// New fields — only hashed when non-zero to preserve backward-compatible hash.
	if cfg.Theme != "" {
		h.Write([]byte(cfg.Theme))
	}
	if cfg.ColorTint != (color.RGBA{}) {
		h.Write([]byte{cfg.ColorTint.R, cfg.ColorTint.G, cfg.ColorTint.B, cfg.ColorTint.A})
	}
	if cfg.Background != (color.RGBA{}) {
		h.Write([]byte{cfg.Background.R, cfg.Background.G, cfg.Background.B, cfg.Background.A})
	}
	if len(cfg.GradientStops) > 0 {
		for _, stop := range cfg.GradientStops {
			var buf [8]byte
			binary.LittleEndian.PutUint64(buf[:], math.Float64bits(stop.Position))
			h.Write(buf[:])
			h.Write([]byte{stop.Color.R, stop.Color.G, stop.Color.B, stop.Color.A})
		}
	}
	if cfg.GlowRadius != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.GlowRadius))
	}
	if cfg.GlowColor != (color.RGBA{}) {
		h.Write([]byte{cfg.GlowColor.R, cfg.GlowColor.G, cfg.GlowColor.B, cfg.GlowColor.A})
	}
	if cfg.PulseCycle != 0 {
		binary.Write(h, binary.LittleEndian, int64(cfg.PulseCycle))
	}
	if cfg.ScanSpeed != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.ScanSpeed))
	}
	if cfg.ScanLength != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.ScanLength))
	}
	if cfg.CornerAccent {
		h.Write([]byte{1})
	}
	if cfg.CornerFlash {
		h.Write([]byte{1})
	}
	if cfg.FlashDuration != 0 {
		binary.Write(h, binary.LittleEndian, int64(cfg.FlashDuration))
	}
	if cfg.FlashInterval != 0 {
		binary.Write(h, binary.LittleEndian, int64(cfg.FlashInterval))
	}
	if cfg.InnerBorder {
		h.Write([]byte{1})
	}
	if cfg.InnerOffset != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.InnerOffset))
	}
	if cfg.InnerColor != (color.RGBA{}) {
		h.Write([]byte{cfg.InnerColor.R, cfg.InnerColor.G, cfg.InnerColor.B, cfg.InnerColor.A})
	}
	if cfg.InnerTileSet != "" {
		h.Write([]byte(cfg.InnerTileSet))
	}
	if cfg.SegmentReveal {
		h.Write([]byte{1})
	}
	if cfg.RevealSpeed != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.RevealSpeed))
	}
	if cfg.RevealOrigin != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.RevealOrigin))
	}
	if cfg.NotchInterval != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.NotchInterval))
	}
	if cfg.NotchLength != 0 {
		binary.Write(h, binary.LittleEndian, int32(cfg.NotchLength))
	}
	if cfg.Opacity != nil {
		h.Write([]byte{1})
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(*cfg.Opacity))
		h.Write(buf[:])
	}
	if cfg.ShowBorder != nil {
		h.Write([]byte{1})
		if *cfg.ShowBorder {
			h.Write([]byte{1})
		} else {
			h.Write([]byte{0})
		}
	}

	return h.Sum64()
}

// signWithState extends the public Sign hash with quantized animation state.
// When all animation state is zero (no animations active), this returns the same
// value as Sign(cfg).
func (w *borderframeWidget) signWithState(cfg Config) uint64 {
	base := Sign(cfg)

	// Only incorporate animation state when animations are active.
	if !hasAnimation(cfg) {
		return base
	}

	// Quantize animation state values.
	pulseQ := int64(w.anim.pulsePhase * 100)   // 1/100th precision
	scanQ := int64(w.anim.scanPosition)        // pixel precision
	flashQ := int64(w.cornerFlashAlpha() * 10) // 1/10th precision
	revealQ := int64(w.anim.revealedTiles)     // tile count

	// If all quantized state is zero, return base hash unchanged.
	if pulseQ == 0 && scanQ == 0 && flashQ == 10 && revealQ == 0 {
		// flashQ == 10 means alpha 1.0 (no modulation), which is "zero state" for flash.
		// However, if hasAnimation is true, there IS animation configured, so we should
		// still hash the state to differentiate frames.
	}

	// XOR with a hash of the animation state to extend the base.
	ah := fnv.New64a()
	binary.Write(ah, binary.LittleEndian, pulseQ)
	binary.Write(ah, binary.LittleEndian, scanQ)
	binary.Write(ah, binary.LittleEndian, flashQ)
	binary.Write(ah, binary.LittleEndian, revealQ)
	stateHash := ah.Sum64()

	// XOR the state hash into the base.
	return base ^ stateHash
}

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*borderframeWidget)(nil)
	_ widgets.Described    = (*borderframeWidget)(nil)
	_ widgets.Configurable = (*borderframeWidget)(nil)
)

// borderframeWidget is the internal struct implementing Renderable, Described,
// Configurable, and Animated for the borderframe widget.
type borderframeWidget struct {
	cfg      Config
	opts     widgets.OptionSet
	cache    widgets.RenderCache[Config, widgets.Sprite]
	anim     animationState
	einkMode bool // Set externally or via suppression context inspection.
}

// New creates a borderframe widget instance satisfying widgets.Renderable.
// The returned value also implements Described and Configurable.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	set := widgets.ApplyOptions(opts...)

	w := &borderframeWidget{
		cfg:  cfg,
		opts: set,
	}

	if set.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, w.signWithState)
	}

	return w
}

// RenderFrame produces a single composited Sprite for the border frame,
// or nil if bounds are smaller than 16×16.
func (w *borderframeWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the full rendering pipeline and returns a *widgets.Sprite.
func (w *borderframeWidget) renderSprite(cfg Config) *widgets.Sprite {
	width := cfg.Bounds.Dx()
	height := cfg.Bounds.Dy()
	if width < 16 || height < 16 {
		return nil
	}

	// E-ink displays use a simplified static border (no animation).
	if w.isEinkMode() {
		return w.renderEinkSprite(cfg)
	}

	cols := width / tileSize
	rows := height / tileSize
	prefix := resolvePrefix(cfg)

	// --- Step 1: Base tile compositing (corners + edges) ---
	borderImg := image.NewRGBA(image.Rect(0, 0, width, height))
	// Draw corners.
	drawTile(borderImg, prefix+"/corner-tl", 0, 0)
	drawTile(borderImg, prefix+"/corner-tr", (cols-1)*tileSize, 0)
	drawTile(borderImg, prefix+"/corner-bl", 0, (rows-1)*tileSize)
	drawTile(borderImg, prefix+"/corner-br", (cols-1)*tileSize, (rows-1)*tileSize)
	// Draw horizontal edges.
	for col := 1; col < cols-1; col++ {
		x := col * tileSize
		drawTile(borderImg, prefix+"/h", x, 0)
		drawTile(borderImg, prefix+"/h", x, (rows-1)*tileSize)
	}
	// Draw vertical edges.
	for row := 1; row < rows-1; row++ {
		y := row * tileSize
		drawTile(borderImg, prefix+"/v", 0, y)
		drawTile(borderImg, prefix+"/v", (cols-1)*tileSize, y)
	}

	// --- Step 2: Segment reveal masking ---
	if cfg.SegmentReveal && cfg.RevealSpeed > 0 && !w.anim.revealDone {
		applySegmentRevealMask(borderImg, cols, rows, w.revealedTileCount(), cfg.RevealOrigin)
	}

	// --- Step 3: Color tinting or gradient sweep ---
	gradientColors := computePerimeterGradient(cols, rows, cfg.GradientStops)
	if gradientColors != nil {
		applyGradientTint(borderImg, cols, rows, gradientColors)
	} else {
		applyTint(borderImg, cfg.ColorTint, cfg.Background)
	}

	// --- Step 4: Corner accent overlay ---
	renderCornerAccents(borderImg, prefix, cfg)

	// --- Step 5: Corner flash alpha modulation ---
	if cfg.CornerFlash && cfg.CornerAccent {
		alpha := w.cornerFlashAlpha()
		applyCornerFlashAlpha(borderImg, cols, rows, alpha)
	}

	// --- Step 6: Scan line highlight ---
	if cfg.ScanSpeed > 0 {
		renderScanLine(borderImg, cols, rows, cfg, w.anim.scanPosition)
	}

	// --- Step 7: Glow rendering (separate layer, composited BENEATH border) ---
	glowRadius := effectiveGlowRadius(cfg)
	glowLayer := renderGlow(borderImg, cfg.Bounds, glowRadius, cfg.GlowColor, w.pulseIntensity())

	// --- Step 8: Inner border rendering ---
	innerLayer := renderInnerBorder(cfg)

	// --- Step 9: Ticker notch decoration ---
	renderNotches(borderImg, cfg)

	// --- Step 10: Composite layers in z-order: glow → inner border → primary border ---
	finalImg := image.NewRGBA(image.Rect(0, 0, width, height))
	if glowLayer != nil {
		draw.Draw(finalImg, finalImg.Bounds(), glowLayer, image.Point{}, draw.Over)
	}
	if innerLayer != nil {
		draw.Draw(finalImg, finalImg.Bounds(), innerLayer, image.Point{}, draw.Over)
	}
	draw.Draw(finalImg, finalImg.Bounds(), borderImg, image.Point{}, draw.Over)

	// --- Step 11: Opacity multiplication (final step) ---
	applyOpacity(finalImg, cfg.Opacity)

	sprite := &widgets.Sprite{
		Image:    finalImg,
		Position: cfg.Bounds.Min,
		Label:    "borderframe",
	}

	if w.opts.LabelOverride != "" {
		sprite.Label = w.opts.LabelOverride
	}

	return sprite
}

// isEinkMode returns whether the widget is in e-ink rendering mode.
func (w *borderframeWidget) isEinkMode() bool {
	return w.einkMode
}

// renderEinkSprite renders the border in e-ink mode: monochrome alpha tiles only,
// skipping glow, animations, color tinting, and gradient sweep.
func (w *borderframeWidget) renderEinkSprite(cfg Config) *widgets.Sprite {
	result := Render(cfg)
	if result == nil {
		return nil
	}
	if w.opts.LabelOverride != "" {
		result.Label = w.opts.LabelOverride
	}
	return result
}

// Describe returns the widget's metadata for registry and suppression.
// Capabilities are computed dynamically based on current Config.
func (w *borderframeWidget) Describe() widgets.Descriptor {
	caps := w.capabilities()
	return widgets.Descriptor{
		Name:         "borderframe",
		MinWidth:     16,
		MinHeight:    16,
		Capabilities: caps,
	}
}

// capabilities computes the capability tags based on current Config state.
// "eink-safe" is returned only when no glow or animation features are active.
// "glow" is returned when GlowRadius > 0.
// "animated" is returned when any animation feature is active.
func (w *borderframeWidget) capabilities() []string {
	hasGlow := w.cfg.GlowRadius > 0
	hasAnim := w.cfg.PulseCycle > 0 || w.cfg.ScanSpeed > 0 || w.cfg.CornerFlash || w.cfg.SegmentReveal

	var caps []string
	if !hasGlow && !hasAnim {
		caps = append(caps, "eink-safe")
	}
	if hasGlow {
		caps = append(caps, "glow")
	}
	if hasAnim {
		caps = append(caps, "animated")
	}
	return caps
}

// Configure updates the widget's parameters between frames.
// Accepts a Config value; panics on type mismatch (programming error).
// Resets animation state when animation-affecting fields change.
func (w *borderframeWidget) Configure(cfg interface{}) {
	newCfg := cfg.(Config)

	// Detect if animation-affecting fields changed — if so, reset animation state.
	if animationFieldsChanged(w.cfg, newCfg) {
		w.anim = animationState{}
	}

	w.cfg = newCfg
}

// animationFieldsChanged returns true when any animation-affecting Config fields
// differ between old and new, indicating that animation state should be reset.
func animationFieldsChanged(old, new Config) bool {
	return old.PulseCycle != new.PulseCycle ||
		old.ScanSpeed != new.ScanSpeed ||
		old.ScanLength != new.ScanLength ||
		old.CornerFlash != new.CornerFlash ||
		old.FlashDuration != new.FlashDuration ||
		old.FlashInterval != new.FlashInterval ||
		old.SegmentReveal != new.SegmentReveal ||
		old.RevealSpeed != new.RevealSpeed ||
		old.RevealOrigin != new.RevealOrigin
}

// init registers the borderframe widget type in the global widget registry.
func init() {
	widgets.Register("borderframe", func() widgets.Described {
		return &borderframeWidget{
			cfg: Config{
				Bounds: image.Rect(0, 0, 32, 32),
			},
		}
	})
}
