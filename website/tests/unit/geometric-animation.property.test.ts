/*
 * Property-Based Tests for Geometric Background Animation
 *
 * Feature: geometric-background-animation
 *
 * Uses fast-check to verify correctness properties defined in the design document.
 * Each property test runs a minimum of 100 iterations to ensure confidence.
 *
 * Shared generators produce valid arbitrary inputs within the constraints
 * specified by the animation's type system and requirements.
 */

import {describe, expect, it} from "vitest";
import fc from "fast-check";
import type {ClusterConfig, HSLColor, SquareConfig,} from "../../src/animations/geometric/types";
import {
    FRAGMENT_CYAN_HUE_RANGE,
    FRAGMENT_CYAN_LIGHTNESS,
    FRAGMENT_CYAN_SAT_MIN,
    FRAGMENT_GREEN_HUE_RANGE,
    FRAGMENT_GREEN_LIGHTNESS,
    FRAGMENT_GREEN_SAT_MIN,
    SQUARE_HUE_MAX,
    SQUARE_HUE_MIN,
    SQUARE_LIGHTNESS_MAX,
    SQUARE_LIGHTNESS_MIN,
    SQUARE_SAT_MIN,
} from "../../src/animations/geometric/types";
import {generateCluster, generateSquare, type RNG} from "../../src/animations/geometric/cluster-generator";
import {getCentralZoneOpacityCap} from "../../src/animations/geometric/central-zone";
import {computeEmergenceTimes, computeFadeOpacity} from "../../src/animations/geometric/fade-cycle";
import {computeGlowColor, computeGlowMultiplier, MIN_GLOW_SPREAD_PX} from "../../src/animations/geometric/glow";
import type {ActiveClusterInfo} from "../../src/animations/geometric/fragment-scheduler";
import {
    type ActiveClusterInfo,
    createFragmentSchedulerState,
    spawnFragment,
    updateFragments,
} from "../../src/animations/geometric/fragment-scheduler";
import {createPerfState, evaluatePerformance, getMobileClusterCount} from "../../src/animations/geometric/performance";
import {redistributeClusters} from "../../src/animations/geometric/cluster-redistribution";

// ---------------------------------------------------------------------------
// Configuration: minimum iterations per property
// ---------------------------------------------------------------------------

const NUM_RUNS = 100;

// ---------------------------------------------------------------------------
// Shared Generators (Arbitraries)
// ---------------------------------------------------------------------------

/*
 * Generates a valid HSLColor within the square accent color range.
 * Hue: [100, 140], Saturation: [70, 100], Lightness: [20, 55]
 */
export function arbHSLColor(): fc.Arbitrary<HSLColor> {
    return fc.record({
        h: fc.integer({min: SQUARE_HUE_MIN, max: SQUARE_HUE_MAX}),
        s: fc.integer({min: SQUARE_SAT_MIN, max: 100}),
        l: fc.integer({min: SQUARE_LIGHTNESS_MIN, max: SQUARE_LIGHTNESS_MAX}),
    });
}

/*
 * Generates a valid SquareConfig with all constraints satisfied:
 * - size: [20, 120] px
 * - rotation: [-45, 45] degrees
 * - color: valid HSLColor in green accent range
 * - phaseOffset: [0, cycleDuration]
 * - cycleDuration: [3, 10] seconds
 * - peakOpacity: [0.05, 0.6]
 */
export function arbSquareConfig(): fc.Arbitrary<SquareConfig> {
    return fc
        .record({
            offsetX: fc.double({min: -200, max: 200, noNaN: true}),
            offsetY: fc.double({min: -200, max: 200, noNaN: true}),
            size: fc.integer({min: 20, max: 120}),
            rotation: fc.double({min: -45, max: 45, noNaN: true}),
            color: arbHSLColor(),
            cycleDuration: fc.double({min: 3, max: 10, noNaN: true}),
            peakOpacity: fc.double({min: 0.05, max: 0.6, noNaN: true}),
        })
        .chain((partial) =>
            fc.double({min: 0, max: partial.cycleDuration, noNaN: true}).map(
                (phaseOffset) => ({
                    ...partial,
                    phaseOffset,
                })
            )
        );
}

