# Adafruit 1.28" Round TFT

The [Adafruit 1.28" 240×240 Round TFT LCD with MicroSD](https://www.adafruit.com/product/6178) is a round IPS SPI display based on the GC9A01A controller — the same round 240×240 panel widely used in smartwatch-style projects. It has an EYESPI FPC connector and an onboard microSD slot.

CyberHUD drives the display path only; the microSD is not used, and the board has no input. The panel is physically round, but the framebuffer is a full 240×240 square — the corner pixels are addressable, they are simply hidden behind the bezel, so lay out content with the visible circle in mind.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.28-round
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.28-round"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 240 × 240 pixels (round; square framebuffer) |
| Color format | RGB565 |
| Controller | GC9A01A |
| Interface | SPI |
| Input | None |

## Pin Assignments

This board is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC and RST as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK/CLK, MOSI/SI) and its chip-select (CS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire the control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| DC | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |
| BL | `display.bl` | `-display-bl` | Backlight (optional; tie high for always-on, or set `bl` to `"none"`) |

The board's microSD chip-select is not used by CyberHUD. Example config for a board wired with DC→GPIO25, RST→GPIO27, backlight→GPIO24:

```json
{
  "display": {
    "panel": "adafruit-1.28-round",
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
| Colors appear inverted or wrong | MADCTL mismatch for your board | Override MADCTL in JSON config (`display.madctl`) to adjust orientation and colour order |
| Image is rotated for how the panel is mounted | Round panels have no inherent "up" | Override MADCTL (`0x08` / `0x68` / `0xC8` / `0xA8` for 0/90/180/270°) to rotate |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
