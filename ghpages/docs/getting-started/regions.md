# Multi-Region Display Model

CyberHUD uses a **region-based display model** that allows each physical screen to be independently controlled. This page explains how displays are organized, how to address them, and how to control modes on individual regions.

## Concepts

### Surfaces

A **surface** is CyberHUD's abstraction for a physical screen. Each surface has a lowercase name (e.g., `main`, `left-aux`, `right-aux`) that identifies it in commands.

### Regions

A **region** is an addressable display unit within the coordinator. Every surface maps to a region identified by a **coordinator index** — a non-negative integer starting at 0. Regions are what you target when switching modes or querying status.

### Region Identifiers

Regions are referenced using `<surface>.<index>` notation:

```
<surface>.<index>
```

- **surface** — a lowercase name matching the pattern `[a-z][a-z0-9-]*`
- **index** — a non-negative integer (typically `0`)

Examples:

| Region ID | Description |
|-----------|-------------|
| `main.0` | Primary display surface |
| `left-aux.0` | Left auxiliary display |
| `right-aux.0` | Right auxiliary display |

## Addressing Modes

### `<surface>.<index>` Notation

The canonical way to address a region uses dot notation:

```bash
cyberhudctl display set 0 clock
cyberhudctl display set left-aux.0 stemma
cyberhudctl display set right-aux.0 gpio
```

### Bare Integer Addressing

As a shorthand, you can reference a region by its **coordinator index** alone — a bare non-negative integer:

| Bare Integer | Equivalent Region |
|--------------|-------------------|
| `0` | Coordinator panel at index 0 (e.g., `main.0`) |
| `1` | Coordinator panel at index 1 (e.g., `left-aux.0`) |
| `2` | Coordinator panel at index 2 (e.g., `right-aux.0`) |

```bash
# These are equivalent (for the waveshare-triple-screen setup):
cyberhudctl display set 0 clock
cyberhudctl display set main.0 clock
```

Use `cyberhudctl display regions` to see which surfaces map to which indices.

## Discovering Available Regions

The `display regions` command lists all configured regions:

```bash
cyberhudctl display regions
```

The output includes each region's:

- **Surface name** — the human-readable identifier (e.g., `main`, `left-aux`)
- **Coordinator index** — the numeric position (e.g., `0`, `1`, `2`)
- **Current mode** — the display mode currently active on that region
- **Available modes** — the set of modes that can be assigned to that region

## Cycling Modes with `display next` / `display prev`

You can cycle through available modes on a region without specifying a mode name:

```bash
# Advance to the next mode on main.0
cyberhudctl display next main.0

# Go back to the previous mode on main.0
cyberhudctl display prev main.0
```

### Wrap-Around Behavior

Cycling wraps around at both ends of the mode list:

- **`display next`** — when the current mode is the **last** in the list, it wraps around to the **first** mode
- **`display prev`** — when the current mode is the **first** in the list, it wraps around to the **last** mode

This means you can repeatedly cycle in either direction without hitting a dead end.

## Error Handling: Non-Existent Regions

If you reference a region that does not exist, the daemon responds with an error listing the available regions:

```
ERR unknown region "foo.0"; available: main.0, left-aux.0, right-aux.0
```

This applies to `display set`, `display next`, `display prev`, and any other command that requires a region identifier. The error message always includes the full list of configured regions so you can correct your command.

## Example: `waveshare-triple-screen` Multi-Panel Setup

The `waveshare-triple-screen` panel registers three surfaces:

| Surface | Coordinator Index | Screen | Resolution |
|---------|-------------------|--------|------------|
| `main` | 0 | Center 1.3" ST7789 | 240×240 |
| `left-aux` | 1 | Left 0.96" ST7735S | 160×80 |
| `right-aux` | 2 | Right 0.96" ST7735S | 160×80 |

### Setting Up Independent Modes

After starting the daemon with the triple-screen panel:

```bash
sudo ./cyberhudd -panel waveshare-triple-screen
```

Set different modes on each region:

```bash
# Main display: interactive menu
cyberhudctl display set 0 menu

# Left auxiliary: STEMMA/I2C sensor readout
cyberhudctl display set left-aux.0 stemma

# Right auxiliary: GPIO pin status
cyberhudctl display set right-aux.0 gpio
```

### Using Bare Integer Addressing

The same commands using bare integer shorthand:

```bash
cyberhudctl display set 0 menu
cyberhudctl display set 1 stemma
cyberhudctl display set 2 gpio
```

### Configuring Policy on Individual Regions

Apply mode-specific policy settings to a single region:

```bash
# Switch main to attract_bokeh with custom speed
cyberhudctl display set 0 attract_bokeh speed=2.0 density=0.8

# Adjust policy on left-aux without switching modes
cyberhudctl display config left-aux.0 interval=5s
```

### Cycling Modes per Region

```bash
# Advance the main display to the next mode
cyberhudctl display next main.0

# Move the right auxiliary display back to the previous mode
cyberhudctl display prev right-aux.0

# Same operations with bare integers
cyberhudctl display next 0
cyberhudctl display prev 2
```

### Querying Region Status

```bash
# List all regions and their current modes
cyberhudctl display regions

# Query the policy active on a specific region
cyberhudctl display policy main.0
```

### Error Example

Referencing a region that doesn't exist:

```bash
cyberhudctl display set back.0 clock
```

Response:

```
ERR unknown region "back.0"; available: main.0, left-aux.0, right-aux.0
```

## Next Steps

- [CLI Reference](../reference/cli.md) — full command documentation including multi-command syntax
- [Waveshare Zero LCD HAT (A)](waveshare-zero-lcd-hat-a.md) — hardware details for the triple-screen panel
- [Policy System](../reference/policy.md) — how to configure mode behavior per region
