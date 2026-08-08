#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "What is your GitHub username/organization and new repo name?"
read -p "Enter module path (e.g., github.com/username/project-name): " MODULE_PATH

if [ -z "$MODULE_PATH" ]; then
  echo "Error: Module path cannot be empty."
  exit 1
fi

echo "Initializing Go module..."
go mod init "$MODULE_PATH"

echo "Creating project directory structure..."
mkdir -p backend/middleware backend/migrations backend/models backend/db backend/templates backend/routes backend/handlers backend/services backend/static/js

echo "Downloading local copy of HTMX..."
curl -s https://unpkg.com/htmx.org@2.0.10/dist/htmx.min.js -o backend/static/js/htmx.min.js

echo "Installing essential dependencies (SQLite, Bcrypt, Stripe)..."
go get github.com/mattn/go-sqlite3
go get golang.org/x/crypto/bcrypt
go get github.com/stripe/stripe-go/v76

echo "Tidying up modules..."
go mod tidy

echo "Resetting Git repository for your new project..."
rm -rf .git
git init
git add .
git commit -m "Initial commit from SaaS boilerplate"

echo "Setup complete! Your micro-SaaS boilerplate is ready. Run 'make run' to start."