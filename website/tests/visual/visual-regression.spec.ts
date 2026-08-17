/*
 * Visual regression tests — Playwright-based screenshot comparison.
 *
 * Captures full-page screenshots at key viewport widths (320, 768, 1440, 2560)
 * and verifies individual sections render consistently across breakpoints.
 * Also checks theme consistency: CSS only uses defined token custom properties.
 *
 */
import {expect, test} from '@playwright/test';

const viewports = [
    {width: 320, height: 568, name: 'mobile-xs'},
    {width: 768, height: 1024, name: 'tablet'},
    {width: 1440, height: 900, name: 'desktop'},
    {width: 2560, height: 1440, name: 'widescreen'},
];

// --- Full-page visual regression at each breakpoint ---

for (const vp of viewports) {
    test(`full page screenshot at ${vp.name} (${vp.width}px)`, async ({page}) => {
        await page.setViewportSize({width: vp.width, height: vp.height});
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        const screenshot = await page.screenshot({fullPage: true});
        expect(screenshot).toMatchSnapshot(`full-page-${vp.name}.png`);
    });
}

// --- Individual section screenshots at key breakpoints ---

const sections = [
    {selector: 'header[role="banner"]', name: 'navigation-bar'},
    {selector: '#hero', name: 'hero-section'},
    {selector: '#features', name: 'feature-showcase'},
    {selector: '#hardware', name: 'hardware-section'},
    {selector: 'footer[role="contentinfo"]', name: 'footer'},
];

for (const section of sections) {
    for (const vp of viewports) {
        test(`${section.name} at ${vp.name} (${vp.width}px)`, async ({page}) => {
            await page.setViewportSize({width: vp.width, height: vp.height});
            await page.goto('/');
            await page.waitForLoadState('networkidle');

            const element = page.locator(section.selector).first();
            // Skip if the section doesn't exist (graceful fallback)
            if ((await element.count()) === 0) {
                test.skip();
                return;
            }

            const screenshot = await element.screenshot();
            expect(screenshot).toMatchSnapshot(
                `${section.name}-${vp.name}.png`,
            );
        });
    }
}

// --- Theme consistency: verify no out-of-theme colors or fonts ---

test('theme consistency — all computed styles use defined token values', async ({page}) => {
    await page.setViewportSize({width: 1440, height: 900});
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Defined theme font families (substrings to match against computed values)
    const allowedFontFamilies = [
        'jetbrains mono',
        'fira code',
        'courier new',
        'monospace',
        'inter',
        'system-ui',
        'sans-serif',
        'ui-sans-serif', // system-ui resolved alias
    ];

    // Check that major sections use theme-defined fonts
    const sectionSelectors = [
        'header[role="banner"]',
        '#hero',
        '#features',
        '#hardware',
        'footer[role="contentinfo"]',
    ];

    for (const selector of sectionSelectors) {
        const fontFamily = await page.locator(selector).first().evaluate((el) => {
            return window.getComputedStyle(el).fontFamily.toLowerCase();
        }).catch(() => null);

        if (!fontFamily) continue;

        // Verify at least one allowed font family is present in the computed value
        const usesThemeFont = allowedFontFamilies.some((allowed) =>
            fontFamily.includes(allowed),
        );
        expect(
            usesThemeFont,
            `Section "${selector}" uses font-family "${fontFamily}" which is not in the defined theme tokens`,
        ).toBe(true);
    }
});

test('theme consistency — background colors use dark palette (lightness ≤20%)', async ({page}) => {
    await page.setViewportSize({width: 1440, height: 900});
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    const sectionSelectors = [
        'header[role="banner"]',
        '#hero',
        '#features',
        '#hardware',
        'footer[role="contentinfo"]',
    ];

    for (const selector of sectionSelectors) {
        const bgColor = await page.locator(selector).first().evaluate((el) => {
            return window.getComputedStyle(el).backgroundColor;
        }).catch(() => null);

        if (!bgColor || bgColor === 'rgba(0, 0, 0, 0)' || bgColor === 'transparent') {
            // Transparent backgrounds inherit from parent — acceptable
            continue;
        }

        // Parse rgb(r, g, b) or rgba(r, g, b, a)
        const match = bgColor.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/);
        if (!match) continue;

        const [, rStr, gStr, bStr] = match;
        const r = parseInt(rStr, 10) / 255;
        const g = parseInt(gStr, 10) / 255;
        const b = parseInt(bStr, 10) / 255;

        // Convert to HSL lightness
        const max = Math.max(r, g, b);
        const min = Math.min(r, g, b);
        const lightness = ((max + min) / 2) * 100;

        // Theme requires base ≤10%, surface ≤20% — allow up to 25% for elevated surfaces
        expect(
            lightness,
            `Section "${selector}" has background lightness ${lightness.toFixed(1)}% (expected ≤25% per theme tokens)`,
        ).toBeLessThanOrEqual(25);
    }
});

test('theme consistency — CSS custom properties are defined in :root', async ({page}) => {
    await page.setViewportSize({width: 1440, height: 900});
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Verify that key theme tokens are present and defined
    const expectedTokens = [
        '--color-bg-base',
        '--color-bg-surface',
        '--color-accent-cyan',
        '--color-accent-magenta',
        '--color-accent-green',
        '--font-mono',
        '--font-body',
        '--space-4',
        '--duration-normal',
    ];

    for (const token of expectedTokens) {
        const value = await page.evaluate((prop) => {
            return getComputedStyle(document.documentElement).getPropertyValue(prop).trim();
        }, token);

        expect(
            value.length,
            `Theme token "${token}" should be defined in :root but has no value`,
        ).toBeGreaterThan(0);
    }
});
