package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type fileConfig struct {
	Socket   string                     `json:"socket,omitempty"`
	I2C      string                     `json:"i2c,omitempty"`
	Scan     string                     `json:"scan,omitempty"`
	Display  fileDisplayConfig          `json:"display"`
	Policies map[string]json.RawMessage `json:"policies,omitempty"` // modeID → policy JSON
}

type fileDisplayConfig struct {
	Disabled     *bool  `json:"disabled,omitempty"`
	Profile      string `json:"panel,omitempty"`
	DisableInput *bool  `json:"disable_input,omitempty"`

	// PPI overrides the panel-level PPI for this deployment.
	// Zero or omitted means "no config override."
	PPI *float64 `json:"ppi,omitempty"`

	// ─── SCREEN ORIENTATION ───────────────────────────────────────────────
	// Set this to correct for how your display is physically mounted.
	// Accepts either:
	//   - a plain orientation string applied to every screen, e.g. "flip"
	//     (equivalent to the -orientation CLI flag); or
	//   - an object mapping screen names to per-screen orientations, e.g.
	//     {"main": "normal", "left-aux": "cw", "right-aux": "ccw"}.
	// Valid orientation values: "normal", "flip", "cw", "ccw".
	Orientation *orientationField `json:"orientation,omitempty"`

	Width   *int   `json:"width,omitempty"`
	Height  *int   `json:"height,omitempty"`
	MADCTL  string `json:"madctl,omitempty"`
	Rotate  *bool  `json:"rotate,omitempty"`
	XOffset *int   `json:"x_offset,omitempty"`
	YOffset *int   `json:"y_offset,omitempty"`

	DC   string `json:"dc,omitempty"`
	RST  string `json:"rst,omitempty"`
	BL   string `json:"bl,omitempty"`
	Busy string `json:"busy,omitempty"`

	InputKey1  string `json:"input_key1,omitempty"`
	InputKey2  string `json:"input_key2,omitempty"`
	InputKey3  string `json:"input_key3,omitempty"`
	InputUp    string `json:"input_up,omitempty"`
	InputDown  string `json:"input_down,omitempty"`
	InputLeft  string `json:"input_left,omitempty"`
	InputRight string `json:"input_right,omitempty"`
	InputPress string `json:"input_press,omitempty"`
}

// orientationField represents the display.orientation JSON value, which may
// be provided in either of two forms:
//
//   - a plain string, e.g. "flip" — applied uniformly to every screen,
//     equivalent to the -orientation CLI flag; or
//   - an object mapping screen names to orientations, e.g.
//     {"main": "normal", "left-aux": "cw"} — for per-screen overrides on
//     multi-screen panels (or single-screen panels using the "main" key).
type orientationField struct {
	Single    string            // set when JSON was a plain string
	PerScreen map[string]string // set when JSON was an object
}

func (o *orientationField) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("orientation: %w", err)
		}
		o.Single = s
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf(`orientation: must be a string (e.g. "flip") or an object mapping screen names to orientations (e.g. {"main": "flip"}): %w`, err)
	}
	o.PerScreen = m
	return nil
}

func (o orientationField) MarshalJSON() ([]byte, error) {
	if o.PerScreen != nil {
		return json.Marshal(o.PerScreen)
	}
	return json.Marshal(o.Single)
}

func loadConfig(path string) (*fileConfig, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &fileConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func visitedFlags() map[string]bool {
	seen := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})
	return seen
}

