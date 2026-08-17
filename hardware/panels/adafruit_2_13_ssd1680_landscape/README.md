# adafruit-2.13-ssd1680-landscape

- Unit: Adafruit 2.13-inch Monochrome E-Ink Bonnet
- Panel ID: `adafruit-2.13-ssd1680-landscape`
- Controller: SSD1680
- Resolution: 250x122
- Input: 2 buttons (KEY1, KEY2)
- SKU: 4687

## Configuration

Low-power e-ink display ideal for always-on battery-powered setups. Updates are slower but draw minimal current when idle.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "adafruit-2.13-ssd1680-landscape",
    "disable_input": false
  }
}
```

## GPIO Mapping

| Signal | GPIO |
|--------|------|
| SPI CS (CE0) | 8 |
| SPI MOSI | 10 |
| SPI SCLK | 11 |
| DC | 22 (or 25 for some variants) |
| RST | 27 |
| BUSY | 17 |
| KEY1 | 5 |
| KEY2 | 6 |
