-- +goose Up
CREATE TABLE
    Tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        detail TEXT,
        current_status VARCHAR(50) DEFAULT 'pending'
    );

-- +goose Down
DROP TABLE Tasks;