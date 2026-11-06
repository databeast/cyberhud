# Adafruit 1.3" SH1106 OLED (MacroPad module)

The [Adafruit Monochrome 1.3" 128×64 SH1106G SPI OLED](https://www.adafruit.com/product/5228) is the bare OLED module used in the Adafruit MacroPad RP2040. It is a 128×64 white-on-black monochrome OLED based on the SH1106 controller, driven over 4-wire SPI.

!!! warning "Bare module"
    This product is the raw display only — no PCB, level shifting, or supporting circuitry, and a 16-pin 0.5mm FPC tail instead of header pins. You must provide the breakout/level-shifting yourself. If you want a ready-to-use 128×64 OLED, Adafruit's SSD1306 breakout or Bonnet is easier to wire.

CyberHUD drives the display path only; the module has no input, and being an OLED it is self-lit with no backlight.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.3-oled-sh1106
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.3-oled-sh1106"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 pixels |
| Color format | 1-bit monochrome (white on black) |
| Controller | SH1106 |
| Interface | SPI |
| Input | None |

## Pin Assignments

This is a bare module, so you wire it to the Pi yourself through your own carrier board and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the module's SPI lines (SCK, MOSI) and chip-select (CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Module pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| D/C | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |

This is an OLED, so there is no backlight pin — omit `display.bl` (or set it to `"none"`). Example config for a module wired with D/C→GPIO24, RST→GPIO25:

```json
{
  "display": {
    "panel": "adafruit-1.3-oled-sh1106",
    "dc": "GPIO24",
    "rst": "GPIO25"
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
| Display stays black | DC/RST mapped to the wrong GPIO, missing level shifting, or wrong panel name | Confirm your `display.dc` / `display.rst` values and that the bare module has proper 3.3V level shifting; check panel name spelling |
| Image shifted 2 px or wraps at the edge | Column offset mismatch | The SH1106 needs a 2-column offset (set by default); override `display.col_offset` if your integration differs |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
