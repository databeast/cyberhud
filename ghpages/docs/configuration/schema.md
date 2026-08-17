# Schema Reference

This page documents the complete JSON configuration schema for the CyberHUD daemon. All fields are derived from the `fileConfig` and `fileDisplayConfig` Go structs in `cmd/cyberhudd/config.go`. Every field is optional — omitted fields use their zero value or panel defaults.

## Complete Schema

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "policies": {},
  "display": {
    "panel": "waveshare-1.3hat",
    "disabled": false,
    "disable_input": false,
    "ppi": 0,
    "orientation": {},
    "rotate": false,
    "width": 0,
    "height": 0,
    "madctl": "",
    "x_offset": 0,
    "y_offset": 0,
    "dc": "",
    "rst": "",
    "bl": "",
    "busy": "",
    "input_key1": "",
    "input_key2": "",
    "input_key3": "",
    "input_up": "",
    "input_down": "",
    "input_left": "",
    "input_right": "",
    "input_press": ""
  }
}
```

## Top-Level Fields (`fileConfig`)

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `socket` | Socket | string | yes | Unix domain socket path for the console protocol. When omitted, the daemon uses the default path (`/run/cyberhudd/console.sock`). |
| `i2c` | I2C | string | yes | Comma-separated list of I2C device paths to scan for STEMMA QT / QWIIC devices (e.g., `"/dev/i2c-1"` or `"/dev/i2c-1,/dev/i2c-3"`). |
| `scan` | Scan | string | yes | I2C scan interval in Go duration format (e.g., `"2s"`, `"500ms"`, `"100ms"`). Controls how often the daemon re-scans I2C buses for connected STEMMA sensors. |
| `policies` | Policies | map[string]json.RawMessage | yes | Per-mode policy snapshots. Keys are mode IDs (e.g., `"ticker"`, `"attract_bokeh"`), values are JSON objects containing mode-specific policy fields. Populated by the `cyberhudctl freeze policy` command and loaded at daemon startup. |
| `display` | Display | fileDisplayConfig | yes | Display hardware and GPIO configuration object. See [Display Object Fields](#display-object-fields-filedisplayconfig) below. |

## Display Object Fields (`fileDisplayConfig`)

These fields appear inside the `"display"` object and configure the panel hardware, orientation, and GPIO pin assignments.

### Panel and Mode Fields

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `panel` | Profile | string | yes | Panel profile name identifying the display hardware (e.g., `"waveshare-1.3hat"`, `"st7789-320x240"`). When omitted, auto-detection or CLI flag is used. Note: the Go field is named `Profile` but the JSON key is `panel`. |
| `disabled` | Disabled | *bool | yes | When `true`, the daemon skips all display initialization and runs in headless mode. Omit or set `false` for normal operation. |
| `disable_input` | DisableInput | *bool | yes | When `true`, disables button and joystick input handling even if the panel defines input pins. Useful for headless or kiosk deployments. |
| `ppi` | PPI | *float64 | yes | Pixels-per-inch override for this deployment. When omitted or zero, the panel-level PPI is used. Affects font sizing and layout calculations. |

### Screen Orientation Fields

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `orientation` | Orientation | map[string]string | yes | Per-screen orientation overrides. Keys are screen names (e.g., `"main"`, `"left-aux"`, `"right-aux"`). Valid values are orientation strings: `"normal"` (no rotation), `"flip"` (180°), `"cw"` (90° clockwise), `"ccw"` (90° counter-clockwise). |
| `rotate` | Rotate | *bool | yes | Legacy rotation flag. When `true`, applies 180° rotation to the display. Prefer `orientation` for new configurations as it supports per-screen control. |

### Geometry Override Fields

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `width` | Width | *int | yes | Display width override in pixels. When omitted, the panel default resolution is used. |
| `height` | Height | *int | yes | Display height override in pixels. When omitted, the panel default resolution is used. |
| `madctl` | MADCTL | string | yes | Memory Access Data Control (MADCTL) register value override as a hex string. Controls display memory read/write direction. Empty string uses the panel default. |
| `x_offset` | XOffset | *int | yes | Horizontal pixel offset applied to the display framebuffer. Used to correct for panels where the visible area is offset within the controller's memory. When omitted, the panel default is used. |
| `y_offset` | YOffset | *int | yes | Vertical pixel offset applied to the display framebuffer. Used to correct for panels where the visible area is offset within the controller's memory. When omitted, the panel default is used. |

### GPIO Pin Fields

All GPIO pin fields accept a GPIO pin identifier string (e.g., `"GPIO25"`). An empty string means the field is unset and the panel default pin assignment is used. Set to `"none"` to explicitly disable a pin.

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `dc` | DC | string | yes | Data/Command GPIO pin for SPI communication. Required for SPI-based displays to distinguish between data and command bytes. |
| `rst` | RST | string | yes | Hardware reset GPIO pin. Used to perform a hardware reset of the display controller on startup. |
| `bl` | BL | string | yes | Backlight control GPIO pin. Set to `"none"` to disable backlight control on panels without PWM backlight. |
| `busy` | Busy | string | yes | Busy/wait GPIO pin used by e-ink displays to signal when a refresh operation is in progress. |

### Input Pin Fields

Input pin fields configure GPIO pins for physical buttons and joystick controls. All accept a GPIO pin identifier string.

| JSON Key | Go Field | Go Type | Optional | Description |
|----------|----------|---------|----------|-------------|
| `input_key1` | InputKey1 | string | yes | Button 1 (Key1) GPIO pin. Typically mapped to the top or primary action button. |
| `input_key2` | InputKey2 | string | yes | Button 2 (Key2) GPIO pin. Typically mapped to the middle or secondary action button. |
| `input_key3` | InputKey3 | string | yes | Button 3 (Key3) GPIO pin. Typically mapped to the bottom or tertiary action button. |
| `input_up` | InputUp | string | yes | Joystick up direction GPIO pin. |
| `input_down` | InputDown | string | yes | Joystick down direction GPIO pin. |
| `input_left` | InputLeft | string | yes | Joystick left direction GPIO pin. |
| `input_right` | InputRight | string | yes | Joystick right direction GPIO pin. |
| `input_press` | InputPress | string | yes | Joystick press/center button GPIO pin. Acts as the primary select action in menu navigation. |

## Default Value Conventions

- **Omitted fields**: Use the panel default or zero value. All fields carry `omitempty` in their JSON struct tags, so omitting a field is the standard way to accept defaults.
- **Empty string (`""`)**: For string fields, treated as unset; the panel default is used.
- **`nil` pointer**: For `*bool`, `*int`, and `*float64` fields, an omitted value results in a nil pointer, meaning "no config override." The panel default applies.
- **`false` for `*bool` fields**: Explicitly disables the feature (e.g., `"disabled": false` keeps the display active).

These conventions allow partial overrides — you only need to specify the fields you want to change from the panel defaults.

## Example: Minimal Config

```json
{
  "display": {
    "panel": "waveshare-1.3hat"
  }
}
```

## Example: Multi-Screen Orientation

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "display": {
    "panel": "waveshare-triple-screen",
    "orientation": {
      "main": "normal",
      "left-aux": "cw",
      "right-aux": "ccw"
    }
  }
}
```

## Example: Custom GPIO Pins

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "st7789-320x240",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24",
    "disable_input": true
  }
}
```

## Example: Headless Mode

```json
{
  "display": {
    "disabled": true
  }
}
```

## Example: With Persisted Policies

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "policies": {
    "ticker": {"font_tier": "large", "accent": "green"},
    "attract_bokeh": {"speed": "1.5", "density": "0.8"}
  },
  "display": {
    "panel": "waveshare-1.3hat"
  }
}
```

## See Also

- [JSON Config Examples](json-config.md) — Common configuration patterns
- [Pin Assignments](../reference/pin-assignments.md) — Board-specific default pin assignments
- [Systemd Integration](systemd.md) — Running with systemd
