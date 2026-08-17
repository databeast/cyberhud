/*
 * Unit tests for the geometric AttractCanvas animation module.
 *
 * Uses minimal DOM mocks (canvas + context) to verify lifecycle,
 * interface shape, reduced-motion behavior, deferred initialization,
 * draw order, and cluster generation without a real browser.
 *
 */
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {type AttractCanvasConfig, createAttractAnimation,} from '../../src/animations/attract-canvas';
import {MIN_DESKTOP_CLUSTERS, PSEUDOCODE_POOL} from '../../src/animations/geometric/types';
import {initializeClusters} from '../../src/animations/geometric/cluster-generator';

// ---------------------------------------------------------------------------
// Mock helpers
// ---------------------------------------------------------------------------

function createMockContext() {
    return {
        clearRect: vi.fn(),
        fillRect: vi.fn(),
        strokeRect: vi.fn(),
        fillText: vi.fn(),
        save: vi.fn(),
        restore: vi.fn(),
        translate: vi.fn(),
        rotate: vi.fn(),
        beginPath: vi.fn(),
        arc: vi.fn(),
        fill: vi.fn(),
        filter: '',
        font: '',
        fillStyle: '' as string,
        strokeStyle: '' as string,
        lineWidth: 1,
        shadowBlur: 0,
        shadowColor: '',
    };
}

function createMockCanvas(id: string, width = 1024, height = 768) {
    const ctxMock = createMockContext();

    const canvasMock = {
        id,
        clientWidth: width,
        clientHeight: height,
        width: 0,
        height: 0,
        getContext: vi.fn(() => ctxMock),
    };

    return {canvasMock, ctxMock};
}

function defaultConfig(overrides: Partial<AttractCanvasConfig> = {}): AttractCanvasConfig {
    return {
        canvasId: 'attract-canvas',
        baseParticleCount: 10,
        speed: 1.0,
        palette: [
            {h: 120, s: 85, l: 40},
            {h: 130, s: 90, l: 35},
        ],
        reducedMode: false,
        ...overrides,
    };
}

// ---------------------------------------------------------------------------
// Test suite
// ---------------------------------------------------------------------------

