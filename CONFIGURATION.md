# Configuration Guide

This document provides JSON configuration file examples for common CyberHUD setups. Configuration files let you set defaults without editing systemd units or repeating CLI flags.

For board-specific configuration and GPIO pin mappings, see the README in each panel's directory under `hardware/panels/`.

## Using a Config File

Create a JSON config at `/etc/cyberhudd.json` and run:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json
```

Any CLI flag overrides the config file value.

## Common Configurations

### Custom GPIO Pin Overrides

If your panel uses non-standard GPIO pins, override them:

File: `/etc/cyberhudd.json`

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

Supported GPIO overrides:
- `dc` – Data/Command pin
- `rst` – Reset pin
- `bl` – Backlight pin (`"none"` to disable)
- `busy` – E-ink busy pin (e-ink displays only)
- `input_key1`, `input_key2`, `input_key3` – Button pins
- `input_up`, `input_down`, `input_left`, `input_right`, `input_press` – Joystick pins

### Headless Mode (Testing/CI)

Daemon runs without display for testing or integration with external control systems.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "disabled": true
  }
}
```

The daemon will:
- Scan I2C buses for STEMMA QT / QWIIC devices
- Monitor GPIO pins
- Serve the console socket
- Skip all display initialization

Use the console protocol to query GPIO, STEMMA, and system state:

```sh
cyberhudctl status
cyberhudctl gpio status
cyberhudctl stemma status
cyberhudctl gpio pins
```

### Fast I2C Scan

For setups with many STEMMA/QWIIC devices, increase scan frequency:

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1,/dev/i2c-3",
  "scan": "500ms",
  "display": {
    "panel": "waveshare-1.3hat"
  }
}
```

Supported values: `1s`, `500ms`, `100ms`, etc. (Go duration format)

## Complete Schema

Here's the full configuration structure (all fields optional). Field names use the exact JSON keys from the Go struct tags in `cmd/cyberhudd/config.go`.

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

For the full schema reference with Go types and detailed field descriptions, see [Schema Reference](ghpages/docs/configuration/schema.md).

## Systemd Integration

To use a config file with the systemd service, create a drop-in override:

```sh
sudo systemctl edit cyberhudd.service
```

Add or modify:

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/cyberhudd -config /etc/cyberhudd.json
```

Reload and restart:

```sh
sudo systemctl daemon-reload
sudo systemctl restart cyberhudd.service
```

Verify:

```sh
sudo systemctl status cyberhudd.service
sudo journalctl -u cyberhudd.service -f
```

## Troubleshooting Configuration

### Config file not being read

Check the path and verify it's valid JSON:

```sh
# Validate JSON syntax
jq . /etc/cyberhudd.json
```

### CLI flags not overriding config

Ensure you're using explicit flags:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json -panel waveshare-2.2
```

### Partial overrides not working

When using a config file, omitted fields use the panel defaults. You only need to specify fields you want to override:

```json
{
  "display": {
    "panel": "waveshare-1.3hat",
    "dc": "GPIO25"
  }
}
```

Empty strings and omitted pointer fields (`*int`, `*bool`, `*float64`) are treated as "not set" and will not override the panel defaults.

## Display Mode Policy: Ticker

The ticker display mode supports runtime policy configuration via the console command interface. Two fields control font sizing and accent color for panel-appropriate rendering.

### Policy Fields

| Field | Allowed Values | Default | Description |
|-------|---------------|---------|-------------|
| `font_tier` | `auto`, `small`, `normal`, `large`, `fullsize` | `auto` | Controls font size tier for panel-appropriate text rendering. When set to `auto`, the tier is selected based on panel pixel height (large for ≥200px, normal otherwise). |
| `accent` | `cyan`, `green`, `amber`, `red`, `white`, `none` | `cyan` | Controls accent color for color TFT styles (glow layers, border tinting, text foreground). |

### CLI Example

Set ticker policy fields at runtime using `cyberhudctl display ticker policy` with key=value pairs:

```sh
cyberhudctl display ticker policy font_tier=large accent=green
```

Response:

```
OK ticker policy style= font=auto font_tier=large line_mode=truncate direction=vertical auto_scroll_ms=0 accent=green show_border=false show_glow=false
```

Invalid values return an error with allowed options:

```sh
cyberhudctl display ticker policy font_tier=huge
# ERR font_tier: must be one of [auto, small, normal, large, fullsize]
```

```sh
cyberhudctl display ticker policy accent=purple
# ERR accent: must be one of [amber, cyan, green, red, white]
```

## See Also

- [Quick Start Guide](QUICKSTART.md)
- [Main README](README.md) – Full flag reference
- [Display Panels](README.md#supported-displays)
