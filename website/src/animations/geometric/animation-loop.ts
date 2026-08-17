/*
 * Geometric Background Animation — Render loop and state management.
 *
 * Manages the requestAnimationFrame loop with:
 * - Time tracking with delta clamped to 100ms (prevents jumps from background tabs)
 * - Deferred initialization: clusters created over first 3 frames (~1/3 each)
 * - Frame-by-frame wiring of cluster fade cycles, fragment scheduler, and performance monitor
 * - Adaptive scaling (reduce/restore clusters based on performance decisions)
 *
 */

import type {AnimationState, ClusterConfig,} from './types';

import {generateSquare, initializeClusters, type RNG} from './cluster-generator';
import {redistributeClusters} from './cluster-redistribution';
import {
  type ActiveClusterInfo,
  createFragmentSchedulerState,
  spawnFragment,
  updateFragments,
} from './fragment-scheduler';
import {computeFadeOpacity} from './fade-cycle';
import {createPerfState, evaluatePerformance,} from './performance';
import {renderFrame} from './renderer';

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

/* Maximum delta time per frame in milliseconds (prevents opacity jumps) */
const MAX_DELTA_MS = 100;

/* Number of frames over which deferred initialization is spread */
const DEFERRED_INIT_FRAMES = 3;

/* Opacity threshold for considering a square "active" */
const ACTIVE_OPACITY_THRESHOLD = 0.1;

// ---------------------------------------------------------------------------
// Animation Loop Class
// ---------------------------------------------------------------------------

export interface AnimationLoopConfig {
    /* Canvas width in pixels */
    width: number;
    /* Canvas height in pixels */
    height: number;
    /* Base cluster count from configuration */
    baseClusterCount: number;
    /* Whether reduced motion is initially active */
    reducedMotion?: boolean;
    /* Injectable RNG for testing (defaults to Math.random) */
    rng?: RNG;
}

/*
 * Creates and manages the animation render loop.
 *
 * Exposes methods to start, stop, resize, and control reduced motion.
 * Uses deferred initialization to spread cluster creation across
 * the first 3 frames, keeping each frame's init work under 16ms.
 */
export class AnimationLoop {
    private state: AnimationState;
    private lastTimestamp: number | null = null;
    private rng: RNG;
    private baseClusterCount: number;
    private ctx: CanvasRenderingContext2D | null = null;

    /* Viewport width (for central zone cap and fragment positioning) */
    private viewportWidth: number;

    /* Total cluster count determined at init (before deferred batching) */
    private pendingClusters: ClusterConfig[] = [];

    /* Stashed clusters removed by performance scaling (for restoration) */
    private stashedClusters: ClusterConfig[] = [];

    constructor(config: AnimationLoopConfig) {
        this.rng = config.rng ?? Math.random;
        this.baseClusterCount = config.baseClusterCount;
        this.viewportWidth = config.width;

        this.state = {
            running: false,
            reducedMotion: config.reducedMotion ?? false,
            width: config.width,
            height: config.height,
            time: 0,
            clusters: [],
            fragmentState: createFragmentSchedulerState(),
            perfState: createPerfState(0),
            frameCount: 0,
            animationFrameId: null,
        };
    }

    // -------------------------------------------------------------------------
    // Public API
    // -------------------------------------------------------------------------

    /*
     * Starts the animation loop.
     * If already running or reducedMotion is active, this is a no-op.
     */
    startLoop(ctx: CanvasRenderingContext2D): void {
        if (this.state.running || this.state.reducedMotion) {
            return;
        }

        this.ctx = ctx;
        this.state.running = true;
        this.lastTimestamp = null;

        // Pre-generate all clusters for deferred initialization
        this.pendingClusters = initializeClusters(
            this.state.width,
            this.state.height,
            this.baseClusterCount,
            this.rng,
        );

        // Reset frame count and time for fresh deferred init
        this.state.frameCount = 0;
        this.state.time = 0;
        this.state.clusters = [];
        this.stashedClusters = [];

        // Compute total squares for performance tracking
        const totalSquares = this.pendingClusters.reduce(
            (sum, c) => sum + c.squares.length,
            0,
        );
        this.state.perfState = createPerfState(totalSquares);

        this.scheduleFrame();
    }

    /*
     * Stops the animation loop and cancels the pending animation frame.
     * Idempotent — safe to call when already stopped.
     */
    stopLoop(): void {
        this.state.running = false;

        if (this.state.animationFrameId !== null) {
            cancelAnimationFrame(this.state.animationFrameId);
            this.state.animationFrameId = null;
        }

        this.lastTimestamp = null;
        this.ctx = null;
    }

    /*
     * Handles canvas resize. Updates dimensions and redistributes clusters.
     */
    resize(width: number, height: number): void {
        this.state.width = width;
        this.state.height = height;
        this.viewportWidth = width;

        // Redistribute existing clusters for new dimensions
        this.state.clusters = redistributeClusters(
            this.state.clusters,
            width,
            height,
        );
    }

