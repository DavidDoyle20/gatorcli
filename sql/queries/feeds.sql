-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedUsers :many
SELECT 
    feeds.*,
    users.id AS user_id,
    users.name AS user_name
FROM feeds 
LEFT JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * 
FROM feeds
WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    updated_at = CURRENT_TIMESTAMP,
    last_fetched_at = CURRENT_TIMESTAMP
WHERE id = $1
;

-- name: GetNextFeedToFetch :one
SELECT * 
FROM feeds
ORDER BY last_fetched_at NULLS FIRST, last_fetched_at ASC
LIMIT 1
;