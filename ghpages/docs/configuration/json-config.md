# JSON Config Examples

Configuration files let you set defaults without editing systemd units or repeating CLI flags. This page covers common configurations for typical CyberHUD setups.

For the full list of available fields, see the [Schema Reference](schema.md).

## Using a Config File

Create a JSON config at `/etc/cyberhudd.json` and run:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json
```

## CLI Flag vs. Config File Precedence

CLI flags **always override** config file values. The daemon uses a `mergeConfig` function that only applies config file values for flags that were not explicitly set on the command line.

The precedence order is:

1. **CLI flags** (highest priority) — any flag passed on the command line wins
2. **Config file values** — applied only when the corresponding CLI flag is absent
3. **Built-in defaults** — used when neither CLI nor config specifies a value

### How It Works

When the daemon starts, it:

1. Parses all CLI flags and records which flags were explicitly provided
2. Loads the JSON config file (if `-config` is specified)
3. For each config field, checks whether the corresponding CLI flag was set:
   - If the CLI flag **was** set → keeps the CLI value (config is ignored)
   - If the CLI flag **was not** set → applies the config file value

This means you can use a config file for baseline settings and selectively override individual values from the command line without touching the file.

### Precedence Examples

Given this config file:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "display": {
    "panel": "waveshare-2.2"
  }
}
```

| Command | Result |
|---------|--------|
| `sudo ./cyberhudd -config /etc/cyberhudd.json` | Uses `waveshare-2.2` from config |
| `sudo ./cyberhudd -config /etc/cyberhudd.json -panel waveshare-1.3hat` | Uses `waveshare-1.3hat` from CLI (overrides config) |
| `sudo ./cyberhudd -config /etc/cyberhudd.json -socket /tmp/test.sock` | Uses `/tmp/test.sock` from CLI, `waveshare-2.2` from config |

### CLI Flag to Config Field Mapping

| CLI Flag | Config JSON Path | Notes |
|----------|-----------------|-------|
| `-socket` | `socket` | Unix domain socket path |
| `-i2c` | `i2c` | Comma-separated I2C bus paths |
| `-scan` | `scan` | Go duration string (e.g., `"2s"`, `"500ms"`) |
| `-nodisplay` | `display.disabled` | Boolean |
| `-panel` | `display.panel` | Panel profile ID |
| `-noinput` | `display.disable_input` | Boolean |
| `-display-width` | `display.width` | Integer pixels |
| `-display-height` | `display.height` | Integer pixels |
| `-display-madctl` | `display.madctl` | Hex string (e.g., `"0x60"`) |
| `-display-rotate` | `display.rotate` | Boolean |
| `-display-x-offset` | `display.x_offset` | Integer pixels |
| `-display-y-offset` | `display.y_offset` | Integer pixels |
| `-display-dc` | `display.dc` | GPIO pin name |
| `-display-rst` | `display.rst` | GPIO pin name |
| `-display-bl` | `display.bl` | GPIO pin name or `"none"` |
| `-display-busy` | `display.busy` | GPIO pin name (e-ink only) |
| `-input-key1` | `display.input_key1` | GPIO pin name |
| `-input-key2` | `display.input_key2` | GPIO pin name |
| `-input-key3` | `display.input_key3` | GPIO pin name |
| `-input-up` | `display.input_up` | GPIO pin name |
| `-input-down` | `display.input_down` | GPIO pin name |
| `-input-left` | `display.input_left` | GPIO pin name |
| `-input-right` | `display.input_right` | GPIO pin name |
| `-input-press` | `display.input_press` | GPIO pin name |

## Per-Panel Minimal Configurations

Each example below shows the minimal JSON config needed to run CyberHUD with a specific panel. Copy the example, save as `/etc/cyberhudd.json`, and start the daemon with `-config /etc/cyberhudd.json`.

### Waveshare 1.3-inch LCD HAT (Color, 240×240, buttons + joystick)

This is the default panel. If you use this panel with default wiring, you only need:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.3hat"
  }
}
```

### Waveshare 1.44-inch LCD HAT (Color, 128×128, buttons + joystick)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.44"
  }
}
```

### Waveshare 2.2-inch SPI Display (Color, 320×240, no inputs)

