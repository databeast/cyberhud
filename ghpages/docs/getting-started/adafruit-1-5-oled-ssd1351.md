# Adafruit 1.5" Color OLED

The [Adafruit 1.5" 16-bit Color OLED Breakout](https://www.adafruit.com/product/1431) is a 128×128 RGB OLED based on the SSD1351 controller, driven over 4-wire SPI. Being an OLED it is self-lit with very high contrast (true blacks) and no backlight. The board has an onboard microSD slot; later revisions add an EYESPI FPC connector.

CyberHUD drives the display path only; the microSD is not used, and the board has no input.

!!! note "OLED burn-in"
    OLED pixels dim with prolonged static content. For an always-on deployment, prefer modes that change the image over time (clock, ticker, attract modes) rather than a fixed static layout.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.5-oled-ssd1351
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.5-oled-ssd1351"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 128 pixels |
| Color format | RGB565 |
| Controller | SSD1351 |
| Interface | SPI |
| Input | None |

## Pin Assignments

This board is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK/CLK, MOSI/SI) and its chip-select (CS/OCS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| DC | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |

This is an OLED, so there is no backlight pin — omit `display.bl` (or set it to `"none"`). The board's microSD chip-select (SDCS) is not used by CyberHUD. Example config for a board wired with DC→GPIO25, RST→GPIO27:

```json
{
  "display": {
    "panel": "adafruit-1.5-oled-ssd1351",
    "dc": "GPIO25",
    "rst": "GPIO27"
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
| Display stays black | DC/RST mapped to the wrong GPIO, or wrong panel name | Confirm your `display.dc` / `display.rst` values match your wiring; check panel name spelling |
| Colors appear wrong or channel-swapped | Colour-order (remap) mismatch for your board | Confirm this is the SSD1351 1.5" OLED; some modules differ in RGB/BGR remap |
| Faint persistent ghost of a previous image | OLED burn-in from long static content | Use a changing display mode; turn the panel off when idle |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
