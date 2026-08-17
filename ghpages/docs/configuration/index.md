# Configuration

This section covers all configuration options for the CyberHUD daemon. Whether you need to set up a basic configuration file, understand the full schema, map GPIO pins, or run the daemon as a systemd service, you'll find the details here.

## Pages

- [JSON Config Examples](json-config.md) - Example configurations for common setups
- [Schema Reference](schema.md) - Complete schema reference for all config fields
- [Pin Assignments](../reference/pin-assignments.md) - GPIO pin assignments and mappings
- [Systemd Integration](systemd.md) - Running the CyberHUD daemon as a systemd service

## Policy System

The policy system lets you persist runtime mode settings so they survive daemon restarts. Each display mode can define policy fields (speed, style, density, etc.) that you tune via `cyberhudctl` at runtime. The `freeze policy` command writes those settings to the config file, and the daemon restores them automatically on startup.

### The Freeze → Persist → Load Cycle

1. **Set options at runtime** — Use `cyberhudctl` to adjust mode settings:

    ```sh
    cyberhudctl display attract_bokeh speed=2.0 density=0.8
    ```

2. **Freeze to disk** — Run `cyberhudctl freeze policy` to snapshot all current mode policies and write them to the `"policies"` map in your JSON config file:

    ```sh
    cyberhudctl freeze policy
    ```

    This produces a `"policies"` key in the config file containing per-mode settings:

    ```json
    {
      "policies": {
        "attract_bokeh": {
          "speed": 2.0,
          "density": 0.8,
          "size_variance": 0.5,
          "saturation": 0.7
        },
        "clock": {
          "style": "digital",
          "show_seconds": true
        }
      }
    }
    ```

3. **Automatic restore on startup** — When the daemon starts, it reads the `"policies"` key from `fileConfig` and calls each mode's `RestorePolicy()` method to apply the saved settings. No manual intervention is needed after a reboot.

### Key Points

- Policies are stored per-mode (keyed by mode ID) and apply globally, not per-region.
- The freeze command snapshots **all** registered modes at once — you cannot selectively freeze a single mode.
- Modes that define no policy fields use a stub snapshotter and produce empty entries (or are omitted).
- Values are validated on restore using the same rules as runtime set commands. Out-of-range values are clamped to valid bounds.

### Further Reading

For the full reference — including mode registration, validation rules, query/response formats, error handling, and the snapshotter interface — see the [Policy System Reference](../reference/policy.md).
