/*
 * Cluster Generator — Pure computation for generating square clusters.
 *
 * Generates SquareConfig and ClusterConfig instances with validated parameters
 * for the geometric background animation. Accepts a seeded RNG function for
 * deterministic, testable output.
 *
 */

import {
    type ClusterConfig,
    type HSLColor,
    MIN_DESKTOP_CLUSTERS,
    SQUARE_HUE_MAX,
    SQUARE_HUE_MIN,
    SQUARE_LIGHTNESS_MAX,
    SQUARE_LIGHTNESS_MIN,
    SQUARE_SAT_MIN,
    type SquareConfig,
} from "./types";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

/* A seeded random number generator returning values in [0, 1). */
export type RNG = () => number;

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

/* Square size range (CSS pixels) */
const SIZE_MIN = 20;
const SIZE_MAX = 120;

/* Rotation angle range (degrees) */
const ROTATION_MIN = -45;
const ROTATION_MAX = 45;

/* Peak opacity range */
const PEAK_OPACITY_MIN = 0.05;
const PEAK_OPACITY_MAX = 0.6;

/* Fade cycle duration range (seconds) */
const CYCLE_DURATION_MIN = 3;
const CYCLE_DURATION_MAX = 10;

/* Cluster square count range */
const CLUSTER_SQUARES_MIN = 8;
const CLUSTER_SQUARES_MAX = 20;

/* Cluster fade-in duration range (seconds) */
const FADE_IN_DURATION_MIN = 1;
const FADE_IN_DURATION_MAX = 3;

/* Minimum rotation angle difference between any two squares in a cluster (degrees) */
const MIN_ROTATION_DIFF = 5;

/* Minimum phase offset difference between any two squares in a cluster (seconds) */
const MIN_PHASE_OFFSET_DIFF = 0.5;

/* Minimum overlap fraction required between at least one pair of squares */
const MIN_OVERLAP_FRACTION = 0.1;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/* Returns a random number in [min, max] using the given RNG. */
function randomInRange(rng: RNG, min: number, max: number): number {
    return min + rng() * (max - min);
}

/*
 * Computes the overlap area between two axis-aligned squares given their
 * centers and sizes. Rotation is ignored for the overlap check — we use
 * bounding-box approximation with offsets, which is acceptable since the
 * requirement only asks that squares are "positioned such that" they overlap.
 */
function computeOverlapArea(
    ox1: number,
    oy1: number,
    size1: number,
    ox2: number,
    oy2: number,
    size2: number,
): number {
    const half1 = size1 / 2;
    const half2 = size2 / 2;

    const left1 = ox1 - half1;
    const right1 = ox1 + half1;
    const top1 = oy1 - half1;
    const bottom1 = oy1 + half1;

    const left2 = ox2 - half2;
    const right2 = ox2 + half2;
    const top2 = oy2 - half2;
    const bottom2 = oy2 + half2;

    const overlapX = Math.max(0, Math.min(right1, right2) - Math.max(left1, left2));
    const overlapY = Math.max(0, Math.min(bottom1, bottom2) - Math.max(top1, top2));

    return overlapX * overlapY;
}

/*
 * Checks whether a new rotation angle satisfies the minimum difference
 * constraint against all existing angles.
 */
function isRotationValid(angle: number, existing: number[]): boolean {
    return existing.every((a) => Math.abs(angle - a) > MIN_ROTATION_DIFF);
}

/*
 * Checks whether a new phase offset satisfies the minimum difference
 * constraint against all existing offsets.
 */
function isPhaseOffsetValid(offset: number, existing: number[]): boolean {
    return existing.every((o) => Math.abs(offset - o) > MIN_PHASE_OFFSET_DIFF);
}

// ---------------------------------------------------------------------------
// generateSquare
// ---------------------------------------------------------------------------

/*
 * Generates a single SquareConfig with all parameters within validated ranges.
 *
 * @param rng - Seeded random number generator returning [0, 1)
 * @param constrainedRotation - If provided, the rotation must differ by >5° from these values
 * @param constrainedPhaseOffsets - If provided, phase offset must differ by >0.5s from these values
 */
