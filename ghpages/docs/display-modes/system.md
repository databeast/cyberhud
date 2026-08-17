# System

System information display showing host details such as hostname, uptime, CPU temperature, memory usage, and network addresses. Provides a quick overview of the host's operational state on any panel.

## Quick Start

```sh
cyberhudctl display set 0 system
```

## How It Works

The system mode gathers host information from the operating system and renders it as a multi-row text display showing hostname, OS architecture, uptime, and IP addresses. The layout arranges rows vertically with automatic font selection to fit the available panel width, and content is centered within the panel area.

The display is static between cache-key evaluations — it redraws when the underlying data changes (a new IP address appears, uptime ticks over) or when the policy changes. For the `cores` and `top` styles, the mode actively samples CPU or process data on each render-cache-key check, making those styles effectively event-driven by the render loop's tick rate.

The data sources are system APIs: `os.Hostname()` for the hostname, `runtime.GOOS/GOARCH` for OS info, `/proc/uptime` for uptime (Linux-only; shows "n/a" elsewhere), and `net.Interfaces()` for IP addresses of active non-loopback interfaces. When no IP addresses are available (no active network interfaces), the mode displays "IP: (none)" as a placeholder.

The `style` option controls layout and data selection:

- **default** — shows hostname, OS/architecture, uptime, and all active IP addresses as separate rows. Balanced overview suitable for most panel sizes.
- **compact** — condensed view for small panels, showing only the most critical metrics in fewer rows with tighter spacing.
- **cores** — focuses on per-core CPU utilization percentages sampled from `/proc/stat`, rendering each core as a row with a usage value. Redraws on each sample.
- **top** — shows top processes ranked by CPU and memory usage in an htop-like list, reading from the process table. Redraws on each sample.

!!! info
    On systems without a CPU temperature sensor (e.g. some VMs or containers), the temperature field may be omitted or show as unavailable.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| style | string | Layout and data density | default | default, compact, cores, top |
| font | string | Font for rendering system text | auto | registered font ID or auto |

Configure options via the CLI:

```sh
cyberhudctl display system style=<value> font=<value>
```

### Style Variants

- **default** — shows hostname, uptime, CPU temp, memory, and primary network address. A balanced overview suitable for most panel sizes.
- **compact** — condensed single-line or two-line summary for small panels. Shows only the most critical metrics (hostname and one or two key values).
- **cores** — focuses on per-core CPU usage and temperatures. Useful for monitoring thermal behavior across all processor cores.
- **top** — shows top processes by CPU and memory usage in an htop-like view. Provides a live process list ranked by resource consumption.

!!! tip
    Use `style=compact` on narrow panels where the full system readout would be truncated, or `style=cores` when diagnosing thermal throttling.

## CLI Examples

Activate the system mode on the main region:

```sh
cyberhudctl display set 0 system
```

Show per-core CPU utilization:

```sh
cyberhudctl display system style=cores
```

Use a specific font for rendering:

```sh
cyberhudctl display system font=ascii-5x7
```

Compact layout with a small font for narrow panels:

```sh
cyberhudctl display system style=compact font=mono-16px
```

Show the top processes view:

```sh
cyberhudctl display system style=top
```

Reset font to automatic selection:

```sh
cyberhudctl display system font=auto
```

Query current system mode settings:

```sh
cyberhudctl display system
```

## Panel Compatibility

The system mode is non-interactive and works on all panels. No input controls are required — it renders passively and updates on a periodic refresh cycle. On monochrome panels, all system metrics render in the native foreground color without color-coded health indicators. The selected style adapts its content density to fit the available display area.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — compact system metrics in white text, periodic refresh |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static system info frame, updates on refresh interval |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale metric text with periodic refresh |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static grayscale system info, updates on refresh interval |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — color-coded health indicators (green/amber/red) with periodic refresh |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static color system info, updates on refresh interval |

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Dashboard](dashboard.md) — consolidated status view that includes system health data
- [Cycle](cycle.md) — auto-cycles through modes including system


