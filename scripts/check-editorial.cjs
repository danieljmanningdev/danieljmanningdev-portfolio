/* Audit regressions on the actual Go app. Development/CI only; no browser bundle.
   Keep the existing browser suite: these are additional content-stress checks. */
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { execFileSync } = require('node:child_process');
const { chromium, firefox, webkit } = require('playwright');
const AxeBuilder = require('@axe-core/playwright').default;

const base = process.env.REVIEW_BASE_URL || '';
const database = path.resolve(process.env.REVIEW_DATABASE_PATH || '.');
const temporary = path.resolve(process.env.RUNNER_TEMP || '.');
if (process.env.CI !== 'true' || base !== 'http://127.0.0.1:18760' ||
    !database.startsWith(temporary + path.sep) || !fs.existsSync(database)) {
  throw new Error('Editorial checks require CI, localhost and a disposable runner-temporary database.');
}
const output = path.resolve('artifacts/browser/editorial');
fs.mkdirSync(output, { recursive: true });
const longTitle = 'Designing understandable interfaces for complicated workflows: the decisions behind a server-rendered product, its accessible states and maintainable implementation';
const longURL = 'https://example.test/' + 'a-very-long-but-valid-path-segment-'.repeat(9);
const markdown = [
  '## A deliberately long section heading about readable content, predictable interaction and real implementation constraints',
  'Readable **strong text**, *emphasis* and [a service link](/web-development/).',
  '- First item\n  - Nested item\n- Second item',
  '1. Understand\n2. Design\n3. Build',
  '> A quoted design decision should remain readable, not fade into the background.',
  '[' + longURL + '](' + longURL + ')',
  '~~~go\nexample := "' + 'long-content-'.repeat(30) + '"\nfmt.Println(example)\n~~~',
  '| Layer | Responsibility |\n| --- | --- |\n| Interface | Semantic HTML and understandable states |\n| Persistence | Explicit boundaries and parameterised queries |',
].join('\n\n');
execFileSync('python3', ['-c', `
import sqlite3, sys
with sqlite3.connect(sys.argv[1]) as db:
    db.execute("INSERT INTO blog_posts (title,slug,excerpt,content,status,published_at) VALUES (?,?,?,?,?,CURRENT_TIMESTAMP) ON CONFLICT(slug) DO UPDATE SET title=excluded.title,excerpt=excluded.excerpt,content=excluded.content,status=excluded.status", (sys.argv[2], 'editorial-long-title', 'A disposable regression fixture for long titles, code, links and structured article content.', sys.argv[3], 'published'))
`, database, longTitle, markdown], { stdio: 'pipe' });

const routes = ['/', '/blog/', '/blog/editorial-long-title', '/ui-ux-design/',
  '/web-development/', '/software-development/', '/web-design/', '/web-design-leeds/',
  '/work/portfolio', '/work/salon-rebuild/'];
