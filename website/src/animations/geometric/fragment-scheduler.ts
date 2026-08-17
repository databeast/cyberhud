/*
 * Geometric Background Animation — Fragment Scheduler.
 *
 * Pure functional state machine managing pseudocode fragment lifecycle.
 * Handles spawning (with all constraints), opacity computation during
 * fade-in/hold/fade-out phases, and expiration of completed fragments.
 *
 * All functions accept and return immutable state for testability.
 * Randomization is injected via an RNG function parameter.
 *
 */

import type {ActiveFragment, FragmentSchedulerState, HSLColor,} from './types';
import {
    FRAGMENT_CYAN_HUE_RANGE,
    FRAGMENT_CYAN_LIGHTNESS,
    FRAGMENT_CYAN_SAT_MIN,
    FRAGMENT_GREEN_HUE_RANGE,
    FRAGMENT_GREEN_LIGHTNESS,
    FRAGMENT_GREEN_SAT_MIN,
    PSEUDOCODE_POOL,
} from './types';

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

/* Maximum number of simultaneously active fragments (Requirement 3.2) */
const MAX_ACTIVE_FRAGMENTS = 3;

/* Minimum time between spawns in seconds (Requirement 3.6) */
const MIN_SPAWN_INTERVAL = 3;

/* Proximity radius in pixels for fragment positioning (Requirement 3.7) */
const CLUSTER_PROXIMITY_RADIUS = 100;

/* Fragment fade-in duration range in seconds (Requirement 3.4) */
const FADE_IN_RANGE: [number, number] = [1, 2];

/* Fragment hold duration range in seconds (Requirement 3.4) */
const HOLD_RANGE: [number, number] = [2, 5];

/* Fragment fade-out duration range in seconds (Requirement 3.4) */
const FADE_OUT_RANGE: [number, number] = [1, 2];

/* Font size range in px (Requirements 3.5, 6.2) */
const FONT_SIZE_RANGE: [number, number] = [10, 16];

/* Peak opacity range (Requirements 3.5, 6.2) */
const PEAK_OPACITY_RANGE: [number, number] = [0.3, 0.8];

// ---------------------------------------------------------------------------
// Active cluster info passed by caller
// ---------------------------------------------------------------------------

/*
 * Represents a cluster that is currently "active" (at least one square
 * has opacity >= 0.1). The caller computes this from cluster configs and
 * fade cycle state, then passes the relevant info to the scheduler.
 */
export interface ActiveClusterInfo {
    /* Absolute X position of cluster center in pixels */
    centerX: number;
    /* Absolute Y position of cluster center in pixels */
    centerY: number;
}

// ---------------------------------------------------------------------------
// State factory
// ---------------------------------------------------------------------------

/*
 * Creates a fresh FragmentSchedulerState.
 */
export function createFragmentSchedulerState(): FragmentSchedulerState {
    return {
        activeFragments: [],
        lastSpawnTime: -Infinity,
        lastSpawnedText: null,
    };
}

// ---------------------------------------------------------------------------
// Fragment opacity computation
// ---------------------------------------------------------------------------

/*
 * Computes the current opacity of a fragment based on its lifecycle phase.
 *
 * Lifecycle: fade-in → hold → fade-out
 *
 * @param fragment - The active fragment
 * @param time - Current animation time in seconds
 * @returns Current opacity in [0, peakOpacity], or -1 if expired
 */
export function computeFragmentOpacity(
    fragment: ActiveFragment,
    time: number,
): number {
    const elapsed = time - fragment.startTime;

    if (elapsed < 0) {
        return 0;
    }

    const {fadeInDuration, holdDuration, fadeOutDuration, peakOpacity} = fragment;
    const totalDuration = fadeInDuration + holdDuration + fadeOutDuration;

    if (elapsed >= totalDuration) {
        // Fragment has expired
        return -1;
    }

    if (elapsed < fadeInDuration) {
        // Fade-in phase: linear ramp from 0 to peakOpacity
        return peakOpacity * (elapsed / fadeInDuration);
    }

    if (elapsed < fadeInDuration + holdDuration) {
        // Hold phase: constant at peakOpacity
        return peakOpacity;
    }

    // Fade-out phase: linear ramp from peakOpacity to 0
    const fadeOutElapsed = elapsed - fadeInDuration - holdDuration;
    return peakOpacity * (1 - fadeOutElapsed / fadeOutDuration);
}

// ---------------------------------------------------------------------------
// Update fragments (remove expired)
// ---------------------------------------------------------------------------

/*
 * Updates fragment state by computing current opacities and removing expired
 * fragments. Returns the new state with only active (non-expired) fragments.
 *
 * @param state - Current scheduler state
 * @param time - Current animation time in seconds
 * @returns Updated state with expired fragments removed
 */
export function updateFragments(
    state: FragmentSchedulerState,
    time: number,
): FragmentSchedulerState {
    const activeFragments = state.activeFragments.filter((fragment) => {
        const opacity = computeFragmentOpacity(fragment, time);
        return opacity >= 0; // -1 means expired
    });

    return {
        ...state,
        activeFragments,
    };
}

// ---------------------------------------------------------------------------
// Spawn fragment
// ---------------------------------------------------------------------------

