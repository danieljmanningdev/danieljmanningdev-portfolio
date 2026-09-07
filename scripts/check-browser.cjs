/* Development-only browser regression checks; no frontend dependency. */
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');
const { execFileSync } = require('node:child_process');
const { chromium, firefox, webkit } = require('playwright');
const AxeBuilder = require('@axe-core/playwright').default;

const base = process.env.REVIEW_BASE_URL || 'http://127.0.0.1:18760';
const origin = new URL(base);
const database = path.resolve(process.env.REVIEW_DATABASE_PATH || '.');
const temporary = path.resolve(process.env.RUNNER_TEMP || '.');
if (process.env.CI !== 'true' || origin.origin !== 'http://127.0.0.1:18760' ||
    !database.startsWith(temporary + path.sep) || !fs.existsSync(database)) {
  throw new Error('Browser checks require CI and an isolated runner-temporary database.');
}
const reportDirectory = path.resolve('artifacts/browser');
fs.mkdirSync(reportDirectory, { recursive: true });
const token = crypto.randomBytes(32).toString('base64url');
const hash = crypto.createHash('sha256').update(token).digest('hex');
// Only disposable local fixtures: this never contacts or seeds production.
execFileSync('python3', ['-c', `
import sqlite3, sys
with sqlite3.connect(sys.argv[1]) as db:
    db.execute("INSERT INTO admins (email, password_hash, display_name) VALUES (?, ?, ?)",
               ('review@example.test', '!login-disabled-fixture', 'Review fixture'))
    admin_id = db.execute('SELECT last_insert_rowid()').fetchone()[0]
    db.execute("INSERT INTO admin_sessions (admin_id, token_hash, expires_at, last_seen_at) VALUES (?, ?, datetime('now', '+1 hour'), CURRENT_TIMESTAMP)", (admin_id, sys.argv[2]))
    db.execute("INSERT INTO blog_posts (title,slug,excerpt,content,status,published_at) VALUES (?,?,?,?,?,CURRENT_TIMESTAMP)",
               ('Review fixture: design & development', 'review-fixture', 'Disposable browser test content.',
                '## Readable content\\n\\nA [service link](/web-development/).\\n\\n- One item\\n- Another item\\n\\n~~~go\\nfmt.Println("Hello")\\n~~~', 'published'))
    db.execute("INSERT INTO blog_posts (title,slug,excerpt,content,status) VALUES (?,?,?,?,?)",
               ('Private draft', 'review-private-draft', 'Must not be public.', 'Draft body', 'draft'))
`, database, hash], { stdio: 'pipe' });

const pages = [
  ['/', 200], ['/web-design/', 200], ['/web-design-leeds/', 200],
  ['/web-development/', 200], ['/software-development/', 200], ['/ui-ux-design/', 200],
  ['/work/portfolio', 200], ['/work/salon-rebuild/', 200], ['/blog/', 200],
  ['/blog/review-fixture', 200], ['/login', 200], ['/review-missing-page', 404],
];
const widths = [320, 375, 390, 768, 1024, 1440];
const results = [];
const failures = [];
async function check(name, callback) {
  try { await callback(); results.push({ name, passed: true }); }
  catch (error) { failures.push({ name, error: String(error.stack || error) }); results.push({ name, passed: false }); }
}

function overflowInPage() {
  const width = document.documentElement.clientWidth;
  const bad = [];
  for (const el of document.querySelectorAll('h1,h2,h3,p,a,button,input,textarea,label,img,pre,ul,ol,dl')) {
    const box = el.getBoundingClientRect();
    if (el.matches('.skip-link:not(:focus)')) continue;
    if (box.width < 1 || box.height < 1 || el.closest('[aria-hidden="true"]') ||
        getComputedStyle(el).visibility === 'hidden') continue;
    if (el.getClientRects().length === 0) continue;
    if (box.right > width + 1 || box.left < -1) {
      bad.push({ element: el.tagName, class: el.className, left: box.left, right: box.right });
    }
  }
  return { width, scrollWidth: document.documentElement.scrollWidth, bad };
}

