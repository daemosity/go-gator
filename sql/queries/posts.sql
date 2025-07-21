-- name: CreatePost :one
INSERT INTO posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: GetPostsForUser :many
WITH followed_feeds AS (
    SELECT feed_follows.feed_id AS feed_id
    FROM feed_follows
    WHERE feed_follows.user_id = $1
)
SELECT *
FROM posts
WHERE posts.feed_id IN (
        SELECT feed_id
        FROM followed_feeds
    )
ORDER BY posts.published_at DESC
LIMIT $2;
-- name: GetSortedOrFilteredPostsForUser :many
SELECT p.*
FROM posts p
    INNER JOIN feed_follows ff ON p.feed_id = ff.feed_id
    INNER JOIN feeds f ON p.feed_id = f.id
WHERE ff.user_id = $1
    AND (
        -- optional. If the feed_name parameter is an empty string, the first part of the 'OR' is true and it shows posts from ALL followed feeds
        -- If feed_name is provided, it filters for posts where f.name matches
        sqlc.arg(feed_name)::text = ''
        OR f.name = @feed_name
    )
ORDER BY -- This helps with dynamic sorting by passing boolean flags rather than injecting ASC or DESC
    CASE
        WHEN sqlc.arg(sort_asc)::bool THEN p.published_at
    END ASC,
    CASE
        WHEN sqlc.arg(sort_desc)::bool THEN p.published_at
    END DESC
LIMIT $2;