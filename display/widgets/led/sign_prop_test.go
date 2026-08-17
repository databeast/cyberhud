package led

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/widgets/gradient"
	"pgregory.net/rapid"
)

// --- Feature: led-redesign, Property 31: Sign determinism and animation sensitivity ---

// TestPropertySignDeterminismAndAnimationSensitivity verifies that:
// 1. For any two Config values identical in all fields including animElapsed, sign produces identical uint64 values.
// 2. For any two Config values differing only in animElapsed, sign produces different uint64 values.
//

func TestPropertySignDeterminismAndAnimationSensitivity(t *testing.T) {
	t.Run("determinism", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			s1 := sign(cfg)
			s2 := sign(cfg)
			if s1 != s2 {
				t.Fatalf("sign not deterministic: got %d and %d for same config", s1, s2)
			}
		})
	})

	t.Run("animation_sensitivity", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			// Ensure animElapsed values differ
			elapsed1 := time.Duration(rapid.Int64Range(0, 1_000_000_000).Draw(t, "elapsed1"))
			elapsed2 := time.Duration(rapid.Int64Range(1_000_000_001, 2_000_000_000).Draw(t, "elapsed2"))

			cfg.animElapsed = elapsed1
			s1 := sign(cfg)

			cfg.animElapsed = elapsed2
			s2 := sign(cfg)

			if s1 == s2 {
				t.Fatalf("sign should differ when animElapsed differs: elapsed1=%v elapsed2=%v both produced %d", elapsed1, elapsed2, s1)
			}
		})
	})
}

// --- Feature: led-redesign, Property 32: Sign sensitivity to all config fields ---

// TestPropertySignSensitivityToAllConfigFields verifies that for any two Config values
// differing in exactly one field, sign produces different uint64 values.
//

