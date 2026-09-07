# Owner actions after the portfolio quality review

## #47 — Journal

Four editorial drafts are in [journal-drafts](../journal-drafts/README.md). They cover architecture, security, accessible design implementation and the Salon retrospective. Publication is deliberately not automated. Review overlap with existing articles, voice, examples and claims, then publish through the admin interface. Real publication dates must be set by the application.

## #48 and #49 — local and open-source authority

The homepage now highlights `djm-cli`, the related Go web libraries as one stack, and `go-jsonld-schema`. Case studies continue to link to public source. README setup instructions are being corrected to the actual vanilla-CSS application.

Research shortlist, not completed submissions:

- [Leeds Digital](https://leedsdigital.org/share-your-event/): a genuine event/workshop proposal, not a generic business-directory submission. Confirm current eligibility and submission dates.
- [West & North Yorkshire Chamber](https://wnychamber.co.uk/membership/): review membership cost and business value before joining; do not buy membership solely for a backlink.
- [Show HN](https://news.ycombinator.com/showhn.html): only when a personally developed, usable tool satisfies the current guidelines and has been tested end to end. Do not submit a brochure or coordinate votes.

Use the public identity **Daniel J. Manning — Digital Product Designer & Developer**, canonical website and public contact email consistently. Do not invent an office address, review, client relationship or case-study result.

Outreach draft, not sent: “Hello, I'm Daniel J. Manning, a Digital Product Designer & Developer working with Go and server-rendered interfaces. Would a practical session on building and testing a small Go website be relevant to your community? I can propose a specific agenda and share a working example for review. My public work is at danieljmanningdev.com.” Only send after selecting a relevant recipient and agreeing to actually run the session.

Keep contact details and private notes outside this public repository. Track opportunity, eligibility, contact date, actual action, result URL and next review date. Mark citations complete only when public evidence exists. Future client testimonials require genuine work and permission.

## #50 — measurement

A real fixed-period Search Console baseline was retrieved through the connected Windsor.ai account during this review and supplied in the private owner handoff, not committed publicly. The property is `sc-domain:danieljmanningdev.com`; web search, 8 August–4 September 2026, fresh data excluded. The sitemap was also inspected read-only. No recurring export or account change was made.

Compare like-for-like periods monthly. Do not add page-level impressions and treat them as property totals, or treat named query rows as a complete query inventory. Missing dates/rows are not invented. A small sample is not proof of stable rankings or a conversion improvement.

Analytics-provider and consent choices are still outstanding. No third-party tracker or subscription is enabled. Where consent is required, test that rejecting or withdrawing consent prevents optional collection. Do not collect admin/workspace traffic, emails, session tokens, enquiry text or arbitrary query strings. A `mailto:` click means intent, not a sent enquiry. Record actual qualified enquiries separately and privately.

## #58 — images and release validation

AVIF/WebP sources at 480, 960 and 1600px are integrated with a PNG fallback and intrinsic dimensions. The original salon PNG was 4,254,849 bytes; the 960px AVIF is 34,951 bytes and WebP is 36,652 bytes. These are asset sizes, not live Core Web Vitals.

The corrupt Open Graph PNG was replaced by a valid 1200×630 image. Existing screenshots were losslessly recompressed. Original filenames remain because production SQLite Markdown may reference them; do not delete assets without a read-only inventory of published image URLs.

Review the PR's actual Lighthouse/browser artifacts, then run the public PageSpeed, security-header, rich-result and device checks after deployment. Localhost benchmarks do not model Cloudflare, Fly cold starts or actual visitor networks. Automated accessibility checks do not replace keyboard, zoom or assistive-technology review.

## Backend changes to know about

Money accepts non-negative decimal input with at most two fractional digits and rejects NaN/Inf, exponents and overflow. Request bodies are bounded to 4 MiB; Journal content is limited to 1 MiB, title/slug to 200 and excerpt to 500 characters (content limit is bytes). Existing oversized posts may require trimming before saving.

Logout returns an error rather than claiming success if token revocation fails. Login throttling is bounded but remains process-local and based on the immediate peer; trusted proxy identity/shared limiting need explicit design before scaling. Backup publication requires filesystem hard-link support. Restore must be offline, and refuses SQLite WAL/SHM/journal sidecars rather than deleting potential recovery data.

A successful reachable-vulnerability scan does not imply every dependency module is advisory-free. Inspect its verbose output and maintain dependencies; do not dismiss future advisories just because today's code does not reach them.
