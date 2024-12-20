// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
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
    RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
)
SELECT inserted_post.id, inserted_post.created_at, inserted_post.updated_at, inserted_post.title, inserted_post.url, inserted_post.description, inserted_post.published_at, inserted_post.feed_id
FROM inserted_post
LEFT JOIN feeds ON inserted_post.feed_id = feeds.id
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	FeedID      uuid.UUID
}

type CreatePostRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i CreatePostRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsFromUser = `-- name: GetPostsFromUser :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at
LIMIT $2
`

type GetPostsFromUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsFromUser(ctx context.Context, arg GetPostsFromUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsFromUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
