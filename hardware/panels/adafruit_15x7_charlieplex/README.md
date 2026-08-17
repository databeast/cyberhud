# adafruit-15x7-charlieplex

- Unit: Adafruit 15x7 CharliePlex LED Matrix Display FeatherWing
- Panel ID: `adafruit-15x7-charlieplex`
- Controller: IS31FL3731
- Resolution: 15x7
- Bus: I2C (default address 0x74)
- Input: None
- SKU: 2965

## Configuration

Compact 15×7 LED matrix for status indicators or scrolling text. Communicates over I2C — no SPI or GPIO pins required beyond the I2C bus. The IS31FL3731's non-linear charlieplex pixel wiring is handled automatically by the driver (`charlie-wing` layout).

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-15x7-charlieplex",
    "disable_input": true
  }
}
```

### Custom I2C address

If you've soldered the address select pads on the FeatherWing, specify the address:

```json
{
  "display": {
    "panel": "adafruit-15x7-charlieplex",
    "address": "0x77"
  }
}
```

Available addresses: 0x74 (default), 0x75, 0x76, 0x77.
