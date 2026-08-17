// Package gauges renders user-supplied numeric data as progress bars, rings,
// arcs, and dial-style widgets sized for the active panel.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions are published via init()
//     with typed option metadata, summaries, defaults, and allowed-value lists.
//   - Command handling: the "gauges" verb accepts set/get/policy operations for
//     pushing data and adjusting policy in a single mode-local command surface.
//   - Policy definition and normalization: a typed Policy struct carries the
//     rendering knobs (style, shape, labels, accent, grid size, padding) and is
//     normalized before rendering.
//   - Data model: Gauge and GaugeSet provide a canonical internal representation
//     for one or many values with min/max range metadata and a normalized percent.
//   - Data-signature change detection: RenderCacheKey combines policy and gauge
//     fingerprints so unchanged payloads do not trigger unnecessary redraws.
//   - Interface-based style dispatch: a style.NewRegistry[GaugeSet, Policy]
//     holds capability-ordered styles, with the first registered style serving as
//     the default.
//   - Panel hints and fitness evaluation: PanelHints keeps per-region geometry
//     isolated, and style fitness selection chooses the best matching layout for
//     the active panel.
//   - Widget sprite compositing: progressbar widgets render the fills and rings,
//     while textlabel sprites render optional per-gauge labels and empty-state
//     placeholders.
//
// Practical use cases:
//
//   - Home server dashboard: disk usage, CPU load, memory, and temperature.
//   - Maker bench panel: printer progress, nozzle temperature, enclosure humidity.
//   - Solar/battery setup: charge percentage, inverter load, battery voltage.
//   - Aquarium or greenhouse: water level, pump duty cycle, soil moisture.
//
// Payloads can be supplied directly through the gauges command as a single number,
// a JSON object, or a JSON array of gauge objects.
package gauges
