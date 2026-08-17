# hardware/driver

Hardware display driver registry and text/layout helpers.

- Driver packages live under `hardware/driver/<chipset>/`.
- Each driver package registers itself in `init()`.
- `hardware/driver/registry.go` holds the generic registry and config types.
- `hardware/driver` also owns renderer, text-hint, and virtual-display primitives.

