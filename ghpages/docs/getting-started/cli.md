# CLI Usage

!!! prerequisite "Daemon must be running"
    The commands on this page require the `cyberhudd` daemon to be running. If you haven't started it yet, complete the [Setup](setup.md) guide first.

CyberHUD provides `cyberhudctl` as the primary command-line interface for interacting with the daemon. You can use it to query hardware, control the display, read sensors, and manage the daemon process.

## Identify Hardware

Discover available display modes and check current display status:

```sh
# List all available display modes
cyberhudctl display modes

# Check current display status
cyberhudctl display status
```

## Control the Display

Check display status, switch between display modes, and see what modes are available:

```sh
# See status of all panels
cyberhudctl display status

# Switch to a different display mode
cyberhudctl display set 0 gpio
cyberhudctl display set 0 stemma
cyberhudctl display set 0 system

# See what modes are available
cyberhudctl display modes
```

For comprehensive documentation on each display mode, see the [Display Modes](../display-modes/index.md) section.

## Read Sensors and GPIO

Query I2C devices, inspect GPIO pin state, and check for pin conflicts:

```sh
# List detected I2C devices (STEMMA/QWIIC)
cyberhudctl stemma status

# List GPIO pins and their state
cyberhudctl gpio status

# Read a single GPIO pin
cyberhudctl gpio in 4

# Set a GPIO pin output
cyberhudctl gpio set 4 1

# Check for pin conflicts
cyberhudctl gpio pins
```

## Daemon Control

Start, configure, and monitor the cyberhud daemon:

```sh
# Start the daemon manually
sudo ./cyberhudd

# Start with a specific panel
sudo ./cyberhudd -panel waveshare-2.2

# Use a JSON config file
sudo ./cyberhudd -config /etc/cyberhudd.json

# Run in headless mode (no display output)
sudo ./cyberhudd -nodisplay

# View daemon logs (follow mode)
sudo journalctl -u cyberhudd.service -f
```

For details on JSON configuration files, see the [Configuration](../configuration/json-config.md) section.

## Programmatic Control

Beyond interactive CLI usage, you can control the daemon programmatically:

- **cyberhudctl CLI** — Script commands in shell scripts or automation pipelines
- **Unix socket** — Connect directly to the daemon socket at `/run/cyberhudd/console.sock` for programmatic communication
- **Display mode switching** — Automate transitions between GPIO, STEMMA, system, and other display modes

You can persist mode settings across daemon restarts using `cyberhudctl freeze policy` — see [Configuration](../configuration/json-config.md) for details on the policy system.

```sh
# Example: query status from a script
cyberhudctl display status
cyberhudctl stemma status
cyberhudctl gpio status
```

## Command Reference Summary

| Command | Description |
|---------|-------------|
| `cyberhudctl display status` | Show status of all panels |
| `cyberhudctl display modes` | List available display modes |
| `cyberhudctl display set <region> <mode>` | Switch display mode |
| `cyberhudctl display next <region>` | Cycle to next mode |
| `cyberhudctl display prev <region>` | Cycle to previous mode |
| `cyberhudctl stemma status` | List detected I2C devices |
| `cyberhudctl gpio status` | List GPIO pins and state |
| `cyberhudctl gpio in <pin>` | Read a GPIO pin |
| `cyberhudctl gpio set <pin> <value>` | Set a GPIO pin output |
| `cyberhudctl gpio pins` | Check for pin conflicts |

## Next Steps

- [Display Modes](../display-modes/index.md) — Comprehensive guide to all display modes
- [Configuration](../configuration/json-config.md) — Set defaults with a JSON config file
- [Pin Assignments](../reference/pin-assignments.md) — Pin assignments for supported displays
- [Troubleshooting](troubleshooting.md) — Common issues and solutions