const sizes = [
  [320, 568], [375, 667], [390, 844], [640, 400], [768, 900],
  [1024, 768], [1440, 900], [1920, 1080],
];
const results = [];
const failures = [];
async function check(name, callback) {
  try { await callback(); results.push({ name, passed: true }); }
  catch (error) { failures.push({ name, error: String(error.stack || error) }); results.push({ name, passed: false }); }
}
function layoutEvidence() {
  const width = document.documentElement.clientWidth;
  const bad = [];
  const clipped = [];
  for (const el of document.querySelectorAll('h1,h2,h3,p,a,button,summary,dl,pre,table,img')) {
    if (el.matches('.skip-link:not(:focus)') || el.closest('[aria-hidden="true"]') ||
        !el.getClientRects().length || getComputedStyle(el).visibility === 'hidden') continue;
    const r = el.getBoundingClientRect();
    if (r.width < 1 || r.height < 1) continue;
    if (r.left < -1 || r.right > width + 1) bad.push({ tag: el.tagName, class: el.className, left: r.left, right: r.right });
    if (el.matches('h1,h2,h3') && (el.scrollHeight > el.clientHeight + 2 || el.scrollWidth > el.clientWidth + 2)) clipped.push(el.className);
  }
  const sample = document.createElement('div');
  sample.style.background = 'var(--color-ink-850)';
  document.body.append(sample);
  const ink = getComputedStyle(sample).backgroundColor;
  sample.remove();
  const header = document.querySelector('.site-header');
  const logo = document.querySelector('.brand-link');
  const hr = header.getBoundingClientRect(), lr = logo.getBoundingClientRect();
  return { width, scrollWidth: document.documentElement.scrollWidth, bad, clipped,
    ink, header: getComputedStyle(header).backgroundColor,
    footer: getComputedStyle(document.querySelector('.site-footer')).backgroundColor,
    flatHeader: lr.top >= hr.top - 1 && lr.bottom <= hr.bottom + 1 };
}
async function assertLayout(page) {
  const evidence = await page.evaluate(layoutEvidence);
  assert.ok(evidence.scrollWidth <= evidence.width + 1, JSON.stringify(evidence));
  assert.deepEqual(evidence.bad, [], 'Visible content outside the viewport');
  assert.deepEqual(evidence.clipped, [], 'A heading is clipped or overflows its box');
  assert.equal(evidence.header, evidence.ink, 'The header must use Ink 850, not pure black');
  assert.equal(evidence.footer, evidence.ink, 'The footer must share the header palette');
  assert.ok(evidence.flatHeader, 'The brand must stay inside the header');
}
async function imagesReady(page) {
  for (const image of await page.locator('img').all()) {
    await image.scrollIntoViewIfNeeded();
    await image.evaluate(el => new Promise((resolve, reject) => {
      if (el.complete) return el.naturalWidth ? resolve() : reject(new Error('Broken image: ' + el.src));
      el.addEventListener('load', resolve, { once: true });
      el.addEventListener('error', () => reject(new Error('Broken image: ' + el.src)), { once: true });
    }));
  }
  await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'instant' }));
}
(async () => {
  for (const [engine, type] of Object.entries({ chromium, firefox, webkit })) {
    const browser = await type.launch({ headless: true });
    try {
      const context = await browser.newContext({ colorScheme: 'dark', reducedMotion: 'reduce' });
      const page = await context.newPage();
      for (const route of routes) {
        for (const [width, height] of sizes) {
          await check(`${engine} editorial ${route} ${width}x${height}`, async () => {
            await page.setViewportSize({ width, height });
            const response = await page.goto(base + route, { waitUntil: 'load' });
            assert.equal(response.status(), 200);
            await page.evaluate(() => document.fonts.ready);
            await assertLayout(page);
            assert.equal(await page.locator('main').count(), 1);
            assert.equal(await page.locator('h1').count(), 1);
            assert.equal(await page.locator('link[rel="canonical"]').count(), 1);
            for (const text of await page.locator('script[type="application/ld+json"]').allTextContents()) assert.ok(JSON.parse(text)['@context']);
            if (width === 390 || width === 1440) {
              const name = `${engine}-${route.replace(/[^a-z0-9]/gi, '_')}-${width}`;
              await imagesReady(page);
              await page.screenshot({ path: path.join(output, name + '.png'), fullPage: true });
              const scan = await new AxeBuilder({ page }).withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa','wcag22aa']).analyze();
              fs.writeFileSync(path.join(output, name + '.axe.json'), JSON.stringify({ violations: scan.violations, incomplete: scan.incomplete }, null, 2));
              assert.deepEqual(scan.violations.map(v => ({ id: v.id, nodes: v.nodes.map(n => n.target) })), []);
            }
          });
        }
      }
      await check(`${engine} content stress and text enlargement`, async () => {
        await page.setViewportSize({ width: 390, height: 844 });
        await page.goto(base + '/blog/editorial-long-title');
        assert.equal(await page.locator('h1').textContent(), longTitle);
        assert.ok(await page.locator('.article-prose pre code').count());
        await page.locator('h1').evaluate(el => { el.textContent += ' ' + 'UnbrokenTitle'.repeat(12); });
        await assertLayout(page);
        await page.evaluate(() => { document.documentElement.style.fontSize = '200%'; });
        await assertLayout(page);
        // This supplements, not substitutes for, manual browser zoom testing.
      });
      await check(`${engine} ordinary phone hero actions are visible`, async () => {
        await page.setViewportSize({ width: 390, height: 844 });
        for (const route of ['/', '/ui-ux-design/', '/web-development/', '/software-development/']) {
          await page.goto(base + route);
          const action = page.locator('.editorial-hero__cta--primary, .service-actions .button-primary').first();
          const box = await action.boundingBox();
          assert.ok(box && box.y + box.height <= 844, route + ': primary action below first viewport');
        }
      });
      await check(`${engine} case previews retain their intrinsic ratio`, async () => {
        for (const route of ['/work/portfolio', '/work/salon-rebuild/']) {
          await page.goto(base + route);
          await imagesReady(page);
          const info = await page.locator('.case-media img').evaluate(el => ({
            natural: el.naturalWidth / el.naturalHeight,
            displayed: el.clientWidth / el.clientHeight,
            fit: getComputedStyle(el).objectFit,
          }));
          assert.equal(info.fit, 'contain');
          assert.ok(Math.abs(info.natural - info.displayed) < 0.03, JSON.stringify(info));
        }
      });
      await check(`${engine} journal and services work without JavaScript`, async () => {
        const noJS = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 390, height: 844 } });
        try {
          const p = await noJS.newPage();
          for (const route of ['/blog/', '/blog/editorial-long-title', '/ui-ux-design/']) {
            await p.goto(base + route, { waitUntil: 'load' });
            await p.locator('h1').waitFor({ state: 'visible' });
            await p.locator('.mobile-nav summary').click();
            await p.locator('.mobile-nav-link').first().waitFor({ state: 'visible' });
          }
        } finally { await noJS.close(); }
      });
      await context.close();
    } finally { await browser.close(); }
    fs.writeFileSync(path.join(output, 'results.json'), JSON.stringify({ results, failures }, null, 2));
  }
  console.log(`${results.filter(r => r.passed).length} editorial checks passed; ${failures.length} failed`);
  if (failures.length) { console.error(JSON.stringify(failures, null, 2)); process.exitCode = 1; }
})().catch(error => { console.error(error); process.exitCode = 1; });
