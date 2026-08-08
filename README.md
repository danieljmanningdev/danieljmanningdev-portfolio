# Go + HTMX Micro-SaaS Boilerplate

A zero-bloat, high-performance Go backend template featuring a clean architecture, local HTMX integration, SQLite database storage, built-in Authentication (`bcrypt`), and Stripe billing hooks. Designed for rapid micro-app development and deployment.

---

## 🚀 Features

* **Backend:** Go (Standard Library `net/http` multiplexer)
* **Frontend:** Server-side rendered HTML templates powered locally by **HTMX** (zero Node.js or Tailwind compiler overhead)
* **Database:** Embedded **SQLite** with automated migration setups
* **Security & Auth:** Secure user registration & login flows using `bcrypt` password hashing
* **Billing:** Pre-configured Stripe checkout session wrappers and webhook-ready billing routes

---

## 🛠️ Project Structure

```text
.
├── backend/
│   ├── db/            # Database connection & migration runners
│   ├── handlers/      # HTTP request handlers (Auth, Billing, etc.)
│   ├── middleware/    # Auth and request middleware
│   ├── migrations/    # SQL database migration files
│   ├── models/        # Data structs (User, etc.)
│   ├── routes/        # Router configuration and path multiplexing
│   ├── services/      # Third-party integrations (Stripe)
│   ├── static/        # Local static assets (HTMX, CSS, JS)
│   └── templates/     # HTML templates (index.html)
├── .gitignore
├── LICENSE
├── Makefile
├── main.go            # Application entrypoint
└── setup.sh           # Automated project initialization script
```

---

## ⚡ Getting Started (New Project Setup)

When creating a new repository from this template, follow these steps to initialize your project:

### 1. Remove the Old Module File (If Present)
If the template repository contains an older `go.mod` file from testing, make sure to delete it before running the setup script so it doesn't conflict with your new module path:
```bash
rm -f go.mod go.sum
```

### 2. Grant Permissions to the Setup Script
Make the automated setup script executable in your terminal:
```bash
chmod +x setup.sh
```

### 3. Run the Setup Script
Execute the script to initialize your Go module, create the full folder hierarchy, download a local copy of HTMX, and install all essential dependencies (SQLite, Bcrypt, Stripe):
```bash
./setup.sh
```
*(You will be prompted to enter your module path, e.g., `github.com/yourusername/your-new-app`)*

### 4. Run the Application
Start your development server:
```bash
go run main.go
```
Your app will be live at `http://localhost:8080`.

---

## 💳 Environment Variables

To use Stripe checkout features, ensure you set your environment variables before running the application:

```bash
export STRIPE_SECRET_KEY="your_stripe_secret_key"
export STRIPE_PRICE_ID="your_stripe_price_id"
```
