# Adafruit 2.13" E-Paper (SSD1680)

The Adafruit 2.13-inch Monochrome E-Ink Bonnet is a 122×250 monochrome e-Paper display based on the SSD1680 controller. It communicates over SPI and includes two physical buttons (KEY1–KEY2). E-Paper displays retain their image with no power draw, making this panel well-suited for low-refresh dashboards and status screens.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-2.13-ssd1680
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-2.13-ssd1680"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 122 × 250 pixels (250 × 122 landscape) |
| Color format | Monochrome (B/W) |
| Controller | SSD1680 |
| Interface | SPI |
| Input | 2 buttons |

## Pin Assignments

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI MOSI | GPIO10 | SPI0 data out |
| SPI SCLK | GPIO11 | SPI0 clock |
| SPI CS | GPIO8 | SPI0 chip-select (CE0) |
| DC (Data/Command) | GPIO22 | High = data, Low = command |
| RST (Reset) | GPIO27 | Active-low reset |
| Busy | GPIO17 | High = busy (display is refreshing) |
| KEY1 | GPIO5 | Physical button 1 |
| KEY2 | GPIO6 | Physical button 2 |

SPI communication uses `/dev/spidev0.0` (CE0).

## Input Details

| Input | GPIO Pin | Location / Label |
|-------|----------|-----------------|
| KEY1 | GPIO5 | Button 1 |
| KEY2 | GPIO6 | Button 2 |

Both inputs are active-low with internal pull-ups enabled by the bonnet hardware. These buttons can be used for mode switching, menu navigation, or user-defined actions via CyberHUD's input mapping.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| Display never updates (stays blank) | Busy pin not wired or wrong GPIO | Verify GPIO17 is connected to the Busy pin on the bonnet; check `cyberhudd` logs for busy-timeout errors |
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Ghosting or shadow of previous image | Normal for e-Paper after many partial refreshes | CyberHUD performs a full refresh every 20 partial updates automatically; force one by restarting the service |
| Refresh is very slow (several seconds) | Expected for e-Paper technology | E-Paper refreshes are inherently slow (~2–4 s for full); use partial refresh modes where possible |
| Image appears inverted (white-on-black when expecting black-on-white) | Display polarity mismatch | Check panel orientation settings in your JSON config; try `"orientation": "flip"` to rotate the display 180° |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
