# waveshare-triple-screen

- Unit: Waveshare Zero LCD HAT A (triple-screen)
- Panel ID: `waveshare-triple-screen`
- Controllers: ST7789 (main), ST7735S (left/right)
- Resolution: main 240x240, left 160x80, right 160x80
- Input: 2 keys
- SKU: 25586

## Configuration

Three coordinated displays: one main (240×240) and two auxiliary (160×80). Each display can be independently controlled.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-triple-screen"
  }
}
```

After startup, set different modes for each display:

```sh
cyberhudctl display set main.0 menu          # Interactive menu
cyberhudctl display set left-aux.0 gpio      # GPIO status
cyberhudctl display set right-aux.0 stemma   # STEMMA/I2C devices
```
