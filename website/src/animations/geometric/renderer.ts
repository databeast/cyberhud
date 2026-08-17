/*
 * Geometric Background Animation — Canvas 2D Rendering Functions.
 *
 * Responsible for all drawing operations: squares with rotation,
 * glow layers with blur, pseudocode text fragments, and per-frame
 * orchestration with correct draw order.
 *
 */

import type {ActiveFragment, ClusterConfig, HSLColor, SquareConfig} from './types';
import {computeGlowColor, computeGlowMultiplier, MIN_GLOW_SPREAD_PX} from './glow';
import {computeFadeOpacity} from './fade-cycle';
import {getCentralZoneOpacityCap} from './central-zone';
import {computeFragmentOpacity} from './fragment-scheduler';

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/* Converts degrees to radians. */
function degreesToRadians(degrees: number): number {
    return (degrees * Math.PI) / 180;
}

/* Formats an HSLColor with a given alpha into a CSS hsla() string. */
function hsla(color: HSLColor, alpha: number): string {
    return `hsla(${color.h}, ${color.s}%, ${color.l}%, ${alpha})`;
}

// ---------------------------------------------------------------------------
// Individual element renderers
// ---------------------------------------------------------------------------

/*
 * Renders a single filled square with rotation and opacity.
 *
 * Draws a rotated rectangle centered at (squareAbsX, squareAbsY) using
 * the square's color at the computed opacity.
 *
 * @param ctx - Canvas 2D rendering context
 * @param squareAbsX - Absolute X position of the square center (px)
 * @param squareAbsY - Absolute Y position of the square center (px)
 * @param square - Square configuration (size, rotation, color)
 * @param opacity - Final computed opacity to render with
 */
export function renderSquare(
    ctx: CanvasRenderingContext2D,
    squareAbsX: number,
    squareAbsY: number,
    square: SquareConfig,
    opacity: number,
): void {
    if (opacity <= 0) return;

    const w = square.size * (square.aspect || 1);
    const h = square.size;
    const halfW = w / 2;
    const halfH = h / 2;

    ctx.save();
    ctx.translate(squareAbsX, squareAbsY);
    ctx.rotate(degreesToRadians(square.rotation));
    ctx.strokeStyle = hsla(square.color, opacity);
    ctx.lineWidth = 2;
    ctx.strokeRect(-halfW, -halfH, w, h);
    ctx.restore();
}

/*
 * Renders the glow layer for a square — a blurred secondary rectangle
 * with lightness elevated by at least 20 points.
 *
 * Uses shadowBlur for the glow effect (spread >= 4px). This is more
 * performant than ctx.filter='blur(...)' which forces full-canvas
 * recomposition and can cause visible top-to-bottom repaint artifacts.
 *
 * @param ctx - Canvas 2D rendering context
 * @param squareAbsX - Absolute X position of the square center (px)
 * @param squareAbsY - Absolute Y position of the square center (px)
 * @param square - Square configuration (size, rotation, color)
 * @param opacity - Final computed opacity to render the glow with
 */
export function renderGlowLayer(
    ctx: CanvasRenderingContext2D,
    squareAbsX: number,
    squareAbsY: number,
    square: SquareConfig,
    opacity: number,
): void {
    if (opacity <= 0) return;

    const glowColor = computeGlowColor(square.color);
    const w = square.size * (square.aspect || 1);
    const h = square.size;
    const halfW = w / 2;
    const halfH = h / 2;

    ctx.save();
    ctx.translate(squareAbsX, squareAbsY);
    ctx.rotate(degreesToRadians(square.rotation));

    // Use shadowBlur for GPU-accelerated glow instead of ctx.filter
    ctx.shadowBlur = MIN_GLOW_SPREAD_PX * 2;
    ctx.shadowColor = hsla(glowColor, opacity);
    ctx.strokeStyle = hsla(glowColor, opacity);
    ctx.lineWidth = 2;
    ctx.strokeRect(-halfW, -halfH, w, h);

    // Reset shadow to avoid bleeding into subsequent draws
    ctx.shadowBlur = 0;
    ctx.shadowColor = 'transparent';
    ctx.restore();
}

/*
 * Renders a pseudocode text fragment using JetBrains Mono at the
 * specified font size and opacity.
 *
 * @param ctx - Canvas 2D rendering context
 * @param fragment - The active fragment with position, text, style info
 * @param opacity - Final computed opacity to render with
 */
export function renderFragment(
    ctx: CanvasRenderingContext2D,
    fragment: ActiveFragment,
    opacity: number,
): void {
    if (opacity <= 0) return;

    ctx.font = `${fragment.fontSize}px 'JetBrains Mono', monospace`;
    ctx.fillStyle = hsla(fragment.color, opacity);
    ctx.fillText(fragment.text, fragment.x, fragment.y);
}

// ---------------------------------------------------------------------------
// Per-frame orchestrator
// ---------------------------------------------------------------------------

