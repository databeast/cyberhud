/*
 * AttractCanvas — Geometric background animation for the hero section.
 *
 * Renders clusters of overlapping green squares that fade in and out on
 * independent timing phases with angular rotations and center-bright glow,
 * inspired by Necron tomb architecture. Pseudocode text fragments materialize
 * over the clusters to reinforce the coding identity.
 *
 * This module replaces the previous bokeh/particle animation while retaining
 * the same export signatures (HSLColor, AttractCanvasConfig, AttractAnimation,
 * createAttractAnimation) for seamless integration.
 *
 */

import {AnimationLoop} from './geometric/animation-loop';
import {MIN_DESKTOP_CLUSTERS} from './geometric/types';

// ---------------------------------------------------------------------------
// Public interfaces (unchanged from previous implementation)
// ---------------------------------------------------------------------------

export interface HSLColor {
    /* Hue 0-360 */
    h: number;
    /* Saturation 0-100 */
    s: number;
    /* Lightness 0-100 */
    l: number;
}

export interface AttractCanvasConfig {
    /* Target canvas element ID */
    canvasId: string;
    /* Base particle count — mapped to cluster count (floor of MIN_DESKTOP_CLUSTERS) */
    baseParticleCount: number;
    /* Animation speed multiplier (accepted but not directly used; fade durations are internal) */
    speed: number;
    /* Color palette (accepted but ignored; internal green palette is used) */
    palette: HSLColor[];
    /* Whether to start in reduced-motion mode */
    reducedMode: boolean;
}

export interface AttractAnimation {
    start(): void;

    stop(): void;

    resize(width: number, height: number): void;

    setReducedMotion(reduced: boolean): void;
}

// ---------------------------------------------------------------------------
// No-op implementation for error cases
// ---------------------------------------------------------------------------

function createNoOpAnimation(): AttractAnimation {
    return {
        start() {
        },
        stop() {
        },
        resize() {
        },
        setReducedMotion() {
        },
    };
}

// ---------------------------------------------------------------------------
// Mobile cluster count calculation
// ---------------------------------------------------------------------------

const MOBILE_BREAKPOINT = 768;

/*
 * Computes the cluster count for the current viewport width.
 * - Desktop (>=768px): max(baseCount, MIN_DESKTOP_CLUSTERS)
 * - Mobile (<768px): max(3, floor(desktopCount * 0.6))
 */
function getClusterCount(baseCount: number, viewportWidth: number): number {
    const desktopCount = Math.max(baseCount, MIN_DESKTOP_CLUSTERS);
    if (viewportWidth < MOBILE_BREAKPOINT) {
        return Math.max(3, Math.floor(desktopCount * 0.6));
    }
    return desktopCount;
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

/*
 * Creates a geometric background animation bound to the specified canvas.
 *
 * @param config - Configuration matching the existing AttractCanvasConfig signature
 * @returns An AttractAnimation interface with start/stop/resize/setReducedMotion
 */
export function createAttractAnimation(config: AttractCanvasConfig): AttractAnimation {
    // --- Error case: canvas element not found ---
    const canvas = document.getElementById(config.canvasId) as HTMLCanvasElement | null;
    if (!canvas) {
        return createNoOpAnimation();
    }

    // --- Error case: getContext('2d') returns null ---
    const ctx = canvas.getContext('2d');
    if (!ctx) {
        return createNoOpAnimation();
    }

    // --- Derive initial dimensions ---
    let width = canvas.clientWidth || window.innerWidth;
    let height = canvas.clientHeight || window.innerHeight;

    // --- Map baseParticleCount to cluster count ---
    const baseClusterCount = Math.max(
        Math.floor(config.baseParticleCount),
        MIN_DESKTOP_CLUSTERS,
    );

    // --- Create the animation loop ---
    const loop = new AnimationLoop({
        width,
        height,
        baseClusterCount,
        reducedMotion: config.reducedMode,
    });

    // --- Track running state locally for idempotent start/stop ---
    let running = false;
    let reducedMotion = config.reducedMode;

    // -------------------------------------------------------------------------
    // Public API
    // -------------------------------------------------------------------------

    function start(): void {
        // Idempotent: no-op if already running
        if (running) return;
        // Cannot start if reduced motion is active
        if (reducedMotion) return;

        // Set canvas dimensions
        canvas!.width = width;
        canvas!.height = height;

        running = true;
        loop.startLoop(ctx!);
    }

    function stop(): void {
        // Idempotent: no-op if already stopped
        if (!running) return;

        running = false;
        loop.stopLoop();
    }

    function resize(newWidth: number, newHeight: number): void {
        width = newWidth;
        height = newHeight;

        // Update canvas element dimensions
        canvas!.width = width;
        canvas!.height = height;

        // Delegate to the animation loop's resize (handles redistributeClusters internally)
        loop.resize(width, height);
    }

    function setReducedMotion(reduced: boolean): void {
        reducedMotion = reduced;

        if (reduced && running) {
            // Immediately stop the render loop
            running = false;
            loop.stopLoop();
        }

        // Propagate to the loop so future startLoop() calls respect the flag
        loop.setReducedMotion(reduced);

        // When set to false: allow next start() to work but do NOT auto-start
    }

    return {start, stop, resize, setReducedMotion};
}
