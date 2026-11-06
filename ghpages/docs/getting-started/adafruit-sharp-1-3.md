# Adafruit 1.3" Sharp Memory LCD

The [Adafruit SHARP Memory Display Breakout 1.3" 168×144 Monochrome](https://www.adafruit.com/product/3502) is a Sharp Memory LCD — an unusual display that sits between e-paper and a normal LCD: it has the ultra-low power draw of e-paper but the fast refresh of an LCD, showing crisp black pixels on a gray background. It is reflective and daylight-readable, with no backlight (you would light the area externally for use in the dark).

CyberHUD drives the display path only; the board has no input. Unlike the TFT and OLED panels, this is a write-only 1-bit display with a bespoke line-addressed protocol rather than a standard controller command set.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-sharp-1.3
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-sharp-1.3"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 168 × 144 pixels |
| Color format | 1-bit monochrome (black on gray) |
| Controller | Sharp Memory LCD protocol |
| Interface | SPI (write-only) |
| Input | None |

The panel refreshes quickly like an LCD, but it is write-only — the host holds the framebuffer and streams whole frames. There is no backlight: the display is reflective.

## Pin Assignments

This board is a breakout, so you wire it to the Pi yourself. It needs only the SPI lines plus two Sharp-specific control signals; there is no data/command, reset, or backlight line, so the only software-configured pins are the SPI bus itself.

Wire SCLK and MOSI to a Raspberry Pi SPI bus (`SPI0` = `/dev/spidev0.0`). Two wiring quirks are specific to Sharp Memory LCDs and are handled in hardware, not in CyberHUD config:

- **CS is active HIGH** on these panels — the opposite of normal SPI chip-select. Depending on your Pi SPI setup you may need to drive CS from a GPIO held high during transfers, or invert the hardware CS line, rather than using the bus's native active-low CE.
- **VCOM / EXTCOMIN** must toggle to prevent a DC bias building up in the panel. Strap `EXTMODE` high and drive `EXTCOMIN` from a steady ~1–60 Hz source (a spare PWM/GPIO or a 555) so the panel self-refreshes its COM inversion in hardware. This matters when the image is static for long periods.
- `DISP` must be high for the display to show anything; tie it high (or to a GPIO you hold high).

Because this panel has no DC/RST/backlight lines, the `display.dc` / `display.rst` / `display.bl` overrides do not apply to it.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Nothing displays at all | `DISP` line not high | Tie `DISP` high (or to a GPIO you drive high) |
| Image appears then slowly fades or develops artifacts | `EXTCOMIN` / VCOM not toggling | Drive `EXTCOMIN` from a steady ~1 Hz+ source with `EXTMODE` high |
| Garbage / inverted pixels | CS polarity wrong (active-high panel) | Ensure CS is high during transfers; a native active-low CE will not work directly |
| Image looks mirrored or rotated | Panel mounted in a different orientation | Rotate the logical layout in config to suit the mounting |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
