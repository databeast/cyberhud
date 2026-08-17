# Troubleshooting

This page covers common issues you might encounter when setting up or running CyberHUD, along with solutions and diagnostic steps.

## Display Initialization Failures

If your display doesn't initialize or shows nothing:

1. Verify SPI is enabled on your Raspberry Pi:

    ```sh
    ls /dev/spidev*
    ```

    If no devices appear, enable SPI via `raspi-config` or by adding `dtparam=spi=on` to `/boot/config.txt`.

2. Confirm you're using the correct panel name:

    ```sh
    cyberhudctl display modes
    ```

3. Check the daemon logs for initialization errors:

    ```sh
    sudo journalctl -u cyberhudd.service -n 50
    ```

4. Try running in headless mode to isolate whether the issue is display-specific:

    ```sh
    sudo ./cyberhudd -nodisplay
    ```

## Input Not Working (Buttons/Joystick)

If buttons or joystick input is unresponsive:

1. Verify the correct panel is loaded — different displays use different input pins:

    ```sh
    cyberhudctl display status
    ```

2. Check for GPIO pin conflicts:

    ```sh
    cyberhudctl gpio pins
    ```

3. If using custom GPIO wiring, ensure your [configuration file](../configuration/json-config.md) specifies the correct button/joystick pins (`key1`, `key2`, `key3`, `up`, `down`, `left`, `right`, `press`).

4. You can disable input entirely if not needed:

    ```json
    {
      "display": {
        "disable_input": true
      }
    }
    ```

## Optional Accessory Device Detection

If STEMMA/QWIIC devices aren't being detected:

1. Confirm I2C is enabled:

    ```sh
    ls /dev/i2c*
    ```

2. Check which devices the daemon sees:

    ```sh
    cyberhudctl stemma status
    ```

3. Verify the I2C bus path in your config matches your hardware. Multiple buses can be specified:

    ```json
    {
      "i2c": "/dev/i2c-1,/dev/i2c-3"
    }
    ```

4. Adjust the scan interval if devices appear intermittently:

    ```json
    {
      "scan": "500ms"
    }
    ```

## Permission Errors

If you see permission denied errors:

1. The daemon typically needs root or appropriate group membership to access SPI, I2C, and GPIO:

    ```sh
    sudo ./cyberhudd
    ```

2. For systemd service deployment, ensure the service user has the required group memberships (`spi`, `i2c`, `gpio`). See [Systemd Integration](../configuration/systemd.md) for details.

3. Check socket permissions if `cyberhudctl` can't connect:

    ```sh
    ls -la /run/cyberhudd/console.sock
    ```

## GPIO Conflicts

If you suspect pin conflicts between the display driver and other software:

1. Run the pin conflict checker:

    ```sh
    cyberhudctl gpio pins
    ```

2. Review the [Pin Assignments](../reference/pin-assignments.md) for your display panel to identify which pins are in use.

3. If conflicts exist, use [custom GPIO overrides](../configuration/json-config.md) to remap pins:

    ```json
    {
      "display": {
        "dc": "GPIO25",
        "rst": "GPIO27",
        "bl": "GPIO24"
      }
    }
    ```

## Checking Systemd Logs

When running as a systemd service, logs are available via `journalctl`:

```sh
# View the last 50 log entries
sudo journalctl -u cyberhudd.service -n 50

# Follow logs in real time
sudo journalctl -u cyberhudd.service -f

# Check service status
sudo systemctl status cyberhudd.service
```

## Running Headless Mode

Headless mode starts the daemon without display initialization, which is useful for isolating display-related issues:

```sh
sudo ./cyberhudd -nodisplay
```

In headless mode the daemon will still:

- Scan I2C buses for STEMMA devices
- Monitor GPIO pins
- Serve the console socket

You can verify core functionality works by querying the daemon:

```sh
cyberhudctl gpio status
cyberhudctl stemma status
cyberhudctl gpio pins
```

If headless mode works correctly, the issue is likely display-specific — check your panel selection and GPIO wiring.

## Configuration Troubleshooting

### Config File Not Being Read

Verify the path is correct and the file contains valid JSON:

```sh
# Validate JSON syntax
jq . /etc/cyberhudd.json
```

Ensure you're passing the config flag when starting the daemon:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json
```

### CLI Flags Not Overriding Config

CLI flags should override config file values. Make sure you're using explicit flags alongside the config flag:

```sh
sudo ./cyberhudd -config /etc/cyberhudd.json -panel waveshare-2.2
```

### Partial Overrides Not Working

When using a config file, only include fields you want to override. Omitted fields automatically keep their panel defaults — you don't need to set them to empty or zero:

```json
{
  "display": {
    "panel": "waveshare-1.3hat",
    "dc": "GPIO22"
  }
}
```

This overrides only the DC pin. All other fields (rst, bl, offsets, etc.) remain at the panel's built-in defaults.

## Next Steps

- Review the full [Configuration Schema](../configuration/schema.md) for all available settings
- Set up cyberhud to start at boot: [Systemd Integration](../configuration/systemd.md)
- Browse available [Display Modes](../display-modes/index.md) to customize your setup
- Check the [FAQ](../faq.md) for answers to common questions

## Need More Help?

1. Check the [CLI Usage](cli.md) page for full command reference
2. Review [Configuration](../configuration/json-config.md) for config file setup
3. Read the [Pin Assignments](../reference/pin-assignments.md) for your specific display
4. Try headless mode to isolate display issues: `sudo ./cyberhudd -nodisplay`
5. Open an issue on [GitHub](https://github.com/databeast/cyberhud/issues) with logs and hardware details
