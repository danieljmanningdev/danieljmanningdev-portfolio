# Quality baseline

This document records a representative automated Lighthouse CI baseline for the redesigned public website.

The audit runs against a locally started, production-shaped application in GitHub Actions after Tailwind CSS has been compiled. It covers the three principal public routes:

- `/`
- `/work/portfolio`
- `/blog/`

The values below were recorded on 29 August 2026 during one of the validation runs attached to pull request #33.

## Lighthouse results

| Route | Performance | Accessibility | Best practices | SEO | CLS | LCP | TBT |
|---|---:|---:|---:|---:|---:|---:|---:|
| `/` | 91 | 100 | 93 | 100 | 0 | 1,919 ms | 338 ms |
| `/work/portfolio` | 99 | 100 | 93 | 100 | 0 | 2,102 ms | 0 ms |
| `/blog/` | 99 | 100 | 93 | 100 | 0 | 1,802 ms | 0 ms |

Lighthouse measurements vary between runs and should be treated as regression signals rather than absolute laboratory guarantees. In this representative run, the homepage Total Blocking Time exceeded its warning target by 38 ms; later validation runs also passed and produced faster measurements.

## Enforced thresholds

The repository currently applies these assertions:

| Measurement | Threshold | Behaviour |
|---|---:|---|
| Accessibility | 95 or higher | Fails the workflow |
| Best practices | 90 or higher | Fails the workflow |
| SEO | 95 or higher | Fails the workflow |
| Cumulative Layout Shift | 0.10 or lower | Fails the workflow |
| Performance | 85 or higher | Warning |
| Largest Contentful Paint | 2,500 ms or lower | Warning |
| Total Blocking Time | 300 ms or lower | Warning |

Accessibility, best-practice, SEO and layout-stability regressions are release-blocking. Performance scores and individual timings remain warning-based because shared CI hardware introduces variance, while the recorded values still provide a useful regression baseline.

## Reports

Each Lighthouse workflow run uploads its HTML and JSON reports as a GitHub Actions artifact retained for 14 days.

The configuration lives in:

```text
lighthouserc.cjs
.github/workflows/lighthouse.yml
```

When the public information architecture changes, add any important new public route to the audited URL list. Authenticated workspace pages require a separate scripted-login audit and are intentionally not represented by this public baseline.
