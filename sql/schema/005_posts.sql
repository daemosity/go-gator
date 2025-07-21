-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text UNIQUE,
    url text UNIQUE NOT NULL,
    description text,
    published_at timestamp NOT NULL,
    feed_id uuid NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE posts;