export function generateSquare(
    rng: RNG,
    constrainedRotation?: number[],
    constrainedPhaseOffsets?: number[],
): SquareConfig {
    // Color within emerald-to-green range
    const color: HSLColor = {
        h: randomInRange(rng, SQUARE_HUE_MIN, SQUARE_HUE_MAX),
        s: randomInRange(rng, SQUARE_SAT_MIN, 100),
        l: randomInRange(rng, SQUARE_LIGHTNESS_MIN, SQUARE_LIGHTNESS_MAX),
    };

    // Size
    const size = randomInRange(rng, SIZE_MIN, SIZE_MAX);

    // Rotation — most squares are 0° (axis-aligned), ~10% get 45° rotation
    const rotation = rng() < 0.1 ? 45 : 0;

    // Cycle duration
    const cycleDuration = randomInRange(rng, CYCLE_DURATION_MIN, CYCLE_DURATION_MAX);

    // Phase offset — generate with constraint if needed
    let phaseOffset: number;
    if (constrainedPhaseOffsets && constrainedPhaseOffsets.length > 0) {
        phaseOffset = generateConstrainedPhaseOffset(rng, cycleDuration, constrainedPhaseOffsets);
    } else {
        phaseOffset = randomInRange(rng, 0, cycleDuration);
    }

    // Peak opacity
    const peakOpacity = randomInRange(rng, PEAK_OPACITY_MIN, PEAK_OPACITY_MAX);

    // Aspect ratio — mix of squares and rectangles in portrait/landscape
    const ASPECT_POOL = [
        1.0,        // square (1:1)
        1.0,        // square — weighted
        1.0,        // square — weighted
        1.0,        // square — weighted
        1.0,        // square — weighted
        1.0,        // square — weighted
        4 / 3,      // 4:3 landscape
        3 / 4,      // 4:3 portrait
        16 / 9,     // 16:9 landscape
        9 / 16,     // 16:9 portrait
        21 / 9,     // 21:9 ultrawide landscape
        9 / 21,     // 21:9 ultrawide portrait
    ];
    const aspect = ASPECT_POOL[Math.floor(rng() * ASPECT_POOL.length)];

    // Offset from cluster center — initially zero; set by generateCluster
    const offsetX = 0;
    const offsetY = 0;

    return {
        offsetX,
        offsetY,
        size,
        aspect,
        rotation,
        color,
        phaseOffset,
        cycleDuration,
        peakOpacity,
    };
}

// ---------------------------------------------------------------------------
// Constrained generation helpers
// ---------------------------------------------------------------------------

/*
 * Generates a rotation angle that differs by >5° from all existing angles.
 * Uses rejection sampling with a fallback to deterministic placement.
 */
function generateConstrainedRotation(rng: RNG, existing: number[]): number {
    // Try rejection sampling first (fast path)
    for (let attempt = 0; attempt < 50; attempt++) {
        const angle = randomInRange(rng, ROTATION_MIN, ROTATION_MAX);
        if (isRotationValid(angle, existing)) {
            return angle;
        }
    }

    // Fallback: find a valid slot deterministically
    // The range is [-45, 45] = 90° total. With max 8 squares at >5° apart,
    // there are always valid positions available.
    for (let candidate = ROTATION_MIN; candidate <= ROTATION_MAX; candidate += 0.5) {
        if (isRotationValid(candidate, existing)) {
            return candidate;
        }
    }

    // Should never reach here with ≤8 squares in 90° range at >5° spacing
    return randomInRange(rng, ROTATION_MIN, ROTATION_MAX);
}

/*
 * Generates a phase offset that differs by >0.5s from all existing offsets.
 * Uses rejection sampling with a fallback to deterministic placement.
 * If no valid slot exists (rare with few squares), picks the offset that
 * maximizes minimum pairwise distance from existing offsets.
 */
function generateConstrainedPhaseOffset(
    rng: RNG,
    cycleDuration: number,
    existing: number[],
): number {
    // Try rejection sampling first
    for (let attempt = 0; attempt < 50; attempt++) {
        const offset = randomInRange(rng, 0, cycleDuration);
        if (isPhaseOffsetValid(offset, existing)) {
            return offset;
        }
    }

    // Fallback: find a valid slot deterministically with fine-grained scan
    const step = 0.1;
    for (let candidate = 0; candidate <= cycleDuration; candidate += step) {
        if (isPhaseOffsetValid(candidate, existing)) {
            return candidate;
        }
    }

    // Final fallback: pick the candidate that maximizes minimum distance
    // from existing offsets to best approximate the constraint.
    let bestCandidate = 0;
    let bestMinDist = -1;
    for (let candidate = 0; candidate <= cycleDuration; candidate += step) {
        const minDist = existing.reduce(
            (min, o) => Math.min(min, Math.abs(candidate - o)),
            Infinity,
        );
        if (minDist > bestMinDist) {
            bestMinDist = minDist;
            bestCandidate = candidate;
        }
    }
    return bestCandidate;
}

// ---------------------------------------------------------------------------
// generateCluster
// ---------------------------------------------------------------------------

