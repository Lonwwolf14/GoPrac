-- +goose Up
CREATE TABLE
    Tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        detail TEXT,
        current_status BOOLEAN NOT NULL DEFAULT FALSE
    );

-- +goose Down
DROP TABLE Tasks;