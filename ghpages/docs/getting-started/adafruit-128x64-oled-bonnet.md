# Adafruit 128x64 OLED Bonnet

The [Adafruit 128x64 OLED Bonnet for Raspberry Pi](https://www.adafruit.com/product/3531) is a 128×64 monochrome OLED (SSD1306 controller, I2C) in a Pi bonnet form factor, with a 5-way joystick and two front buttons. It plugs straight onto the Pi's 40-pin header with a fixed pinout — no wiring needed. Being an OLED it is self-lit with high contrast and no backlight.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-128x64-oled-bonnet
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-128x64-oled-bonnet"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 64 pixels |
| Color format | 1-bit monochrome (white OLED) |
| Controller | SSD1306 |
| Interface | I2C |
| Input | 5-way joystick + 2 buttons |

## Pin Assignments

This is a Pi bonnet with a fixed pinout — the GPIOs below are wired on the board and need no configuration. The display is reached over the Pi's I2C bus at address `0x3C`; there are no DC/RST/SPI control pins.

| Function | Pin | Notes |
|----------|-----|-------|
| SDA / SCL | I2C0 | Display data (address 0x3C) |

## Buttons & Joystick

The bonnet's inputs are on fixed GPIOs (active-low, with pull-ups), usable by the input-driven modes (menu, pager, GPIO control):

| Control | GPIO Pin |
|---------|----------|
| Button A | GPIO5 |
| Button B | GPIO6 |
| Joystick Up | GPIO17 |
| Joystick Down | GPIO22 |
| Joystick Left | GPIO27 |
| Joystick Right | GPIO23 |
| Joystick Center (press) | GPIO4 |

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `i2c device ... open=error` | I2C not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → I2C → Enable, then reboot |
| No display but no error logged | Wrong I2C address or loose seating | Run `i2cdetect -y 1` to confirm the device appears at `0x3C`; reseat the bonnet |
| Display flickers or shows noise | I2C bus noise or contention | Ensure no other device shares address 0x3C; keep the bonnet firmly seated |
| Buttons or joystick do nothing | Input not enabled, or a pin claimed by another device | Confirm no other HAT/overlay uses GPIO4/5/6/17/22/23/27 |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
