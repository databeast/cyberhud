# Adafruit 2.7" Grayscale OLED

The [Adafruit Monochrome 2.7" 128×64 OLED Graphic Display](https://www.adafruit.com/product/2674) is a 128×64 white OLED driven by an SSD1325 controller over SPI. Although the glass is a single colour, the SSD1325 drives it with 4-bit (16-level) grayscale, so CyberHUD treats it as a grayscale panel and can render smooth brightness ramps, not just on/off pixels.

CyberHUD drives the display path only; being an OLED it is self-lit with no backlight, and the module has no input.

!!! note "OLED burn-in"
    OLED pixels dim with prolonged static content. For an always-on deployment, prefer modes that change the image over time, or turn the panel off when idle.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-2.7-oled-ssd1325
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-2.7-oled-ssd1325"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 pixels |
| Color format | 4-bit grayscale (16 levels, white OLED) |
| Controller | SSD1325 |
| Interface | SPI |
| Input | None |

## Pin Assignments

This module is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the module's SPI lines (SCK/CLK, MOSI/DATA) and chip-select (CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Module pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| D/C | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |

This is an OLED, so there is no backlight pin — omit `display.bl` (or set it to `"none"`). Note the module may ship in 8-bit parallel mode; it must be set to SPI (a solder-jumper change on the back) to work with CyberHUD. Example config for a module wired with D/C→GPIO25, RST→GPIO27:

```json
{
  "display": {
    "panel": "adafruit-2.7-oled-ssd1325",
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
| Display stays black | Module still in 8-bit mode, or DC/RST on the wrong GPIO | Set the module's interface jumpers to SPI; confirm your `display.dc` / `display.rst` values |
| Image is halved, doubled, or interleaved | Grayscale nibble packing / remap mismatch | Confirm this is an SSD1325 2.7" module; the remap is tuned for it |
| Faint persistent ghost of a previous image | OLED burn-in from static content | Use a changing display mode; turn the panel off when idle |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
