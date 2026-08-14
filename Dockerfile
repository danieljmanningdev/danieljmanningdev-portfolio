# Build Tailwind CSS
FROM node:22-bookworm-slim AS css-builder

WORKDIR /src

COPY package.json package-lock.json ./
RUN npm ci

COPY web ./web

RUN npm run build:css


# Build Go application
FROM golang:1.26.5-bookworm AS go-builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 go build -v -o /run-app ./cmd/server
RUN CGO_ENABLED=1 go build -v -o /adminctl ./cmd/adminctl

# Production image
FROM debian:bookworm

WORKDIR /app

COPY --from=go-builder /run-app /run-app
COPY --from=go-builder /adminctl /adminctl

# Copy web files INCLUDING the Tailwind-generated app.css
COPY --from=css-builder /src/web /app/web

COPY --from=go-builder /usr/src/app/migrations /app/migrations

CMD ["/run-app"]