# Adafruit 13×9 RGB LED Matrix

The [Adafruit IS31FL3741 13×9 PWM RGB LED Matrix](https://www.adafruit.com/product/5201) is a compact full-colour LED display: 117 RGB LEDs in a 13×9 grid, driven by the ISSI IS31FL3741 over I2C with 8-bit PWM per channel (24-bit colour). It has STEMMA QT / Qwiic connectors for solderless wiring.

CyberHUD provides a single panel definition for this board. It is best suited to colourful status indicators, small icons, and scrolling text rather than detailed graphics.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-13x9-is31fl3741
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-13x9-is31fl3741"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 13 × 9 pixels |
| Color format | 24-bit RGB (8-bit PWM per channel) |
| Controller | IS31FL3741 |
| Interface | I2C (STEMMA QT / Qwiic) |
| Input | None |

The board ships at I2C address `0x30`. Up to four boards can share one bus by strapping the address jumpers to `0x30`, `0x31`, `0x33`, or `0x34`; CyberHUD drives a single matrix at the default address. Because every pixel is written over I2C, refresh is modest — fine for simple animations, not video.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `i2c device ... open=error` | I2C not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → I2C → Enable, then reboot |
| No LEDs light up but no error logged | Wrong I2C address or loose connection | Verify wiring and run `i2cdetect -y 1` to confirm the device appears at `0x30` |
| Colors look wrong (red/blue swapped) | Board colour order differs from the BGR default | Confirm this is the STEMMA QT matrix; other IS31FL3741 boards may use a different channel order |
| LEDs flicker or display corrupts | I2C bus noise or insufficient pull-ups | Use short wires, verify a stable supply (5V recommended for best green/blue), and ensure I2C pull-ups are present |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