A display-only panel with no onboard buttons. Control the daemon entirely via CLI:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-2.2",
    "disable_input": true
  }
}
```

### ST7789 320×240 Generic (Color, no inputs)

For bare ST7789 breakout boards wired to custom GPIO pins:

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

### ST7789 240×240 Generic (Color, no inputs)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "st7789-240x240",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24",
    "disable_input": true
  }
}
```

### ST7789 240×135 Generic (Color, no inputs)

A compact color display commonly used in mini-TFT breakouts:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "st7789-240x135",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24",
    "disable_input": true
  }
}
```

### Waveshare 1.3-inch OLED HAT (Monochrome, 128×64, buttons + joystick)

A monochrome OLED with full input controls. No backlight pin is needed (OLEDs are self-emitting):

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.3-oled-hat"
  }
}
```

### Waveshare 2.23-inch OLED HAT — I2C Variant (Monochrome, 128×32, no inputs)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-2.23-oled-hat-i2c",
    "disable_input": true
  }
}
```

### Waveshare 2.23-inch OLED HAT — SPI Variant (Monochrome, 128×32, no inputs)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-2.23-oled-hat-spi",
    "disable_input": true
  }
}
```

### Waveshare 4.26-inch E-Paper HAT (Monochrome, 800×480, no inputs)

E-paper panels use deferred refresh and require a `busy` pin instead of a backlight pin:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-4.26-epaper",
    "disable_input": true
  }
}
```

### Adafruit 2.13-inch E-Ink Bonnet (Monochrome, 122×250, two buttons)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-2.13-ssd1680"
  }
}
```

### Adafruit 15×7 CharliePlex LED Matrix (Grayscale, 15×7, no inputs)

An I2C-driven LED matrix — no SPI or GPIO display pins needed:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-15x7-charlieplex",
    "disable_input": true
  }
}
```

### Adafruit 16×8 CharliePlex LED Matrix Bonnet (Grayscale, 16×8, no inputs)

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-4120-charlieplex",
    "disable_input": true
  }
}
```

### Waveshare Triple Screen (Color, 240×240 main + 2×160×80 aux, two buttons)

A multi-screen panel with three independent displays. Use the [Regions](../getting-started/regions.md) system to address each screen individually:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-triple-screen"
  }
}
```

## Common Configurations

### Custom GPIO Pin Overrides

If your panel uses non-standard GPIO pins, override them in the config file:

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

| Field | Description |
|-------|-------------|
| `dc` | Data/Command pin |
| `rst` | Reset pin |
| `bl` | Backlight pin (`"none"` to disable) |
| `busy` | E-ink busy pin (e-ink displays only) |
| `input_key1`, `input_key2`, `input_key3` | Button pins |
| `input_up`, `input_down`, `input_left`, `input_right`, `input_press` | Joystick pins |

For detailed Pin Assignments and board-specific assignments, see [Pin Assignments](../reference/pin-assignments.md).

### Headless Mode (Testing/CI)

Run the daemon without a display for testing or integration with external control systems:

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

In headless mode the daemon will:

- Scan I2C buses for STEMMA QT / QWIIC devices
- Monitor GPIO pins
- Serve the console socket
- Skip all display initialization

Use the console protocol to query state:

```sh
cyberhudctl status
cyberhudctl gpio status
cyberhudctl stemma status
cyberhudctl gpio pins
```

### Fast I2C Scan

For setups with many STEMMA QT / QWIIC devices, increase the scan frequency:

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

The `scan` field accepts any Go duration format value. Common choices:

| Value | Use Case |
|-------|----------|
| `2s` | Default, suitable for most setups |
| `1s` | Moderate scan frequency |
| `500ms` | Fast scanning for many devices |
| `100ms` | Very fast, higher CPU usage |

### Multi-Display Orientation

For panels mounted upside-down or rotated, use the `orientation` map to correct the display without rewiring:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-triple-screen",
    "orientation": {
      "main": "flip",
      "left-aux": "normal",
      "right-aux": "normal"
    }
  }
}
```

Valid orientation values: `"normal"`, `"flip"`, `"cw"`, `"ccw"`.

### Persisted Mode Policies

Save mode-specific settings so they survive daemon restarts. Use `cyberhudctl freeze policy` to persist the current mode configuration, or write policies directly:

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.3hat"
  },
  "policies": {
    "clock": {"style": "digital", "show_seconds": "true"},
    "thermal": {"layout": "overview"}
  }
}
```

See [Policy System](index.md) for details on how policies are frozen, persisted, and loaded at startup.
