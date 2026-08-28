# Daniel J. Manning — Portfolio & Client Workspace

A fast, server-rendered developer portfolio and authenticated internal client workspace built with Go, HTMX, Tailwind CSS and SQLite.

The public website presents my UI/UX design and full-stack development work. Behind it, the application includes a private administrative workspace for managing clients, projects and contracts without duplicating broader operational tools such as scheduling, team communication or time tracking.

## What this demonstrates

• Production-style Go application architecture
• Authentication and secure server-side sessions
• SQLite persistence and migrations
• HTMX server-rendered interfaces
• Automated testing
• Security-focused middleware

## Screenshots

### Public portfolio

![Daniel J. Manning's portfolio home screen](web/static/images/djmdev-homescreen.png)

### Internal workspace

![Daniel J. Manning's internal workspace](web/static/images/internal-workspace.png)

## Current Features

### Public portfolio

- Responsive portfolio homepage
- Selected-work section
- About and capabilities sections
- Contact call to action
- Shared layout, header and footer templates
- Tailwind CSS visual system
- Server-rendered pages using Go templates

### Internal workspace

- Authenticated administrative dashboard
- Admin login and logout
- Server-side session persistence
- Session idle and absolute expiry
- HttpOnly session cookies
- Secure cookies in production
- Protected `/dashboard/*` routes
- No-store caching policy for authenticated pages
- Cross-origin request protection
- Central HTTP security headers
- Production-only HSTS
- Hardened HTTP server timeouts and header limits
- Client listing and detail pages
- Create and edit clients
- Delete clients with HTMX confirmation
- Active and inactive client statuses
- Project listing and detail pages
- Create and edit projects
- Archive projects
- Associate projects with clients
- Planned, active, completed and archived project statuses
- Optional project start and due dates
- Contract listing and detail pages
- Create and edit contracts
- Cancel contracts while preserving the record
- Associate contracts with clients
- Optional project association for contracts
- Contract-to-project ownership validation
- Draft, sent, accepted, completed and cancelled contract statuses
- Contract start and end dates
- Contract values stored as integer currency amounts
- Contract notes
- Server-side form validation
- SQLite persistence with foreign-key enforcement
- Automatic database migrations
- Structured application logging
- HTTP request logging with response status and duration
- Automated HTTP, auth, database and repository tests
- HTMX-aware redirects

## Technology

- **Backend:** Go and the standard-library `net/http` package
- **Templates:** Go `html/template`
- **Interactivity:** HTMX
- **Styling:** Tailwind CSS v4
- **Database:** SQLite using `modernc.org/sqlite`
- **Authentication:** bcrypt passwords and server-side sessions
- **Logging:** Go `log/slog`
- **Architecture:** HTTP handler → repository → SQLite
- **Testing:** Go `testing`, `httptest` and in-memory SQLite databases

The application deliberately avoids a heavy frontend framework. Most pages are rendered on the server, while HTMX is used for focused interactions where it provides a clear benefit.

## Project Structure

```text
.
├── cmd
│   ├── adminctl
│   │   └── main.go
│   └── server
│       ├── main.go
│       └── main_test.go
├── internal
│   ├── auth
│   ├── config
│   ├── database
│   ├── http
│   ├── logging
│   ├── models
│   └── repository
├── migrations
│   ├── 001_initial.sql
│   └── 002_auth.sql
├── web
│   ├── static
│   │   ├── css
│   │   │   └── input.css
│   │   └── images
│   │       ├── djmdev-homescreen.png
│   │       └── internal-workspace.png
│   └── templates
│       ├── components
│       ├── layouts
│       └── pages
├── go.mod
├── go.sum
├── package.json
├── package-lock.json
├── LICENSE
└── README.md
```

## Application Architecture

The application is intentionally organised into small, understandable layers.

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
HTTP handler
      ↓
Repository
      ↓
SQLite
      ↓
Go model
      ↓
HTML template
      ↓
Browser response
```

Each layer has a specific responsibility:

- `cmd/server` opens the database, constructs handlers, registers routes and starts the HTTP server.
- `cmd/adminctl` provides an explicit local command for creating administrative accounts.
- `internal/auth` contains password hashing, token generation and session-service logic.
- `internal/http` handles routing, HTTP methods, forms, validation, responses, authentication middleware, security middleware and template rendering.
- `internal/repository` contains parameterised SQL queries and persistence logic.
- `internal/models` defines the Go data structures used by the application.
- `internal/database` opens SQLite, enables foreign keys and runs migrations.
- `internal/config` loads application configuration from environment variables.
- `internal/logging` configures structured application logging.
- `web/templates` contains page, layout and shared component templates.
- `web/static` contains Tailwind source files and repository screenshots.

The core commercial relationship is:

```text
Client
  ↓
Project
  ↓
Contract
```

A contract always belongs to a client and may optionally belong to a project.

When a project is selected for a contract, the application validates that the project belongs to the same client as the contract.

This keeps the application focused on client and commercial records while operational workflows such as scheduling and time tracking remain outside the application.

## Authentication Flow

```text
GET /dashboard/
      ↓
RequireAdmin middleware
      ↓
No valid session
      ↓
303 /login
```

Successful login:

```text
POST /login
      ↓
Admin lookup
      ↓
bcrypt password verification
      ↓
cryptographically random session token
      ↓
SHA-256 token hash stored in SQLite
      ↓
HttpOnly session cookie
      ↓
303 /dashboard/
```

Session behaviour:

- Raw session tokens are not stored in SQLite.
- Only SHA-256 hashes of session tokens are persisted.
- Sessions have a 24-hour absolute lifetime.
- Sessions expire after 30 minutes of inactivity.
- Activity timestamps are periodically refreshed.
- Sessions are revoked when an admin becomes inactive.
- Logout revokes the server-side session and clears the browser cookie.
- Production cookies use the `Secure` flag.
- Authenticated dashboard responses use `Cache-Control: no-store`.

## Security

The application currently includes:

- bcrypt password hashing
- Generic login failure messages
- Server-side session persistence
- Hashed session tokens
- Session idle expiry
- Session absolute expiry
- Admin deactivation checks
- HttpOnly cookies
- SameSite cookies
- Secure production cookies
- Authorization middleware for all `/dashboard/*` routes
- Go `CrossOriginProtection`
- `Cache-Control: no-store` for authenticated workspace responses
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `Content-Security-Policy: frame-ancestors 'none'`
- `Referrer-Policy: strict-origin-when-cross-origin`
- Restricted `Permissions-Policy`
- Production-only `Strict-Transport-Security`
- Request-header, read, write and idle server limits
- Contract client/project ownership validation

The CSP is intentionally limited to framing protection for now. A broader CSP should be introduced after the public site's external asset requirements are reviewed.

## Local Development

### Requirements

- Go 1.26.5 or newer
- Node.js
- npm
- Git
- SQLite command-line tools are optional but useful for inspecting the development database

### Clone the repository

```bash
git clone https://github.com/danieljmanningdev/danieljmanningdev-portfolio.git
cd danieljmanningdev-portfolio
```

### Install frontend tooling

```bash
npm ci
```

### Build Tailwind CSS

```bash
npm run build:css
```

The generated stylesheet is written to:

```text
web/static/css/app.css
```

The generated CSS file is ignored by Git and should be rebuilt after cloning the repository or before deployment.

### Watch Tailwind during development

Run this in a separate terminal:

```bash
npm run dev:css
```

### Start the Go server

```bash
go run ./cmd/server
```

The development server starts at:

```text
http://localhost:8080
```

## Create the First Admin

There is no public registration route.

Administrative accounts are created explicitly using `adminctl`:

```bash
go run ./cmd/adminctl \
  -email "admin@example.com" \
  -name "Daniel Manning"
```

The command prompts for the password without echoing it to the terminal.

```text
Password:
Confirm password:
Admin created successfully
```

Only the bcrypt password hash is stored in SQLite.

## Useful Routes

| Route | Purpose | Access |
|---|---|---|
| `/` | Public portfolio | Public |
| `/health` | JSON health check | Public |
| `/login` | Admin login | Public |
| `/logout` | Revoke admin session | Authenticated POST |
| `/dashboard/` | Internal workspace overview | Admin |
| `/dashboard/clients` | Client listing | Admin |
| `/dashboard/clients/new` | Create-client form | Admin |
| `/dashboard/clients/{id}` | Client detail page | Admin |
| `/dashboard/clients/{id}/edit` | Edit-client form | Admin |
| `/dashboard/projects` | Project listing | Admin |
| `/dashboard/projects/new` | Create-project form | Admin |
| `/dashboard/projects/{id}` | Project detail page | Admin |
| `/dashboard/projects/{id}/edit` | Edit-project form | Admin |
| `/dashboard/projects/{id}/archive` | Archive project | Admin |
| `/dashboard/contracts` | Contract listing | Admin |
| `/dashboard/contracts/new` | Create-contract form | Admin |
| `/dashboard/contracts/{id}` | Contract detail page | Admin |
| `/dashboard/contracts/{id}/edit` | Edit-contract form | Admin |
| `/dashboard/contracts/{id}/cancel` | Cancel contract | Admin |

## Configuration

The application reads configuration from environment variables and falls back to development defaults when values are not provided.

| Variable | Default | Purpose |
|---|---|---|
| `APP_ENV` | `development` | Application environment |
| `APP_PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `./data/app.db` | SQLite database path |
| `TEMPLATE_DIR` | `web/templates` | Go template directory |
| `LOG_LEVEL` | `info` | Structured log level |

Example:

```bash
APP_ENV=development \
APP_PORT=9090 \
DATABASE_PATH=./data/development.db \
TEMPLATE_DIR=web/templates \
LOG_LEVEL=debug \
go run ./cmd/server
```

Local environment files and database files are excluded from Git.

## Database

SQLite is opened through `modernc.org/sqlite`.

The database layer:

- Creates the database directory when required
- Enables SQLite foreign-key enforcement
- Verifies that foreign-key enforcement is active
- Uses a single connection for in-memory test databases
- Runs versioned SQL migrations when the application starts
- Tracks applied migrations in `schema_migrations`

Current migrations:

```text
001_initial.sql
    clients
    projects
    contracts
    earlier foundation tables

002_auth.sql
    admins
    admin_sessions
```

The currently implemented application domains are:

- Clients
- Projects
- Contracts
- Admin authentication
- Admin sessions

Contracts require a client relationship. Their project relationship is optional, allowing both project-specific contracts and broader client agreements.

Application-level validation also ensures a selected project belongs to the contract's selected client.

### Inspect the development database

```bash
sqlite3 data/app.db
```

Useful SQLite commands:

```sql
.tables

.schema clients
.schema projects
.schema contracts
.schema admins
.schema admin_sessions

SELECT * FROM clients;
SELECT * FROM projects;
SELECT * FROM contracts;
SELECT id, email, display_name, is_active FROM admins;
SELECT id, admin_id, expires_at, last_seen_at FROM admin_sessions;
SELECT * FROM schema_migrations;
```

Never expose password hashes or session-token hashes unnecessarily when inspecting or logging authentication data.

## Client Management

The client feature supports:

```text
Create
Read
Update
Delete
```

Client forms validate:

- Required name
- Required email
- Valid email format
- Allowed client statuses
- Maximum field lengths
- Whitespace trimming

Valid statuses:

```text
active
inactive
```

## Project Management

Projects belong to clients.

The project feature supports:

```text
Create
Read
Update
Archive
```

Project forms validate:

- Required client
- Required project name
- Maximum name and description lengths
- Allowed project statuses
- Optional start and due dates
- `YYYY-MM-DD` date format
- Due dates cannot be earlier than start dates
- Whitespace trimming

Valid statuses:

```text
planned
active
completed
archived
```

## Contract Management

Contracts represent commercial agreements with clients.

Every contract belongs to a client and may optionally be associated with a project.

The contract feature supports:

```text
Create
Read
Update
Cancel
```

Contracts are cancelled rather than deleted so commercial records remain available.

Contract forms validate:

- Required client
- Required contract title
- Maximum title length
- Maximum notes length
- Allowed contract statuses
- Optional project association
- Selected project ownership
- Optional start and end dates
- `YYYY-MM-DD` date format
- End date cannot be earlier than start date
- Non-negative contract values
- Whitespace trimming

Valid statuses:

```text
draft
sent
accepted
completed
cancelled
```

## Logging

The application uses structured logging through Go's `log/slog`.

Development uses human-readable text logs, while production can use JSON output.

HTTP request middleware records:

- Request method
- Request path
- Response status
- Request duration

Example:

```text
time=2026-08-12T13:49:25+01:00 level=INFO msg="http request" method=POST path=/login status=303 duration=269ms
```

Authentication credentials, raw session tokens and password hashes should never be logged.

## Testing

### Run all tests

```bash
go test ./...
```

### Run tests verbosely

```bash
go test -v ./...
```

### Run package-specific tests

```bash
go test ./internal/auth
go test ./internal/database
go test ./internal/repository
go test ./internal/http
go test ./cmd/server
```

### Run static analysis

```bash
go vet ./...
```

### Check formatting

```bash
gofmt -w cmd internal
```

### Check for whitespace errors

```bash
git diff --check
```

### Full local verification

```bash
gofmt -w cmd internal
go test ./...
go vet ./...
npm run build:css
git diff --check
```

The test suite covers areas including:

- Configuration defaults and environment overrides
- Log-level configuration
- Database opening and connectivity
- Automatic database-directory creation
- Migration ordering and idempotency
- Failed migration handling
- SQLite foreign-key enforcement
- Client creation, updates and deletion
- Client validation
- HTMX deletion redirects
- Project repository create/read/update/archive flows
- Project form parsing and validation
- Project date handling
- Contract repository create/read/update/cancel flows
- Contract listing and client filtering
- Contract project/client ownership
- Missing contract clients and projects
- Contract form parsing and validation
- Contract date validation
- Contract value parsing
- Admin repository behaviour
- Admin session repository behaviour
- bcrypt password hashing and verification
- Session-token generation and hashing
- Absolute session expiry
- Idle session expiry
- Session activity updates
- Inactive-admin rejection
- Session revocation
- Login success and failure
- Secure-cookie behaviour
- Logout revocation
- Authenticated admin context
- Protected dashboard routing
- Security headers
- No-store cache behaviour
- Cross-origin protection
- HTTP router behaviour
- HTTP request logging
- Response-status capture
- Health endpoint behaviour
- Homepage rendering

## Roadmap

### Completed foundation

- [x] Public portfolio foundation
- [x] Client management
- [x] Project management
- [x] Contract management
- [x] Client → project relationships
- [x] Client → contract relationships
- [x] Optional project → contract relationships
- [x] Contract project ownership validation
- [x] SQLite persistence
- [x] Repository layer
- [x] Server-side validation
- [x] Structured application logging
- [x] HTTP request middleware
- [x] Admin authentication
- [x] Secure server-side sessions
- [x] HttpOnly session cookies
- [x] Secure production cookies
- [x] Authorization middleware
- [x] Cross-origin protection
- [x] Security headers
- [x] Authenticated-page no-store policy
- [x] HTTP server timeouts and header limits
- [x] Automated tests

### Next launch work

- [ ] Decide whether the client portal is required for initial launch
- [ ] Public portfolio content polish
- [ ] Case-study pages
- [ ] Real contact and social links
- [ ] SEO and Open Graph metadata
- [ ] Custom 404 page
- [ ] Accessibility review
- [ ] Responsive/browser QA
- [ ] Production deployment
- [ ] Persistent SQLite storage
- [ ] Automated database backups
- [ ] Backup restoration test
- [ ] Production smoke test

### Future workspace improvements

- [ ] Client search
- [ ] Client status filtering
- [ ] Prefer client archiving over permanent deletion
- [ ] Improve project filtering and sorting
- [ ] Improve contract filtering and sorting
- [ ] Add selected workspace summaries to the dashboard
- [ ] Review and remove unused legacy schema tables

### Future commercial workflow

Potential future additions may include:

- Invitation-only client access
- Client deliverables
- Contract document generation
- Contract acceptance or electronic signing
- Invoicing
- Payment status
- Payment processing

These features should extend the existing client → project → contract relationship rather than duplicate operational tools used for scheduling, tasks or time tracking.

## Design Philosophy

This project is built around a small set of principles:

- Prefer server-rendered HTML when a JavaScript application is unnecessary.
- Use HTMX for focused interactions rather than introducing a large frontend framework.
- Keep SQL and persistence logic out of HTTP handlers.
- Keep HTTP concerns out of repositories.
- Use small, understandable layers.
- Use structured logging for observable application behaviour.
- Build fast interfaces without unnecessary frontend bloat.
- Keep the internal workspace focused on clients, projects and commercial agreements.
- Avoid rebuilding scheduling, team communication and time-tracking tools that already exist elsewhere.
- Preserve important commercial records rather than deleting them unnecessarily.
- Treat authentication and authorization as server-side concerns.
- Store session-token hashes rather than raw session credentials.
- Add abstractions only when repeated application behaviour justifies them.
- Keep the public portfolio polished while building the internal workspace incrementally.

## Deployment Status

The application now includes the core authentication and HTTP security controls required to protect the internal administrative workspace.

Production launch still requires:

- persistent SQLite storage
- backup automation
- backup restoration testing
- HTTPS deployment
- production environment configuration
- final security verification
- public-site QA

The internal workspace should not be considered production-ready until those deployment controls are also verified.

## Licence

This project is licensed under the MIT Licence.

Copyright © 2026 Daniel J. Manning.
