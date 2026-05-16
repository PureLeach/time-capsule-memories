FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy module manifests first to leverage build cache
COPY go.mod go.sum ./
RUN go mod download

# Install goose for running migrations on startup. Pinned: newer goose
# releases require Go 1.25+; bump together with the builder image.
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

COPY . .
RUN go build -o time-capsule-memories ./cmd


FROM alpine:latest

# postgresql-client gives entrypoint.sh pg_isready
RUN apk add --no-cache postgresql-client && \
    adduser -D -u 1000 app

WORKDIR /app

COPY --from=builder --chown=app:app /app/time-capsule-memories /app/time-capsule-memories
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --chown=app:app migrations /app/migrations
COPY --chown=app:app entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

USER app

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8000/healthz || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
