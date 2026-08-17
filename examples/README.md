# CyberHUD Configuration Examples

Example `cyberhudd.json` configuration files for common deployment scenarios. Copy one as a starting point and customize to match your hardware setup.

## Usage

```sh
sudo cp examples/waveshare-1.3hat-desktop.json /etc/cyberhudd.json
sudo cyberhudd -config /etc/cyberhudd.json
```

Any CLI flag overrides the corresponding config file value.

## Examples

| File | Panel | Use Case |
|------|-------|----------|
| `waveshare-1.3hat-desktop.json` | Waveshare 1.3" LCD HAT | Desktop cyberdeck with menu navigation and attract modes |
| `waveshare-triple-screen.json` | Waveshare Triple Screen | Multi-display cyberdeck with per-screen orientation |
| `waveshare-2.23-oled-status.json` | Waveshare 2.23" OLED (I2C) | Compact monochrome status ticker |
| `eink-kiosk.json` | Adafruit 2.13" E-Ink | Low-power e-ink dashboard with slow refresh |
| `headless-sensor-hub.json` | (none) | No display, I2C sensor monitoring only |
| `generic-320x240-attract.json` | ST7789 320x240 | Full-screen attract mode showcase |
| `waveshare-1.44-clock.json` | Waveshare 1.44" LCD | Dedicated clock display with custom styling |
| `multi-bus-lab.json` | ST7789 240x240 | Multi-I2C-bus lab bench sensor hub |
| `policies-showcase.json` | Waveshare 1.3" LCD HAT | All mode policies configured in one file |

## Schema Reference

See [CONFIGURATION.md](../CONFIGURATION.md) for the full schema and field descriptions.
