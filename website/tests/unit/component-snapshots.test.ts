/*
 * Component snapshot tests — Vitest-based HTML and data snapshot tests.
 *
 * Since Astro components render at build time and require the Astro compiler,
 * we snapshot the static data structures that feed into components and verify
 * the built HTML contains expected structural patterns.
 *
 */
import {describe, expect, it} from 'vitest';
import {features} from '../../src/data/features';
import {supportedPanels} from '../../src/data/hardware';
import {siteConfig} from '../../src/data/navigation';
import {existsSync, readFileSync} from 'node:fs';
import {resolve} from 'node:path';

// --- Data structure snapshots (component inputs) ---

describe('Component data snapshots', () => {
    it('features data snapshot matches expected structure', () => {
        expect(features).toMatchSnapshot();
    });

    it('hardware panels data snapshot matches expected structure', () => {
        expect(supportedPanels).toMatchSnapshot();
    });

    it('site config data snapshot matches expected structure', () => {
        expect(siteConfig).toMatchSnapshot();
    });

    it('navigation links derived from config snapshot', () => {
        const navLinks = [
            {label: 'Features', href: '#features', external: false, anchor: true},
            {label: 'Hardware', href: '#hardware', external: false, anchor: true},
            {label: 'GitHub', href: siteConfig.links.github, external: true, anchor: false},
            {label: 'Docs', href: siteConfig.links.docs, external: true, anchor: false},
        ];
        expect(navLinks).toMatchSnapshot();
    });
});

// --- Built HTML structure verification ---

describe('Built HTML structure verification', () => {
    const distPath = resolve(__dirname, '../../dist/index.html');
    const htmlExists = existsSync(distPath);

    // These tests require `astro build` to have been run first.
    // They verify the built output contains expected component structures.

    it.skipIf(!htmlExists)('contains NavigationBar with expected structure', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Nav bar should have sticky/fixed positioning class and proper landmark
        expect(html).toContain('role="banner"');
        expect(html).toContain('aria-label');

        // Contains navigation links
        expect(html).toContain('#features');
        expect(html).toContain('#hardware');
        expect(html).toContain('target="_blank"');
        expect(html).toContain('rel="noopener noreferrer"');
    });

    it.skipIf(!htmlExists)('contains HeroSection with h1 and CTA buttons', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Single h1 with project name
        expect(html).toContain('<h1');
        expect(html).toContain('CyberHUD');

        // Tagline present
        expect(html).toContain(siteConfig.tagline);

        // CTA buttons for install and GitHub
        expect(html).toContain(siteConfig.links.install);
        expect(html).toContain(siteConfig.links.github);

        // Canvas element for animation
        expect(html).toContain('<canvas');
    });

    it.skipIf(!htmlExists)('contains FeatureShowcase with all feature cards', () => {
        const html = readFileSync(distPath, 'utf-8');

        // All four features should be present
        for (const feature of features) {
            expect(html).toContain(feature.title);
            expect(html).toContain(feature.description);
            expect(html).toContain(feature.mediaAlt);
        }

        // Features section landmark
        expect(html).toContain('id="features"');
    });

    it.skipIf(!htmlExists)('contains HardwareSection grouped by chipset family', () => {
        const html = readFileSync(distPath, 'utf-8').replaceAll('&quot;', '"');

        // Every distinct chipset family and controller should be represented,
        // rather than every individual product (visitors already know their
        // own hardware's specs; what matters is whether their chip is supported).
        const families = [...new Set(supportedPanels.map((p) => p.family))];
        const controllers = [...new Set(supportedPanels.flatMap((p) => p.controller.split(' + ')))];
        for (const family of families) {
            expect(html).toContain(family);
        }
        for (const controller of controllers) {
            expect(html).toContain(controller);
        }

        // Hardware section landmark
        expect(html).toContain('id="hardware"');

        // Generic panel override / fallback statement
        expect(html).toContain('automatic fallback if the primary panel fails to init');

        // No-code value proposition should be present
        expect(html).toContain('no code required');
    });

    it.skipIf(!htmlExists)('contains Footer with external links and credits', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Footer landmark
        expect(html).toContain('role="contentinfo"');

        // Footer links
        expect(html).toContain(siteConfig.links.github);
        expect(html).toContain(siteConfig.links.docs);

        // Credits
        expect(html).toContain('Astro');
        expect(html).toContain(siteConfig.projectName);
    });

    it.skipIf(!htmlExists)('built HTML uses theme CSS custom properties', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Astro extracts styles into a separate CSS file during build.
        // Check the compiled CSS file for theme token references.
        const distDir = resolve(__dirname, '../../dist/_astro');
        const cssFiles = require('node:fs').readdirSync(distDir).filter((f: string) => f.endsWith('.css'));
        expect(cssFiles.length).toBeGreaterThan(0);

        const cssContent = cssFiles
            .map((f: string) => readFileSync(resolve(distDir, f), 'utf-8'))
            .join('\n');

        // Combine HTML and CSS for token verification (tokens may appear in either)
        const combined = html + '\n' + cssContent;

        const themeTokens = [
            '--color-bg-base',
            '--color-accent-cyan',
            '--font-mono',
            '--space-',
        ];

        for (const token of themeTokens) {
            expect(
                combined,
                `Expected built output (HTML + CSS) to reference theme token "${token}"`,
            ).toContain(token);
        }
    });

    it.skipIf(!htmlExists)('no section introduces fonts outside theme tokens', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Extract all font-family declarations from inline styles and style blocks
        const fontDeclarations = html.match(/font-family:\s*([^;}"]+)/g) || [];

        // Allowed font references (as they appear in CSS)
        const allowedFontPatterns = [
            'var(--font-mono)',
            'var(--font-body)',
            'jetbrains mono',
            'fira code',
            'courier new',
            'monospace',
            'inter',
            'system-ui',
            'sans-serif',
            'inherit',
        ];

        for (const declaration of fontDeclarations) {
            const value = declaration.replace('font-family:', '').trim().toLowerCase();
            const isAllowed = allowedFontPatterns.some((pattern) =>
                value.includes(pattern.toLowerCase()),
            );
            expect(
                isAllowed,
                `Found non-theme font declaration: "${declaration}"`,
            ).toBe(true);
        }
    });

    it.skipIf(!htmlExists)('snapshot of complete built HTML structure', () => {
        const html = readFileSync(distPath, 'utf-8');

        // Snapshot the structural skeleton (strip dynamic content like year, long text)
        // Extract just the tag structure for a lightweight snapshot
        const structuralElements = html.match(/<(header|nav|main|section|footer|h[1-6]|canvas)[^>]*>/g) || [];
        expect(structuralElements).toMatchSnapshot();
    });
});