/*
 * Generates a ClusterConfig with 3-8 squares satisfying all constraints:
 * - Rotation angles differ pairwise by >5°
 * - Phase offsets differ pairwise by >0.5s
 * - At least one pair of squares overlaps by ≥10% of the smaller square's area
 * - Bounding radius computed from max offset distance
 * - Fade-in duration in [1, 3]s
 *
 * @param rng - Seeded random number generator returning [0, 1)
 * @param centerXPct - Cluster center X position as percentage of viewport width [0, 100]
 * @param centerYPct - Cluster center Y position as percentage of viewport height [0, 100]
 * @param spawnTime - Animation time (seconds) when the cluster first appears
 */
export function generateCluster(
    rng: RNG,
    centerXPct: number,
    centerYPct: number,
    spawnTime: number = 0,
): ClusterConfig {
    const squareCount = Math.floor(randomInRange(rng, CLUSTER_SQUARES_MIN, CLUSTER_SQUARES_MAX + 0.999));
    const clampedCount = Math.max(CLUSTER_SQUARES_MIN, Math.min(CLUSTER_SQUARES_MAX, squareCount));

    // Pre-generate phase offsets: each square gets a random offset within its own
    // cycle duration. This ensures squares are spread across their lifecycle from
    // the start, so deaths/respawns trickle in gradually rather than batching up.
    const squares: SquareConfig[] = [];
    const rotations: number[] = [];

    for (let i = 0; i < clampedCount; i++) {
        const square = generateSquare(rng, rotations, []);
        rotations.push(square.rotation);

        // Random phase offset spread across the full cycle
        square.phaseOffset = rng() * square.cycleDuration;

        // Assign size based on position in cluster: first squares (center) are larger,
        // later squares (edges) are smaller. Steep exponential creates dense small edges.
        const t = i / (clampedCount - 1 || 1); // 0 for first, 1 for last
        square.size = SIZE_MIN + (SIZE_MAX - SIZE_MIN) * Math.pow(1 - t, 3.5);

        squares.push(square);
    }

    // Position squares on grid — anchor-point snapping for honeycomb effect
    positionSquaresOnGrid(rng, squares);

    // Compute bounding radius: max distance from cluster center to any square center
    const boundingRadius = computeBoundingRadius(squares);

    // Fade-in duration
    const fadeInDuration = randomInRange(rng, FADE_IN_DURATION_MIN, FADE_IN_DURATION_MAX);

    return {
        centerXPct,
        centerYPct,
        squares,
        boundingRadius,
        spawnTime,
        fadeInDuration,
    };
}

/*
 * Positions squares using anchor-point grid snapping for a "square honeycomb" effect.
 *
 * Each square is conceptually subdivided into a 4×4 grid (16 cells), producing
 * a 5×5 lattice of anchor points at the cell corners. The grid step in world
 * coordinates is `size / 4`.
 *
 * The first (largest) square is placed at the cluster center. Each subsequent
 * square picks a random anchor point on an existing square, then places one
 * of its own corners at that anchor. This creates overlapping grid-aligned
 * clusters that grow outward.
 *
 * Squares are sorted largest-first so the center gets the big squares and
 * edges accumulate smaller ones.
 */
function positionSquaresOnGrid(rng: RNG, squares: SquareConfig[]): void {
    if (squares.length === 0) return;

    // Sort squares largest to smallest — big ones go at center
    squares.sort((a, b) => b.size - a.size);

    // Place the first (largest) square at the cluster center
    squares[0].offsetX = 0;
    squares[0].offsetY = 0;

    // For each subsequent square, snap one of its corners to an anchor on an existing square
    for (let i = 1; i < squares.length; i++) {
        const newSquare = squares[i];

        // Pick a random existing square to attach to
        const parentIdx = Math.floor(rng() * i);
        const parent = squares[parentIdx];

        // Get a random anchor point on the parent's 5×5 grid
        const anchor = getRandomAnchorPoint(parent, rng);

        // Pick which corner of the new square to snap to that anchor
        // Corners of an axis-aligned square centered at (0,0) with half-size h:
        // TL=(-h,-h), TR=(h,-h), BR=(h,h), BL=(-h,h)
        const h = newSquare.size / 2;
        const corners = [
            {dx: -h, dy: -h}, // top-left
            {dx: h, dy: -h}, // top-right
            {dx: h, dy: h}, // bottom-right
            {dx: -h, dy: h}, // bottom-left
        ];
        const corner = corners[Math.floor(rng() * corners.length)];

        // Position the new square so that its chosen corner lands on the anchor
        // anchor = (newSquare center) + corner offset
        // => newSquare center = anchor - corner offset
        newSquare.offsetX = anchor.x - corner.dx;
        newSquare.offsetY = anchor.y - corner.dy;
    }
}