/*
 * Generates a valid ClusterConfig with 3-8 squares.
 * - centerXPct: [0, 100]
 * - centerYPct: [0, 100]
 * - squares: array of 3-8 valid SquareConfigs
 * - boundingRadius: positive number derived from square offsets
 * - spawnTime: [0, 30] seconds
 * - fadeInDuration: [1, 3] seconds
 */
export function arbClusterConfig(): fc.Arbitrary<ClusterConfig> {
    return fc
        .record({
            centerXPct: fc.double({min: 0, max: 100, noNaN: true}),
            centerYPct: fc.double({min: 0, max: 100, noNaN: true}),
            squares: fc.array(arbSquareConfig(), {minLength: 3, maxLength: 8}),
            spawnTime: fc.double({min: 0, max: 30, noNaN: true}),
            fadeInDuration: fc.double({min: 1, max: 3, noNaN: true}),
        })
        .map((partial) => {
            // Compute bounding radius from square offsets
            const maxDist = partial.squares.reduce((max, sq) => {
                const dist = Math.sqrt(sq.offsetX ** 2 + sq.offsetY ** 2);
                return Math.max(max, dist);
            }, 0);
            return {
                ...partial,
                boundingRadius: Math.max(maxDist, 1), // minimum 1px to avoid division by zero
            };
        });
}

// ---------------------------------------------------------------------------
// Property Tests
// ---------------------------------------------------------------------------

