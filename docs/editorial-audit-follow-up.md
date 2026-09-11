# Public editorial design — audit follow-up

Basis: the owner's Product Quality Audit dated 9 September 2026, applied to the
existing dark editorial homepage branch on 11 September 2026. The historical
scores in that audit are not new test results.

## Scope and design decisions

The public shell, Journal archive, all Markdown article pages, UI/UX design,
web development, software development, the existing web-design/Leeds landing
pages, both case studies and related-content links now share the same editorial
measure, type roles, spacing rhythm, buttons and focus treatment. The private
workspace, backend logic, production data and deployment settings are unchanged.

The header and footer use the existing `--color-ink-850` foundation, rather than
pure black. The flattened centred brand remains. Text and accent colours stay
within Ink / Mist / Signal / Electric. Decorative glows, glass and false card
hover affordances are removed from the revised public page styles.

Real existing case-study image assets are shown without cropping. The Salon
Rebuild remains explicitly a fictional retrospective, not a commissioned client
redesign. Existing explanatory copy, links, metadata, CSP nonce handling and
server-rendered navigation are preserved.

## Findings addressed

| Audit finding | Implementation | Regression evidence |
| --- | --- | --- |
| UI-001 — hero hierarchy | Shared fluid hero scale, compact short-window/mobile spacing, primary actions before explanatory rail | 390 × 844 primary-action visibility assertions on homepage and three core services |
| UI-002 — long headings | Content-driven sizing, `minmax(0, 1fr)`, no fixed title heights, defensive wrapping and readable line heights | Deliberately long published test title, unbroken title token, large headings and 200% root text-size test |
| UI-003 — mobile rhythm | Shared mobile section rhythm, simpler fact rows, fewer nested panels and narrower reading column | Eight viewport sizes including narrow portrait and 640 × 400 landscape |
| A11Y-001 — subdued text | Reuse brighter Mist text roles for repeated public labels and body copy; retain clear focus/link treatment | Automated WCAG AA checks in three engines, plus existing keyboard suite |
| CSS-001 — token adoption | Public editorial semantic roles derived from existing palette/type/space tokens; replace page sheets rather than add another override stack | Shared `editorial.css`; old service/case polish layers removed |
| PERF-001 — opportunistic cleanup | No new runtime dependency, font, image asset or JavaScript bundle; remove obsolete CSS layers and version changed stylesheets | Existing application build and Lighthouse workflow retained |
| SEO-001 — preserve structure | Canonical/social/JSON-LD generation retained; meaningful headings and image dimensions preserved | Canonical/JSON-LD checks on revised routes; existing Go/SEO regression suite retained |

The original ten AA-passing/AAA-failing colour pairs were not enumerated in the
provided audit. This change improves repeated subdued roles but does not claim
that all ten historical combinations have been reconstructed or that the site
meets AAA. Automated checks are not a formal accessibility certification.

## Verification

`check-editorial.cjs` adds content-stress tests to the existing browser suite. It
runs only in CI against the fixed localhost origin and an isolated temporary
SQLite database. It never writes production articles or contacts production
services. The fixture is deliberately long and includes lists, a quote, a long
URL, code and a table.

The workflow runs Chromium, Firefox and WebKit, retains full-page screenshots
at 390px and 1440px, and checks layout, image proportions, semantic landmarks,
metadata, automated accessibility, text enlargement and no-JavaScript navigation.
Results are retained under `artifacts/browser/editorial/`; consult the actual run
for pass/fail status. No historical Lighthouse score is carried forward as a new
measurement.

Local offline layout fixtures were checked at 320, 375, 390, 640, 768, 1024,
1440 and 1920px. These are layout checks, not a substitute for the actual Go
application and downloaded project images. The changed Go templates also parse
and render in a standard-library harness; the additional browser script passes
`node --check`.

## Boundaries and follow-up

- Retain manual visual approval on the owner's Mac/font stack and testing with
  real published content. Complete screen-reader/browser coverage remains manual.
- The optional CSP/COEP/HSTS suggestions were not blindly enabled. No security
  policy or operational architecture was changed as part of visual polish.
- The historical HTTP redirect discrepancy needs separate network-level
  verification; this work does not relabel it as a confirmed defect or resolved
  issue. Search Console monitoring and external authority remain ongoing work.
- Existing race tests, vulnerability checks, Docker smoke checks and non-root
  verification remain intact. This branch must not be treated as deployed until
  explicitly reviewed and merged through the normal process.
