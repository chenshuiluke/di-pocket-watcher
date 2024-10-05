-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
  users (
    id BIGSERIAL PRIMARY KEY,
    email text NOT NULL,
    password text
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
