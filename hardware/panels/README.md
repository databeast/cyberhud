# hardware/panels

Panel product registry and configuration.

- Each product lives in its own package under `hardware/panels/<product>/`.
- Product packages call `panels.Register(...)` in `init()`.
- `hardware/panels/all` imports all built-ins for daemon startup.
- `modes.go`, `resolve.go`, and `pins.go` provide runtime helpers used by `cyberhudd`.

