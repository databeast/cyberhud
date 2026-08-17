# CLI Reference

`cyberhudctl` is the command-line client for controlling the CyberHUD daemon (`cyberhudd`). It communicates over a Unix socket using a line-oriented text protocol.

## Connection

```bash
cyberhudctl [flags] <command>
```

| Flag | Default | Description |
|------|---------|-------------|
| `-socket` | `/run/cyberhudd/console.sock` | Path to the cyberhudd Unix socket |
| `-timeout` | `2s` | Socket read/write timeout |

## Region Addressing

CyberHUD organizes displays into **surfaces**, each corresponding to a physical screen. Every surface has a name and a numeric index. Commands that operate on a specific display accept a **region identifier** to target a surface.

### `<surface>.<index>` Notation

The canonical form for addressing a region uses dot-separated notation:

```
<surface>.<index>
```

- **surface** — a lowercase name matching the pattern `[a-z][a-z0-9-]*` (e.g., `main`, `left-aux`, `right-aux`)
- **index** — a non-negative integer (e.g., `0`, `1`, `2`)

Examples:

| Region ID | Meaning |
|-----------|---------|
| `main.0` | Primary display, index 0 |
| `left-aux.0` | Left auxiliary display, index 0 |
| `right-aux.0` | Right auxiliary display, index 0 |

### Bare Integer Addressing

As a shorthand, you can reference a region by its **coordinator index** alone — a bare non-negative integer:

| Region ID | Equivalent |
|-----------|------------|
| `0` | Coordinator panel at index 0 |
| `1` | Coordinator panel at index 1 |
| `2` | Coordinator panel at index 2 |

Bare integer addressing resolves to the coordinator panel at that index. Use `cyberhudctl display regions` to see which surfaces map to which indices.

### Discovering Available Regions

```bash
cyberhudctl display regions
```

This lists all configured regions with their surface names and coordinator indices. If you reference a non-existent region, the daemon responds with an error indicating the region is not configured and lists available regions.

---

## System Commands

### `status`

Query the overall daemon status.

```bash
cyberhudctl status
```

Returns an `OK` response with system state information.

---

## GPIO Commands

These commands interact with the GPIO subsystem for reading and writing individual pins.

### `gpio status`

Query the overall GPIO subsystem status.

```bash
cyberhudctl gpio status
```

Returns status information about the GPIO hardware state.

### `gpio pins`

List all configured GPIO pins.

```bash
cyberhudctl gpio pins
```

Returns the set of GPIO pins known to the daemon, including their current assignments.

### `gpio set <pin> <0|1>`

Set a GPIO output pin to a specified level.

```bash
cyberhudctl gpio set <pin> <0|1>
```

- **pin** — integer pin number
- **0|1** — output level (0 = low, 1 = high)

**Example:**

```bash
# Set pin 25 high
cyberhudctl gpio set 25 1

# Set pin 25 low
cyberhudctl gpio set 25 0
```

### `gpio in <pin>`

Read the current input level of a GPIO pin.

```bash
cyberhudctl gpio in <pin>
```

- **pin** — integer pin number

**Example:**

```bash
# Read the level of pin 17
cyberhudctl gpio in 17
```

---

## STEMMA Commands

These commands query the STEMMA QT / QWIIC I2C sensor subsystem.

### `stemma status`

Query the STEMMA subsystem status.

```bash
cyberhudctl stemma status
```

Returns information about detected STEMMA QT / QWIIC devices on the configured I2C buses.

---

## Display Region Commands

These commands query information about the display system without modifying state.

### `display regions`

List all configured display regions.

```bash
cyberhudctl display regions
```

Returns each region's surface name, coordinator index, current mode, and available modes.

### `display status`

Query the current display state. (`display list` is accepted as a synonym.)

```bash
cyberhudctl display status
```

Returns status information for all active display regions.

### `display modes`

List all available display modes.

```bash
cyberhudctl display modes
```

Returns the set of display modes registered in the mode catalog that can be assigned to regions.


---

## Display Control Commands

