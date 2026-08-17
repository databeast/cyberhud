# Waveshare Zero LCD HAT (A)

The Waveshare Zero LCD HAT (A) is a triple-screen display board designed for Raspberry Pi Zero. It features a 1.3" 240×240 main display (ST7789) and two 0.96" 160×80 auxiliary displays (ST7735S), plus two physical buttons. CyberHUD registers this board under the panel ID `waveshare-triple-screen`.

This is the first multi-screen panel supported by cyberhud. Each screen is independently addressable and can run a different display mode simultaneously.

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

| Screen | Resolution | Controller | SPI Bus | Default Mode |
|--------|-----------|-----------|---------|--------------|
| main | 240 × 240 | ST7789 | SPI1.0 | menu |
| left-aux | 160 × 80 | ST7735S | SPI0.0 | stemma |
| right-aux | 160 × 80 | ST7735S | SPI0.1 | gpio |

All three screens use RGB565 color format.

## Pin Assignments

| Function | GPIO Pin | Screen | Notes |
|----------|----------|--------|-------|
| DC (Data/Command) | GPIO22 | main | High = data, Low = command |
| RST (Reset) | GPIO27 | main | Active-low reset |
| Backlight | GPIO19 | main | PWM-capable; high = on |
| DC (Data/Command) | GPIO4 | left-aux | High = data, Low = command |
| RST (Reset) | GPIO24 | left-aux | Active-low reset |
| Backlight | GPIO13 | left-aux | Conflicts with 3-pin connector |
| DC (Data/Command) | GPIO5 | right-aux | High = data, Low = command |
| RST (Reset) | GPIO23 | right-aux | Active-low reset |
| Backlight | GPIO12 | right-aux | PWM-capable; high = on |
| KEY1 | GPIO25 | — | Physical button 1 |
| KEY2 | GPIO26 | — | Physical button 2 |

This panel requires **two separate SPI buses**:

- **SPI0** (enabled by default) drives both auxiliary screens — left-aux on CE0 (`/dev/spidev0.0`) and right-aux on CE1 (`/dev/spidev0.1`)
- **SPI1** (must be explicitly enabled via `dtoverlay=spi1-1cs`) drives the main screen on CE0 (`/dev/spidev1.0`)

The dual-bus architecture allows all three screens to operate without contention. SPI0 uses GPIO10 (MOSI), GPIO11 (SCLK), GPIO8 (CE0), and GPIO7 (CE1). SPI1 uses GPIO20 (MOSI), GPIO21 (SCLK), and GPIO18 (CE0). See [SPI1 Enablement](#spi1-enablement) below for setup instructions.

## Independent Screen Control

Each virtual screen can be switched to a different display mode independently using `cyberhudctl`:

```sh
cyberhudctl display set 0 clock
```

```sh
cyberhudctl display set left-aux.0 stemma
```

```sh
cyberhudctl display set right-aux.0 gpio
```

You can run any combination of display modes across the three screens simultaneously.

## SPI1 Enablement

!!! warning "SPI1 is not enabled by default"

    The Raspberry Pi does not enable SPI1 by default. The main display will fail to initialize unless you explicitly enable SPI1 via device tree overlay. A reboot is required after making this change.

Add the following line to `/boot/config.txt`:

```ini
dtoverlay=spi1-1cs
```

After rebooting, verify SPI1 is available:

```sh
ls /dev/spidev1.*
```

Expected output:

```
/dev/spidev1.0
```

If this file is not present, the overlay was not loaded — double-check the spelling in `config.txt` and confirm you rebooted.

## Pin Conflict

!!! warning "Both cyberhud 3-pin GPIO connectors are unavailable"

    When using this panel, **both** cyberhud 3-pin output connectors are claimed:

    - **GPIO13** is used by the left-aux backlight control
    - **GPIO18** is used as the SPI1.0 chip-select (CE0)

    No external peripherals can be connected to either 3-pin connector while this panel is active.

**Consequence:** Any device connected to the 3-pin GPIO13 or GPIO18 headers will not function correctly, and may interfere with display operation.

Run `cyberhudctl gpio pins` to see a live report of pin conflicts for your active panel.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi1 device ... open=error` | SPI1 not enabled on the Raspberry Pi | Add `dtoverlay=spi1-1cs` to `/boot/config.txt` and reboot |
| One screen stays blank while others work | Pin conflict with another overlay or GPIO export claiming DC/RST pin | Check `cyberhudctl gpio pins` output; remove conflicting overlays or GPIO exports |
| Backlight on but screen shows no image | BL pin active but DC or RST not toggling; GPIO export conflict | Verify no other process has exported the screen's DC or RST GPIO pin |
| Screens initialize in wrong order or one times out | Multi-screen SPI bus contention during init sequence | Ensure SPI0 and SPI1 are on separate buses (not shared); restart the service to retry initialization |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
