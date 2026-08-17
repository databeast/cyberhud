# cyberhud

Adafruit, Waveshare, Elecrow and friends sell cheap SPI/I2C display boards that do nothing on their own — you write the driver, wire up GPIO, and build your own rendering loop before you see a single pixel. CyberHUD is the daemon that skips all of that: plug in the panel, point `cyberhudd` at it, and you get working dashboards, system monitors, and status screens immediately. No Python dependencies, no per-project driver code, no pinout guesswork.

If you've ever bought a display HAT and spent an afternoon fighting Python libraries, broken wikis, and hardcoded GPIO numbers just to show some text on screen — this replaces all of that with a single binary and a one-line config. It isn't a library you build against; it's a running service with a control protocol, so any script or userspace tool can push content to the screen without writing a driver.

## Getting Started

- **[QUICKSTART.md](QUICKSTART.md)** – Identify your display, enable SPI, first run in 15 minutes
- **[CONFIGURATION.md](CONFIGURATION.md)** – Copy-paste JSON configs for common setups, GPIO pin tables

## What It Does

`cyberhudd` is a background daemon that drives SPI/I2C displays on any Linux GPIO device — Raspberry Pi, cyberdeck build, or other embedded board. One screen or three. Interactive with buttons, or passive on a shelf.

- **Instant hardware support** – native drivers for common HATs, generic overrides for any panel on a supported controller chip — no code to write
- **Multi-display** – drive up to 3 panels independently, each with its own mode
- **Display modes** – dashboard, GPIO status, I2C device scanner, system info, clock, serial monitor, USB bench, ticker, image display
- **Runtime control** – switch modes, push text, query state via Unix socket or CLI — from your own scripts or userspace tools, no driver code required
- **Input handling** – buttons and joystick mapped per-panel, or disabled for passive setups
- **I2C scanning** – auto-detects STEMMA QT / QWIIC sensors and breakouts on I2C buses
- **GPIO monitoring** – tracks pin state across the 40-pin header
- **Fallback panels** – if your display doesn't init, it tries alternatives automatically

## Supported Displays

Named panels below get pin/button mappings out of the box. Support is really driven by the controller chip, not the product name — any board using one of these chips works with a generic override, no driver code needed.

| Panel | Controller | Size | Input |
|---|---|---|---|
| `waveshare-1.3hat` | ST7789 | 240×240 | 3 keys + joystick |
| `waveshare-1.44` | ST7735S | 128×128 | 3 keys + joystick |
| `waveshare-2.2` | ST7789 | 320×240 | none (passive) |
| `waveshare-1.3-oled-hat` | SH1106 | 128×64 | 3 keys + D-pad |
| `waveshare-triple-screen` | ST7789 + 2×ST7735S | 240×240 + 2×160×80 | 2 keys |
| `adafruit-2.13-ssd1680-landscape` | SSD1680 | 250×122 e-ink | 2 buttons |
| `st7789-240x240` | ST7789 | 240×240 | none |
| `st7789-240x135` | ST7789 | 240×135 | none |
| `st7789-320x240` | ST7789 | 320×240 | none |

Any panel using ST7789, ST7735S, SH1106, or SSD1680 should work — from Adafruit, Waveshare, Elecrow, or anywhere else — use a generic panel name and override pins/dimensions as needed. See [CONFIGURATION.md](CONFIGURATION.md) for GPIO pin tables and override options.

## Building

```sh
go build -o cyberhudd ./cmd/cyberhudd
go build -o cyberhudctl ./cmd/cyberhudctl
```

Requirements: Go 1.26+, SPI enabled on the Pi, user in `gpio`/`spi`/`i2c` groups.

### Debian package

```sh
make deb VERSION=0.1.0
sudo apt install ./dist/cyberhud_0.1.0-1_arm64.deb
```

The `.deb` installs both binaries, a systemd service, creates a `cyberhud` system user, and starts the daemon. The GitHub Actions workflow at `.github/workflows/deb.yml` cross-builds arm64 packages automatically on tags.

## Running

```sh
sudo ./cyberhudd -panel waveshare-1.3hat
```

On startup, the display shows boot progress until `multi-user.target` is reached, then switches to the dashboard.

### Key flags

| Flag | Default | Description |
|---|---|---|
| `-panel` | `waveshare-1.3hat` | Which display to drive |
| `-config` | | JSON config file (see CONFIGURATION.md) |
| `-nodisplay` | `false` | Headless mode — no LCD |
| `-noinput` | `false` | Passive mode — no buttons |
| `-socket` | `/run/cyberhudd/console.sock` | Control socket path |
| `-i2c` | `/dev/i2c-1` | I2C bus(es) for STEMMA QT / QWIIC scanning |
| `-scan` | `2s` | I2C rescan interval |

All display parameters (width, height, MADCTL, offsets, GPIO pins) are overridable via flags. Run `cyberhudd -help` for the full list, or use a JSON config file.

## Control

### cyberhudctl

```sh
cyberhudctl display status                              # see all panels and their modes
cyberhudctl display modes                               # list available modes per panel
cyberhudctl display set main.0 dashboard                # switch by name
cyberhudctl display set 0 gpio                          # switch by index
cyberhudctl display next main.0                         # cycle to next mode
cyberhudctl display ticker set "line one | line two"    # push text to ticker mode
cyberhudctl stemma status                               # show detected STEMMA QT / QWIIC devices
cyberhudctl gpio status                                 # show GPIO pin state
cyberhudctl gpio pins                                   # pin conflict report
```

Regions are addressed using `<surface>.<index>` notation (e.g., `main.0`, `left-aux.0`) or bare integer coordinator indices (`0`, `1`, `2`). Run `cyberhudctl display regions` to see what's available.

### Raw socket

```sh
nc -U /run/cyberhudd/console.sock
```

The protocol is line-oriented text. Every response starts with `OK` or `ERR`. Send `quit` to disconnect.

## Systemd

The packaged service starts automatically. The socket is group-accessible — add users to the `cyberhud` group for `cyberhudctl` access without root.

For custom flags, use a systemd drop-in:

```sh
sudo systemctl edit cyberhudd
```

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/cyberhudd -panel waveshare-2.2
```

## Project Structure

```
cmd/cyberhudd/             daemon entry point
cmd/cyberhudctl/           CLI client
display/modes/             mode implementations (dashboard, gpio, stemma, clock, etc.)
display/catalog/           mode registry, state tracking, command dispatch
display/surface/           RGBA drawing surface, font rendering
display/widgets/           reusable UI components (icons, borders, progress bars)
hardware/driver/           display controller drivers (ST7789, SSD1680, SH1106, etc.)
hardware/gpio/             GPIO pin manager
hardware/input/            button and joystick input handler
hardware/panels/           panel definitions, override resolution, pin reporting
runtime/console/           Unix socket server (line protocol)
runtime/ui/                UI loop, mode engine, input mapping
```

To add a new display mode: create a package under `display/modes/<name>/`, implement the mode interface, register it in `display/modes/registry.go`. No other wiring needed.
