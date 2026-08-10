# Daniel J. Manning — Portfolio & Client Workspace

A fast, server-rendered developer portfolio and internal client workspace built with Go, HTMX, Tailwind CSS and SQLite.

The public website presents my UI/UX design and full-stack development work. Behind it, the application includes an internal workspace for managing clients and projects without duplicating broader operational tools such as scheduling, team communication or time tracking.

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

- Dashboard overview
- Client listing
- Client detail pages
- Create clients
- Edit clients
- Delete clients with HTMX confirmation
- Active and inactive client statuses
- Project listing
- Project detail pages
- Create projects
- Edit projects
- Archive projects
- Associate projects with clients
- Planned, active, completed and archived project statuses
- Optional project start and due dates
- Server-side form validation
- SQLite persistence with foreign-key enforcement
- Automatic database migrations
- Structured application logging
- HTTP request logging with response status and duration
- Automated HTTP, database and repository tests
- HTMX-aware redirects

## Technology

- **Backend:** Go and the standard-library `net/http` package
- **Templates:** Go `html/template`
- **Interactivity:** HTMX
- **Styling:** Tailwind CSS v4
- **Database:** SQLite using `modernc.org/sqlite`
- **Logging:** Go `log/slog`
- **Architecture:** HTTP handler → repository → SQLite
- **Testing:** Go `testing`, `httptest` and in-memory SQLite databases

The application deliberately avoids a heavy frontend framework. Most pages are rendered on the server, while HTMX is used for focused interactions where it provides a clear benefit.

## Project Structure

```text
.
├── cmd
│   └── server
│       ├── main.go
│       └── main_test.go
├── internal
│   ├── config
│   ├── database
│   ├── http
│   ├── logging
│   ├── models
│   └── repository
├── migrations
│   └── 001_initial.sql
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
HTTP middleware
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

- `cmd/server` opens the database, constructs handlers, registers routes and starts the server.
- `internal/http` handles routing, HTTP methods, forms, validation, responses, middleware and template rendering.
- `internal/repository` contains parameterised SQL queries and persistence logic.
- `internal/models` defines the Go data structures used by the application.
- `internal/database` opens SQLite, enables foreign keys and runs migrations.
- `internal/config` loads application configuration from environment variables.
- `internal/logging` configures structured application logging.
- `web/templates` contains page, layout and shared component templates.
- `web/static` contains Tailwind source files and repository screenshots.

### Example client-update flow

```text
POST /dashboard/clients/5
      ↓
ClientsHandler.handlePOST
      ↓
ClientsHandler.updateClient
      ↓
ClientRepository.Update
      ↓
UPDATE clients
      ↓
303 redirect to /dashboard/clients/5
```

### Example project-create flow

```text
POST /dashboard/projects/new
      ↓
ProjectsHandler.handleProjectPOST
      ↓
ProjectsHandler.createProject
      ↓
ProjectRepository.Create
      ↓
INSERT INTO projects
      ↓
303 redirect to /dashboard/projects/{id}
```

### Example request-logging flow

```text
HTTP request
      ↓
RequestLogger middleware
      ↓
Application handler
      ↓
Response status captured
      ↓
