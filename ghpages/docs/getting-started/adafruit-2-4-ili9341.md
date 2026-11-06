# Adafruit 2.4" TFT Touchscreen Breakout

The [Adafruit 2.4" TFT LCD with Touchscreen Breakout](https://www.adafruit.com/product/2478) is a 240×320 color SPI display based on the ILI9341 controller. The board carries a resistive touchscreen and a microSD socket in addition to the display.

CyberHUD drives the display path only — the resistive touchscreen is not yet wired, so this panel is output-only. The same ILI9341 driver also covers the many generic 2.4"/2.8"/3.2" 240×320 SPI TFT modules built around this controller.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-2.4-ili9341
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-2.4-ili9341"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 240 × 320 pixels (portrait) |
| Color format | RGB565 |
| Controller | ILI9341 |
| Interface | SPI |
| Input | None (touchscreen not yet supported) |

## Pin Assignments

This board is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK/CLK, MOSI/SI) and its chip-select (CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| D/C | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |
| LITE | `display.bl` | `-display-bl` | Backlight (optional; tie to 3.3V for always-on, or set `bl` to `"none"`) |

The board's microSD chip-select (CCS) and card-detect are not used by CyberHUD. Example config for a board wired with D/C→GPIO25, RST→GPIO27, backlight→GPIO24:

```json
{
  "display": {
    "panel": "adafruit-2.4-ili9341",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24"
  }
}
```

Those GPIO numbers are an example, not a requirement — use whatever you wired. The LITE/backlight pin can be left tied high in hardware; if you do, omit `bl` or set it to `"none"`.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| `pin not found` at startup | DC or RST override not supplied | This panel has no default pins — set `display.dc` and `display.rst` to the GPIOs you wired |
| Display is white/blank but backlight is on | DC/RST mapped to the wrong GPIO, or wrong panel name | Confirm your `display.dc` / `display.rst` values match your wiring; check panel name spelling |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your module | Override MADCTL in JSON config to adjust color order bit |
| Image is rotated or mirrored | Orientation differs from the portrait default | Override MADCTL in JSON config to flip the row/column exchange and mirror bits |
| Image is offset or clipped at edges | X/Y offset mismatch for your module | Override `xOffset` and `yOffset` in JSON config |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
