/*
 * Accessibility end-to-end tests for the CyberHUD marketing site.
 *
  *
 * - axe-core automated WCAG 2.1 AA checks
 * - ARIA landmarks presence
 * - Keyboard navigation through interactive elements
 * - Focus indicator visibility (≥2px outline)
 * - Heading hierarchy (no skipping levels)
 */
import {expect, test} from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Accessibility — WCAG 2.1 AA compliance', () => {
    test.beforeEach(async ({page}) => {
        await page.goto('/');
    });

    /*
     * Run axe-core with WCAG 2.1 AA rules and assert zero violations.
     */
    test('page has no WCAG AA violations', async ({page}) => {
        const results = await new AxeBuilder({page})
            .withTags(['wcag2a', 'wcag2aa', 'wcag21aa'])
            .analyze();

        expect(results.violations).toEqual([]);
    });

    /*
     * Verify ARIA landmark roles: banner, navigation, main, contentinfo.
     */
    test('ARIA landmarks are present', async ({page}) => {
        const banner = page.locator('[role="banner"]');
        const navigation = page.locator('[role="navigation"]');
        const main = page.locator('[role="main"]');
        const contentinfo = page.locator('[role="contentinfo"]');

        await expect(banner).toBeVisible();
        await expect(navigation).toBeVisible();
        await expect(main).toBeVisible();
        await expect(contentinfo).toBeVisible();
    });

    /*
     * Tab through the page and verify focus moves to interactive elements
     * in logical document order (logo → nav links → hero CTAs → footer links).
     */
    test('keyboard navigation follows logical tab order', async ({page}) => {
        // Start focus at the top of the document
        await page.keyboard.press('Tab');

        // First focusable element should be the logo/home link
        const firstFocused = await page.evaluate(() => {
            const el = document.activeElement;
            return el?.getAttribute('aria-label') || el?.textContent?.trim() || null;
        });
        expect(firstFocused).toContain('CyberHUD');

        // Tab through the navigation links
        const expectedNavLabels = ['Features', 'Hardware', 'GitHub', 'Docs'];
        for (const label of expectedNavLabels) {
            await page.keyboard.press('Tab');
            const focusedText = await page.evaluate(() => {
                const el = document.activeElement;
                return el?.textContent?.trim() || el?.getAttribute('aria-label') || '';
            });
            expect(focusedText).toContain(label);
        }

        // Continue tabbing — next interactive elements should be reachable
        // (e.g., hero CTAs, then eventually footer links)
        await page.keyboard.press('Tab');
        const afterNav = await page.evaluate(() => {
            const el = document.activeElement;
            return {
                tagName: el?.tagName?.toLowerCase() || '',
                role: el?.getAttribute('role') || '',
            };
        });
        // The next element after nav should be focusable (a, button, etc.)
        expect(['a', 'button', 'input', 'select', 'textarea']).toContain(
            afterNav.tagName,
        );
    });

    /*
     * Tab to an interactive element and verify the focus indicator is visible
     * with at least a 2px outline.
     */
    test('focus indicators are visible with minimum 2px outline', async ({
                                                                             page,
                                                                         }) => {
        // Tab to the first interactive element
        await page.keyboard.press('Tab');

        const outlineInfo = await page.evaluate(() => {
            const el = document.activeElement;
            if (!el) return {width: '0px', style: 'none', color: ''};
            const computed = window.getComputedStyle(el);
            return {
                width: computed.outlineWidth,
                style: computed.outlineStyle,
                color: computed.outlineColor,
            };
        });

        // Outline must be at least 2px and not "none"
        const outlineWidth = parseFloat(outlineInfo.width);
        expect(outlineWidth).toBeGreaterThanOrEqual(2);
        expect(outlineInfo.style).not.toBe('none');
    });

    /*
     * Headings follow a logical sequence — no skipping levels.
     * e.g., h1 → h2 → h3 is valid; h1 → h3 is not.
     */
    test('heading hierarchy does not skip levels', async ({page}) => {
        const headingLevels = await page.evaluate(() => {
            const headings = Array.from(
                document.querySelectorAll('h1, h2, h3, h4, h5, h6'),
            );
            return headings.map((h) => parseInt(h.tagName.replace('H', ''), 10));
        });

        expect(headingLevels.length).toBeGreaterThan(0);

        // First heading should be h1
        expect(headingLevels[0]).toBe(1);

        // No heading should jump more than one level deeper than the previous
        for (let i = 1; i < headingLevels.length; i++) {
            const current = headingLevels[i];
            const previous = headingLevels[i - 1];
            // Going deeper: can only go one level at a time (h2→h3 OK, h2→h4 not)
            // Going shallower: any jump back up is fine (h3→h1 OK)
            if (current > previous) {
                expect(current - previous).toBeLessThanOrEqual(1);
            }
        }
    });
});