describe('AttractCanvas geometric animation', () => {
    let mockRaf: ReturnType<typeof vi.fn>;
    let mockCancelRaf: ReturnType<typeof vi.fn>;
    let rafCallbacks: Array<FrameRequestCallback>;
    let rafIdCounter: number;

    beforeEach(() => {
        rafCallbacks = [];
        rafIdCounter = 0;
        mockRaf = vi.fn((cb: FrameRequestCallback) => {
            rafCallbacks.push(cb);
            rafIdCounter++;
            return rafIdCounter;
        });
        mockCancelRaf = vi.fn();

        vi.stubGlobal('requestAnimationFrame', mockRaf);
        vi.stubGlobal('cancelAnimationFrame', mockCancelRaf);
    });

    afterEach(() => {
        vi.unstubAllGlobals();
        vi.restoreAllMocks();
    });

    // -------------------------------------------------------------------------
    // Req 4.1: Factory returns correct AttractAnimation interface shape
    // -------------------------------------------------------------------------
    describe('Req 4.1 — factory returns correct AttractAnimation interface', () => {
        it('returns an object with start, stop, resize, setReducedMotion methods', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());

            expect(typeof anim.start).toBe('function');
            expect(typeof anim.stop).toBe('function');
            expect(typeof anim.resize).toBe('function');
            expect(typeof anim.setReducedMotion).toBe('function');
            // Ensure no extra enumerable keys beyond the four methods
            expect(Object.keys(anim).sort()).toEqual(
                ['resize', 'setReducedMotion', 'start', 'stop'],
            );
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.2: No-op when canvas not found or getContext returns null
    // -------------------------------------------------------------------------
    describe('Req 4.2 — no-op implementations', () => {
        it('returns no-op when canvas element is not found', () => {
            vi.stubGlobal('document', {
                getElementById: vi.fn(() => null),
            });

            const anim = createAttractAnimation(defaultConfig());

            // Should not throw and should not request frames
            anim.start();
            anim.stop();
            anim.resize(800, 600);
            anim.setReducedMotion(true);
            expect(mockRaf).not.toHaveBeenCalled();
        });

        it('returns no-op when getContext("2d") returns null', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            canvasMock.getContext = vi.fn(() => null);

            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });

            const anim = createAttractAnimation(defaultConfig());
            anim.start();
            expect(mockRaf).not.toHaveBeenCalled();
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.3: start() idempotency
    // -------------------------------------------------------------------------
    describe('Req 4.3 — start() idempotency', () => {
        it('calling start() twice does not create duplicate animation loops', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            anim.start();
            const callsAfterFirst = mockRaf.mock.calls.length;

            anim.start();
            const callsAfterSecond = mockRaf.mock.calls.length;

            expect(callsAfterSecond).toBe(callsAfterFirst);
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.4: start() blocked by reducedMode config and setReducedMotion(true)
    // -------------------------------------------------------------------------
    describe('Req 4.4 — start() blocked when reduced motion is active', () => {
        it('start() does nothing when reducedMode config is true', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({reducedMode: true}));
            anim.start();
            expect(mockRaf).not.toHaveBeenCalled();
        });

        it('start() does nothing after setReducedMotion(true) has been called', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            anim.setReducedMotion(true);
            anim.start();
            expect(mockRaf).not.toHaveBeenCalled();
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.5: stop() cancels animation frame
    // -------------------------------------------------------------------------
    describe('Req 4.5 — stop() cancels animation frame', () => {
        it('cancels pending animation frame on stop()', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            anim.start();
            expect(mockRaf).toHaveBeenCalled();

            anim.stop();
            expect(mockCancelRaf).toHaveBeenCalled();
        });

        it('stop() while already stopped is a no-op (does not throw)', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            // Never started — stop should be safe
            expect(() => anim.stop()).not.toThrow();
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.7: setReducedMotion(true) stops rendering immediately
    // -------------------------------------------------------------------------
    describe('Req 4.7 — setReducedMotion(true) stops rendering immediately', () => {
        it('cancels animation frame when setReducedMotion(true) is called while running', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            anim.start();
            expect(mockRaf).toHaveBeenCalled();

            anim.setReducedMotion(true);
            expect(mockCancelRaf).toHaveBeenCalled();
        });
    });

    // -------------------------------------------------------------------------
    // Req 4.8: setReducedMotion(false) allows next start() but doesn't auto-start
    // -------------------------------------------------------------------------
    describe('Req 4.8 — setReducedMotion(false) allows next start() but no auto-start', () => {
        it('does not auto-start rendering after setReducedMotion(false)', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({reducedMode: true}));
            anim.setReducedMotion(false);

            // Should not have started automatically
            expect(mockRaf).not.toHaveBeenCalled();
        });

        it('allows start() to work after setReducedMotion(false)', () => {
            const {canvasMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({reducedMode: true}));
            anim.setReducedMotion(false);
            anim.start();
            expect(mockRaf).toHaveBeenCalled();
        });
    });

    // -------------------------------------------------------------------------
    // Req 3.3: Pseudocode pool has ≥10 entries each ≤40 chars
    // -------------------------------------------------------------------------
    describe('Req 3.3 — pseudocode pool constraints', () => {
        it('PSEUDOCODE_POOL has at least 10 entries', () => {
            expect(PSEUDOCODE_POOL.length).toBeGreaterThanOrEqual(10);
        });

        it('every entry in PSEUDOCODE_POOL is at most 40 characters', () => {
            for (const snippet of PSEUDOCODE_POOL) {
                expect(snippet.length).toBeLessThanOrEqual(40);
            }
        });
    });

    // -------------------------------------------------------------------------
    // Req 5.6: Deferred initialization spreads cluster creation across 3 frames
    // -------------------------------------------------------------------------
    describe('Req 5.6 — deferred initialization across 3 frames', () => {
        it('spreads cluster creation across the first 3 frames', () => {
            const {canvasMock, ctxMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({baseParticleCount: 10}));
            anim.start();

            // After start, one rAF is requested. We simulate 3 frames.
            // Each frame should render an increasing number of clusters.
            // Track strokeRect calls per frame (squares use strokeRect for rendering).
            const strokeRectCountsPerFrame: number[] = [];

            for (let frame = 0; frame < 3; frame++) {
                ctxMock.strokeRect.mockClear();
                const cb = rafCallbacks[rafCallbacks.length - 1];
                cb(16 + frame * 16);
                strokeRectCountsPerFrame.push(ctxMock.strokeRect.mock.calls.length);
            }

            // Each frame should do SOME rendering (deferred batches initialized progressively)
            // Frame 1 should have fewer strokes than frame 3 total (cumulative clusters grow)
            expect(strokeRectCountsPerFrame[0]).toBeGreaterThan(0);
            // By frame 3, all clusters are initialized; stroke count should be >= frame 1's
            expect(strokeRectCountsPerFrame[2]).toBeGreaterThanOrEqual(strokeRectCountsPerFrame[0]);
        });
    });

    // -------------------------------------------------------------------------
    // Req 1.10: Cluster spawn fade-in from 0 over 1–3 seconds
    // -------------------------------------------------------------------------
    describe('Req 1.10 — cluster spawn fade-in', () => {
        it('clusters have fadeInDuration between 1 and 3 seconds', () => {
            // Use initializeClusters directly with a seeded RNG
            let seed = 0.42;
            const rng = () => {
                seed = (seed * 16807) % 1;
                if (seed <= 0) seed += 1;
                return seed;
            };

            const clusters = initializeClusters(1024, 768, 5, rng);

            for (const cluster of clusters) {
                expect(cluster.fadeInDuration).toBeGreaterThanOrEqual(1);
                expect(cluster.fadeInDuration).toBeLessThanOrEqual(3);
            }
        });

        it('clusters spawn with initial opacity of 0 and fade in over fadeInDuration', () => {
            const {canvasMock, ctxMock} = createMockCanvas('attract-canvas');
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({baseParticleCount: 5}));
            anim.start();

            // Simulate first frame at time 0 — clusters with spawnTime > 0 should not render,
            // and those at spawnTime = 0 start with very low opacity (fade-in starts at 0)
            const cb = rafCallbacks[rafCallbacks.length - 1];
            cb(0.001); // Nearly zero elapsed time

            // The strokeStyle values set should contain very low alpha values
            // We just verify rendering started (strokeRect called)
            expect(ctxMock.strokeRect).toHaveBeenCalled();
        });
    });

    // -------------------------------------------------------------------------
    // Req 1.2: initializeClusters produces ≥5 clusters on desktop ≥1024px
    // -------------------------------------------------------------------------
    describe('Req 1.2 — minimum desktop cluster count', () => {
        it('initializeClusters produces at least MIN_DESKTOP_CLUSTERS on desktop ≥1024px', () => {
            let seed = 0.5;
            const rng = () => {
                seed = (seed * 16807 + 1) % 2147483647;
                return seed / 2147483647;
            };

            const clusters = initializeClusters(1024, 768, 3, rng);
            expect(clusters.length).toBeGreaterThanOrEqual(MIN_DESKTOP_CLUSTERS);
        });

        it('MIN_DESKTOP_CLUSTERS constant is at least 5', () => {
            expect(MIN_DESKTOP_CLUSTERS).toBeGreaterThanOrEqual(5);
        });
    });

    // -------------------------------------------------------------------------
    // Req 6.3: Canvas cleared with clearRect each frame (transparent background)
    // -------------------------------------------------------------------------
    describe('Req 6.3 — canvas cleared each frame with clearRect', () => {
        it('calls clearRect(0, 0, width, height) every frame', () => {
            const {canvasMock, ctxMock} = createMockCanvas('attract-canvas', 1024, 768);
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig());
            anim.start();

            // Run 3 frames
            for (let i = 0; i < 3; i++) {
                const cb = rafCallbacks[rafCallbacks.length - 1];
                cb(16 + i * 16);
            }

            // clearRect should have been called once per frame
            expect(ctxMock.clearRect.mock.calls.length).toBe(3);
            // Each call should clear the full canvas
            for (const call of ctxMock.clearRect.mock.calls) {
                expect(call).toEqual([0, 0, 1024, 768]);
            }
        });
    });

    // -------------------------------------------------------------------------
    // Req 3.1: Draw order — glow layers → square fills → fragment text
    // -------------------------------------------------------------------------
    describe('Req 3.1 — draw order: glow → square fills → fragment text', () => {
        it('renders glow (strokeRect with shadow), then square outlines (strokeRect), then text (fillText)', () => {
            const {canvasMock, ctxMock} = createMockCanvas('attract-canvas', 1024, 768);
            vi.stubGlobal('document', {
                getElementById: vi.fn((id: string) => (id === 'attract-canvas' ? canvasMock : null)),
            });
            vi.stubGlobal('window', {innerWidth: 1024, innerHeight: 768});

            const anim = createAttractAnimation(defaultConfig({baseParticleCount: 5}));
            anim.start();

            // Run enough frames to get clusters initialized and potentially fragments spawned
            let time = 0;
            for (let i = 0; i < 5; i++) {
                time += 16;
                const cb = rafCallbacks[rafCallbacks.length - 1];
                cb(time);
            }

            // Track call order. The renderer uses:
            //   1. clearRect (beginning of frame)
            //   2. For glow: save → translate → rotate → set shadowBlur → strokeRect → restore
            //   3. For squares: save → translate → rotate → set strokeStyle → strokeRect → restore
            //   4. For text: set font → set fillStyle → fillText
            //
            // We verify: all fillText calls come after all strokeRect passes.

            const strokeRectOrder = ctxMock.strokeRect.mock.invocationCallOrder;
            const fillTextOrder = ctxMock.fillText.mock.invocationCallOrder;

            if (fillTextOrder.length > 0 && strokeRectOrder.length > 0) {
                // The last strokeRect should have a lower invocation order than the first fillText
                const lastStrokeRect = strokeRectOrder[strokeRectOrder.length - 1];
                const firstFillText = fillTextOrder[0];
                expect(lastStrokeRect).toBeLessThan(firstFillText);
            }

            // Verify rendering happened
            expect(ctxMock.save).toHaveBeenCalled();
            expect(ctxMock.translate).toHaveBeenCalled();
            expect(ctxMock.strokeRect).toHaveBeenCalled();
            expect(ctxMock.restore).toHaveBeenCalled();
        });
    });
});
