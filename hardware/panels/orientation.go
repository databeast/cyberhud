package panels

import (
	"strings"

	"github.com/databeast/cyberhud/hardware/driver"
)

// applyOrientation applies a named orientation to a panel definition.
// For single-screen panels with a top-level Orientations map, it applies directly
// to def.Config. For multi-screen (virtual) panels, it applies the orientation to
// every virtual screen that has a matching entry in its own Orientations map.
func applyOrientation(def *Definition, orient driver.Orientation) {
	// Single-screen path: top-level Orientations map.
	if def.Orientations != nil {
		if cfg, ok := def.Orientations[orient]; ok {
			def.Config.MADCTL = cfg.MADCTL
			def.Config.XOffset = cfg.XOffset
			def.Config.YOffset = cfg.YOffset
			if cfg.Width > 0 {
				def.Config.Width = cfg.Width
			}
			if cfg.Height > 0 {
				def.Config.Height = cfg.Height
			}
			if cfg.Rotate180 {
				def.Config.Rotate180 = true
			}
		}
	}

	// Multi-screen path: apply to each virtual screen that supports this orientation.
	for i := range def.Virtual {
		if def.Virtual[i].Orientations == nil {
			continue
		}
		cfg, ok := def.Virtual[i].Orientations[orient]
		if !ok {
			continue
		}
		def.Virtual[i].Config.MADCTL = cfg.MADCTL
		def.Virtual[i].Config.XOffset = cfg.XOffset
		def.Virtual[i].Config.YOffset = cfg.YOffset
		if cfg.Rotation != 0 {
			def.Virtual[i].Rotation = cfg.Rotation
		}
		if cfg.MirrorX {
			def.Virtual[i].MirrorX = true
		}
		if cfg.Rotate180 {
			def.Virtual[i].Config.Rotate180 = true
		}
	}
}

// applyScreenOrientation applies a named orientation to a specific screen in a
// multi-screen panel definition. If the screen has an Orientations map with an
// entry for the given orientation, the MADCTL and offsets are applied.
func applyScreenOrientation(def *Definition, screenName string, orient driver.Orientation) {
	// Multi-screen path: find matching virtual screen.
	for i := range def.Virtual {
		if strings.EqualFold(def.Virtual[i].Name, screenName) {
			if def.Virtual[i].Orientations == nil {
				return
			}
			cfg, ok := def.Virtual[i].Orientations[orient]
			if !ok {
				return
			}
			def.Virtual[i].Config.MADCTL = cfg.MADCTL
			def.Virtual[i].Config.XOffset = cfg.XOffset
			def.Virtual[i].Config.YOffset = cfg.YOffset
			if cfg.Rotation != 0 {
				def.Virtual[i].Rotation = cfg.Rotation
			}
			if cfg.MirrorX {
				def.Virtual[i].MirrorX = true
			}
			return
		}
	}

	// Single-screen fallback: if no virtual screens matched, apply to the
	// top-level definition when the screen name is "main" or matches the panel name.
	if len(def.Virtual) == 0 && (strings.EqualFold(screenName, "main") || strings.EqualFold(screenName, def.Name)) {
		if def.Orientations == nil {
			return
		}
		cfg, ok := def.Orientations[orient]
		if !ok {
			return
		}
		def.Config.MADCTL = cfg.MADCTL
		def.Config.XOffset = cfg.XOffset
		def.Config.YOffset = cfg.YOffset
		if cfg.Width > 0 {
			def.Config.Width = cfg.Width
		}
		if cfg.Height > 0 {
			def.Config.Height = cfg.Height
		}
		if cfg.Rotate180 {
			def.Config.Rotate180 = true
		}
	}
}
