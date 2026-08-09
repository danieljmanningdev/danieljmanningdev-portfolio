# Daniel J. Manning — Portfolio & Client Workspace

A fast, server-rendered developer portfolio and internal client-management workspace built with Go, HTMX, Tailwind CSS and SQLite.

The public website presents my UI/UX design and full-stack development work. Behind it, the application includes an internal workspace for managing clients and is being expanded to support projects, contracts, messaging, time tracking and an invitation-only customer portal.

## Current Features

### Public portfolio

- Responsive portfolio homepage
- Selected-work section
- About and capabilities sections
- Contact call to action
- Shared layout, header and footer templates
- Tailwind CSS visual system

### Internal workspace

- Dashboard overview
- Client listing
- Client detail pages
- Create clients
- Edit clients
- Delete clients with HTMX confirmation
- Active and inactive client statuses
- Server-side form validation
- SQLite persistence with foreign-key enforcement
- Automated HTTP, database and repository tests

## Technology

- **Backend:** Go and the standard-library `net/http` package
- **Templates:** Go `html/template`
- **Interactivity:** HTMX
- **Styling:** Tailwind CSS v4
- **Database:** SQLite using `modernc.org/sqlite`
- **Architecture:** HTTP handler → repository → SQLite
- **Testing:** Go `testing` and `httptest`

The application deliberately avoids a heavy frontend framework. Most pages are rendered on the server, while HTMX is used for focused interactions such as deleting clients and redirecting back to the client-management page.

## Project Structure

```text
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── config
│   ├── database
│   ├── http
│   ├── models
│   └── repository
├── migrations
│   └── 001_initial.sql
├── web
│   ├── static
│   │   └── css
│   │       └── input.css
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