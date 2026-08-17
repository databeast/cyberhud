package led

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/databeast/cyberhud/display/widgets"
)

// renderGroup renders a group of LEDs as a single composite Sprite.
// Each entry is rendered independently using renderSingle, with group-level
// settings inherited for any zero-value override fields in the entry.
//
// Special cases:
//   - Zero entries → nil (already caught by validate, defensive)
//   - Single entry → merge with group config, render via renderSingle
//   - >32 entries → truncated to first 32 (already handled by validate)
//   - Diameter < 3 → nil (already caught by validate)
//
// Group glow handling: each LED is rendered within a Diameter×Diameter cell.
// Glow that extends beyond the cell boundary is clipped to the cell area.
func renderGroup(cfg Config) *widgets.Sprite {
	// Defensive: nil or zero entries → nil.
	if cfg.Group == nil || len(cfg.Group) == 0 {
		return nil
	}

	// Defensive: Diameter < 3 → nil.
	if cfg.Diameter < 3 {
		return nil
	}

	// Truncate to 32 entries (defensive; validate already does this).
	entries := cfg.Group
	if len(entries) > 32 {
		entries = entries[:32]
	}

	count := len(entries)

	// Single-entry group: merge entry overrides with group config and render
	// via the standard single-LED path.
	if count == 1 {
		entryCfg := buildEntryCfg(cfg, entries[0])
		resolveColors(&entryCfg)
		effectiveBrightness := resolveBrightness(entryCfg)
		effectiveBrightness = resolveAnimation(entryCfg, effectiveBrightness)
		return renderSingle(entryCfg, effectiveBrightness)
	}

	// Multi-entry group: compute total output dimensions.
	spacing := cfg.Spacing
	diameter := cfg.Diameter

	var totalWidth, totalHeight int
	if cfg.Orientation == Vertical {
		totalWidth = diameter
		totalHeight = count*diameter + (count-1)*spacing
	} else {
		// Horizontal (default)
		totalWidth = count*diameter + (count-1)*spacing
		totalHeight = diameter
	}

	// Allocate output image (initialized to transparent).
	output := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Render each entry independently.
	for i, entry := range entries {
		entryCfg := buildEntryCfg(cfg, entry)
		resolveColors(&entryCfg)
		effectiveBrightness := resolveBrightness(entryCfg)
		effectiveBrightness = resolveAnimation(entryCfg, effectiveBrightness)

		// Render the single LED (may include glow, which expands the image).
		ledSprite := renderSingle(entryCfg, effectiveBrightness)
		if ledSprite == nil || ledSprite.Image == nil {
			continue
		}

		// Compute the cell position in the group output.
		var cellX, cellY int
		if cfg.Orientation == Vertical {
			cellX = 0
			cellY = i * (diameter + spacing)
		} else {
			cellX = i * (diameter + spacing)
			cellY = 0
		}

		// The cell is Diameter × Diameter. The rendered LED may be larger
		// (if glow is enabled). We clip to the cell boundary by only copying
		// the central Diameter×Diameter portion of the rendered LED image.
		ledImg := ledSprite.Image
		ledBounds := ledImg.Bounds()
		ledWidth := ledBounds.Dx()
		ledHeight := ledBounds.Dy()

		// Compute the offset to center the Diameter×Diameter region within
		// the rendered LED image (glow expands equally on all sides).
		offsetX := (ledWidth - diameter) / 2
		offsetY := (ledHeight - diameter) / 2

		// Source rect: the central Diameter×Diameter area of the LED image.
		srcRect := image.Rect(
			ledBounds.Min.X+offsetX,
			ledBounds.Min.Y+offsetY,
			ledBounds.Min.X+offsetX+diameter,
			ledBounds.Min.Y+offsetY+diameter,
		)

		// Destination point in the group output.
		dstPt := image.Pt(cellX, cellY)

		// Draw (copy with alpha compositing) the clipped LED onto the output.
		draw.Draw(output, image.Rect(dstPt.X, dstPt.Y, dstPt.X+diameter, dstPt.Y+diameter),
			ledImg, srcRect.Min, draw.Over)
	}

	return &widgets.Sprite{
		Image:    output,
		Position: cfg.Bounds.Min,
		Label:    "led/group",
	}
}

// buildEntryCfg constructs a per-entry Config by merging a GroupEntry's
// overrides with the group-level Config. Zero-value fields in the entry
// inherit from the group-level config.
//
// Special handling:
//   - State: always use entry.State directly (On=0 is iota, so we can't
//     distinguish "explicit On" from "inherit group default"). Accept this
//     design limitation — entry State is always used directly.
//   - GlowEnabled: always inherit from group-level (bool zero = false makes
//     per-entry override unreliable without a separate flag).
func buildEntryCfg(cfg Config, entry GroupEntry) Config {
	return Config{
		Shape:        mergeShape(entry.Shape, cfg.Shape),
		State:        entry.State,
		Brightness:   cfg.Brightness,
		Diameter:     cfg.Diameter,
		Bounds:       cfg.Bounds,
		Foreground:   mergeColor(entry.Foreground, cfg.Foreground),
		Background:   cfg.Background,
		WarningColor: mergeColor(entry.WarningColor, cfg.WarningColor),
		Gradient:     cfg.Gradient,
		GlowEnabled:  cfg.GlowEnabled,
		GlowRadius:   mergeInt(entry.GlowRadius, cfg.GlowRadius),
		BorderWidth:  mergeInt(entry.BorderWidth, cfg.BorderWidth),
		BorderColor:  mergeColor(entry.BorderColor, cfg.BorderColor),
		ShineStyle:   cfg.ShineStyle,
		ShineOpacity: cfg.ShineOpacity,
		Animation:    cfg.Animation,
		Group:        nil, // Entries are not recursive groups
		Orientation:  cfg.Orientation,
		Spacing:      cfg.Spacing,
		animElapsed:  cfg.animElapsed,
	}
}

// mergeShape returns the entry shape if non-zero (not Circle), otherwise
// falls back to the group-level shape.
// Note: since Circle is iota=0, an entry that explicitly wants Circle
// cannot be distinguished from "inherit". Accept Circle as the fallback.
func mergeShape(entryShape, groupShape Shape) Shape {
	if entryShape != 0 {
		return entryShape
	}
	return groupShape
}

// mergeColor returns the entry color if non-zero-value, otherwise the group color.
func mergeColor(entryColor, groupColor color.RGBA) color.RGBA {
	zero := color.RGBA{}
	if entryColor != zero {
		return entryColor
	}
	return groupColor
}

// mergeInt returns the entry value if non-zero, otherwise the group value.
func mergeInt(entryVal, groupVal int) int {
	if entryVal != 0 {
		return entryVal
	}
	return groupVal
}
