# Adafruit 1.3" Color TFT Bonnet

The [Adafruit 1.3" Color TFT Bonnet for Raspberry Pi](https://www.adafruit.com/product/4506) is a 240×240 IPS colour TFT (ST7789 controller, SPI) in a Raspberry Pi bonnet form factor, with a 5-way joystick and two front buttons. Unlike the bare breakouts, it is a HAT that plugs straight onto the Pi's 40-pin header with a fixed pinout — no wiring needed.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.3-tft-bonnet
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.3-tft-bonnet"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 240 × 240 pixels |
| Color format | RGB565 |
| Controller | ST7789 |
| Interface | SPI |
| Input | 5-way joystick + 2 buttons |

## Pin Assignments

This is a Pi bonnet with a fixed pinout — the GPIOs below are wired on the board and need no configuration.

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI SCLK / MOSI / CS | SPI0 (CE0) | Display data; MISO is unused (write-only) |
| DC (Data/Command) | GPIO25 | High = data, Low = command |
| Backlight | GPIO26 | High = on |
| RST (Reset) | GPIO24 | Not wired on the bonnet; reset is handled by the controller's SPI software reset |

The bonnet does not break out a hardware reset line, so CyberHUD relies on the ST7789's software reset during initialization. A STEMMA QT / Qwiic I2C connector is present for sensors but is not used by the display.

## Buttons & Joystick

The bonnet's inputs are on fixed GPIOs (active-low, with pull-ups), usable by the input-driven modes (menu, pager, GPIO control):

| Control | GPIO Pin |
|---------|----------|
| Button A | GPIO5 |
| Button B | GPIO6 |
| Joystick Up | GPIO17 |
| Joystick Down | GPIO22 |
| Joystick Left | GPIO27 |
| Joystick Right | GPIO23 |
| Joystick Center (press) | GPIO4 |

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Display is white/blank but backlight is on | SPI not enabled, or another overlay claiming the pins | Verify SPI is enabled and no conflicting display overlay is loaded |
| Colors appear inverted (BGR vs RGB) | MADCTL mismatch for your board batch | Override MADCTL (`display.madctl`) to adjust the colour-order bit |
| Buttons or joystick do nothing | Input not enabled, or a pin claimed by another device | Confirm no other HAT/overlay uses GPIO4/5/6/17/22/23/27 |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