func TestPropertySignSensitivityToAllConfigFields(t *testing.T) {
	t.Run("Shape", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			// Pick a different shape
			newShape := Shape(rapid.IntRange(0, 3).Draw(t, "newShape"))
			if newShape == cfg.Shape {
				newShape = (cfg.Shape + 1) % 4
			}
			cfg.Shape = newShape
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Shape changes: original shape=%d, new shape=%d", cfg.Shape, newShape)
			}
		})
	})

	t.Run("State", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newState := State(rapid.IntRange(0, 2).Draw(t, "newState"))
			if newState == cfg.State {
				newState = (cfg.State + 1) % 3
			}
			cfg.State = newState
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when State changes")
			}
		})
	})

	t.Run("Brightness", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			// Pick a brightness that differs
			newBrightness := rapid.Float64Range(0.0, 1.0).Draw(t, "newBrightness")
			if newBrightness == cfg.Brightness {
				newBrightness = cfg.Brightness + 0.01
				if newBrightness > 1.0 {
					newBrightness = cfg.Brightness - 0.01
				}
			}
			cfg.Brightness = newBrightness
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Brightness changes")
			}
		})
	})

	t.Run("Diameter", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newDiameter := rapid.IntRange(3, 64).Draw(t, "newDiameter")
			if newDiameter == cfg.Diameter {
				newDiameter = cfg.Diameter + 1
			}
			cfg.Diameter = newDiameter
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Diameter changes")
			}
		})
	})

	t.Run("Foreground", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			cfg.Foreground = drawDifferentColor(t, cfg.Foreground, "fg")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Foreground changes")
			}
		})
	})

	t.Run("Background", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			cfg.Background = drawDifferentColor(t, cfg.Background, "bg")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Background changes")
			}
		})
	})

	t.Run("WarningColor", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			cfg.WarningColor = drawDifferentColor(t, cfg.WarningColor, "warn")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when WarningColor changes")
			}
		})
	})

	t.Run("GlowEnabled", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			cfg.GlowEnabled = !cfg.GlowEnabled
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when GlowEnabled changes")
			}
		})
	})

	t.Run("GlowRadius", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newRadius := rapid.IntRange(0, 32).Draw(t, "newRadius")
			if newRadius == cfg.GlowRadius {
				newRadius = (cfg.GlowRadius + 1) % 33
			}
			cfg.GlowRadius = newRadius
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when GlowRadius changes")
			}
		})
	})

	t.Run("BorderWidth", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newWidth := rapid.IntRange(0, 4).Draw(t, "newWidth")
			if newWidth == cfg.BorderWidth {
				newWidth = (cfg.BorderWidth + 1) % 5
			}
			cfg.BorderWidth = newWidth
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when BorderWidth changes")
			}
		})
	})

	t.Run("BorderColor", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			cfg.BorderColor = drawDifferentColor(t, cfg.BorderColor, "border")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when BorderColor changes")
			}
		})
	})

	t.Run("ShineStyle", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newStyle := ShineStyle(rapid.IntRange(0, 2).Draw(t, "newShine"))
			if newStyle == cfg.ShineStyle {
				newStyle = (cfg.ShineStyle + 1) % 3
			}
			cfg.ShineStyle = newStyle
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when ShineStyle changes")
			}
		})
	})

	t.Run("ShineOpacity", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newOpacity := uint8(rapid.IntRange(0, 255).Draw(t, "newOpacity"))
			if newOpacity == cfg.ShineOpacity {
				newOpacity = cfg.ShineOpacity + 1
			}
			cfg.ShineOpacity = newOpacity
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when ShineOpacity changes")
			}
		})
	})

	t.Run("AnimationType", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newType := Animation(rapid.IntRange(0, 3).Draw(t, "newAnimType"))
			if newType == cfg.Animation.Type {
				newType = (cfg.Animation.Type + 1) % 4
			}
			cfg.Animation.Type = newType
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Animation.Type changes")
			}
		})
	})

	t.Run("AnimationPeriod", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newPeriod := time.Duration(rapid.Int64Range(100_000_000, 5_000_000_000).Draw(t, "newPeriod"))
			if newPeriod == cfg.Animation.Period {
				newPeriod = cfg.Animation.Period + time.Millisecond
			}
			cfg.Animation.Period = newPeriod
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Animation.Period changes")
			}
		})
	})

	t.Run("AnimationMinBrightness", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newMin := rapid.Float64Range(0.0, 0.99).Draw(t, "newMin")
			if newMin == cfg.Animation.MinBrightness {
				newMin = cfg.Animation.MinBrightness + 0.01
				if newMin > 0.99 {
					newMin = cfg.Animation.MinBrightness - 0.01
				}
			}
			cfg.Animation.MinBrightness = newMin
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Animation.MinBrightness changes")
			}
		})
	})

	t.Run("Orientation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			if cfg.Orientation == Horizontal {
				cfg.Orientation = Vertical
			} else {
				cfg.Orientation = Horizontal
			}
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Orientation changes")
			}
		})
	})

	t.Run("Spacing", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			newSpacing := rapid.IntRange(0, 8).Draw(t, "newSpacing")
			if newSpacing == cfg.Spacing {
				newSpacing = (cfg.Spacing + 1) % 9
			}
			cfg.Spacing = newSpacing
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Spacing changes")
			}
		})
	})

	t.Run("GroupEntryState", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			// Change state of first group entry
			idx := 0
			newState := State(rapid.IntRange(0, 2).Draw(t, "entryState"))
			if newState == cfg.Group[idx].State {
				newState = (cfg.Group[idx].State + 1) % 3
			}
			cfg.Group[idx].State = newState
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry State changes")
			}
		})
	})

	t.Run("GroupEntryForeground", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			// Change foreground of first group entry
			cfg.Group[0].Foreground = drawDifferentColor(t, cfg.Group[0].Foreground, "entryFg")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry Foreground changes")
			}
		})
	})

	t.Run("GroupEntryWarningColor", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			cfg.Group[0].WarningColor = drawDifferentColor(t, cfg.Group[0].WarningColor, "entryWarn")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry WarningColor changes")
			}
		})
	})

	t.Run("GroupEntryShape", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			newShape := Shape(rapid.IntRange(0, 3).Draw(t, "entryShape"))
			if newShape == cfg.Group[0].Shape {
				newShape = (cfg.Group[0].Shape + 1) % 4
			}
			cfg.Group[0].Shape = newShape
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry Shape changes")
			}
		})
	})

	t.Run("GroupEntryBorderWidth", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			newWidth := rapid.IntRange(0, 4).Draw(t, "entryBW")
			if newWidth == cfg.Group[0].BorderWidth {
				newWidth = (cfg.Group[0].BorderWidth + 1) % 5
			}
			cfg.Group[0].BorderWidth = newWidth
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry BorderWidth changes")
			}
		})
	})

	t.Run("GroupEntryBorderColor", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			cfg.Group[0].BorderColor = drawDifferentColor(t, cfg.Group[0].BorderColor, "entryBC")
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry BorderColor changes")
			}
		})
	})

	t.Run("GroupEntryGlowEnabled", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			cfg.Group[0].GlowEnabled = !cfg.Group[0].GlowEnabled
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry GlowEnabled changes")
			}
		})
	})

	t.Run("GroupEntryGlowRadius", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfigWithGroup(t)
			original := sign(cfg)

			newRadius := rapid.IntRange(0, 32).Draw(t, "entryGR")
			if newRadius == cfg.Group[0].GlowRadius {
				newRadius = (cfg.Group[0].GlowRadius + 1) % 33
			}
			cfg.Group[0].GlowRadius = newRadius
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when a Group entry GlowRadius changes")
			}
		})
	})

	t.Run("Gradient", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			// Ensure gradient is set
			cfg.Gradient = &GradientConfig{
				Stops: []gradient.ColorStop{
					{Position: 0.0, Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
					{Position: 1.0, Color: color.RGBA{R: 0, G: 0, B: 255, A: 255}},
				},
			}
			original := sign(cfg)

			// Modify the first stop color
			cfg.Gradient.Stops[0].Color.R = cfg.Gradient.Stops[0].Color.R + 1
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Gradient stop color changes")
			}
		})
	})

	t.Run("Bounds", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			cfg := drawConfig(t)
			original := sign(cfg)

			// Modify bounds
			cfg.Bounds.Min.X = cfg.Bounds.Min.X + 1
			modified := sign(cfg)

			if original == modified {
				t.Fatalf("sign should differ when Bounds changes")
			}
		})
	})
}