These commands modify display state — switching modes, updating policy, and cycling through modes on a specific region.

### `display set <region> <mode> [key=value ...]`

Switch the specified region to a new display mode. Optionally pass inline policy key=value pairs to configure the mode at switch time.

```bash
cyberhudctl display set <region> <mode> [key=value ...]
```

- **region** — target region using `<surface>.<index>` notation or a bare integer
- **mode** — the mode ID to switch to (e.g., `clock`, `attract_plasma`)
- **key=value** — zero or more policy parameters applied to the new mode

The mode must not be empty. If inline policy keys are provided, they are applied immediately after the mode switch.

**Examples:**

```bash
# Switch main.0 to the clock mode
cyberhudctl display set 0 clock

# Switch to attract_bokeh with inline policy
cyberhudctl display set 0 attract_bokeh speed=2.0 density=0.8

# Using bare integer addressing
cyberhudctl display set 0 thermal
```

### `display config <region> [key=value ...]`

Update policy fields on the currently active mode for a region without switching modes. When called with no key=value pairs, queries the current policy state.

```bash
cyberhudctl display config <region> [key=value ...]
```

- **region** — target region using `<surface>.<index>` notation or a bare integer
- **key=value** — zero or more policy parameters to set on the active mode

**Examples:**

```bash
# Update policy on the currently active mode for main.0
cyberhudctl display config main.0 speed=1.5

# Set multiple policy fields at once
cyberhudctl display config left-aux.0 brightness=0.8 interval=5s

# Query current policy (no key=value pairs)
cyberhudctl display config main.0
```

### `display next <region>`

Cycle to the next available mode on the specified region. Wraps around from the last mode to the first.

```bash
cyberhudctl display next <region>
```

- **region** — target region using `<surface>.<index>` notation or a bare integer

**Example:**

```bash
# Advance main.0 to the next mode in the cycle
cyberhudctl display next main.0

# Using bare integer addressing
cyberhudctl display next 0
```

### `display prev <region>`

Cycle to the previous available mode on the specified region. Wraps around from the first mode to the last.

```bash
cyberhudctl display prev <region>
```

- **region** — target region using `<surface>.<index>` notation or a bare integer

**Example:**

```bash
# Move main.0 to the previous mode in the cycle
cyberhudctl display prev main.0

# Using bare integer addressing
cyberhudctl display prev 0
```

---

## Policy Commands

These commands query and inspect mode policy state.

### `display policy <mode|region>`

Query the current policy for a specific mode or region. When given a mode ID, returns that mode's policy fields. When given a region identifier, returns the policy for the mode currently active on that region.

```bash
cyberhudctl display policy <mode|region>
```

- **mode|region** — a mode ID (e.g., `attract_bokeh`) or a region identifier (e.g., `main.0`, `0`)

**Examples:**

```bash
# Query policy for the attract_bokeh mode
cyberhudctl display policy attract_bokeh

# Query policy for the mode active on main.0
cyberhudctl display policy main.0
```

### `policy dump`

Dump all mode policies currently held by the daemon. Returns the full set of per-mode policy snapshots.

```bash
cyberhudctl policy dump
```

**Example:**

```bash
# View all mode policies
cyberhudctl policy dump
```

---

## Advanced Commands

### `help modes`

Query mode command metadata from the daemon. Returns structured information about each registered mode command, including its verb, usage pattern, summary, scope, and available policy options with types and defaults.

```bash
cyberhudctl help modes
```

The output is formatted as a human-readable list of mode commands with their options:

```
Mode Commands:
- ticker (Ticker, scope=global)
  Scrolling text ticker display
  usage: display ticker set <text>
  options:
    - speed (float, default=1.0): Scroll speed multiplier
    - font_size (int, default=16): Font size in pixels
```

**Example:**

```bash
# List all mode commands and their options
cyberhudctl help modes
```

### `raw <line...>`

Send a raw protocol line directly to the daemon without any client-side parsing or validation. Intended for debugging and development.

```bash
cyberhudctl raw <line...>
```

- **line** — one or more arguments joined with spaces and sent verbatim to the daemon

