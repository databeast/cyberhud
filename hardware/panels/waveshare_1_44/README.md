# waveshare-1.44

- Unit: Waveshare 1.44-inch LCD HAT
- Panel ID: `waveshare-1.44`
- Controller: ST7735S
- Resolution: 128x128
- Input: 3 keys + joystick
- SKU: 13891

## Configuration

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.44"
  }
}
```

### Display Rotation/Correction

Some panel batches ship with different MADCTL (rotation/color) settings. Correct them in config:

```json
{
  "display": {
    "panel": "waveshare-1.44",
    "madctl": "0xC8"
  }
}
```

Or with X/Y offsets for panels with shifted framebuffer:

```json
{
  "display": {
    "panel": "waveshare-1.44",
    "x_offset": 2,
    "y_offset": 3
  }
}
```

## GPIO Mapping

| Signal | GPIO | RPi pin |
|--------|------|---------|
| SPI MOSI | 10 | 19 |
| SPI SCLK | 11 | 23 |
| SPI CS | 8 | 24 |
| DC | 25 | 22 |
| RST | 27 | 13 |
| Backlight | 24 | 18 |
| KEY1 | 21 | 40 |
| KEY2 | 20 | 38 |
| KEY3 | 16 | 36 |
| Joy-Up | 6 | 31 |
| Joy-Down | 19 | 35 |
| Joy-Left | 5 | 29 |
| Joy-Right | 26 | 37 |
| Joy-Press | 13 | 33 |
