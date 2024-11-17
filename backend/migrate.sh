#!/bin/bash

# Проверка, существует ли файл .env
if [ -f .env ]; then
    # Чтение переменных окружения из .env и игнорирование комментариев
    export $(grep -v '^#' .env | xargs)
fi

# Сборка строки подключения
DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
# DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# Выполнение команды миграции
migrate -path migrations -database "$DATABASE_URL" up
