# Pin Assignments

Each panel profile has built-in GPIO pin assignments that match the manufacturer's reference wiring. When you select a panel (via `-panel` flag or `display.panel` in your config), the daemon automatically uses the correct pins for that hardware.

This page documents the default assignments and the connector conflict diagnostic.

## Panel Default Pins

### Waveshare SPI Displays

| Panel | DC | RST | BL | Bus |
|-------|-----|------|-----|------|
| waveshare-1.3hat | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |
| waveshare-1.44 | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |
| waveshare-2.2 | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |
| waveshare-1.3-oled-hat | GPIO24 | GPIO25 | — | SPI0.0 |
| waveshare-2.23-oled-hat-spi | GPIO24 | GPIO25 | GPIO18 | SPI0.0 |

### Waveshare E-Paper

| Panel | DC | RST | Busy | Bus |
|-------|-----|------|------|------|
| waveshare-4.26-epaper | GPIO25 | GPIO17 | GPIO24 | SPI0.0 |

### Adafruit Displays

| Panel | DC | RST | Busy | Bus |
|-------|-----|------|------|------|
| adafruit-2.13-ssd1680 | GPIO22 | GPIO27 | GPIO17 | SPI0.0 |
| adafruit-15x7-charlieplex | — | — | — | I2C |
| adafruit-4120-charlieplex | — | — | — | I2C |

### Generic ST7789

| Panel | DC | RST | BL | Bus |
|-------|-----|------|-----|------|
| st7789-240x135 | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |
| st7789-240x240 | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |
| st7789-320x240 | GPIO25 | GPIO27 | GPIO24 | SPI0.0 |

### Waveshare Triple Screen

The triple-screen panel uses three SPI displays on separate chip-selects:

| Screen | DC | RST | BL | Bus |
|--------|-----|------|-----|------|
| main (240×240) | GPIO22 | GPIO27 | GPIO19 | SPI1.0 |
| left-aux (160×80) | GPIO4 | GPIO24 | GPIO13 | SPI0.0 |
| right-aux (160×80) | GPIO5 | GPIO23 | GPIO12 | SPI0.1 |

## Input Pin Defaults

Panels with onboard buttons or joysticks have these default input pin mappings:

### waveshare-1.3hat (3 buttons + joystick)

| Input | Pin |
|-------|------|
| Key1 | GPIO5 |
| Key2 | GPIO6 |
| Key3 | GPIO13 |
| Up | GPIO19 |
| Down | GPIO21 |
| Left | GPIO16 |
| Right | GPIO20 |
| Press | GPIO26 |

### waveshare-1.44 (3 buttons + joystick)

| Input | Pin |
|-------|------|
| Key1 | GPIO21 |
| Key2 | GPIO20 |
| Key3 | GPIO16 |
| Up | GPIO6 |
| Down | GPIO19 |
| Left | GPIO5 |
| Right | GPIO26 |
| Press | GPIO13 |

### waveshare-1.3-oled-hat (3 buttons + joystick)

| Input | Pin |
|-------|------|
| Key1 | GPIO21 |
| Key2 | GPIO20 |
| Key3 | GPIO16 |
| Up | GPIO19 |
| Down | GPIO6 |
| Left | GPIO5 |
| Right | GPIO26 |
| Press | GPIO13 |

### waveshare-triple-screen (2 buttons)

| Input | Pin |
|-------|------|
| Key1 | GPIO25 |
| Key2 | GPIO26 |

### adafruit-2.13-ssd1680 (2 buttons)

| Input | Pin |
|-------|------|
| Key1 | GPIO5 |
| Key2 | GPIO6 |

## SPI Bus Pin Usage

SPI buses claim fixed GPIO pins for data transfer in addition to the per-panel DC/RST/BL pins:

| Bus | MOSI | SCLK | CS |
|-----|------|------|----|
| SPI0.0 | GPIO10 | GPIO11 | GPIO8 |
| SPI0.1 | GPIO10 | GPIO11 | GPIO7 |
| SPI1.0 | GPIO20 | GPIO21 | GPIO18 |
| SPI1.1 | GPIO20 | GPIO21 | GPIO17 |

## Connector Conflict Diagnostic

The daemon tracks which GPIO pins are in use by the active panel and reports conflicts with CyberHUD's onboard output connectors. Query this with:

```bash
cyberhudctl gpio pins
```

Example output:

```
OK pin report
  panel=waveshare-2.23-oled-hat-spi controller=ssd1306
  display_mode=monochrome
  display_pins:
    display_bl:  GPIO18
    display_dc:  GPIO24
    display_rst: GPIO25
  connectors:
    3-pin connector GPIO13:  kind=GPIO pins=GND, 5V, GPIO13 status=free
    3-pin connector GPIO18:  kind=GPIO pins=GND, 5V, GPIO18 status=conflict conflict_with=GPIO18->display_bl
    STEMMA QT/Qwiic #1:     kind=I2C1 pins=GND, 3V3, GPIO2/SDA, GPIO3/SCL status=shared (shared I2C bus)
    STEMMA QT/Qwiic #2:     kind=I2C1 pins=GND, 3V3, GPIO2/SDA, GPIO3/SCL status=shared (shared I2C bus)
```

A **conflict** means the panel's display is using a pin that an output connector also needs. In this case, the 3-pin GPIO18 connector can't be used for output because GPIO18 is driving the display backlight.

## Overriding Pin Assignments

If your hardware uses non-standard wiring, override pins in the config file or via CLI flags. Only specify the pins that differ from the panel defaults — omitted fields keep their defaults.

See [Schema Reference](../configuration/schema.md#gpio-pin-fields) for the full list of override fields, or [JSON Config Examples](../configuration/json-config.md#custom-gpio-pin-overrides) for a worked example.
