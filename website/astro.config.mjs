import {defineConfig} from 'astro/config';
import sitemap from '@astrojs/sitemap';

export default defineConfig({
    output: 'static',
    site: 'https://cyberhud.io',
    integrations: [
        sitemap({
            lastmod: new Date(),
        }),
    ],
});
