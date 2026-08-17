/*
 * Geometric Background Animation — Glow computation utilities.
 *
 * Pure functions for computing the cluster glow opacity multiplier
 * and the glow layer color for individual squares.
 *
 */

import type {HSLColor} from "./types";

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

/* Minimum glow spread in pixels (Requirement 1.7) */
export const MIN_GLOW_SPREAD_PX = 4;

// ---------------------------------------------------------------------------
// Glow multiplier (Requirement 1.6)
// ---------------------------------------------------------------------------

/*
 * Computes the opacity multiplier for a square based on its distance from
 * the cluster center. Linearly interpolates from 1.0 at the center to 0.6
 * at the bounding radius.
 *
 * Formula: multiplier = 1.0 - 0.4 * (d / R)
 *
 * @param distanceFromCenter - Distance from the cluster center (d), in pixels.
 * @param boundingRadius - The cluster's bounding radius (R), in pixels.
 * @returns A multiplier clamped to [0.6, 1.0].
 */
export function computeGlowMultiplier(
    distanceFromCenter: number,
    boundingRadius: number,
): number {
    if (boundingRadius <= 0) {
        return 1.0;
    }

    const ratio = distanceFromCenter / boundingRadius;
    const multiplier = 1.0 - 0.4 * ratio;

    // Clamp to [0.6, 1.0]
    return Math.max(0.6, Math.min(1.0, multiplier));
}

// ---------------------------------------------------------------------------
// Glow color (Requirement 1.7)
// ---------------------------------------------------------------------------

/*
 * Computes the glow layer color for a square. The glow color has the same
 * hue and saturation as the base color but with lightness boosted by at least
 * 20 percentage points (capped at 100).
 *
 * @param baseColor - The square's base HSL color.
 * @returns The HSL color for the glow layer.
 */
export function computeGlowColor(baseColor: HSLColor): HSLColor {
    return {
        h: baseColor.h,
        s: baseColor.s,
        l: Math.min(baseColor.l + 20, 100),
    };
}