func mergeConfig(
	cfg *fileConfig,
	seen map[string]bool,
	socketPath *string,
	i2cBuses *string,
	scanInterval *time.Duration,
	noDisplay *bool,
	noInput *bool,
	displayProfileName *string,
	displayWidth *int,
	displayHeight *int,
	displayMADCTL *string,
	displayRotate *string,
	displayXOffset *int,
	displayYOffset *int,
	displayDC *string,
	displayRST *string,
	displayBL *string,
	displayBusy *string,
	inputKey1 *string,
	inputKey2 *string,
	inputKey3 *string,
	inputUp *string,
	inputDown *string,
	inputLeft *string,
	inputRight *string,
	inputPress *string,
) error {
	if cfg == nil {
		return nil
	}
	if !seen["socket"] && strings.TrimSpace(cfg.Socket) != "" {
		*socketPath = strings.TrimSpace(cfg.Socket)
	}
	if !seen["i2c"] && strings.TrimSpace(cfg.I2C) != "" {
		*i2cBuses = strings.TrimSpace(cfg.I2C)
	}
	if !seen["scan"] && strings.TrimSpace(cfg.Scan) != "" {
		d, err := time.ParseDuration(strings.TrimSpace(cfg.Scan))
		if err != nil {
			return fmt.Errorf("config scan: %w", err)
		}
		*scanInterval = d
	}

	d := cfg.Display
	if !seen["nodisplay"] && d.Disabled != nil {
		*noDisplay = *d.Disabled
	}
	if !seen["panel"] && strings.TrimSpace(d.Profile) != "" {
		*displayProfileName = strings.TrimSpace(d.Profile)
	}
	if !seen["noinput"] && d.DisableInput != nil {
		*noInput = *d.DisableInput
	}
	if !seen["display-width"] && d.Width != nil {
		*displayWidth = *d.Width
	}
	if !seen["display-height"] && d.Height != nil {
		*displayHeight = *d.Height
	}
	if !seen["display-madctl"] && strings.TrimSpace(d.MADCTL) != "" {
		*displayMADCTL = strings.TrimSpace(d.MADCTL)
	}
	if !seen["display-rotate"] && d.Rotate != nil {
		if *d.Rotate {
			*displayRotate = "true"
		} else {
			*displayRotate = "false"
		}
	}
	if !seen["display-x-offset"] && d.XOffset != nil {
		*displayXOffset = *d.XOffset
	}
	if !seen["display-y-offset"] && d.YOffset != nil {
		*displayYOffset = *d.YOffset
	}
	if !seen["display-dc"] && strings.TrimSpace(d.DC) != "" {
		*displayDC = strings.TrimSpace(d.DC)
	}
	if !seen["display-rst"] && strings.TrimSpace(d.RST) != "" {
		*displayRST = strings.TrimSpace(d.RST)
	}
	if !seen["display-bl"] && strings.TrimSpace(d.BL) != "" {
		*displayBL = strings.TrimSpace(d.BL)
	}
	if !seen["display-busy"] && strings.TrimSpace(d.Busy) != "" {
		*displayBusy = strings.TrimSpace(d.Busy)
	}
	if !seen["input-key1"] && strings.TrimSpace(d.InputKey1) != "" {
		*inputKey1 = strings.TrimSpace(d.InputKey1)
	}
	if !seen["input-key2"] && strings.TrimSpace(d.InputKey2) != "" {
		*inputKey2 = strings.TrimSpace(d.InputKey2)
	}
	if !seen["input-key3"] && strings.TrimSpace(d.InputKey3) != "" {
		*inputKey3 = strings.TrimSpace(d.InputKey3)
	}
	if !seen["input-up"] && strings.TrimSpace(d.InputUp) != "" {
		*inputUp = strings.TrimSpace(d.InputUp)
	}
	if !seen["input-down"] && strings.TrimSpace(d.InputDown) != "" {
		*inputDown = strings.TrimSpace(d.InputDown)
	}
	if !seen["input-left"] && strings.TrimSpace(d.InputLeft) != "" {
		*inputLeft = strings.TrimSpace(d.InputLeft)
	}
	if !seen["input-right"] && strings.TrimSpace(d.InputRight) != "" {
		*inputRight = strings.TrimSpace(d.InputRight)
	}
	if !seen["input-press"] && strings.TrimSpace(d.InputPress) != "" {
		*inputPress = strings.TrimSpace(d.InputPress)
	}
	return nil
}

// orientationFromConfig extracts the uniform (all-screens) orientation string
// from the config file, i.e. when display.orientation was given as a plain
// string rather than a per-screen object. Returns "" if no config, no
// orientation specified, or the orientation was given as a per-screen object.
func orientationFromConfig(cfg *fileConfig) string {
	if cfg == nil || cfg.Display.Orientation == nil {
		return ""
	}
	return cfg.Display.Orientation.Single
}

// screenOrientationsFromConfig extracts per-screen orientation overrides from
// the config file. Returns nil if no config or the orientation was given as a
// plain string rather than a per-screen object.
func screenOrientationsFromConfig(cfg *fileConfig) map[string]string {
	if cfg == nil || cfg.Display.Orientation == nil {
		return nil
	}
	if len(cfg.Display.Orientation.PerScreen) == 0 {
		return nil
	}
	return cfg.Display.Orientation.PerScreen
}
