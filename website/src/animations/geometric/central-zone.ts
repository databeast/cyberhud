/*
 * Geometric Background Animation — Central Zone Opacity Cap.
 *
 * Limits the maximum rendered opacity of elements (squares and fragments)
 * that fall within the central 60% of the viewport width, ensuring overlaid
 * text (project name, tagline, CTA) remains legible.
 *
 */

/*
 * Returns the effective opacity cap for an element at horizontal position `x`.
 *
 * The central zone is defined as the middle 60% of the viewport width,
 * i.e., x values between 20% and 80% of `viewportWidth`.
 *
 * - Elements inside the central zone are capped at 0.4 opacity.
 * - Elements outside the central zone are uncapped (returns 1.0).
 *
 * @param x - The horizontal position of the element in pixels.
 * @param viewportWidth - The total viewport width in pixels.
 * @returns 0.4 if the element is within the central 60%, 1.0 otherwise.
 */
export function getCentralZoneOpacityCap(x: number, viewportWidth: number): number {
    const leftBound = viewportWidth * 0.2;
    const rightBound = viewportWidth * 0.8;

    if (x >= leftBound && x <= rightBound) {
        return 0.4;
    }

    return 1.0;
}
