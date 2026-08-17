# Frequently Asked Questions

## How do I switch display modes?

Use the `cyberhudctl display set` command followed by the region address and the mode ID. This immediately activates the chosen mode on the specified display panel.

```sh
cyberhudctl display set 0 clock
```

See [CLI Usage](getting-started/cli.md) for the full list of available commands and regions.

## How do I configure mode options?

You can set mode options from the CLI using the `cyberhudctl display <mode-id> key=value` syntax. For persistent configuration, add a `"policies"` key to your JSON config file mapping mode IDs to their option objects.

```sh
cyberhudctl display clock style=big-digit time_format=12h
```

See [JSON Configuration](configuration/json-config.md) for details on persisting settings across restarts.

## Which panels support which modes?

Every mode documents its panel compatibility using six capability classes: MonoFast, MonoSlow, GrayscaleFast, GrayscaleSlow, ColorFast, and ColorSlow. These classes describe the display's color depth and refresh speed, which determine how a mode renders — color panels get full output, mono panels receive dithered or simplified visuals, and slow-refresh panels disable animation.

See the [Display Modes Index](display-modes/index.md) for a complete listing of all modes and their panel support.

## How do I update CyberHUD?

For source builds, pull the latest code and rebuild:

```sh
cd /path/to/cyberhud && go build ./cmd/...
```

For Debian packages, install the updated `.deb` from the releases:

```sh
sudo apt install ./dist/cyberhud_*.deb
```

## How do I report bugs or request features?

Open an issue on the GitHub repository at [github.com/databeast/cyberhud/issues](https://github.com/databeast/cyberhud/issues). Include your panel model, CyberHUD version, and steps to reproduce for bug reports.

## How do I run CyberHUD at boot?

CyberHUD ships with a systemd unit file that starts the daemon automatically on boot. See [Systemd Integration](configuration/systemd.md) for setup instructions and service management commands.

## How do I use multiple displays?

CyberHUD supports multiple panels through its region addressing model, where each physical display is assigned a named region (e.g., `main.0`, `left-aux.0`). You configure regions in your JSON config and address them independently via the CLI.

See [Regions](getting-started/regions.md) for the full multi-display setup guide.
