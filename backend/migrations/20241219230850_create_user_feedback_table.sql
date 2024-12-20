-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_feedback (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message VARCHAR(4096) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_feedback;
-- +goose StatementEnd
