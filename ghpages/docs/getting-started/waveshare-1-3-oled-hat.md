# Waveshare 1.3" OLED HAT

The Waveshare 1.3-inch OLED HAT is a 128×64 monochrome display based on the SH1106 controller. It communicates over SPI and includes three physical buttons (KEY1–KEY3) plus a 5-way joystick (up, down, left, right, press). As an OLED panel it is self-emitting (white on black) and requires no backlight.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel waveshare-1.3-oled-hat
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "waveshare-1.3-oled-hat"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 pixels |
| Color format | Monochrome (white on black) |
| Controller | SH1106 |
| Interface | SPI |
| Input | 3 buttons + 5-way joystick |

## Pin Assignments

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI MOSI | GPIO10 | SPI0 data out |
| SPI SCLK | GPIO11 | SPI0 clock |
| SPI CS | GPIO8 | SPI0 chip-select (CE0) |
| DC (Data/Command) | GPIO24 | High = data, Low = command |
| RST (Reset) | GPIO25 | Active-low reset |
| KEY1 | GPIO21 | Physical button 1 |
| KEY2 | GPIO20 | Physical button 2 |
| KEY3 | GPIO16 | Physical button 3 |
| Joystick Up | GPIO19 | Active-low |
| Joystick Down | GPIO6 | Active-low |
| Joystick Left | GPIO5 | Active-low |
| Joystick Right | GPIO26 | Active-low |
| Joystick Press | GPIO13 | Active-low; see pin conflict note below |

SPI communication uses `/dev/spidev0.0` (CE0). No backlight pin is required — OLED pixels are self-emitting.

## Input Details

### Buttons & Joystick

| Input | GPIO Pin | Location / Label |
|-------|----------|-----------------|
| KEY1 | GPIO21 | Top-left button (closest to USB ports) |
| KEY2 | GPIO20 | Middle button |
| KEY3 | GPIO16 | Bottom-left button |
| Joystick Up | GPIO19 | Joystick pushed upward |
| Joystick Down | GPIO6 | Joystick pushed downward |
| Joystick Left | GPIO5 | Joystick pushed left |
| Joystick Right | GPIO26 | Joystick pushed right |
| Joystick Press | GPIO13 | Joystick pressed inward (center click) |

All inputs are active-low with internal pull-ups enabled by the HAT hardware.

## GPIO13 Pin Conflict

!!! warning "Joystick Press conflicts with the cyberhud 3-pin connector"

    GPIO13 is used by both the joystick press input on this HAT **and** the cyberhud 3-pin output connector. When this panel is active, the 3-pin GPIO13 connector is **unavailable** for external peripherals.

**Workaround:** If you need the 3-pin GPIO13 connector for another peripheral, you can disable joystick press input in your configuration and use only the three buttons and four joystick directions. The 3-pin GPIO18 connector remains free and can be used without conflict.

Run `cyberhudctl gpio pins` to see a live report of pin conflicts for your active panel.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Display is completely dark | DC or RST pin not accessible, or wrong panel name | Verify GPIO24 and GPIO25 are not claimed by another overlay; check panel name spelling |
| Display shows garbled pixels | Column offset mismatch | The panel definition includes a 2-pixel column offset for the SH1106; ensure you are using the latest build |
| Partial image or shifted content | SH1106 has 132-column RAM but 128-column display | This is handled automatically; update to the latest release if you see misalignment |
| Joystick directions seem swapped | Using an older build with the pin mapping bug | Update to the latest cyberhud release which corrects JoyUp/JoyDown assignments |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
