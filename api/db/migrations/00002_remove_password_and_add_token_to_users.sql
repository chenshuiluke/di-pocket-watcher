-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN token text NOT NULL;

ALTER TABLE users
DROP COLUMN password;

-- +goose StatementEnd