    /*
     * Sets the reduced motion preference.
     * When true, stops the loop immediately.
     * When false, allows the next startLoop() call to work.
     */
    setReducedMotion(reduced: boolean): void {
        this.state.reducedMotion = reduced;
        if (reduced && this.state.running) {
            this.stopLoop();
        }
    }

    /* Returns the current animation state (for testing/inspection). */
    getState(): Readonly<AnimationState> {
        return this.state;
    }

    // -------------------------------------------------------------------------
    // Frame loop internals
    // -------------------------------------------------------------------------

    /*
     * Schedules the next animation frame.
     */
    private scheduleFrame(): void {
        this.state.animationFrameId = requestAnimationFrame((ts) => this.tick(ts));
    }

    /*
     * Main frame tick — called once per animation frame.
     *
     * @param timestamp - High-resolution timestamp from requestAnimationFrame
     */
    private tick(timestamp: number): void {
        if (!this.state.running || !this.ctx) {
            return;
        }

        // -----------------------------------------------------------------------
        // Step 1: Compute delta time, clamped to MAX_DELTA_MS
        // -----------------------------------------------------------------------
        let dtMs: number;
        if (this.lastTimestamp === null) {
            dtMs = 16; // Assume ~60fps for first frame
        } else {
            dtMs = timestamp - this.lastTimestamp;
        }
        dtMs = Math.min(dtMs, MAX_DELTA_MS);
        this.lastTimestamp = timestamp;

        const dtSec = dtMs / 1000;
        this.state.time += dtSec * 2; // 2x speed multiplier

        // -----------------------------------------------------------------------
        // Step 2: Deferred initialization — spread cluster creation over 3 frames
        // -----------------------------------------------------------------------
        if (this.state.frameCount < DEFERRED_INIT_FRAMES && this.pendingClusters.length > 0) {
            const remaining = DEFERRED_INIT_FRAMES - this.state.frameCount;
            const batchSize = Math.ceil(this.pendingClusters.length / remaining);
            const batch = this.pendingClusters.splice(0, batchSize);
            this.state.clusters.push(...batch);

            // Update perfState square count to reflect currently initialized clusters
            const currentSquareCount = this.state.clusters.reduce(
                (sum, c) => sum + c.squares.length,
                0,
            );
            this.state.perfState = {
                ...this.state.perfState,
                currentSquareCount,
                originalSquareCount: this.state.perfState.originalSquareCount,
            };
        }

        // -----------------------------------------------------------------------
        // Step 3: Determine active clusters (at least one square opacity >= 0.1)
        // -----------------------------------------------------------------------
        const activeClusters: ActiveClusterInfo[] = [];

        for (const cluster of this.state.clusters) {
            // Check spawn fade-in factor
            const elapsed = this.state.time - cluster.spawnTime;
            if (elapsed <= 0) continue;

            const fadeInFactor = elapsed >= cluster.fadeInDuration
                ? 1
                : elapsed / cluster.fadeInDuration;

            let hasActiveSquare = false;
            for (const square of cluster.squares) {
                const baseOpacity = computeFadeOpacity(
                    this.state.time,
                    square.phaseOffset,
                    square.cycleDuration,
                    square.peakOpacity,
                );
                const effectiveOpacity = baseOpacity * fadeInFactor;
                if (effectiveOpacity >= ACTIVE_OPACITY_THRESHOLD) {
                    hasActiveSquare = true;
                    break;
                }
            }

            if (hasActiveSquare) {
                activeClusters.push({
                    centerX: (cluster.centerXPct / 100) * this.state.width,
                    centerY: (cluster.centerYPct / 100) * this.state.height,
                });
            }
        }

        // -----------------------------------------------------------------------
        // Step 3b: Square lifecycle — remove completed squares, spawn replacements
        // A square is "done" when it has lived through at least one full cycle
        // and its opacity is back at (or near) 0.
        // Limit replacements per frame to avoid visual bursts.
        // -----------------------------------------------------------------------
        let replacementsThisFrame = 0;
        const MAX_REPLACEMENTS_PER_FRAME = 1;

        for (const cluster of this.state.clusters) {
            if (replacementsThisFrame >= MAX_REPLACEMENTS_PER_FRAME) break;

            const elapsed = this.state.time - cluster.spawnTime;
            if (elapsed <= cluster.fadeInDuration) continue; // cluster still fading in

            for (let i = cluster.squares.length - 1; i >= 0; i--) {
                if (replacementsThisFrame >= MAX_REPLACEMENTS_PER_FRAME) break;

                const sq = cluster.squares[i];
                // Square has lived at least one full cycle
                const squareAge = this.state.time - cluster.spawnTime;
                if (squareAge < sq.cycleDuration) continue;

                const opacity = computeFadeOpacity(
                    this.state.time,
                    sq.phaseOffset,
                    sq.cycleDuration,
                    sq.peakOpacity,
                );

                // If opacity is near zero after at least one cycle, it's done
                if (opacity < 0.01) {
                    // Remember the dead square's position before removing
                    const deadX = sq.offsetX;
                    const deadY = sq.offsetY;

                    // Remove it
                    cluster.squares.splice(i, 1);

                    const newSquare = generateSquare(this.rng);
                    // Varied size for visual freshness
                    newSquare.size = 15 + this.rng() * 40;

                    // Position it on the grid of a random remaining square,
                    // but reject positions too close to where the dead square was.
                    if (cluster.squares.length > 0) {
                        let placed = false;
                        for (let attempt = 0; attempt < 10; attempt++) {
                            const parent = cluster.squares[Math.floor(this.rng() * cluster.squares.length)];
                            const step = parent.size / 4;
                            const halfSize = parent.size / 2;
                            const gridX = Math.floor(this.rng() * 5);
                            const gridY = Math.floor(this.rng() * 5);
                            const localX = -halfSize + gridX * step;
                            const localY = -halfSize + gridY * step;
                            const h = newSquare.size / 2;
                            const corners = [
                                {dx: -h, dy: -h},
                                {dx: h, dy: -h},
                                {dx: h, dy: h},
                                {dx: -h, dy: h},
                            ];
                            const corner = corners[Math.floor(this.rng() * 4)];
                            const candidateX = parent.offsetX + localX - corner.dx;
                            const candidateY = parent.offsetY + localY - corner.dy;

                            // Reject if too close to the dead square's position
                            const dist = Math.sqrt((candidateX - deadX) ** 2 + (candidateY - deadY) ** 2);
                            if (dist > newSquare.size * 0.5) {
                                newSquare.offsetX = candidateX;
                                newSquare.offsetY = candidateY;
                                placed = true;
                                break;
                            }
                        }
                        // Fallback: pick an outward position from cluster center
                        if (!placed) {
                            const angle = this.rng() * Math.PI * 2;
                            const radius = 60 + this.rng() * 40;
                            newSquare.offsetX = Math.cos(angle) * radius;
                            newSquare.offsetY = Math.sin(angle) * radius;
                        }
                    }

                    // Random phase offset so it fades in fresh
                    newSquare.phaseOffset = this.rng() * newSquare.cycleDuration;
                    cluster.squares.push(newSquare);
                    replacementsThisFrame++;
                }
            }
        }

        // -----------------------------------------------------------------------
        // Step 4: Fragment scheduler — spawn and update
        // -----------------------------------------------------------------------
        this.state.fragmentState = spawnFragment(
            this.state.fragmentState,
            this.state.time,
            activeClusters,
            this.state.width,
            this.state.height,
            this.rng,
        );

        this.state.fragmentState = updateFragments(
            this.state.fragmentState,
            this.state.time,
        );

        // -----------------------------------------------------------------------
        // Step 5: Performance evaluation
        // -----------------------------------------------------------------------
        const perfDecision = evaluatePerformance(this.state.perfState, dtMs);
        this.state.perfState = perfDecision.newState;

        // Step 6: Apply performance scaling decisions
        if (perfDecision.shouldReduce) {
            // Remove clusters from the end of the array
            const targetSquares = perfDecision.newState.currentSquareCount;
            this.reduceClustersBySquareCount(targetSquares);
        } else if (perfDecision.shouldRestore) {
            // Restore clusters from the stash
            const targetSquares = perfDecision.newState.currentSquareCount;
            this.restoreClustersBySquareCount(targetSquares);
        }

        // -----------------------------------------------------------------------
        // Step 7: Render the frame
        // -----------------------------------------------------------------------
        renderFrame(
            this.ctx,
            this.state.clusters,
            this.state.fragmentState.activeFragments,
            this.state.time,
            this.viewportWidth,
            this.state.width,
            this.state.height,
        );

        // -----------------------------------------------------------------------
        // Step 8: Increment frame count and schedule next frame
        // -----------------------------------------------------------------------
        this.state.frameCount++;
        this.scheduleFrame();
    }

