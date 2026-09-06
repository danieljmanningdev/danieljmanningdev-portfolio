<div align="center">

<img src="web/static/djmdev-svg-assets/djmdev-horizontal-dark.svg" alt="Daniel J. Manning" width="420">

# Portfolio & Client Workspace

**A public portfolio, technical Journal and private administrative workspace, built with Go, HTMX, vanilla CSS and SQLite.**

[Live website](https://danieljmanningdev.com/) · [Case study](https://danieljmanningdev.com/work/portfolio) · [Journal](https://danieljmanningdev.com/blog/)

[![CI](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/ci.yml/badge.svg)](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/ci.yml)

</div>

## Product and architecture

The public site presents services, two selected projects, reusable developer tooling and technical writing. The private workspace manages clients, projects, contracts and Journal publishing.

Go renders HTML with `html/template`. HTMX enhances selected interactions; ordinary documents, forms and links remain the foundation. Stylesheets are tracked vanilla CSS using custom-property tokens. **There is no Tailwind build step and no Node/npm requirement to run the application.**

This is a **personal administrative application**, not a public-registration or multi-tenant SaaS product. The companion Go web libraries are related reusable work; this repository still has its own application-specific packages under `internal/`.

| Area | Implementation |
| --- | --- |
| Server | Go 1.27.1+, standard-library `net/http` |
| Persistence | SQLite via `modernc.org/sqlite`, numbered migrations |
| Journal | Goldmark Markdown, draft/published states |
| Structured data | `go-jsonld-schema` |
| Identity | bcrypt passwords, random bearer tokens stored as hashes |
| Operations | `log/slog`, request IDs, Docker, GitHub Actions, Fly.io |

```text
cmd/server/       HTTP server and routing
cmd/adminctl/     Explicit administrator creation
cmd/dbctl/        Database backup, verification and offline restore
internal/auth/    Sessions, passwords, CSRF and login throttling
internal/http/    Public pages and middleware
internal/blog/    Journal validation, Markdown and publishing
internal/clients/    Client workflows
internal/projects/   Project workflows
internal/contracts/  Contracts and exact integer money handling
internal/repository/ Parameterised database operations
internal/database/   SQLite, migrations and backup safety
internal/rendering/  Template loading
migrations/       Versioned SQL changes
web/templates/    Public, authentication and workspace HTML
web/static/       Vanilla CSS, HTMX, logos and responsive images
scripts/          Development-only image and browser tooling
docs/             Operational guidance and editorial drafts
```

## Run locally

Install Go 1.27.1 or later and Git:

```bash
git clone https://github.com/danieljmanningdev/danieljmanningdev-portfolio.git
cd danieljmanningdev-portfolio
go mod download
go run ./cmd/server
```

Open `http://localhost:8080/`. The application creates the local database and applies migrations. Public pages do not need an administrator. Edit CSS directly; restart the application after template edits because templates are parsed at startup.

Create your own administrator in a separate terminal:

```bash
go run ./cmd/adminctl -email "you@example.com" -name "Your name"
```

The command prompts without echoing the password. There is no public registration or committed demo credential.

| Variable | Default |
| --- | --- |
| `APP_ENV` | `development` |
| `APP_PORT` | `8080` |
| `DATABASE_PATH` | `./data/app.db` |
| `TEMPLATE_DIR` | `web/templates` |
| `LOG_LEVEL` | `info` |

Use `APP_ENV=production` and HTTPS in production. Keep the SQLite file on persistent storage. Never commit credentials, environment secrets, database files, backups or private client records.

## Quality checks

```bash
gofmt -w cmd internal
go vet ./...
go test -race ./...
go build ./...
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
git diff --check
```

CI checks formatting, static analysis, tests, the production image, startup and its non-root process. Companion workflows exercise browser/accessibility regressions and reachable known vulnerabilities. Browser tools are development-only and do not enter the production image.

Passing checks are evidence for those checks, not a guarantee of complete security or WCAG conformance. Direct test coverage varies by package. Vulnerability results depend on the toolchain, dependency versions and advisory database; a successful reachable-code scan can still report module-level advisories.

## Responsive assets

The homepage uses AVIF/WebP variants at 480, 960 and 1600px with intrinsic dimensions, a PNG fallback and below-the-fold lazy loading. Existing screenshot URLs remain available because published SQLite Markdown may reference them.

Optional regeneration (not an application build step):

```bash
python3 -m venv .venv
. .venv/bin/activate
python -m pip install Pillow==12.3.0
python scripts/build_responsive_images.py responsive
```

Review generated assets before committing. The `social` batch requires the DejaVu Sans development font and regenerates the 1200×630 social preview. `lossless-0` and `lossless-1` recompress existing PNGs without changing decoded pixels. Do not commit `.venv/` or Python cache files.

## Security and operational boundaries

Implemented controls include parameterised SQL, server-side validation, session expiry/revocation, secure production cookies, CSRF/cross-origin protection, request size limits, bounded login throttling, a restrictive CSP and private-response no-store/noindex policies.

The review adds exact contract-money parsing, buffered public HTML, strict public route matching, protection against static directory listings/dotfiles, explicit revocation failure handling, bounded Journal input, safer database publication/restore and correct process failure on a server bind error.

Operational limits remain important:

- The login limiter is process-local and uses the immediate peer. Trusted proxy identity and shared/edge limiting require deliberate design before scaling.
- Request bodies are capped at 4 MiB; Journal content at 1 MiB. Existing oversized content may need editing before it can be saved.
- Restore must be offline. WAL/SHM/journal sidecars are not safe to delete to bypass a check. Backup publication requires hard-link support on the destination filesystem.

See [database operations](docs/database-operations.md). Back up production data and rehearse recovery before a release.

## Deployment and remaining product work

The production Docker build uses Go 1.27.1 and runs the application as a non-root user. The Fly deployment workflow remains restricted to `main`, behind its CI job. Feature-branch work does not deploy to the live website.

[Owner actions](docs/seo/review-actions.md) document the remaining analytics, publishing, external authority and live-validation tasks. [Journal drafts](docs/journal-drafts/README.md) are not automatically published. Actual posts live in SQLite.

Search Console measurements are supplied separately to the owner rather than committed as public analytics. Genuine testimonials, client work, directory memberships and outreach cannot be manufactured by code.

## Author and licence

Designed and developed by **Daniel J. Manning — Digital Product Designer & Developer**.

[Website](https://danieljmanningdev.com/) · [GitHub](https://github.com/danieljmanningdev) · [Email](mailto:daniel@danieljmanningdev.com)

Released under the [MIT licence](LICENSE).
