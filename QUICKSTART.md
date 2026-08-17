# Quick Start Guide

Welcome! This guide will get you from unboxing to displaying system status in minutes.

## Step 1: Identify Your Display

Look at the display board you're connecting to your Raspberry Pi and find it in the list below:

### Common Waveshare Displays (with buttons)

| Display | Panel | Size | Input |
|---------|-------|------|-------|
| **Waveshare 1.3-inch LCD HAT** | `waveshare-1.3hat` | 240×240 | 3 buttons + joystick |
| Waveshare 1.44-inch LCD HAT | `waveshare-1.44` | 128×128 | 3 buttons + joystick |
| Waveshare 1.3-inch OLED HAT | `waveshare-1.3-oled-hat` | 128×64 | 3 buttons + D-pad |
| Waveshare 2.2-inch SPI LCD | `waveshare-2.2` | 320×240 | None (display only) |

### Adafruit Displays

| Display | Panel | Size | Input |
|---------|-------|------|-------|
| Adafruit 2.13-inch SSD1680 e-ink | `adafruit-2.13-ssd1680-landscape` | 250×122 | 2 buttons |
| Adafruit 2.13-inch e-ink (2-button) | `adafruit-2.13-ssd1680-landscape-2btn` | 250×122 | 2 buttons |

### Generic ST7789 Panels (no buttons)

| Display | Panel | Size |
|---------|-------|------|
| Generic 240×240 | `st7789-240x240` | 240×240 |
| Generic 240×135 | `st7789-240x135` | 240×135 |
| Generic 320×240 | `st7789-320x240` | 320×240 |

### Multi-Display Setups

| Setup | Panel | Details |
|-------|-------|---------|
| Waveshare Zero Triple-Screen | `waveshare-triple-screen` | 1×240×240 (main) + 2×160×80 (aux) |

**Can't find your display?** See [Troubleshooting](#troubleshooting) below.

## Step 2: Install & Enable Required Interfaces

On your Raspberry Pi, enable SPI first. Only enable I2C if you plan to use optional accessory devices:

```sh
sudo raspi-config
```

Navigate to: **Interface Options** → Enable:
- **SPI** (Serial Peripheral Interface)
- **I2C** (optional, only for extra sensors or breakouts)

Press **Finish** and reboot:

```sh
sudo reboot
```

## Step 3: Build (or Install Package)

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

