# Adafruit 1.54" Tri-Color eInk Breakout

The [Adafruit 1.54" Tri-Color eInk / ePaper Display with SRAM](https://www.adafruit.com/product/4868) is a 200×200 black/red/white e-paper breakout based on the SSD1681 controller. It is the tri-color sibling of the [monochrome 1.54" board](adafruit-1-54-ssd1681.md), and carries an onboard SRAM buffer, a microSD socket, and an 18-pin EYESPI FPC connector.

CyberHUD drives the display path only. The board has no buttons, so this panel is output-only, and the onboard SRAM and microSD are not used — CyberHUD buffers frames on the host and streams the black and red ink planes to the controller directly.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel adafruit-1.54-ssd1681-tricolor
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "adafruit-1.54-ssd1681-tricolor"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 200 × 200 pixels (square) |
| Color format | Tri-color (black / red / white) |
| Controller | SSD1681 |
| Interface | SPI |
| Input | None |

Tri-color e-paper refreshes slowly — a full update takes several seconds — and does not support partial refresh, so every changed frame triggers a full refresh. Pixels are classified into three inks: a strongly red pixel becomes red, an otherwise dark pixel becomes black, and everything else stays white.

## Pin Assignments

This board is a breakout, not a Raspberry Pi HAT, so there is no fixed pinout — you wire it to the Pi yourself and tell CyberHUD which GPIOs you used. The panel definition ships with **no default pins**; you must supply DC, RST, and BUSY as config overrides, or the daemon will exit at startup with a `pin not found` error.

Wire the board's SPI lines (SCK, MOSI) and its chip-select (ECS) to a Raspberry Pi SPI bus — `SPI0` CE0 is `/dev/spidev0.0`. Chip-select is selected by the SPI bus you choose, not by a config key. Then wire these three control lines to any free GPIOs and declare them:

| Board pin | Config key | CLI flag | Purpose |
|-----------|-----------|----------|---------|
| D/C | `display.dc` | `-display-dc` | Data/command select (high = data) |
| RST | `display.rst` | `-display-rst` | Active-low reset |
| BUSY | `display.busy` | `-display-busy` | Reads high while the controller refreshes |

The board's SRAM (SRCS) and microSD (SDCS) chip-selects are not used by CyberHUD and can be left unconnected. Example config for a board wired with D/C→GPIO22, RST→GPIO27, BUSY→GPIO17:

```json
{
  "display": {
    "panel": "adafruit-1.54-ssd1681-tricolor",
    "dc": "GPIO22",
    "rst": "GPIO27",
    "busy": "GPIO17"
  }
}
```

Those GPIO numbers are an example, not a requirement — use whatever you wired.

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| `ssd1681-tricolor busy timeout` | BUSY pin not configured/accessible, or wrong `busy_high` polarity | Confirm `display.busy` points at the GPIO you wired to the board's BUSY line; some modules idle the busy line inverted |
| `pin not found` at startup | DC, RST, or BUSY override not supplied | This panel has no default pins — set `display.dc`, `display.rst`, and `display.busy` to the GPIOs you wired |
| Display never updates | DC/RST/BUSY mapped to the wrong GPIO, or wrong panel name | Confirm your `display.dc` / `display.rst` / `display.busy` values match your wiring; check panel name spelling |
| Red and black inks are swapped or missing | Content maps to the wrong ink plane, or a mono-only module | Confirm the board is the tri-color (4868) variant, not the monochrome (4196) one |
| Image appears rotated or mirrored | Orientation differs from the square default | Rotate the logical layout in config, or re-seat the EYESPI cable orientation |

## Related Pages

- [1.54" Monochrome eInk](adafruit-1-54-ssd1681.md) — the black/white sibling board
- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
