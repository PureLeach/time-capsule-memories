FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Manifests first, so the dependency layer survives source changes.
COPY go.mod go.sum ./
RUN go mod download

# Applies migrations on startup. Pinned: newer goose needs Go 1.25+, so bump it
# together with the builder image.
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o time-capsule-memories ./cmd


FROM alpine:3.21

# pg_isready for entrypoint.sh; busybox wget backs the healthcheck.
RUN apk add --no-cache postgresql-client && \
    adduser -D -u 1000 app

WORKDIR /app

COPY --from=builder --chown=app:app /app/time-capsule-memories ./time-capsule-memories
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --chown=app:app migrations ./migrations
COPY --chown=app:app --chmod=755 entrypoint.sh ./entrypoint.sh

USER app

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8000/healthz || exit 1

ENTRYPOINT ["./entrypoint.sh"]
