# Policy System

The policy system is a per-mode key=value configuration mechanism. Each display mode registers its own set of named fields with validation rules, settable at runtime via the console protocol and persistable to the JSON config file.

## Overview

Every display mode can define **policy fields** — named parameters that control rendering behavior at runtime. Policy fields are:

- Declared at mode registration time with type, default value, and validation rules
- Queried and set via the console protocol (through `cyberhudctl`)
- Persisted to the JSON config file with `freeze policy`
- Restored automatically on daemon startup

---

## Mode Registration

Modes register policy fields through the catalog system during `init()`. Each mode provides:

1. **A catalog `Definition`** — declares the mode's available options (key, type, default, allowed values, summary)
2. **A `CommandDefinition`** — registers a mode-specific command verb with a handler
3. **A `PolicySnapshotter`** — handles serialization and restoration of policy state

### Catalog Definition Options

When a mode calls `catalog.Register()`, it includes an `Options` slice describing each configurable field:

| Field | Description |
|-------|-------------|
| `Key` | The policy key name (e.g., `speed`, `density`, `style`) |
| `Type` | Value type: `float`, `int`, `string`, or `bool` |
| `Summary` | One-sentence description of what the field controls |
| `Default` | Default value as a string |
| `Allowed` | List of allowed values (for enum-style fields) |

### Command Handler Pattern

Each mode with policy fields registers a `CmdHandler` (from `display/catalog/cmdutil`) that defines:

- **Keys** — recognized key names with validation functions
- **Get** — returns the current value for a key
- **Apply** — sets a validated value for a key
- **PostApply** (optional) — hook for additional actions after policy changes (e.g., fitness notes)

### Validation Rules

The `cmdutil` package provides standard validators:

| Validator | Behavior | Error Example |
|-----------|----------|---------------|
| Float range | Accepts floats in `[min, max]` | `must be in [0.1, 10.0]` |
| `AllowedValidator` | Accepts values from a fixed list (case-insensitive) | `must be one of [bold, compact, dense]` |
| `BoolValidator` | Accepts `true/false/yes/no/on/off/1/0` | `must be true or false` |
| `IntValidator` | Accepts integers >= minimum | `must be >= 500` |

Modes can also define custom validators for special cases (e.g., timezone validation, endpoint URL validation).

---

## Querying and Setting Policy

### `display config` — Query or Set on Active Mode

Query or update policy fields on the currently active mode for a region without switching modes:

```bash
cyberhudctl display config <region> [key=value ...]
```

- With **no key=value pairs**: returns current policy state
- With **key=value pairs**: validates all pairs, then applies atomically

**Examples:**

```bash
# Query current policy for the active mode on main.0
cyberhudctl display config main.0

# Set a single field
cyberhudctl display config main.0 speed=1.5

# Set multiple fields atomically
cyberhudctl display config main.0 speed=2.0 density=0.8
```

### `display set` — Switch Mode with Inline Policy

Switch a region to a new mode, optionally applying policy fields in a single operation:

```bash
cyberhudctl display set <region> <mode> [key=value ...]
```

**Examples:**

```bash
# Switch to attract_bokeh with default policy
cyberhudctl display set 0 attract_bokeh

# Switch to attract_bokeh with inline policy overrides
cyberhudctl display set 0 attract_bokeh speed=2.0 density=0.8
```

### `display <mode>` — Mode-Specific Command

Each mode registers a command verb matching its mode ID. This provides a direct interface for querying and setting that mode's policy:

```bash
cyberhudctl display <mode> [key=value ...]
```

**Examples:**

```bash
# Query attract_bokeh policy
cyberhudctl display attract_bokeh

# Set attract_bokeh policy fields
cyberhudctl display attract_bokeh speed=1.5 density=0.8

# Query clock policy
cyberhudctl display clock

# Set clock policy fields
cyberhudctl display clock style=digital show_seconds=true time_format=24h
```

---

## Query Response Format

When a mode command is invoked with **zero arguments**, it returns the current policy state:

```
OK <mode> key1=value1 key2=value2 ...
```

Keys appear in the order they were declared in the mode's `CmdHandler.Keys` slice.

**Example responses:**

```
OK attract_bokeh speed=1 density=0.5 size_variance=0.5 saturation=0.7
OK clock style=digital show_seconds=false time_format=12h date_format=short timezone=Local show_weekday=true blink_colon=true fgcolor=white show_led=false seconds_bar=none show_daybar=false show_border=false border_color=white
```

