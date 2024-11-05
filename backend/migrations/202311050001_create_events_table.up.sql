-- migrations/202311050001_create_events_table.up.sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY, -- Поле для ID события
    sender_name VARCHAR(100) NOT NULL,
    message_date TIMESTAMP NOT NULL,
    message TEXT NOT NULL,
    open_date TIMESTAMP NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    telegram_nick VARCHAR(100),
    photos TEXT -- Поле для хранения пути к фото или URL-адресов
);
