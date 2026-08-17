# Waveshare 2.23" OLED HAT

The Waveshare 2.23-inch OLED HAT is a 128×32 monochrome display based on the SSD1305 controller. It supports both SPI (default) and I2C communication, selectable via solder-pad resistors on the board.

CyberHUD provides two panel definitions — one for each wiring mode.

## Panel Names

| Panel Name | Interface | Use When |
|------------|-----------|----------|
| `waveshare-2.23-oled-hat-spi` | SPI | Board is in default factory configuration |
| `waveshare-2.23-oled-hat-i2c` | I2C | Resistors resoldered for I2C mode |

## Quick Start

=== "SPI (default)"

    ```sh
    sudo ./cyberhudd -panel waveshare-2.23-oled-hat-spi
    ```

=== "I2C"

    ```sh
    sudo ./cyberhudd -panel waveshare-2.23-oled-hat-i2c
    ```

Or in your `/etc/cyberhudd.json`:

=== "SPI"

    ```json
    {
      "display": {
        "panel": "waveshare-2.23-oled-hat-spi"
      }
    }
    ```

=== "I2C"

    ```json
    {
      "display": {
        "panel": "waveshare-2.23-oled-hat-i2c"
      }
    }
    ```

## Display Characteristics

| Property | Value |
|----------|-------|
| Resolution | 128 × 32 pixels |
| Color | Monochrome (white on black) |
| Controller | SSD1305 |
| Interface | SPI (8 MHz) or I2C (address 0x3C) |
| Input | None (display only) |

This is a display-only panel with no buttons or joystick. Control it entirely via `cyberhudctl`:

```sh
# Switch display mode
cyberhudctl display set 0 clock

# Check status
cyberhudctl display status
```

## Pin Assignments

### SPI Mode

| Function | GPIO Pin | Notes |
|----------|----------|-------|
| DC (Data/Command) | GPIO24 | Physical pin 18 |
| RST (Reset) | GPIO25 | Physical pin 22 |
| BL (Backlight) | GPIO18 | Not typically used (OLED is self-emitting); active-low |
| BUSY | — | Not used |

SPI communication uses `/dev/spidev0.0` at 1 MHz, SPI Mode 3 (CPOL=1, CPHA=1).

### I2C Mode

| Parameter | Value |
|-----------|-------|
| Bus | `/dev/i2c-1` |
| Device address | `0x3C` |
| GPIO pins | None required |

No GPIO pins are needed in I2C mode — all communication occurs over the I2C bus.

## Hardware Setup

### Identifying Your Interface Mode

The Waveshare 2.23" OLED HAT ships with **SPI as the default communication mode**. I2C mode requires physically resoldering resistor jumpers on the back of the board — it cannot be switched in software.

- **SPI mode (factory default)**: Zero-ohm resistors connect BS1/BS2/DIN/CLK/CS/DC to the SPI signals. No modification needed.
- **I2C mode (requires soldering)**: You must move the zero-ohm resistors to the I2C pads. This connects BS1 to 3V3, DIN to SDA, CLK to SCL, and ties CS/DC to GND. See the [Waveshare wiki](https://www.waveshare.com/wiki/2.23inch_OLED_HAT) for the soldering diagram.

!!! warning
    Do NOT use the I2C panel name (`waveshare-2.23-oled-hat-i2c`) unless you have physically resoldered the jumpers. The board will not respond over I2C in its factory SPI configuration.

!!! tip
    If you haven't modified the board, use the SPI panel: `waveshare-2.23-oled-hat-spi`

### Enabling I2C on the Raspberry Pi

If using I2C mode, ensure I2C is enabled:

```sh
sudo raspi-config
# Navigate to: Interface Options → I2C → Enable
```

Verify the device is visible:

```sh
i2cdetect -y 1
```

You should see address `0x3C` in the output.

### Enabling SPI on the Raspberry Pi

If using SPI mode, ensure SPI is enabled:

```sh
sudo raspi-config
# Navigate to: Interface Options → SPI → Enable
```

Verify the SPI device exists:

```sh
ls /dev/spidev0.*
```

## Troubleshooting

### Display Shows Nothing

CyberHUD emits detailed diagnostic logs during panel initialization. Check them with:

```sh
sudo journalctl -u cyberhudd.service -n 50
```

The logs report each step: bus probing, pin resolution, driver factory invocation, and success/failure outcomes. Example output for a successful I2C init:

```
display: init panel="waveshare-2.23-oled-hat-i2c" controller=ssd1305 mode=I2C
display: available buses i2c=[/dev/i2c-1] spi=[SPI0.0 SPI0.1]
display: pin assignment DC=unassigned RST=unassigned BL=unassigned BUSY=unassigned
display: i2c bus="/dev/i2c-1" addr=0x3C open=ok
display: driver factory id=ssd1305 type=I2C config=128x32
display: waveshare-2.23-oled-hat-i2c ready (128x32 SSD1305)
```

### Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `i2c bus ... open=error` | I2C not enabled or wrong bus | Run `raspi-config` to enable I2C |
| `spi device ... open=error` | SPI not enabled | Run `raspi-config` to enable SPI |
| `required pin missing name=GPIO8 type=DC` | GPIO pin not accessible | Check for pin conflicts with `cyberhudctl gpio pins` |
| `unsupported display controller "ssd1305"` | Driver not registered | Ensure you're running a build with SSD1305 support |
| Display inverted or garbled | Wrong interface mode selected | Verify board resistor pads match your panel choice |

### Switching Between SPI and I2C

If you change the board's interface resistors, update your config accordingly. Using the wrong panel name for your wiring will result in initialization failures — the logs will indicate whether the SPI port or I2C bus couldn't be opened.

## Related Pages

- [Hardware Identification](hardware.md) — find your panel name
- [Pin Assignments](../reference/pin-assignments.md) — custom pin overrides
- [JSON Config Examples](../configuration/json-config.md) — full configuration file setup
- [Troubleshooting](troubleshooting.md) — general diagnostic steps
