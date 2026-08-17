import {type ConsoleMessage, expect, test} from '@playwright/test';

test.describe('CyberHUD Marketing Site E2E', () => {
    test.describe('Critical content without JavaScript', () => {
        test('page loads h1, nav, and feature cards with JS disabled', async ({page}) => {
            await page.setJavaScriptEnabled(false);
            await page.goto('/');

            // h1 must be visible (rendered by HeroSection)
            const h1 = page.locator('h1');
            await expect(h1).toBeVisible();
            await expect(h1).toHaveText(/CyberHUD/i);

            // Navigation links present
            const nav = page.locator('nav[aria-label="Main navigation"]');
            await expect(nav).toBeVisible();

            const navLinks = nav.locator('a');
            await expect(navLinks).not.toHaveCount(0);

            // Feature cards are present in the DOM (they may not be "visible" due to
            // CSS animation requiring JS IntersectionObserver, but they must exist)
            const featureSection = page.locator('#features');
            await expect(featureSection).toBeAttached();

            const featureCards = page.locator('.feature-card');
            const cardCount = await featureCards.count();
            expect(cardCount).toBeGreaterThanOrEqual(1);
        });
    });

    test.describe('Sticky navigation bar', () => {
        test('nav bar remains visible after scrolling down', async ({page}) => {
            await page.goto('/');

            // Scroll down past the hero section
            await page.evaluate(() => window.scrollBy(0, 1200));
            await page.waitForTimeout(300);

            const header = page.locator('header.site-header');
            const boundingBox = await header.boundingBox();
            expect(boundingBox).not.toBeNull();

            // The top of the header should still be within or at the viewport top
            // (sticky positioning keeps it at top: 0)
            expect(boundingBox!.y).toBeGreaterThanOrEqual(0);
            expect(boundingBox!.y).toBeLessThan(100);
        });
    });

    test.describe('Mobile hamburger menu', () => {
        test('hamburger menu opens and shows all links at 375px viewport', async ({page}) => {
            await page.setViewportSize({width: 375, height: 667});
            await page.goto('/');

            // Hamburger toggle should be visible on mobile
            const toggle = page.locator('#mobile-menu-toggle');
            await expect(toggle).toBeVisible();

            // Open the menu
            await toggle.click();

            // Panel should become visible
            const panel = page.locator('#mobile-menu-panel');
            await expect(panel).toBeVisible();

            // All navigation links should be exposed
            const menuLinks = panel.locator('.mobile-menu-link');
            const linkTexts = await menuLinks.allTextContents();

            // Verify expected links are present (Features, Hardware, GitHub, Docs)
            const expectedLabels = ['Features', 'Hardware', 'GitHub', 'Docs'];
            for (const label of expectedLabels) {
                const matching = linkTexts.some((text) => text.includes(label));
                expect(matching, `Expected link "${label}" to be in mobile menu`).toBe(true);
            }

            // Each link should be visible
            for (let i = 0; i < (await menuLinks.count()); i++) {
                await expect(menuLinks.nth(i)).toBeVisible();
            }
        });
    });

    test.describe('External links', () => {
        test('all external links have target="_blank"', async ({page}) => {
            await page.goto('/');

            // Query all anchor elements
            const allLinks = page.locator('a[href]');
            const count = await allLinks.count();

            for (let i = 0; i < count; i++) {
                const link = allLinks.nth(i);
                const href = await link.getAttribute('href');

                if (!href) continue;

                // External links start with http:// or https:// (not internal anchors or relative)
                const isExternal = href.startsWith('http://') || href.startsWith('https://');
                if (isExternal) {
                    const target = await link.getAttribute('target');
                    expect(
                        target,
                        `External link "${href}" should have target="_blank"`
                    ).toBe('_blank');
                }
            }
        });
    });

    test.describe('Canvas animation', () => {
        test('canvas element exists and is visible when hero is in viewport', async ({page}) => {
            await page.goto('/');

            // Canvas element should exist
            const canvas = page.locator('#attract-canvas');
            await expect(canvas).toBeAttached();

            // When hero is visible (page just loaded, hero is in viewport), canvas should be visible
            await expect(canvas).toBeVisible();

            // Verify canvas has dimensions (is rendered)
            const box = await canvas.boundingBox();
            expect(box).not.toBeNull();
            expect(box!.width).toBeGreaterThan(0);
            expect(box!.height).toBeGreaterThan(0);
        });
    });

    test.describe('Lazy-loaded images', () => {
        test('feature card images have loading="lazy" attribute', async ({page}) => {
            await page.goto('/');

            // Feature card images should use lazy loading
            const featureImages = page.locator('.feature-card__image');
            const imgCount = await featureImages.count();
            expect(imgCount).toBeGreaterThan(0);

            for (let i = 0; i < imgCount; i++) {
                const img = featureImages.nth(i);
                const loadingAttr = await img.getAttribute('loading');
                expect(
                    loadingAttr,
                    `Feature card image ${i} should have loading="lazy"`
                ).toBe('lazy');
            }
        });
    });

    test.describe('No console errors', () => {
        test('page loads without console errors', async ({page}) => {
            const errors: ConsoleMessage[] = [];

            page.on('console', (msg) => {
                if (msg.type() === 'error') {
                    errors.push(msg);
                }
            });

            await page.goto('/');

            // Wait a moment for any async errors to surface
            await page.waitForTimeout(1000);

            const errorTexts = errors.map((e) => e.text());
            expect(
                errorTexts,
                `Expected no console errors, but found: ${errorTexts.join(', ')}`
            ).toHaveLength(0);
        });
    });
});
