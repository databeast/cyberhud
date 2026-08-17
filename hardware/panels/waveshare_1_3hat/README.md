# waveshare-1.3hat

- Unit: Waveshare 1.3-inch LCD HAT
- Panel ID: `waveshare-1.3hat`
- Controller: ST7789
- Resolution: 240x240
- Input: 3 keys + joystick
- SKU: 14972

## Configuration

Interactive display with buttons and joystick. Recommended for most users.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.3hat"
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
| KEY1 | 5 | 29 |
| KEY2 | 6 | 31 |
| KEY3 | 13 | 33 |
| Joy-Up | 21 | 40 |
| Joy-Down | 19 | 35 |
| Joy-Left | 16 | 36 |
| Joy-Right | 20 | 38 |
| Joy-Press | 26 | 37 |
