# SEO Operations

This document records the recurring checks that sit outside the application's automated test suite. It is intentionally lightweight: implementation belongs in code, while deployment validation and performance observations belong here.

## Structured-data validation matrix

Validate a representative URL for each template after structured-data changes reach production.

| Page type | Representative URL | Expected primary types |
|---|---|---|
| Homepage | `/` | `Person`, `WebSite`, `ProfessionalService` |
| Service page | `/web-design/` | `Service`, `BreadcrumbList` |
| Local service page | `/web-design-leeds/` | `Service`, `BreadcrumbList` |
| Case study | `/work/salon-rebuild/` | `CreativeWork`, `BreadcrumbList` |
| Journal archive | `/blog/` | `Blog`, `BreadcrumbList` |
| Journal post | `/blog/{slug}` | `BlogPosting`, `BreadcrumbList` |

For each representative URL:

1. Run Google's Rich Results Test and resolve critical errors for supported result types.
2. Run Schema.org Validator and resolve vocabulary or shape errors.
3. Inspect the live URL in Google Search Console to confirm Google can fetch the rendered JSON-LD.
4. Request indexing only for important new or materially changed URLs.
5. Recheck Search Console enhancement reports after recrawling.

A successful validator result does not guarantee a rich result. Markup must remain accurate, visible in the page's subject matter and consistent with the canonical URL.

## Analytics decision record

Do not add analytics merely because it is available. Select a provider only after answering:

- What business question will the data answer?
- What personal data, cookies or persistent identifiers are collected?
- Is consent required for the proposed configuration?
- Can internal/development traffic be excluded?
- What is the performance cost?
- How long is data retained?
- Can the same goal be met with less data?

The primary conversion is a genuine project enquiry. Secondary interactions may include opening a case study, visiting a live project or following a source-code link, but only track them when the information will change a decision.

## Search Console baseline

Record a baseline immediately before sustained outreach or a significant publishing cycle.

| Metric | Baseline date | Value | Comparison date | Value | Notes |
|---|---|---:|---|---:|---|
| Total impressions |  |  |  |  |  |
| Total clicks |  |  |  |  |  |
| Click-through rate |  |  |  |  |  |
| Average position |  |  |  |  |  |

Also record:

- top queries;
- top landing pages;
- pages receiving impressions but few clicks;
- queries with pages ranking around positions 5–20;
- branded versus non-branded discovery;
- enquiry count for the same period.

## Monthly review

Once per month:

1. Compare Search Console performance with the previous period.
2. Identify pages or queries gaining impressions.
3. Review titles and descriptions only where click-through data supports a change.
4. Improve pages ranking around positions 5–20 by addressing search intent, content depth and internal linking.
5. Check for indexing, Core Web Vitals or structured-data regressions.
6. Record what changed and the date, then allow enough time for a meaningful comparison.

Avoid reacting to day-to-day movement or rewriting pages without enough data.

## Content priorities

Current publishing order:

1. Salon Rebuild design and engineering retrospective.
2. Server-rendered architecture: where it helps and where it does not.
3. Authentication and security patterns in Go.
4. UI/UX and engineering as one workflow.
5. Practical accessible web-development decisions.
6. Building and maintaining a small design system.

Each article should link to the most relevant service and proof of work, but links should be editorially useful rather than inserted solely for keywords.

## Local and authority work

Local SEO and backlinks are operational work rather than application features. Prioritise:

- accurate Leeds context where it is genuinely relevant;
- legitimate local business and freelancer listings;
- useful contributions to design/developer communities;
- references earned by articles, case studies and open-source packages;
- genuine client reviews after successful work.

Do not use manufactured reviews, mass-directory submissions, reciprocal-link schemes or purchased backlinks.
