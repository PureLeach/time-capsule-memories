FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Копируем только go.mod и go.sum, чтобы использовать кэширование
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Устанавливаем goose для миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем весь исходный код
COPY . .

# Компилируем проект
RUN go build -o time-capsule-memories ./cmd

# Используем легковесный образ для финального контейнера
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk add --no-cache postgresql-client


# Копируем собранное приложение из builder
COPY --from=builder /app/time-capsule-memories /app/time-capsule-memories

# Копируем goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose


# Копируем entrypoint
COPY migrations /app/migrations
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh


# Устанавливаем рабочую директорию
WORKDIR /app

# Экспонируем порт
EXPOSE 8080


ENTRYPOINT ["/app/entrypoint.sh"]
