-- +goose Up
CREATE TABLE feeds(
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name text UNIQUE NOT NULL,
    -- The name of the feed
    url text UNIQUE NOT NULL,
    -- The URL of the feed
    user_id uuid NOT NULL,
    -- the ID of the user who added the feed
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE feeds;