After successfully applying key=value pairs, the handler also returns this format — showing the updated state.

---

## Error Responses

When a policy field receives an invalid value, the system responds with:

```
ERR key: reason
```

- **key** — the field name that failed validation
- **reason** — a human-readable description of why the value was rejected

All key=value pairs are validated **before** any are applied. If any pair fails validation, none are applied — the operation is atomic.

### Error Examples

```
ERR speed: must be in [0.1, 10.0]
ERR style: must be one of [bold, compact, dense]
ERR show_seconds: must be true or false
ERR refresh_ms: must be >= 500
ERR unknown key foo
```

An unrecognized key produces:

```
ERR unknown key <name>
```

A bare token (no `=` separator) produces:

```
ERR unknown key <token>
```

---

## Persistence

### Snapshotter Interface

Every registered mode implements the `PolicySnapshotter` interface:

```go
type PolicySnapshotter interface {
    SnapshotPolicy() map[string]interface{}
    RestorePolicy(data map[string]interface{}) error
}
```

- **`SnapshotPolicy()`** — returns the current policy as a JSON-serializable map with snake_case keys
- **`RestorePolicy(data)`** — applies policy values from a JSON map, applying normalization (clamping out-of-range values)

### Persisting Policy

The `freeze policy` command tells the daemon to snapshot all mode policies and write them to the `policies` section of the JSON config file:

```bash
cyberhudctl freeze policy
```

On next startup, the daemon reads the `policies` map and calls `RestorePolicy()` on each mode's snapshotter.

### Config File Structure

The `policies` field in the JSON config stores per-mode policy snapshots:

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
      "show_seconds": true,
      "time_format": "24h"
    }
  }
}
```

### Policy Change Notification

When any `CmdHandler` successfully applies key=value pairs, the system:

1. Snapshots the mode's current policy via `SnapshotPolicy()`
2. Serializes it to JSON
3. Forwards the snapshot to the daemon's `PolicyStore` via the `OnPolicyChange` callback

This keeps the in-memory policy store synchronized with runtime changes.

---

## Stub Snapshotters

Modes with **no configurable policy fields** use a `StubSnapshotter`:

```go
type StubSnapshotter struct{}

func (StubSnapshotter) SnapshotPolicy() map[string]interface{}     { return map[string]interface{}{} }
func (StubSnapshotter) RestorePolicy(map[string]interface{}) error { return nil }
```

A stub snapshotter:

- Returns an empty map on snapshot (no fields to persist)
- Accepts any restore call as a no-op (nothing to apply)
- Allows the mode to participate in the `freeze policy` mechanism without special-casing

The test/utility modes (`testpattern`, `testfonts`, `testicons`, `testwidgets`) use stub snapshotters since they have no user-configurable parameters. Modes with configurable fields may also use stub snapshotters when their configuration is transient or not meaningful to persist.

---

## Complete Example

This walkthrough demonstrates the full policy lifecycle for `attract_bokeh`:

```bash
# 1. Query current policy (zero-argument invocation)
$ cyberhudctl display attract_bokeh
OK attract_bokeh speed=1 density=0.5 size_variance=0.5 saturation=0.7

# 2. Set speed and density
$ cyberhudctl display attract_bokeh speed=2.5 density=0.8
OK attract_bokeh speed=2.5 density=0.8 size_variance=0.5 saturation=0.7

# 3. Attempt an invalid value
$ cyberhudctl display attract_bokeh speed=99
ERR speed: must be in [0.1, 10.0]

# 4. Use display config on the active mode for a region
$ cyberhudctl display config main.0 speed=1.5
OK attract_bokeh speed=1.5 density=0.8 size_variance=0.5 saturation=0.7

# 5. Switch mode with inline policy
$ cyberhudctl display set 0 attract_geometric speed=3.0 density=0.6
OK attract_geometric speed=3 density=0.6 glow_intensity=1 fragment_rate=1

# 6. Persist all policies to disk
$ cyberhudctl freeze policy
OK

# 7. On next daemon restart, policies are restored from config file
```

---

## Mode Policy Reference

Use `help modes` to query all registered mode commands and their available policy options at runtime:

```bash
cyberhudctl help modes
```

This returns each mode's command verb, summary, usage pattern, and the list of options with types, defaults, and descriptions. See individual [display mode pages](../display-modes/index.md) for per-mode policy field tables.
