FROM golang:1.26.5-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 go build -v -o /run-app ./cmd/server


FROM debian:bookworm

WORKDIR /app

COPY --from=builder /run-app /run-app

COPY web /app/web
COPY migrations /app/migrations

CMD ["/run-app"]