(async () => {
  for (const [engine, browserType] of Object.entries({ chromium, firefox, webkit })) {
    const browser = await browserType.launch({ headless: true });
    const context = await browser.newContext({ reducedMotion: 'reduce', colorScheme: 'dark' });
    const page = await context.newPage();
    const scriptErrors = [];
    const resourceErrors = [];
    page.on('pageerror', error => scriptErrors.push(error.message));
    page.on('response', response => {
      if (response.status() >= 400 && response.request().resourceType() !== 'document') {
        resourceErrors.push(`${response.status()} ${response.url()}`);
      }
    });
    for (const [route, status] of pages) {
      for (const width of widths) {
        await check(`${engine} ${route} ${width}px`, async () => {
          await page.setViewportSize({ width, height: 900 });
          const response = await page.goto(base + route, { waitUntil: 'load' });
          assert.equal(response.status(), status, 'Unexpected document status');
          await page.evaluate(() => document.fonts.ready);
          if (route === '/') {
            // A full-page screenshot alone does not trigger every lazy image.
            // Scroll and decode all homepage assets, then restore the top.
            for (const image of await page.locator('img').all()) {
              await image.scrollIntoViewIfNeeded();
              await image.evaluate(el => el.decode());
            }
            await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'instant' }));
          }
          await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))));
          const overflow = await page.evaluate(overflowInPage);
          assert.ok(overflow.scrollWidth <= width + 1, JSON.stringify(overflow));
          assert.deepEqual(overflow.bad, [], 'Visible content exceeds the viewport');
          assert.equal(await page.locator('main').count(), 1, 'One main landmark required');
          assert.equal(await page.locator('h1').count(), 1, 'One page heading required');
          const jsonld = await page.locator('script[type="application/ld+json"]').allTextContents();
          for (const text of jsonld) {
            const data = JSON.parse(text);
            assert.ok(data['@context'], 'JSON-LD context missing');
          }
          if (width === 390 || width === 1440) {
            const scan = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']).analyze();
            const name = `${engine}-${route.replace(/[^a-z0-9]/gi, '_') || 'home'}-${width}`;
            fs.writeFileSync(path.join(reportDirectory, `${name}.axe.json`), JSON.stringify({ url: route, violations: scan.violations, incomplete: scan.incomplete }, null, 2));
            // Keep screenshots even when an accessibility assertion fails.
            if (route === '/' || (route === '/work/portfolio' && width === 390)) {
              await page.screenshot({ path: path.join(reportDirectory, `${name}.png`), fullPage: true });
            }
            assert.deepEqual(scan.violations.map(v => ({ id: v.id, impact: v.impact, nodes: v.nodes.map(n => n.target) })), [], 'Automated accessibility findings');
          }
        });
      }
    }
    await check(`${engine} keyboard and native mobile navigation`, async () => {
      await page.setViewportSize({ width: 390, height: 900 });
      await page.goto(base + '/');
      await page.keyboard.press('Tab');
      assert.equal(await page.evaluate(() => document.activeElement.className), 'skip-link');
      await page.keyboard.press('Enter');
      assert.equal(await page.evaluate(() => document.activeElement.id), 'main-content');
      const menu = page.locator('.mobile-nav');
      const summary = menu.locator('summary');
      await summary.focus();
      await page.keyboard.press('Enter');
      assert.equal(await menu.getAttribute('open'), '');
      await page.keyboard.press('Tab');
      assert.equal(await page.evaluate(() => document.activeElement.getAttribute('href')), '/#work');
      const focus = await page.evaluate(() => getComputedStyle(document.activeElement).outlineStyle);
      assert.notEqual(focus, 'none', 'Keyboard focus must remain visible');
      await summary.focus();
      await page.keyboard.press('Enter');
      assert.equal(await menu.getAttribute('open'), null);
    });
    await check(`${engine} images decode and use responsive sources`, async () => {
      await page.goto(base + '/');
      for (const selector of ['.editorial-feature__media img', '.editorial-project--lead .editorial-project__media img']) {
        const image = page.locator(selector);
        await image.scrollIntoViewIfNeeded();
        await image.evaluate(el => el.decode());
        const info = await image.evaluate(el => ({ source: el.currentSrc, width: el.naturalWidth, height: el.naturalHeight }));
        assert.ok(info.width > 0 && info.height > 0);
        assert.match(info.source, /salon-rebuild-home-(480|960|1600)\.(avif|webp)$/);
        results.push({ name: `${engine} responsive image selected: ${selector}`, ...info, passed: true });
      }
    });
    for (const width of widths) {
      await check(`${engine} dark homepage ending ${width}px`, async () => {
        await page.setViewportSize({ width, height: 900 });
        await page.goto(base + '/', { waitUntil: 'load' });
        assert.equal(await page.locator('main[data-theme="light"]').count(), 1);
        assert.equal(await page.locator('.editorial-hero__avatar').count(), 0, 'Removed avatar must not return');
        assert.equal(await page.locator('.footer-cta').count(), 0, 'Removed duplicate CTA must not return');
        assert.equal(await page.locator('#contact').count(), 1);
        await page.evaluate(() => window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'instant' }));
        await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))));
        const ending = await page.evaluate(() => {
          const header = document.querySelector('.site-header');
          const panel = document.querySelector('.site-ending--home');
          const footer = document.querySelector('.site-footer');
          return {
            headerBottom: header.getBoundingClientRect().bottom,
            panelTop: panel.getBoundingClientRect().top,
            footerTop: footer.getBoundingClientRect().top,
            footerBottom: footer.getBoundingClientRect().bottom,
            viewport: window.innerHeight,
            colours: [header, panel, document.querySelector('#contact'), footer].map(el => getComputedStyle(el).backgroundColor),
          };
        });
        assert.ok(ending.panelTop <= ending.headerBottom + 1, 'Light content must not show between header and ending');
        assert.ok(Math.abs(ending.footerBottom - ending.viewport) <= 1, 'Footer must reach the viewport bottom');
        assert.ok(ending.footerTop > ending.viewport / 2, 'Footer belongs at the end, below the contact content');
        assert.equal(new Set(ending.colours).size, 1, 'Header, contact and footer must share the shell colour token');
        if (width === 390 || width === 1440) {
          await page.screenshot({ path: path.join(reportDirectory, `${engine}-home-ending-${width}.png`) });
        }
      });
    }
    await check(`${engine} no JavaScript content and navigation`, async () => {
      const noJS = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 390, height: 900 }, reducedMotion: 'reduce' });
      try {
        const noJSPage = await noJS.newPage();
        await noJSPage.goto(base + '/', { waitUntil: 'domcontentloaded' });
        // DOM readiness can precede stylesheet layout, especially in Firefox.
        await noJSPage.getByRole('heading', { name: 'Portfolio & Client Workspace', exact: true }).waitFor({ state: 'visible' });
        await noJSPage.getByRole('heading', { name: 'go-jsonld-schema', exact: true }).waitFor({ state: 'visible' });
        await noJSPage.locator('.mobile-nav summary').click();
        await noJSPage.locator('.mobile-nav-link').first().waitFor({ state: 'visible' });
      } finally {
        await noJS.close();
      }
    });
    await check(`${engine} private routes require authentication and drafts stay private`, async () => {
      const privateResponse = await context.request.get(base + '/dashboard/', { maxRedirects: 0 });
      assert.ok([302, 303].includes(privateResponse.status()));
      const draft = await context.request.get(base + '/blog/review-private-draft');
      assert.equal(draft.status(), 404);
    });
    await check(`${engine} authenticated workspace privacy`, async () => {
      await context.addCookies([{ name: 'djm_admin_session', value: token, url: base, httpOnly: true, sameSite: 'Lax' }]);
      await page.goto(base + '/dashboard/');
      assert.match(page.url(), /\/dashboard\/$/);
      const response = await context.request.get(base + '/dashboard/');
      assert.match(response.headers()['cache-control'] || '', /no-store/);
      assert.match(response.headers()['x-robots-tag'] || '', /noindex/);
      await context.clearCookies();
    });
    await check(`${engine} no script or asset failures`, async () => {
      assert.deepEqual(scriptErrors, []);
      assert.deepEqual(resourceErrors, []);
    });
    await context.close();
    await browser.close();
    fs.writeFileSync(path.join(reportDirectory, 'results.json'), JSON.stringify({ results, failures }, null, 2));
  }
  console.log(`${results.filter(r => r.passed).length} checks passed; ${failures.length} failed`);
  if (failures.length) { console.error(JSON.stringify(failures, null, 2)); process.exitCode = 1; }
})().catch(error => { console.error(error); process.exitCode = 1; });