(See [Building](README.md#building) for package build instructions.)

## Step 4: First Run

### Simplest Start (1.3-inch Waveshare HAT)

```sh
sudo ./cyberhudd
```

The daemon will:
1. Show `Booting...` and system status on the display
2. Initialize GPIO, I2C, and display drivers
3. Start an interactive menu when fully booted
4. Listen on `/run/cyberhudd/console.sock` for remote commands

**Expected output:**

```
[cyberhudd] STEMMA QT / QWIIC scanner started (buses: /dev/i2c-1, interval: 2s)
[cyberhudd] GPIO manager started
[cyberhudd] console socket: /run/cyberhudd/console.sock
[cyberhudd] using panel: waveshare-1.3hat (ST7789, 240x240, input enabled)
[cyberhudd] display initialized and ready
```

### Using a Different Display

```sh
sudo ./cyberhudd -panel waveshare-2.2
```

Replace `waveshare-2.2` with your panel name from Step 1.

### Using a JSON Config File

Create `/etc/cyberhudd.json`:

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

## Step 5: Control via CLI or Menu

### What You'll See on Boot

When CyberHUD starts with a display:
1. Display shows boot progress (`Booting...` / `Boot Complete`)
2. Automatically switches to **passive status dashboard** (if no buttons) or **interactive menu** (if buttons present)

#### Dashboard Display (No Buttons)
```
▌CYBERHUD ONLINE

Host: my-rpi4
IP: 192.168.1.100

GPIO: ready
GPIO PINS: 27 available

Optional accessory devices:
  0x76 BME280 Temperature/Pressure...

USAGE:
  cyberhudctl status       [check status]
  cyberhudctl gpio status  [GPIO state]
  cyberhudctl stemma status[I2C devices]
```

The dashboard shows:
- **System status** – Welcoming "CYBERHUD ONLINE" message
- **Network info** – Hostname and IP address (for SSH access)
- **Hardware summary** – GPIO status first, plus any optional accessories you connected
- **Usage hints** – Common CLI commands to try

Perfect for **wall-mounted displays** or **always-on status dashboards**.

#### Interactive Menu (With Buttons/Joystick)
```
         CYBERHUD
    ┌─────────────────┐
    │ Dashboard       │
    │►GPIO Pins       │  ← Cursor here
    │ STEMMA Devices  │
    │ System Info     │
    │ Pin Planner     │
    └─────────────────┘
  [K1]open [K2]▲ [K3]▼
```

Browse modes with joystick/buttons and select with KEY1.

### On-Display Menu (if buttons/joystick present)

- **Joystick/KEY2** – Move cursor up
- **Joystick Down/KEY3** – Move cursor down
- **Joystick Press/KEY1** – Select menu item
- See [Control](README.md#control) for complete key map

### Remote Control (CLI or socket)

**Display status:**

```sh
cyberhudctl display status
```

**Switch to a different mode:**

```sh
cyberhudctl display set main.0 gpio          # Show GPIO status
cyberhudctl display set main.0 stemma        # Show STEMMA QT / QWIIC devices
cyberhudctl display set main.0 system        # Show system info
```

**List available modes for a panel:**

```sh
cyberhudctl display modes
```

**GPIO control:**

```sh
cyberhudctl gpio set 4 1              # Set GPIO4 to HIGH
cyberhudctl gpio in 17                # Read GPIO17 state
```

**GPIO pin mapping:**

```sh
cyberhudctl gpio pins
```

See the full [Control](README.md#control) section for more commands.

## Common Setups

### Setup A: Waveshare 1.3-inch HAT (Interactive)

Fully interactive dashboard with menu buttons. This is the default and recommended starting point.

```sh
sudo ./cyberhudd
```

### Setup B: Waveshare 2.2-inch (Passive Dashboard)

Display-only setup—no buttons—refreshes automatically. Control via CLI only.

```sh
sudo ./cyberhudd -panel waveshare-2.2
```

### Setup C: Adafruit e-ink (Low Power, Slow Updates)

Perfect for always-on status display with minimal power draw.

```sh
sudo ./cyberhudd -panel adafruit-2.13-ssd1680-landscape
```

See [Adafruit 2.13-inch panel README](hardware/panels/adafruit_2_13_ssd1680_landscape/README.md) for GPIO wiring notes.

### Setup D: Multi-Display (Waveshare Triple-Screen)

Three coordinated displays: one main (240×240) and two auxiliary (160×80 each).

```sh
sudo ./cyberhudd -panel waveshare-triple-screen
```

Set mode for each panel independently:

```sh
cyberhudctl display set main.0 menu
cyberhudctl display set left-aux.0 gpio
cyberhudctl display set right-aux.0 stemma
```

### Setup E: Custom GPIO Wiring

If your panel uses different GPIO pins than the default:

```sh
sudo ./cyberhudd \
  -panel waveshare-2.2 \
  -display-dc GPIO25 \
  -display-rst GPIO27 \
  -display-bl GPIO24
```

Use `cyberhudctl help modes` to see all available display modes and their options.

## Enable as Systemd Service

To auto-start on boot:

```sh
sudo systemctl enable cyberhudd.service
sudo systemctl start cyberhudd.service
```

Check status:

```sh
sudo systemctl status cyberhudd.service
```

View logs:

```sh
sudo journalctl -u cyberhudd.service -f
```

### Customize Service (Different Panel)

Create a systemd override:

```sh
sudo systemctl edit cyberhudd.service
```

Add (or modify):

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/cyberhudd -panel waveshare-2.2
```

Save and reload:

```sh
sudo systemctl daemon-reload
sudo systemctl restart cyberhudd.service
```

## Troubleshooting

### "Display failed to initialize"

**Problem:** Daemon logs `panel init failed: ...` and tries fallbacks.

**Solution:**

1. **Verify wiring** – Check GPIO pin assignments match your hardware
2. **Check SPI is enabled** – Run `raspi-config` and confirm SPI is on
3. **List available modes** – `cyberhudctl display modes` to verify setup
4. **Override GPIO pins** – If wiring differs from default, use `-display-dc`, `-display-rst`, `-display-bl`
5. **Check device files** – Verify `/dev/spidev0.0` and `/dev/i2c-1` exist

Example diagnostics:

```sh
# List available display modes
cyberhudctl display modes

# Try with verbose logging
sudo ./cyberhudd -panel waveshare-1.3hat 2>&1 | head -20

# Check for SPI/I2C devices
ls -la /dev/spi* /dev/i2c*
```

### "No input detected"

**Problem:** Buttons/joystick don't work on a panel with input pins.

**Solution:**

1. **Check input pins are wired** – Verify physical connections
2. **Override if needed** – Use `-input-key1`, `-input-key2`, etc.
3. **Fall back to display-only** – Use `waveshare-2.2` or `-noinput` flag

Example with custom key pins:

```sh
sudo ./cyberhudd \
  -panel waveshare-1.3hat \
  -input-key1 GPIO5 \
  -input-key2 GPIO6 \
  -input-key3 GPIO13
```

### "STEMMA QT / QWIIC devices not appearing"

**Problem:** `cyberhudctl stemma status` shows no devices.

**Solution:**

1. **Check I2C is enabled** – Run `raspi-config` → Interface Options → I2C
2. **Verify sensor wiring** – Power, GND, SDA, SCL connections correct
3. **Scan I2C bus manually:**

   ```sh
   sudo i2cdetect -l
   sudo i2cdetect -y 1
   ```

3. **Adjust scan interval** – By default, checks every 2s. Speed up with:

   ```sh
   sudo ./cyberhudd -scan 500ms
   ```

### "Error: permission denied (/dev/gpiomem, /dev/spidev0.0)"

**Problem:** User not in required groups.

**Solution:**

```sh
# Add current user to GPIO, SPI, I2C groups
for grp in gpio spi i2c; do
  sudo usermod -a -G $grp $USER
done

# Apply new group membership
newgrp gpio

# Or reboot to refresh group membership
sudo reboot
```

After package install, the installer automatically adds the `cyberhud` service user to these groups.

### "Display shows garbage / incorrect colors"

**Problem:** Colors are inverted or image appears shifted.

**Solution:**

1. **Adjust MADCTL (rotation/color mode)** – Some panel batches differ. Try:

   ```sh
   sudo ./cyberhudd -display-madctl 0xC8
   ```

2. **Adjust offsets** – Some panels ship with different row/column offsets:

   ```sh
   sudo ./cyberhudd -display-y-offset 3
   ```

   See [Waveshare 1.44-inch panel README](hardware/panels/waveshare_1_44/README.md) for more details.

### "Can't connect to console socket"

**Problem:** `cyberhudctl` returns connection refused.

**Solution:**

1. **Check daemon is running** – `sudo systemctl status cyberhudd.service`
2. **Verify socket exists** – `ls -la /run/cyberhudd/console.sock`
3. **Check socket permissions** – User must be in `cyberhud` group:

   ```sh
   sudo usermod -a -G cyberhud $USER
   newgrp cyberhud
   ```

4. **Specify socket path explicitly** – If custom location:

   ```sh
   cyberhudctl -socket /custom/path/console.sock status
   ```

## Next Steps

- Explore display modes: `cyberhudctl display modes`
- Read [Control](README.md#control) for advanced control
- See [CONFIGURATION.md](CONFIGURATION.md) for GPIO pin configuration
- See [Project Structure](README.md#project-structure) to understand how modes work
- Review panel-specific READMEs under `hardware/panels/<panel>/README.md` for detailed GPIO maps

## Getting Help

- Check logs: `sudo journalctl -u cyberhudd.service -n 100`
- List modes: `cyberhudctl display modes`
- Try headless mode to test daemon without display: `sudo ./cyberhudd -nodisplay`
- Open an issue on GitHub with logs and hardware details

Good luck! 

