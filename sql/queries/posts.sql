-- name: CreatePost :one
WITH inserted_post AS (
    INSERT INTO posts (id, title, url, description, published_at, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
    RETURNING *
)
SELECT inserted_post.*
FROM inserted_post
LEFT JOIN feeds ON inserted_post.feed_id = feeds.id
;

-- name: GetPostsFromUser :many
SELECT posts.*
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at
LIMIT $2
;