-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET
  updated_at = $1,
  last_fetched_at = $1
WHERE feeds.id = $2
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY
  last_fetched_at ASC NULLS FIRST,
  updated_at ASC
LIMIT 1;

-- name: DeleteAllFeeds :exec
-- DELETE FROM feeds;

