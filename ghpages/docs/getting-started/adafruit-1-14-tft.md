# Adafruit 1.14" Color TFT (Newxie)

The [Adafruit 1.14" 240×135 Color "Newxie" TFT](https://www.adafruit.com/product/6113) is a thumbnail-sized 240×135 IPS colour TFT based on the ST7789 controller, driven over 4-wire SPI with a backlight. The breakout is designed for vertical "nixie-tube" stacking on a 0.1" grid. It is a minimized board — no microSD and no reset line — needing only 7 pins.

CyberHUD drives the display path only; the board has no input.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.14-tft
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.14-tft"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 240 × 135 pixels (landscape) |
| Color format | RGB565 |
| Controller | ST7789 |
| Interface | SPI |
| Input | None |

## Pin Assignments

This is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC (and RST, see below) as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK, SDA/MOSI) and chip-select (CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| DC | `display.dc` | `-display-dc` | Data/command select (high = data) |
| BL | `display.bl` | `-display-bl` | Backlight (optional; tie high for always-on, or set `bl` to `"none"`) |

**Reset:** this minimized board does **not** break out a reset line. CyberHUD's ST7789 driver still requires a reset pin to be configured, so set `display.rst` to any spare, unused GPIO — it is harmless (nothing is wired to it) and the ST7789's SPI software reset handles controller initialization. Example config with DC→GPIO25, backlight→GPIO24, and RST pointed at an unused GPIO27:

```json
{
  "display": {
    "panel": "adafruit-1.14-tft",
    "dc": "GPIO25",
    "rst": "GPIO27",
    "bl": "GPIO24"
  }
}
```

Those GPIO numbers are an example, not a requirement — use whatever you wired. The `rst` GPIO just needs to be free; it does not connect to the board.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| `pin not found` at startup | DC or RST override not supplied | Set `display.dc` to your wired DC pin and `display.rst` to any free GPIO (no reset line is broken out) |
| Display is white/blank but backlight is on | DC mapped to the wrong GPIO, or wrong panel name | Confirm your `display.dc` value matches your wiring; check panel name spelling |
| Image is shifted or wraps at an edge | Framebuffer offset mismatch | Tune `display.x_offset` / `display.y_offset` (defaults target the standard 240×135 window) |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your module | Override MADCTL (`display.madctl`) to adjust colour order |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
