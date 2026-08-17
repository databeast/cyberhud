# Hardware Identification

This page helps you identify your display hardware and find the correct panel name to use with cyberhud. If you haven't assembled your hardware yet, start with the [Setup](setup.md) guide for wiring instructions and first-boot steps.

## Supported Display Panels

The following table lists every panel registered in the panel registry (`panels.Names()`). Use the **Panel ID** value as the `-panel` argument or `"panel"` JSON config field.

| Panel ID | Controller | Resolution | Color Mode | Refresh | Inputs |
|----------|------------|------------|------------|---------|--------|
| [`adafruit-15x7-charlieplex`](adafruit-15x7-charlieplex.md) | IS31FL3731 | 15×7 | grayscale | continuous | none |
| [`adafruit-2.13-ssd1680`](adafruit-2-13-ssd1680.md) | SSD1680 | 122×250 | monochrome | deferred | keys |
| [`adafruit-4120-charlieplex`](adafruit-4120-charlieplex.md) | IS31FL3731 | 16×8 | grayscale | continuous | none |
| [`st7789-240x135`](st7789-240x135.md) | ST7789 | 240×135 | color | continuous | none |
| [`st7789-240x240`](st7789-240x240.md) | ST7789 | 240×240 | color | continuous | none |
| [`st7789-320x240`](st7789-320x240.md) | ST7789 | 320×240 | color | continuous | none |
| [`waveshare-1.3-oled-hat`](waveshare-1-3-oled-hat.md) | SH1106 | 128×64 | monochrome | continuous | keys + joystick |
| [`waveshare-1.3hat`](waveshare-1-3hat.md) | ST7789 | 240×240 | color | continuous | keys + joystick |
| [`waveshare-1.44`](waveshare-1-44-lcd-hat.md) | ST7735S | 128×128 | color | continuous | keys + joystick |
| [`waveshare-2.2`](waveshare-2-2.md) | ST7789 | 320×240 | color | continuous | none |
| [`waveshare-2.23-oled-hat-i2c`](waveshare-2-23-oled-hat.md) | SSD1305 | 128×32 | monochrome | continuous | none |
| [`waveshare-2.23-oled-hat-spi`](waveshare-2-23-oled-hat.md) | SSD1305 | 128×32 | monochrome | continuous | none |
| [`waveshare-4.26-epaper`](waveshare-4-26-epaper.md) | EPD4in26 | 800×480 | monochrome | deferred | none |
| [`waveshare-triple-screen`](waveshare-triple-screen.md) | ST7789 + 2×ST7735S | 240×240 (main) + 2×160×80 (aux) | color | continuous | keys |

### Column Definitions

- **Controller**: The display controller chip driving the panel hardware.
- **Resolution**: Native pixel dimensions (width×height). For multi-screen panels, each screen's resolution is listed.
- **Color Mode**: `color` for RGB/full-color TFT, `monochrome` for single-bit (B/W or B/W/R e-ink, OLED), `grayscale` for PWM-driven LED matrices.
- **Refresh**: `continuous` for panels that accept frame updates at any time; `deferred` for e-paper panels that batch updates with a visible refresh cycle.
- **Inputs**: `keys` = button GPIO pins assigned; `joystick` = directional + press GPIO pins assigned; `keys + joystick` = both; `none` = no onboard input hardware.

### Input Details

Panels with `InputPins` defined provide physical controls for on-device navigation (menu mode, mode cycling):

| Panel ID | Keys | Joystick |
|----------|------|----------|
| `adafruit-2.13-ssd1680` | Key1, Key2 | — |
| `waveshare-1.3-oled-hat` | Key1, Key2, Key3 | Up, Down, Left, Right, Press |
| `waveshare-1.3hat` | Key1, Key2, Key3 | Up, Down, Left, Right, Press |
| `waveshare-1.44` | Key1, Key2, Key3 | Up, Down, Left, Right, Press |
| `waveshare-triple-screen` | Key1, Key2 | — |

## Hardware Identification Commands

Use the `cyberhudctl` CLI to discover and identify hardware:

### List all available display modes

```sh
cyberhudctl display modes
```

This prints every available display mode registered in the mode catalog.

### Check current display status

```sh
cyberhudctl display status
```

Shows the current state of all configured display regions, including active modes and available alternatives.

## Choosing a Panel

- **If your display has buttons or a joystick**, use the matching HAT panel (e.g., `waveshare-1.3hat`). This enables the on-display interactive menu.
- **If your display has no buttons**, use a display-only panel (e.g., `waveshare-2.2`, `st7789-240x240`). The daemon runs a passive dashboard and you control it via [CLI](cli.md).
- **If you have multiple displays**, use the `waveshare-triple-screen` panel and address each screen independently. See [Regions](regions.md) for how to target individual screens.

## Custom GPIO Wiring

If your panel uses different GPIO pins than the default, override them at startup:

```sh
sudo ./cyberhudd \
  -panel waveshare-2.2 \
  -display-dc GPIO25 \
  -display-rst GPIO27 \
  -display-bl GPIO24
```

For full GPIO pin mapping details, see [Pin Assignments](../reference/pin-assignments.md).

## Next Steps

- Ready to set up? Follow the [Setup](setup.md) guide
- Once running, control your cyberhud from the terminal: [CLI Usage](cli.md)
- Persist your settings with a [JSON config file](../configuration/json-config.md)
- Running into issues? See [Troubleshooting](troubleshooting.md)
