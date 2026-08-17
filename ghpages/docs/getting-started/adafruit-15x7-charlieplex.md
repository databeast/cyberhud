# Adafruit 15×7 CharliePlex

The Adafruit 15×7 CharliePlex LED Matrix FeatherWing is a compact monochrome green LED display driven by the IS31FL3731 controller. It communicates over I2C and provides a 15×7 pixel grid suitable for scrolling text, simple icons, and status indicators.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-15x7-charlieplex
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-15x7-charlieplex"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 15 × 7 pixels |
| Color format | Monochrome (green LED) |
| Controller | IS31FL3731 |
| Interface | I2C |
| Input | None |

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `i2c device ... open=error` | I2C not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → I2C → Enable, then reboot |
| No LEDs light up but no error logged | Wrong I2C address or loose connection | Verify wiring and run `i2cdetect -y 1` to confirm device appears at address 0x74 |
| LEDs flicker or display corrupts | I2C bus noise or insufficient pull-ups | Use short wires, verify 3.3V supply is stable, and ensure I2C pull-up resistors are present |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
