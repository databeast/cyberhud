import {type ConsoleMessage, expect, test} from '@playwright/test';

test.describe('Geometric Background Animation', () => {
    const viewports = [
        {name: 'mobile', width: 320, height: 568},
        {name: 'tablet', width: 768, height: 1024},
        {name: 'desktop', width: 1440, height: 900},
    ];

    test.describe('Hero section renders without crash across viewports', () => {
        for (const vp of viewports) {
            test(`renders at ${vp.width}px (${vp.name})`, async ({page}) => {
                await page.setViewportSize({width: vp.width, height: vp.height});
                await page.goto('/');

                const canvas = page.locator('#attract-canvas');
                await expect(canvas).toBeAttached();
                await expect(canvas).toBeVisible();

                const box = await canvas.boundingBox();
                expect(box).not.toBeNull();
                expect(box!.width).toBeGreaterThan(0);
                expect(box!.height).toBeGreaterThan(0);
            });
        }
    });

    test.describe('Performance', () => {
        test('frame rate stays ≥30fps on desktop (1440px) for 5 seconds (Req 5.1)', async ({
                                                                                               page,
                                                                                           }) => {
            await page.setViewportSize({width: 1440, height: 900});
            await page.goto('/');

            // Wait for animation to initialize
            await page.waitForTimeout(1000);

            const fps = await page.evaluate(() => {
                return new Promise<number>((resolve) => {
                    const frameTimes: number[] = [];
                    let lastTime = performance.now();

                    function measure(now: number) {
                        frameTimes.push(now - lastTime);
                        lastTime = now;

                        // Measure for 5 seconds
                        if (frameTimes.length < 150) {
                            requestAnimationFrame(measure);
                        } else {
                            // Calculate average frame time, excluding the first frame (warmup)
                            const relevant = frameTimes.slice(1);
                            const avgFrameTime =
                                relevant.reduce((sum, t) => sum + t, 0) / relevant.length;
                            const averageFps = 1000 / avgFrameTime;
                            resolve(averageFps);
                        }
                    }

                    requestAnimationFrame(measure);
                });
            });

            expect(fps, `Expected ≥30fps but got ${fps.toFixed(1)}fps`).toBeGreaterThanOrEqual(30);
        });
    });

    test.describe('First visible frame timing', () => {
        test('first visible frame renders within 2 seconds of page load (Req 5.5)', async ({
                                                                                               page,
                                                                                           }) => {
            await page.setViewportSize({width: 1440, height: 900});

            const startTime = Date.now();
            await page.goto('/');

            // Poll the canvas for non-transparent content
            const hasContent = await page.evaluate(() => {
                return new Promise<boolean>((resolve) => {
                    const deadline = Date.now() + 2000;

                    function check() {
                        const canvas = document.getElementById('attract-canvas') as HTMLCanvasElement | null;
                        if (!canvas) {
                            if (Date.now() < deadline) {
                                requestAnimationFrame(check);
                            } else {
                                resolve(false);
                            }
                            return;
                        }

                        const ctx = canvas.getContext('2d');
                        if (!ctx) {
                            resolve(false);
                            return;
                        }

                        // Sample pixels to detect non-zero (drawn) content
                        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
                        const data = imageData.data;
                        let hasNonZero = false;

                        // Check every 100th pixel for efficiency
                        for (let i = 3; i < data.length; i += 400) {
                            if (data[i] > 0) {
                                hasNonZero = true;
                                break;
                            }
                        }

                        if (hasNonZero) {
                            resolve(true);
                        } else if (Date.now() < deadline) {
                            requestAnimationFrame(check);
                        } else {
                            resolve(false);
                        }
                    }

                    requestAnimationFrame(check);
                });
            });

            const elapsed = Date.now() - startTime;
            expect(hasContent, `Canvas had no visible content within 2s (elapsed: ${elapsed}ms)`).toBe(
                true
            );
            expect(elapsed).toBeLessThan(2000);
        });
    });

    test.describe('Reduced motion', () => {
        test('with prefers-reduced-motion: reduce — animation does not render, no rAF calls (Req 4.4)', async ({
                                                                                                                   page,
                                                                                                               }) => {
            // Emulate reduced motion preference before navigation
            await page.emulateMedia({reducedMotion: 'reduce'});
            await page.goto('/');

            // Wait for any potential initialization
            await page.waitForTimeout(1000);

            // Track rAF calls for 2 seconds
            const rafCallCount = await page.evaluate(() => {
                return new Promise<number>((resolve) => {
                    let count = 0;
                    const originalRaf = window.requestAnimationFrame;

                    window.requestAnimationFrame = (cb: FrameRequestCallback) => {
                        count++;
                        return originalRaf(cb);
                    };

                    setTimeout(() => {
                        window.requestAnimationFrame = originalRaf;
                        resolve(count);
                    }, 2000);
                });
            });

            // With reduced motion, the animation should not be making rAF calls
            // Allow a small number for other page scripts, but animation loop should not run
            expect(
                rafCallCount,
                `Expected minimal rAF calls with reduced motion, but got ${rafCallCount}`
            ).toBeLessThan(5);

            // Verify canvas has no drawn content
            const hasContent = await page.evaluate(() => {
                const canvas = document.getElementById('attract-canvas') as HTMLCanvasElement | null;
                if (!canvas) return false;

                const ctx = canvas.getContext('2d');
                if (!ctx) return false;

                const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
                const data = imageData.data;

                for (let i = 3; i < data.length; i += 400) {
                    if (data[i] > 0) return true;
                }
                return false;
            });

            expect(hasContent, 'Canvas should have no rendered content with reduced motion').toBe(false);
        });
    });

    test.describe('Console errors', () => {
        test('no console errors during 10-second animation run', async ({page}) => {
            const errors: ConsoleMessage[] = [];

            page.on('console', (msg) => {
                if (msg.type() === 'error') {
                    errors.push(msg);
                }
            });

            await page.setViewportSize({width: 1440, height: 900});
            await page.goto('/');

            // Let the animation run for 10 seconds
            await page.waitForTimeout(10000);

            const errorTexts = errors.map((e) => e.text());
            expect(
                errorTexts,
                `Expected no console errors during 10s animation run, but found: ${errorTexts.join(', ')}`
            ).toHaveLength(0);
        });
    });
});
