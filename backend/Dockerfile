FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy module manifests first to leverage build cache
COPY go.mod go.sum ./
RUN go mod download

# Install goose for running migrations on startup
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -o time-capsule-memories ./cmd


FROM alpine:latest

# Required by entrypoint for pg_isready
RUN apk add --no-cache postgresql-client

COPY --from=builder /app/time-capsule-memories /app/time-capsule-memories
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /app/migrations
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

WORKDIR /app
EXPOSE 8000

ENTRYPOINT ["/app/entrypoint.sh"]
