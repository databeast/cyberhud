# Adafruit 4120 CharliePlex

The Adafruit 4120 is a 16×8 CharliePlex LED matrix bonnet driven by the IS31FL3731 controller. It communicates over I2C and produces a monochrome green display suitable for compact status readouts and scrolling text.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-4120-charlieplex
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-4120-charlieplex"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 16 × 8 pixels |
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
| Display shows no LEDs lit | Wrong I2C address or bus | Verify the bonnet is seated correctly; run `i2cdetect -y 1` to confirm address 0x74 appears |
| Only partial rows illuminate | Incorrect panel name or orientation mismatch | Check panel name spelling; ensure no other I2C device conflicts on the bus |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
