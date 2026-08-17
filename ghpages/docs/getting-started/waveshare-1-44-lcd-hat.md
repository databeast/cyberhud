# Waveshare 1.44" LCD HAT

The Waveshare 1.44-inch LCD HAT is a 128×128 color display based on the ST7735S controller. It communicates over SPI and includes three physical buttons (KEY1–KEY3) plus a 5-way joystick (up, down, left, right, press).

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel waveshare-1.44
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "waveshare-1.44"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 128 pixels |
| Color format | RGB565 |
| Controller | ST7735S |
| Interface | SPI |
| Input | 3 buttons + 5-way joystick |

## Pin Assignments

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI MOSI | GPIO10 | SPI0 data out |
| SPI SCLK | GPIO11 | SPI0 clock |
| SPI CS | GPIO8 | SPI0 chip-select (CE0) |
| DC (Data/Command) | GPIO25 | High = data, Low = command |
| RST (Reset) | GPIO27 | Active-low reset |
| Backlight | GPIO24 | PWM-capable; high = on |
| KEY1 | GPIO21 | Physical button 1 |
| KEY2 | GPIO20 | Physical button 2 |
| KEY3 | GPIO16 | Physical button 3 |
| Joystick Up | GPIO6 | Active-low |
| Joystick Down | GPIO19 | Active-low |
| Joystick Left | GPIO5 | Active-low |
| Joystick Right | GPIO26 | Active-low |
| Joystick Press | GPIO13 | Active-low; see pin conflict note below |

SPI communication uses `/dev/spidev0.0` (CE0).

## Buttons & Joystick

| Input | GPIO Pin | Location / Label |
|-------|----------|-----------------|
| KEY1 | GPIO21 | Top-left button (closest to USB ports) |
| KEY2 | GPIO20 | Middle button |
| KEY3 | GPIO16 | Bottom-left button |
| Joystick Up | GPIO6 | Joystick pushed upward |
| Joystick Down | GPIO19 | Joystick pushed downward |
| Joystick Left | GPIO5 | Joystick pushed left |
| Joystick Right | GPIO26 | Joystick pushed right |
| Joystick Press | GPIO13 | Joystick pressed inward (center click) |

All inputs are active-low with internal pull-ups enabled by the HAT hardware.

## GPIO13 Pin Conflict

!!! warning "Joystick Press conflicts with the cyberhud 3-pin connector"

    GPIO13 is used by both the joystick press input on this HAT **and** the cyberhud 3-pin output connector. When this panel is active, the 3-pin GPIO13 connector is **unavailable** for external peripherals.

**Consequence:** Any device connected to the 3-pin GPIO13 header will not function correctly while the Waveshare 1.44" LCD HAT is in use, and joystick press events may be unreliable if an external load is present on the pin.

**Workaround:** If you need the 3-pin GPIO13 connector for another peripheral, you can disable joystick press input in your configuration and use only the three buttons and four joystick directions. The 3-pin GPIO18 connector remains free and can be used without conflict.

Run `cyberhudctl gpio pins` to see a live report of pin conflicts for your active panel.

## MADCTL & Offset Overrides

!!! note "Display orientation tuning"

    The default MADCTL value and X/Y framebuffer offsets are configured for the most common panel batch. Some batches of the ST7735S may ship with different internal rotation settings. If your display appears mirrored, shifted, or has an incorrect color order, you can override the MADCTL value and X/Y offsets in your JSON configuration file.

See [JSON Config Examples](../configuration/json-config.md) for override syntax.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Display is white/blank but backlight is on | DC or RST pin not accessible, or wrong panel name | Verify GPIO25 and GPIO27 are not claimed by another overlay; check panel name spelling |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your panel batch | Override MADCTL in JSON config to adjust color order bit |
| Image is offset or clipped at edges | X/Y offset mismatch for your panel batch | Override `xOffset` and `yOffset` in JSON config |
| Joystick directions seem swapped | Using an older build with the pin mapping bug | Update to the latest cyberhud release which corrects JoyUp/JoyDown assignments |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
