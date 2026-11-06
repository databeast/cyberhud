# Adafruit 1.8" Color TFT

The [Adafruit 1.8" Color TFT LCD with microSD Breakout](https://www.adafruit.com/product/358) is a 128×160 RGB TFT based on the ST7735R controller, driven over 4-wire SPI. It is a true TFT (full 18-bit colour, fast refresh), with a backlight and an onboard microSD slot; later revisions add an EYESPI FPC connector.

CyberHUD drives the display path only; the microSD is not used, and the board has no input. The ST7735R shares the ST7735 initialisation and command set with the ST7735S, so CyberHUD drives it with the existing ST7735S driver.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.8-tft
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.8-tft"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 160 pixels (portrait) |
| Color format | RGB565 |
| Controller | ST7735R |
| Interface | SPI |
| Input | None |

## Pin Assignments

This board is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK/CLK, MOSI/SI) and its chip-select (TCS/CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| D/C | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |
| LITE | `display.bl` | `-display-bl` | Backlight (optional; tie high for always-on, or set `bl` to `"none"`) |

The board's microSD chip-select (CCS) and card-detect are not used by CyberHUD. Example config for a board wired with D/C→GPIO25, RST→GPIO27, backlight→GPIO24:

```json
{
  "display": {
    "panel": "adafruit-1.8-tft",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24"
  }
}
```

Those GPIO numbers are an example, not a requirement — use whatever you wired.

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
| Image is shifted a few pixels / wraps at an edge | Framebuffer offset differs for your panel batch | Tune `display.x_offset` / `display.y_offset` (these 1.8" panels vary slightly by batch) |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your module | Override MADCTL (`display.madctl`) to adjust colour order |
| Image is rotated or mirrored | Orientation differs from the portrait default | Override MADCTL to flip the row/column exchange and mirror bits |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