**Example:**

```bash
# Send an arbitrary protocol line
cyberhudctl raw display set main.0 clock
```

### Mode-Specific Commands: `display <mode> [args...]`

Any unrecognized `display` subcommand is treated as a mode-specific command. The mode ID and all subsequent arguments are forwarded to the daemon as a single protocol line:

```bash
cyberhudctl display <mode> [args...]
```

This enables mode-registered commands. For example, the ticker mode registers a `set` subcommand and the image mode registers `set` and `clear`:

```bash
# Set ticker text
cyberhudctl display ticker set "Hello World"

# Set an image
cyberhudctl display image set /path/to/image.png

# Clear the image
cyberhudctl display image clear
```

Mode commands are documented on each mode's individual page. Use `cyberhudctl help modes` to discover available mode commands and their options.


---

## Persistence Commands

These commands persist runtime state to disk. Both are executed daemon-side — the CLI sends the command and the daemon writes to its configuration file.

### `freeze`

Persist the current hardware configuration to the config file.

```bash
cyberhudctl freeze
```

Saves the active hardware settings (panel assignments, GPIO mappings, display parameters) so they survive a daemon restart.

### `freeze policy`

Persist all mode policies to the config file.

```bash
cyberhudctl freeze policy
```

Saves every mode's current policy state (key=value fields) to the `policies` section of the JSON config. On next startup, the daemon restores these saved policies automatically.

---

## Multi-Command Syntax

`cyberhudctl` supports sending multiple commands in a single invocation by separating them with standalone semicolons (` ; `). Commands are processed sequentially over a single socket connection.

### Syntax

```bash
cyberhudctl [flags] <cmd1> ; <cmd2> ; <cmd3> ...
```

Each `;` must be a standalone argument — shell quoting or escaping is typically needed:

```bash
# Using shell quoting to prevent semicolons from being interpreted by the shell
cyberhudctl region main.0 ';' mode attract_bokeh ';' config speed=2.0
```

### Scoped Commands

Within a multi-command invocation, the following **scoped commands** are available:

| Scoped Command | Expands To | Description |
|----------------|------------|-------------|
| `region <id>` | *(sets context)* | Set the active region for subsequent scoped commands |
| `mode <mode> [key=value ...]` | `display set <region> <mode> [key=value ...]` | Switch mode on the active region |
| `config [key=value ...]` | `display config <region> [key=value ...]` | Update policy on the active region |
| `next` | `display next <region>` | Cycle to next mode on the active region |
| `prev` | `display prev <region>` | Cycle to previous mode on the active region |
| `status` | `display policy <region>` | Query policy for the active region |

### Region Context Requirement

A `region <id>` command **must precede** any scoped command. Without an active region context, scoped commands (`mode`, `config`, `next`, `prev`, `status`) produce an error:

```
ERR no region context active; use 'region <id>' first
```

The `region` command itself does not generate a protocol command — it only sets the context for subsequent scoped commands in the same invocation.

### Examples

```bash
# Set region, switch mode, and configure policy in one invocation
cyberhudctl region main.0 ';' mode attract_plasma ';' config speed=1.5 density=0.8

# Cycle to next mode on a region
cyberhudctl region left-aux.0 ';' next

# Use bare integer addressing with region context
cyberhudctl region 0 ';' mode clock ';' status

# Mix scoped and non-scoped commands
cyberhudctl region main.0 ';' mode ticker ';' display ticker set "Hello World"
```

### Non-Scoped Commands in Multi-Command Mode

Commands that are not scoped (e.g., `status`, `freeze`, `gpio status`) can also appear in a multi-command invocation. They are routed through the standard command parser and do not require a region context:

```bash
# Query system status, then set a mode on a region
cyberhudctl status ';' region main.0 ';' mode clock
```

### Execution Behavior

- Commands execute sequentially over a single socket connection.
- If any command returns an error response, execution stops immediately and the error is reported.
- The `region` command sets context but produces no protocol traffic.
- All successful responses are printed in order after execution completes.
