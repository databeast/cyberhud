# waveshare-2.2

- Unit: Waveshare 2.2-inch SPI display
- Panel ID: `waveshare-2.2`
- Controller: ST7789
- Resolution: 320x240
- Input: none
- SKU: pending exact confirmation

## Configuration

320×240 passive dashboard with no buttons. Perfect for wall-mounted displays or always-on dashboards.

File: `/etc/cyberhudd.json`

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-2.2",
    "disable_input": true
  }
}
```

Control via CLI from another machine:

```sh
cyberhudctl -socket /tmp/cyberhud.sock display status
cyberhudctl -socket /tmp/cyberhud.sock display set main.0 gpio
```
