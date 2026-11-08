# Schedule

Shows the calendar event that is active right now, read from an iCalendar (`.ics`) source that can be a local file or a remote `http(s)` URL. When nothing is active, the mode previews the next upcoming event and how soon it starts. The display is a single "what's on now" card that adapts to the panel's capability and size.

## Quick Start

```sh
cyberhudctl display set 0 schedule
```

Point the mode at a calendar — a local file:

```sh
cyberhudctl display schedule source=/home/pi/calendar.ics
```

…or a remote URL (for example a published Google/Nextcloud calendar):

```sh
cyberhudctl display schedule source=https://example.com/calendar.ics
```

The mode fetches the calendar in the background and refreshes it every five minutes by default.

## How It Works

Schedule reads an iCalendar source and renders the single event that spans the current moment. A background goroutine owns the fetch lifecycle: on activation it loads the calendar once, then re-fetches on the interval set by `refresh_sec`.

- **Source resolution** — if `source` begins with `http://` or `https://` the calendar is fetched over HTTP; otherwise it is read from the local filesystem path. Changing `source` drops the cached calendar and triggers an immediate re-fetch.
- **Active-event selection** — the parsed events are kept sorted by start time, and the event active *right now* is chosen on every render, not just on every fetch. So the card flips to the next event the moment the current one ends, even between fetches. When events overlap, the one that started most recently wins.
- **Idle preview** — when no event is active, the card shows the next upcoming event and a countdown ("in 20m") to its start.
- **Status states** — an unset source, a calendar still loading, and a fetch or parse error each state themselves plainly in a restrained, unaccented treatment, so a problem never looks like a real event.

On color panels an active event is accented (amber) to read as "live now"; idle and error states stay white so the accent always means something is on.

!!! note
    Recurring events (`RRULE`) are not yet expanded — only the single start/end an event literally declares is considered, so a repeating entry currently appears for its first occurrence only. Expanding recurrences is planned for a later iteration.

## Styles

Schedule has a single visual intent — the "what's on now" card — that adapts to the panel rather than splitting into separate named looks. A style is selected automatically by panel capability and resolution; the `style` option lets you pin a specific per-panel variant by name (for example `color-fast-320x240` or `mono-slow-400x300`), but the default auto-selection is almost always what you want.

The card sizes the active event's title in the largest font tier the panel can carry, with its time window on the row beneath and, when `show_upcoming` is enabled, a preview of the next event below that. Long titles and locations are truncated with an ellipsis to fit the panel width.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| source | string | Calendar location: an `http(s)` URL or a local path to an `.ics` file (empty disables fetching) | (empty) | a URL or filesystem path, or empty |
| refresh_sec | int | Seconds between calendar re-fetches | 300 | any integer >= 30 |
| time_format | string | Event time display format | 24h | 24h, 12h |
| show_upcoming | bool | Show the next upcoming event beneath the active one | true | true, false |
| font | string | Font for rendering text | auto | auto, or any registered font ID |
| style | string | Per-panel style name (empty auto-selects by panel) | (empty) | empty, or a registered per-panel style name |

Configure options via the CLI:

```sh
cyberhudctl display schedule <key>=<value> [<key>=<value> ...]
```

!!! note
    `refresh_sec` is clamped to a minimum of 30 seconds to avoid hammering a remote calendar host.

## CLI Examples

Query all current settings (zero arguments returns the full policy state):

```sh
cyberhudctl display schedule
```

Set a local calendar file:

```sh
cyberhudctl display schedule source=/home/pi/calendar.ics
```

Use a remote calendar and refresh every two minutes:

```sh
cyberhudctl display schedule source=https://example.com/cal.ics refresh_sec=120
```

Switch to 12-hour event times:

```sh
cyberhudctl display schedule time_format=12h
```

Hide the upcoming-event preview for a cleaner single-line card:

```sh
cyberhudctl display schedule show_upcoming=false
```

Disable fetching (clears the card):

```sh
cyberhudctl display schedule source=
```

## Panel Compatibility

Works on all panels. The card adapts its title size to the panel height and drops the upcoming-event row on panels too short to carry it. On color panels an active event is accented; on monochrome and grayscale panels it renders in the native foreground. On slow-refresh e-paper panels the card renders as a static frame that updates only when the active event changes or a new fetch lands, avoiding unnecessary redraws.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128, 240×240, 320×240) | Fully supported — white text on black |
| MonoSlow | Slow-refresh monochrome e-ink (128×64, 128×128, 200×200, 250×122, 296×128, 400×300, 800×480) | Fully supported — static frame, updates on event change |
| GrayscaleFast | Fast-refresh grayscale (128×128, 240×240, 320×240, 480×320) | Fully supported — grayscale text |
| GrayscaleSlow | Slow-refresh grayscale e-ink (200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static grayscale frame |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 240×320, 320×240, 320×480, 480×320, 480×800, 800×480) | Fully supported — accented active event, full color |
| ColorSlow | Slow-refresh color e-ink (200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static color frame |

| Input | Action |
|-------|--------|
| K1 / Primary | Navigate to the dashboard mode |

The mode is a read-only status view and needs no input controls — it updates itself as events begin and end.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Clock](clock.md) — current time, date, and weekday display
- [Ticker](ticker.md) — scrolling ticker display from external data feeds