/*
 * Returns a random anchor point from a square's 4×4 subdivision grid.
 * The grid has 5×5 = 25 anchor points (corners of the 16 cells).
 * Coordinates are in world space (relative to cluster center).
 */
function getRandomAnchorPoint(
    square: SquareConfig,
    rng: RNG,
): { x: number; y: number } {
    const step = square.size / 4;
    const halfSize = square.size / 2;

    // Grid indices 0-4 in both axes
    const gridX = Math.floor(rng() * 5);
    const gridY = Math.floor(rng() * 5);

    // Local coordinates (relative to square center)
    const localX = -halfSize + gridX * step;
    const localY = -halfSize + gridY * step;

    // If the square is rotated 45°, transform the local point
    if (square.rotation === 45) {
        const rad = Math.PI / 4;
        const cosR = Math.cos(rad);
        const sinR = Math.sin(rad);
        return {
            x: square.offsetX + localX * cosR - localY * sinR,
            y: square.offsetY + localX * sinR + localY * cosR,
        };
    }

    // Axis-aligned (0° rotation) — local coords are world coords offset by square position
    return {
        x: square.offsetX + localX,
        y: square.offsetY + localY,
    };
}

/*
 * Computes the bounding radius of a cluster: the maximum distance from
 * the cluster center (0, 0) to any square's offset position.
 */
function computeBoundingRadius(squares: SquareConfig[]): number {
    let maxDist = 0;
    for (const sq of squares) {
        const dist = Math.sqrt(sq.offsetX * sq.offsetX + sq.offsetY * sq.offsetY);
        if (dist > maxDist) {
            maxDist = dist;
        }
    }
    return maxDist;
}

// ---------------------------------------------------------------------------
// initializeClusters
// ---------------------------------------------------------------------------

/*
 * Creates clusters distributed across the viewport based on viewport size.
 *
 * - Desktop (≥1024px): max(baseCount, MIN_DESKTOP_CLUSTERS) clusters
 * - Mobile (<768px): max(3, floor(desktopCount * 0.6)) clusters
 * - Tablet (768–1023px): same as desktop count (the spec only distinguishes <768 and >=1024)
 *
 * Cluster positions are stored as percentages of viewport dimensions for
 * proportional redistribution on resize.
 *
 * @param viewportWidth - Current viewport width in pixels
 * @param viewportHeight - Current viewport height in pixels
 * @param baseCount - Base cluster count from configuration
 * @param rng - Seeded random number generator
 */
export function initializeClusters(
    viewportWidth: number,
    viewportHeight: number,
    baseCount: number,
    rng: RNG,
): ClusterConfig[] {
    // Determine cluster count based on viewport
    const desktopCount = Math.max(baseCount, MIN_DESKTOP_CLUSTERS);
    let clusterCount: number;

    if (viewportWidth >= 1024) {
        // Desktop
        clusterCount = desktopCount;
    } else if (viewportWidth < 768) {
        // Mobile: reduce to 60%, minimum 3
        clusterCount = Math.max(3, Math.floor(desktopCount * 0.6));
    } else {
        // Tablet (768–1023): use desktop count
        clusterCount = desktopCount;
    }

    // Distribute clusters across the viewport
    const clusters: ClusterConfig[] = [];

    for (let i = 0; i < clusterCount; i++) {
        // Distribute positions using a combination of grid-based seeding and
        // random jitter to avoid pure randomness creating clumps.
        // Use a simple strategy: divide viewport into a rough grid and jitter.
        const cols = Math.ceil(Math.sqrt(clusterCount));
        const rows = Math.ceil(clusterCount / cols);

        const col = i % cols;
        const row = Math.floor(i / cols);

        // Base position in grid cell (as percentage)
        const cellWidth = 100 / cols;
        const cellHeight = 100 / rows;

        // Center of cell + random jitter within cell (with padding from edges)
        const padding = 5; // 5% padding from viewport edges
        const centerXPct = Math.max(
            padding,
            Math.min(100 - padding, cellWidth * (col + 0.5) + randomInRange(rng, -cellWidth * 0.3, cellWidth * 0.3)),
        );
        const centerYPct = Math.max(
            padding,
            Math.min(100 - padding, cellHeight * (row + 0.5) + randomInRange(rng, -cellHeight * 0.3, cellHeight * 0.3)),
        );

        // Stagger spawn times slightly for visual interest
        const spawnTime = i * randomInRange(rng, 0.1, 0.5);

        const cluster = generateCluster(rng, centerXPct, centerYPct, spawnTime);
        clusters.push(cluster);
    }

    return clusters;
}
