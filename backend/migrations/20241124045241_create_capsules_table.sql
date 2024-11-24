-- +goose Up
-- +goose StatementBegin
CREATE TABLE capsules (
    id SERIAL PRIMARY KEY,
    sender_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    send_at TIMESTAMP NOT NULL,
    message VARCHAR(4096) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    recipient_tg_username VARCHAR(50),
    files_folder_UUID UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE capsules;
-- +goose StatementEnd
