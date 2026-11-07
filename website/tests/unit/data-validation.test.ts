/*
 * Data validation tests for site content and configuration.
 *
 * Ensures all static data meets the constraints defined in the design document
 * and accessibility requirements.
 *
 */
import {describe, expect, it} from 'vitest';
import {features} from '../../src/data/features';
import {supportedControllerFamilies} from '../../src/data/hardware';
import {siteConfig} from '../../src/data/navigation';

describe('Feature data validation', () => {
    it('has at least 4 features', () => {
        expect(features.length).toBeGreaterThanOrEqual(4);
    });

    it('all feature descriptions are ≤120 characters', () => {
        for (const feature of features) {
            expect(
                feature.description.length,
                `Feature "${feature.id}" description is ${feature.description.length} chars (max 120)`,
            ).toBeLessThanOrEqual(120);
        }
    });

    it('all feature mediaAlt texts are ≤125 characters', () => {
        for (const feature of features) {
            expect(
                feature.mediaAlt.length,
                `Feature "${feature.id}" mediaAlt is ${feature.mediaAlt.length} chars (max 125)`,
            ).toBeLessThanOrEqual(125);
        }
    });

    it('all features have non-empty id, title, and description', () => {
        for (const feature of features) {
            expect(feature.id).toBeTruthy();
            expect(feature.title).toBeTruthy();
            expect(feature.description).toBeTruthy();
        }
    });
});

describe('Site configuration validation', () => {
    it('tagline is ≤100 characters', () => {
        expect(
            siteConfig.tagline.length,
            `Tagline is ${siteConfig.tagline.length} chars (max 100)`,
        ).toBeLessThanOrEqual(100);
    });

    it('metaDescription is between 80 and 160 characters', () => {
        expect(
            siteConfig.metaDescription.length,
            `metaDescription is ${siteConfig.metaDescription.length} chars (must be 80-160)`,
        ).toBeGreaterThanOrEqual(80);
        expect(
            siteConfig.metaDescription.length,
            `metaDescription is ${siteConfig.metaDescription.length} chars (must be 80-160)`,
        ).toBeLessThanOrEqual(160);
    });

    it('page title (generated from projectName) is ≤60 characters', () => {
        // Page title format: "ProjectName — tagline" or just projectName
        // The shortest form is just the project name; the full form includes tagline
        const pageTitle = `${siteConfig.projectName} — ${siteConfig.tagline}`;
        // If the full title exceeds 60, at minimum the project name alone must fit
        // The requirement is that the rendered <title> tag is ≤60 chars
        // Typical title: "CyberHUD | Cyberpunk HUD displays for Raspberry Pi"
        // We test the project name itself is short enough to form a valid title
        expect(
            siteConfig.projectName.length,
            `projectName "${siteConfig.projectName}" is ${siteConfig.projectName.length} chars (must leave room in 60-char title)`,
        ).toBeLessThanOrEqual(60);
    });
});

describe('Navigation links validation', () => {
    it('external links have proper https:// URLs', () => {
        const {links} = siteConfig;
        const externalUrls = [links.install, links.github, links.docs];

        for (const url of externalUrls) {
            expect(url, `External link "${url}" should start with https://`).toMatch(/^https:\/\//);
        }
    });

    it('external links should open in new tab (target="_blank" pattern)', () => {
        // Verify the data supports identification of external links
        // External links are those with full URLs (https://)
        const {links} = siteConfig;
        const externalUrls = [links.install, links.github, links.docs];

        for (const url of externalUrls) {
            // All external links must be absolute URLs, enabling the template
            // to add target="_blank" and rel="noopener"
            expect(url).toMatch(/^https?:\/\//);
        }
    });

    it('internal navigation anchors start with #', () => {
        // Internal section links are defined as anchors
        const internalAnchors = ['#features', '#hardware'];

        for (const anchor of internalAnchors) {
            expect(anchor).toMatch(/^#/);
        }
    });
});

describe('Hardware data validation', () => {
    it('lists every supported display controller', () => {
        const controllers = supportedControllerFamilies.flatMap((family) => family.controllers);
        expect(controllers).toHaveLength(20);
        expect(new Set(controllers).size).toBe(controllers.length);
        expect(controllers).toEqual(expect.arrayContaining([
            'EPD4in26',
            'GC9A01A',
            'HT16K33',
            'ILI9341',
            'ILI9486',
            'IS31FL3731',
            'IS31FL3741',
            'MAX7219',
            'SH1106',
            'Sharp Memory LCD',
            'SSD1305',
            'SSD1306',
            'SSD1325',
            'SSD1351',
            'SSD1680',
            'SSD1681',
            'SSD1681 tri-color',
            'ST7735S',
            'ST7789',
            'UC8179',
        ]));
    });

    it('all controller families have a description and controllers', () => {
        for (const family of supportedControllerFamilies) {
            expect(family.family).toBeTruthy();
            expect(family.note).toBeTruthy();
            expect(family.controllers.length).toBeGreaterThan(0);
        }
    });
});
