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

RUN CGO_ENABLED=1 go build -trimpath -v -o /run-app ./cmd/server
RUN CGO_ENABLED=1 go build -trimpath -v -o /adminctl ./cmd/adminctl
RUN CGO_ENABLED=1 go build -trimpath -v -o /dbctl ./cmd/dbctl


# Production image
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates gosu \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 10001 app \
    && useradd \
        --uid 10001 \
        --gid app \
        --home-dir /app \
        --shell /usr/sbin/nologin \
        app

WORKDIR /app

COPY --from=go-builder /run-app /run-app
COPY --from=go-builder /adminctl /adminctl
COPY --from=go-builder /dbctl /dbctl

# Copy web files including the Tailwind-generated app.css.
COPY --from=css-builder /src/web /app/web
COPY --from=go-builder /usr/src/app/migrations /app/migrations
COPY --chmod=0755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint

RUN mkdir -p /data \
    && chown -R app:app /app /data

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
CMD ["/run-app"]
