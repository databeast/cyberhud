/*
 * Geometric Background Animation — Performance monitoring and adaptive scaling.
 *
 * Tracks frame times in a rolling 10-frame window and makes scaling decisions:
 * - Reduces total square count by ~50% (minimum 4) when avg frame time > 33ms
 * - Restores 25% of reduced count when avg frame time < 25ms after a reduction
 * - Skips evaluation for clamped frames (>100ms, e.g. tab was backgrounded)
 *
 * Also provides mobile cluster count scaling.
 *
 */

import type {PerfState} from "./types";

/* Maximum rolling window size */
const WINDOW_SIZE = 10;

/* Frame time threshold (ms) above which reduction triggers */
const REDUCTION_THRESHOLD_MS = 33;

/* Frame time threshold (ms) below which restoration triggers */
const RESTORATION_THRESHOLD_MS = 25;

/* Minimum total squares after reduction */
const MIN_SQUARE_COUNT = 4;

/* Maximum frame time considered valid for performance evaluation (ms) */
const CLAMPED_FRAME_THRESHOLD_MS = 100;

export interface PerformanceDecision {
    /* Updated PerfState after evaluation */
    newState: PerfState;
    /* Whether a reduction was triggered this frame */
    shouldReduce: boolean;
    /* Whether a restoration was triggered this frame */
    shouldRestore: boolean;
}

/*
 * Evaluate performance based on frame time and current state.
 *
 * Returns a new PerfState along with scaling decisions. Does NOT mutate
 * the input state.
 *
 * - Skips evaluation entirely if frameTime > 100ms (tab was backgrounded)
 * - Triggers reduction (halve square count, min 4) when rolling avg > 33ms
 * - Triggers restoration (add back 25% of removed count) when rolling avg < 25ms
 *   after a prior reduction, capped at originalSquareCount
 */
export function evaluatePerformance(
    state: PerfState,
    frameTime: number
): PerformanceDecision {
    // Skip evaluation for clamped frames (e.g. tab was backgrounded)
    if (frameTime > CLAMPED_FRAME_THRESHOLD_MS) {
        return {
            newState: {...state},
            shouldReduce: false,
            shouldRestore: false,
        };
    }

    // Add frame time to rolling window (keep last WINDOW_SIZE entries)
    const frameTimes = [...state.frameTimes, frameTime];
    if (frameTimes.length > WINDOW_SIZE) {
        frameTimes.splice(0, frameTimes.length - WINDOW_SIZE);
    }

    // Only evaluate when we have a full window
    if (frameTimes.length < WINDOW_SIZE) {
        return {
            newState: {
                ...state,
                frameTimes,
            },
            shouldReduce: false,
            shouldRestore: false,
        };
    }

    const avg = frameTimes.reduce((sum, t) => sum + t, 0) / frameTimes.length;

    // Check for reduction: avg > 33ms and we haven't already hit the minimum
    if (avg > REDUCTION_THRESHOLD_MS && state.currentSquareCount > MIN_SQUARE_COUNT) {
        const halved = Math.floor(state.currentSquareCount / 2);
        const newSquareCount = Math.max(MIN_SQUARE_COUNT, halved);

        return {
            newState: {
                frameTimes,
                currentSquareCount: newSquareCount,
                originalSquareCount: state.originalSquareCount,
                hasReduced: true,
            },
            shouldReduce: true,
            shouldRestore: false,
        };
    }

    // Check for restoration: avg < 25ms, only if we previously reduced
    if (avg < RESTORATION_THRESHOLD_MS && state.hasReduced && state.currentSquareCount < state.originalSquareCount) {
        const removedCount = state.originalSquareCount - state.currentSquareCount;
        const restoreAmount = Math.max(1, Math.floor(removedCount * 0.25));
        const newSquareCount = Math.min(
            state.originalSquareCount,
            state.currentSquareCount + restoreAmount
        );

        // If fully restored, clear the hasReduced flag
        const hasReduced = newSquareCount < state.originalSquareCount;

        return {
            newState: {
                frameTimes,
                currentSquareCount: newSquareCount,
                originalSquareCount: state.originalSquareCount,
                hasReduced,
            },
            shouldReduce: false,
            shouldRestore: true,
        };
    }

    // No scaling change needed
    return {
        newState: {
            ...state,
            frameTimes,
        },
        shouldReduce: false,
        shouldRestore: false,
    };
}

/*
 * Compute the cluster count for mobile viewports (< 768px).
 *
 * Returns max(3, floor(desktopCount * 0.6)).
 */
export function getMobileClusterCount(desktopCount: number): number {
    return Math.max(3, Math.floor(desktopCount * 0.6));
}

/*
 * Create a fresh PerfState for a given square count.
 * Used when initializing or when viewport category changes.
 */
export function createPerfState(squareCount: number): PerfState {
    return {
        frameTimes: [],
        currentSquareCount: squareCount,
        originalSquareCount: squareCount,
        hasReduced: false,
    };
}

/*
 * Update the originalSquareCount for a new viewport category while
 * preserving reduction state. Used on viewport resize to recompute
 * the restoration cap for the new viewport.
 *
 * If the new original is less than the current count, the current
 * count is also clamped down.
 */
export function updateOriginalForViewport(
    state: PerfState,
    newOriginalSquareCount: number
): PerfState {
    const currentSquareCount = Math.min(state.currentSquareCount, newOriginalSquareCount);
    return {
        ...state,
        originalSquareCount: newOriginalSquareCount,
        currentSquareCount,
        // If current matches original, no longer in reduced state
        hasReduced: currentSquareCount < newOriginalSquareCount ? state.hasReduced : false,
    };
}
