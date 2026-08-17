/*
 * Geometric Background Animation — Fade Cycle Calculator.
 *
 * Pure function computing the base opacity of a square at a given time using
 * sinusoidal easing. The function guarantees output is always clamped to
 * [0, peakOpacity] and produces smooth transitions (no abrupt jumps).
 *
 */

/*
 * Computes the base opacity of a square at time `t`.
 *
 * Uses sinusoidal easing:
 *   opacity = peakOpacity * (0.5 + 0.5 * sin(2π * (t + phaseOffset) / cycleDuration - π/2))
 *
 * This guarantees:
 * - Minimum output: 0 (when sin term equals -1)
 * - Maximum output: peakOpacity (when sin term equals 1)
 * - Smooth transitions with bounded frame-to-frame change
 *
 * @param time - Current animation time in seconds
 * @param phaseOffset - Phase offset in seconds [0, cycleDuration]
 * @param cycleDuration - Full cycle duration in seconds [3, 10]
 * @param peakOpacity - Maximum opacity for this square [0.05, 0.6]
 * @returns Computed opacity in [0, peakOpacity]
 */
export function computeFadeOpacity(
    time: number,
    phaseOffset: number,
    cycleDuration: number,
    peakOpacity: number,
): number {
    const TWO_PI = 2 * Math.PI;
    const HALF_PI = Math.PI / 2;

    const phase = TWO_PI * (time + phaseOffset) / cycleDuration - HALF_PI;
    const raw = peakOpacity * (0.5 + 0.5 * Math.sin(phase));

    // Clamp to [0, peakOpacity] to guard against floating-point edge cases
    return Math.max(0, Math.min(peakOpacity, raw));
}

/*
 * Computes the first time each square's opacity crosses 0.1 rising from 0,
 * and validates that all emergence times fall within a 2-second window.
 * If not, adjusts phase offsets to satisfy the constraint.
 *
 * The emergence time is derived from the sinusoidal fade formula:
 *   opacity = peakOpacity * (0.5 + 0.5 * sin(2π * (t + phaseOffset) / cycleDuration - π/2))
 *
 * Setting opacity = 0.1 and solving for the ascending crossing gives:
 *   t = cycleDuration * (arcsin((0.2 / peakOpacity) - 1) + π/2) / (2π) - phaseOffset
 *
 *
 * @param squares - Array of squares with their phase offsets, cycle durations, and peak opacities
 * @returns Object with emergence times for each square and whether adjustments were made
 */
export function computeEmergenceTimes(
    squares: Array<{ phaseOffset: number; cycleDuration: number; peakOpacity: number }>,
): { emergenceTimes: number[]; adjusted: boolean } {
    const EMERGENCE_THRESHOLD = 0.1;
    const MAX_WINDOW = 2.0; // seconds
    const TWO_PI = 2 * Math.PI;
    const HALF_PI = Math.PI / 2;

    /*
     * Computes the first t >= 0 where opacity crosses EMERGENCE_THRESHOLD rising from 0
     * for a square with the given parameters.
     *
     * Returns Infinity if peakOpacity < EMERGENCE_THRESHOLD (square never reaches threshold).
     */
    function getEmergenceTime(
        phaseOffset: number,
        cycleDuration: number,
        peakOpacity: number,
    ): number {
        // If peak opacity is below threshold, the square can never emerge
        if (peakOpacity < EMERGENCE_THRESHOLD) {
            return Infinity;
        }

        // Solve: peakOpacity * (0.5 + 0.5 * sin(θ)) = EMERGENCE_THRESHOLD
        // sin(θ) = (2 * EMERGENCE_THRESHOLD / peakOpacity) - 1
        const sinValue = (2 * EMERGENCE_THRESHOLD / peakOpacity) - 1;

        // Clamp for floating-point safety (sinValue should be in [-1, 1])
        const clampedSin = Math.max(-1, Math.min(1, sinValue));

        // The ascending crossing occurs at θ = arcsin(clampedSin)
        // where θ = 2π * (t + phaseOffset) / cycleDuration - π/2
        const theta = Math.asin(clampedSin);

        // Solve for t: t = cycleDuration * (theta + π/2) / (2π) - phaseOffset
        let t = cycleDuration * (theta + HALF_PI) / TWO_PI - phaseOffset;

        // Normalize t to be >= 0 by adding full cycles
        while (t < 0) {
            t += cycleDuration;
        }

        return t;
    }

    // Compute emergence times for all squares
    const emergenceTimes = squares.map((sq) =>
        getEmergenceTime(sq.phaseOffset, sq.cycleDuration, sq.peakOpacity),
    );

    // Filter out Infinity values for window check (squares that can never emerge)
    const finiteTimes = emergenceTimes.filter((t) => isFinite(t));

    // If fewer than 2 finite times, no window constraint to enforce
    if (finiteTimes.length < 2) {
        return {emergenceTimes, adjusted: false};
    }

    const minTime = Math.min(...finiteTimes);
    const maxTime = Math.max(...finiteTimes);

    // Check if already within the 2-second window
    if (maxTime - minTime <= MAX_WINDOW) {
        return {emergenceTimes, adjusted: false};
    }

    // Need to adjust phase offsets to bring emergence times within window.
    // Strategy: shift phase offsets so that all emergence times cluster around
    // the median emergence time within a 2-second window.
    const targetCenter = (minTime + maxTime) / 2;
    const targetStart = targetCenter - MAX_WINDOW / 2;
    const targetEnd = targetCenter + MAX_WINDOW / 2;

    for (let i = 0; i < squares.length; i++) {
        if (!isFinite(emergenceTimes[i])) {
            continue;
        }

        const currentEmergence = emergenceTimes[i];

        if (currentEmergence < targetStart || currentEmergence > targetEnd) {
            // Compute the desired emergence time (clamp to window)
            const desiredEmergence = Math.max(targetStart, Math.min(targetEnd, currentEmergence));

            // Reverse-engineer the phase offset needed to achieve this emergence time.
            // From: t = cycleDuration * (theta + π/2) / (2π) - phaseOffset
            // So: phaseOffset = cycleDuration * (theta + π/2) / (2π) - t
            const sq = squares[i];
            const sinValue = (2 * EMERGENCE_THRESHOLD / sq.peakOpacity) - 1;
            const clampedSin = Math.max(-1, Math.min(1, sinValue));
            const theta = Math.asin(clampedSin);

            let newPhaseOffset = sq.cycleDuration * (theta + HALF_PI) / TWO_PI - desiredEmergence;

            // Normalize phase offset to [0, cycleDuration)
            while (newPhaseOffset < 0) {
                newPhaseOffset += sq.cycleDuration;
            }
            while (newPhaseOffset >= sq.cycleDuration) {
                newPhaseOffset -= sq.cycleDuration;
            }

            // Mutate the square's phase offset
            sq.phaseOffset = newPhaseOffset;

            // Update the emergence time
            emergenceTimes[i] = desiredEmergence;
        }
    }

    return {emergenceTimes, adjusted: true};
}
