/*
 * Geometric Background Animation — Cluster redistribution on resize.
 *
 * When the canvas dimensions change, cluster positions (stored as percentages)
 * are preserved. Since ClusterConfig uses centerXPct/centerYPct, the actual
 * pixel positions are computed at render time from these percentages and the
 * current canvas dimensions.
 *
 */

import type {ClusterConfig} from "./types";

/*
 * Redistributes clusters for new canvas dimensions.
 *
 * Since ClusterConfig stores positions as viewport percentages (centerXPct,
 * centerYPct), proportional positioning is inherently preserved. This function
 * validates the new dimensions and returns the clusters unchanged.
 *
 * Edge case: if the new canvas has zero area (width <= 0 or height <= 0),
 * redistribution is skipped and the clusters are returned as-is.
 *
 * @param clusters - Existing cluster configurations
 * @param newWidth - New canvas width in pixels
 * @param newHeight - New canvas height in pixels
 * @returns The clusters with percentage-based positions preserved
 */
export function redistributeClusters(
    clusters: ClusterConfig[],
    newWidth: number,
    newHeight: number
): ClusterConfig[] {
    // Zero-area canvas: skip redistribution, return clusters unchanged
    if (newWidth <= 0 || newHeight <= 0) {
        return clusters;
    }

    // ClusterConfig stores positions as percentages (centerXPct, centerYPct).
    // Percentage-based positions are inherently preserved across resizes —
    // the pixel conversion happens at render time.
    return clusters;
}
