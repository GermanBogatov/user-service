-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id           UUID NOT NULL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    surname      VARCHAR(255) NOT NULL,
    email        VARCHAR(255) UNIQUE NOT NULL,
    password     VARCHAR(255) NOT NULL,
    roles        VARCHAR[],
    created_date  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_date  TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email
    ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_email;
DROP TABLE users;
-- +goose StatementEnd
