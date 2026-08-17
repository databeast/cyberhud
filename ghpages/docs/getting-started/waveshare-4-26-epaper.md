# Waveshare 4.26" E-Paper

The Waveshare 4.26-inch E-Paper HAT is an 800×480 monochrome (B/W) display based on the epd4in26 controller. It communicates over SPI and has no physical input buttons or joystick. Its ultra-low power consumption and wide viewing angle make it well-suited for always-on dashboards, status boards, and low-refresh information displays.

CyberHUD provides a single panel definition for this board.

## Quick Start

=== "CLI"

    ```sh
    sudo ./cyberhudd -panel waveshare-4.26-epaper
    ```

=== "JSON config"

    In your `/etc/cyberhudd.json`:

    ```json
    {
      "display": {
        "panel": "waveshare-4.26-epaper"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 800 × 480 pixels |
| Color format | Monochrome (B/W) |
| Controller | epd4in26 |
| Interface | SPI |
| Input | None |

## Pin Assignments

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| SPI MOSI | GPIO10 | SPI0 data out |
| SPI SCLK | GPIO11 | SPI0 clock |
| SPI CS | GPIO8 | SPI0 chip-select (CE0) |
| DC (Data/Command) | GPIO25 | High = data, Low = command |
| RST (Reset) | GPIO17 | Active-low reset |
| Busy | GPIO24 | High = display busy, wait before sending data |

SPI communication uses `/dev/spidev0.0` (CE0).

## Troubleshooting

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `spi device ... open=error` | SPI not enabled on the Raspberry Pi | Run `sudo raspi-config` → Interface Options → SPI → Enable, then reboot |
| Display does not update (stays blank) | Busy pin not connected or wrong GPIO | Verify GPIO24 is wired to the Busy pin on the HAT; check that no other overlay claims GPIO24 |
| Partial or corrupted image | Data sent while display was busy | Ensure nothing is interrupting the busy-wait loop; check logs for timeout warnings |
| Ghosting or faint previous image visible | Normal for e-Paper after many partial refreshes | Trigger a full refresh cycle by restarting cyberhudd or switching modes |
| Refresh takes several seconds | Expected behavior for e-Paper technology | E-Paper full refresh is inherently slow (2–4 seconds typical); use partial refresh modes where supported |
| Display shows inverse colors (white-on-black instead of black-on-white) | MADCTL or LUT mismatch for panel batch | Override color inversion setting in JSON config |

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
