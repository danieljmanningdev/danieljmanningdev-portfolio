<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="web/static/djmdev-svg-assets/djmdev-horizontal-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="web/static/djmdev-svg-assets/djmdev-horizontal-light.svg">
  <img
    src="web/static/djmdev-svg-assets/djmdev-horizontal-light.svg"
    alt="Daniel J. Manning"
    width="420"
  >
</picture>

# Portfolio & Personal Workspace

**A production-deployed, server-rendered Go application combining a public digital-product portfolio, a Markdown journal, and a private personal workspace.**

[Live website](https://danieljmanningdev.com/) ·
[Case study](https://danieljmanningdev.com/work/portfolio) ·
[Journal](https://danieljmanningdev.com/blog/)

[![CI](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/ci.yml/badge.svg)](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/ci.yml)
[![Fly Deploy](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/fly.yml/badge.svg)](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/actions/workflows/fly.yml)
[![Go](https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-2563EB.svg)](LICENSE)

</div>

---

## Overview

This repository contains the complete source for
[danieljmanningdev.com](https://danieljmanningdev.com/).

It is designed and engineered as one coherent product with two distinct surfaces:

| Surface | Purpose |
|---|---|
| **Public portfolio** | Presents selected work, product-design capability, engineering experience, a detailed case study, and a technical journal. |
| **Personal workspace** | Provides an authenticated area for managing clients, projects, contracts, and journal publishing. |

The application is intentionally **server-rendered first**. Go produces the HTML, while HTMX is used only where a focused interaction benefits from progressive enhancement. There is no client-side router, hydration layer, or duplicated JSON API.

The project demonstrates the overlap between product design and engineering: information architecture, interaction design, visual systems, backend structure, authentication, persistence, testing, deployment, and ongoing product operations all live in the same codebase.

## Highlights

- Distinct public, authentication, and personal-workspace layouts
- Responsive portfolio with selected work, capabilities, process, and contact sections
- Long-form portfolio case study
- Markdown journal powered by Goldmark
- Custom journal administration for drafting, editing, publishing, unpublishing, and deleting posts
- Client, project, and contract management
- Server-side sessions with idle and absolute expiry
- CSRF protection, login throttling, cross-origin protection, and security headers
- Custom accessible 404 experience
- Automated Go checks, tests, template parsing, and CSS compilation in CI
- Multi-stage Docker build and Fly.io deployment with persistent SQLite storage

## Technology

| Area | Technology |
|---|---|
| Language | Go 1.26.5 |
| HTTP | Standard-library `net/http` |
| Templates | `html/template` |
| Progressive enhancement | HTMX |
| Design tokens | CSS-first `@theme` configuration with OKLCH colours |
| Database | SQLite through `modernc.org/sqlite` |
| Markdown | Goldmark |
| Authentication | bcrypt passwords and server-side sessions |
| Logging | `log/slog` |
| Testing | Go `testing`, `httptest`, and in-memory SQLite |
| Packaging | Multi-stage Docker |
| Hosting | Fly.io with a persistent volume |
| Automation | GitHub Actions |

## Architecture

The application keeps request handling, domain behaviour, persistence, and rendering in clear layers:

```text
Browser request
      ↓
Request logging
      ↓
Security headers
      ↓
Cross-origin protection
      ↓
Router
      ↓
Authentication middleware
      ↓
Feature handler
      ↓
Repository
      ↓
SQLite
      ↓
Go model + HTML template
      ↓
Browser response
```

Feature areas are separated into focused packages rather than placed inside one large HTTP package:

```text
internal/
├── auth/        # Passwords, tokens, CSRF, login limiting, sessions
├── blog/        # Personal journal administration and validation
├── clients/     # Client HTTP flows and forms
├── config/      # Environment-driven application configuration
├── contracts/   # Contract HTTP flows, validation, and relationships
├── database/    # SQLite setup and migrations
├── http/        # Public pages, authentication, middleware, and routing helpers
├── logging/     # Structured logging
├── models/      # Shared application models
├── projects/    # Project HTTP flows and forms
├── rendering/   # Public, workspace, and authentication template loading
└── repository/  # Parameterised persistence operations
```

The core workspace relationships are:

```text
Client
├── Projects
└── Contracts
    └── Optional project association
```

A contract always belongs to a client. When a project is selected, the application validates that it belongs to the same client.

## Features

### Public portfolio

- Product-led responsive homepage
- Selected-work presentation
- Detailed portfolio and workspace case study
- Capabilities and working-process sections
- Public Markdown journal
- Individual article pages
- Page descriptions and Open Graph metadata
- Responsive navigation
- Custom noindex 404 page
- Accessible focus states, skip links, and reduced-motion support
- Shared branded header and footer

### Personal workspace

- Authenticated dashboard
- Dedicated workspace navigation
- Client creation, viewing, editing, status management, and deletion
- Project creation, viewing, editing, client association, dates, and archiving
- Contract creation, viewing, editing, client/project relationships, values, dates, and cancellation
- Journal archive
- Markdown post editor
- Draft and published states
- Publish and unpublish actions
- HTMX-enhanced deletion with confirmation
- Server-side form parsing and validation
- Empty states and contextual actions

The workspace is intentionally a **personal administrative tool**, not a public registration system or multi-tenant SaaS product.

### Authentication and security

- No public account-registration route
- bcrypt password hashing
- Generic login-failure responses
- Login attempt throttling
- Cryptographically random session tokens
- SHA-256 session-token hashes stored in SQLite
- 24-hour absolute session lifetime
- 30-minute idle timeout
- Periodic session activity refresh
- Server-side session revocation on logout
- Inactive-admin rejection
- HttpOnly and SameSite cookies
- Secure cookies when `APP_ENV=production`
- CSRF protection for authentication and workspace commands
- Standard-library cross-origin request protection
- `Cache-Control: no-store` for authenticated responses
- Noindex policy for authentication and workspace pages
- Framing, content-type, referrer, permissions, and HSTS response controls
- Hardened HTTP server timeouts and header limits

## Project structure

```text
.
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── fly.yml
├── cmd/
│   ├── adminctl/
│   └── server/
├── docs/
│   └── site-experience-plan.md
├── internal/
│   ├── auth/
│   ├── blog/
│   ├── clients/
│   ├── config/
│   ├── contracts/
│   ├── database/
│   ├── http/
│   ├── logging/
│   ├── models/
│   ├── projects/
│   ├── rendering/
│   └── repository/
├── migrations/
│   ├── 001_initial.sql
│   ├── 002_auth.sql
│   └── 003_blog.sql
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── djmdev-svg-assets/
│   │   ├── images/
│   │   └── scripts/
│   └── templates/
│       ├── components/
│       ├── layouts/
│       └── pages/
├── Dockerfile
├── fly.toml
├── go.mod
├── go.sum
├── package.json
└── package-lock.json
```

Generated CSS, local databases, environment files, and dependency directories are excluded from Git.

## Local development

### Requirements

- Go 1.26.5 or newer
- Node.js 22 or newer recommended
- npm
- Git

The SQLite command-line client is optional but useful for inspecting local data.

### Clone

```bash
git clone https://github.com/danieljmanningdev/danieljmanningdev-portfolio.git
cd danieljmanningdev-portfolio
```

### Install frontend dependencies

```bash
npm ci
```

### Build the stylesheet

```bash
npm run build:css
```

Tailwind writes the generated stylesheet to:

```text
web/static/css/vanilla/main.css
```

That generated file is ignored by Git and is rebuilt in CI and during the Docker build.

### Start the application

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080
```

### Watch CSS during development

Run the Go server in one terminal:

```bash
go run ./cmd/server
```

Run Tailwind in another:

```bash
npm run dev:css
```

## Create an administrator

There is no public registration page. Administrative accounts are created explicitly with `adminctl`.

```bash
go run ./cmd/adminctl \
  -email "admin@example.com" \
  -name "Daniel Manning"
```

The command prompts for a password and confirmation without echoing the password to the terminal.

Only the bcrypt password hash is persisted.

## Configuration

The application reads configuration from environment variables and otherwise uses development defaults.

| Variable | Default | Purpose |
|---|---:|---|
| `APP_ENV` | `development` | Application environment |
| `APP_PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `./data/app.db` | SQLite database path |
| `TEMPLATE_DIR` | `web/templates` | Template root |
| `LOG_LEVEL` | `info` | Structured log level |

Example:

```bash
APP_ENV=development \
APP_PORT=8080 \
DATABASE_PATH=./data/app.db \
TEMPLATE_DIR=web/templates \
LOG_LEVEL=debug \
go run ./cmd/server
```

Set `APP_ENV=production` in production to enable production-only cookie and HSTS behaviour.

## Useful routes

| Route | Purpose | Access |
|---|---|---|
| `/` | Public portfolio | Public |
| `/work/portfolio` | Portfolio and workspace case study | Public |
| `/blog/` | Journal archive | Public |
| `/blog/{slug}` | Published journal article | Public |
| `/health` | Health endpoint | Public |
| `/login` | Administrator sign-in | Public |
| `/dashboard/` | Personal workspace overview | Authenticated |
| `/dashboard/clients` | Client management | Authenticated |
| `/dashboard/projects` | Project management | Authenticated |
| `/dashboard/contracts` | Contract management | Authenticated |
| `/dashboard/blog` | Journal administration | Authenticated |

## Database and migrations

SQLite is opened automatically when the application starts. The database package:

- Creates the database directory when required
- Enables and verifies foreign-key enforcement
- Runs versioned SQL migrations
- Tracks applied versions in `schema_migrations`
- Uses in-memory databases throughout the test suite where appropriate

Current migrations provide:

```text
001_initial.sql  → clients, projects, contracts, and foundation tables
002_auth.sql     → administrators and server-side sessions
003_blog.sql     → journal posts and publishing indexes
```

Local application data is stored under `data/` by default and is intentionally excluded from source control.

## Quality checks

Run the full local verification sequence before opening or merging a pull request:

```bash
gofmt -w cmd internal
go vet ./...
go test ./...
npm run build:css
git diff --check
```

The CI workflow runs on pull requests and pushes to `main`, checking:

- Go formatting
- `go vet`
- All Go tests
- Application-template parsing
- Frontend dependency installation

The test suite covers configuration, migrations, foreign keys, repositories, authentication, session expiry, CSRF behaviour, middleware, routing, form validation, domain relationships, HTMX responses, public rendering, and not-found handling.

## Deployment

The production image uses a multi-stage build:

```text
Node.js builder

Go builder
    ├── Builds the web server
    └── Builds adminctl

Debian runtime
    ├── Application binaries
    ├── Templates and static assets
    └── SQL migrations
```

The application is deployed to Fly.io. Production SQLite data is mounted separately from the application image so posts and workspace records survive new deployments.

A push to `main` triggers the Fly deployment workflow.

Recommended production configuration:

```text
APP_ENV=production
APP_PORT=8080
DATABASE_PATH=/data/portfolio.db
TEMPLATE_DIR=web/templates
LOG_LEVEL=info
```

## Design approach

The current experience is built around a **Clarity / Control** direction:

- Deep navy, cobalt, polo-blue, and restrained icy-cyan accents
- Modern sans-serif typography with precise monospaced details
- Distinct public, authentication, and personal-workspace shells
- CSS-first Tailwind v4 tokens
- Responsive layouts and generous spacing
- Small, purposeful motion
- Progressive enhancement rather than JavaScript dependence
- Accessible interaction states and reduced-motion support

The detailed design rationale and phased improvement plan are documented in:

```text
docs/site-experience-plan.md
```

## Status and roadmap

This is an actively maintained personal product rather than a tutorial repository.

Current and planned work is tracked through
[GitHub Issues](https://github.com/danieljmanningdev/danieljmanningdev-portfolio/issues).

Likely areas of continued refinement include:

- Bringing every client, project, and contract screen to the same visual standard
- Improving active workspace navigation and contextual wayfinding
- Expanding structured metadata and social sharing
- Tightening the Content Security Policy
- Additional browser-level, accessibility, and performance testing
- Carefully chosen HTMX enhancements where they improve an existing server-rendered workflow

## Licence

Released under the [MIT Licence](LICENSE).

## Author

Designed and built by **Daniel J. Manning**.

- Website: [danieljmanningdev.com](https://danieljmanningdev.com/)
- GitHub: [@danieljmanningdev](https://github.com/danieljmanningdev)
- Email: [daniel@danieljmanningdev.com](mailto:daniel@danieljmanningdev.com)
