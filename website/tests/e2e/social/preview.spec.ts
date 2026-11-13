import {expect, test} from '@playwright/test';

test.use({javaScriptEnabled: false});

test('serves a complete social preview without client-side rendering', async ({page, request}) => {
    const response = await page.goto('/');
    expect(response?.ok()).toBe(true);

    const html = await response!.text();
    expect(html).toContain('property="og:title"');
    expect(html).toContain('property="og:description"');
    expect(html).toContain('property="og:image"');
    await expect(page.locator('meta[property="og:title"]')).toHaveAttribute('content', /CyberHUD/);
    await expect(page.locator('meta[property="og:description"]')).toHaveAttribute('content', /CyberHUD/);
    await expect(page.locator('meta[property="og:url"]')).toHaveAttribute('content', 'https://cyberhud.io');
    await expect(page.locator('meta[property="og:type"]')).toHaveAttribute('content', 'website');
    await expect(page.locator('meta[property="og:image:type"]')).toHaveAttribute('content', 'image/png');
    await expect(page.locator('meta[property="og:image:width"]')).toHaveAttribute('content', '1200');
    await expect(page.locator('meta[property="og:image:height"]')).toHaveAttribute('content', '630');
    await expect(page.locator('meta[property="og:image:alt"]')).toHaveAttribute('content', /CyberHUD/);

    const imageUrl = await page.locator('meta[property="og:image"]').getAttribute('content');
    expect(imageUrl).toBe('https://cyberhud.io/images/og-preview.png?v=1');
    await expect(page.locator('meta[name="twitter:image"]')).toHaveAttribute('content', imageUrl!);
    const image = new URL(imageUrl!);
    const imageResponse = await request.get(image.pathname + image.search);
    expect(imageResponse.ok()).toBe(true);
    expect(imageResponse.headers()['content-type']).toContain('image/png');
    const body = await imageResponse.body();
    expect(body.subarray(0, 8)).toEqual(Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]));
    expect(body.readUInt32BE(16)).toBe(1200);
    expect(body.readUInt32BE(20)).toBe(630);
    expect(body.includes(Buffer.from('IDAT'))).toBe(true);
});
