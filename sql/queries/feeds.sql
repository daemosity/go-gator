-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE feeds.url = $1;
-- name: DeleteAllFeeds :exec
DELETE FROM feeds
WHERE 1 = 1;
-- name: ListAllFeeds :many
SELECT feeds.name as "feedName",
    feeds.url as "feedURL",
    users.name as "createdBy"
FROM feeds
    LEFT JOIN users ON users.id = feeds.user_id;