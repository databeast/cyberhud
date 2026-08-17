// Package region is the core display region infrastructure for cyberhud.
//
// It manages virtual display framebuffers, region allocation, mode lifecycle,
// per-region tick scheduling, flush paths to hardware, and the render loop.
//
// # Architectural Layers
//
// The package is organized into five layers with a top-to-bottom data flow:
//
//   - VirtualDisplay: Unified RGBA framebuffer backing all physical screens on a panel.
//   - RegionManager: Owns region lifecycle — allocation, validation, lookup, and mode routing — operating on sub-image views of the VirtualDisplay.
//   - FlushPath: Extracts per-screen rectangles from the VirtualDisplay and pushes pixel data to hardware draw targets.
//   - RenderLoop: Tick-driven renderer that coordinates per-frame rendering and input dispatch across all regions managed by the RegionManager, flushing results through the FlushPath.
//   - ModeSwitch: Command handler that routes runtime mode-change requests through the RegionManager to individual regions.
//
// # How to Wire a Region from Scratch
//
// To set up a fully functional region infrastructure:
//
//  1. Call [ActivatePanel] with a [PanelActivationConfig] containing screen positions,
//     layout, and available modes. This constructs the [VirtualDisplay],
//     [RegionManager], [FlushPath], [RenderLoop], and [ModeSwitch].
//  2. For each [Region], call [Region.SetModeFactory] to inject the [ModeFactory]
//     that constructs [ModeInstance] values for that region's available display modes.
//  3. Call [RenderLoop.Apply] with [WithRenderer] to set the [Renderer] responsible
//     for rendering each region's active mode instance every frame.
//  4. Call [RenderLoop.Apply] with [WithInputDispatcher] to set the [InputDispatcher]
//     that routes input events to the input-active region.
//  5. Call [RenderLoop.Apply] with [WithTickRateResolver] to set the [TickRateResolver]
//     that determines per-region tick intervals (enables per-region deadline scheduling).
//  6. Call [RenderLoop.Run] to start the render loop (blocks until [RenderLoop.Stop]).
//
// # Zero-Copy Sub-Image Relationship
//
// Region surfaces share underlying pixel memory with the VirtualDisplay via
// [VirtualDisplay.SubImage]. When a region is allocated, its surface is backed
// by a sub-image view of the VirtualDisplay's framebuffer. Writes to the
// region's surface are immediately visible in the VirtualDisplay without any
// copy step. This zero-copy design means the FlushPath can read rendered pixels
// directly from the VirtualDisplay framebuffer after regions render in place.
//
// # SetMode Lifecycle
//
// [Region.SetMode] performs a three-phase lifecycle transition with panic recovery:
//
//  1. Construct: The [ModeFactory] is called to build a new [ModeInstance] for the
//     requested mode. If the factory panics, the panic is recovered, the previous
//     mode is retained, and an error is returned.
//  2. Activate: The new instance's Activate method is called. If Activate panics,
//     the panic is recovered, a cleanup Deactivate is attempted on the new instance
//     (itself wrapped in a separate recover), the new instance is discarded, the
//     previous mode is retained, and an error is returned.
//  3. Deactivate old: The previous instance's Deactivate method runs in a goroutine
//     with a 5-second timeout. If deactivation does not complete within 5 seconds,
//     the timeout is logged and the loop continues without blocking.
//
// [Region.SetModeFactory] must be called before SetMode. If no factory is wired,
// SetMode panics to surface the misconfiguration immediately at startup.
//
// # Per-Region Tick Scheduling
//
// Each display mode may declare its own preferred tick interval via the
// [TickRateProvider] interface. Modes register themselves by calling
// [RegisterTickRate] from their package init functions, associating a mode ID
// with a provider.
//
// The [DefaultTickRateResolver] queries the registry of registered providers to
// determine the tick interval for a given mode. If a provider is registered, its
// PreferredTickInterval is clamped to the bounds [1 ms, 10000 ms]. If no provider
// is registered for a mode, the default interval of 1000 ms is used.
//
// When a [TickRateResolver] is configured on the [RenderLoop] via
// [WithTickRateResolver], the loop switches from a single global ticker to
// per-region deadline scheduling. Each region maintains an independent deadline
// and only regions whose deadline has elapsed are rendered each cycle.
package region
