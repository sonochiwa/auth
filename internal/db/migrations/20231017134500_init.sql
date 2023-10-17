-- +goose Up
CREATE TABLE IF NOT EXISTS passwords (
    user_guid VARCHAR(36) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS passwords;