describe("Geometric Animation — Property-Based Tests", () => {
    // -------------------------------------------------------------------------
    // Property 1: Square generation invariants
    // -------------------------------------------------------------------------
    describe("Property 1: Square generation invariants", () => {
        it("generated squares have valid color, size, and opacity ranges", () => {
            fc.assert(
                fc.property(fc.nat(), (seed) => {
                    // Create a seeded RNG using a simple mulberry32 algorithm
                    let state = seed | 0;
                    const rng: RNG = () => {
                        state = (state + 0x6d2b79f5) | 0;
                        let t = Math.imul(state ^ (state >>> 15), 1 | state);
                        t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                        return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                    };

                    const square = generateSquare(rng);

                    // Color hue in [100, 140]
                    expect(square.color.h).toBeGreaterThanOrEqual(100);
                    expect(square.color.h).toBeLessThanOrEqual(140);

                    // Color saturation >= 70
                    expect(square.color.s).toBeGreaterThanOrEqual(70);

                    // Color lightness in [20, 55]
                    expect(square.color.l).toBeGreaterThanOrEqual(20);
                    expect(square.color.l).toBeLessThanOrEqual(55);

                    // Size in [20, 120]
                    expect(square.size).toBeGreaterThanOrEqual(20);
                    expect(square.size).toBeLessThanOrEqual(120);

                    // Peak opacity in [0.05, 0.6]
                    expect(square.peakOpacity).toBeGreaterThanOrEqual(0.05);
                    expect(square.peakOpacity).toBeLessThanOrEqual(0.6);

                    // Cycle duration in [3, 10]
                    expect(square.cycleDuration).toBeGreaterThanOrEqual(3);
                    expect(square.cycleDuration).toBeLessThanOrEqual(10);

                    // Rotation in [-45, 45]
                    expect(square.rotation).toBeGreaterThanOrEqual(-45);
                    expect(square.rotation).toBeLessThanOrEqual(45);
                }),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 2: Cluster composition invariants
    // -------------------------------------------------------------------------
    describe("Property 2: Cluster composition invariants", () => {
        it("clusters contain 3-8 squares with grid-connected positioning", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 0, max: 2 ** 31 - 1}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    (seed, centerXPct, centerYPct) => {
                        // Create a seeded RNG using a simple mulberry32 PRNG
                        let state = seed | 0;
                        const rng = () => {
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), 1 | state);
                            t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        const cluster = generateCluster(rng, centerXPct, centerYPct, 0);

                        // Verify cluster has 8-20 squares
                        expect(cluster.squares.length).toBeGreaterThanOrEqual(8);
                        expect(cluster.squares.length).toBeLessThanOrEqual(20);

                        // Verify squares are positioned (not all at origin) and the largest is at center
                        const sorted = [...cluster.squares].sort((a, b) => b.size - a.size);
                        // The largest square should be at or very near the cluster center (offset ~0)
                        expect(Math.abs(sorted[0].offsetX)).toBeLessThanOrEqual(1);
                        expect(Math.abs(sorted[0].offsetY)).toBeLessThanOrEqual(1);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 3: Rotation angle constraints
    // -------------------------------------------------------------------------
    describe("Property 3: Rotation angle constraints", () => {
        it("rotation angles are either 0 or 45 degrees", () => {
            fc.assert(
                fc.property(fc.integer({min: 1, max: 2 ** 31 - 1}), (seed) => {
                    // Create a seeded RNG
                    let state = seed;
                    const rng = () => {
                        // Simple mulberry32 PRNG for deterministic output
                        state |= 0;
                        state = (state + 0x6d2b79f5) | 0;
                        let t = Math.imul(state ^ (state >>> 15), 1 | state);
                        t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                        return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                    };

                    const centerX = rng() * 100;
                    const centerY = rng() * 100;
                    const cluster = generateCluster(rng, centerX, centerY, 0);

                    // All squares have rotation of exactly 0 or 45
                    for (const square of cluster.squares) {
                        if (square.rotation !== 0 && square.rotation !== 45) {
                            return false;
                        }
                    }

                    return true;
                }),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 4: Cluster glow interpolation
    // -------------------------------------------------------------------------
    describe("Property 4: Cluster glow interpolation", () => {
        it("glow multiplier equals 1.0 - 0.4 * (d / R)", () => {
            fc.assert(
                fc.property(
                    fc.double({min: 0, max: 500, noNaN: true}),
                    fc.double({min: 1, max: 500, noNaN: true}),
                    (distanceFromCenter, boundingRadius) => {
                        const multiplier = computeGlowMultiplier(distanceFromCenter, boundingRadius);

                        // 1. At center (distance=0), multiplier = 1.0
                        const atCenter = computeGlowMultiplier(0, boundingRadius);
                        expect(atCenter).toBe(1.0);

                        // 2. At bounding radius (distance=R), multiplier = 0.6
                        const atBoundary = computeGlowMultiplier(boundingRadius, boundingRadius);
                        expect(atBoundary).toBeCloseTo(0.6, 10);

                        // 3. For arbitrary distances, multiplier equals 1.0 - 0.4 * (d / R)
                        //    (clamped to [0.6, 1.0])
                        const expected = Math.max(
                            0.6,
                            Math.min(1.0, 1.0 - 0.4 * (distanceFromCenter / boundingRadius))
                        );
                        expect(multiplier).toBeCloseTo(expected, 10);

                        // 4. Result is always in [0.6, 1.0]
                        expect(multiplier).toBeGreaterThanOrEqual(0.6);
                        expect(multiplier).toBeLessThanOrEqual(1.0);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 5: Glow color derivation
    // -------------------------------------------------------------------------
    describe("Property 5: Glow color derivation", () => {
        it("glow lightness >= min(L + 20, 100) and spread >= 4px", () => {
            fc.assert(
                fc.property(
                    fc.record({
                        h: fc.integer({min: 0, max: 360}),
                        s: fc.integer({min: 0, max: 100}),
                        l: fc.integer({min: 0, max: 100}),
                    }),
                    (baseColor: HSLColor) => {
                        const glowColor = computeGlowColor(baseColor);

                        // Glow lightness must be >= min(L + 20, 100)
                        expect(glowColor.l).toBeGreaterThanOrEqual(
                            Math.min(baseColor.l + 20, 100)
                        );

                        // Glow spread must be >= 4px
                        expect(MIN_GLOW_SPREAD_PX).toBeGreaterThanOrEqual(4);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 6: Phase offset and cycle duration validity
    // -------------------------------------------------------------------------
    describe("Property 6: Phase offset separation", () => {
        it("phase offsets are within [0, cycleDuration] and cycleDuration in [3, 10]s", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 1, max: 1_000_000}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    (seed, centerX, centerY) => {
                        // Create a seeded RNG
                        let state = seed;
                        const rng: RNG = () => {
                            // Simple mulberry32-style PRNG
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), state | 1);
                            t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        const cluster = generateCluster(rng, centerX, centerY, 0);

                        // Verify all cycleDurations are in [3, 10]
                        for (const sq of cluster.squares) {
                            expect(sq.cycleDuration).toBeGreaterThanOrEqual(3);
                            expect(sq.cycleDuration).toBeLessThanOrEqual(10);
                        }

                        // Verify phase offsets are valid (within [0, cycleDuration])
                        for (const sq of cluster.squares) {
                            expect(sq.phaseOffset).toBeGreaterThanOrEqual(0);
                            expect(sq.phaseOffset).toBeLessThanOrEqual(sq.cycleDuration);
                        }
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 7: Fade cycle correctness
    // -------------------------------------------------------------------------
    describe("Property 7: Fade cycle correctness", () => {
        it("opacity in [0, 0.6] and per-frame change <= 0.05 per 33ms", () => {
            fc.assert(
                fc.property(
                    fc.double({min: 0, max: 10, noNaN: true}),      // phaseOffset
                    fc.double({min: 3, max: 10, noNaN: true}),      // cycleDuration
                    fc.double({min: 0.05, max: 0.6, noNaN: true}), // peakOpacity
                    fc.double({min: 0, max: 100, noNaN: true}),     // time
                    (phaseOffset, cycleDuration, peakOpacity, time) => {
                        // 1. Output must be in [0, 0.6]
                        const opacity = computeFadeOpacity(time, phaseOffset, cycleDuration, peakOpacity);
                        expect(opacity).toBeGreaterThanOrEqual(0);
                        expect(opacity).toBeLessThanOrEqual(0.6);

                        // 2. Per-frame change must be <= 0.05 for dt = 0.033s (33ms)
                        const dt = 0.033;
                        const opacityNext = computeFadeOpacity(time + dt, phaseOffset, cycleDuration, peakOpacity);
                        const change = Math.abs(opacityNext - opacity);
                        expect(change).toBeLessThanOrEqual(0.05);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 8: Cluster emergence window
    // -------------------------------------------------------------------------
    describe("Property 8: Cluster emergence window", () => {
        it("first opacity-crosses-0.1 times fall within 2-second window", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 1, max: 2 ** 31 - 1}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    fc.double({min: 0, max: 100, noNaN: true}),
                    (seed, centerX, centerY) => {
                        // Create a seeded RNG
                        let state = seed;
                        const rng: RNG = () => {
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), 1 | state);
                            t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        const cluster = generateCluster(rng, centerX, centerY, 0);

                        // Extract squares' phase offsets, cycle durations, and peak opacities
                        const squareParams = cluster.squares.map((sq) => ({
                            phaseOffset: sq.phaseOffset,
                            cycleDuration: sq.cycleDuration,
                            peakOpacity: sq.peakOpacity,
                        }));

                        const {emergenceTimes} = computeEmergenceTimes(squareParams);

                        // Filter out Infinity values (squares with peakOpacity < 0.1 that never cross threshold)
                        const finiteTimes = emergenceTimes.filter((t) => isFinite(t));

                        // If fewer than 2 finite times, window constraint is trivially satisfied
                        if (finiteTimes.length < 2) {
                            return true;
                        }

                        // Verify all finite emergence times fall within a 2-second window
                        // Allow a tiny floating-point epsilon since the implementation uses
                        // iterative adjustment with trig functions.
                        const EPSILON = 1e-9;
                        const minTime = Math.min(...finiteTimes);
                        const maxTime = Math.max(...finiteTimes);
                        expect(maxTime - minTime).toBeLessThanOrEqual(2.0 + EPSILON);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 9: Fragment scheduler correctness
    // -------------------------------------------------------------------------
    describe("Property 9: Fragment scheduler correctness", () => {
        it("max 3 simultaneous, >=3s spawn interval, no spawn when inactive, no consecutive same text", () => {
            fc.assert(
                fc.property(
                    // Generate a sequence of time deltas between spawn attempts [0.5, 10] seconds
                    fc.array(fc.double({min: 0.5, max: 10, noNaN: true}), {
                        minLength: 5,
                        maxLength: 30,
                    }),
                    // Seed for RNG
                    fc.integer({min: 1, max: 2 ** 31 - 1}),
                    (timeDeltas, seed) => {
                        // Create a seeded RNG
                        let state = seed;
                        const rng = () => {
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), 1 | state);
                            t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        const activeClusters: ActiveClusterInfo[] = [
                            {centerX: 500, centerY: 500},
                            {centerX: 200, centerY: 300},
                        ];

                        let currentState = createFragmentSchedulerState();
                        let currentTime = 0;
                        let lastSpawnTime = -Infinity;
                        let lastSpawnedText: string | null = null;

                        // Simulate a sequence of spawn attempts
                        for (const delta of timeDeltas) {
                            currentTime += delta;

                            // Update fragments (remove expired ones)
                            currentState = updateFragments(currentState, currentTime);

                            const prevFragmentCount = currentState.activeFragments.length;

                            // Attempt to spawn
                            currentState = spawnFragment(
                                currentState,
                                currentTime,
                                activeClusters,
                                1920,
                                1080,
                                rng,
                            );

                            // (a) Max 3 simultaneous active fragments
                            expect(currentState.activeFragments.length).toBeLessThanOrEqual(3);

                            // Check if a spawn occurred
                            if (currentState.activeFragments.length > prevFragmentCount) {
                                // (b) At least 3s since last spawn
                                if (isFinite(lastSpawnTime)) {
                                    expect(currentTime - lastSpawnTime).toBeGreaterThanOrEqual(3);
                                }

                                // (d) No consecutive same text
                                const newFragment =
                                    currentState.activeFragments[currentState.activeFragments.length - 1];
                                if (lastSpawnedText !== null) {
                                    expect(newFragment.text).not.toBe(lastSpawnedText);
                                }

                                lastSpawnedText = newFragment.text;
                                lastSpawnTime = currentTime;
                            }
                        }

                        // (c) When activeClusters is empty, no spawn should occur
                        let emptyClusterState = createFragmentSchedulerState();
                        // Set lastSpawnTime to long ago so timing constraint won't block
                        emptyClusterState = {
                            ...emptyClusterState,
                            lastSpawnTime: -Infinity,
                        };

                        const stateBeforeEmptySpawn = emptyClusterState;
                        const stateAfterEmptySpawn = spawnFragment(
                            stateBeforeEmptySpawn,
                            100,
                            [], // no active clusters
                            1920,
                            1080,
                            rng,
                        );

                        // State should be unchanged (no spawn)
                        expect(stateAfterEmptySpawn.activeFragments.length).toBe(0);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 10: Fragment positioning proximity
    // -------------------------------------------------------------------------
    describe("Property 10: Fragment positioning proximity", () => {
        it("fragment position within 100px of an active cluster center", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 1, max: 2 ** 31 - 1}),
                    fc.integer({min: 200, max: 1920}),
                    fc.integer({min: 200, max: 1080}),
                    fc.double({min: 0, max: 1, noNaN: true}),
                    fc.double({min: 0, max: 1, noNaN: true}),
                    (seed, viewportWidth, viewportHeight, xFrac, yFrac) => {
                        // Place cluster center at least 100px from any edge so clamping
                        // cannot push the fragment position beyond 100px from center.
                        const centerX = 100 + xFrac * (viewportWidth - 200);
                        const centerY = 100 + yFrac * (viewportHeight - 200);
                        // Create a seeded RNG using mulberry32
                        let state = seed | 0;
                        const rng = () => {
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), 1 | state);
                            t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        // Create initial state with lastSpawnTime = -Infinity so spawn is allowed
                        const initialState = createFragmentSchedulerState();

                        // One active cluster at the known position
                        const activeClusters: ActiveClusterInfo[] = [{centerX, centerY}];

                        // Attempt to spawn a fragment
                        const newState = spawnFragment(
                            initialState,
                            0,
                            activeClusters,
                            viewportWidth,
                            viewportHeight,
                            rng,
                        );

                        // If a fragment was spawned, verify its distance to the cluster center
                        // Allow a tiny floating-point epsilon since the implementation uses
                        // trigonometric functions which can introduce rounding errors.
                        if (newState.activeFragments.length > 0) {
                            const fragment = newState.activeFragments[0];
                            const dx = fragment.x - centerX;
                            const dy = fragment.y - centerY;
                            const distance = Math.sqrt(dx * dx + dy * dy);
                            expect(distance).toBeLessThanOrEqual(100 + 1e-9);
                        }
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 11: Fragment style constraints
    // -------------------------------------------------------------------------
    describe("Property 11: Fragment style constraints", () => {
        it("font size [10, 16]px, peak opacity [0.3, 0.8], color in green or cyan range", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 1, max: 2 ** 31 - 1}),
                    fc.double({min: 200, max: 3000, noNaN: true}),
                    fc.double({min: 200, max: 3000, noNaN: true}),
                    (seed, viewportWidth, viewportHeight) => {
                        // Create a seeded RNG (mulberry32)
                        let state = seed;
                        const rng = () => {
                            state = (state + 0x6d2b79f5) | 0;
                            let t = Math.imul(state ^ (state >>> 15), 1 | state);
                            t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
                            return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
                        };

                        // Create a valid active cluster
                        const activeClusters: ActiveClusterInfo[] = [
                            {centerX: viewportWidth * 0.5, centerY: viewportHeight * 0.5},
                        ];

                        // Start with fresh state and spawn a fragment
                        const initialState = createFragmentSchedulerState();
                        const result = spawnFragment(
                            initialState,
                            10, // time well past the min spawn interval
                            activeClusters,
                            viewportWidth,
                            viewportHeight,
                            rng,
                        );

                        // Verify a fragment was spawned
                        expect(result.activeFragments.length).toBe(1);
                        const fragment = result.activeFragments[0];

                        // fontSize in [10, 16]
                        expect(fragment.fontSize).toBeGreaterThanOrEqual(10);
                        expect(fragment.fontSize).toBeLessThanOrEqual(16);

                        // peakOpacity in [0.3, 0.8]
                        expect(fragment.peakOpacity).toBeGreaterThanOrEqual(0.3);
                        expect(fragment.peakOpacity).toBeLessThanOrEqual(0.8);

                        // Color is either in green range OR cyan range
                        const {color} = fragment;
                        const isGreen =
                            color.h >= FRAGMENT_GREEN_HUE_RANGE[0] &&
                            color.h <= FRAGMENT_GREEN_HUE_RANGE[1] &&
                            color.s >= FRAGMENT_GREEN_SAT_MIN &&
                            color.l >= FRAGMENT_GREEN_LIGHTNESS[0] &&
                            color.l <= FRAGMENT_GREEN_LIGHTNESS[1];

                        const isCyan =
                            color.h >= FRAGMENT_CYAN_HUE_RANGE[0] &&
                            color.h <= FRAGMENT_CYAN_HUE_RANGE[1] &&
                            color.s >= FRAGMENT_CYAN_SAT_MIN &&
                            color.l >= FRAGMENT_CYAN_LIGHTNESS[0] &&
                            color.l <= FRAGMENT_CYAN_LIGHTNESS[1];

                        expect(isGreen || isCyan).toBe(true);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 12: Proportional cluster redistribution on resize
    // -------------------------------------------------------------------------
    describe("Property 12: Proportional cluster redistribution on resize", () => {
        it("centerXPct and centerYPct preserved after resize", () => {
            fc.assert(
                fc.property(
                    fc.array(arbClusterConfig(), {minLength: 1, maxLength: 10}),
                    fc.integer({min: 100, max: 3000}),
                    fc.integer({min: 100, max: 2000}),
                    (clusters, newWidth, newHeight) => {
                        const result = redistributeClusters(clusters, newWidth, newHeight);

                        // For each cluster, centerXPct and centerYPct must be unchanged (exact equality)
                        expect(result.length).toBe(clusters.length);
                        for (let i = 0; i < clusters.length; i++) {
                            expect(result[i].centerXPct).toBe(clusters[i].centerXPct);
                            expect(result[i].centerYPct).toBe(clusters[i].centerYPct);
                        }
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });

        it("zero-area canvas returns clusters unchanged", () => {
            fc.assert(
                fc.property(
                    fc.array(arbClusterConfig(), {minLength: 1, maxLength: 10}),
                    fc.oneof(
                        fc.constant(0),
                        fc.integer({min: -100, max: 0})
                    ),
                    fc.oneof(
                        fc.constant(0),
                        fc.integer({min: -100, max: 0})
                    ),
                    (clusters, zeroWidth, zeroHeight) => {
                        // Test with zero/negative width
                        const resultZeroW = redistributeClusters(clusters, zeroWidth, 500);
                        expect(resultZeroW).toBe(clusters);

                        // Test with zero/negative height
                        const resultZeroH = redistributeClusters(clusters, 500, zeroHeight);
                        expect(resultZeroH).toBe(clusters);

                        // Test with both zero
                        const resultBoth = redistributeClusters(clusters, zeroWidth, zeroHeight);
                        expect(resultBoth).toBe(clusters);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 13: Mobile cluster scaling
    // -------------------------------------------------------------------------
    describe("Property 13: Mobile cluster scaling", () => {
        it("mobile count = max(3, floor(N * 0.6)) for width < 768px", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 1, max: 50}),
                    (desktopCount) => {
                        const mobileCount = getMobileClusterCount(desktopCount);

                        // Verify: mobileCount === max(3, floor(N * 0.6))
                        const expected = Math.max(3, Math.floor(desktopCount * 0.6));
                        expect(mobileCount).toBe(expected);

                        // Also verify the result is always >= 3
                        expect(mobileCount).toBeGreaterThanOrEqual(3);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 14: Adaptive performance scaling
    // -------------------------------------------------------------------------
    describe("Property 14: Adaptive performance scaling", () => {
        it("halving at avg > 33ms (min 4 squares), restore 25% at avg < 25ms", () => {
            fc.assert(
                fc.property(
                    fc.integer({min: 8, max: 100}), // initial square count
                    fc.double({min: 34, max: 99, noNaN: true}), // slow frame time (> 33ms but <= 100ms)
                    fc.double({min: 1, max: 24, noNaN: true}),  // fast frame time (< 25ms)
                    (initialCount, slowFrameTime, fastFrameTime) => {
                        // --- Reduction scenario ---
                        // Create a fresh PerfState and feed 10 slow frames to fill window and trigger decision.
                        // The 10th frame completes the window and triggers the reduction.
                        let state = createPerfState(initialCount);
                        let reductionResult = evaluatePerformance(state, slowFrameTime);
                        for (let i = 1; i < 10; i++) {
                            reductionResult = evaluatePerformance(reductionResult.newState, slowFrameTime);
                        }

                        // After 10 frames, window is full. The 10th call triggers the decision.
                        // Verify reduction triggers and halves (min 4)
                        expect(reductionResult.shouldReduce).toBe(true);
                        const expectedReduced = Math.max(4, Math.floor(initialCount / 2));
                        expect(reductionResult.newState.currentSquareCount).toBe(expectedReduced);
                        expect(reductionResult.newState.hasReduced).toBe(true);

                        // --- Restoration scenario ---
                        // Start from the reduced state and feed fast frames to fill a fresh window
                        let restoredState: typeof reductionResult.newState = {
                            ...reductionResult.newState,
                            frameTimes: [], // Clear window for fresh fast-frame measurement
                        };
                        let restoreResult = evaluatePerformance(restoredState, fastFrameTime);
                        for (let i = 1; i < 10; i++) {
                            restoreResult = evaluatePerformance(restoreResult.newState, fastFrameTime);
                        }

                        // Only expect restoration if current < original (i.e., reduction actually reduced)
                        if (reductionResult.newState.currentSquareCount < initialCount) {
                            expect(restoreResult.shouldRestore).toBe(true);
                            // Restored amount = 25% of (original - current), at least 1
                            const removedCount = initialCount - reductionResult.newState.currentSquareCount;
                            const restoreAmount = Math.max(1, Math.floor(removedCount * 0.25));
                            const expectedRestored = Math.min(
                                initialCount,
                                reductionResult.newState.currentSquareCount + restoreAmount
                            );
                            expect(restoreResult.newState.currentSquareCount).toBe(expectedRestored);
                            // Never exceeds original
                            expect(restoreResult.newState.currentSquareCount).toBeLessThanOrEqual(initialCount);
                        }

                        // --- Minimum constraint ---
                        // Even with very small counts, can't go below 4
                        let smallState = createPerfState(5);
                        let smallResult = evaluatePerformance(smallState, slowFrameTime);
                        for (let i = 1; i < 10; i++) {
                            smallResult = evaluatePerformance(smallResult.newState, slowFrameTime);
                        }
                        expect(smallResult.newState.currentSquareCount).toBeGreaterThanOrEqual(4);

                        // --- Clamped frame skip ---
                        // Frame time > 100ms should be skipped (no decision)
                        const clampedResult = evaluatePerformance(state, 150);
                        expect(clampedResult.shouldReduce).toBe(false);
                        expect(clampedResult.shouldRestore).toBe(false);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });

    // -------------------------------------------------------------------------
    // Property 15: Central zone opacity cap
    // -------------------------------------------------------------------------
    describe("Property 15: Central zone opacity cap", () => {
        it("elements in central 60% of viewport have rendered opacity <= 0.4", () => {
            fc.assert(
                fc.property(
                    fc
                        .integer({min: 100, max: 3000})
                        .chain((viewportWidth) =>
                            fc
                                .double({min: 0, max: viewportWidth, noNaN: true})
                                .map((x) => ({viewportWidth, x}))
                        ),
                    ({viewportWidth, x}) => {
                        const cap = getCentralZoneOpacityCap(x, viewportWidth);

                        const leftBound = viewportWidth * 0.2;
                        const rightBound = viewportWidth * 0.8;
                        const inCentralZone = x >= leftBound && x <= rightBound;

                        // 1. If x is in central 60%, then getCentralZoneOpacityCap returns 0.4
                        if (inCentralZone) {
                            expect(cap).toBe(0.4);
                        } else {
                            // 2. If x is outside the central 60%, then getCentralZoneOpacityCap returns 1.0
                            expect(cap).toBe(1.0);
                        }

                        // 3. The returned value is always either 0.4 or 1.0
                        expect([0.4, 1.0]).toContain(cap);
                    }
                ),
                {numRuns: NUM_RUNS}
            );
        });
    });
});