/*
 * Attempts to spawn a new pseudocode fragment, enforcing all scheduling
 * constraints. Returns the updated state (which may be unchanged if
 * constraints prevent spawning).
 *
 * Constraints checked:
 * - Max 3 simultaneous fragments (Requirement 3.2)
 * - At least 3 seconds since last spawn (Requirement 3.6)
 * - At least one cluster must be active (Requirement 3.8)
 * - No consecutive same text (Requirement 3.9)
 *
 * @param state - Current scheduler state
 * @param time - Current animation time in seconds
 * @param activeClusters - List of currently active clusters with absolute positions
 * @param viewportWidth - Viewport width in pixels
 * @param viewportHeight - Viewport height in pixels
 * @param rng - Random number generator function returning [0, 1)
 * @returns Updated state (possibly with a new fragment added)
 */
export function spawnFragment(
    state: FragmentSchedulerState,
    time: number,
    activeClusters: ActiveClusterInfo[],
    viewportWidth: number,
    viewportHeight: number,
    rng: () => number,
): FragmentSchedulerState {
    // Constraint: max 3 simultaneous fragments
    if (state.activeFragments.length >= MAX_ACTIVE_FRAGMENTS) {
        return state;
    }

    // Constraint: at least 3 seconds since last spawn
    if (time - state.lastSpawnTime < MIN_SPAWN_INTERVAL) {
        return state;
    }

    // Constraint: at least one cluster must be active
    if (activeClusters.length === 0) {
        return state;
    }

    // Select text from pool (no consecutive same text)
    const text = selectText(state.lastSpawnedText, rng);

    // Select a random active cluster to position near
    const clusterIndex = Math.floor(rng() * activeClusters.length);
    const cluster = activeClusters[clusterIndex];

    // Position within 100px of the cluster center
    const {x, y} = computeFragmentPosition(
        cluster.centerX,
        cluster.centerY,
        viewportWidth,
        viewportHeight,
        rng,
    );

    // Generate timing parameters
    const fadeInDuration = randomInRange(FADE_IN_RANGE[0], FADE_IN_RANGE[1], rng);
    const holdDuration = randomInRange(HOLD_RANGE[0], HOLD_RANGE[1], rng);
    const fadeOutDuration = randomInRange(FADE_OUT_RANGE[0], FADE_OUT_RANGE[1], rng);

    // Generate style parameters
    const fontSize = randomInRange(FONT_SIZE_RANGE[0], FONT_SIZE_RANGE[1], rng);
    const peakOpacity = randomInRange(PEAK_OPACITY_RANGE[0], PEAK_OPACITY_RANGE[1], rng);
    const color = generateFragmentColor(rng);

    const newFragment: ActiveFragment = {
        text,
        x,
        y,
        startTime: time,
        fadeInDuration,
        holdDuration,
        fadeOutDuration,
        fontSize,
        color,
        peakOpacity,
    };

    return {
        activeFragments: [...state.activeFragments, newFragment],
        lastSpawnTime: time,
        lastSpawnedText: text,
    };
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

/*
 * Selects a text snippet from the pool, avoiding consecutive repeats.
 */
function selectText(lastText: string | null, rng: () => number): string {
    const pool = PSEUDOCODE_POOL;

    if (lastText === null) {
        // No previous text, pick any
        return pool[Math.floor(rng() * pool.length)];
    }

    // Filter out the last text to avoid consecutive repeats
    const available = pool.filter((t) => t !== lastText);

    // Fallback: if pool has only one entry (shouldn't happen with 10+), allow repeat
    if (available.length === 0) {
        return pool[Math.floor(rng() * pool.length)];
    }

    return available[Math.floor(rng() * available.length)];
}

/*
 * Computes a fragment position within CLUSTER_PROXIMITY_RADIUS of a cluster center,
 * clamped to viewport bounds.
 */
function computeFragmentPosition(
    clusterCenterX: number,
    clusterCenterY: number,
    viewportWidth: number,
    viewportHeight: number,
    rng: () => number,
): { x: number; y: number } {
    // Random angle and distance within the proximity radius
    const angle = rng() * 2 * Math.PI;
    const distance = rng() * CLUSTER_PROXIMITY_RADIUS;

    let x = clusterCenterX + Math.cos(angle) * distance;
    let y = clusterCenterY + Math.sin(angle) * distance;

    // Clamp to viewport bounds
    x = Math.max(0, Math.min(viewportWidth, x));
    y = Math.max(0, Math.min(viewportHeight, y));

    return {x, y};
}

/*
 * Generates a random fragment color in either the green or cyan range.
 * Approximately 50/50 chance between green and cyan.
 */
function generateFragmentColor(rng: () => number): HSLColor {
    const useGreen = rng() < 0.5;

    if (useGreen) {
        return {
            h: randomInRange(FRAGMENT_GREEN_HUE_RANGE[0], FRAGMENT_GREEN_HUE_RANGE[1], rng),
            s: randomInRange(FRAGMENT_GREEN_SAT_MIN, 100, rng),
            l: randomInRange(FRAGMENT_GREEN_LIGHTNESS[0], FRAGMENT_GREEN_LIGHTNESS[1], rng),
        };
    }

    return {
        h: randomInRange(FRAGMENT_CYAN_HUE_RANGE[0], FRAGMENT_CYAN_HUE_RANGE[1], rng),
        s: randomInRange(FRAGMENT_CYAN_SAT_MIN, 100, rng),
        l: randomInRange(FRAGMENT_CYAN_LIGHTNESS[0], FRAGMENT_CYAN_LIGHTNESS[1], rng),
    };
}

/*
 * Returns a random number in [min, max] using the provided RNG.
 */
function randomInRange(min: number, max: number, rng: () => number): number {
    return min + rng() * (max - min);
}
