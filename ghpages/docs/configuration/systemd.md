# Systemd Integration

Run the CyberHUD daemon as a systemd service with a custom configuration file. This page covers creating a systemd override, reloading the service, and verifying it's running correctly.

## Creating a Systemd Override

To use a config file with the systemd service, create a drop-in override:

```sh
sudo systemctl edit cyberhudd.service
```

This opens an editor where you add the override content. Enter the following to point the service at your configuration file:

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/cyberhudd -config /etc/cyberhudd.json
```

!!! note
    The first empty `ExecStart=` line clears the default command. The second line sets the new command with your config file path.

## Reload and Restart

After saving the override, reload the systemd daemon and restart the service:

```sh
sudo systemctl daemon-reload
sudo systemctl restart cyberhudd.service
```

## Verification

Confirm the service is running with your configuration:

```sh
sudo systemctl status cyberhudd.service
```

Follow the live logs to check for errors:

```sh
sudo journalctl -u cyberhudd.service -f
```

## Related Pages

- [JSON Config Examples](json-config.md) - Configuration file examples to use with systemd
- [Schema Reference](schema.md) - Complete configuration schema
- [Pin Assignments](../reference/pin-assignments.md) - GPIO pin override fields