Structured log entry
```

Example development log:

```text
level=INFO msg="http request" method=GET path=/dashboard/projects status=200 duration=705µs
```

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

Tailwind will rebuild the stylesheet when template classes change.

### Start the Go server

```bash
go run ./cmd/server
```

The development server starts at:

```text
http://localhost:8080
```

## Useful Routes

| Route | Purpose |
|---|---|
| `/` | Public portfolio |
| `/health` | JSON health check |
| `/dashboard/` | Internal workspace overview |
| `/dashboard/clients` | Client listing |
| `/dashboard/clients/new` | Create-client form |
| `/dashboard/clients/{id}` | Client detail page |
| `/dashboard/clients/{id}/edit` | Edit-client form |
| `/dashboard/projects` | Project listing |
| `/dashboard/projects/new` | Create-project form |
| `/dashboard/projects/{id}` | Project detail page |
| `/dashboard/projects/{id}/edit` | Edit-project form |
| `/dashboard/projects/{id}/archive` | Archive project |

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

The currently implemented application domains are:

- Clients
- Projects

The initial schema also contains some earlier foundation tables that may be removed or deferred as the workspace remains intentionally focused.

Relationships between records are protected using SQLite foreign keys.

### Inspect the development database

```bash
sqlite3 data/app.db
```

Useful SQLite commands:

```sql
.tables
.schema clients
.schema projects
SELECT * FROM clients;
SELECT * FROM projects;
SELECT * FROM schema_migrations;
```

Exit SQLite with:

```text
.quit
```

## Client Management

The client feature supports the basic CRUD lifecycle:

```text
Create
Read
Update
Delete
```

### Repository operations

```go
List()
GetByID()
Create()
Update()
Delete()
```

### HTTP operations

```text
GET     /dashboard/clients
GET     /dashboard/clients/new
POST    /dashboard/clients/new
GET     /dashboard/clients/{id}
GET     /dashboard/clients/{id}/edit
POST    /dashboard/clients/{id}
DELETE  /dashboard/clients/{id}
```

### Current validation

Client forms validate:

- Required name
- Required email
- Valid email format
- Allowed client statuses
- Whitespace trimming

Valid statuses are:

```text
active
inactive
```

## Project Management

Projects belong to clients and are intentionally focused on project records rather than full operational task management.

The project feature supports:

```text
Create
Read
Update
Archive
```

### Repository operations

```go
List()
ListByClientID()
GetByID()
Create()
Update()
Archive()
```

### HTTP operations

```text
GET   /dashboard/projects
GET   /dashboard/projects/new
POST  /dashboard/projects/new
GET   /dashboard/projects/{id}
GET   /dashboard/projects/{id}/edit
POST  /dashboard/projects/{id}
POST  /dashboard/projects/{id}/archive
```

### Current validation

Project forms validate:

- Required client
- Required project name
- Maximum name and description lengths
- Allowed project statuses
- Optional start and due dates
- `YYYY-MM-DD` date format
- Due dates cannot be earlier than start dates
- Whitespace trimming

Valid statuses are:

```text
planned
active
completed
archived
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
time=2026-08-10T12:40:49+01:00 level=INFO msg="http request" method=GET path=/ status=200 duration=596µs
```

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
go test ./internal/database
go test ./internal/repository
go test ./internal/http
```

### Run static analysis

```bash
go vet ./...
```

### Check formatting

```bash
gofmt -w .
```

### Check for whitespace errors

```bash
git diff --check
```

### Full local verification

```bash
gofmt -w .
go test ./...
go vet ./...
npm run build:css
git diff --check
```

The test suite currently covers areas including:

- Configuration defaults and environment overrides
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
- HTTP router behaviour
- HTTP request logging
- Response-status capture
- Health endpoint behaviour
- Homepage rendering

## Roadmap

### Completed foundation

- [x] Public portfolio
- [x] Client management
- [x] Project management
- [x] SQLite persistence
- [x] Repository layer
- [x] Server-side validation
- [x] Structured application logging
- [x] HTTP request middleware
- [x] Automated tests

### Near-term improvements

- [ ] Client search
- [ ] Client status filtering
- [ ] Prefer client archiving over permanent deletion
- [ ] Improve project filtering and sorting
- [ ] Add selected project summaries to the dashboard
- [ ] Improve validation feedback and form polish
- [ ] Review and remove unused legacy schema tables

### Security before public dashboard deployment

- [ ] Admin authentication
- [ ] Secure server-side sessions
- [ ] HTTP-only secure cookies
- [ ] CSRF protection
- [ ] Authorization middleware
- [ ] Production security headers
- [ ] Appropriate production server timeouts
- [ ] Persistent deployment storage for SQLite

The internal dashboard is currently intended for local/private use until those security controls are implemented.

## Design Philosophy

This project is built around a small set of principles:

- Prefer server-rendered HTML when a JavaScript application is unnecessary.
- Use HTMX for focused interactions rather than introducing a large frontend framework.
- Keep SQL and persistence logic out of HTTP handlers.
- Keep HTTP concerns out of repositories.
- Use small, understandable layers.
- Use structured logging for observable application behaviour.
- Build fast interfaces without unnecessary frontend bloat.
- Keep the internal workspace focused on client and project records.
- Avoid rebuilding scheduling, team communication and time-tracking tools that already exist elsewhere.
- Add abstractions only when repeated application behaviour justifies them.
- Keep the public portfolio polished while building the internal workspace incrementally.

## Deployment Status

The public portfolio can be prepared for deployment, but the internal dashboard should remain private until authentication and related security controls are implemented.

## Licence

This project is licensed under the MIT Licence.

Copyright © 2026 Daniel J. Manning.