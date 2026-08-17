# Waveshare Triple Screen

The Waveshare Triple Screen is a composite panel combining three individual displays into a single logical unit: a 1.3-inch ST7789 main screen (240×240) flanked by two 0.96-inch ST7735S auxiliary screens (160×80 each). This setup spans two SPI buses (SPI0 and SPI1), giving CyberHUD a wide multi-panel workspace for simultaneous data views.

CyberHUD treats the three screens as a single virtual panel. Each screen can display a different mode independently.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel waveshare-triple-screen
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "waveshare-triple-screen"
      }
    }
    ```

## Display Characteristics

| Property | Screen 0 (main) | Screen 1 (left-aux) | Screen 2 (right-aux) |
|----------|-----------------|----------------------|----------------------|
| Resolution | 240 × 240 pixels | 160 × 80 pixels | 160 × 80 pixels |
| Size | 1.3 inch | 0.96 inch | 0.96 inch |
| Color format | RGB565 | RGB565 (inverted) | RGB565 (inverted) |
| Controller | ST7789 | ST7735S | ST7735S |
| Interface | SPI1.0 | SPI0.0 | SPI0.1 |
| Default mode | menu | stemma | gpio |

## Pin Assignments

### Screen 0 — Main (ST7789, SPI1.0)

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI1 MOSI | GPIO20 | SPI1 data out |
| SPI1 SCLK | GPIO21 | SPI1 clock |
| SPI1 CS (CE0) | GPIO18 | SPI1 chip-select 0 |
| DC (Data/Command) | GPIO22 | High = data, Low = command |
| RST (Reset) | GPIO27 | Active-low reset |
| Backlight | GPIO19 | PWM-capable; high = on |

### Screen 1 — Left Aux (ST7735S, SPI0.0)

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI0 MOSI | GPIO10 | SPI0 data out |
| SPI0 SCLK | GPIO11 | SPI0 clock |
| SPI0 CS (CE0) | GPIO8 | SPI0 chip-select 0 |
| DC (Data/Command) | GPIO4 | High = data, Low = command |
| RST (Reset) | GPIO24 | Active-low reset |
| Backlight | GPIO13 | PWM-capable; high = on |

### Screen 2 — Right Aux (ST7735S, SPI0.1)

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI0 MOSI | GPIO10 | SPI0 data out (shared with Screen 1) |
| SPI0 SCLK | GPIO11 | SPI0 clock (shared with Screen 1) |
| SPI0 CS (CE1) | GPIO7 | SPI0 chip-select 1 |
| DC (Data/Command) | GPIO5 | High = data, Low = command |
| RST (Reset) | GPIO23 | Active-low reset |
| Backlight | GPIO12 | PWM-capable; high = on |

### Input Keys

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| KEY1 | GPIO25 | Physical button 1 |
| KEY2 | GPIO26 | Physical button 2 |

## Input Details

| Input | GPIO Pin | Description |
|-------|----------|-------------|
| KEY1 | GPIO25 | User button 1, active-low |
| KEY2 | GPIO26 | User button 2, active-low |

Both inputs are active-low with internal pull-ups enabled by the HAT hardware.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` for SPI1 | SPI1 overlay not loaded | Add `dtoverlay=spi1-1cs` to `/boot/config.txt` and reboot |
| Only main screen works, aux screens blank | SPI0 not enabled | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Only aux screens work, main screen blank | SPI1 not enabled or wrong CS pin | Verify `dtoverlay=spi1-1cs` is present and GPIO18 is not claimed by another overlay |
| One aux screen works but the other is blank | CE1 not available on SPI0 | Ensure no other overlay claims GPIO7 (SPI0 CE1) |
| Colors appear inverted on aux screens | Expected behavior | The ST7735S aux screens use `InvertColors: true` by default; override MADCTL in JSON config if your batch differs |
| Image is offset or clipped on aux screens | X/Y offset mismatch | Override `xOffset` and `yOffset` in JSON config (defaults: X=1, Y=26) |
| Display appears mirrored | MADCTL mismatch for your panel batch | Override MADCTL in JSON config to adjust mirror and color order bits |

### Multi-SPI Bus Notes

This panel is unique in that it uses **both SPI0 and SPI1** simultaneously:

- **SPI0** (`/dev/spidev0.0` and `/dev/spidev0.1`) drives the two auxiliary screens. They share MOSI and SCLK lines but use separate chip-selects (CE0 and CE1).
- **SPI1** (`/dev/spidev1.0`) drives the main screen independently.

Both SPI buses must be enabled in your Raspberry Pi configuration. If you only see partial output, verify that the relevant `/dev/spidevX.Y` device nodes exist.

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
