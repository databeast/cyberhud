# adafruit-4120-charlieplex

- Unit: Adafruit 8x16 CharliePlex LED Matrix Bonnet (Green)
- Panel ID: `adafruit-4120-charlieplex`
- Controller: IS31FL3731
- Resolution: 16x8
- Bus: I2C (default address 0x74)
- Input: None
- SKU: 4120

## Configuration

Compact 16×8 LED matrix bonnet for status indicators or scrolling text. Communicates over I2C — no SPI or GPIO pins required beyond the I2C bus. The IS31FL3731's non-linear charlieplex pixel wiring is handled automatically by the driver (`charlie-bonnet` layout).

For a bonnet mounted rotated 90° clockwise (portrait 8×16), set `"orientation": "cw"` in the display config.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-4120-charlieplex",
    "disable_input": true
  }
}
```

### Custom I2C address

If you've soldered the address select pads on the bonnet, specify the address:

```json
{
  "display": {
    "panel": "adafruit-4120-charlieplex",
    "address": "0x77"
  }
}
```

Available addresses: 0x74 (default), 0x75, 0x76, 0x77.
