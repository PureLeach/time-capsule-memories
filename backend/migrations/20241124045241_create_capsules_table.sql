-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_type AS ENUM ('waiting', 'in progress', 'done');

CREATE TABLE capsules (
    id SERIAL PRIMARY KEY,
    sender_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    send_at TIMESTAMP NOT NULL,
    message VARCHAR(4096) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    files_folder_UUID UUID,
    status status_type NOT NULL DEFAULT 'waiting'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS capsules;
DROP TYPE IF EXISTS status_type;
-- +goose StatementEnd