/*
 * Orchestrates a complete frame render with correct draw order:
 *   1. Clear canvas (transparent background)
 *   2. Render all glow layers
 *   3. Render all square fills
 *   4. Render all fragment text
 *
 * Applies:
 * - Fade cycle base opacity (sinusoidal)
 * - Glow multiplier (distance-based within cluster)
 * - Cluster spawn fade-in factor (0→1 over fadeInDuration from spawnTime)
 * - Central zone opacity cap as final multiplier
 *
 * @param ctx - Canvas 2D rendering context
 * @param clusters - Active cluster configurations
 * @param fragments - Active pseudocode fragments
 * @param time - Current animation time in seconds
 * @param viewportWidth - Viewport width in pixels (for central zone cap)
 * @param canvasWidth - Canvas pixel width for clearRect
 * @param canvasHeight - Canvas pixel height for clearRect
 */
export function renderFrame(
    ctx: CanvasRenderingContext2D,
    clusters: ClusterConfig[],
    fragments: ActiveFragment[],
    time: number,
    viewportWidth: number,
    canvasWidth: number,
    canvasHeight: number,
): void {
    // Step 1: Clear canvas (transparent background — Requirement 6.3)
    ctx.clearRect(0, 0, canvasWidth, canvasHeight);

    // Pre-compute per-square render data for two-pass rendering (glow then fill)
    interface SquareRenderData {
        absX: number;
        absY: number;
        square: SquareConfig;
        opacity: number;
    }

    const squareRenderList: SquareRenderData[] = [];

    for (const cluster of clusters) {
        // Compute cluster spawn fade-in factor (Requirement 1.10)
        const elapsed = time - cluster.spawnTime;
        let fadeInFactor: number;
        if (elapsed <= 0) {
            fadeInFactor = 0;
        } else if (elapsed >= cluster.fadeInDuration) {
            fadeInFactor = 1;
        } else {
            fadeInFactor = elapsed / cluster.fadeInDuration;
        }

        // Skip cluster entirely if it hasn't spawned yet
        if (fadeInFactor <= 0) continue;

        const clusterAbsX = (cluster.centerXPct / 100) * canvasWidth;
        const clusterAbsY = (cluster.centerYPct / 100) * canvasHeight;

        for (const square of cluster.squares) {
            const squareAbsX = clusterAbsX + square.offsetX;
            const squareAbsY = clusterAbsY + square.offsetY;

            // Compute base opacity from fade cycle
            const baseOpacity = computeFadeOpacity(
                time,
                square.phaseOffset,
                square.cycleDuration,
                square.peakOpacity,
            );

            // Apply glow multiplier based on distance from cluster center
            const distanceFromCenter = Math.sqrt(
                square.offsetX * square.offsetX + square.offsetY * square.offsetY,
            );
            const glowMultiplier = computeGlowMultiplier(
                distanceFromCenter,
                cluster.boundingRadius,
            );

            // Combine: base * glow * fadeIn
            let opacity = baseOpacity * glowMultiplier * fadeInFactor;

            // Apply central zone opacity cap as final multiplier (Requirement 6.4)
            const cap = getCentralZoneOpacityCap(squareAbsX, viewportWidth);
            opacity = Math.min(opacity, cap);

            squareRenderList.push({absX: squareAbsX, absY: squareAbsY, square, opacity});
        }
    }

    // Step 2: Draw glow on the largest square per cluster only (performance-safe)
    const glowCandidates = new Set<number>();
    for (const cluster of clusters) {
        let largestIdx = -1;
        let largestSize = 0;
        for (let i = 0; i < squareRenderList.length; i++) {
            const item = squareRenderList[i];
            // Check if this square belongs to this cluster by matching offsets
            if (item.square.size > largestSize && item.opacity > 0.05) {
                // Find if this item is part of the current cluster
                const clusterAbsX = (cluster.centerXPct / 100) * canvasWidth;
                const clusterAbsY = (cluster.centerYPct / 100) * canvasHeight;
                const expectedX = clusterAbsX + item.square.offsetX;
                const expectedY = clusterAbsY + item.square.offsetY;
                if (Math.abs(expectedX - item.absX) < 0.01 && Math.abs(expectedY - item.absY) < 0.01) {
                    largestSize = item.square.size;
                    largestIdx = i;
                }
            }
        }
        if (largestIdx >= 0) {
            glowCandidates.add(largestIdx);
        }
    }
    for (const idx of glowCandidates) {
        const item = squareRenderList[idx];
        renderGlowLayer(ctx, item.absX, item.absY, item.square, item.opacity);
    }

    // Step 3: Draw all square outlines
    for (const item of squareRenderList) {
        renderSquare(ctx, item.absX, item.absY, item.square, item.opacity);
    }

    // Step 4: Draw all fragment text (Requirement 3.1 — text painted after all squares)
    for (const fragment of fragments) {
        // Compute fragment opacity from its lifecycle phase
        let opacity = computeFragmentOpacity(fragment, time);

        // Skip expired or invisible fragments
        if (opacity <= 0) continue;

        // Apply central zone opacity cap (Requirement 6.4)
        const cap = getCentralZoneOpacityCap(fragment.x, viewportWidth);
        opacity = Math.min(opacity, cap);

        renderFragment(ctx, fragment, opacity);
    }
}
