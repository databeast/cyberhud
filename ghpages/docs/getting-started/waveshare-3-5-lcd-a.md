# Waveshare 3.5" RPi LCD (A)

The Waveshare 3.5-inch RPi LCD (A) is a 480×320 color SPI display based on the ILI9486 controller, with an XPT2046 resistive touch controller. The board is designed to sit directly on the Raspberry Pi 40-pin header.

CyberHUD drives the display path only — the resistive touch (XPT2046) input is not yet wired, so this panel is output-only.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel waveshare-3.5-lcd-a
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "waveshare-3.5-lcd-a"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 480 × 320 pixels (landscape) |
| Color format | RGB565 |
| Controller | ILI9486 |
| Interface | SPI |
| Input | None (touch not yet supported) |

## Pin Assignments

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI MOSI | GPIO10 | SPI0 data out |
| SPI SCLK | GPIO11 | SPI0 clock |
| SPI CS | GPIO8 | SPI0 chip-select (CE0) |
| DC (Data/Command) | GPIO24 | High = data, Low = command |
| RST (Reset) | GPIO25 | Active-low reset |
| Backlight | — | Hardwired on; no GPIO control on the (A) revision |

SPI communication uses `/dev/spidev0.0` (CE0). The pin mapping follows the standard `waveshare35a` device-tree overlay used by this board.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Display is white/blank but backlight is on | DC or RST pin not accessible, or wrong panel name | Verify GPIO24 and GPIO25 are not claimed by another overlay; check panel name spelling |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your board revision | Override MADCTL in JSON config to adjust color order bit |
| Image is rotated or mirrored | Orientation differs across board revisions | Override MADCTL in JSON config to flip the row/column exchange and mirror bits |
| Image is offset or clipped at edges | X/Y offset mismatch for your panel batch | Override `xOffset` and `yOffset` in JSON config |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
