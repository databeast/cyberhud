# Setup

This guide walks you through setting up cyberhud for the first time, from identifying your display to running the daemon and controlling it via CLI.

**Time estimate:** 10–15 minutes

## Step 1: Identify Your Display

Before starting, determine which display panel matches your hardware. See the [Hardware Identification](hardware.md) page for a full list of supported displays and detection commands.

## Step 2: Enable SPI on Your Raspberry Pi

CyberHUD communicates with displays over SPI. Enable it using `raspi-config`:

```sh
sudo raspi-config
```

Navigate to **Interface Options** and enable:

- **SPI** (required for all displays)
- **I2C** (optional — only needed for STEMMA QT / QWIIC accessory sensors)

Reboot after making changes:

```sh
sudo reboot
```

## Step 3: Build and Install

### From Source

```sh
cd /path/to/cyberhud
go build -o cyberhudd ./cmd/cyberhudd
go build -o cyberhudctl ./cmd/cyberhudctl
```

### From Debian Package

```sh
sudo apt install ./dist/cyberhud_0.1.0-1_arm64.deb
```

## Step 4: First Run

Start the daemon with the default panel (Waveshare 1.3-inch HAT):

```sh
sudo ./cyberhudd
```

The daemon will:

1. Show `Booting...` and system status on the display
2. Initialize GPIO, I2C, and display drivers
3. Start an interactive menu (if buttons are present) or a passive dashboard
4. Listen on `/run/cyberhudd/console.sock` for remote commands

**Expected output:**

```
[cyberhudd] STEMMA QT / QWIIC scanner started (buses: /dev/i2c-1, interval: 2s)
[cyberhudd] GPIO manager started
[cyberhudd] console socket: /run/cyberhudd/console.sock
[cyberhudd] using panel: waveshare-1.3hat (ST7789, 240x240, input enabled)
[cyberhudd] display initialized and ready
```

!!! success "What You'll See"

    On the physical display you will see `Booting...` followed by system status text as
    drivers initialize. Once boot completes, the display enters **menu mode** (interactive
    navigation via buttons) if your panel has input controls, or **dashboard mode** (a
    passive system-status view) if no buttons are present. In the terminal, the line
    `console socket: /run/cyberhudd/console.sock` confirms that the daemon is ready to
    accept commands from `cyberhudctl`.

### Using a Different Display

Specify your panel from [Hardware Identification](hardware.md):

```sh
sudo ./cyberhudd -panel waveshare-2.2
```

### Using a JSON Config File

Create `/etc/cyberhudd.json` (see [JSON Config Examples](../configuration/json-config.md) for details):

```json
{
  "socket": "/run/cyberhudd/console.sock",
  "i2c": "/dev/i2c-1",
  "scan": "2s",
  "display": {
    "panel": "waveshare-1.3hat"
  }
}
```

Then run:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json
```

## Daemon Control

### Start the daemon manually

```sh
sudo ./cyberhudd
```

### Start with a specific panel

```sh
sudo ./cyberhudd -panel waveshare-2.2
```

### Use a JSON config file

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json
```

### Run in headless mode (no display)

```sh
sudo ./cyberhudd -nodisplay
```

### View daemon logs

```sh
sudo journalctl -u cyberhudd.service -f
```

## Next Steps

- Learn to control cyberhud from the terminal: [CLI Usage](cli.md)
- Identify your hardware and supported panels: [Hardware Identification](hardware.md)
- Configure cyberhud with a JSON file: [JSON Config Examples](../configuration/json-config.md)
- Run into issues? Check [Troubleshooting](troubleshooting.md)