    // -------------------------------------------------------------------------
    // Performance scaling helpers
    // -------------------------------------------------------------------------

    /*
     * Removes clusters from the end of the active list until the total
     * square count is at or below the target. Stashes removed clusters
     * for potential restoration later.
     */
    private reduceClustersBySquareCount(targetSquares: number): void {
        let currentSquares = this.state.clusters.reduce(
            (sum, c) => sum + c.squares.length,
            0,
        );

        while (currentSquares > targetSquares && this.state.clusters.length > 1) {
            const removed = this.state.clusters.pop()!;
            this.stashedClusters.push(removed);
            currentSquares -= removed.squares.length;
        }
    }

    /*
     * Restores clusters from the stash until the total square count
     * reaches the target (or the stash is empty).
     */
    private restoreClustersBySquareCount(targetSquares: number): void {
        let currentSquares = this.state.clusters.reduce(
            (sum, c) => sum + c.squares.length,
            0,
        );

        while (
            currentSquares < targetSquares &&
            this.stashedClusters.length > 0
            ) {
            const restored = this.stashedClusters.pop()!;
            if (currentSquares + restored.squares.length <= targetSquares) {
                this.state.clusters.push(restored);
                currentSquares += restored.squares.length;
            } else {
                // Put it back if it would exceed target
                this.stashedClusters.push(restored);
                break;
            }
        }
    }
}
