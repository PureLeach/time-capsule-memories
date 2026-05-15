-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_capsules_due
    ON capsules (send_at)
    WHERE status = 'waiting';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_capsules_due;
-- +goose StatementEnd
