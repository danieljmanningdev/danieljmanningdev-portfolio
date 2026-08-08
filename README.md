# Daniel Manning Portfolio & Micro-SaaS Engine

A high-performance, minimalist personal portfolio and client portal built using Go, HTMX, and Tailwind CSS. Designed to showcase custom software engineering services with zero heavy JavaScript frameworks and zero unnecessary cloud bloat.

---

## Architecture & Tech Stack

* Backend: Go (Standard Library net/http router)
* Frontend Interactivity: HTMX (hx-boost, reactive partial swaps)
* Styling: Tailwind CSS v4 (CDN / Utility-first)
* Database: SQLite (app.db) managed via native Go drivers and custom migrations
* Content Engine: Flat-file Markdown parser (goldmark) with frontmatter extraction for engineering notes/blogs
* Payments/Billing: Stripe integration stubbed for project scoping and deposits

---

## Project Structure

```text
.
├── app.db                # SQLite database
├── backend
│   ├── db                # Database connection & migration scripts
│   ├── handlers          # HTTP handlers (Auth, Billing, Blog)
│   ├── middleware        # Request middleware & auth guards
│   ├── models            # Data structs (User, etc.)
│   ├── routes            # Route definitions and multiplexer wiring
│   ├── services          # Third-party integrations (Stripe)
│   ├── static            # Static assets (HTMX, scripts)
│   └── templates         # HTML templates (Index, Blog, Login, Portal)
├── content               # Markdown blog posts & articles
├── main.go               # Application entry point
├── Makefile              # Build and run commands
├── setup.sh              # Environment bootstrap script
└── README.md
```

---

## Getting Setup & Running

### Prerequisites
* Go 1.22+ installed locally.

### Local Development Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/danieljmanningdev/danieljmanningdev-portfolio.git
   cd danieljmanningdev-portfolio
   ```

2. Run the server:
   ```bash
   make run
   ```
   (Or run via Go directly: go run main.go)

3. Open in your browser:
   Navigate to http://localhost:8080 to view the portfolio, interactive project estimator, and markdown blog engine.