// --- Generators ---

// drawConfig generates a random valid Config for property testing of the sign function.
func drawConfig(t *rapid.T) Config {
	shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
	state := State(rapid.IntRange(0, 2).Draw(t, "state"))
	brightness := rapid.Float64Range(-1.0, 1.0).Draw(t, "brightness")
	diameter := rapid.IntRange(3, 64).Draw(t, "diameter")

	minX := rapid.IntRange(0, 200).Draw(t, "minX")
	minY := rapid.IntRange(0, 200).Draw(t, "minY")
	bounds := image.Rect(minX, minY, minX+diameter, minY+diameter)

	fg := drawColor(t, "fg")
	bg := drawColor(t, "bg")
	warnColor := drawColor(t, "warn")
	borderColor := drawColor(t, "borderColor")

	glowEnabled := rapid.Bool().Draw(t, "glowEnabled")
	glowRadius := rapid.IntRange(0, 32).Draw(t, "glowRadius")
	borderWidth := rapid.IntRange(0, 4).Draw(t, "borderWidth")
	shineStyle := ShineStyle(rapid.IntRange(0, 2).Draw(t, "shineStyle"))
	shineOpacity := uint8(rapid.IntRange(0, 255).Draw(t, "shineOpacity"))

	animType := Animation(rapid.IntRange(0, 3).Draw(t, "animType"))
	animPeriod := time.Duration(rapid.Int64Range(100_000_000, 5_000_000_000).Draw(t, "animPeriod"))
	animMinBright := rapid.Float64Range(0.0, 0.99).Draw(t, "animMinBright")

	orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))
	spacing := rapid.IntRange(0, 8).Draw(t, "spacing")
	elapsed := time.Duration(rapid.Int64Range(0, 10_000_000_000).Draw(t, "elapsed"))

	cfg := Config{
		Shape:        shape,
		State:        state,
		Brightness:   brightness,
		Diameter:     diameter,
		Bounds:       bounds,
		Foreground:   fg,
		Background:   bg,
		WarningColor: warnColor,
		GlowEnabled:  glowEnabled,
		GlowRadius:   glowRadius,
		BorderWidth:  borderWidth,
		BorderColor:  borderColor,
		ShineStyle:   shineStyle,
		ShineOpacity: shineOpacity,
		Animation: AnimationConfig{
			Type:          animType,
			Period:        animPeriod,
			MinBrightness: animMinBright,
		},
		Orientation: orientation,
		Spacing:     spacing,
		animElapsed: elapsed,
	}

	// Optionally add a gradient
	if rapid.Bool().Draw(t, "hasGradient") {
		numStops := rapid.IntRange(2, 5).Draw(t, "numStops")
		stops := make([]gradient.ColorStop, numStops)
		for i := range stops {
			stops[i] = gradient.ColorStop{
				Position: rapid.Float64Range(0.0, 1.0).Draw(t, "stopPos"),
				Color:    drawColor(t, "stopColor"),
			}
		}
		cfg.Gradient = &GradientConfig{Stops: stops}
	}

	return cfg
}

// drawConfigWithGroup generates a Config with a non-empty Group for testing group entry sensitivity.
func drawConfigWithGroup(t *rapid.T) Config {
	cfg := drawConfig(t)

	numEntries := rapid.IntRange(2, 5).Draw(t, "numEntries")
	entries := make([]GroupEntry, numEntries)
	for i := range entries {
		entries[i] = GroupEntry{
			State:        State(rapid.IntRange(0, 2).Draw(t, "entryState")),
			Foreground:   drawColor(t, "entryFg"),
			WarningColor: drawColor(t, "entryWarn"),
			Shape:        Shape(rapid.IntRange(0, 3).Draw(t, "entryShape")),
			BorderWidth:  rapid.IntRange(0, 4).Draw(t, "entryBW"),
			BorderColor:  drawColor(t, "entryBC"),
			GlowEnabled:  rapid.Bool().Draw(t, "entryGlow"),
			GlowRadius:   rapid.IntRange(0, 32).Draw(t, "entryGR"),
		}
	}
	cfg.Group = entries
	return cfg
}

// drawColor generates a random RGBA color.
func drawColor(t *rapid.T, prefix string) color.RGBA {
	return color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(t, prefix+"R")),
		G: uint8(rapid.IntRange(0, 255).Draw(t, prefix+"G")),
		B: uint8(rapid.IntRange(0, 255).Draw(t, prefix+"B")),
		A: uint8(rapid.IntRange(0, 255).Draw(t, prefix+"A")),
	}
}

// drawDifferentColor generates an RGBA color guaranteed to differ from the given color.
func drawDifferentColor(t *rapid.T, c color.RGBA, prefix string) color.RGBA {
	newColor := drawColor(t, prefix+"New")
	if newColor == c {
		// Flip the R channel to ensure difference
		newColor.R = c.R + 1
	}
	return newColor
}
