-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id           UUID NOT NULL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    surname      VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    roles        VARCHAR[],
    createdDate  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updatedDate  TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_email
    ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_email;
DROP TABLE users;
-- +goose StatementEnd
