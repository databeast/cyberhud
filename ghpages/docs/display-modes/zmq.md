# ZMQ

The zmq display mode subscribes to ZeroMQ message streams and renders incoming data on the panel display. It supports PUB/SUB and PULL socket types, optional JSON field extraction, and configurable line buffering for real-time message monitoring.

## Quick Start

```sh
cyberhudctl display set 0 zmq
```

## How It Works

The zmq mode connects to a ZeroMQ endpoint (SUB or PULL socket) and ingests messages into a bounded ring buffer. Each incoming message is rendered as a text line on the panel surface, with the newest messages appearing at the bottom of the display. The mode updates the display each time a new message arrives, providing real-time streaming output.

When JSON field filtering is configured via `json_fields`, the mode parses each incoming message as a JSON object and extracts only the specified fields, displaying them as "key: value" lines. Messages that are not valid JSON or that lack the specified fields are displayed verbatim. When no messages have been received, the display shows a waiting state (style-dependent: "Waiting for messages..." or "(no messages)").

The mode supports multiple visual styles that are automatically selected based on your panel's resolution and capability class. Each style fits content to the available panel area, truncating long lines and showing the most recent messages that fit vertically.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `endpoint` | string | ZMQ connection endpoint URL where the mode dials to receive messages (e.g. tcp://localhost:5556) | | any valid string |
| `socket_type` | string | ZeroMQ socket type determining the messaging pattern used to receive data from the endpoint | sub | sub, pull |
| `topic` | string | Subscription topic filter applied to SUB sockets to select which messages are received (ignored for PULL) | | any valid string |
| `max_lines` | int | Ring buffer capacity controlling how many message lines are retained and available for display rendering | 24 | 1–1000 |
| `json_fields` | string | Comma-separated list of JSON field names to extract and display from incoming structured messages | | any valid string |
| `font` | string | Font selection for text rendering, using auto for automatic best-fit or a registered font ID | auto | any valid string |
| `style` | string | Visual presentation style controlling the layout, font sizing, and spatial arrangement of rendered messages | | Resolution-specific style names |

## CLI Examples

Activate the zmq mode on the main display region:

```sh
cyberhudctl display set 0 zmq
```

Configure the ZMQ endpoint and socket type for a PUB/SUB stream:

```sh
cyberhudctl display zmq endpoint=tcp://localhost:5556 socket_type=sub topic=sensor
```

Use PULL socket with JSON field filtering to show only specific data:

```sh
cyberhudctl display zmq endpoint=tcp://192.168.1.10:5557 socket_type=pull json_fields=temperature,humidity
```

Increase the message buffer to retain more history:

```sh
cyberhudctl display zmq max_lines=100
```

Query current zmq mode settings:

```sh
cyberhudctl display zmq
```

## Panel Compatibility

The zmq mode is a text-rendering mode with no interactive controls required. It works across all panel capability classes, adapting its output to the available display surface.

| Capability | Description | Behavior |
|------------|-------------|----------|
| MonoFast | Monochrome OLED with fast refresh | Fully supported. Text rendered in native foreground color with automatic font sizing. |
| MonoSlow | Monochrome e-ink with slow refresh | Fully supported. Sharp border frame on large panels. Updates only on new messages to minimize redraws. |
| GrayscaleFast | Grayscale display with fast refresh | Fully supported. Text rendered in white/light gray with automatic layout. |
| GrayscaleSlow | Grayscale e-ink with slow refresh | Fully supported. Static rendering with centered text block. Updates on new messages only. |
| ColorFast | Color TFT with fast refresh | Fully supported. Full-color text rendering with automatic layout. |
| ColorSlow | Color e-ink with slow refresh | Fully supported. Static rendering mode avoids unnecessary redraws between messages. |

No input controls (buttons or joystick) are required. The mode has no minimum resolution requirement — content is automatically truncated and fitted to any panel size.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Ticker](ticker.md) — external data feed display with similar streaming text
- [Pager](pager.md) — file/pipe tailing mode with similar text rendering


<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/zmq/color-slow-122x250_0001.png" alt="color-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-128x128_0001.png" alt="color-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-160x80_0001.png" alt="color-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-slow-176x264_0001.png" alt="color-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-240x135_0001.png" alt="color-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-240x240_0001.png" alt="color-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-240x320_0001.png" alt="color-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-320x240_0001.png" alt="color-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-320x480_0001.png" alt="color-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-480x320_0001.png" alt="color-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-480x800_0001.png" alt="color-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-fast-800x480_0001.png" alt="color-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/zmq/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/zmq/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-160x80_0001.png" alt="grayscale-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-240x135_0001.png" alt="grayscale-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-240x240_0001.png" alt="grayscale-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-320x240_0001.png" alt="grayscale-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-480x320_0001.png" alt="grayscale-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/zmq/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/zmq/mono-slow-64x128_0001.png" alt="mono-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-104x212_0001.png" alt="mono-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-122x250_0001.png" alt="mono-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-128x32_0001.png" alt="mono-fast-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-128x64_0001.png" alt="mono-fast-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-128x64_0001.png" alt="mono-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-128x128_0001.png" alt="mono-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-160x80_0001.png" alt="mono-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-176x264_0001.png" alt="mono-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-200x200_0001.png" alt="mono-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-212x104_0001.png" alt="mono-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-240x135_0001.png" alt="mono-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-240x240_0001.png" alt="mono-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-250x122_0001.png" alt="mono-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-264x176_0001.png" alt="mono-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-296x128_0001.png" alt="mono-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-300x400_0001.png" alt="mono-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-320x240_0001.png" alt="mono-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-400x300_0001.png" alt="mono-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-480x320_0001.png" alt="mono-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-480x800_0001.png" alt="mono-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-fast-800x480_0001.png" alt="mono-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/zmq/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->

